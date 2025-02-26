package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func New() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The base URL for the AWX/Tower instance (e.g., https://awx.example.com). This URL will be used for all API calls.",
				DefaultFunc: schema.EnvDefaultFunc("AWX_HOST", nil),
			},
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The OAuth2 token or Personal Access Token for authenticating with AWX/Tower. This token must have sufficient permissions to perform the requested operations.",
				DefaultFunc: schema.EnvDefaultFunc("AWX_TOKEN", nil),
				Sensitive:   true,
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"awx_credential_types": dataSourceCredentialTypes(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"awx_credentials":              ResourceCredentials(),
			"awx_inventory":                ResourceInventory(),
			"awx_inventory_host":           ResourceInventoryHost(),
			"awx_project":                  ResourceProject(),
			"awx_job_template":             ResourceJobTemplate(),
			"awx_job_template_schedule":    ResourceJobTemplateSchedule(),
			"awx_job_template_launch":      ResourceJobTemplateLaunch(),
			"awx_job_template_credentials": ResourceJobTemplateCredential(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	host := d.Get("host").(string)
	token := d.Get("token").(string)

	clientInstance, err := NewClient(host, token)
	if err != nil {
		return nil, fmt.Errorf("failed to create AWX  %s", err)
	}
	return clientInstance, nil
}
