package dglite

import (
	"mooncamp.com/dgx/gql"
)

type reader struct {
	schemas  []Schema
	database database
	filter   filter
}

func (rdr *reader) read(q gql.GraphQuery) []map[string]interface{} {
	var res []map[string]interface{}

	switch q.Func.Name {
	case "uid":
		for _, e := range q.UID {
			if q.Filter != nil {
				if !rdr.filter.filter(e, *q.Filter) {
					continue
				}
			}
			res = append(res, rdr.resolveChildren(e, q.Children))
		}
	case "has":
		uids := rdr.database.ReversePredicate(q.Func.Attr)
		cln := gql.CopyGraphQuery(q)
		cln.Func = &gql.Function{Name: "uid"}
		cln.UID = uids
		return rdr.read(cln)
	case "eq":
		args := []string{}
		for _, e := range q.Func.Args {
			args = append(args, e.Value)
		}

		uids := rdr.database.ReverseObject(q.Func.Attr, args)
		cln := gql.CopyGraphQuery(q)
		cln.Func = &gql.Function{Name: "uid"}
		cln.UID = uids
		return rdr.read(cln)
	case "allofterms":
		args := []string{}
		for _, e := range q.Func.Args {
			args = append(args, e.Value)
		}

		uids := rdr.database.ReverseObjectMatchAll(q.Func.Attr, args)
		cln := gql.CopyGraphQuery(q)
		cln.Func = &gql.Function{Name: "uid"}
		cln.UID = uids
		return rdr.read(cln)
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

			if schema.Type == "uid" {
				if q.Filter != nil {
					if !rdr.filter.filter(rdf.Object.(uint64), *q.Filter) {
						continue
					}
				}

				if schema.Many {
					if _, ok := res[q.Attr]; !ok {
						res[q.Attr] = []interface{}{}
					}

					res[q.Attr] = append(
						res[q.Attr].([]interface{}),
						rdr.resolveChildren(rdf.Object.(uint64), q.Children),
					)
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
