package dglite

import (
	"mooncamp.com/dgx/gql"
)

type reader struct {
	schemas []Schema
}

func (rdr *reader) read(q gql.GraphQuery, rdfs []RDF) []map[string]interface{} {
	var res []map[string]interface{}

	switch q.Func.Name {
	case "uid":
		for _, e := range q.UID {
			res = append(res, rdr.resolveChildren(e, q.Children, rdfs))
		}
	default:
	}
	return res
}

func (rdr *reader) resolveChildren(uid uint64, qs []gql.GraphQuery, rdfs []RDF) map[string]interface{} {
	res := make(map[string]interface{})
	for _, q := range qs {
		if q.Attr == "uid" {
			res["uid"] = uid
			continue
		}

		for _, rdf := range rdfs {
			if rdf.Subject != uid {
				continue
			}

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

					res[q.Attr] = append(res[q.Attr].([]interface{}), rdr.resolveChildren(rdf.Object.(uint64), q.Children, rdfs))
					continue
				}

				res[q.Attr] = rdr.resolveChildren(rdf.Object.(uint64), q.Children, rdfs)
				continue
			}

			res[q.Attr] = rdf.Object
		}
	}
	return res
}

func (rdr *reader) findSchema(predicate string) (Schema, bool) {
	for _, schema := range rdr.schemas {
		if schema.Predicate == predicate {
			return schema, true
		}
	}

	return Schema{}, false
}
