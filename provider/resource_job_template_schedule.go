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
		Description: "Manages a schedule for an Ansible AWX/Tower job template. This resource allows you to create, " +
			"update, and delete scheduled runs of job templates. You can configure various parameters including the " +
			"execution schedule (using RRULE format), playbook options, and variables.",

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of this schedule. Used to identify the schedule in the AWX/Tower interface.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional description of the schedule. Can be used to provide more context about the schedule's purpose.",
			},
			"job_template_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the job template this schedule is associated with. This defines which job template will be executed on the schedule.",
			},
			"job_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "run",
				Description: "The type of job run. Can be either 'run' for normal execution or 'check' for check mode.",
			},
			"inventory_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the inventory to use for this scheduled job. If specified, this will override the inventory set in the job template.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the project containing the playbook to execute. If specified, this will override the project set in the job template.",
			},
			"playbook": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the playbook to execute. If specified, this will override the playbook set in the job template.",
			},
			"scm_branch": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specific branch, tag or commit to checkout from SCM before running the playbook.",
			},
			"forks": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Number of parallel processes to use while executing the playbook. Default of 0 uses the ansible default.",
			},
			"limit": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Limit the execution to specific hosts or groups. Corresponds to ansible's --limit parameter.",
			},
			"verbosity": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Control the level of output ansible will produce during execution. Higher numbers mean more output.",
			},
			"extra_vars": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A JSON or YAML string containing extra variables to pass to the playbook. These variables will be merged with any survey variables.",
			},
			"job_tags": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specify which tagged tasks from the playbook to execute. Only tasks with specified tags will be run.",
			},
			"rrule": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A recurrence rule (RRULE) string that defines when the schedule will run. Uses the iCal RRULE format (e.g., FREQ=DAILY;INTERVAL=1).",
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
	d.Set("playbook", resp["playbook"])
	d.Set("scm_branch", resp["scm_branch"].(string))
	d.Set("forks", resp["forks"])
	d.Set("limit", resp["limit"].(string))
	d.Set("verbosity", resp["verbosity"])
	d.Set("extra_vars", resp["extra_vars"])
	d.Set("job_tags", resp["job_tags"].(string))
	d.Set("rrule", resp["rrule"].(string))
	if resp["inventory"] != nil {
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
