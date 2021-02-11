package ksql

import (
	"fmt"
	"strings"

	"github.com/Mongey/ksql/ksql"
)

func dropObjectWithTerminateDeps(c *ksql.Client, typeAndNameame string) error {
	dropObjKSQL := fmt.Sprintf("DROP %s;", typeAndNameame)
	_, err := c.Do(ksql.Request{KSQL: dropObjKSQL})
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
			_, err := c.Do(ksql.Request{KSQL: dropObjKSQL})
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
