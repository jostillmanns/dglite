package tour

import (
	"encoding/json"
	"fmt"
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

func Test_step_2(t *testing.T) {
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

func Test_step_3(t *testing.T) {
	q := gql.GraphQuery{
		Func: &gql.Function{Name: "allofterms", Attr: "name", Args: []gql.Arg{{Value: "Michael"}}},
		Children: []gql.GraphQuery{
			{Attr: "name"}, {Attr: "age"},
			{
				Attr:     "friend",
				Filter:   &gql.FilterTree{Func: &gql.Function{Name: "ge", Attr: "age", Args: []gql.Arg{{Value: "27"}}}},
				Children: []gql.GraphQuery{{Attr: "name"}, {Attr: "age"}},
			},
		},
	}

	expected := []Person{
		{
			Name: "Michael",
			Age:  39,
			Friend: []Person{
				{
					Name: "Artyom",
					Age:  35,
				},
				{
					Name: "Amit",
					Age:  35,
				},
				{
					Name: "Sarah",
					Age:  55,
				},
			},
		},
	}

	dgl := database(t)
	var actual []Person
	dgl.Read([]gql.GraphQuery{q}, &actual)

	js, _ := json.MarshalIndent(actual, "", "  ")
	fmt.Println(string(js))

	require.ElementsMatch(t, expected[0].Friend, actual[0].Friend)
	require.Equal(t, expected[0].Name, actual[0].Name)
	require.Equal(t, expected[0].Age, actual[0].Age)
}

func Test_step_4(t *testing.T) {
	q := gql.GraphQuery{
		Func: &gql.Function{Name: "allofterms", Attr: "name", Args: []gql.Arg{{Value: "Michael"}}},
		Children: []gql.GraphQuery{
			{Attr: "name"}, {Attr: "age"},
			{
				Attr: "friend",
				Filter: &gql.FilterTree{
					Op: "and",
					Child: []gql.FilterTree{
						{Func: &gql.Function{Name: "ge", Attr: "age", Args: []gql.Arg{{Value: "27"}}}},
						{Func: &gql.Function{Name: "le", Attr: "age", Args: []gql.Arg{{Value: "48"}}}},
					},
				},
				Children: []gql.GraphQuery{{Attr: "name"}, {Attr: "age"}},
			},
		},
	}

	expected := []Person{
		{
			Name: "Michael",
			Age:  39,
			Friend: []Person{
				{
					Name: "Artyom",
					Age:  35,
				},
				{
					Name: "Amit",
					Age:  35,
				},
			},
		},
	}

	dgl := database(t)
	var actual []Person
	dgl.Read([]gql.GraphQuery{q}, &actual)

	js, _ := json.MarshalIndent(actual, "", "  ")
	fmt.Println(string(js))

	require.ElementsMatch(t, expected[0].Friend, actual[0].Friend)
	require.Equal(t, expected[0].Name, actual[0].Name)
	require.Equal(t, expected[0].Age, actual[0].Age)
}

func Test_step_5(t *testing.T) {
	q := gql.GraphQuery{
		Func: &gql.Function{Name: "allofterms", Attr: "name", Args: []gql.Arg{{Value: "Michael"}}},
		Children: []gql.GraphQuery{
			{Attr: "name"}, {Attr: "age"},
			{
				Attr:     "friend",
				Order:    []gql.Order{{Attr: "age"}},
				Children: []gql.GraphQuery{{Attr: "name"}, {Attr: "age"}},
			},
		},
	}

	expected := []Person{
		{
			Name: "Michael",
			Age:  39,
			Friend: []Person{
				{
					Name: "Catalina",
					Age:  19,
				},
				{
					Name: "Sang Hyun",
					Age:  24,
				},
				{
					Name: "Amit",
					Age:  35,
				},
				{
					Name: "Artyom",
					Age:  35,
				},
				{
					Name: "Sarah",
					Age:  55,
				},
			},
		},
	}

	dgl := database(t)
	var actual []Person
	dgl.Read([]gql.GraphQuery{q}, &actual)

	js, _ := json.MarshalIndent(actual, "", "  ")
	fmt.Println(string(js))

	require.Equal(t, expected, actual)

	require.ElementsMatch(t, expected[0].Friend, actual[0].Friend)
	require.Equal(t, expected[0].Name, actual[0].Name)
	require.Equal(t, expected[0].Age, actual[0].Age)
}
