package plusserver

import (
	"context"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/plusserver/terraform-provider-plusserver/api/dns"
	"strconv"
)

func resourceDomain() *schema.Resource {
	return &schema.Resource{
		Description: "Use the plusserver DNS API to create/modify/delete a domain.",
		CreateContext: resourceDomainCreate,
		ReadContext:   resourceDomainRead,
		UpdateContext: resourceDomainUpdate,
		DeleteContext: resourceDomainDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"unicode_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Description: "Domain name in unicode",
				Required: true,
			},
			"domain_id": {
				Type: schema.TypeInt,
				Description: "Exported ID of the domain. Same as \"id\"",
				Computed: true,
				Optional: true,
			},
			"company_id": {
				Type:     schema.TypeString,
				Description: "Company ID to set with the domain as metadata",
				Optional: true,
			},
			"protected": {
				Type:     schema.TypeBool,
				Default:  false,
				ForceNew: true,
				Description: "Protects the domain from accidental deletion. " +
					"However this provider will be unable to delete the domain if set to true. " +
					"The default Value is false",
				Optional: true,
			},
			"replication_type": {
				Type:     schema.TypeString,
				Description: "Domain replication type. Can be either Master, Slave, Native or None",
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{"Master", "Slave", "Native", "None"}, false),
				Required: true,
			},
			"dns_nameserver_pair_name": {
				Type: schema.TypeString,
				ForceNew: true,
				Default: "ns1.plusserver.com",
				Description: "Domain pair identifier. The default value is ns1.plusserver.com",
				Optional: true,
			},
			"replication_master_ip_address_list": {
				Type: schema.TypeList,
				Required: true,
				Description: "List of replication servers",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"contract_id": {
				Type:     schema.TypeString,
				Description: "Contract ID to set with the domain as metadata",
				Optional: true,
			},
		},
	}
}

func resourceDomainUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*dns.Client)

	domainId := d.Id()
	companyId := d.Get("company_id").(string)
	protected := d.Get("protected").(bool)
	contractId := d.Get("contract_id").(string)
	replicationMasterIpAddressList := d.Get("replication_master_ip_address_list").([]interface{})

	var ips []string
	for _, ipAddress := range replicationMasterIpAddressList {
		elem := ipAddress.(string)
		ips = append(ips, elem)
	}

	resp, err := client.UpdateDomain(ctx, domainId, &dns.UpdateDomain{DnsDomain: struct {
		ReplicationMasterIpAddressList []string `json:"replicationMasterIpAddressList"`
		Protected                      bool     `json:"protected"`
		CompanyId                      string   `json:"companyId"`
		ContractId                     string   `json:"contractId"`
	}(struct {
		ReplicationMasterIpAddressList []string
		Protected                      bool
		CompanyId                      string
		ContractId                     string
	}{ReplicationMasterIpAddressList: ips, Protected: protected, CompanyId: companyId, ContractId: contractId})})
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("domain_id", resp.DnsDomain.DnsDomainId)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(resp.DnsDomain.DnsDomainId))

	return resourceDomainRead(ctx, d, m)
}

func resourceDomainDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*dns.Client)
	var diags diag.Diagnostics

	domainId := d.Id()

	_, err := client.DeleteDomain(ctx, domainId)
	if err != nil {
		return diag.FromErr(err)
	}

	// implied but we explicitly set it here
	d.SetId("")

	return diags
}


func resourceDomainRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*dns.Client)
	var diags diag.Diagnostics
	var result error

	domainId := d.Id()


	resp, err := client.GetDomainById(ctx, domainId)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("unicode_name", resp.DnsDomain.Name); err != nil {
		result = multierror.Append(result, err)
	}
	if err = d.Set("protected", resp.DnsDomain.Protected); err != nil {
		result = multierror.Append(result, err)
	}
	if err = d.Set("replication_type", resp.DnsDomain.ReplicationType); err != nil {
		result = multierror.Append(result, err)
	}
	if err = d.Set("replication_master_ip_address_list", resp.DnsDomain.ReplicationMasterIpAddressList); err != nil {
		result = multierror.Append(result, err)
	}
	if err = d.Set("contract_id", resp.DnsDomain.ContractId); err != nil {
		result = multierror.Append(result, err)
	}

	if result != nil {
		return diag.FromErr(result)
	}


	return diags
}

func resourceDomainCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*dns.Client)
	var diags diag.Diagnostics

	unicodeName := d.Get("unicode_name").(string)
	companyId := d.Get("company_id").(string)
	protected := d.Get("protected").(bool)
	replicationType := d.Get("replication_type").(string)
	contractId := d.Get("contract_id").(string)
	dnsNameserverPairName := d.Get("dns_nameserver_pair_name").(string)

	replicationMasterIpAddressList := d.Get("replication_master_ip_address_list").([]interface{})

	var ips []string
	for _, ipAddress := range replicationMasterIpAddressList {
		elem := ipAddress.(string)
		ips = append(ips, elem)
	}

	data := dns.CreateDomain{DnsDomain: struct {
		UnicodeName                    string   `json:"unicodeName"`
		CompanyId                      string   `json:"companyId"`
		DnsNameserverPairName          string   `json:"dnsNameserverPairName"`
		Protected                      bool     `json:"protected"`
		ReplicationType                string   `json:"replicationType"`
		ReplicationMasterIpAddressList []string `json:"replicationMasterIpAddressList"`
		ContractId                     string   `json:"contractId"`
	}(struct {
		UnicodeName                    string
		CompanyId                      string
		DnsNameserverPairName          string
		Protected                      bool
		ReplicationType                string
		ReplicationMasterIpAddressList []string
		ContractId                     string
	}{UnicodeName: unicodeName, CompanyId: companyId, DnsNameserverPairName: dnsNameserverPairName, Protected: protected,
		ReplicationType: replicationType, ReplicationMasterIpAddressList: ips, ContractId: contractId})}

	resp, err := client.CreateDomain(ctx, &data)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("domain_id", resp.DnsDomain.DnsDomainId)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(resp.DnsDomain.DnsDomainId))


	return diags
}
