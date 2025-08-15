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

	fmt.Println("=== ORIGINAL MODULE ===")
	printModuleAST(originalModule)
	
	fmt.Println("\n=== EXPECTED MODULE ===")
	printModuleAST(expectedModule)
	
	fmt.Println("\n=== ACTUAL MODULE ===")
	printModuleAST(actualModule)

	// Compare expected vs actual
	fmt.Println("\n=== COMPARISON ===")
	if expectedModule.Equal(actualModule) {
		fmt.Println("Modules are equal!")
	} else {
		fmt.Println("Modules are NOT equal!")
		
		// Detailed comparison
		fmt.Println("\nExpected module rules:", len(expectedModule.Rules))
		fmt.Println("Actual module rules:", len(actualModule.Rules))
		
		for i, rule := range expectedModule.Rules {
			if i < len(actualModule.Rules) {
				fmt.Printf("\nRule %d comparison:\n", i)
				fmt.Printf("Expected: %v\n", rule)
				fmt.Printf("Actual:   %v\n", actualModule.Rules[i])
				if rule.Equal(actualModule.Rules[i]) {
					fmt.Printf("Rule %d: EQUAL\n", i)
				} else {
					fmt.Printf("Rule %d: NOT EQUAL\n", i)
					// Compare specific parts
					if rule.Head.Equal(actualModule.Rules[i].Head) {
						fmt.Printf("  Head: EQUAL\n")
					} else {
						fmt.Printf("  Head: NOT EQUAL\n")
						fmt.Printf("  Expected Head: %v\n", rule.Head)
						fmt.Printf("  Actual Head:   %v\n", actualModule.Rules[i].Head)
					}
					if rule.Body.Equal(actualModule.Rules[i].Body) {
						fmt.Printf("  Body: EQUAL\n")
					} else {
						fmt.Printf("  Body: NOT EQUAL\n")
						fmt.Printf("  Expected Body: %v\n", rule.Body)
						fmt.Printf("  Actual Body:   %v\n", actualModule.Rules[i].Body)
					}
				}
			}
		}
	}
}

func printModuleAST(module *ast.Module) {
	fmt.Printf("Package: %v\n", module.Package)
	fmt.Printf("Rules (%d):\n", len(module.Rules))
	for i, rule := range module.Rules {
		fmt.Printf("  Rule %d:\n", i)
		fmt.Printf("    Head: %v\n", rule.Head)
		fmt.Printf("    Body: %v\n", rule.Body)
	}
}
