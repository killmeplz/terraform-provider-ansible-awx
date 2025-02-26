package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceInventoryHost() *schema.Resource {
	return &schema.Resource{
		Create: resourceInventoryHostCreate,
		Read:   resourceInventoryHostRead,
		Update: resourceInventoryHostUpdate,
		Delete: resourceInventoryHostDelete,
		Description: "Manages a host within an Ansible AWX/Tower inventory. A host represents a managed node that " +
			"Ansible can configure and manage. Hosts can have variables specific to that host and can be enabled " +
			"or disabled to control whether they are available for running jobs.",

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of this host. This can be either a DNS name, IP address, or any other name used to identify the host.",
			},
			"inventory_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The ID of the inventory this host belongs to. Hosts must be associated with an inventory to be managed by AWX/Tower.",
				ValidateFunc: StringIsID,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional description of this host. Can be used to provide additional context about the host's purpose or configuration.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "If enabled (true), this host can be used in jobs. If disabled (false), this host will not be used in jobs even if included in the inventory.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The instance ID for this host if it is managed through a cloud provider. This helps track the host across IP or DNS changes.",
			},
			"variables": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Host variables in JSON or YAML format. These variables will be available to playbooks when running against this specific host and will override inventory variables.",
			},
		},
	}
}

func resourceInventoryHostCreate(d *schema.ResourceData, m interface{}) error {
	clientInstance := m.(*Client)
	data := map[string]interface{}{
		"name":        d.Get("name").(string),
		"description": d.Get("description").(string),
		"enabled":     d.Get("enabled"),
		"instance_id": d.Get("instance_id").(string),
		"variables":   d.Get("variables").(string),
	}

	resp, err := clientInstance.Post(fmt.Sprintf("/api/v2/inventories/%s/hosts", d.Get("inventory_id").(string)), data)
	if err != nil {
		return fmt.Errorf("failed to create AWX inventory host: %s", err)
	}

	id, ok := resp["id"].(float64)
	if !ok {
		return fmt.Errorf("AWX API did not return an id %v", resp)
	}
	d.SetId(fmt.Sprintf("%.0f", id))
	return resourceInventoryHostRead(d, m)
}

func resourceInventoryHostRead(d *schema.ResourceData, m interface{}) error {
	clientInstance := m.(*Client)
	id := d.Id()

	resp, err := clientInstance.Get(fmt.Sprintf("/api/v2/hosts/%s", id))
	if err != nil {
		if clientInstance.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed to read AWX inventory host: %s", err)
	}

	d.Set("name", resp["name"].(string))
	d.Set("description", resp["description"].(string))
	d.Set("inventory_id", F64ToStr(resp["inventory"]))
	d.Set("enabled", resp["enabled"])
	d.Set("instance_id", resp["instance_id"].(string))
	d.Set("variables", resp["variables"].(string))
	return nil
}

func resourceInventoryHostUpdate(d *schema.ResourceData, m interface{}) error {
	clientInstance := m.(*Client)
	id := d.Id()

	data := map[string]interface{}{}
	data["name"] = d.Get("name").(string)
	data["description"] = d.Get("description").(string)
	data["inventory"] = IfaceToInt(d.Get("inventory_id"))
	data["enabled"] = d.Get("enabled")
	data["instance_id"] = d.Get("instance_id").(string)
	data["variables"] = d.Get("variables").(string)

	_, err := clientInstance.Put(fmt.Sprintf("/api/v2/hosts/%s", id), data)
	if err != nil {
		return fmt.Errorf("failed to update AWX inventory host: %s, %v", err, data)
	}
	return resourceInventoryHostRead(d, m)
}

func resourceInventoryHostDelete(d *schema.ResourceData, m interface{}) error {
	clientInstance := m.(*Client)
	id := d.Id()

	err := clientInstance.Delete(fmt.Sprintf("/api/v2/hosts/%s", id))
	if err != nil {
		return fmt.Errorf("failed to delete AWX inventory host: %s", err)
	}
	d.SetId("")
	return nil
}
