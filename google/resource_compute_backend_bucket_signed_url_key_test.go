package google

import (
	"fmt"
	"testing"

	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccComputeBackendBucketSignedUrlKey_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeBackendBucketSignedUrlKeyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeBackendBucketSignedUrlKey_basic(context),
				Check:  testAccCheckComputeBackendBucketSignedUrlKeyCreatedProducer(t),
			},
		},
	})
}

func testAccComputeBackendBucketSignedUrlKey_basic(context map[string]interface{}) string {
	return Nprintf(`
resource "google_compute_backend_bucket_signed_url_key" "backend_key" {
  name           = "test-key-%{random_suffix}"
  key_value      = "iAmAFakeKeyRandomBytes=="
  backend_bucket = google_compute_backend_bucket.test_backend.name
}

resource "google_compute_backend_bucket" "test_backend" {
  name        = "test-signed-backend-bucket-%{random_suffix}"
  description = "Contains beautiful images"
  bucket_name = google_storage_bucket.bucket.name
  enable_cdn  = true
}

resource "google_storage_bucket" "bucket" {
  name     = "test-storage-bucket-%{random_suffix}"
  location = "EU"
}
`, context)
}

func testAccCheckComputeBackendBucketSignedUrlKeyDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		exists, err := checkComputeBackendBucketSignedUrlKeyExists(t, s)
		if err != nil && !IsGoogleApiErrorWithCode(err, 404) {
			return err
		}
		if exists {
			return fmt.Errorf("ComputeBackendBucketSignedUrlKey still exists")
		}
		return nil
	}
}

func testAccCheckComputeBackendBucketSignedUrlKeyCreatedProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		exists, err := checkComputeBackendBucketSignedUrlKeyExists(t, s)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("expected ComputeBackendBucketSignedUrlKey to have been created")
		}
		return nil
	}
}

func checkComputeBackendBucketSignedUrlKeyExists(t *testing.T, s *terraform.State) (bool, error) {
	for name, rs := range s.RootModule().Resources {
		if rs.Type != "google_compute_backend_bucket_signed_url_key" {
			continue
		}
		if strings.HasPrefix(name, "data.") {
			continue
		}

		config := GoogleProviderConfig(t)
		keyName := rs.Primary.Attributes["name"]

		url, err := replaceVarsForTest(config, rs, "{{ComputeBasePath}}projects/{{project}}/global/backendBuckets/{{backend_bucket}}")
		if err != nil {
			return false, err
		}

		res, err := SendRequest(config, "GET", "", url, config.UserAgent, nil)
		if err != nil {
			return false, err
		}
		policyRaw, ok := res["cdnPolicy"]
		if !ok {
			return false, nil
		}

		policy := policyRaw.(map[string]interface{})
		keyNames, ok := policy["signedUrlKeyNames"]
		if !ok {
			return false, nil
		}

		// Because the sensitive key value is not returned, all we can do is verify a
		// key with this name exists and assume the key value hasn't been changed.
		for _, k := range keyNames.([]interface{}) {
			if k.(string) == keyName {
				// Just return empty map to indicate key was found
				return true, nil
			}
		}
	}

	return false, nil
}
