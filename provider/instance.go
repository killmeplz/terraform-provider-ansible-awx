package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/killmeplz/terraform-provider-ansible-awx/client"
)

// ResourceInstance defines the AWX instance resource.
func ResourceInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceInstanceCreate,
		Read:   resourceInstanceRead,
		Update: resourceInstanceUpdate,
		Delete: resourceInstanceDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A human-readable name for the AWX instance.",
			},
			"hostname": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The hostname of the instance.",
			},
			// Add additional fields based on AWX API documentation as needed.
		},
	}
}

func resourceInstanceCreate(d *schema.ResourceData, m interface{}) error {
	clientInstance := m.(*client.Client)
	instanceData := map[string]interface{}{
		"name":     d.Get("name").(string),
		"hostname": d.Get("hostname").(string),
		// Map additional fields as required.
	}

	// POST to /api/v2/instances/ to create a new instance.
	resp, err := clientInstance.Post("/api/v2/instances/", instanceData)
	if err != nil {
		return fmt.Errorf("failed to create AWX instance: %s", err)
	}

	// Assuming the API returns an "id" field for the created instance.
	id, ok := resp["id"].(string)
	if !ok || id == "" {
		return fmt.Errorf("AWX API did not return an instance id")
	}
	d.SetId(id)
	return resourceInstanceRead(d, m)
}

func resourceInstanceRead(d *schema.ResourceData, m interface{}) error {
	clientInstance := m.(*client.Client)
	id := d.Id()

	// GET from /api/v2/instances/{id}/ to fetch the instance details.
	resp, err := clientInstance.Get(fmt.Sprintf("/api/v2/instances/%s/", id))
	if err != nil {
		if clientInstance.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed to read AWX instance: %s", err)
	}

	// Update state with the response data.
	d.Set("name", resp["name"])
	d.Set("hostname", resp["hostname"])
	// Update additional fields as needed.
	return nil
}

func resourceInstanceUpdate(d *schema.ResourceData, m interface{}) error {
	clientInstance := m.(*client.Client)
	id := d.Id()

	updateData := map[string]interface{}{}
	if d.HasChange("name") {
		updateData["name"] = d.Get("name").(string)
	}
	if d.HasChange("hostname") {
		updateData["hostname"] = d.Get("hostname").(string)
	}
	// Process additional fields if they have changed.

	// PATCH to /api/v2/instances/{id}/ to update the instance.
	_, err := clientInstance.Patch(fmt.Sprintf("/api/v2/instances/%s/", id), updateData)
	if err != nil {
		return fmt.Errorf("failed to update AWX instance: %s", err)
	}
	return resourceInstanceRead(d, m)
}

func resourceInstanceDelete(d *schema.ResourceData, m interface{}) error {
	clientInstance := m.(*client.Client)
	id := d.Id()

	// DELETE the instance via /api/v2/instances/{id}/.
	err := clientInstance.Delete(fmt.Sprintf("/api/v2/instances/%s/", id))
	if err != nil {
		return fmt.Errorf("failed to delete AWX instance: %s", err)
	}
	d.SetId("")
	return nil
}
