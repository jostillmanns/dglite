package dglite

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"mooncamp.com/dgx/gql"
)

func Test_it_supports_time_objects(t *testing.T) {
	tm := time.Now()
	user := User{
		Name:     "user a",
		Birthday: tm,
	}

	dgl := New(testSchema)
	uids, err := dgl.Write([]User{user})
	require.NoError(t, err)

	q := gql.GraphQuery{
		Func:     &gql.Function{Name: "uid"},
		UID:      uids,
		Children: []gql.GraphQuery{{Attr: "uid"}, {Attr: "user.birthday"}},
	}

	var actual []User
	err = dgl.Read([]gql.GraphQuery{q}, &actual)
	require.NoError(t, err)

	require.Equal(t, tm.Round(time.Second).Local(), actual[0].Birthday.Round(time.Second).Local())
}
