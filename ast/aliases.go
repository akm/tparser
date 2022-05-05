package ast

import "github.com/akm/tparser/ast/astcore"

type (
	Ident     = astcore.Ident
	IdentList = astcore.IdentList
	Node      = astcore.Node
	Nodes     = astcore.Nodes

	Position = astcore.Position
	Location = astcore.Location
)

var (
	NewIdent     = astcore.NewIdent
	NewIdentFrom = astcore.NewIdentFrom
	NewIdentList = astcore.NewIdentList
	NewLocation  = astcore.NewLocation
)
