package tinydgraph

type writer struct {
	ids chan uint64
}

func newWriter() *writer {
	wr := &writer{
		ids: make(chan uint64),
	}

	go func() {
		c := uint64(0)
		for {
			c := c + 1
			wr.ids <- c
		}
	}()

	return wr
}

func (wr *writer) rdfify(node map[string]interface{}) ([]RDF, uint64) {
	res := []RDF{}
	var uid uint64

	if _, ok := node["uid"]; ok {
		uid = uint64(node["uid"].(float64))
	} else {
		uid = <-wr.ids
	}

	for k, v := range node {
		if k == "uid" {
			continue
		}

		switch actual := v.(type) {
		case map[string]interface{}:
			next, id := wr.rdfify(actual)
			if len(next) == 0 {
				continue
			}

			res = append(res, RDF{Subject: uid, Predicate: k, Object: id, Type: "uid"})
			res = append(res, next...)
		case []interface{}:
			if len(actual) == 0 {
				continue
			}

			if _, ok := actual[0].(map[string]interface{}); !ok {
				res = append(res, RDF{Subject: uid, Predicate: k, Object: actual})
				continue
			}

			for _, e := range actual {
				next, id := wr.rdfify(e.(map[string]interface{}))
				if len(next) == 0 {
					continue
				}

				res = append(res, RDF{Subject: uid, Predicate: k, Object: id, Type: "uid"})
				res = append(res, next...)
			}
		default:
			res = append(res, RDF{Subject: uid, Predicate: k, Object: actual, Type: "string"})
		}
	}

	return res, uid
}
