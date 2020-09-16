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
		nSplits := 4
		if len(strings.Split(resType, " ")) > 1 {
			nSplits = 5
		}

		// Parse the kSQL and set the name and query
		queryFull = queryFullVal.(string)
		splits := strings.SplitN(queryFull, " ", nSplits)
		if len(splits) != nSplits {
			return fmt.Errorf("expected valid query but got %s", queryFull)
		}
		ksqlType := splits[1]
		if nSplits == 5 {
			ksqlType = splits[1] + " " + splits[2]
		}
		if strings.ToUpper(ksqlType) != resType {
			return fmt.Errorf("not a %s create statement %s", resType, queryFull)
		}
		name = splits[nSplits-2]
		d.Set("name", name)
		query = splits[nSplits-1]
		d.Set("query", query)

	}
	if nameSpec && querySpec {
		name = nameVal.(string)
		query = queryVal.(string)
		// compose queryFull
		queryFull = fmt.Sprintf("CREATE %s %s %s", resType, name, query)
		d.Set("ksql", queryFull)
	}

	log.Printf("[DEBUG] Creating a %s: %s with %s", resType, name, query)
	c := meta.(*ksql.Client)
	log.Printf("[DEBUG] Query %s", queryFull)

	r := ksql.Request{
		KSQL: queryFull,
	}

	resp, err := c.Do(r)
	log.Printf("[DEBUG] %v", resp)
	if err != nil {
		log.Printf("[ERROR] filed to execute query : %s : %v", queryFull, err)
		return err
	}
	d.SetId(name)
	return nil
}
