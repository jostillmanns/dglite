package dglite

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"mooncamp.com/dgx/gql"
)

type node struct {
	ID   uint64 `json:"uid"`
	Name string `json:"name"`

	Child    *node  `json:"child"`
	Children []node `json:"children"`
}

var testNodeSchema = []Schema{
	{Predicate: "children", Many: true, Type: "uid"},
	{Predicate: "child", Many: false, Type: "uid"},
	{Predicate: "name", Many: false, Type: "string"},
}

func Test_it_supports_children(t *testing.T) {
	in := node{
		ID:   1,
		Name: "A",
		Children: []node{
			{ID: 2, Name: "B"},
			{ID: 3, Name: "C"},
		},
	}

	js, err := json.Marshal(in)
	require.NoError(t, err)

	var n map[string]interface{}
	err = json.Unmarshal(js, &n)
	require.NoError(t, err)

	wr := newWriter(testNodeSchema)
	rdfs, uid := wr.rdfify(n)

	q := gql.GraphQuery{
		UID:  []uint64{uid},
		Func: &gql.Function{Name: "uid"},
		Children: []gql.GraphQuery{
			{Attr: "uid"}, {Attr: "name"},
			{
				Attr: "children",
				Children: []gql.GraphQuery{
					{Attr: "uid"}, {Attr: "name"},
				},
			},
		},
	}

	rdr := &reader{
		schemas:  testNodeSchema,
		database: newMapDB(testNodeSchema),
	}

	rdr.database.Write(rdfs)

	a := rdr.read(q)
	require.Len(t, a, 1)

	js, err = json.Marshal(a[0])
	require.NoError(t, err)

	var actual node
	err = json.Unmarshal(js, &actual)
	require.NoError(t, err)

	require.Equal(t, in, actual)
}

func Test_it_reads(t *testing.T) {
	in := node{
		ID:   1,
		Name: "A",
		Child: &node{
			ID:   2,
			Name: "B",
			Child: &node{
				ID:   3,
				Name: "B",
			},
		},
	}

	js, err := json.Marshal(in)
	require.NoError(t, err)

	var n map[string]interface{}
	err = json.Unmarshal(js, &n)
	require.NoError(t, err)

	wr := newWriter(testNodeSchema)
	rdfs, uid := wr.rdfify(n)

	q := gql.GraphQuery{
		UID:  []uint64{uid},
		Func: &gql.Function{Name: "uid"},
		Children: []gql.GraphQuery{
			{Attr: "name"}, {Attr: "uid"},
			{
				Attr: "child",
				Children: []gql.GraphQuery{
					{Attr: "name"}, {Attr: "uid"},
					{
						Attr: "child",
						Children: []gql.GraphQuery{
							{Attr: "name"}, {Attr: "uid"},
						},
					},
				},
			},
		},
	}

	rdr := &reader{
		schemas:  testNodeSchema,
		database: newMapDB(testNodeSchema),
	}
	rdr.database.Write(rdfs)
	a := rdr.read(q)
	require.Len(t, a, 1)

	js, err = json.Marshal(a[0])
	require.NoError(t, err)

	var actual node
	err = json.Unmarshal(js, &actual)
	require.NoError(t, err)

	require.Equal(t, in, actual)
}
