package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/killmeplz/terraform-provider-ansible-awx/client"
)

// Provider returns a terraform.ResourceProvider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The base URL for the AWX instance.",
				DefaultFunc: schema.EnvDefaultFunc("AWX_HOST", nil),
			},
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The API token for authenticating with AWX.",
				DefaultFunc: schema.EnvDefaultFunc("AWX_TOKEN", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"ansible_awx_instance": ResourceInstance(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	host := d.Get("host").(string)
	token := d.Get("token").(string)

	clientInstance, err := client.NewClient(host, token)
	if err != nil {
		return nil, fmt.Errorf("failed to create AWX client: %s", err)
	}
	return clientInstance, nil
}
