package token

import (
	"github.com/akm/opparser/ext"
)

var Directives = ext.Strings{
	// "absolute", // Used in VarDecl
	"abstract",
	// "assembler", // Not found in Grammer
	// "automated", // Not found in Grammer
	"cdecl",
	// "contains", // Used in ContainsClause
	// "default",  // Used in PropertySpecifiers
	// "deprecated", // ==> PortabilityDirective
	// "dispid", // Not found in Grammer
	"dynamic",
	"export",
	"external",
	"far",
	"forward",
	// "implements", // Used in PropertySpecifiers
	// "index", // Used in ExportsItem or PropertySpecifiers
	// "library", // Used in library / ==> PortabilityDirective
	"local",
	"message",
	// "name", // Used in ExportsItem
	"near",
	// "nodefault", // Used in PropertySpecifiers
	"overload",
	"override",
	// "package", // Used in Package
	"pascal",
	// "platform", // ==> PortabilityDirective
	// "private", // Used in ClassVisibility
	// "protected", // Used in ClassVisibility
	// "public", // Used in ClassVisibility
	// "published", // Used in ClassVisibility
	// "read", // Used in PropertySpecifiers
	// "readonly", // Not found in Grammer
	"register",
	"reintroduce",
	// "requires", // Used in RequiresClause
	// "resident", // Not found in Grammer
	"safecall",
	"stdcall",
	// "stored", // Used in PropertySpecifiers
	"varargs",
	// "virtual", // Used in MethodList
	// "write", // Used in PropertySpecifiers
	// "writeonly", // Not found in Grammer
}.Set()

var PortabilityDirectives = ext.Strings{
	"platform",
	"deprecated",
	"library",
}.Set()
