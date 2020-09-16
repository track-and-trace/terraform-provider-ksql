package ksql

import (
	"fmt"
	"log"

	"github.com/Mongey/ksql/ksql"
	"github.com/hashicorp/terraform/helper/schema"
)

func ksqlSourceConnectorResource() *schema.Resource {
	return &schema.Resource{
		Create: sourceConnectorCreate,
		Read:   sourceConnectorRead,
		Delete: sourceConnectorDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Description:   "The name of the sourceConnector",
				ConflictsWith: []string{"ksql"},
				Computed:      true,
			},
			"query": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Description:   "The query after CREATE SOURCE CONNECTOR[name]",
				ConflictsWith: []string{"ksql"},
				Computed:      true,
			},
			"ksql": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Description:   "The full query along with CREATE SOURCE CONNECTOR [name] infront",
				ConflictsWith: []string{"name", "query"},
				Computed:      true,
			},
		},
	}
}

func sourceConnectorCreate(d *schema.ResourceData, meta interface{}) error {
	return createKSQLResource(d, meta, "SOURCE CONNECTOR")
}

func sourceConnectorRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ksql.Client)
	name := d.Get("name").(string)
	log.Printf("[TRACE] Searching for sourceConnector %s", name)
	sourceConnectors, err := c.ListConnectors()
	if err != nil {
		return err
	}
	for _, s := range sourceConnectors {
		//d.Set("query")
		log.Printf("[TRACE] Found %s: %v", s.Name, s)
		if s.Name == name && s.Type == "source" {
			return nil
		}
	}
	return fmt.Errorf("source connector %s was not found", name)
}

func sourceConnectorDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ksql.Client)
	name := d.Get("name").(string)
	log.Printf("[INFO] Deleting sourceConnector %s", name)
	_, err := c.Do(ksql.Request{KSQL: fmt.Sprintf("DROP CONNECTOR %s;", name)})
	return err
}
