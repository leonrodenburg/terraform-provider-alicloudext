package alicloudssl

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceCertificateCreate,
		Read:   resourceCertificateRead,
		Update: resourceCertificateUpdate,
		Delete: resourceCertificateDelete,

		Schema: map[string]*schema.Schema{
			"domain": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceCertificateCreate(d *schema.ResourceData, m interface{}) error {
	domain := d.Get("domain").(string)
	d.SetId(domain)
	return resourceCertificateRead(d, m)
}

func resourceCertificateRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceCertificateUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceCertificateRead(d, m)
}

func resourceCertificateDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
