package alicloudssl

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/bobesa/go-domain-util/domainutil"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceApiGatewayDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceApiGatewayDomainCreate,
		Read:   resourceApiGatewayDomainRead,
		Delete: resourceApiGatewayDomainDelete,

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
			"record_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceApiGatewayDomainCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(Configuration).Client

	groupId := d.Get("group_id").(string)
	group, err := fetchApiGroup(client, groupId)
	if err != nil {
		return err
	}

	record, err := ensureCnameForApiGroup(client, d.Get("domain").(string), group.SubDomain)
	if err != nil {
		return err
	}

	req := cloudapi.CreateSetDomainRequest()
	req.GroupId = groupId
	req.DomainName = d.Get("domain").(string)
	res := cloudapi.CreateSetDomainResponse()
	err = client.DoAction(req, res)
	if err != nil {
		return err
	}

	d.SetId(res.RequestId)
	_ = d.Set("record_id", record.RecordId)

	return resourceApiGatewayDomainRead(d, m)
}

func resourceApiGatewayDomainRead(d *schema.ResourceData, m interface{}) error {
	client := m.(Configuration).Client
	req := cloudapi.CreateDescribeDomainRequest()
	req.GroupId = d.Get("group_id").(string)
	req.DomainName = d.Get("domain").(string)
	res := cloudapi.CreateDescribeDomainResponse()

	err := client.DoAction(req, res)
	if err != nil {
		d.SetId("")
		return err
	}

	_ = d.Set("domain", res.DomainName)
	_ = d.Set("group_id", res.GroupId)

	return nil
}

func resourceApiGatewayDomainDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(Configuration).Client
	dnsReq := alidns.CreateDeleteDomainRecordRequest()
	dnsReq.RecordId = d.Get("record_id").(string)
	dnsRes := alidns.CreateDeleteDomainRecordResponse()
	err := client.DoAction(dnsReq, dnsRes)
	if err != nil {
		return err
	}

	req := cloudapi.CreateDeleteDomainRequest()
	req.GroupId = d.Get("group_id").(string)
	req.DomainName = d.Get("domain").(string)
	res := cloudapi.CreateDeleteDomainResponse()
	err = client.DoAction(req, res)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func ensureCnameForApiGroup(client *sdk.Client, domain string, value string) (*alidns.Record, error) {
	subdomain := "@"
	if domainutil.HasSubdomain(domain) {
		subdomain = domainutil.Subdomain(domain)
	}
	naked := domainutil.Domain(domain)

	dnsReq := alidns.CreateDescribeDomainRecordsRequest()
	dnsReq.DomainName = naked
	dnsReq.RRKeyWord = subdomain
	dnsReq.Type = "CNAME"
	dnsReq.ValueKeyWord = value
	dnsRes := alidns.CreateDescribeDomainRecordsResponse()
	err := client.DoAction(dnsReq, dnsRes)
	if dnsRes.TotalCount > 0 {
		return &dnsRes.DomainRecords.Record[0], nil
	}

	req := alidns.CreateAddDomainRecordRequest()
	req.DomainName = naked
	req.RR = subdomain
	req.Type = "CNAME"
	req.Value = value
	res := alidns.CreateAddDomainRecordResponse()
	err = client.DoAction(req, res)
	if err != nil {
		return nil, err
	}
	return &alidns.Record{RecordId: res.RecordId}, nil
}

func fetchApiGroup(client *sdk.Client, groupId string) (*cloudapi.DescribeApiGroupResponse, error) {
	req := cloudapi.CreateDescribeApiGroupRequest()
	req.GroupId = groupId
	res := cloudapi.CreateDescribeApiGroupResponse()
	err := client.DoAction(req, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
