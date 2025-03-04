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

func TestAccTPUNode_tpuNodeBasicExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckTPUNodeDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTPUNode_tpuNodeBasicExample(context),
			},
			{
				ResourceName:            "google_tpu_node.tpu",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"zone"},
			},
		},
	})
}

func testAccTPUNode_tpuNodeBasicExample(context map[string]interface{}) string {
	return Nprintf(`

data "google_tpu_tensorflow_versions" "available" {
}

resource "google_tpu_node" "tpu" {
  name = "tf-test-test-tpu%{random_suffix}"
  zone = "us-central1-b"

  accelerator_type   = "v3-8"
  tensorflow_version = data.google_tpu_tensorflow_versions.available.versions[0]
  cidr_block         = "10.2.0.0/29"
}
`, context)
}

func TestAccTPUNode_tpuNodeFullExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckTPUNodeDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTPUNode_tpuNodeFullExample(context),
			},
			{
				ResourceName:            "google_tpu_node.tpu",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"zone"},
			},
		},
	})
}

func testAccTPUNode_tpuNodeFullExample(context map[string]interface{}) string {
	return Nprintf(`
data "google_tpu_tensorflow_versions" "available" {
}


resource "google_tpu_node" "tpu" {
  name = "tf-test-test-tpu%{random_suffix}"
  zone = "us-central1-b"

  accelerator_type = "v3-8"

  tensorflow_version = data.google_tpu_tensorflow_versions.available.versions[0]

  description = "Terraform Google Provider test TPU"
  use_service_networking = true
  network = google_service_networking_connection.private_service_connection.network

  labels = {
    foo = "bar"
  }

  scheduling_config {
    preemptible = true
  }
}

data "google_compute_network" "network" {
  name = "default"
}

resource "google_compute_global_address" "service_range" {
  name          = "tf-test%{random_suffix}"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = data.google_compute_network.network.id
}

resource "google_service_networking_connection" "private_service_connection" {
  network                 = data.google_compute_network.network.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.service_range.name]
}
`, context)
}

func testAccCheckTPUNodeDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_tpu_node" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := GoogleProviderConfig(t)

			url, err := replaceVarsForTest(config, rs, "{{TPUBasePath}}projects/{{project}}/locations/{{zone}}/nodes/{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = SendRequest(config, "GET", billingProject, url, config.UserAgent, nil)
			if err == nil {
				return fmt.Errorf("TPUNode still exists at %s", url)
			}
		}

		return nil
	}
}
