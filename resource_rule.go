package main

import (
	"fmt"
	"github.com/atlassian/go-vtm"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func resourceRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceRuleCreate,
		Read:   resourceRuleRead,
		Update: resourceRuleUpdate,
		Delete: resourceRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"content": &schema.Schema{
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: compareWithNoteContent,
			},
			"note": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
		},
	}
}

func compareWithNoteContent(k, old, new string, d *schema.ResourceData) bool {
	newContent := fmt.Sprintf("#=-%v\n#=-\n", d.Get("note")) + string(new)
	log.Printf(old, newContent)
	return hashString(old) == hashString(newContent)
}

func resourceRuleCreate(d *schema.ResourceData, meta interface{}) error {
	err := resourceRuleSet(d, meta)
	if err != nil {
		return err
	}

	return resourceRuleRead(d, meta)
}

func resourceRuleRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*providerConfig).client

	r, resp, err := c.GetRule(d.Get("name").(string))
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			// The resource doesn't exist anymore
			d.SetId("")

			return nil
		}

		return fmt.Errorf("Error reading resource: %s", err)
	}

	d.Set("content", r.String())
	d.Set("note", r.GetNote())

	return nil
}

func resourceRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	err := resourceRuleSet(d, meta)
	if err != nil {
		return err
	}

	return resourceRuleRead(d, meta)
}

func resourceRuleDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*providerConfig).client
	r := stingray.NewRule(d.Id())

	_, err := c.Delete(r)
	if err != nil {
		return err
	}

	return nil
}

func resourceRuleSet(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*providerConfig).client
	r := stingray.NewRule(d.Get("name").(string))

	r.Content = []byte(d.Get("content").(string))
	r.Note = d.Get("note").(string)

	_, err := c.Set(r)
	if err != nil {
		return err
	}

	d.SetId(d.Get("name").(string))

	return nil
}
