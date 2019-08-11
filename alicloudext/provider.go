package alicloudext

import (
	"os"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"alicloudext_certificate": resourceCertificate(),
			"alicloudext_api_gateway_domain": resourceApiGatewayDomain(),
			"alicloudext_api_gateway_domain_certificate": resourceApiGatewayDomainCertificate(),
		},
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALICLOUD_ACCESS_KEY", os.Getenv("ALICLOUD_ACCESS_KEY")),
				Description: "Access key of your account",
			},
			"secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALICLOUD_SECRET_KEY", os.Getenv("ALICLOUD_SECRET_KEY")),
				Description: "Secret key of your account",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALICLOUD_REGION", os.Getenv("ALICLOUD_REGION")),
				Description: "Region to deploy resources in",
			},
		},
		ConfigureFunc: providerConfigure,
	}
}

type Configuration struct {
	AccessKey string
	SecretKey string
	Region    string
	Client    *sdk.Client
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	accessKey := d.Get("access_key").(string)
	secretKey := d.Get("secret_key").(string)
	region := d.Get("region").(string)

	client, err := sdk.NewClientWithAccessKey(region, accessKey, secretKey)
	if err != nil {
		return nil, err
	}
	return Configuration{
		AccessKey: accessKey,
		SecretKey: secretKey,
		Region:    region,
		Client:    client,
	}, nil
}
