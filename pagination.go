package dglite

import (
	"strconv"

	"mooncamp.com/dgx/gql"
)

type pagination struct{}

func (pg *pagination) paginateNodes(q gql.GraphQuery, nodes []map[string]interface{}) []map[string]interface{} {
	var res []map[string]interface{}
	for _, e := range nodes {
		res = append(res, pg.paginate(q, e))
	}
	return res
}

func (pg *pagination) paginate(q gql.GraphQuery, node map[string]interface{}) map[string]interface{} {
	res := map[string]interface{}{}
	for k, v := range node {
		switch actual := v.(type) {
		case map[string]interface{}:
			res[k] = pg.paginate(findQuery(q.Children, k), actual)
		case []interface{}:
			if len(actual) == 0 {
				res[k] = v
			}

			if _, ok := actual[0].(map[string]interface{}); !ok {
				res[k] = v
				continue
			}

			nq := findQuery(q.Children, k)
			var ns []interface{}
			for _, e := range actual {
				ns = append(ns, pg.paginate(nq, e.(map[string]interface{})))
			}

			if len(nq.Args) == 0 {
				res[k] = ns
				continue
			}

			if c, ok := nq.Args["offset"]; ok {
				i, err := strconv.Atoi(c)
				if err == nil {
					ns = ns[i:]
				}
			}

			if c, ok := nq.Args["first"]; ok {
				i, err := strconv.Atoi(c)
				if err == nil {
					ns = ns[:i]
				}
			}

			res[k] = ns
		default:
			res[k] = v
		}
	}
	return res
}
