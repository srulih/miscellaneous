package calculator

import "testing"

func TestEvaluator(t *testing.T) {
	p := BuildParser()
	eval := CreateEvaluator()

	mathExprs := map[string]int{
		//"3*4+2-3+5":                  16,
		"2+2":   4,
		"4-6+3": 1,
		"-2":    -2,
		//"2*7*11":                     154,
		/*`set abc = 23
		set cde = 34
		set hello = 42
		abc + cde + hello`: 99,*/
		`func hello(a, b) { print(a+b) }
		hello(2,3)`: 0,
	}

	for expr, res := range mathExprs {
		parsed := p.Parse(expr)
		calcRes := eval.Eval(parsed)
		if calcRes != res {
			t.Errorf("expected %s to evaluate to %d got %d", expr, res, calcRes)
		}
	}
}
