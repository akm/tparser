package ast

import "github.com/akm/tparser/ast/astcore"

type (
	Ident     = astcore.Ident
	IdentRef  = astcore.IdentRef
	IdentList = astcore.IdentList
	Node      = astcore.Node
	Nodes     = astcore.Nodes

	Position = astcore.Position
	Location = astcore.Location
)

var (
	NewIdent     = astcore.NewIdent
	NewIdentRef  = astcore.NewIdentRef
	NewIdentFrom = astcore.NewIdentFrom
	NewIdentList = astcore.NewIdentList
	NewLocation  = astcore.NewLocation
)
