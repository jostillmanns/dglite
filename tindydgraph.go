package tinydgraph

import "fmt"

type RDF struct {
	Subject         uint64
	Object          interface{}
	Predicate, Type string
}

func (rdf RDF) String() string {
	return fmt.Sprintf("<%d> <%s> <%v>:%s", rdf.Subject, rdf.Predicate, rdf.Object, rdf.Type)
}

type RDFs []RDF

func (s RDFs) Len() int           { return len(s) }
func (s RDFs) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s RDFs) Less(i, j int) bool { return s[i].Subject < s[j].Subject }

type Schema struct {
	Predicate string
	Many      bool
	Type      string
}
