package alicloudssl

import (
	"strconv"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cas"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/leonrodenburg/terraform-provider-alicloudssl/pkg/certificates"
)

func resourceCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceCertificateCreate,
		Read:   resourceCertificateRead,
		Delete: resourceCertificateDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"email": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"certificate": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certificate_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certificate_stable_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCertificateCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(Configuration)
	client := config.Client
	certUser := certificates.User{
		Email: d.Get("email").(string),
	}

	certClient, err := certificates.CreateCertClientForUser(&certUser)
	if err != nil {
		return err
	}

	dnsProvider, err := certificates.CreateAlicloudDnsProvider(
		config.AccessKey,
		config.SecretKey,
		config.Region,
	)
	if err != nil {
		return err
	}

	issuedCert, err := certificates.RequestCertificateForDomainUsingDns(
		d.Get("domain").(string),
		certClient,
		dnsProvider,
	)
	if err != nil {
		return err
	}

	rawCertificate := sanitizePem(string(issuedCert.Certificate))
	rawPrivateKey := sanitizePem(string(issuedCert.PrivateKey))

	req := cas.CreateCreateUserCertificateRequest()
	req.Cert = rawCertificate
	req.Key = rawPrivateKey
	req.Name = d.Get("name").(string)
	res := cas.CreateCreateUserCertificateResponse()

	err = client.DoAction(req, res)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(res.CertId))
	_ = d.Set("certificate", rawCertificate)
	_ = d.Set("private_key", rawPrivateKey)
	_ = d.Set("certificate_url", issuedCert.CertURL)
	_ = d.Set("certificate_stable_url", issuedCert.CertStableURL)

	return resourceCertificateRead(d, m)
}

func sanitizePem(s string) string {
	return strings.Replace(s, "\n\n", "\n", -1)
}

func resourceCertificateRead(d *schema.ResourceData, m interface{}) error {
	client := m.(Configuration).Client

	req := cas.CreateDescribeUserCertificateDetailRequest()
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	req.CertId = requests.NewInteger(id)
	res := cas.CreateDescribeUserCertificateDetailResponse()

	err = client.DoAction(req, res)
	if err != nil {
		d.SetId("")
		return err
	}

	_ = d.Set("name", res.Name)
	_ = d.Set("domain", res.Sans)

	return nil
}

func resourceCertificateDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(Configuration).Client

	req := cas.CreateDeleteUserCertificateRequest()
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	req.CertId = requests.NewInteger(id)
	res := cas.CreateDeleteUserCertificateResponse()

	err = client.DoAction(req, res)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
