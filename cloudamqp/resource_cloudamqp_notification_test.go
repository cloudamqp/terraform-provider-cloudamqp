package cloudamqp

import (
	"fmt"
	"os"
	"testing"

	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/configuration"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccNotification_Basic: Create CPU alarm, import and change values.
func TestAccNotification_Basic(t *testing.T) {
	t.Parallel()

	var (
		fileNames                = []string{"instance", "notification"}
		instanceResourceName     = "cloudamqp_instance.instance"
		notificationResourceName = "cloudamqp_notification.recipient"

		params = map[string]string{
			"InstanceName":   "TestAccNotification_Basic",
			"InstanceID":     fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":   "bunny-1",
			"RecipientType":  "email",
			"RecipientValue": "notification@example.com",
			"RecipientName":  "notification",
		}

		paramsUpdated = map[string]string{
			"InstanceName":   "TestAccNotification_Basic",
			"InstanceID":     fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":   "bunny-1",
			"RecipientType":  "email",
			"RecipientValue": "test@example.com",
			"RecipientName":  "test",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(notificationResourceName, "type", params["RecipientType"]),
					resource.TestCheckResourceAttr(notificationResourceName, "value", params["RecipientValue"]),
					resource.TestCheckResourceAttr(notificationResourceName, "name", params["RecipientName"]),
				),
			},
			{
				ResourceName:      notificationResourceName,
				ImportStateIdFunc: testAccImportCombinedStateIdFunc(instanceResourceName, notificationResourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, paramsUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(notificationResourceName, "type", paramsUpdated["RecipientType"]),
					resource.TestCheckResourceAttr(notificationResourceName, "value", paramsUpdated["RecipientValue"]),
					resource.TestCheckResourceAttr(notificationResourceName, "name", paramsUpdated["RecipientName"]),
				),
			},
		},
	})
}

func TestAccNotification_Email(t *testing.T) {
	t.Parallel()

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: `
          resource "cloudamqp_notification" "email_recipient" {
            instance_id = 1091
            type        = "email"
            value       = "alarm@example.com"
            name        = "alarm"
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cloudamqp_notification.email_recipient", "type", "email"),
					resource.TestCheckResourceAttr("cloudamqp_notification.email_recipient", "value", "alarm@example.com"),
					resource.TestCheckResourceAttr("cloudamqp_notification.email_recipient", "name", "alarm"),
				),
			},
			{
				ResourceName:      "cloudamqp_notification.email_recipient",
				ImportStateIdFunc: testAccImportCombinedIdFunc("1091", "cloudamqp_notification.email_recipient"),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNotification_Opsgenie(t *testing.T) {
	t.Parallel()

	// Set sanitized value for playback and use real value for recording
	testRecipientOpsgenieValue := "RECIPIENT_OPSGENIE_VALUE"
	if os.Getenv("CLOUDAMQP_RECORD") != "" {
		testRecipientOpsgenieValue = os.Getenv("RECIPIENT_OPSGENIE_VALUE")
	}

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
          resource "cloudamqp_notification" "opsgenie_recipient" {
            instance_id = 1091
            type        = "opsgenie"
            value       = "%s"
            name        = "OpsGenie"
            responders {
              type = "team"
              id   = "209402c4-a92f-4b91-a1e8-d48a8ea621b9"
            }
            responders {
              type      = "user"
              username  = "alarm@example.com"
            }
	        }`, testRecipientOpsgenieValue),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cloudamqp_notification.opsgenie_recipient", "type", "opsgenie"),
					resource.TestCheckResourceAttr("cloudamqp_notification.opsgenie_recipient", "value", testRecipientOpsgenieValue),
					resource.TestCheckResourceAttr("cloudamqp_notification.opsgenie_recipient", "name", "OpsGenie"),
					resource.TestCheckResourceAttr("cloudamqp_notification.opsgenie_recipient", "responders.#", "2"),
					resource.TestCheckResourceAttr("cloudamqp_notification.opsgenie_recipient", "responders.0.type", "team"),
					resource.TestCheckResourceAttr("cloudamqp_notification.opsgenie_recipient", "responders.0.id", "209402c4-a92f-4b91-a1e8-d48a8ea621b9"),
					resource.TestCheckResourceAttr("cloudamqp_notification.opsgenie_recipient", "responders.1.type", "user"),
					resource.TestCheckResourceAttr("cloudamqp_notification.opsgenie_recipient", "responders.1.username", "alarm@example.com"),
				),
			},
			{
				ResourceName:      "cloudamqp_notification.opsgenie_recipient",
				ImportStateIdFunc: testAccImportCombinedIdFunc("1091", "cloudamqp_notification.opsgenie_recipient"),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNotification_PagerDuty(t *testing.T) {
	t.Parallel()

	// Set sanitized value for playback and use real value for recording
	testRecipientPagerDutyValue := "RECIPIENT_PAGERDUTY_VALUE"
	if os.Getenv("CLOUDAMQP_RECORD") != "" {
		testRecipientPagerDutyValue = os.Getenv("RECIPIENT_PAGERDUTY_VALUE")
	}

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
          resource "cloudamqp_notification" "pagerduty_recipient" {
            instance_id = 1091
            type        = "pagerduty"
            value       = "%s"
            name        = "PagerDuty"
            options     = {
              "dedupkey" = "DEDUPKEY"
            }
          }`, testRecipientPagerDutyValue),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cloudamqp_notification.pagerduty_recipient", "type", "pagerduty"),
					resource.TestCheckResourceAttr("cloudamqp_notification.pagerduty_recipient", "value", testRecipientPagerDutyValue),
					resource.TestCheckResourceAttr("cloudamqp_notification.pagerduty_recipient", "name", "PagerDuty"),
					resource.TestCheckResourceAttr("cloudamqp_notification.pagerduty_recipient", "options.%", "1"),
					resource.TestCheckResourceAttr("cloudamqp_notification.pagerduty_recipient", "options.dedupkey", "DEDUPKEY"),
				),
			},
			{
				ResourceName:      "cloudamqp_notification.pagerduty_recipient",
				ImportStateIdFunc: testAccImportCombinedIdFunc("1091", "cloudamqp_notification.pagerduty_recipient"),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNotification_Signl4(t *testing.T) {
	t.Parallel()

	// Set sanitized value for playback and use real value for recording
	testRecipientSignl4Value := "RECIPIENT_SIGNL4_VALUE"
	if os.Getenv("CLOUDAMQP_RECORD") != "" {
		testRecipientSignl4Value = os.Getenv("RECIPIENT_SIGNL4_VALUE")
	}

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
          resource "cloudamqp_notification" "signl4_recipient" {
            instance_id = 1091
            type        = "signl4"
            value       = "%s"
            name        = "Signl4"
          }`, testRecipientSignl4Value),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cloudamqp_notification.signl4_recipient", "type", "signl4"),
					resource.TestCheckResourceAttr("cloudamqp_notification.signl4_recipient", "value", testRecipientSignl4Value),
					resource.TestCheckResourceAttr("cloudamqp_notification.signl4_recipient", "name", "Signl4"),
				),
			},
			{
				ResourceName:      "cloudamqp_notification.signl4_recipient",
				ImportStateIdFunc: testAccImportCombinedIdFunc("1091", "cloudamqp_notification.signl4_recipient"),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNotification_Slack(t *testing.T) {
	t.Parallel()

	// Set sanitized value for playback and use real value for recording
	testRecipientSlackValue := "RECIPIENT_SLACK_VALUE"
	if os.Getenv("CLOUDAMQP_RECORD") != "" {
		testRecipientSlackValue = os.Getenv("RECIPIENT_SLACK_VALUE")
	}

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
          resource "cloudamqp_notification" "slack_recipient" {
            instance_id = 1091
            type        = "slack"
            value       = "%s"
            name        = "Slack"
          }`, testRecipientSlackValue),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cloudamqp_notification.slack_recipient", "type", "slack"),
					resource.TestCheckResourceAttr("cloudamqp_notification.slack_recipient", "value", testRecipientSlackValue),
					resource.TestCheckResourceAttr("cloudamqp_notification.slack_recipient", "name", "Slack"),
				),
			},
			{
				ResourceName:      "cloudamqp_notification.slack_recipient",
				ImportStateIdFunc: testAccImportCombinedIdFunc("1091", "cloudamqp_notification.slack_recipient"),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNotification_Teams(t *testing.T) {
	t.Parallel()

	// Set sanitized value for playback and use real value for recording
	testRecipientTeamsValue := "RECIPIENT_TEAMS_VALUE"
	if os.Getenv("CLOUDAMQP_RECORD") != "" {
		testRecipientTeamsValue = os.Getenv("RECIPIENT_TEAMS_VALUE")
	}

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
          resource "cloudamqp_notification" "teams_recipient" {
            instance_id = 1091
            type        = "teams"
            value       = "%s"
            name        = "Teams"
          }`, testRecipientTeamsValue),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cloudamqp_notification.teams_recipient", "type", "teams"),
					resource.TestCheckResourceAttr("cloudamqp_notification.teams_recipient", "value", testRecipientTeamsValue),
					resource.TestCheckResourceAttr("cloudamqp_notification.teams_recipient", "name", "Teams"),
				),
			},
			{
				ResourceName:      "cloudamqp_notification.teams_recipient",
				ImportStateIdFunc: testAccImportCombinedIdFunc("1091", "cloudamqp_notification.teams_recipient"),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
