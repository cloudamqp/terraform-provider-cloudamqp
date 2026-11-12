package cloudamqp

import (
	"fmt"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// teamMemberAcceptWindow is how long TestAccTeamMember_AcceptedMember sleeps
// between inviting and expecting "member" status, giving you time to accept
// the invite via the web UI before the next step's Read runs. Only applies
// when CLOUDAMQP_RECORD is set (recording), never during replay. Raise this
// and re-record if it isn't enough time in practice.
var teamMemberAcceptWindow = 60 * time.Second

// TestAccTeamMember_PendingInvite: invite a member, change the pending
// invite's tags (which replaces the invite, a pending invite can't be
// updated), import it, then let the implicit destroy revoke it, the invite is
// never accepted in this test. Exercises POST /api/team/invite and
// DELETE /api/team/invites/:id. Does NOT run in parallel with the other
// TestAccTeamMember_* tests, they share the one test@example.com mailbox and
// must run one at a time.
func TestAccTeamMember_PendingInvite(t *testing.T) {
	memberResourceName := "cloudamqp_team_member.this"
	var firstInviteID string

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: `
					resource "cloudamqp_team_member" "this" {
						email = "test@example.com"
						tags  = ["tf-acc-test-member-tag-a"]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(memberResourceName, "email", "test@example.com"),
					resource.TestCheckResourceAttr(memberResourceName, "id", "test@example.com"),
					resource.TestCheckResourceAttr(memberResourceName, "role", "member"),
					resource.TestCheckResourceAttr(memberResourceName, "status", "pending"),
					resource.TestCheckResourceAttrSet(memberResourceName, "invite_id"),
					resource.TestCheckNoResourceAttr(memberResourceName, "user_id"),
					resource.TestCheckResourceAttr(memberResourceName, "tags.#", "1"),
					resource.TestCheckTypeSetElemAttr(memberResourceName, "tags.*", "tf-acc-test-member-tag-a"),
					resource.TestCheckResourceAttrWith(memberResourceName, "invite_id", func(value string) error {
						firstInviteID = value
						return nil
					}),
				),
			},
			{
				Config: `
					resource "cloudamqp_team_member" "this" {
						email = "test@example.com"
						tags  = ["tf-acc-test-member-tag-a", "tf-acc-test-member-tag-b"]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(memberResourceName, "status", "pending"),
					resource.TestCheckResourceAttr(memberResourceName, "tags.#", "2"),
					resource.TestCheckResourceAttrWith(memberResourceName, "invite_id", func(value string) error {
						if value == firstInviteID {
							return fmt.Errorf("expected the pending invite to be replaced, invite_id is still %s", value)
						}
						return nil
					}),
				),
			},
			{
				ResourceName:      memberResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccTeamMember_AcceptedMember: invite a member, manually accept it (see
// teamMemberAcceptWindow above), confirm the resource picks up "member"
// status, update the member's tags, import it, then let the implicit destroy
// remove the member. Exercises POST /invite, GET /api/team?email=, PUT
// /api/team/:user_id and POST /api/team/remove. Does NOT run in parallel with
// the other TestAccTeamMember_* tests, see teamMemberAcceptWindow.
func TestAccTeamMember_AcceptedMember(t *testing.T) {
	memberResourceName := "cloudamqp_team_member.this"

	pendingConfig := `
		resource "cloudamqp_team_member" "this" {
			email = "test@example.com"
			role  = "member"
			tags  = ["tf-acc-test-member-tag-a"]
		}
	`

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: pendingConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(memberResourceName, "status", "pending"),
					resource.TestCheckResourceAttrSet(memberResourceName, "invite_id"),
				),
			},
			{
				// Same config as the previous step, this step exists only to force a
				// fresh Read after the invite has (hopefully) been accepted.
				PreConfig: func() {
					if os.Getenv("CLOUDAMQP_RECORD") != "" {
						t.Logf("Accept the invite sent to test@example.com via the web UI now, waiting %s...",
							teamMemberAcceptWindow)
						time.Sleep(teamMemberAcceptWindow)
					}
				},
				Config: pendingConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(memberResourceName, "status", "member"),
					resource.TestCheckResourceAttrSet(memberResourceName, "user_id"),
					resource.TestCheckNoResourceAttr(memberResourceName, "invite_id"),
					resource.TestCheckResourceAttr(memberResourceName, "tags.#", "1"),
				),
			},
			{
				Config: `
					resource "cloudamqp_team_member" "this" {
						email = "test@example.com"
						role  = "member"
						tags  = ["tf-acc-test-member-tag-a", "tf-acc-test-member-tag-b"]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(memberResourceName, "status", "member"),
					resource.TestCheckResourceAttr(memberResourceName, "tags.#", "2"),
				),
			},
			{
				Config: `
					resource "cloudamqp_team_member" "this" {
						email = "test@example.com"
						role  = "member"
						tags  = ["tf-acc-test-member-tag-a"]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(memberResourceName, "status", "member"),
					resource.TestCheckResourceAttr(memberResourceName, "tags.#", "1"),
				),
			},
			{
				Config: `
					resource "cloudamqp_team_member" "this" {
						email = "test@example.com"
						role  = "billing manager"
						tags  = []
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(memberResourceName, "status", "member"),
					resource.TestCheckResourceAttr(memberResourceName, "role", "billing manager"),
					resource.TestCheckResourceAttr(memberResourceName, "tags.#", "0"),
				),
			},
			{
				ResourceName:      memberResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccTeamMember_InvalidEmail: inviting a malformed email is rejected.
func TestAccTeamMember_InvalidEmail(t *testing.T) {
	t.Parallel()

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: `
					resource "cloudamqp_team_member" "member" {
						email = "not-an-email"
						tags  = ["tf-acc-test-invalid-email"]
					}
				`,
				ExpectError: regexp.MustCompile(`(?i)invalid email`),
			},
		},
	})
}

// TestAccTeamMember_TagCase: the second invite uses a tag that only differs in letter case from
// the tag the first invite created. The backend assigns the existing tag, and the configured
// spelling is kept in state, so there's no diff. Both invites stay pending, nothing to accept.
func TestAccTeamMember_TagCase(t *testing.T) {
	t.Parallel()

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: `
					resource "cloudamqp_team_member" "first" {
						email = "tag-case-first@example.com"
						tags  = ["TF-Acc-Test-Case"]
					}

					resource "cloudamqp_team_member" "second" {
						email      = "tag-case-second@example.com"
						tags       = ["tf-acc-test-case"]
						depends_on = [cloudamqp_team_member.first]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemAttr("cloudamqp_team_member.first", "tags.*", "TF-Acc-Test-Case"),
					resource.TestCheckTypeSetElemAttr("cloudamqp_team_member.second", "tags.*", "tf-acc-test-case"),
				),
			},
		},
	})
}

// TestAccTeamMember_TaggedRoleRequiresTags: the "member" and "monitor" roles
// require at least one tag, whether the role is set explicitly or left to
// its default.
func TestAccTeamMember_TaggedRoleRequiresTags(t *testing.T) {
	t.Parallel()

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: `
					resource "cloudamqp_team_member" "this" {
						email = "test@example.com"
						role  = "monitor"
						tags  = []
					}
				`,
				ExpectError: regexp.MustCompile(`(?i)"monitor" role requires at least one tag`),
			},
			{
				Config: `
					resource "cloudamqp_team_member" "this" {
						email = "test@example.com"
						tags  = []
					}
				`,
				ExpectError: regexp.MustCompile(`(?i)"member" role requires at least one tag`),
			},
		},
	})
}

// TestAccTeamMember_UntaggedRoleRejectsTags: roles other than "member" and
// "monitor" can't have tags, the backend clears them on role change.
func TestAccTeamMember_UntaggedRoleRejectsTags(t *testing.T) {
	t.Parallel()

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: `
					resource "cloudamqp_team_member" "this" {
						email = "test@example.com"
						role  = "admin"
						tags  = ["tf-acc-test-member-tag-a"]
					}
				`,
				ExpectError: regexp.MustCompile(`(?i)"admin" role can't have tags`),
			},
		},
	})
}
