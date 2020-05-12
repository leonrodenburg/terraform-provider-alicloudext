package alicloudext

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cas"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/leonrodenburg/terraform-provider-alicloudext/pkg/certificates"
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
			"certificate_path": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"private_key_path": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
	client, _ := createCasClient(config)
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

	rawCertificate := certificates.Sanitize(string(issuedCert.Certificate))
	rawPrivateKey := certificates.Sanitize(string(issuedCert.PrivateKey))

	if certificatePath, ok := d.GetOk("certificate_path"); ok {
		err := writeToFile(certificatePath.(string), rawCertificate)
		if err != nil {
			return err
		}
	}
	if privateKeyPath, ok := d.GetOk("private_key_path"); ok {
		err := writeToFile(privateKeyPath.(string), rawPrivateKey)
		if err != nil {
			return err
		}
	}

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
	_ = d.Set("certificate_url", issuedCert.CertURL)
	_ = d.Set("certificate_stable_url", issuedCert.CertStableURL)

	return resourceCertificateRead(d, m)
}

func resourceCertificateRead(d *schema.ResourceData, m interface{}) error {
	client, _ := createCasClient(m.(Configuration))

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
	client, _ := createCasClient(m.(Configuration))

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

func createCasClient(config Configuration) (*cas.Client, error) {
	client, err := cas.NewClientWithAccessKey(config.Region, config.AccessKey, config.SecretKey)
	if err != nil {
		return nil, err
	}
	return client, err
}

func writeToFile(path string, contents string) error {
	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, []byte(contents), 0644)
}
