package tour

import (
	"testing"

	"github.com/stretchr/testify/require"
	"mooncamp.com/dglite"
	"mooncamp.com/dgx/gql"
)

var schema = []dglite.Schema{
	{Predicate: "name", Type: "string"},
	{Predicate: "age", Type: "int"},
	{Predicate: "friend", Type: "uid", Many: true},
	{Predicate: "owns_pet", Type: "uid"},
}

type Person struct {
	ID      uint64   `json:"uid"`
	Name    string   `json:"name"`
	Age     int      `json:"age"`
	Friend  []Person `json:"friend"`
	OwnsPet *Pet     `json:"owns_pet"`
}

type Pet struct {
	ID   uint64 `json:"uid"`
	Name string `json:"name"`
}

func database(t *testing.T) dglite.DGLite {
	dgl := dglite.New(schema)
	dgl.WriteRDF(friendRDFS)

	return dgl
}

func Test_step_1(t *testing.T) {
	q := gql.GraphQuery{
		Func:     &gql.Function{Name: "eq", Attr: "name", Args: []gql.Arg{{Value: "Michael"}}},
		Children: []gql.GraphQuery{{Attr: "name"}, {Attr: "age"}},
	}

	expected := []Person{
		{
			Age:  39,
			Name: "Michael",
		},
	}

	dgl := database(t)

	var actual []Person
	dgl.Read([]gql.GraphQuery{q}, &actual)

	require.Equal(t, expected, actual)
}

func Test_stop_2(t *testing.T) {
	q := gql.GraphQuery{
		Func: &gql.Function{Name: "eq", Attr: "name", Args: []gql.Arg{{Value: "Michael"}}},
		Children: []gql.GraphQuery{
			{Attr: "name"}, {Attr: "age"},
			{Attr: "friend", Children: []gql.GraphQuery{{Attr: "name"}}},
		},
	}

	expected := []Person{
		{
			Age:  39,
			Name: "Michael",
			Friend: []Person{
				{
					Name: "Sarah",
				},
				{
					Name: "Sang Hyun",
				},
				{
					Name: "Artyom",
				},
				{
					Name: "Amit",
				},
				{
					Name: "Catalina",
				},
			},
		},
	}

	dgl := database(t)

	var actual []Person
	dgl.Read([]gql.GraphQuery{q}, &actual)

	require.ElementsMatch(t, expected[0].Friend, actual[0].Friend)
}
