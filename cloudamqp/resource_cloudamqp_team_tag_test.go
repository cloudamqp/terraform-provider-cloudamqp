package cloudamqp

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccTeamTag_Import: Test team tag creation and import. "id" is the tag's
// name (not its underlying numeric id - see "tag_id").
func TestAccTeamTag_Import(t *testing.T) {
	t.Parallel()

	tagResourceName := "cloudamqp_team_tag.tag"

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: `
					resource "cloudamqp_team_tag" "tag" {
						name = "tf-acc-test-tag-import"
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(tagResourceName, "name", "tf-acc-test-tag-import"),
					resource.TestCheckResourceAttr(tagResourceName, "id", "tf-acc-test-tag-import"),
					resource.TestCheckResourceAttrSet(tagResourceName, "tag_id"),
				),
			},
			{
				ResourceName:      tagResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccTeamTag_ReplacementOnNameChange: tags can't be renamed, so changing
// "name" must force a new resource (DELETE the old tag, POST a new one).
func TestAccTeamTag_ReplacementOnNameChange(t *testing.T) {
	t.Parallel()

	tagResourceName := "cloudamqp_team_tag.tag"

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: `
					resource "cloudamqp_team_tag" "tag" {
						name = "tf-acc-test-tag-replace-a"
					}
				`,
				Check: resource.TestCheckResourceAttr(tagResourceName, "name", "tf-acc-test-tag-replace-a"),
			},
			{
				Config: `
					resource "cloudamqp_team_tag" "tag" {
						name = "tf-acc-test-tag-replace-b"
					}
				`,
				Check: resource.TestCheckResourceAttr(tagResourceName, "name", "tf-acc-test-tag-replace-b"),
			},
		},
	})
}

// TestAccTeamTag_InvalidName: blank and over-length tag names are rejected at
// plan time by the "name" validator (mirrors Services::TeamTags::MAX_NAME_LENGTH
// and the BlankTagNameError check in customer-console).
func TestAccTeamTag_InvalidName(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		tag  string
	}{
		{"Blank", ""},
		{"TooLong", strings.Repeat("a", 501)},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cloudamqpResourceTest(t, resource.TestCase{
				PreCheck: func() { testAccPreCheck(t) },
				Steps: []resource.TestStep{
					{
						Config: fmt.Sprintf(`
							resource "cloudamqp_team_tag" "tag" {
								name = "%s"
							}
						`, c.tag),
						ExpectError: regexp.MustCompile(`string length must be between 1 and 500`),
					},
				},
			})
		})
	}
}
