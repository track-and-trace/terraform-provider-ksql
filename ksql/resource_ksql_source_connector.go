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
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the source connector",
			},
			"query": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The query",
			},
		},
	}
}

func sourceConnectorCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	query := d.Get("query").(string)
	log.Printf("[WARN] Creating a sourceConnector: %s with %s", name, query)
	c := meta.(*ksql.Client)
	q := fmt.Sprintf("CREATE SOURCE CONNECTOR %s %s", name, query)
	log.Printf("[WARN] Query %s", q)

	r := ksql.Request{
		KSQL: q,
	}
	resp, err := c.Do(r)
	log.Printf("[RESP] %v", resp)
	if err != nil {
		return err
	}
	d.SetId(name)
	return nil
}

func sourceConnectorRead(d *schema.ResourceData, meta interface{}) error {
	return fmt.Errorf("Source connector read is not implemented yet.")
	// TODO: Read the connectors
	//c := meta.(*ksql.Client)
	//name := d.Get("name").(string)
	//log.Printf("[ERROR] Searching for sourceConnector %s", name)
	// sourceConnectors, err := c.ListTables()
	// if err != nil {
	// 	return err
	// }
	// for _, t := range sourceConnectors {
	// 	//d.Set("query")
	// 	log.Printf("[INFO] Found %s: %v", t.Name, t)
	// }
	// return nil
}

func sourceConnectorDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ksql.Client)
	name := d.Get("name").(string)
	log.Printf("[INFO] Deleting sourceConnector %s", name)
	// TODO: Implement drop of source connector in ksql package
	// err := c.DropTable(&ksql.DropTableRequest{Name: name})
	_, err := c.Do(ksql.Request{KSQL: fmt.Sprintf("DROP CONNECTOR %s", name)})
	return err
}
