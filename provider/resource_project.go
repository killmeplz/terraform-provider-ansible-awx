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
			"local_path": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Local path (relative to PROJECTS_ROOT) containing playbooks and related files for this project.",
			},
			"scm_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the source control system used to store the project.",
			},
			"scm_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The location where the project is stored.",
			},
			"scm_branch": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specific branch, tag or commit to checkout.",
			},
			"scm_refspec": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "For git projects, an additional refspec to fetch.",
			},
			"scm_clean": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Discard any local changes before syncing the project.",
			},
			"scm_track_submodules": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Track submodules latest commits on defined branch.",
			},
			"scm_delete_on_update": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Delete the project before syncing.",
			},
			"credential_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "",
				Description:  "SCM credential id",
				ValidateFunc: StringIsID,
			},
			"scm_update_on_launch": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Update the project when a job is launched that uses the project.",
			},
			"allow_override": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Allow changing the SCM branch or revision in a job template that uses this project.",
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
