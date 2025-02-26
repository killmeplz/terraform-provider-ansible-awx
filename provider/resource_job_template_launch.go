package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceJobTemplateLaunch() *schema.Resource {
	return &schema.Resource{
		Create: resourceJobTemplateLaunchCreateOrUpdate,
		Read:   resourceJobTemplateLaunchRead,
		Update: resourceJobTemplateLaunchCreateOrUpdate,
		Delete: resourceJobTemplateLaunchDelete,

		Schema: map[string]*schema.Schema{
			"job_template_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Optional description",
			},
			"inventory_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Inventory id",
			},
			"extra_vars": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
		},
	}
}

func resourceJobTemplateLaunchCreateOrUpdate(d *schema.ResourceData, m interface{}) error {
	clientInstance := m.(*Client)
	job_template_id := d.Get("job_template_id")

	data := map[string]interface{}{}
	data["job_template_id"] = IfaceToInt(d.Get("job_template_id"))
	data["extra_vars"] = d.Get("extra_vars").(string)
	if d.Get("inventory_id") != 0 {
		data["inventory_id"] = IfaceToInt(d.Get("inventory_id"))
	}

	resp, err := clientInstance.Post(fmt.Sprintf("/api/v2/job_templates/%d/launch", job_template_id), data)
	if err != nil {
		return fmt.Errorf("failed to create AWX job template launch: %s", err)
	}

	id, ok := resp["id"].(float64)
	if !ok {
		return fmt.Errorf("AWX API did not return an id %v", resp)
	}
	d.SetId(fmt.Sprintf("%.0f", id))
	return resourceJobTemplateLaunchRead(d, m)
}

func resourceJobTemplateLaunchRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceJobTemplateLaunchDelete(d *schema.ResourceData, m interface{}) error {
	d.SetId("")
	return nil
}
