package ksql

import (
	"fmt"
	"log"
	"strings"

	"github.com/Mongey/ksql/ksql"
	"github.com/hashicorp/terraform/helper/schema"
)

func createKSQLResource(d *schema.ResourceData, meta interface{}, resType string) error {
	nameVal, nameSpec := d.GetOk("name")
	queryVal, querySpec := d.GetOk("query")
	queryFullVal, ksqlSpec := d.GetOk("ksql")

	var name, query, queryFull string

	if !((nameSpec && querySpec) || ksqlSpec) {
		return fmt.Errorf("either (name and query) or ksql should be specified")
	}

	if ksqlSpec {
		// Parse the kSQL and set the name and query
		queryFull := queryFullVal.(string)
		splits := strings.SplitN(queryFull, " ", 4)
		if len(splits) != 4 {
			return fmt.Errorf("expected valid query but got %s", queryFull)
		}
		if strings.ToUpper(splits[1]) != resType {
			return fmt.Errorf("not a %s create statement %s", resType, queryFull)
		}
		name = splits[2]
		d.Set("name", name)
		query = splits[3]
		d.Set("query", query)

	}
	if nameSpec && querySpec {
		name = nameVal.(string)
		query = queryVal.(string)
		// compose queryFull
		queryFull = fmt.Sprintf("CREATE %s %s %s", resType, name, query)
		d.Set("ksql", queryFull)
	}

	log.Printf("[WARN] Creating a %s: %s with %s", resType, name, query)
	c := meta.(*ksql.Client)
	log.Printf("[WARN] Query %s", queryFull)

	r := ksql.Request{
		KSQL: queryFull,
	}

	resp, err := c.Do(r)
	log.Printf("[RESP] %v", resp)
	if err != nil {
		return err
	}
	d.SetId(name)
	return nil
}
