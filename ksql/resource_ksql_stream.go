package ksql

import (
	"fmt"
	"log"
	"strings"

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
	_, err := c.Do(ksql.Request{KSQL: fmt.Sprintf("DROP STREAM %s;", name)})
	if err != nil {
		// For streams we have the case with Queries termination
		// The following queries read from this source: [YYY].
		// The following queries write into this source: [XXX].
		// You need to terminate them before dropping ZZZ.
		errMsg := fmt.Sprintf("%v", err)
		if strings.Contains(errMsg, "terminate them before") {
			depQueries := getDepQueries(errMsg)
			for _, qn := range depQueries {
				_, err := c.Do(ksql.Request{KSQL: fmt.Sprintf("TERMINATE %s;", qn)})
				if err != nil {
					return err
				}
			}
			// Try again
			_, err := c.Do(ksql.Request{KSQL: fmt.Sprintf("DROP STREAM %s;", name)})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func getDepQueries(errMsg string) []string {
	queries := []string{}
	lines := strings.Split(errMsg, "\n")
	for _, l := range lines {
		qrs := strings.Split(l, "[")
		if len(qrs) > 1 {
			qrs = strings.Split(qrs[1], "]")
			if len(qrs) > 0 && len(qrs[0]) > 0 {
				queries = append(queries, strings.Split(qrs[0], ",")...)
			}
		}
	}

	return queries
}
