package dglite

func findSchema(schemas []Schema, predicate string) (Schema, bool) {
	for _, schema := range schemas {
		if schema.Predicate == predicate {
			return schema, true
		}
	}

	return Schema{}, false
}
