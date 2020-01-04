package tour

import "mooncamp.com/dglite"

var friendRDFS = []dglite.PlaceholderRDF{
	{PlaceHolder: "_:michael", RDF: dglite.RDF{Predicate: "name", Object: "Michael"}},
	{PlaceHolder: "_:michael", RDF: dglite.RDF{Predicate: "dgraph.type", Object: "Person"}},
	{PlaceHolder: "_:michael", RDF: dglite.RDF{Predicate: "age", Object: 39}},
	{PlaceHolder: "_:michael", RDF: dglite.RDF{Predicate: "friend", Object: "_:amit"}},
	{PlaceHolder: "_:michael", RDF: dglite.RDF{Predicate: "friend", Object: "_:sarah"}},
	{PlaceHolder: "_:michael", RDF: dglite.RDF{Predicate: "friend", Object: "_:sang"}},
	{PlaceHolder: "_:michael", RDF: dglite.RDF{Predicate: "friend", Object: "_:catalina"}},
	{PlaceHolder: "_:michael", RDF: dglite.RDF{Predicate: "friend", Object: "_:artyom"}},
	{PlaceHolder: "_:michael", RDF: dglite.RDF{Predicate: "owns_pet", Object: "_:rammy"}},

	{PlaceHolder: "_:amit", RDF: dglite.RDF{Predicate: "name", Object: "अमित"}},
	{PlaceHolder: "_:amit", RDF: dglite.RDF{Predicate: "name", Object: "অমিত"}},
	{PlaceHolder: "_:amit", RDF: dglite.RDF{Predicate: "name", Object: "Amit"}},
	{PlaceHolder: "_:amit", RDF: dglite.RDF{Predicate: "dgraph.type", Object: "Person"}},
	{PlaceHolder: "_:amit", RDF: dglite.RDF{Predicate: "age", Object: 35}},
	{PlaceHolder: "_:amit", RDF: dglite.RDF{Predicate: "friend", Object: "_:michael"}},
	{PlaceHolder: "_:amit", RDF: dglite.RDF{Predicate: "friend", Object: "_:sang"}},
	{PlaceHolder: "_:amit", RDF: dglite.RDF{Predicate: "friend", Object: "_:artyom"}},

	{PlaceHolder: "_:luke", RDF: dglite.RDF{Predicate: "name", Object: "Luke"}},
	{PlaceHolder: "_:luke", RDF: dglite.RDF{Predicate: "dgraph.type", Object: "Person"}},
	{PlaceHolder: "_:luke", RDF: dglite.RDF{Predicate: "name", Object: "Łukasz"}},
	{PlaceHolder: "_:luke", RDF: dglite.RDF{Predicate: "age", Object: 77}},

	{PlaceHolder: "_:artyom", RDF: dglite.RDF{Predicate: "name", Object: "Артём"}},
	{PlaceHolder: "_:artyom", RDF: dglite.RDF{Predicate: "name", Object: "Artyom"}},
	{PlaceHolder: "_:artyom", RDF: dglite.RDF{Predicate: "dgraph.type", Object: "Person"}},
	{PlaceHolder: "_:artyom", RDF: dglite.RDF{Predicate: "age", Object: 35}},

	{PlaceHolder: "_:sarah", RDF: dglite.RDF{Predicate: "name", Object: "Sarah"}},
	{PlaceHolder: "_:sarah", RDF: dglite.RDF{Predicate: "dgraph.type", Object: "Person"}},
	{PlaceHolder: "_:sarah", RDF: dglite.RDF{Predicate: "age", Object: 55}},

	{PlaceHolder: "_:sang", RDF: dglite.RDF{Predicate: "name", Object: "상현"}},
	{PlaceHolder: "_:sang", RDF: dglite.RDF{Predicate: "name", Object: "Sang Hyun"}},
	{PlaceHolder: "_:sang", RDF: dglite.RDF{Predicate: "dgraph.type", Object: "Person"}},
	{PlaceHolder: "_:sang", RDF: dglite.RDF{Predicate: "age", Object: "24"}},
	{PlaceHolder: "_:sang", RDF: dglite.RDF{Predicate: "friend", Object: "_:amit"}},
	{PlaceHolder: "_:sang", RDF: dglite.RDF{Predicate: "friend", Object: "_:catalina"}},
	{PlaceHolder: "_:sang", RDF: dglite.RDF{Predicate: "friend", Object: "_:hyung"}},
	{PlaceHolder: "_:sang", RDF: dglite.RDF{Predicate: "owns_pet", Object: "_:goldie"}},

	{PlaceHolder: "_:hyung", RDF: dglite.RDF{Predicate: "name", Object: "형신"}},
	{PlaceHolder: "_:hyung", RDF: dglite.RDF{Predicate: "name", Object: "Hyung Sin"}},
	{PlaceHolder: "_:hyung", RDF: dglite.RDF{Predicate: "dgraph.type", Object: "Person"}},
	{PlaceHolder: "_:hyung", RDF: dglite.RDF{Predicate: "friend", Object: "_:sang"}},

	{PlaceHolder: "_:catalina", RDF: dglite.RDF{Predicate: "name", Object: "Catalina"}},
	{PlaceHolder: "_:catalina", RDF: dglite.RDF{Predicate: "dgraph.type", Object: "Person"}},
	{PlaceHolder: "_:catalina", RDF: dglite.RDF{Predicate: "age", Object: 19}},
	{PlaceHolder: "_:catalina", RDF: dglite.RDF{Predicate: "friend", Object: "_:sang"}},
	{PlaceHolder: "_:catalina", RDF: dglite.RDF{Predicate: "owns_pet", Object: "_:perro"}},

	{PlaceHolder: "_:rammy", RDF: dglite.RDF{Predicate: "name", Object: "Rammy the sheep"}},
	{PlaceHolder: "_:rammy", RDF: dglite.RDF{Predicate: "dgraph.type", Object: "Animal"}},

	{PlaceHolder: "_:goldie", RDF: dglite.RDF{Predicate: "name", Object: "Goldie"}},
	{PlaceHolder: "_:goldie", RDF: dglite.RDF{Predicate: "dgraph.type", Object: "Animal"}},

	{PlaceHolder: "_:perro", RDF: dglite.RDF{Predicate: "name", Object: "Perro"}},
	{PlaceHolder: "_:perro", RDF: dglite.RDF{Predicate: "dgraph.type", Object: "Animal"}},
}
