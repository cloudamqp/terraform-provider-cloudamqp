package cloudamqp

import (
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

// TestAccTeamMember_PendingInvite: invite a member, update the pending
// invite's tags, import it, then let the implicit destroy revoke it, the
// invite is never accepted in this test. Exercises POST /api/team/invite and
// PUT/DELETE /api/team/invites/:id. Does NOT run in parallel with the other
// TestAccTeamMember_* tests, they share the one test@example.com mailbox and
// must run one at a time.
func TestAccTeamMember_PendingInvite(t *testing.T) {
	memberResourceName := "cloudamqp_team_member.member"

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: `
					resource "cloudamqp_team_tag" "member_tag_a" {
						name = "tf-acc-test-member-tag-a"
					}

					resource "cloudamqp_team_tag" "member_tag_b" {
						name       = "tf-acc-test-member-tag-b"
						depends_on = [cloudamqp_team_tag.member_tag_a]
					}

					resource "cloudamqp_team_member" "member" {
						email = "test@example.com"
						tags  = []
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(memberResourceName, "email", "test@example.com"),
					resource.TestCheckResourceAttr(memberResourceName, "id", "test@example.com"),
					resource.TestCheckResourceAttr(memberResourceName, "role", "member"),
					resource.TestCheckResourceAttr(memberResourceName, "status", "pending"),
					resource.TestCheckResourceAttrSet(memberResourceName, "invite_id"),
					resource.TestCheckNoResourceAttr(memberResourceName, "user_id"),
					resource.TestCheckResourceAttr(memberResourceName, "tags.#", "0"),
				),
			},
			{
				Config: `
					resource "cloudamqp_team_tag" "member_tag_a" {
						name = "tf-acc-test-member-tag-a"
					}

					resource "cloudamqp_team_tag" "member_tag_b" {
						name       = "tf-acc-test-member-tag-b"
						depends_on = [cloudamqp_team_tag.member_tag_a]
					}

					resource "cloudamqp_team_member" "member" {
						email = "test@example.com"
						tags  = [cloudamqp_team_tag.member_tag_a.name]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(memberResourceName, "status", "pending"),
					resource.TestCheckResourceAttr(memberResourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(memberResourceName, "tags.0", "tf-acc-test-member-tag-a"),
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
// remove the member. Exercises POST /invite, GET /api/team/member, PUT
// /api/team/:user_id and POST /api/team/remove. Does NOT run in parallel with
// the other TestAccTeamMember_* tests - see teamMemberAcceptWindow.
func TestAccTeamMember_AcceptedMember(t *testing.T) {
	memberResourceName := "cloudamqp_team_member.member"

	pendingConfig := `
		resource "cloudamqp_team_tag" "member_tag_a" {
			name = "tf-acc-test-member-tag-a"
		}

		resource "cloudamqp_team_tag" "member_tag_b" {
			name       = "tf-acc-test-member-tag-b"
			depends_on = [cloudamqp_team_tag.member_tag_a]
		}

		resource "cloudamqp_team_member" "member" {
			email = "test@example.com"
			role  = "member"
			tags  = [cloudamqp_team_tag.member_tag_a.name]
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
						t.Logf("Accept the invite sent to test@example.com via the web UI now - waiting %s...",
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
					resource "cloudamqp_team_tag" "member_tag_a" {
						name = "tf-acc-test-member-tag-a"
					}

					resource "cloudamqp_team_tag" "member_tag_b" {
						name       = "tf-acc-test-member-tag-b"
						depends_on = [cloudamqp_team_tag.member_tag_a]
					}

					resource "cloudamqp_team_member" "member" {
						email = "test@example.com"
						role  = "member"
						tags  = [cloudamqp_team_tag.member_tag_a.name, cloudamqp_team_tag.member_tag_b.name]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(memberResourceName, "status", "member"),
					resource.TestCheckResourceAttr(memberResourceName, "tags.#", "2"),
				),
			},
			{
				// Clearing an accepted member's tags to [] via PUT /api/team/:user_id.
				Config: `
					resource "cloudamqp_team_tag" "member_tag_a" {
						name = "tf-acc-test-member-tag-a"
					}

					resource "cloudamqp_team_tag" "member_tag_b" {
						name       = "tf-acc-test-member-tag-b"
						depends_on = [cloudamqp_team_tag.member_tag_a]
					}

					resource "cloudamqp_team_member" "member" {
						email = "test@example.com"
						role  = "member"
						tags  = []
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(memberResourceName, "status", "member"),
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
						tags  = []
					}
				`,
				ExpectError: regexp.MustCompile(`(?i)invalid email`),
			},
		},
	})
}
