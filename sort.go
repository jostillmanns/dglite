package dglite

import (
	"sort"

	"mooncamp.com/dgx/gql"
)

type sorter struct{}

func (st *sorter) sortNodes(q gql.GraphQuery, nodes []map[string]interface{}) []map[string]interface{} {
	var res []map[string]interface{}
	for _, e := range nodes {
		res = append(res, st.sort(q, e))
	}
	return res
}

func (st *sorter) sort(q gql.GraphQuery, node map[string]interface{}) map[string]interface{} {
	res := map[string]interface{}{}
	for k, v := range node {
		switch actual := v.(type) {
		case map[string]interface{}:
			res[k] = st.sort(findQuery(q.Children, k), actual)
		case []interface{}:
			if len(actual) == 0 {
				res[k] = v
				continue
			}

			if _, ok := actual[0].(map[string]interface{}); !ok {
				res[k] = v
				continue
			}

			nq := findQuery(q.Children, k)
			var ns []interface{}
			for _, e := range actual {
				ns = append(ns, st.sort(nq, e.(map[string]interface{})))
			}

			if nq.Order == nil {
				res[k] = ns
				continue
			}

			sn := &sortedNodes{nodes: ns, predicate: nq.Order[0].Attr}
			sort.Sort(sn)

			res[k] = sn.nodes

		default:
			res[k] = v
		}
	}
	return res
}

func findQuery(qs []gql.GraphQuery, predicate string) gql.GraphQuery {
	for _, e := range qs {
		if e.Attr == predicate {
			return e
		}
	}

	return gql.GraphQuery{}
}

type sortedNodes struct {
	nodes     []interface{}
	predicate string
}

func (sn *sortedNodes) Len() int {
	return len(sn.nodes)
}

func (sn *sortedNodes) Swap(i, j int) {
	sn.nodes[i], sn.nodes[j] = sn.nodes[j], sn.nodes[i]
}

func (sn *sortedNodes) Less(i, j int) bool {
	switch a := sn.nodes[i].(map[string]interface{})[sn.predicate].(type) {
	case int:
		b := sn.nodes[j].(map[string]interface{})[sn.predicate].(int)
		return a < b
	case string:
		b := sn.nodes[j].(map[string]interface{})[sn.predicate].(string)
		return a < b
	}

	return false
}
