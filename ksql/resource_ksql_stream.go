package ksql

import (
	"fmt"
	"log"

	"github.com/Mongey/ksql/ksql"
	"github.com/hashicorp/terraform/helper/schema"
)

func ksqlStreamResource() *schema.Resource {
	return &schema.Resource{
		Create: streamCreate,
		Read:   streamRead,
		Delete: streamDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Description:   "The name of the stream",
				ConflictsWith: []string{"ksql"},
				Computed:      true,
			},
			"query": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Description:   "The query after CREATE STREAM [name]",
				ConflictsWith: []string{"ksql"},
				Computed:      true,
			},
			"ksql": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Description:   "The full query along with CREATE STREAM [name] infront",
				ConflictsWith: []string{"name", "query"},
				Computed:      true,
			},
		},
	}
}

func streamCreate(d *schema.ResourceData, meta interface{}) error {
	return createKSQLResource(d, meta, "STREAM")
}

func streamRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ksql.Client)
	name := d.Get("name").(string)
	log.Printf("[ERROR] Searching for stream %s", name)
	streams, err := c.ListStreams()
	if err != nil {
		return err
	}
	for _, s := range streams {
		//d.Set("query")
		log.Printf("[INFO] Found %s: %v", s.Name, s)
		if s.Name == name {
			return nil
		}
	}
	return fmt.Errorf("did not found stream %s", name)
}

func streamDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ksql.Client)
	name := d.Get("name").(string)
	log.Printf("[INFO] Deleting stream %s", name)
	return dropObjectWithTerminateDeps(c, fmt.Sprintf("STREAM %s", name))
}
