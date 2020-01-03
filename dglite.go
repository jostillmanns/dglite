package dglite

import (
	"encoding/json"
	"fmt"

	"mooncamp.com/dgx/gql"
)

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

type DGLite interface {
	Write(interface{}) ([]uint64, error)
	Read([]gql.GraphQuery, interface{}) error
}

type dglite struct {
	reader *reader
	writer *writer
}

func New(schema []Schema) DGLite {
	return &dglite{
		reader: &reader{schemas: schema, database: newMapDB(schema)},
		writer: newWriter(),
	}
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
		nodes = append(next, nodes...)
	}

	js, err := json.MarshalIndent(nodes, "", "  ")
	if err != nil {
		return err
	}

	return json.Unmarshal(js, in)
}
