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

// type_test_runner.go
type BaseTestRunner = runners.BaseTestRunner

// type_section_test_runner.go
type TypeSectionTestRunner = runners.TypeSectionTestRunner

var NewTypeSectionTestRunner = runners.NewTypeSectionTestRunner

// var_section_test_runner.go
var RunVarSectionTest = runners.RunVarSectionTest

// type_test_runner.go
type TypeTestRunner = runners.TypeTestRunner

var (
	NewTypeTestRunner = runners.NewTypeTestRunner
	RunTypeTest       = runners.RunTypeTest
)

// unit_test_runner.go
type UnitTestRunner = runners.UnitTestRunner

var (
	NewUnitTestRunner = runners.NewUnitTestRunner
	RunUnitTest       = runners.RunUnitTest
)

// program_test_runner.go
type ProgramTestRunner = runners.ProgramTestRunner

var (
	NewProgramTestRunner = runners.NewProgramTestRunner
	RunProgramTest       = runners.RunProgramTest
)
