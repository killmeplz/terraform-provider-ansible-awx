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
		Description: "Manages an Ansible AWX/Tower inventory. An inventory is a collection of hosts against which jobs " +
			"may be launched, the same as an Ansible inventory file. Inventories are divided into groups and these " +
			"groups contain the actual hosts. Groups may be sourced manually, by entering host names into Tower, or " +
			"from one of Ansible Tower's supported cloud providers.",

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of this inventory. Used to identify the inventory in the AWX/Tower interface.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional description of this inventory. Can be used to provide more context about the inventory's purpose or contents.",
			},
			"organization": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringIsID,
				Description:  "The organization the inventory belongs to. Inventories must be associated with an organization for role-based access control.",
			},
			"kind": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The kind of inventory being represented. Choices include: '' (regular inventory), 'smart' (smart inventory), or 'constructed' (constructed inventory).",
			},
			"host_filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter that will be applied to the hosts of this inventory. Only used when kind=smart.",
			},
			"variables": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Inventory variables in JSON or YAML format. These variables will be available to all hosts in this inventory.",
			},
			"prevent_instance_group_fallback": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If enabled, the inventory will prevent falling back to instance groups defined at the organization or tower level. When disabled, the inventory will use instance groups from the organization or tower level if no inventory-specific instance groups are defined.",
			},
		},
	}
}

func resourceInventoryCreate(d *schema.ResourceData, m interface{}) error {
	clientInstance := m.(*Client)
	data := map[string]interface{}{
		"name":                            d.Get("name").(string),
		"description":                     d.Get("description").(string),
		"organization":                    IfaceToInt(d.Get("organization")),
		"kind":                            d.Get("kind").(string),
		"host_filter":                     d.Get("host_filter").(string),
		"variables":                       d.Get("variables").(string),
		"prevent_instance_group_fallback": d.Get("prevent_instance_group_fallback"),
	}

	resp, err := clientInstance.Post("/api/v2/inventories/", data)
	if err != nil {
		return fmt.Errorf("failed to create AWX inventory: %s", err)
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
		return fmt.Errorf("failed to read AWX inventory: %s", err)
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

	data := map[string]interface{}{}
	data["name"] = d.Get("name").(string)
	data["description"] = d.Get("description").(string)
	data["organization"] = IfaceToInt(d.Get("organization"))
	data["kind"] = d.Get("kind").(string)
	data["host_filter"] = d.Get("host_filter").(string)
	data["variables"] = d.Get("variables").(string)
	data["prevent_instance_group_fallback"] = d.Get("prevent_instance_group_fallback")

	_, err := clientInstance.Put(fmt.Sprintf("/api/v2/inventories/%s/", id), data)
	if err != nil {
		return fmt.Errorf("failed to update AWX inventory: %s, %v", err, data)
	}
	return resourceInventoryRead(d, m)
}

func resourceInventoryDelete(d *schema.ResourceData, m interface{}) error {
	clientInstance := m.(*Client)
	id := d.Id()

	err := clientInstance.Delete(fmt.Sprintf("/api/v2/inventories/%s/", id))
	if err != nil {
		return fmt.Errorf("failed to delete AWX inventory: %s", err)
	}
	d.SetId("")
	return nil
}
