package typedcalculator

import (
	"fmt"
	"math"
	//"reflect"
)

func CreateEvaluator(p bool) *Eval {
	parser := BuildParser()

	return &Eval{
		parser:    parser,
		pr:        p,
		env:       make(map[string]Number),
		PrintVals: make([]Number, 0),
	}
}

type Eval struct {
	parser    *Parser
	current   int
	res       Number
	env       map[string]Number
	pr        bool     //whether to print the value
	PrintVals []Number //values that are printed used for testing
}

func (e *Eval) Run(program string) {
	node := e.parser.Parse(program)
	//printValue("", reflect.ValueOf(node), make(map[interface{}]bool))
	node.accept(e)
}

func (e *Eval) visitProgramStmt(f *Program) {
	for _, l := range f.Lines {
		l.accept(e)
	}
}

func (e *Eval) visitLineStmt(f *Line) {
	f.Stmt.accept(e)
}

func (e *Eval) visitAssignmentStmt(f *Assignment) {
	f.Expr.accept(e)
	e.res.Fixed = true
	e.env[f.Identifier] = e.res
}

func (e *Eval) visitPrintStmt(f *Print) {
	f.Expr.accept(e)
	if e.pr {
		fmt.Println(e.res)
	} else {
		e.PrintVals = append(e.PrintVals, e.res)
	}
}

func (e *Eval) visitResetStmt(f *Reset) {}
func (e *Eval) visitBinary2Stmt(f *Binary2) {
	f.Lhs.accept(e)
	lhs := e.res

	var rhs Number
	if f.Op != NOOP {
		f.Rhs.accept(e)
		rhs = e.res
	}

	rt := inferType(lhs, rhs)

	//change the types of the lhs and the rhs to the inferred type
	lhs = changeType(lhs, rt)
	rhs = changeType(rhs, rt)

	switch o := f.Op; o {
	case PLUS:
		{
			e.res = AddNums(lhs, rhs)
		}
	case MINUS:
		{
			e.res = SubNums(lhs, rhs)
		}
	case MULTIPLY:
		{
			e.res = MultiplyNums(lhs, rhs)
		}
	case DIVIDE:
		{
			e.res = DivideNums(lhs, rhs)
		}
	case POWER:
		{
			e.res = PowerNums(lhs, rhs)
		}
	case NOOP:
		{
			break
		}
	}

	e.res.Type = rt

}

func (e *Eval) visitIdentifierStmt(f *Identifier) {
	if v, ok := e.env[f.Val]; !ok {
		panic(fmt.Sprintf("unknown identifier %s", f.Val))
	} else {
		e.res = v
	}
}

func (e *Eval) visitNumberStmt(f *Number) {
	e.res = *f
}

//at this point both Numbers should have the same type
func AddNums(lhs Number, rhs Number) Number {
	if lhs.Type != rhs.Type {
		panic(fmt.Sprintf("lhs type %s is not equal to rhs type %s", TypeStringMap[lhs.Type], TypeStringMap[rhs.Type]))
	}
	if lhs.Type == INT {
		return Number{
			Type: lhs.Type,
			Num:  lhs.Num + rhs.Num,
		}
	}

	if lhs.Type == FLOAT {
		return Number{
			Type: lhs.Type,
			Flt:  lhs.Flt + rhs.Flt,
		}
	}
	panic(fmt.Sprintf("Unknown type %s", TypeStringMap[lhs.Type]))
	return Number{}
}

func SubNums(lhs Number, rhs Number) Number {
	if lhs.Type != rhs.Type {
		panic(fmt.Sprintf("lhs type %s is not equal to rhs type %s", TypeStringMap[lhs.Type], TypeStringMap[rhs.Type]))
	}
	if lhs.Type == INT {
		return Number{
			Type: lhs.Type,
			Num:  lhs.Num - rhs.Num,
		}
	}

	if lhs.Type == FLOAT {
		return Number{
			Type: lhs.Type,
			Flt:  lhs.Flt - rhs.Flt,
		}
	}
	panic(fmt.Sprintf("Unknown type %s", TypeStringMap[lhs.Type]))
	return Number{}
}

func MultiplyNums(lhs Number, rhs Number) Number {
	if lhs.Type != rhs.Type {
		panic(fmt.Sprintf("lhs type %s is not equal to rhs type %s", TypeStringMap[lhs.Type], TypeStringMap[rhs.Type]))
	}
	if lhs.Type == INT {
		return Number{
			Type: lhs.Type,
			Num:  lhs.Num * rhs.Num,
		}
	}

	if lhs.Type == FLOAT {
		return Number{
			Type: lhs.Type,
			Flt:  lhs.Flt * rhs.Flt,
		}
	}
	panic(fmt.Sprintf("Unknown type %d", lhs.Type))
	return Number{}
}

func DivideNums(lhs Number, rhs Number) Number {
	if lhs.Type != rhs.Type {
		panic(fmt.Sprintf("lhs type %d is not equal to rhs type %d", lhs.Type, rhs.Type))
	}
	if lhs.Type == INT {
		return Number{
			Type: lhs.Type,
			Num:  lhs.Num / rhs.Num,
		}
	}

	if lhs.Type == FLOAT {
		return Number{
			Type: lhs.Type,
			Flt:  lhs.Flt / rhs.Flt,
		}
	}
	panic(fmt.Sprintf("Unknown type %d", lhs.Type))
	return Number{}
}

func PowerNums(lhs Number, rhs Number) Number {
	if lhs.Type != rhs.Type {
		panic(fmt.Sprintf("lhs type %d is not equal to rhs type %d", lhs.Type, rhs.Type))
	}
	if lhs.Type == INT {
		return Number{
			Type: lhs.Type,
			Num:  intPow(lhs.Num, rhs.Num),
		}
	}

	if lhs.Type == FLOAT {
		return Number{
			Type: lhs.Type,
			Flt:  math.Pow(lhs.Flt, rhs.Flt),
		}
	}
	panic(fmt.Sprintf("Unknown type %d", lhs.Type))
	return Number{}
}

func intPow(x, y int) int {
	return int(math.Pow(float64(x), float64(y)))
}

// TYPE CHECKING

func inferType(n1 Number, n2 Number) Type {
	var err error
	var t Type
	//TODO should we allow a type to change if it fixed
	if (n1.Fixed && n2.Fixed) && (n1.Type != n2.Type) {
		panic(fmt.Sprintf("Unmatched types %s and %s", TypeStringMap[n1.Type], TypeStringMap[n2.Type]))
	} else if n1.Fixed && n2.Fixed {
		return n1.Type
	} else if n1.Fixed {
		t, err = convertible(n2, n1.Type)
		if err != nil {
			panic(err)
		}
		return t
	} else if n2.Fixed {
		t, err = convertible(n1, n2.Type)
		if err != nil {
			panic(err)
		}
		return t
	}
	return maxType(n1.Type, n2.Type)
}

func convertible(n Number, t Type) (Type, error) {
	if n.Type <= t {
		return t, nil
	}

	//TODO a float can be converted to a int if it has no decimal part eg 3.0 -> 3

	return NOTYPE, fmt.Errorf("cannot change %d to %d", t, n.Type)
}

func maxType(t1 Type, t2 Type) Type {
	if t1 >= t2 {
		return t1
	}
	return t2
}

func changeType(n1 Number, t Type) Number {
	if n1.Type == t {
		return n1
	}

	switch t {
	case INT:
		{
			return Number{
				Type: INT,
				Num:  int(n1.Flt),
			}
		}
	case FLOAT:
		{
			return Number{
				Type: FLOAT,
				Flt:  float64(n1.Num),
			}
		}

	}

	panic(fmt.Sprintf("Unknown type %s", TypeStringMap[t]))
	return Number{}
}
