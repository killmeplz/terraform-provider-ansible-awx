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

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of this credential.",
			},
			"inventory_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The value used by the remote inventory source to uniquely identify the host.",
				ValidateFunc: StringIsID,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional description of this credential.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Is this host online and available for running jobs? ",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The value used by the remote inventory source to uniquely identify the host.",
			},
			"variables": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specify the type of credential you want to create. Refer to the documentation for details on each type.",
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
		return fmt.Errorf("failed to create AWX resource: %s", err)
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
		return fmt.Errorf("failed to read AWX resource: %s", err)
	}

	d.Set("name", resp["name"].(string))
	d.Set("description", resp["description"].(string))
	d.Set("inventory_id", fmt.Sprintf("%.0f", resp["inventory"].(float64)))
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
		return fmt.Errorf("failed to update AWX resource: %s, %v", err, data)
	}
	return resourceInventoryHostRead(d, m)
}

func resourceInventoryHostDelete(d *schema.ResourceData, m interface{}) error {
	clientInstance := m.(*Client)
	id := d.Id()

	err := clientInstance.Delete(fmt.Sprintf("/api/v2/hosts/%s", id))
	if err != nil {
		return fmt.Errorf("failed to delete AWX resource: %s", err)
	}
	d.SetId("")
	return nil
}
