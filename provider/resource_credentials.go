package provider

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceCredentials() *schema.Resource {
	return &schema.Resource{
		Create: resourceCredentialsCreate,
		Read:   resourceCredentialsRead,
		Update: resourceCredentialsUpdate,
		Delete: resourceCredentialsDelete,

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
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Inherit permissions from organization roles. If provided on creation, do not give either user or team.",
				ValidateFunc: StringIsID,
			},
			"credential_type": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Specify the type of credential you want to create. Refer to the documentation for details on each type.",
				ValidateFunc: StringIsID,
			},
			"inputs": {
				Type:        schema.TypeMap,
				Required:    true,
				Sensitive:   true,
				Description: "Specify the type of credential you want to create. Refer to the documentation for details on each type.",
			},
		},
	}
}

func resourceCredentialsCreate(d *schema.ResourceData, m interface{}) error {
	clientInstance := m.(*Client)
	data := map[string]interface{}{
		"name":        d.Get("name").(string),
		"description": d.Get("description").(string),
		"inputs":      d.Get("inputs"),
	}
	if d.Get("organization") != "" {
		data["organization"] = IfaceToInt(d.Get("organization"))
	}
	data["credential_type"] = IfaceToInt(d.Get("credential_type"))

	resp, err := clientInstance.Post("/api/v2/credentials/", data)
	if err != nil {
		return fmt.Errorf("failed to create AWX resource: %s", err)
	}

	id, ok := resp["id"].(float64)
	if !ok {
		return fmt.Errorf("AWX API did not return an id %v", resp)
	}
	d.SetId(fmt.Sprintf("%.0f", id))
	return resourceCredentialsRead(d, m)
}

func resourceCredentialsRead(d *schema.ResourceData, m interface{}) error {
	clientInstance := m.(*Client)
	id := d.Id()

	resp, err := clientInstance.Get(fmt.Sprintf("/api/v2/credentials/%s/", id))
	if err != nil {
		if clientInstance.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed to read AWX resource: %s", err)
	}

	d.Set("name", resp["name"])
	d.Set("description", resp["description"])
	d.Set("organization", resp["organization"])
	return nil
}

func resourceCredentialsUpdate(d *schema.ResourceData, m interface{}) error {
	clientInstance := m.(*Client)
	id := d.Id()

	data := map[string]interface{}{
		"name":        d.Get("name").(string),
		"description": d.Get("description").(string),
		"inputs":      d.Get("inputs"),
	}
	if d.Get("organization") != "" {
		organization_id, _ := strconv.Atoi(d.Get("organization").(string))
		data["organization"] = organization_id
	}
	credential_type, _ := strconv.Atoi(d.Get("credential_type").(string))
	data["credential_type"] = credential_type

	_, err := clientInstance.Put(fmt.Sprintf("/api/v2/credentials/%s/", id), data)
	if err != nil {
		return fmt.Errorf("failed to update AWX resource: %s", err)
	}
	return resourceCredentialsRead(d, m)
}

func resourceCredentialsDelete(d *schema.ResourceData, m interface{}) error {
	clientInstance := m.(*Client)
	id := d.Id()

	err := clientInstance.Delete(fmt.Sprintf("/api/v2/credentials/%s/", id))
	if err != nil {
		return fmt.Errorf("failed to delete AWX resource: %s", err)
	}
	d.SetId("")
	return nil
}
