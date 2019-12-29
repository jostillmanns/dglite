package dglite

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"mooncamp.com/dgx/gql"
)

func Test_it_resolves_variables(t *testing.T) {
	in := node{
		Name: "root",
		ID:   1,
		Children: []node{
			{Name: "a", ID: 2},
			{Name: "b", ID: 3},
		},
	}

	js, err := json.Marshal(in)
	require.NoError(t, err)

	var n map[string]interface{}
	err = json.Unmarshal(js, &n)
	require.NoError(t, err)

	wr := newWriter()
	rdfs, uid := wr.rdfify(n)

	qs := []gql.GraphQuery{
		{
			Alias: "var",
			UID:   []uint64{uid},
			Func:  &gql.Function{Name: "uid"},
			Children: []gql.GraphQuery{
				{
					Attr:     "children",
					Children: []gql.GraphQuery{{Attr: "uid", Var: "root_children"}},
				},
			},
		},
		{
			Alias:    "n",
			Func:     &gql.Function{Name: "uid"},
			NeedsVar: []gql.VarContext{{Name: "root_children"}},
			Children: []gql.GraphQuery{{Attr: "uid"}},
		},
	}

	rdr := &reader{
		schemas: []Schema{
			{Predicate: "children", Many: true},
			{Predicate: "child", Many: false},
			{Predicate: "name", Many: false},
		},
	}
	actual := rdr.resolveVariables(qs, rdfs)

	expected := []gql.GraphQuery{
		{
			Alias: "var",
			UID:   []uint64{uid},
			Func:  &gql.Function{Name: "uid"},
			Children: []gql.GraphQuery{
				{
					Attr:     "children",
					Children: []gql.GraphQuery{{Attr: "uid", Var: "root_children"}},
				},
			},
		},
		{
			Alias:    "n",
			Func:     &gql.Function{Name: "uid"},
			UID:      []uint64{2, 3},
			Children: []gql.GraphQuery{{Attr: "uid"}},
		},
	}

	require.Equal(t, expected, actual)
}
