package dglite

type database interface {
	Get(uid uint64) []RDF
	Write(rdfs []RDF)
}

type mapdb struct {
	rdfs   map[uint64][]RDF
	schema []Schema
}

func newMapDB(schema []Schema) database {
	return &mapdb{rdfs: make(map[uint64][]RDF), schema: schema}
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
