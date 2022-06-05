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
type TypeSectionTestRunner = runners.TypeSectionTestRunner

var NewTypeSectionTestRunner = runners.NewTypeSectionTestRunner

// type_test_runner.go
type TypeTestRunner = runners.TypeTestRunner

var NewTypeTestRunner = runners.NewTypeTestRunner
