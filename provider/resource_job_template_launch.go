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
		Description: "Launches an Ansible AWX/Tower job template. This resource allows you to execute job templates and " +
			"optionally override certain parameters such as inventory and variables. The job will be launched when this " +
			"resource is created or updated.",

		Schema: map[string]*schema.Schema{
			"job_template_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringIsID,
				Description:  "The ID of the job template to launch. This is a required field that specifies which AWX/Tower job template should be executed.",
			},
			"inventory_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringIsID,
				Description:  "The ID of the inventory to use for this job launch. If specified, this will override the inventory set in the job template.",
			},
			"extra_vars": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A JSON or YAML string containing extra variables to pass to the job template. These variables will be merged with any survey variables defined in the job template.",
			},
		},
	}
}

func resourceJobTemplateLaunchCreateOrUpdate(d *schema.ResourceData, m interface{}) error {
	clientInstance := m.(*Client)

	data := map[string]interface{}{}
	data["job_template_id"] = IfaceToInt(d.Get("job_template_id"))
	data["extra_vars"] = d.Get("extra_vars").(string)
	if d.Get("inventory_id") != "" {
		data["inventory_id"] = IfaceToInt(d.Get("inventory_id"))
	}

	resp, err := clientInstance.Post(fmt.Sprintf("/api/v2/job_templates/%s/launch", d.Get("job_template_id")), data)
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
