package dglite

type writer struct {
	ids    chan uint64
	schema []Schema
}

func newWriter(schema []Schema) *writer {
	wr := &writer{
		ids:    make(chan uint64),
		schema: schema,
	}

	go func() {
		c := uint64(0)
		for {
			c = c + 1
			wr.ids <- c
		}
	}()

	return wr
}

func (wr *writer) rdfify(node map[string]interface{}) ([]RDF, uint64) {
	res := []RDF{}
	var uid uint64

	if actual, ok := node["uid"]; ok && actual != float64(0) {
		uid = uint64(node["uid"].(float64))
	} else {
		uid = <-wr.ids
	}

	for k, v := range node {
		if k == "uid" {
			continue
		}

		if v == nil {
			continue
		}

		schema, ok := findSchema(wr.schema, k)
		if !ok {
			continue
		}

		switch schema.Type {
		case "uid":
			if schema.Many {
				actual := v.([]interface{})

				for _, e := range actual {
					next, id := wr.rdfify(e.(map[string]interface{}))
					if len(next) == 0 {
						continue
					}

					res = append(res, RDF{Subject: uid, Predicate: k, Object: id})
					res = append(res, next...)
				}
				continue
			}

			actual := v.(map[string]interface{})
			next, id := wr.rdfify(actual)
			if len(next) == 0 {
				continue
			}

			res = append(res, RDF{Subject: uid, Predicate: k, Object: id})
			res = append(res, next...)
		default:
			res = append(res, RDF{Subject: uid, Predicate: k, Object: v})
		}
	}

	return res, uid
}
