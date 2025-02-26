package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceJobTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceJobTemplateCreate,
		Read:   resourceJobTemplateRead,
		Update: resourceJobTemplateUpdate,
		Delete: resourceJobTemplateDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of this job template.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional description",
			},
			"job_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "run",
				Description: "Choose between run and check.",
			},
			"inventory_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Inventory id",
				ValidateFunc: StringIsID,
			},
			"project_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Project ID",
				ValidateFunc: StringIsID,
			},
			"playbook": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Playbook to use",
			},
			"scm_branch": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specific branch, tag or commit to checkout.",
			},
			"forks": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "",
			},
			"limit": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "",
			},
			"verbosity": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Output verbosity",
			},
			"extra_vars": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
			"job_tags": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
			"ask_inventory_on_launch": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "",
			},
		},
	}
}

func resourceJobTemplateCreate(d *schema.ResourceData, m interface{}) error {
	clientInstance := m.(*Client)
	data := map[string]interface{}{}
	data["name"] = d.Get("name").(string)
	data["description"] = d.Get("description").(string)
	data["job_type"] = d.Get("job_type").(string)
	data["playbook"] = d.Get("playbook").(string)
	data["scm_branch"] = d.Get("scm_branch").(string)
	data["forks"] = d.Get("forks")
	data["limit"] = d.Get("limit").(string)
	data["verbosity"] = d.Get("verbosity")
	data["extra_vars"] = d.Get("extra_vars").(string)
	data["job_tags"] = d.Get("job_tags").(string)
	data["ask_inventory_on_launch"] = d.Get("ask_inventory_on_launch")
	data["inventory"] = IfaceToInt(d.Get("inventory_id"))
	data["project"] = IfaceToInt(d.Get("project_id"))

	resp, err := clientInstance.Post("/api/v2/job_templates/", data)
	if err != nil {
		return fmt.Errorf("failed to create AWX instance: %s", err)
	}

	id, ok := resp["id"].(float64)
	if !ok {
		return fmt.Errorf("AWX API did not return an id %v", resp)
	}
	d.SetId(fmt.Sprintf("%.0f", id))
	return resourceJobTemplateRead(d, m)
}

func resourceJobTemplateRead(d *schema.ResourceData, m interface{}) error {
	clientInstance := m.(*Client)
	id := d.Id()

	resp, err := clientInstance.Get(fmt.Sprintf("/api/v2/job_templates/%s/", id))
	if err != nil {
		if clientInstance.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed to read AWX instance: %s", err)
	}

	d.Set("name", resp["name"].(string))
	d.Set("description", resp["description"].(string))
	d.Set("job_type", resp["job_type"].(string))
	d.Set("inventory_id", resp["inventory"])
	d.Set("project_id", resp["project"])
	d.Set("playbook", resp["playbook"].(string))
	d.Set("scm_branch", resp["scm_branch"].(string))
	d.Set("forks", resp["forks"])
	d.Set("limit", resp["limit"].(string))
	d.Set("verbosity", resp["verbosity"])
	d.Set("extra_vars", resp["extra_vars"].(string))
	d.Set("job_tags", resp["job_tags"].(string))
	d.Set("ask_inventory_on_launch", resp["ask_inventory_on_launch"])
	return nil
}

func resourceJobTemplateUpdate(d *schema.ResourceData, m interface{}) error {
	clientInstance := m.(*Client)
	id := d.Id()

	data := map[string]interface{}{}
	data["name"] = d.Get("name").(string)
	data["description"] = d.Get("description").(string)
	data["job_type"] = d.Get("job_type").(string)
	data["inventory"] = IfaceToInt(d.Get("inventory_id"))
	data["project"] = IfaceToInt(d.Get("project_id"))
	data["playbook"] = d.Get("playbook").(string)
	data["scm_branch"] = d.Get("scm_branch").(string)
	data["forks"] = d.Get("forks")
	data["limit"] = d.Get("limit").(string)
	data["verbosity"] = d.Get("verbosity")
	data["extra_vars"] = d.Get("extra_vars").(string)
	data["job_tags"] = d.Get("job_tags").(string)
	data["ask_inventory_on_launch"] = d.Get("ask_inventory_on_launch")

	_, err := clientInstance.Put(fmt.Sprintf("/api/v2/job_templates/%s/", id), data)
	if err != nil {
		return fmt.Errorf("failed to update AWX instance: %s, %v", err, data)
	}
	return resourceJobTemplateRead(d, m)
}

func resourceJobTemplateDelete(d *schema.ResourceData, m interface{}) error {
	clientInstance := m.(*Client)
	id := d.Id()

	err := clientInstance.Delete(fmt.Sprintf("/api/v2/job_templates/%s/", id))
	if err != nil {
		return fmt.Errorf("failed to delete AWX instance: %s", err)
	}
	d.SetId("")
	return nil
}
