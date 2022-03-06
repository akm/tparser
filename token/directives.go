package token

import (
	"strings"

	"github.com/akm/opparser/ext"
)

func isDirective(w string) bool {
	return directives.Include(strings.ToUpper(w))
}

var directives = ext.Strings{
	// "ABSOLUTE", // Used in VarDecl
	"ABSTRACT",
	// "ASSEMBLER", // Not found in Grammer
	// "AUTOMATED", // Not found in Grammer
	"CDECL",
	// "CONTAINS", // Used in ContainsClause
	// "DEFAULT",  // Used in PropertySpecifiers
	// "DEPRECATED", // ==> PortabilityDirective
	// "DISPID", // Not found in Grammer
	"DYNAMIC",
	"EXPORT",
	"EXTERNAL",
	"FAR",
	"FORWARD",
	// "IMPLEMENTS", // Used in PropertySpecifiers
	// "INDEX", // Used in ExportsItem or PropertySpecifiers
	// "LIBRARY", // Used in library / ==> PortabilityDirective
	"LOCAL",
	"MESSAGE",
	// "NAME", // Used in ExportsItem
	"NEAR",
	// "NODEFAULT", // Used in PropertySpecifiers
	"OVERLOAD",
	"OVERRIDE",
	// "PACKAGE", // Used in Package
	"PASCAL",
	// "PLATFORM", // ==> PortabilityDirective
	// "PRIVATE", // Used in ClassVisibility
	// "PROTECTED", // Used in ClassVisibility
	// "PUBLIC", // Used in ClassVisibility
	// "PUBLISHED", // Used in ClassVisibility
	// "READ", // Used in PropertySpecifiers
	// "READONLY", // Not found in Grammer
	"REGISTER",
	"REINTRODUCE",
	// "REQUIRES", // Used in RequiresClause
	// "RESIDENT", // Not found in Grammer
	"SAFECALL",
	"STDCALL",
	// "STORED", // Used in PropertySpecifiers
	"VARARGS",
	// "VIRTUAL", // Used in MethodList
	// "WRITE", // Used in PropertySpecifiers
	// "WRITEONLY", // Not found in Grammer
}.Set()

func isPortabilityDirective(w string) bool {
	return portabilityDirectives.Include(strings.ToUpper(w))
}

var portabilityDirectives = ext.Strings{
	"PLATFORM",
	"DEPRECATED",
	"LIBRARY",
}.Set()
