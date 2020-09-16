package ksql

import (
	"fmt"
	"log"

	"github.com/Mongey/ksql/ksql"
	"github.com/hashicorp/terraform/helper/schema"
)

func ksqlSinkConnectorResource() *schema.Resource {
	return &schema.Resource{
		Create: sinkConnectorCreate,
		Read:   sinkConnectorRead,
		Delete: sinkConnectorDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Description:   "The name of the sinkConnector",
				ConflictsWith: []string{"ksql"},
				Computed:      true,
			},
			"query": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Description:   "The query after CREATE SINK CONNECTOR[name]",
				ConflictsWith: []string{"ksql"},
				Computed:      true,
			},
			"ksql": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Description:   "The full query along with CREATE SINK CONNECTOR [name] infront",
				ConflictsWith: []string{"name", "query"},
				Computed:      true,
			},
		},
	}
}

func sinkConnectorCreate(d *schema.ResourceData, meta interface{}) error {
	return createKSQLResource(d, meta, "SINK CONNECTOR")
}

func sinkConnectorRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ksql.Client)
	name := d.Get("name").(string)
	log.Printf("[TRACE] Searching for sinkConnector %s", name)
	sinkConnectors, err := c.ListConnectors()
	if err != nil {
		return err
	}
	for _, s := range sinkConnectors {
		//d.Set("query")
		log.Printf("[TRACE] Found %s: %v", s.Name, s)
		if s.Name == name && s.Type == "sink" {
			return nil
		}
	}
	return fmt.Errorf("sink connector %s was not found", name)
}

func sinkConnectorDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ksql.Client)
	name := d.Get("name").(string)
	log.Printf("[INFO] Deleting sinkConnector %s", name)
	_, err := c.Do(ksql.Request{KSQL: fmt.Sprintf("DROP CONNECTOR %s;", name)})
	return err
}
