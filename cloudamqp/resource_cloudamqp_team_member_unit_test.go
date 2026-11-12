package cloudamqp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Unit tests for cloudamqp_team_member against a fake CloudAMQP API. The fake API mimics the
// customer-console team endpoints (response shapes taken from the TestAccTeamMember_* cassettes)
// and records every request, so the tests can assert which calls the provider makes in each
// state: pending invite, accepted member, and removed outside Terraform.

const (
	teamMemberTestEmail    = "test@example.com"
	teamMemberTestUserID   = "8849c6b5-a975-4a3e-8991-f31cc49d5784"
	teamMemberResourceName = "cloudamqp_team_member.this"
)

type fakeTeamEntry struct {
	id    string
	email string
	role  string
	tags  []string
}

type fakeTeamAPI struct {
	mu           sync.Mutex
	nextInviteID int
	invites      []*fakeTeamEntry
	members      []*fakeTeamEntry
	teamTags     []string
	reads        []string // GET requests
	writes       []string // POST, PUT and DELETE requests
	lastPutBody  string

	// Simulates the invite being accepted right before the provider revokes it.
	acceptOnRevoke bool
}

func newFakeTeamAPI(t *testing.T) (*fakeTeamAPI, *httptest.Server) {
	fake := &fakeTeamAPI{nextInviteID: 1}
	server := httptest.NewServer(http.HandlerFunc(fake.handle))
	t.Cleanup(server.Close)
	return fake, server
}

func (f *fakeTeamAPI) handle(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()

	request := r.Method + " " + r.URL.RequestURI()
	if r.Method == http.MethodGet {
		f.reads = append(f.reads, request)
	} else {
		f.writes = append(f.writes, request)
	}
	w.Header().Set("Content-Type", "application/json")

	path := r.URL.Path
	email := r.URL.Query().Get("email")
	switch {
	case r.Method == http.MethodPost && path == "/api/team/invites":
		f.invite(w, r)
	case r.Method == http.MethodPost && path == "/api/team/remove":
		f.remove(w, r)
	case r.Method == http.MethodGet && path == "/api/team/invites":
		result := []map[string]any{}
		for _, invite := range f.invites {
			if email == "" || invite.email == email {
				result = append(result, inviteJSON(invite))
			}
		}
		writeJSON(w, http.StatusOK, result)
	case strings.HasPrefix(path, "/api/team/invites/"):
		f.handleInvite(w, r, strings.TrimPrefix(path, "/api/team/invites/"))
	case r.Method == http.MethodGet && path == "/api/team":
		result := []map[string]any{}
		for _, member := range f.members {
			if email == "" || member.email == email {
				body := memberJSON(member)
				body["roles"] = append([]string{member.role}, member.tags...)
				result = append(result, body)
			}
		}
		writeJSON(w, http.StatusOK, result)
	case strings.HasPrefix(path, "/api/team/"):
		f.handleMember(w, r, strings.TrimPrefix(path, "/api/team/"))
	default:
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "not found"})
	}
}

func (f *fakeTeamAPI) invite(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email string   `json:"email"`
		Role  string   `json:"role"`
		Tags  []string `json:"tags"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	if f.find(f.members, body.Email) != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "User is already in this team"})
		return
	}
	if f.find(f.invites, body.Email) != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "User is already invited"})
		return
	}
	role := body.Role
	if role == "" {
		role = "member"
	}
	invite := &fakeTeamEntry{
		id:    strconv.Itoa(f.nextInviteID),
		email: body.Email,
		role:  role,
		tags:  f.teamTagsFor(body.Tags),
	}
	f.nextInviteID++
	f.invites = append(f.invites, invite)
	writeJSON(w, http.StatusCreated, inviteJSON(invite))
}

func (f *fakeTeamAPI) handleInvite(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method == http.MethodDelete && f.acceptOnRevoke {
		f.acceptLocked(teamMemberTestEmail)
	}
	index := slices.IndexFunc(f.invites, func(invite *fakeTeamEntry) bool { return invite.id == id })
	if index < 0 {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": fmt.Sprintf("Pending invite %s not found.", id)})
		return
	}
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, http.StatusOK, inviteJSON(f.invites[index]))
	case http.MethodDelete:
		f.invites = slices.Delete(f.invites, index, index+1)
		writeJSON(w, http.StatusOK, map[string]any{"message": "Invite revoked"})
	default:
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "not found"})
	}
}

func (f *fakeTeamAPI) handleMember(w http.ResponseWriter, r *http.Request, id string) {
	index := slices.IndexFunc(f.members, func(member *fakeTeamEntry) bool { return member.id == id })
	if index < 0 {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "User not found"})
		return
	}
	member := f.members[index]
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, http.StatusOK, memberJSON(member))
	case http.MethodPut:
		f.updateMember(w, r, member)
	default:
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "not found"})
	}
}

// updateMember mimics PUT /api/team/:user_id: the role is set first, and switching to a role
// other than member/monitor clears all tags. An empty tags list leaves the tags unchanged.
func (f *fakeTeamAPI) updateMember(w http.ResponseWriter, r *http.Request, member *fakeTeamEntry) {
	var body map[string]any
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	raw, _ := json.Marshal(body)
	f.lastPutBody = string(raw)

	role, hasRole := body["role"].(string)
	var tags []string
	if list, ok := body["tags"].([]any); ok {
		for _, tag := range list {
			if name, _ := tag.(string); name != "" {
				tags = append(tags, name)
			}
		}
	}
	if !hasRole && len(tags) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "Nothing to update."})
		return
	}
	if hasRole {
		if !slices.Contains(assignableTeamRoleNames, role) {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "Role does not exist."})
			return
		}
		if !slices.Contains(taggedTeamRoleNames, role) {
			member.tags = nil
		}
		member.role = role
	}
	if len(tags) > 0 {
		member.tags = f.teamTagsFor(tags)
	}
	writeJSON(w, http.StatusOK, map[string]any{"message": "Role and tags updated for user."})
}

// teamTagsFor mimics how the backend matches tags to the team's tags: an exact match first, then
// one that only differs in letter case. A tag that matches none is created.
func (f *fakeTeamAPI) teamTagsFor(tags []string) []string {
	result := []string{}
	for _, tag := range tags {
		index := slices.Index(f.teamTags, tag)
		if index < 0 {
			index = slices.IndexFunc(f.teamTags, func(t string) bool { return strings.EqualFold(t, tag) })
		}
		if index < 0 {
			f.teamTags = append(f.teamTags, tag)
			index = len(f.teamTags) - 1
		}
		if !slices.Contains(result, f.teamTags[index]) {
			result = append(result, f.teamTags[index])
		}
	}
	return result
}

func (f *fakeTeamAPI) remove(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	index := slices.IndexFunc(f.members, func(member *fakeTeamEntry) bool { return member.email == body.Email })
	if index < 0 {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "No such user in team"})
		return
	}
	f.members = slices.Delete(f.members, index, index+1)
	writeJSON(w, http.StatusOK, map[string]any{"message": "The user has been deleted from team"})
}

func (f *fakeTeamAPI) find(entries []*fakeTeamEntry, email string) *fakeTeamEntry {
	for _, entry := range entries {
		if entry.email == email {
			return entry
		}
	}
	return nil
}

func inviteJSON(invite *fakeTeamEntry) map[string]any {
	id, _ := strconv.Atoi(invite.id)
	return map[string]any{"id": id, "email": invite.email, "role": invite.role,
		"tags": nonNilTags(invite.tags), "created": "2026-10-16 09:30:59 UTC"}
}

func memberJSON(member *fakeTeamEntry) map[string]any {
	return map[string]any{"id": member.id, "email": member.email, "tfa_auth_enabled": false,
		"role": member.role, "tags": nonNilTags(member.tags)}
}

// nonNilTags returns the tags sorted by name, like the backend does
func nonNilTags(tags []string) []string {
	sorted := slices.Clone(tags)
	if sorted == nil {
		sorted = []string{}
	}
	slices.Sort(sorted)
	return sorted
}

// accept turns the pending invite for email into an accepted member, like the invitee
// accepting it in the web UI.
func (f *fakeTeamAPI) accept(email string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.acceptLocked(email)
}

func (f *fakeTeamAPI) acceptLocked(email string) {
	index := slices.IndexFunc(f.invites, func(invite *fakeTeamEntry) bool { return invite.email == email })
	if index < 0 {
		return
	}
	invite := f.invites[index]
	f.invites = slices.Delete(f.invites, index, index+1)
	f.members = append(f.members, &fakeTeamEntry{
		id: teamMemberTestUserID, email: invite.email, role: invite.role, tags: invite.tags,
	})
}

// removeOutsideTerraform removes the invite and member for email, like a team admin doing it
// in the web UI.
func (f *fakeTeamAPI) removeOutsideTerraform(email string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	hasEmail := func(entry *fakeTeamEntry) bool { return entry.email == email }
	f.invites = slices.DeleteFunc(f.invites, hasEmail)
	f.members = slices.DeleteFunc(f.members, hasEmail)
}

func (f *fakeTeamAPI) resetReads() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.reads = nil
}

// requireOnlyReads fails the test unless every GET since the last reset is the expected one.
func (f *fakeTeamAPI) requireOnlyReads(t *testing.T, expected string) {
	t.Helper()
	f.mu.Lock()
	defer f.mu.Unlock()
	if len(f.reads) == 0 {
		t.Fatalf("expected GET requests to %s, got none", expected)
	}
	for _, read := range f.reads {
		if read != expected {
			t.Fatalf("expected only %s, got %v", expected, f.reads)
		}
	}
}

// checkWrites asserts the POST, PUT and DELETE requests since the last check, in order.
func (f *fakeTeamAPI) checkWrites(expected ...string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		f.mu.Lock()
		defer f.mu.Unlock()
		writes := f.writes
		f.writes = nil
		if len(expected) == 0 && len(writes) == 0 {
			return nil
		}
		if !reflect.DeepEqual(expected, writes) {
			return fmt.Errorf("unexpected requests\nwant: %v\n got: %v", expected, writes)
		}
		return nil
	}
}

func (f *fakeTeamAPI) checkLastPutBody(expected string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		f.mu.Lock()
		defer f.mu.Unlock()
		var want, got any
		if err := json.Unmarshal([]byte(expected), &want); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(f.lastPutBody), &got); err != nil {
			return err
		}
		if !reflect.DeepEqual(want, got) {
			return fmt.Errorf("unexpected PUT body\nwant: %s\n got: %s", expected, f.lastPutBody)
		}
		return nil
	}
}

func teamMemberUnitTest(t *testing.T, steps []resource.TestStep) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: muxProviderFactories(http.DefaultClient),
		Steps:                    steps,
	})
}

func teamMemberUnitTestConfig(baseURL, attributes string) string {
	return fmt.Sprintf(`
provider "cloudamqp" {
  apikey  = "test"
  baseurl = %q
}

resource "cloudamqp_team_member" "this" {
  email = %q
  %s
}
`, baseURL, teamMemberTestEmail, attributes)
}

// TestTeamMember_TagOrder: tags are a set, so tags in another order than the backend's sorted
// order give no diff, neither for a pending invite nor for a member.
func TestTeamMember_TagOrder(t *testing.T) {
	t.Parallel()
	fake, server := newFakeTeamAPI(t)
	config := teamMemberUnitTestConfig(server.URL, `tags = ["b", "a"]`)

	teamMemberUnitTest(t, []resource.TestStep{
		{
			Config: config,
			Check: resource.ComposeAggregateTestCheckFunc(
				fake.checkWrites("POST /api/team/invites"),
				resource.TestCheckTypeSetElemAttr(teamMemberResourceName, "tags.*", "a"),
				resource.TestCheckTypeSetElemAttr(teamMemberResourceName, "tags.*", "b"),
			),
		},
		{
			PreConfig: func() { fake.accept(teamMemberTestEmail) },
			Config:    config,
			Check:     resource.TestCheckResourceAttr(teamMemberResourceName, "status", "member"),
		},
		{
			Config:   config,
			PlanOnly: true,
		},
	})
}

// TestTeamMember_TagCase: a tag that only differs in letter case from a team tag is assigned the
// team's tag by the backend. The configured spelling is kept in state, so there's no diff, both
// for a pending invite and after it's accepted.
func TestTeamMember_TagCase(t *testing.T) {
	t.Parallel()
	fake, server := newFakeTeamAPI(t)
	fake.teamTags = []string{"Production"}

	teamMemberUnitTest(t, []resource.TestStep{
		{
			Config: teamMemberUnitTestConfig(server.URL, `tags = ["production"]`),
			Check: resource.ComposeAggregateTestCheckFunc(
				fake.checkWrites("POST /api/team/invites"),
				resource.TestCheckTypeSetElemAttr(teamMemberResourceName, "tags.*", "production"),
				func(_ *terraform.State) error {
					if tags := fake.invites[0].tags; !slices.Equal(tags, []string{"Production"}) {
						return fmt.Errorf("expected the team's tag Production to be assigned, got %v", tags)
					}
					return nil
				},
			),
		},
		{
			PreConfig: func() { fake.accept(teamMemberTestEmail) },
			Config:    teamMemberUnitTestConfig(server.URL, `tags = ["production"]`),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(teamMemberResourceName, "status", "member"),
				resource.TestCheckTypeSetElemAttr(teamMemberResourceName, "tags.*", "production"),
			),
		},
		{
			Config: teamMemberUnitTestConfig(server.URL, `tags = ["production", "staging"]`),
			Check: resource.ComposeAggregateTestCheckFunc(
				fake.checkWrites("PUT /api/team/"+teamMemberTestUserID),
				fake.checkLastPutBody(`{"role":"member","tags":["production","staging"]}`),
				resource.TestCheckTypeSetElemAttr(teamMemberResourceName, "tags.*", "production"),
				resource.TestCheckTypeSetElemAttr(teamMemberResourceName, "tags.*", "staging"),
			),
		},
	})
}

// TestTeamMember_PendingInviteTagCaseChange: changing only the letter case of a tag on a pending
// invite updates the state, it doesn't revoke and resend the invite.
func TestTeamMember_PendingInviteTagCaseChange(t *testing.T) {
	t.Parallel()
	fake, server := newFakeTeamAPI(t)
	fake.teamTags = []string{"Production"}

	teamMemberUnitTest(t, []resource.TestStep{
		{
			Config: teamMemberUnitTestConfig(server.URL, `tags = ["production"]`),
			Check:  fake.checkWrites("POST /api/team/invites"),
		},
		{
			Config: teamMemberUnitTestConfig(server.URL, `tags = ["Production"]`),
			Check: resource.ComposeAggregateTestCheckFunc(
				fake.checkWrites(),
				resource.TestCheckResourceAttr(teamMemberResourceName, "invite_id", "1"),
				resource.TestCheckTypeSetElemAttr(teamMemberResourceName, "tags.*", "Production"),
			),
		},
	})
}

// TestTeamMember_PendingInviteReplace: changing tags or role on a pending invite revokes it and
// sends a new invite, it's never updated in place.
func TestTeamMember_PendingInviteReplace(t *testing.T) {
	t.Parallel()
	fake, server := newFakeTeamAPI(t)

	teamMemberUnitTest(t, []resource.TestStep{
		{
			Config: teamMemberUnitTestConfig(server.URL, `tags = ["a"]`),
			Check: resource.ComposeAggregateTestCheckFunc(
				fake.checkWrites("POST /api/team/invites"),
				resource.TestCheckResourceAttr(teamMemberResourceName, "status", "pending"),
				resource.TestCheckResourceAttr(teamMemberResourceName, "invite_id", "1"),
				resource.TestCheckResourceAttr(teamMemberResourceName, "role", "member"),
			),
		},
		{
			Config: teamMemberUnitTestConfig(server.URL, `tags = ["a", "b"]`),
			Check: resource.ComposeAggregateTestCheckFunc(
				fake.checkWrites("DELETE /api/team/invites/1", "POST /api/team/invites"),
				resource.TestCheckResourceAttr(teamMemberResourceName, "status", "pending"),
				resource.TestCheckResourceAttr(teamMemberResourceName, "invite_id", "2"),
				resource.TestCheckResourceAttr(teamMemberResourceName, "tags.#", "2"),
			),
		},
		{
			Config: teamMemberUnitTestConfig(server.URL, `
  role = "monitor"
  tags = ["a", "b"]
`),
			Check: resource.ComposeAggregateTestCheckFunc(
				fake.checkWrites("DELETE /api/team/invites/2", "POST /api/team/invites"),
				resource.TestCheckResourceAttr(teamMemberResourceName, "invite_id", "3"),
				resource.TestCheckResourceAttr(teamMemberResourceName, "role", "monitor"),
			),
		},
	})
}

// TestTeamMember_AcceptedInvite: once the invite is accepted, the next refresh picks up the
// member, and role/tag changes are updated in place.
func TestTeamMember_AcceptedInvite(t *testing.T) {
	t.Parallel()
	fake, server := newFakeTeamAPI(t)
	memberPath := "/api/team/" + teamMemberTestUserID

	teamMemberUnitTest(t, []resource.TestStep{
		{
			Config: teamMemberUnitTestConfig(server.URL, `tags = ["a"]`),
			Check:  fake.checkWrites("POST /api/team/invites"),
		},
		{
			PreConfig: func() { fake.accept(teamMemberTestEmail) },
			Config:    teamMemberUnitTestConfig(server.URL, `tags = ["a"]`),
			Check: resource.ComposeAggregateTestCheckFunc(
				fake.checkWrites(),
				resource.TestCheckResourceAttr(teamMemberResourceName, "status", "member"),
				resource.TestCheckResourceAttr(teamMemberResourceName, "user_id", teamMemberTestUserID),
				resource.TestCheckNoResourceAttr(teamMemberResourceName, "invite_id"),
			),
		},
		{
			Config: teamMemberUnitTestConfig(server.URL, `tags = ["a", "b"]`),
			Check: resource.ComposeAggregateTestCheckFunc(
				fake.checkWrites("PUT "+memberPath),
				fake.checkLastPutBody(`{"role": "member", "tags": ["a", "b"]}`),
				resource.TestCheckResourceAttr(teamMemberResourceName, "status", "member"),
				resource.TestCheckResourceAttr(teamMemberResourceName, "tags.#", "2"),
			),
		},
		{
			// The backend clears the tags when switching to a role other than member/monitor
			Config: teamMemberUnitTestConfig(server.URL, `
  role = "admin"
  tags = []
`),
			Check: resource.ComposeAggregateTestCheckFunc(
				fake.checkWrites("PUT "+memberPath),
				fake.checkLastPutBody(`{"role": "admin", "tags": []}`),
				resource.TestCheckResourceAttr(teamMemberResourceName, "role", "admin"),
				resource.TestCheckResourceAttr(teamMemberResourceName, "tags.#", "0"),
			),
		},
	})
}

// TestTeamMember_RefreshLooksUpByID: a refresh does a single lookup by id, the email search is
// only used when the invite has just been accepted.
func TestTeamMember_RefreshLooksUpByID(t *testing.T) {
	t.Parallel()
	fake, server := newFakeTeamAPI(t)
	config := teamMemberUnitTestConfig(server.URL, `tags = ["a"]`)

	teamMemberUnitTest(t, []resource.TestStep{
		{
			Config: config,
		},
		{
			PreConfig: fake.resetReads,
			Config:    config,
			PlanOnly:  true,
		},
		{
			PreConfig: func() {
				fake.requireOnlyReads(t, "GET /api/team/invites/1")
				fake.accept(teamMemberTestEmail)
			},
			Config: config,
			Check:  resource.TestCheckResourceAttr(teamMemberResourceName, "status", "member"),
		},
		{
			PreConfig: fake.resetReads,
			Config:    config,
			PlanOnly:  true,
		},
		{
			PreConfig: func() { fake.requireOnlyReads(t, "GET /api/team/"+teamMemberTestUserID) },
			Config:    config,
			PlanOnly:  true,
		},
	})
}

// TestTeamMember_RemovedOutsideTerraform: an invite or member removed outside Terraform is
// dropped from state, so the next plan creates it again instead of failing.
func TestTeamMember_RemovedOutsideTerraform(t *testing.T) {
	t.Parallel()

	for name, accepted := range map[string]bool{"pending invite": false, "member": true} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			fake, server := newFakeTeamAPI(t)
			config := teamMemberUnitTestConfig(server.URL, `tags = ["a"]`)

			teamMemberUnitTest(t, []resource.TestStep{
				{
					Config: config,
				},
				{
					PreConfig: func() {
						if accepted {
							fake.accept(teamMemberTestEmail)
						}
					},
					Config: config,
				},
				{
					PreConfig:          func() { fake.removeOutsideTerraform(teamMemberTestEmail) },
					Config:             config,
					PlanOnly:           true,
					ExpectNonEmptyPlan: true,
				},
			})
		})
	}
}

// TestTeamMember_AcceptedBetweenPlanAndApply: the invite is accepted after the plan to replace
// it, so the revoke finds no invite. The apply stops without sending a new invite and keeps the
// state, and the next apply updates the member in place.
func TestTeamMember_AcceptedBetweenPlanAndApply(t *testing.T) {
	t.Parallel()
	fake, server := newFakeTeamAPI(t)
	config := teamMemberUnitTestConfig(server.URL, `tags = ["a", "b"]`)

	teamMemberUnitTest(t, []resource.TestStep{
		{
			Config: teamMemberUnitTestConfig(server.URL, `tags = ["a"]`),
			Check:  fake.checkWrites("POST /api/team/invites"),
		},
		{
			PreConfig: func() {
				fake.mu.Lock()
				fake.acceptOnRevoke = true
				fake.mu.Unlock()
			},
			Config:      config,
			ExpectError: regexp.MustCompile(`(?s)invite for "test@example.com" was accepted after the plan`),
		},
		{
			Config: config,
			Check: resource.ComposeAggregateTestCheckFunc(
				fake.checkWrites("DELETE /api/team/invites/1", "PUT /api/team/"+teamMemberTestUserID),
				fake.checkLastPutBody(`{"role": "member", "tags": ["a", "b"]}`),
				resource.TestCheckResourceAttr(teamMemberResourceName, "status", "member"),
				resource.TestCheckResourceAttr(teamMemberResourceName, "tags.#", "2"),
			),
		},
	})
}
