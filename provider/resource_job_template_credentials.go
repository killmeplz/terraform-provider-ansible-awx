package provider

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceJobTemplateCredential() *schema.Resource {
	return &schema.Resource{
		Create: resourceJobTemplateCredentialCreate,
		Read:   resourceJobTemplateCredentialRead,
		Update: resourceJobTemplateCredentialUpdate,
		Delete: resourceJobTemplateCredentialDelete,

		Schema: map[string]*schema.Schema{
			"job_template_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Optional description",
			},
			"credentials_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "run",
				Description: "Choose between run and check.",
			},
		},
	}
}

func resourceJobTemplateCredentialCreate(d *schema.ResourceData, m interface{}) error {
	clientInstance := m.(*Client)
	data := map[string]interface{}{
		"id": IfaceToInt(d.Get("credentials_id")),
	}

	_, err := clientInstance.Post(fmt.Sprintf("/api/v2/job_templates/%s/credentials", d.Get("job_template_id")), data)
	if err != nil {
		return fmt.Errorf("failed to associate credentials with AWX job template: %s", err)
	}

	d.SetId(fmt.Sprintf("%s_%s", d.Get("job_template_id"), d.Get("credentials_id")))
	return resourceJobTemplateCredentialRead(d, m)
}

func resourceJobTemplateCredentialRead(d *schema.ResourceData, m interface{}) error {
	clientInstance := m.(*Client)
	ids := strings.Split(d.Id(), "_")

	resp, err := clientInstance.Get(fmt.Sprintf("/api/v2/job_templates/%s/credentials", d.Get("job_template_id")))
	if err != nil {
		if clientInstance.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed to read AWX job template schedule: %s", err)
	}

	res := false
	if len(resp["results"].([]interface{})) == 0 {
		d.SetId("")
		return nil
	}
	for _, result := range resp["results"].([]interface{}) {
		if fmt.Sprintf("%0.f", result.(map[string]interface{})["id"]) == ids[1] {
			res = true
		}
	}
	if !res {
		d.SetId("")
		return nil
	}
	return nil
}

func resourceJobTemplateCredentialUpdate(d *schema.ResourceData, m interface{}) error {
	clientInstance := m.(*Client)
	data := map[string]interface{}{
		"id": IfaceToInt(d.Get("credentials_id")),
	}

	_, err := clientInstance.Post(fmt.Sprintf("/api/v2/job_templates/%s/credentials", d.Get("job_template_id")), data)
	if err != nil {
		return fmt.Errorf("failed to associate credentials with AWX job template: %s", err)
	}

	d.SetId(fmt.Sprintf("%s_%s", d.Get("job_template_id"), d.Get("credentials_id")))
	return resourceJobTemplateCredentialRead(d, m)
}

func resourceJobTemplateCredentialDelete(d *schema.ResourceData, m interface{}) error {
	clientInstance := m.(*Client)
	data := map[string]interface{}{
		"id":           IfaceToInt(d.Get("credentials_id")),
		"disassociate": true,
	}

	_, err := clientInstance.Post(fmt.Sprintf("/api/v2/job_templates/%s/credentials", d.Get("job_template_id")), data)
	if err != nil {
		return fmt.Errorf("failed to dicassociate credentials with AWX job template: %s", err)
	}
	d.SetId("")
	return nil
}
