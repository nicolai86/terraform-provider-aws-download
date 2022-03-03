package aws_download

import (
	"context"
	"log"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	awsbase "github.com/hashicorp/aws-sdk-go-base/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "AWS Region",
				DefaultFunc: schema.EnvDefaultFunc("AWS_REGION", "us-west-2"),
			},
			"role_arn": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Amazon Resource Name of an IAM Role to assume prior to making API calls.",
				DefaultFunc: schema.EnvDefaultFunc("AWS_ROLE_ARN", nil),
			},
			"session_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "An identifier for the assumed role session.",
				DefaultFunc: schema.EnvDefaultFunc("AWS_SESSION_NAME", "terraform"),
				ValidateFunc: validation.All(
					validation.StringLenBetween(2, 64),
					validation.StringMatch(regexp.MustCompile(`[\w+=,.@\-]*`), ""),
				),
			},
		},
		ResourcesMap: map[string]*schema.Resource{},
		DataSourcesMap: map[string]*schema.Resource{
			"aws-download_s3_object": dataSourceAwsS3DownloadObject(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

type Config struct {
	region      string
	sessionName string
	roleArn     string
}

type AWSClient struct {
	S3Conn *s3.Client
}

func (c *Config) Client(ctx context.Context) (interface{}, diag.Diagnostics) {
	// runs with credentials provided in ENV
	basecfg := awsbase.Config{
		AssumeRole: &awsbase.AssumeRole{
			RoleARN:     c.roleArn,
			SessionName: c.sessionName,
		},
		Region: c.region,
	}
	cfg, err := awsbase.GetAwsConfig(ctx, &basecfg)
	if err != nil {
		return nil, diag.Errorf("error configuring Terraform AWS Provider: %s", err)
	}
	log.Println("[INFO] got AWS aconfig")

	// runs with assumed role
	client := &AWSClient{
		S3Conn: s3.NewFromConfig(cfg),
	}
	return client, nil
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	config := Config{
		region:      d.Get("region").(string),
		roleArn:     d.Get("role_arn").(string),
		sessionName: d.Get("session_name").(string),
	}

	return config.Client(ctx)
}
