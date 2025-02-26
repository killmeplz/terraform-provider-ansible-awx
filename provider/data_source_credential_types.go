package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCredentialTypes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCredentialTypesRead,

		Description: "Retrieves all available credential types from AWX/Tower. Credential types define the various " +
			"authentication mechanisms available for use in credentials. Each type specifies what information is required " +
			"(like username/password, SSH keys, API tokens, etc.) and how that information should be used for " +
			"authentication. This data source is useful when you need to reference credential type IDs in credential resources.",

		Schema: map[string]*schema.Schema{
			"types": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Map of credential type names to their corresponding IDs. Common types include 'Machine' for SSH credentials, " +
					"'Source Control' for VCS access, 'Amazon Web Services', 'OpenStack', 'VMware vCenter', etc.",
				Computed: true,
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
