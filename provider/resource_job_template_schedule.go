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
				Type:        schema.TypeString,
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
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Inventory id",
			},
			"project_id": {
				Type:        schema.TypeString,
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
	data := map[string]interface{}{
		"name":        d.Get("name").(string),
		"description": d.Get("description").(string),
		"job_type":    d.Get("job_type").(string),
		"project":     IfaceToInt(d.Get("project_id")),
		"playbook":    d.Get("playbook").(string),
		"scm_branch":  d.Get("scm_branch").(string),
		"forks":       d.Get("forks"),
		"limit":       d.Get("limit").(string),
		"verbosity":   d.Get("verbosity"),
		"extra_vars":  d.Get("extra_vars").(string),
		"job_tags":    d.Get("job_tags").(string),
		"rrule":       d.Get("rrule").(string),
	}
	if d.Get("inventory_id") != "" {
		data["inventory"] = IfaceToInt(d.Get("inventory_id"))
	}

	resp, err := clientInstance.Post(fmt.Sprintf("/api/v2/job_templates/%s/schedules", d.Get("job_template_id")), data)
	if err != nil {
		return fmt.Errorf("failed to create AWX job template schedule: %s", err)
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
		return fmt.Errorf("failed to read AWX job template schedule: %s", err)
	}

	d.Set("name", resp["name"].(string))
	d.Set("description", resp["description"].(string))
	d.Set("job_type", resp["job_type"].(string))
	d.Set("project", F64ToStr(resp["project"]))
	d.Set("playbook", resp["playbook"])
	d.Set("scm_branch", resp["scm_branch"].(string))
	d.Set("forks", resp["forks"])
	d.Set("limit", resp["limit"].(string))
	d.Set("verbosity", resp["verbosity"])
	d.Set("extra_vars", resp["extra_vars"])
	d.Set("job_tags", resp["job_tags"].(string))
	d.Set("rrule", resp["rrule"].(string))
	if resp["inventory"] != 0 {
		d.Set("inventory_id", F64ToStr(resp["inventory"]))
	}
	return nil
}

func resourceJobTemplateScheduleUpdate(d *schema.ResourceData, m interface{}) error {
	clientInstance := m.(*Client)
	id := d.Id()

	data := map[string]interface{}{}
	data["name"] = d.Get("name").(string)
	data["description"] = d.Get("description").(string)
	data["unified_job_template"] = d.Get("job_template_id")
	data["job_type"] = d.Get("job_type").(string)
	data["project"] = d.Get("project_id")
	data["playbook"] = d.Get("playbook").(string)
	data["scm_branch"] = d.Get("scm_branch").(string)
	data["forks"] = d.Get("forks")
	data["limit"] = d.Get("limit").(string)
	data["verbosity"] = d.Get("verbosity")
	data["extra_vars"] = d.Get("extra_vars").(string)
	data["job_tags"] = d.Get("job_tags").(string)
	data["rrule"] = d.Get("rrule").(string)
	if d.Get("inventory_id") != "" {
		data["inventory"] = IfaceToInt(d.Get("inventory_id"))
	}

	_, err := clientInstance.Put(fmt.Sprintf("/api/v2/schedules/%s/", id), data)
	if err != nil {
		return fmt.Errorf("failed to update AWX job template schedule: %s, %v", err, data)
	}
	return resourceJobTemplateScheduleRead(d, m)
}

func resourceJobTemplateScheduleDelete(d *schema.ResourceData, m interface{}) error {
	clientInstance := m.(*Client)
	id := d.Id()

	err := clientInstance.Delete(fmt.Sprintf("/api/v2/schedules/%s/", id))
	if err != nil {
		return fmt.Errorf("failed to delete AWX job template schedule: %s", err)
	}
	d.SetId("")
	return nil
}
