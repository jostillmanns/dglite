package dglite

import (
	"strconv"

	"mooncamp.com/dgx/gql"
)

type filter struct {
	database database
}

func (fl *filter) value(uid uint64, predicate string) (interface{}, bool) {
	rdfs := fl.database.Get(uid)
	for _, e := range rdfs {
		if e.Predicate != predicate {
			continue
		}

		return e.Object, true
	}

	return nil, false
}

func (fl *filter) filter(uid uint64, tree gql.FilterTree) bool {
	if len(tree.Child) > 0 {
		all := []bool{}
		for _, e := range tree.Child {
			all = append(all, fl.filter(uid, e))
		}

		switch tree.Op {
		case "and":
			for _, e := range all {
				if !e {
					return false
				}
			}
			return true
		case "or":
			for _, e := range all {
				if e {
					return true
				}
			}
			return false
		case "not":
			return !all[0]
		default:
			return false
		}
	}

	switch tree.Func.Name {
	case "uid":
		for _, e := range tree.Func.UID {
			if e == uid {
				return true
			}
		}
		return false
	case "ge":
		value, ok := fl.value(uid, tree.Func.Attr)
		if !ok {
			return false
		}

		val, err := strconv.Atoi(tree.Func.Args[0].Value)
		if err != nil {
			return false
		}

		if value.(int) >= val {
			return true
		}

		return false

	case "le":
		value, ok := fl.value(uid, tree.Func.Attr)
		if !ok {
			return false
		}

		val, err := strconv.Atoi(tree.Func.Args[0].Value)
		if err != nil {
			return false
		}

		if value.(int) <= val {
			return true
		}

		return false
	default:
		return false
	}
}
