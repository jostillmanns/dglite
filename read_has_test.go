package dglite

import (
	"testing"

	"github.com/stretchr/testify/require"
	"mooncamp.com/dgx/gql"
)

func Test_it_supports_has_function(t *testing.T) {
	user := []User{{
		Name: "francest",
	}}

	location := []Location{{
		City: "cologne",
	}}

	dgl := New(testSchema)
	_, err := dgl.Write(user)
	require.NoError(t, err)
	_, err = dgl.Write(location)
	require.NoError(t, err)

	q := []gql.GraphQuery{{
		Func:     &gql.Function{Name: "has", Attr: "user.name"},
		Children: []gql.GraphQuery{{Attr: "user.name"}},
	}}

	var actual []User
	err = dgl.Read(q, &actual)
	require.NoError(t, err)

	require.Equal(t, user, actual)
}

func Test_it_supports_uid_filter(t *testing.T) {
	user := []User{
		{
			Name: "francesc",
		},
		{
			Name: "manish",
		},
	}

	dgl := New(testSchema)
	uids, err := dgl.Write(user)
	require.NoError(t, err)

	q := []gql.GraphQuery{{
		Func:     &gql.Function{Name: "has", Attr: "user.name"},
		Filter:   &gql.FilterTree{Func: &gql.Function{Name: "uid", UID: []uint64{uids[0]}}},
		Children: []gql.GraphQuery{{Attr: "uid"}},
	}}

	var actual []User
	err = dgl.Read(q, &actual)
	require.NoError(t, err)

	require.Equal(t, []User{{ID: uids[0]}}, actual)
}

func Test_it_supports_filter_trees(t *testing.T) {
	dgl := New(testSchema)

	users := []User{{Name: "user a"}, {Name: "user b"}, {Name: "user c"}}
	userIDs, err := dgl.Write(users)
	require.NoError(t, err)

	q := []gql.GraphQuery{{
		Func: &gql.Function{Name: "has", Attr: "user.name"},
		Filter: &gql.FilterTree{
			Op: "or",
			Child: []gql.FilterTree{
				{Func: &gql.Function{UID: []uint64{userIDs[0]}, Name: "uid"}},
				{Func: &gql.Function{UID: []uint64{userIDs[1]}, Name: "uid"}},
			},
		},
		Children: []gql.GraphQuery{{Attr: "uid"}},
	}}

	var actual []User
	err = dgl.Read(q, &actual)
	require.NoError(t, err)

	require.Equal(t, []User{{ID: userIDs[0]}, {ID: userIDs[1]}}, actual)
}

func Test_it_supports_uid_filter_on_nested(t *testing.T) {
	dgl := New(testSchema)

	teams := []Team{{Name: "team a"}, {Name: "team b"}}
	teamIDs, err := dgl.Write(teams)
	require.NoError(t, err)

	user := []User{{Name: "user", Teams: []Team{{ID: teamIDs[0]}, {ID: teamIDs[1]}}}}
	userIDs, err := dgl.Write(user)
	require.NoError(t, err)

	q := []gql.GraphQuery{{
		UID:  userIDs,
		Func: &gql.Function{Name: "uid"},
		Children: []gql.GraphQuery{
			{Attr: "uid"},
			{
				Attr:     "user.teams",
				Children: []gql.GraphQuery{{Attr: "team.name"}},
			},
		},
	}}

	var notFiltered []User
	err = dgl.Read(q, &notFiltered)
	require.NoError(t, err)
	require.Len(t, notFiltered, 1)
	require.Len(t, notFiltered[0].Teams, 2)

	q = []gql.GraphQuery{{
		UID:  userIDs,
		Func: &gql.Function{Name: "uid"},
		Children: []gql.GraphQuery{
			{
				Attr:     "user.teams",
				Filter:   &gql.FilterTree{Func: &gql.Function{UID: []uint64{teamIDs[0]}, Name: "uid"}},
				Children: []gql.GraphQuery{{Attr: "team.name"}},
			},
		},
	}}

	var filtered []User
	err = dgl.Read(q, &filtered)
	require.NoError(t, err)
	require.Len(t, filtered[0].Teams, 1)
}
