package ast

type InterfaceType interface {
	IsInterfaceType() bool
	// implements
	Type
}

// - InterfaceType
//   ```
//   INTERFACE
//   [InterfaceHeritage]
//   [ClassMethodList]
//   [ClassPropertyList]
//   ...
//   END
//   ```
type CustomInterfaceType struct {
	InterfaceHeritage *InterfaceHeritage
	MethodList        ClassMethodList
	PropertyList      ClassPropertyList
}

// - InterfaceHeritage
//   ```
//   '(' IdentList ')'
//   ```
type InterfaceHeritage struct {
	IdentList IdentList
}
