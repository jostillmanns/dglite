package dglite

import (
	"mooncamp.com/dgx/gql"
)

type reader struct {
	schemas  []Schema
	database database
}

func (rdr *reader) read(q gql.GraphQuery) []map[string]interface{} {
	var res []map[string]interface{}

	switch q.Func.Name {
	case "uid":
		for _, e := range q.UID {
			res = append(res, rdr.resolveChildren(e, q.Children))
		}
	default:
	}
	return res
}

func (rdr *reader) resolveChildren(uid uint64, qs []gql.GraphQuery) map[string]interface{} {
	res := make(map[string]interface{})
	for _, q := range qs {
		if q.Attr == "uid" {
			res["uid"] = uid
			continue
		}

		for _, rdf := range rdr.database.Get(uid) {
			if rdf.Predicate != q.Attr {
				continue
			}

			schema, ok := rdr.findSchema(q.Attr)
			if !ok {
				continue
			}

			if rdf.Type == "uid" {
				if schema.Many {
					if _, ok := res[q.Attr]; !ok {
						res[q.Attr] = []interface{}{}
					}

					res[q.Attr] = append(res[q.Attr].([]interface{}), rdr.resolveChildren(rdf.Object.(uint64), q.Children))
					continue
				}

				res[q.Attr] = rdr.resolveChildren(rdf.Object.(uint64), q.Children)
				continue
			}

			res[q.Attr] = rdf.Object
		}
	}
	return res
}

func (rdr *reader) findSchema(predicate string) (Schema, bool) {
	return findSchema(rdr.schemas, predicate)
}
