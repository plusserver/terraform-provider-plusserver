package plusserver

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/plusserver/terraform-provider-plusserver/api"
	"github.com/plusserver/terraform-provider-plusserver/api/dns"
)

func Provider() *schema.Provider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"client_id": {
				Type: schema.TypeString,
				Required: true,
				DefaultFunc: schema.EnvDefaultFunc("CLIENT_ID", ""),
				Description: "the client id of the keycloak app",
			},
			"client_secret": {
				Type: schema.TypeString,
				Required: true,
				Sensitive: true,
				DefaultFunc: schema.EnvDefaultFunc("CLIENT_SECRET", ""),
				Description: "the client secret of the keycloak app",
			},
			"username": {
				Type: schema.TypeString,
				Required: true,
				DefaultFunc: schema.EnvDefaultFunc("USERNAME", ""),
				Description: "the username to authenticate against keycloak",
			},
			"password": {
				Type: schema.TypeString,
				Required: true,
				Sensitive: true,
				DefaultFunc: schema.EnvDefaultFunc("PASSWORD", ""),
				Description: "the password to authenticate against keycloak",
			},
			"token_url": {
				Type: schema.TypeString,
				Required: true,
				DefaultFunc: schema.EnvDefaultFunc("TOKEN_URL", ""),
				Description: "the keycloak token url",
			},
			"env": {
				Type: schema.TypeString,
				Optional: true,
				DefaultFunc: schema.EnvDefaultFunc("API_ENV", "test"),
				Description: "api environment",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"plusserver_domain": dataSourceDomain(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"plusserver_domain":        resourceDomain(),
			"plusserver_domain_record": resourceDomainRecord(),
		},
		ConfigureContextFunc: providerConfigure,
	}

	return p
}

func buildAPI(d *schema.ResourceData) string {
	var env string
	if d.Get("env").(string) == "prod" {
		// prod has no prefix
		return fmt.Sprintf("https://tool.ps-intern.de/api-gateway-legacy/gateway/entity")
	} else {
		env = d.Get("env").(string)
		return fmt.Sprintf("https://tool-%s.ps-intern.de/api-gateway-legacy/gateway/entity", env)
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	dnsClient, err := dns.NewDNSClient(&api.OAuthConfig{
		ClientID:     d.Get("client_id").(string),
		ClientSecret: d.Get("client_secret").(string),
		Username:     d.Get("username").(string),
		Password:     d.Get("password").(string),
		TokenURL:     d.Get("token_url").(string),
	}, buildAPI(d))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create api client",
			Detail:   err.Error(),
		})

		return nil, diags
	}

	return dnsClient, diags
}
