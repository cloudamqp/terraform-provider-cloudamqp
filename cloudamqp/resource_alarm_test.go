package cloudamqp

// import (
// 	"fmt"
// 	"log"
// 	"testing"

// 	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
// )

// func TestAccAlarm_Basic(t *testing.T) {
// 	//instance_id := 195
// 	resource_name := "cloudamqp_alarm.connection_01"

// 	resource.ParallelTest(t, resource.TestCase{
// 		PreCheck:     func() { testAccPrecheck(t) },
// 		Providers:    testAccProviders,
// 		CheckDestroy: testAccCheckInstanceDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccInstanceConfig_basic(name),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckInstanceExists(resource_name),
// 					resource.TestCheckResourceAttr(resource_name, "type", connection),
// 					resource.TestCheckResourceAttr(resource_name, "value_threshold", 0),
// 					resource.TestCheckResourceAttr(resource_name, "time_threshold", 60),
// 					resource.TestCheckResourceAttr(resource_name, "queue_regex", nil),
// 					resource.TestCheckResourceAttr(resource_name, "vhost_regex", nil),
// 					resource.TestCheckResourceAttr(resource_name, "notification_ids", []),
// 				),
// 			},
// 		},
// 	})
// }

// func testAccCheckAlarmExist(instance_name, resource_name string) resource.TestCheckFunc {
// 	log.Printf("[DEBUG] resource_alarm::testAccCheckAlarmExist resource: %s", resource_name)

// 	return func(state *terraform.State) error {
// 		rs, ok := state.RootModule().Resources[resource_name]
// 		if !ok {
// 			return fmt.Errorf("Resource: %s not found", resource_name)
// 		}
// 		if rs.Primary.ID == "" {
// 			return fmt.Errorf("No Record ID is set")
// 		}
// 		id := rs.Primary.ID

// 		rs, ok = state.RootModule().Resources[instance_name]
// 		if !ok {
// 			return fmt.Errorf("Instance resource not found")
// 		}
// 		if rs.Primary.ID == "" {
// 			return fmt.Errorf("No Record ID is set for instance")
// 		}
// 		instance_id, _ := strconv.Atoi(rs.Primary.ID)

// 		api := testAccProvider.Meta().(*api.API)
// 		data, err := api.ReadAlarm(instance_id, id)
// 		log.Printf("[DEBUG] resource_alarm::testAccCheckAlarmExist data: %v", data)
// 		if err != nil {
// 			return fmt.Errorf("Error fetching item with resource %s. %s", resource_name, err)
// 		}
// 		return nil
// 	}
// }

// func testAccCheckAlarmDestroy(s *terraform.State) error {
// 	instance_name := "cloudamqp_instance.instance_03"
// 	resource_name := "cloudamqp_alarm.alarm"

// 	return nil
// }

// /*
// func testAccCheckAlarmExists(id int, alarm *cloudamqp.alarm) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {
// 		rs, ok := s.RootModule.Resource[id]
// 		if !ok {
// 			return fmt.Errorf("alarm not created: %v", id)
// 		}
// 		if rs.Primary.ID == "" {
// 			return fmt.Errorf("no alarm ID is set")
// 		}

// 		api := testAccProvider.Meta().(*api.API)

// 		foundAlarm = &cloudamqp.alarm{}

// 	}
// }
// */

// /*
// func testAccCheckAlarmDestroy(s *terraform.State) error {
// 	log.Printf("resource_instance::testAccCheckInstanceDestroy")
// 	api := testAccProvider.Meta().(*api.API)
// 	// Read out instance id, needed to make the API call
// 	// How to make sure I get the correct instance?
// 	instance_id := 0
// 	for _, rs := range s.RootModule().Resource {
// 		if rs.Type != "cloudamqp_instance" {
// 			instance = rs.Primary.ID
// 		}
// 	}

// 	if instance_id == 0 {
// 		return fmt.Errorf("No instance found")
// 	}
// 	// Read out alarm
// 	for _, rs := range s.RootModule().Resources {
// 		if rs.Type != "cloudamqp_alarm" {
// 			continue
// 		}

// 		// Make API call to make check if the instance still exists.
// 		_, err := api.ReadAlarm(instance_id, rs.Primary.ID)
// 		if err == nil {
// 			return fmt.Errorf("Alert still exists")
// 		}
// 		notFoundErr := "not found"
// 		expectedErr := regexp.MustCompile(notFoundErr)
// 		if !expectedErr.Match([]byte(err.Error())) {
// 			return fmt.Errorf("expected %s, got %s", notFoundErr, err)
// 		}
// 	}
// 	return nil
// }
// */

// func testAccNotification_Recipient() string {
// 	log.Printf("[DEBUG] resource_notification::testAccNotification_Recipient")
// 	return fmt.Sprintf(`
// 		resource "cloudamqp_instance" "instance" {
// 			name 				= "terraform-alarm-test"
// 			nodes 			= 1
// 			plan  			= "bunny"
// 			region 			= "amazon-web-services::eu-north-1"
// 			rmq_version = "3.8.2"
// 			tags 				= ["terraform"]
// 		}

// 		resource "cloudamqp_amqp" "connection_alarm" {
// 			instance_id 			= cloudamqp_instance.instance.id
// 			type 							= connection
// 			value_threshold 	= 0
// 			time_threshold 		= 60
// 		}
// 		`)
// }
