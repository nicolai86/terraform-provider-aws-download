package aws_download

import (
	"context"
	"fmt"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
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
			},
			"role_arn": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Amazon Resource Name of an IAM Role to assume prior to making API calls.",
			},
			"session_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "An identifier for the assumed role session.",
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
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(c.region))
	if err != nil {
		return nil, diag.FromErr(fmt.Errorf("unable to load AWS config: %w", err))
	}

	stsClient := sts.New(sts.Options{
		Credentials: cfg.Credentials,
	})
	resp, err := stsClient.AssumeRole(ctx, &sts.AssumeRoleInput{
		RoleArn:         &c.roleArn,
		RoleSessionName: &c.sessionName,
	})
	if err != nil {
		return nil, diag.Errorf("unable to assume role (%s): %w", c.roleArn, err)
	}

	// runs with assumed role
	client := &AWSClient{
		S3Conn: s3.New(s3.Options{
			Credentials: credentials.NewStaticCredentialsProvider(*resp.Credentials.AccessKeyId, *resp.Credentials.SecretAccessKey, *resp.Credentials.SessionToken),
			Region:      c.region,
		}),
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
