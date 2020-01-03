package dglite

type database interface {
	Get(uid uint64) []RDF
	Write(rdfs []RDF)
	ReversePredicate(predicate string) []uint64
}

type mapdb struct {
	rdfs   map[uint64][]RDF
	schema []Schema
}

func newMapDB(schema []Schema) database {
	return &mapdb{rdfs: make(map[uint64][]RDF), schema: schema}
}

func (db *mapdb) ReversePredicate(predicate string) []uint64 {
	res := make([]uint64, 0)
	for k, v := range db.rdfs {
		if db.hasPredicate(v, predicate) {
			res = append(res, k)
		}
	}
	return res
}

func (db *mapdb) hasPredicate(rdfs []RDF, predicate string) bool {
	for _, e := range rdfs {
		if e.Predicate == predicate {
			return true
		}
	}

	return false
}

func (db *mapdb) Get(uid uint64) []RDF {
	return db.rdfs[uid]
}

func (db *mapdb) Write(rdfs []RDF) {
	for _, rdf := range rdfs {
		uid := rdf.Subject

		next := make([]RDF, 0, len(db.rdfs[uid]))

		schema, ok := findSchema(db.schema, rdf.Predicate)
		if !ok {
			continue
		}

		switch schema.Type {
		case "string":
			if rdf.Object.(string) == "" {
				continue
			}
		case "uid":
			if rdf.Object == nil {
				continue
			}
		}

		if !schema.Many {
			for _, e := range db.rdfs[uid] {
				if e.Predicate == rdf.Predicate {
					continue
				}

				next = append(next, e)
			}
		} else {
			next = db.rdfs[uid]
		}

		next = append(next, rdf)
		db.rdfs[uid] = next
	}
}
