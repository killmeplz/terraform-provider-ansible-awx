package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCredentialTypes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCredentialTypesRead,

		Description: "Get all credentials types from AWX API",

		Schema: map[string]*schema.Schema{
			"types": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "List of all available credential types",
				Computed:    true,
			},
		},
	}
}

func dataSourceCredentialTypesRead(d *schema.ResourceData, m interface{}) error {
	clientInstance := m.(*Client)
	resp, err := clientInstance.Get("/api/v2/credential_types?page_size=100")
	if err != nil {
		if clientInstance.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed to read AWX instance: %s", err)
	}

	mapTypes := make(map[string]int)
	for _, credentialType := range resp["results"].([]interface{}) {
		ct := credentialType.(map[string]interface{})
		mapTypes[ct["name"].(string)] = int(ct["id"].(float64))
	}
	d.SetId("1")
	d.Set("types", mapTypes)

	return nil
}
