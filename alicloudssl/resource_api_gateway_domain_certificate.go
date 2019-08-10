package alicloudssl

import (
	"errors"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cas"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/leonrodenburg/terraform-provider-alicloudssl/pkg/certificates"
)

func resourceApiGatewayDomainCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceApiGatewayDomainCertificateCreate,
		Read:   resourceApiGatewayDomainCertificateRead,
		Delete: resourceApiGatewayDomainCertificateDelete,

		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"certificate_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"system_certificate_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceApiGatewayDomainCertificateCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(Configuration).Client

	certRes, err := fetchCertificateById(d.Get("certificate_id").(int), client)
	if err != nil {
		return err
	}

	req := cloudapi.CreateSetDomainCertificateRequest()
	req.GroupId = d.Get("group_id").(string)
	req.DomainName = d.Get("domain").(string)
	req.CertificateName = certRes.Name
	req.CertificateBody = certificates.Sanitize(certRes.Cert)
	req.CertificatePrivateKey = certificates.Sanitize(certRes.Key)
	res := cloudapi.CreateSetDomainCertificateResponse()
	err = client.DoAction(req, res)
	if err != nil {
		return err
	}

	d.SetId(res.RequestId)

	return resourceApiGatewayDomainCertificateRead(d, m)
}

func resourceApiGatewayDomainCertificateRead(d *schema.ResourceData, m interface{}) error {
	client := m.(Configuration).Client

	req := cloudapi.CreateDescribeDomainRequest()
	req.GroupId = d.Get("group_id").(string)
	req.DomainName = d.Get("domain").(string)
	res := cloudapi.CreateDescribeDomainResponse()
	err := client.DoAction(req, res)
	if err != nil {
		return err
	}

	if len(res.CertificateId) < 0 {
		d.SetId("")
		return errors.New("no certificate bound to API Gateway group")
	}

	_ = d.Set("system_certificate_id", res.CertificateId)

	return nil
}

func resourceApiGatewayDomainCertificateDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(Configuration).Client

	req := cloudapi.CreateDeleteDomainCertificateRequest()
	req.GroupId = d.Get("group_id").(string)
	req.DomainName = d.Get("domain").(string)
	req.CertificateId = d.Get("system_certificate_id").(string)
	res := cloudapi.CreateDeleteDomainCertificateResponse()
	err := client.DoAction(req, res)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func fetchCertificateById(certId int, client *sdk.Client) (*cas.DescribeUserCertificateDetailResponse, error) {
	certReq := cas.CreateDescribeUserCertificateDetailRequest()
	certReq.CertId = requests.NewInteger(certId)
	certRes := cas.CreateDescribeUserCertificateDetailResponse()
	err := client.DoAction(certReq, certRes)
	if err != nil {
		return nil, err
	}
	return certRes, nil
}
