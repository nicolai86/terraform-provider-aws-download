package aws_download

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAwsS3DownloadObject() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAwsS3DownloadObjectRead,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"bucket": &schema.Schema{
				Type:     schema.TypeString,
				Computed: false,
			},
			"key": &schema.Schema{
				Type:     schema.TypeString,
				Computed: false,
			},
			"filename": &schema.Schema{
				Type:     schema.TypeString,
				Computed: false,
			},
		},
	}
}

func dataSourceAwsS3DownloadObjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	conn := m.(*AWSClient).S3Conn

	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)

	input := s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	log.Printf("[DEBUG] Reading S3 Object: %s", input)
	out, err := conn.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed getting S3 Bucket (%s) Object (%s): %w", bucket, key, err))
	}

	filename := d.Get("filename").(string)
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to open file (%s) for writing: %w", filename, err))
	}
	if _, err := io.Copy(f, out.Body); err != nil {
		return diag.FromErr(fmt.Errorf("failed to copy s3 content: %w", err))
	}
	defer f.Close()

	uniqueId := bucket + "/" + key
	d.SetId(uniqueId)

	return nil
}
