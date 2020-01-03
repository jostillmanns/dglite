package dglite

import "mooncamp.com/dgx/gql"

type filter struct{}

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
	default:
		return false
	}
}
