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

	fmt.Println("=== DETAILED AST COMPARISON ===")
	
	// Compare the modules using the Compare method
	compareResult := expectedModule.Compare(actualModule)
	fmt.Printf("Module Compare result: %d\n", compareResult)
	
	// Compare package
	fmt.Printf("Package Compare: %d\n", expectedModule.Package.Compare(actualModule.Package))
	
	// Compare rules
	fmt.Printf("Rules count - Expected: %d, Actual: %d\n", len(expectedModule.Rules), len(actualModule.Rules))
	
	// Detailed rule comparison
	for i, expectedRule := range expectedModule.Rules {
		if i < len(actualModule.Rules) {
			actualRule := actualModule.Rules[i]
			fmt.Printf("\nRule %d detailed comparison:\n", i)
			fmt.Printf("Rule Compare: %d\n", expectedRule.Compare(actualRule))
			fmt.Printf("Head Compare: %d\n", expectedRule.Head.Compare(actualRule.Head))
			fmt.Printf("Body Compare: %d\n", expectedRule.Body.Compare(actualRule.Body))
			
			// Compare body expressions one by one
			fmt.Printf("Body expressions comparison:\n")
			minLen := len(expectedRule.Body)
			if len(actualRule.Body) < minLen {
				minLen = len(actualRule.Body)
			}
			
			for j := 0; j < minLen; j++ {
				exprCompare := expectedRule.Body[j].Compare(actualRule.Body[j])
				fmt.Printf("  Expression %d Compare: %d\n", j, exprCompare)
				if exprCompare != 0 {
					fmt.Printf("    Expected: %v\n", expectedRule.Body[j])
					fmt.Printf("    Actual:   %v\n", actualRule.Body[j])
				}
			}
		}
	}
	
	// Let's also check if there are any differences in the string representation
	fmt.Println("\n=== STRING COMPARISON ===")
	expectedStr2 := expectedModule.String()
	actualStr2 := actualModule.String()
	
	if expectedStr2 == actualStr2 {
		fmt.Println("String representations are identical")
	} else {
		fmt.Println("String representations differ:")
		fmt.Printf("Expected:\n%s\n", expectedStr2)
		fmt.Printf("Actual:\n%s\n", actualStr2)
	}
}
