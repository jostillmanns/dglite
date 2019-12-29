module mooncamp.com/dglite

go 1.13

require (
	github.com/dgraph-io/dgraph v1.1.1
	github.com/stretchr/testify v1.4.0
	mooncamp.com/dgx v0.0.0
)

replace mooncamp.com/dgx => ../awkward/dgx
