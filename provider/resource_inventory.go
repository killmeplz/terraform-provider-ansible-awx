package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceInventory() *schema.Resource {
	return &schema.Resource{
		Create: resourceInventoryCreate,
		Read:   resourceInventoryRead,
		Update: resourceInventoryUpdate,
		Delete: resourceInventoryDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of this credential.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional description of this credential.",
			},
			"organization": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Inherit permissions from organization roles. If provided on creation, do not give either user or team.",
			},
			"kind": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Kind of inventory being represented. (choice)",
			},
			"host_filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter that will be applied to the hosts of this inventory.",
			},
			"variables": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specify the type of credential you want to create. Refer to the documentation for details on each type.",
			},
			"prevent_instance_group_fallback": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If enabled, the inventory will prevent adding any organization instance groups to the list of preferred instances groups to run associated job templates on.If this setting is enabled and you provided an empty list, the global instance groups will be applied.",
			},
		},
	}
}

func resourceInventoryCreate(d *schema.ResourceData, m interface{}) error {
	clientInstance := m.(*Client)
	data := map[string]interface{}{
		"name":                            d.Get("name").(string),
		"description":                     d.Get("description").(string),
		"organization":                    d.Get("organization"),
		"kind":                            d.Get("kind").(string),
		"host_filter":                     d.Get("host_filter").(string),
		"variables":                       d.Get("variables").(string),
		"prevent_instance_group_fallback": d.Get("prevent_instance_group_fallback"),
	}

	resp, err := clientInstance.Post("/api/v2/inventories/", data)
	if err != nil {
		return fmt.Errorf("failed to create AWX instance: %s", err)
	}

	id, ok := resp["id"].(float64)
	if !ok {
		return fmt.Errorf("AWX API did not return an id %v", resp)
	}
	d.SetId(fmt.Sprintf("%.0f", id))
	return resourceInventoryRead(d, m)
}

func resourceInventoryRead(d *schema.ResourceData, m interface{}) error {
	clientInstance := m.(*Client)
	id := d.Id()

	resp, err := clientInstance.Get(fmt.Sprintf("/api/v2/inventories/%s/", id))
	if err != nil {
		if clientInstance.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed to read AWX instance: %s", err)
	}

	d.Set("name", resp["name"].(string))
	d.Set("description", resp["description"].(string))
	d.Set("organization", resp["organization"])
	d.Set("kind", resp["kind"].(string))
	d.Set("host_filter", resp["host_filter"])
	d.Set("variables", resp["variables"].(string))
	d.Set("prevent_instance_group_fallback", resp["prevent_instance_group_fallback"])
	return nil
}

func resourceInventoryUpdate(d *schema.ResourceData, m interface{}) error {
	clientInstance := m.(*Client)
	id := d.Id()

	updateData := map[string]interface{}{}
	updateData["name"] = d.Get("name").(string)
	updateData["description"] = d.Get("description").(string)
	updateData["organization"] = d.Get("organization")
	updateData["kind"] = d.Get("kind").(string)
	updateData["host_filter"] = d.Get("host_filter").(string)
	updateData["variables"] = d.Get("variables").(string)
	updateData["prevent_instance_group_fallback"] = d.Get("prevent_instance_group_fallback")

	_, err := clientInstance.Put(fmt.Sprintf("/api/v2/inventories/%s/", id), updateData)
	if err != nil {
		return fmt.Errorf("failed to update AWX instance: %s, %v", err, updateData)
	}
	return resourceInventoryRead(d, m)
}

func resourceInventoryDelete(d *schema.ResourceData, m interface{}) error {
	clientInstance := m.(*Client)
	id := d.Id()

	err := clientInstance.Delete(fmt.Sprintf("/api/v2/inventories/%s/", id))
	if err != nil {
		return fmt.Errorf("failed to delete AWX instance: %s", err)
	}
	d.SetId("")
	return nil
}
