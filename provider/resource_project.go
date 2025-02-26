package provider

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectCreate,
		Read:   resourceProjectRead,
		Update: resourceProjectUpdate,
		Delete: resourceProjectDelete,
		Description: "Manages an Ansible AWX/Tower project. A project is a logical collection of Ansible playbooks, " +
			"represented in Tower. You can manage playbooks and playbook directories by either placing them manually " +
			"under the Project Base Path on your Tower server, or by placing your playbooks into a source code " +
			"management (SCM) system supported by Tower, including Git, Subversion, and Mercurial.",

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of this project. Used to identify the project in the AWX/Tower interface.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional description of this project. Can be used to provide more context about the project's purpose.",
			},
			"organization": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The organization the project belongs to. Projects must be associated with an organization for role-based access control.",
			},
			"local_path": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The local path (relative to PROJECTS_ROOT) on the AWX/Tower server where playbooks are stored. Used when scm_type is set to 'manual'.",
			},
			"scm_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Type of source control management system. Valid options include: 'manual', 'git', 'svn', 'hg', and others as supported by AWX/Tower.",
			},
			"scm_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The source control URL for the project. Required when scm_type is set to a valid SCM system.",
			},
			"scm_branch": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The branch, tag, or commit to checkout from the SCM system. Default is the default branch of the SCM repository.",
			},
			"scm_refspec": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "For git projects, an additional refspec to fetch. Can be used to retrieve additional branches or pull requests.",
			},
			"scm_clean": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If enabled, the project directory will be cleared before each update, removing any untracked files.",
			},
			"scm_track_submodules": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If enabled, git submodules will be tracked and updated when the project is updated.",
			},
			"scm_delete_on_update": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If enabled, the project directory will be deleted and recreated with each project update.",
			},
			"credential_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "",
				Description:  "The ID of the credential to use for authenticating with the SCM system.",
				ValidateFunc: StringIsID,
			},
			"scm_update_on_launch": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If enabled, the project will update from its SCM source before each job using this project is run.",
			},
			"allow_override": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If enabled, users can override the SCM branch or revision in job templates that use this project.",
			},
		},
	}
}

func resourceProjectCreate(d *schema.ResourceData, m interface{}) error {
	clientInstance := m.(*Client)
	data := map[string]interface{}{
		"name":                 d.Get("name").(string),
		"description":          d.Get("description").(string),
		"organization":         d.Get("organization"),
		"local_path":           d.Get("local_path").(string),
		"scm_type":             d.Get("scm_type").(string),
		"scm_url":              d.Get("scm_url").(string),
		"scm_branch":           d.Get("scm_branch").(string),
		"scm_refspec":          d.Get("scm_refspec").(string),
		"scm_clean":            d.Get("scm_clean"),
		"scm_track_submodules": d.Get("scm_track_submodules"),
		"scm_delete_on_update": d.Get("scm_delete_on_update"),
		"scm_update_on_launch": d.Get("scm_update_on_launch"),
		"allow_override":       d.Get("allow_override"),
	}
	if d.Get("credential_id") != "" {
		credential_id, _ := strconv.Atoi(d.Get("credential_id").(string))
		data["credential_id"] = credential_id
	}
	resp, err := clientInstance.Post("/api/v2/projects/", data)
	if err != nil {
		return fmt.Errorf("failed to create AWX project: %s", err)
	}

	id, ok := resp["id"].(float64)
	if !ok {
		return fmt.Errorf("AWX API did not return an id %v", resp)
	}
	d.SetId(fmt.Sprintf("%.0f", id))
	return resourceProjectRead(d, m)
}

func resourceProjectRead(d *schema.ResourceData, m interface{}) error {
	clientInstance := m.(*Client)
	id := d.Id()

	resp, err := clientInstance.Get(fmt.Sprintf("/api/v2/projects/%s/", id))
	if err != nil {
		if clientInstance.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed to read AWX project: %s", err)
	}

	d.Set("name", resp["name"].(string))
	d.Set("description", resp["description"].(string))
	d.Set("organization", resp["organization"])
	d.Set("scm_type", resp["scm_type"].(string))
	d.Set("scm_url", resp["scm_url"].(string))
	d.Set("scm_branch", resp["scm_branch"].(string))
	d.Set("scm_refspec", resp["scm_refspec"].(string))
	d.Set("scm_clean", resp["scm_clean"])
	d.Set("scm_track_submodules", resp["scm_track_submodules"])
	d.Set("scm_delete_on_update", resp["scm_delete_on_update"])
	d.Set("credential_id", resp["credential"])
	d.Set("scm_update_on_launch", resp["scm_update_on_launch"])
	d.Set("allow_override", resp["allow_override"])
	return nil
}

func resourceProjectUpdate(d *schema.ResourceData, m interface{}) error {
	clientInstance := m.(*Client)
	id := d.Id()

	updateData := map[string]interface{}{}
	updateData["name"] = d.Get("name").(string)
	updateData["description"] = d.Get("description").(string)
	updateData["organization"] = d.Get("organization")
	updateData["local_path"] = d.Get("local_path").(string)
	updateData["scm_type"] = d.Get("scm_type").(string)
	updateData["scm_url"] = d.Get("scm_url").(string)
	updateData["scm_branch"] = d.Get("scm_branch").(string)
	updateData["scm_refspec"] = d.Get("scm_refspec").(string)
	updateData["scm_clean"] = d.Get("scm_clean")
	updateData["scm_track_submodules"] = d.Get("scm_track_submodules")
	updateData["scm_delete_on_update"] = d.Get("scm_delete_on_update")
	updateData["credential_id"] = d.Get("credential_id")
	updateData["scm_update_on_launch"] = d.Get("scm_update_on_launch")
	updateData["allow_override"] = d.Get("allow_override")

	_, err := clientInstance.Put(fmt.Sprintf("/api/v2/projects/%s/", id), updateData)
	if err != nil {
		return fmt.Errorf("failed to update AWX project: %s, %v", err, updateData)
	}
	return resourceProjectRead(d, m)
}

func resourceProjectDelete(d *schema.ResourceData, m interface{}) error {
	clientInstance := m.(*Client)
	id := d.Id()

	err := clientInstance.Delete(fmt.Sprintf("/api/v2/projects/%s/", id))
	if err != nil {
		return fmt.Errorf("failed to delete AWX project: %s", err)
	}
	d.SetId("")
	return nil
}
