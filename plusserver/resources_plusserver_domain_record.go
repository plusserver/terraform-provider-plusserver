package plusserver

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/plusserver/terraform-provider-plusserver/api/dns"
	"strconv"
	"strings"
)

func resourceDomainRecord() *schema.Resource {
	return &schema.Resource{
		Description: "Use the plusserver DNS API to create/modify/delete a domain record.",
		CreateContext: resourceDomainRecordCreate,
		ReadContext:   resourceDomainRecordRead,
		UpdateContext: resourceDomainRecordUpdate,
		DeleteContext: resourceDomainRecordDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDomainRecordImport,
		},

		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type: schema.TypeInt,
				Description: "Exported ID of the domain record. Same as \"id\"",
				Required: true,
			},
			"name": {
				Type: schema.TypeString,
				Description: "Domain record name without TLD or second-level domain",
				Required: true,
			},
			"type": {
				Type: schema.TypeString,
				Default: "A",
				Description: "Domain record type. Can be either one of \"A\", \"AAAA\", \"CAA\", \"CNAME\", \"MX\", \"NS\", \"PTR\", \"SRV\", \"TXT\" or \"SOA\"." +
					"The default Value is \"A\".",
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{"A", "AAAA", "CAA", "CNAME", "MX", "NS", "PTR", "SRV", "TXT", "SOA"}, false),
				Optional: true,
			},
			"ttl": {
				Type: schema.TypeInt,
				Description: "Domain record time to live in seconds. The default value is 300 seconds",
				Default: 300,
				Optional: true,
			},
			"content": {
				Type: schema.TypeString,
				Description: "Domain record content. For example the IP address of the A record",
				Required: true,
			},
		},
	}
}

func resourceDomainRecordImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), ":", 2)
	client := meta.(*dns.Client)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return []*schema.ResourceData{}, fmt.Errorf("unexpected format of ID (%s), expected domainId:recordId", parts)
	}

	// Split ID separate by domainId and recordId
	domainId, convErr := strconv.Atoi(parts[0])
	convErr = d.Set("domain_id", domainId)
	if convErr != nil {
		return []*schema.ResourceData{}, convErr
	}
	// Set recordId to resource
	d.SetId(parts[1])

	// Get all records in rrset
	records, err := client.GetRecords(ctx, domainId)
	if err != nil {
		return []*schema.ResourceData{}, err
	}
	// Get record by Id
	domainRecord := getRecordResourceById(d.Id(), records.DnsResourceRecordList)
	if domainRecord == nil {
		return []*schema.ResourceData{}, errors.New("could not find record by id in rrset domain")
	}
	err = d.Set("content", domainRecord.Content)
	err = d.Set("name", domainRecord.Name)
	err = d.Set("ttl", domainRecord.Ttl)
	err = d.Set("type", domainRecord.Type)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}

func resourceDomainRecordUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*dns.Client)
	var diags diag.Diagnostics

	content := d.Get("content").(string)
	domainId := d.Get("domain_id").(int)
	name := d.Get("name").(string)
	ttl := d.Get("ttl").(int)

	_, err := client.UpdateRecord(ctx, domainId, d.Id(), content, ttl)
	if err != nil {
		return diag.FromErr(err)
	}

	// scan again to get the updated record_id
	records, err := client.GetRecords(ctx, domainId)
	if err != nil {
		return diag.FromErr(err)
	}

	recordId := getRecordResourceID(name, content, records.DnsResourceRecordList)
	if recordId == "" {
		//NOTE: If this happens we are all screwed since I can't remove the record because I don't have the ID
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to save resource id",
			Detail:   "the created record was not found in the domain",
		})
		return diags
	}

	d.SetId(recordId)

	return diags
}

func resourceDomainRecordDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*dns.Client)
	var diags diag.Diagnostics

	domainId := d.Get("domain_id").(int)

	_, err := client.DeleteRecord(ctx, domainId, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// implied but we explicitly set it here
	d.SetId("")

	return diags
}

func resourceDomainRecordRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*dns.Client)
	var diags diag.Diagnostics

	content := d.Get("content").(string)
	domainId := d.Get("domain_id").(int)
	name := d.Get("name").(string)

	records, err := client.GetRecords(ctx, domainId)

	if err != nil {
		return diag.FromErr(err)
	}

	recordId := getRecordResourceID(name, content, records.DnsResourceRecordList)
	if recordId == "" {
		//NOTE: If this happens we are all screwed since I can't remove the record because I don't have the ID
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to save resource id",
			Detail:   "the created record was not found in the domain",
		})
		return diags
	}

	return diags
}


func resourceDomainRecordCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*dns.Client)
	var diags diag.Diagnostics

	content := d.Get("content").(string)
	domainId := d.Get("domain_id").(int)
	dnsType := d.Get("type").(string)
	name := d.Get("name").(string)
	ttl := d.Get("ttl").(int)

	_, err := client.CreateRecord(ctx, &dns.RecordCreateRequest{DnsResourceRecordList: []struct {
		Content     string `json:"content"`
		DnsDomainId int    `json:"dnsDomainId"`
		Type        string `json:"type"`
		Name        string `json:"name"`
		Ttl         int    `json:"ttl"`
	}{{content, domainId, dnsType, name, ttl}}})

	if err != nil {
		return diag.FromErr(err)
	}

	records, err := client.GetRecords(ctx, domainId)
	if err != nil {
		return diag.FromErr(err)
	}

	recordId := getRecordResourceID(name, content, records.DnsResourceRecordList)
	if recordId == "" {
		//NOTE: If this happens we are all screwed since I can't remove the stupid record because I don't have the ID
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to save resource id",
			Detail:   "the created record was not found in the domain",
		})
		return diags
	}

	d.SetId(recordId)

	return diags
}

func getRecordResourceID(name string, content string, items []*dns.RecordsResponseEntry) string {
	if items != nil {
		for _, recordItem := range items {
			if recordItem.Name == name && recordItem.Content == content {
				return recordItem.DnsResourceRecordId
			}
		}
		return ""
	}

	return ""
}

func getRecordResourceById(id string, items []*dns.RecordsResponseEntry) *dns.RecordsResponseEntry {
	if items != nil {
		for _, recordItem := range items {
			if recordItem.DnsResourceRecordId == id {
				return recordItem
			}
		}
		return nil
	}

	return nil
}