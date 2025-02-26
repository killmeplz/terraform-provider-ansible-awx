package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceJobTemplateSchedule() *schema.Resource {
	return &schema.Resource{
		Create: resourceJobTemplateScheduleCreate,
		Read:   resourceJobTemplateScheduleRead,
		Update: resourceJobTemplateScheduleUpdate,
		Delete: resourceJobTemplateScheduleDelete,

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
			"job_template_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Optional description",
			},
			"job_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "run",
				Description: "Choose between run and check.",
			},
			"inventory_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Inventory id",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Project ID",
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
			"rrule": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
		},
	}
}

func resourceJobTemplateScheduleCreate(d *schema.ResourceData, m interface{}) error {
	clientInstance := m.(*Client)
	job_template_id := d.Get("job_template_id")
	data := map[string]interface{}{
		"name":        d.Get("name").(string),
		"description": d.Get("description").(string),
		"job_type":    d.Get("job_type").(string),
		"project":     d.Get("project_id"),
		"playbook":    d.Get("playbook").(string),
		"scm_branch":  d.Get("scm_branch").(string),
		"forks":       d.Get("forks"),
		"limit":       d.Get("limit").(string),
		"verbosity":   d.Get("verbosity"),
		"extra_vars":  d.Get("extra_vars").(string),
		"job_tags":    d.Get("job_tags").(string),
		"rrule":       d.Get("rrule").(string),
	}

	resp, err := clientInstance.Post(fmt.Sprintf("/api/v2/job_templates/%d/schedules", job_template_id), data)
	if err != nil {
		return fmt.Errorf("failed to create AWX instance: %s", err)
	}

	id, ok := resp["id"].(float64)
	if !ok {
		return fmt.Errorf("AWX API did not return an id %v", resp)
	}
	d.SetId(fmt.Sprintf("%.0f", id))
	return resourceJobTemplateScheduleRead(d, m)
}

func resourceJobTemplateScheduleRead(d *schema.ResourceData, m interface{}) error {
	clientInstance := m.(*Client)
	id := d.Id()

	resp, err := clientInstance.Get(fmt.Sprintf("/api/v2/schedules/%s/", id))
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
	d.Set("project", resp["project"])
	d.Set("playbook", resp["playbook"])
	d.Set("scm_branch", resp["scm_branch"].(string))
	d.Set("forks", resp["forks"])
	d.Set("limit", resp["limit"].(string))
	d.Set("verbosity", resp["verbosity"])
	d.Set("extra_vars", resp["extra_vars"])
	d.Set("job_tags", resp["job_tags"].(string))
	d.Set("rrule", resp["rrule"].(string))
	return nil
}

func resourceJobTemplateScheduleUpdate(d *schema.ResourceData, m interface{}) error {
	clientInstance := m.(*Client)
	id := d.Id()

	updateData := map[string]interface{}{}
	updateData["name"] = d.Get("name").(string)
	updateData["description"] = d.Get("description").(string)
	updateData["unified_job_template"] = d.Get("job_template_id")
	updateData["job_type"] = d.Get("job_type").(string)
	updateData["project"] = d.Get("project_id")
	updateData["playbook"] = d.Get("playbook").(string)
	updateData["scm_branch"] = d.Get("scm_branch").(string)
	updateData["forks"] = d.Get("forks")
	updateData["limit"] = d.Get("limit").(string)
	updateData["verbosity"] = d.Get("verbosity")
	updateData["extra_vars"] = d.Get("extra_vars").(string)
	updateData["job_tags"] = d.Get("job_tags").(string)
	updateData["rrule"] = d.Get("rrule").(string)

	_, err := clientInstance.Put(fmt.Sprintf("/api/v2/schedules/%s/", id), updateData)
	if err != nil {
		return fmt.Errorf("failed to update AWX instance: %s, %v", err, updateData)
	}
	return resourceJobTemplateScheduleRead(d, m)
}

func resourceJobTemplateScheduleDelete(d *schema.ResourceData, m interface{}) error {
	clientInstance := m.(*Client)
	id := d.Id()

	err := clientInstance.Delete(fmt.Sprintf("/api/v2/schedules/%s/", id))
	if err != nil {
		return fmt.Errorf("failed to delete AWX instance: %s", err)
	}
	d.SetId("")
	return nil
}
