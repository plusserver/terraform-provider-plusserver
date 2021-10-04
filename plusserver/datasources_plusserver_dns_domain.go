package plusserver

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/plusserver/terraform-provider-plusserver/api/dns"
	"log"
	"strconv"
	"time"
)

func dataSourceDomain() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDomainRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type: schema.TypeString,
				Required: true,
			},
			"domain_id": {
				Type: schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"domains": {
				Type: schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain_id": {
							Type: schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type: schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}


func dataSourceDomainRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*dns.Client)

	var diags diag.Diagnostics

	domains, err := client.SearchDomains(ctx, d.Get("name").(string))

	log.Printf("[INFO] Received: %+v", domains)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(domains.DnsDomainList) > 1 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "ambiguous domain selector",
			Detail:   "the api has returned more than one domain that matched your criteria",
		})
		return diags
	} else if len(domains.DnsDomainList) == 0{
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "no domain found",
			Detail:   "the api has returned no domain that matched your criteria",
		})
		return diags
	}

	domainItems := flattenDomainsData(&domains.DnsDomainList)
	log.Printf("[INFO] domainItems: %+v", domainItems)
	if err = d.Set("domains", domainItems); err != nil {
		return diag.FromErr(err)
	}

	// select first element and set the domain_id as domains[0].domain_id
	err = d.Set("domain_id", domains.DnsDomainList[0].DnsDomainId)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "error setting domain_id",
			Detail:   err.Error(),
		})
		return diags
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func flattenDomainsData(items *[]dns.DomainListResponse) []interface{} {
	if items != nil {
		domainItems := make([]interface{}, len(*items), len(*items))

		for i, domainItem := range *items {
			di := make(map[string]interface{})
			di["domain_id"] = domainItem.DnsDomainId
			di["name"] = domainItem.Name

			domainItems[i] = di
		}
		return domainItems
	}

	return make([]interface{}, 0)
}