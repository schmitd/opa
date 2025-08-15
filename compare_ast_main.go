package main

import (
	"fmt"

	"github.com/open-policy-agent/opa/v1/ast"
)

func main() {
	// The test case from the failing test
	moduleStr := `package test
	f(_) = {"a": [1, 2, 3], "b": [4, 5, 6], "c": [7, 8, 9]}
	p := [ v | m := {l | l := f(true)[k]}[_]; v := m[_]; print(v)]
	`

	expectedStr := `package test
	f(__local0__) = {"a": [1, 2, 3], "b": [4, 5, 6], "c": [7, 8, 9]} if { true }
	p := __local7__ if { true; __local7__ = [__local3__ | __local6__ = {__local1__ | data.test.f(true, __local5__); __local1__ = __local5__[k]}; __local2__ = __local6__[_]; __local3__ = __local2__[_]; __local8__ = {__local4__ | __local4__ = __local3__}; internal.print([__local8__])] }
	`

	// Parse the original module
	originalModule, err := ast.ParseModule("test.rego", moduleStr)
	if err != nil {
		fmt.Printf("Error parsing original module: %v\n", err)
		return
	}

	// Parse the expected module
	expectedModule, err := ast.ParseModule("expected.rego", expectedStr)
	if err != nil {
		fmt.Printf("Error parsing expected module: %v\n", err)
		return
	}

	// Compile the original module to get the actual result
	c := ast.NewCompiler().WithEnablePrintStatements(true)
	c.Compile(map[string]*ast.Module{
		"test.rego": originalModule,
	})
	if c.Failed() {
		fmt.Printf("Compilation failed: %v\n", c.Errors)
		return
	}

	actualModule := c.Modules["test.rego"]

	fmt.Println("=== AST COMPARISON ===")

	// Compare the specific rule that's failing
	expectedRule := expectedModule.Rules[1]
	actualRule := actualModule.Rules[1]

	fmt.Printf("Expected Rule: %v\n", expectedRule)
	fmt.Printf("Actual Rule:   %v\n", actualRule)

	// Compare the body expressions
	fmt.Printf("\nExpected Body length: %d\n", len(expectedRule.Body))
	fmt.Printf("Actual Body length: %d\n", len(actualRule.Body))

	for i := 0; i < len(expectedRule.Body) && i < len(actualRule.Body); i++ {
		fmt.Printf("\nExpression %d:\n", i)
		fmt.Printf("Expected: %v\n", expectedRule.Body[i])
		fmt.Printf("Actual:   %v\n", actualRule.Body[i])

		// Compare using the Compare method
		compareResult := expectedRule.Body[i].Compare(actualRule.Body[i])
		fmt.Printf("Compare result: %d\n", compareResult)

		if compareResult != 0 {
			fmt.Printf("*** DIFFERENCE FOUND in expression %d ***\n", i)

			// Try to get more details about the difference
			expectedTerms := expectedRule.Body[i].Operands()
			actualTerms := actualRule.Body[i].Operands()

			fmt.Printf("Expected terms: %v\n", expectedTerms)
			fmt.Printf("Actual terms:   %v\n", actualTerms)

			if len(expectedTerms) == len(actualTerms) {
				for j, expectedTerm := range expectedTerms {
					actualTerm := actualTerms[j]
					termCompare := expectedTerm.Compare(actualTerm)
					fmt.Printf("Term %d compare: %d\n", j, termCompare)
					if termCompare != 0 {
						fmt.Printf("  Expected term %d: %v\n", j, expectedTerm)
						fmt.Printf("  Actual term %d:   %v\n", j, actualTerm)
					}
				}
			}
		}
	}
}
