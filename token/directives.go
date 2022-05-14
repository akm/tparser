package token

import (
	"fmt"
	"strings"

	"github.com/akm/tparser/ext"
)

type DirectivePredicator struct {
	nameSet ext.StringSet
}

func NewDirectivePredicator(nameSet ext.StringSet) *DirectivePredicator {
	return &DirectivePredicator{nameSet: nameSet}
}

func (m *DirectivePredicator) Name() string {
	return fmt.Sprintf("Directive(%s)", strings.Join(m.nameSet.Slice(), ","))
}

func (m *DirectivePredicator) Predicate(t *Token) bool {
	return m.nameSet.Include(strings.ToUpper(t.RawString()))
}

func Directives(names ...string) Predicator {
	for _, name := range names {
		if !directives.Include(name) {
			panic(fmt.Sprintf("%s is not a directive", name))
		}
	}
	return NewDirectivePredicator(ext.NewStringSet(names...))
}

var Directive = NewDirectivePredicator(directives)

func PortabilityDirectives(names ...string) Predicator {
	for _, name := range names {
		if !portabilityDirectives.Include(name) {
			panic(fmt.Sprintf("%s is not a portable directive", name))
		}
	}
	return NewDirectivePredicator(ext.NewStringSet(names...))
}

var PortabilityDirective = NewDirectivePredicator(portabilityDirectives)

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
	"INDEX", // Used in ExportsItem or PropertySpecifiers. Or in EXTERNAL Options
	// "LIBRARY", // Used in library / ==> PortabilityDirective
	"LOCAL",
	"MESSAGE",
	"NAME", // Used in ExportsItem. . Or in EXTERNAL Options
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

var portabilityDirectives = ext.Strings{
	"PLATFORM",
	"DEPRECATED",
	"LIBRARY",
}.Set()
