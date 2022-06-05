package parsertest

import (
	"github.com/akm/tparser/parser/parsertest/runners"
)

// base.go
var (
	NewTestParser         = runners.NewTestParser
	NewTestProgramParser  = runners.NewTestProgramParser
	NewTestUnitParser     = runners.NewTestUnitParser
	NewTestProgramContext = runners.NewTestProgramContext
	NewTestUnitContext    = runners.NewTestUnitContext
)

// type_section_test_runner.go
type TypeSectionTest = runners.TypeSectionTest

var NewTypeSectionTest = runners.NewTypeSectionTest

// type_test_runner.go
type TypeTest = runners.TypeTest

var NewTypeTest = runners.NewTypeTest
