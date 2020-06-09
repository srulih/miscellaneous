package typedcalculator

import (
	"math"
	"testing"
)

const float64EqualityThreshold = 1e-9

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= float64EqualityThreshold
}

func NumberSliceEqual(s1 []Number, s2 []Number) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i, _ := range s1 {
		if s1[i].Type != s2[i].Type {
			return false
		}
		if s1[i].Type == INT && s1[i].Num != s2[i].Num {
			return false
		}
		if s1[i].Type == FLOAT && !almostEqual(s1[i].Flt, s2[i].Flt) {
			return false
		}
	}
	return true
}

//TODO when reset is implemented the tests can be refactored so we dont
//create a new evaluator for every test.
func TestEval(t *testing.T) {

	intTable := map[string][]Number{
		"int val = 23; print val+2": []Number{Number{Type: INT, Num: 25}},
		"print 2+3+4":               []Number{Number{Type: INT, Num: 9}},
		"print 2*5+23":              []Number{Number{Type: INT, Num: 33}},
		"print 2**3+5":              []Number{Number{Type: INT, Num: 13}},
		"print 2**3*2+43-21":        []Number{Number{Type: INT, Num: 38}},
		"print 2**3**2":             []Number{Number{Type: INT, Num: 512}},
	}

	for pr, res := range intTable {
		iptr := CreateEvaluator(false)
		iptr.Run(pr)
		if !NumberSliceEqual(iptr.PrintVals, res) {
			t.Errorf("expected %+v got %+v", res, iptr.PrintVals)
		}
	}

	floatTable := map[string][]Number{
		"float v = 3.12; float p = 3.14; print v+p": []Number{Number{Type: FLOAT, Flt: 6.26}},
	}

	for pr, res := range floatTable {
		iptr := CreateEvaluator(false)
		iptr.Run(pr)
		if !NumberSliceEqual(iptr.PrintVals, res) {
			t.Errorf("expected %+v got %+v", res, iptr.PrintVals)
		}
	}

}
