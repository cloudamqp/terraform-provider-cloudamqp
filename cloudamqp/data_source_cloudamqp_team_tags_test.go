package cloudamqp

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccTeamTagsDataSource_Basic: lists team tags and confirms two
// newly-created tags are present in the result. Only checks for presence of
// its own tags, not an exact count, so it's safe to run alongside other
// parallel tests that also create/delete team tags.
func TestAccTeamTagsDataSource_Basic(t *testing.T) {
	t.Parallel()

	dataSourceName := "data.cloudamqp_team_tags.tags"

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: `
					resource "cloudamqp_team_tag" "ds_tag_a" {
						name = "tf-acc-test-ds-tag-a"
					}

					resource "cloudamqp_team_tag" "ds_tag_b" {
						name       = "tf-acc-test-ds-tag-b"
						depends_on = [cloudamqp_team_tag.ds_tag_a]
					}

					data "cloudamqp_team_tags" "tags" {
						depends_on = [cloudamqp_team_tag.ds_tag_a, cloudamqp_team_tag.ds_tag_b]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "tags.#"),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceName, "tags.*",
						map[string]string{"name": "tf-acc-test-ds-tag-a"}),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceName, "tags.*",
						map[string]string{"name": "tf-acc-test-ds-tag-b"}),
				),
			},
		},
	})
}
