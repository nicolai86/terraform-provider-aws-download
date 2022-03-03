package aws_download

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"aws-download": testAccProvider,
	}
}

func testAccPreCheck(t *testing.T) {
	if err := os.Getenv("AWS_ACCESS_KEY_ID"); err == "" {
		t.Fatal("AWS_ACCESS_KEY_ID must be set for acceptance tests")
	}
	if err := os.Getenv("AWS_SECRET_ACCESS_KEY"); err == "" {
		t.Fatal("AWS_SECRET_ACCESS_KEY must be set for acceptance tests")
	}
}
