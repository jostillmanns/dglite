package dglite

import (
	"mooncamp.com/dgx/gql"
)

func (rdr *reader) resolveVariables(qs []gql.GraphQuery) []gql.GraphQuery {
	res := []gql.GraphQuery{}
	for _, q := range qs {
		if q.NeedsVar == nil {
			res = append(res, q)
			continue
		}

		res = append(res, rdr.resolveVariable(q, qs))
	}

	return res
}

func (rdr *reader) resolveVariable(on gql.GraphQuery, qs []gql.GraphQuery) gql.GraphQuery {
	cln := gql.CopyGraphQuery(on)
	for _, q := range qs {
		if !rdr.hasVariable(on.NeedsVar[0].Name, q) {
			continue
		}

		if q.NeedsVar != nil {
			q = rdr.resolveVariable(q, qs)
		}

		nodes := rdr.read(q)

		uints := []uint64{}
		for _, n := range nodes {
			uints = append(uints, rdr.grabVariableValue(n)...)
		}

		cln.UID = uints
		cln.NeedsVar = nil
	}
	return cln
}

func (rdr *reader) grabVariableValue(node map[string]interface{}) []uint64 {
	for _, v := range node {
		switch actual := v.(type) {
		case map[string]interface{}:
			return rdr.grabVariableValue(actual)

		case []interface{}:
			if len(actual) == 0 {
				return nil
			}

			if _, ok := actual[0].(map[string]interface{}); !ok {
				res := []uint64{}
				for _, e := range actual {
					res = append(res, e.(uint64))
				}
				return res
			}

			res := []uint64{}
			for _, e := range actual {
				res = append(res, rdr.grabVariableValue(e.(map[string]interface{}))...)
			}
			return res
		default:
			return []uint64{actual.(uint64)}
		}
	}

	return nil
}

func (rdr *reader) hasVariable(v string, q gql.GraphQuery) bool {
	for _, c := range q.Children {
		if rdr.hasVariable(v, c) {
			return true
		}
	}

	return q.Var == v
}
