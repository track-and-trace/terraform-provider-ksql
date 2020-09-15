package ksql

import (
	"log"

	"github.com/Mongey/ksql/ksql"
	"github.com/hashicorp/terraform/helper/schema"
)

func ksqlTableResource() *schema.Resource {
	return &schema.Resource{
		Create: tableCreate,
		Read:   tableRead,
		Delete: tableDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Description:   "The name of the table",
				ConflictsWith: []string{"ksql"},
				Computed:      true,
			},
			"query": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Description:   "The query after CREATE TABLE [name]",
				ConflictsWith: []string{"ksql"},
				Computed:      true,
			},
			"ksql": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Description:   "The full query along with CREATE TABLE [name] infront",
				ConflictsWith: []string{"name", "query"},
				Computed:      true,
			},
		},
	}
}

func tableCreate(d *schema.ResourceData, meta interface{}) error {
	return createKSQLResource(d, meta, "TABLE")
}

func tableRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ksql.Client)
	name := d.Get("name").(string)
	log.Printf("[ERROR] Searching for table %s", name)
	tables, err := c.ListTables()
	if err != nil {
		return err
	}
	for _, t := range tables {
		//d.Set("query")
		log.Printf("[INFO] Found %s: %v", t.Name, t)
	}
	return nil
}

func tableDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ksql.Client)
	name := d.Get("name").(string)
	log.Printf("[INFO] Deleting table %s", name)
	err := c.DropTable(&ksql.DropTableRequest{Name: name})
	return err
}
