package dglite

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_it_creates_rdfs(t *testing.T) {
	in := map[string]interface{}{
		"name": "francesc",
		"uid":  float64(1),
		"teams": []interface{}{
			map[string]interface{}{"uid": float64(2), "name": "operation", "location": "cologne"},
			map[string]interface{}{"uid": float64(3), "name": "development", "location": "berlin"},
		},
		"office": map[string]interface{}{"uid": float64(4), "name": "cologne"},
	}

	expected := []RDF{
		{Subject: 1, Predicate: "office", Object: uint64(4), Type: "uid"},
		{Subject: 1, Predicate: "name", Object: "francesc", Type: "string"},
		{Subject: 1, Predicate: "teams", Object: uint64(2), Type: "uid"},
		{Subject: 1, Predicate: "teams", Object: uint64(3), Type: "uid"},
		{Subject: 2, Predicate: "name", Object: "operation", Type: "string"},
		{Subject: 2, Predicate: "location", Object: "cologne", Type: "string"},
		{Subject: 3, Predicate: "name", Object: "development", Type: "string"},
		{Subject: 3, Predicate: "location", Object: "berlin", Type: "string"},
		{Subject: 4, Predicate: "name", Object: "cologne", Type: "string"},
	}

	wr := newWriter()
	actual, _ := wr.rdfify(in)
	sort.Sort(RDFs(actual))
	sort.Sort(RDFs(expected))

	require.ElementsMatch(t, expected, actual)
}
