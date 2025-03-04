// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package google

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataprocAutoscalingPolicy_dataprocAutoscalingPolicyBasicExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataprocAutoscalingPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocAutoscalingPolicy_dataprocAutoscalingPolicyBasicExample(context),
			},
			{
				ResourceName:            "google_dataproc_autoscaling_policy.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
		},
	})
}

func testAccDataprocAutoscalingPolicy_dataprocAutoscalingPolicyBasicExample(context map[string]interface{}) string {
	return Nprintf(`
resource "google_dataproc_autoscaling_policy" "basic" {
  policy_id = "tf-test-dataproc-policy%{random_suffix}"
  location  = "us-central1"

  worker_config {
    max_instances = 3
  }

  basic_algorithm {
    yarn_config {
      graceful_decommission_timeout = "30s"

      scale_up_factor   = 0.5
      scale_down_factor = 0.5
    }
  }
}
`, context)
}

func TestAccDataprocAutoscalingPolicy_dataprocAutoscalingPolicyExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataprocAutoscalingPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocAutoscalingPolicy_dataprocAutoscalingPolicyExample(context),
			},
			{
				ResourceName:            "google_dataproc_autoscaling_policy.asp",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
		},
	})
}

func testAccDataprocAutoscalingPolicy_dataprocAutoscalingPolicyExample(context map[string]interface{}) string {
	return Nprintf(`
resource "google_dataproc_cluster" "basic" {
  name     = "tf-test-dataproc-policy%{random_suffix}"
  region   = "us-central1"

  cluster_config {
    autoscaling_config {
      policy_uri = google_dataproc_autoscaling_policy.asp.name
    }
  }
}

resource "google_dataproc_autoscaling_policy" "asp" {
  policy_id = "tf-test-dataproc-policy%{random_suffix}"
  location  = "us-central1"

  worker_config {
    max_instances = 3
  }

  basic_algorithm {
    yarn_config {
      graceful_decommission_timeout = "30s"

      scale_up_factor   = 0.5
      scale_down_factor = 0.5
    }
  }
}
`, context)
}

func testAccCheckDataprocAutoscalingPolicyDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_dataproc_autoscaling_policy" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := GoogleProviderConfig(t)

			url, err := replaceVarsForTest(config, rs, "{{DataprocBasePath}}projects/{{project}}/locations/{{location}}/autoscalingPolicies/{{policy_id}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = SendRequest(config, "GET", billingProject, url, config.UserAgent, nil)
			if err == nil {
				return fmt.Errorf("DataprocAutoscalingPolicy still exists at %s", url)
			}
		}

		return nil
	}
}
