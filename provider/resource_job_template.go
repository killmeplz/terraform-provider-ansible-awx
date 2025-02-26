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
		Description: "Manages an Ansible AWX/Tower job template. A job template is a definition and set of parameters for running " +
			"an Ansible job. Job templates are useful to execute the same job many times. Job templates can contain specifications " +
			"for: the inventory to run the job against, the project and playbook to use, credentials, extra variables, and various " +
			"other parameters.",

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of this job template. Used to identify the template in the AWX/Tower interface.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional description of the job template. Can be used to provide more context about the template's purpose.",
			},
			"job_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "run",
				Description: "The type of job to run. Can be either 'run' for normal execution or 'check' for check mode (dry run).",
			},
			"inventory_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The ID of the inventory to be used by this job template. Defines which hosts the playbook will be run against.",
				ValidateFunc: StringIsID,
			},
			"project_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The ID of the project containing the playbook to be used by this job template.",
				ValidateFunc: StringIsID,
			},
			"playbook": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the playbook to be run. The playbook must exist in the project specified by project_id.",
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
				Description: "Control the level of output Ansible will produce during execution. Higher numbers mean more verbose output (0-4).",
			},
			"extra_vars": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A JSON or YAML string containing extra variables to pass to the playbook. These variables will be available to the playbook and any surveys.",
			},
			"job_tags": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specify which tagged tasks from the playbook to execute. Only tasks with the specified tags will be run.",
			},
			"ask_inventory_on_launch": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If enabled, users will be prompted to select an inventory when the job template is launched.",
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
		return fmt.Errorf("failed to create AWX job template: %s", err)
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
		return fmt.Errorf("failed to read AWX job template: %s", err)
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
		return fmt.Errorf("failed to update AWX job template: %s, %v", err, data)
	}
	return resourceJobTemplateRead(d, m)
}

func resourceJobTemplateDelete(d *schema.ResourceData, m interface{}) error {
	clientInstance := m.(*Client)
	id := d.Id()

	err := clientInstance.Delete(fmt.Sprintf("/api/v2/job_templates/%s/", id))
	if err != nil {
		return fmt.Errorf("failed to delete AWX job template: %s", err)
	}
	d.SetId("")
	return nil
}
