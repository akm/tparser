package ast

import "github.com/akm/tparser/ast/astcore"

// implemented by Program and Unit.
type Namespace interface {
	GetIdent() *Ident

	// Declarations in the namespace.
	// They are called by just only its ident or by Namespace.(ident)
	// InsideDecls doesn't include unit names in uses clause. because
	// they can't be called by namespace.(unit name)
	GetDeclMap() astcore.DeclMap
}

/*

## Reference in Project1 (USES Unit1)

- OK DeclInUnit1
- OK Unit1.DeclInUnit1
- OK DeclInProject1
- OK Project1.DeclInProject1
- NG Unit1
- NG Unit1.Unit1.DeclInUnit1
- NG Project1
- NG Project1.Project1.DeclInProject1

## Reference in Unit1 (USES Unit2)

- OK DeclInUnit1
- OK Unit1.DeclInUnit1
- OK DeclInUnit2
- OK Unit2.DeclInUnit2
- NG Unit1
- NG Unit2
- NG DeclInProject1
- NG Project1.DeclInProject1
- NG Unit1.Unit1.DeclInUnit1
- NG Unit2.Unit2.DeclInUnit2
- NG Unit1.Unit2.DeclInUnit2

## Steps

1. Get Ident
2. If Ident is a namespace, search with GetDeclMap.
3. If Ident is not a namespace, use it without namespace.

Context.DeclMap on parsing Proram includes
- Itself (Program)     ==> YES
- UsesClauseItem       ==> YES
- Decls in DeclSection ==> YES

Program.GetDeclMap (itself via Context.DeclMap) on parsing Program includes
- Itself (Program)     ==> NO
- UsesClauseItem       ==> NO
- Decls in DeclSection ==> YES

Context.DeclMap on parsing Unit1 includes
- Itself (Unit1)       ==> YES
- UsesClauseItem       ==> YES
- Decls in InterfaceSection ==> YES
- Decls in ImplementationSection ==> YES

Unit1.GetDeclMap (itself via Context.DeclMap) on parsing Unit1 includes
- Itself (Unit1)       ==> NO
- UsesClauseItem       ==> MO
- Decls in InterfaceSection ==> YES
- Decls in ImplementationSection ==> YES

Unit2.GetDeclMap on parsing Unit1 includes
- Itself (Unit2)       ==> NO
- UsesClauseItem       ==> NO
- Decls in InterfaceSection ==> YES
- Decls in ImplementationSection ==> NO

*/
