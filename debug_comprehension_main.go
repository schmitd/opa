package main

import (
	"fmt"

	"github.com/open-policy-agent/opa/v1/ast"
)

func main() {
	expectedStr := `package test
	f(__local0__) = {"a": [1, 2, 3], "b": [4, 5, 6], "c": [7, 8, 9]} if { true }
	p := __local7__ if { true; __local7__ = [__local3__ | __local6__ = {__local1__ | data.test.f(true, __local5__); __local1__ = __local5__[k]}; __local2__ = __local6__[_]; __local3__ = __local2__[_]; __local8__ = {__local4__ | __local4__ = __local3__}; internal.print([__local8__])] }
	`

	// Parse the expected module
	expectedModule, err := ast.ParseModule("expected.rego", expectedStr)
	if err != nil {
		fmt.Printf("Error parsing expected module: %v\n", err)
		return
	}

	// Get the second rule and examine the array comprehension
	rule := expectedModule.Rules[1]
	fmt.Printf("Rule: %v\n", rule)
	fmt.Printf("Rule Head: %v\n", rule.Head)
	fmt.Printf("Rule Body length: %d\n", len(rule.Body))

	// Examine the second expression (index 1) which contains the array comprehension
	expr := rule.Body[1]
	fmt.Printf("Expression 1: %v\n", expr)
	fmt.Printf("Expression 1 type: %T\n", expr)

	// Check if it's an equality expression
	if len(expr.Terms) >= 3 {
		fmt.Printf("Expression has %d terms\n", len(expr.Terms))
		for i, term := range expr.Terms {
			fmt.Printf("Term %d: %v (type: %T)\n", i, term, term)
		}

		// The right side should be an array comprehension
		if len(expr.Terms) == 3 {
			rightSide := expr.Terms[2]
			fmt.Printf("Right side: %v (type: %T)\n", rightSide, rightSide)

			// Check if it's an array comprehension
			if arrayComp, ok := rightSide.Value.(*ast.ArrayComprehension); ok {
				fmt.Printf("Array comprehension found!\n")
				fmt.Printf("Term: %v\n", arrayComp.Term)
				fmt.Printf("Body: %v\n", arrayComp.Body)
				fmt.Printf("Body length: %d\n", len(arrayComp.Body))

				// Examine each expression in the comprehension body
				for i, compExpr := range arrayComp.Body {
					fmt.Printf("Comprehension expression %d: %v\n", i, compExpr)
					fmt.Printf("Comprehension expression %d type: %T\n", i, compExpr)
				}
			} else {
				fmt.Printf("Right side is not an array comprehension, it's: %T\n", rightSide.Value)
			}
		}
	}
}
