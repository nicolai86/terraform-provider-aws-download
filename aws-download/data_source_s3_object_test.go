package aws_download

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAWSDownloadS3Object(t *testing.T) {
	bucket := os.Getenv("BUCKET_NAME")
	key := os.Getenv("BUCKET_KEY")
	target := "/tmp/target"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckS3Object(bucket, key, target),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDownloadExists("data.aws-download_s3_object.new"),
				),
			},
		},
	})
}

func testAccCheckS3Object(bucket, key, target string) string {
	return fmt.Sprintf(`
	data "aws-download_s3_object" "new" {
		bucket = %q
		key = %q
		filename = %q
	}
	`, bucket, key, target)
}

func testAccCheckDownloadExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No OrderID set")
		}

		_, err := os.Stat(rs.Primary.Attributes["filename"])

		return err
	}
}
