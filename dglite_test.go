package dglite

import (
	"testing"

	"github.com/stretchr/testify/require"
	"mooncamp.com/dgx/gql"
)

type User struct {
	ID     uint64  `json:"uid"`
	Name   string  `json:"user.name"`
	Office *Office `json:"user.office"`
	Teams  []Team  `json:"user.teams"`
}

type Office struct {
	ID       uint64    `json:"uid"`
	Name     string    `json:"office.name"`
	Location *Location `json:"office.location"`
}

type Team struct {
	ID   uint64 `json:"uid"`
	Name string `json:"team.name"`
}

type Location struct {
	City    string `json:"location.city"`
	Country string `json:"location.country"`
}

var testSchema = []Schema{
	{Predicate: "user.name", Type: "string", Many: false},
	{Predicate: "user.office", Type: "uid", Many: false},
	{Predicate: "user.teams", Type: "uid", Many: true},
	{Predicate: "user.name", Type: "string", Many: false},
	{Predicate: "office.name", Type: "string", Many: false},
	{Predicate: "office.location", Type: "uid", Many: false},
	{Predicate: "location.city", Type: "string", Many: false},
	{Predicate: "location.country", Type: "string", Many: false},
	{Predicate: "team.name", Type: "string", Many: false},
}

func Test_it_finds_all_locations(t *testing.T) {
	in := []User{
		{
			Name: "francesc",
			Office: &Office{
				Name: "berlin office",
				Location: &Location{
					City:    "Berlin",
					Country: "Germany",
				},
			},
		},
		{
			Name: "manish",
			Office: &Office{
				Name: "amsterdam office",
				Location: &Location{
					City:    "Amsterdam",
					Country: "Netherlands",
				},
			},
		},
	}

	dgl := New(testSchema)
	uids, err := dgl.Write(in)
	require.NoError(t, err)

	require.Len(t, uids, 2)

	q := []gql.GraphQuery{
		{
			Alias: "var",
			UID:   uids,
			Func:  &gql.Function{Name: "uid"},
			Children: []gql.GraphQuery{
				{
					Attr: "user.office",
					Children: []gql.GraphQuery{
						{
							Attr: "office.location",
							Children: []gql.GraphQuery{
								{Attr: "uid", Var: "office_locations"},
							},
						},
					},
				},
			},
		},
		{
			NeedsVar: []gql.VarContext{{Name: "office_locations"}},
			Func:     &gql.Function{Name: "uid"},
			Children: []gql.GraphQuery{
				{Attr: "location.city"}, {Attr: "location.country"},
			},
		},
	}

	var actual []Location
	err = dgl.Read(q, &actual)
	require.NoError(t, err)

	expected := []Location{
		{
			City:    "Berlin",
			Country: "Germany",
		},
		{
			City:    "Amsterdam",
			Country: "Netherlands",
		},
	}

	require.ElementsMatch(t, expected, actual)
}
