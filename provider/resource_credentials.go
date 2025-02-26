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
		Description: "Manages credentials in Ansible AWX/Tower. Credentials are utilized by Tower for authentication " +
			"when launching jobs against machines, synchronizing with inventory sources, and importing project content from " +
			"version control systems. Different credential types support different authentication methods (SSH keys, " +
			"username/password, API tokens, etc.) depending on the service they connect to.",

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of this credential. Used to identify the credential in the AWX/Tower interface.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional description of this credential. Can be used to provide more context about the credential's purpose or usage.",
			},
			"organization": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The organization the credential belongs to. If provided, the credential will inherit permissions from organization roles. Cannot be specified together with user or team.",
				ValidateFunc: StringIsID,
			},
			"credential_type": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The type of credential being created (e.g., SSH, AWS, GitHub, etc.). This determines what authentication fields are required in the inputs parameter.",
				ValidateFunc: StringIsID,
			},
			"inputs": {
				Type:        schema.TypeMap,
				Required:    true,
				Sensitive:   true,
				Description: "A map of inputs required by the credential type. The specific inputs required depend on the credential_type. For example, an SSH credential might need 'username' and 'password' or 'ssh_key_data'.",
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
		return fmt.Errorf("failed to create AWX credentials: %s", err)
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
		return fmt.Errorf("failed to read AWX credentials: %s", err)
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
		return fmt.Errorf("failed to update AWX credentials: %s", err)
	}
	return resourceCredentialsRead(d, m)
}

func resourceCredentialsDelete(d *schema.ResourceData, m interface{}) error {
	clientInstance := m.(*Client)
	id := d.Id()

	err := clientInstance.Delete(fmt.Sprintf("/api/v2/credentials/%s/", id))
	if err != nil {
		return fmt.Errorf("failed to delete AWX credentials: %s", err)
	}
	d.SetId("")
	return nil
}
