package token

import (
	"github.com/akm/opparser/ext"
)

var Directives = ext.Strings{
	"absolute",
	"abstract",
	"assembler",
	"automated",
	"cdecl",
	"contains",
	"default",
	"deprecated",
	"dispid",
	"dynamic",
	"export",
	"external",
	"far",
	"forward",
	"implements",
	"index",
	"library",
	"local",
	"message",
	"name",
	"near",
	"nodefault",
	"overload",
	"override",
	"package",
	"pascal",
	"platform",
	"private",
	"protected",
	"public",
	"published",
	"read",
	"readonly",
	"register",
	"reintroduce",
	"requires",
	"resident",
	"safecall",
	"stdcall",
	"stored",
	"varargs",
	"virtual",
	"write",
	"writeonly",
}.Set()