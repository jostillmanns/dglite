package dglite

import (
	"encoding/json"
	"fmt"

	"mooncamp.com/dgx/gql"
)

type RDF struct {
	Subject   uint64
	Object    interface{}
	Predicate string
}

func (rdf RDF) String() string {
	return fmt.Sprintf("<%d> <%s> <%v>", rdf.Subject, rdf.Predicate, rdf.Object)
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

type DGLite interface {
	Write(interface{}) ([]uint64, error)
	Read([]gql.GraphQuery, interface{}) error
	WriteRDF([]PlaceholderRDF)
}

type dglite struct {
	reader *reader
	writer *writer
	sorter *sorter
}

func New(schema []Schema) DGLite {
	db := newMapDB(schema)
	return &dglite{
		reader: &reader{schemas: schema, database: db, filter: filter{database: db}},
		writer: newWriter(schema),
		sorter: &sorter{},
	}
}

func (dgl *dglite) WriteRDF(rdfs []PlaceholderRDF) {
	actualRDFS := dgl.writer.resolvePlaceholder(rdfs)
	dgl.reader.database.Write(actualRDFS)
}

func (dgl *dglite) Write(in interface{}) ([]uint64, error) {
	js, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}

	var ns []map[string]interface{}
	if err := json.Unmarshal(js, &ns); err != nil {
		return nil, err
	}

	uids := []uint64{}
	for _, n := range ns {
		rdfs, uid := dgl.writer.rdfify(n)
		dgl.reader.database.Write(rdfs)
		uids = append(uids, uid)
	}

	return uids, nil
}

func (dgl *dglite) Read(qs []gql.GraphQuery, in interface{}) error {
	resolved := dgl.reader.resolveVariables(qs)
	var nodes []map[string]interface{}

	for _, e := range resolved {
		if e.Alias == "var" {
			continue
		}

		next := dgl.reader.read(e)
		next = dgl.sorter.sortNodes(e, next)
		nodes = append(next, nodes...)
	}

	js, err := json.MarshalIndent(nodes, "", "  ")
	if err != nil {
		return err
	}

	return json.Unmarshal(js, in)
}
