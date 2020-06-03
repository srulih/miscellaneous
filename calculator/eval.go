package calculator

import (
	"fmt"
	"math"
)

func intPow(x, y int) int {
	return int(math.Pow(float64(x), float64(y)))
}

func CreateEvaluator() *Evaluator {
	ev := Evaluator{
		env: Env{frames: make([]Frame, 0)},
	}

	ev.env.CreateFrame()
	return &ev
}

type Evaluator struct {
	env Env
}

func (e *Evaluator) Eval(node Node) int {

	switch n := node.(type) {
	case *programStmt:
		{
			var res int
			for _, dec := range n.declarations {
				res = e.Eval(dec)
			}
			return res
		}
	case *blockStmt:
		{
			e.env.CreateFrame()
			e.Eval(n.program)
			e.env.RemoveFrame()
		}
	case *assignStmt:
		{
			iden := n.identifier
			res := e.Eval(n.expr)
			e.env.AddVar(iden, res)
		}
	case *funcStmt:
		{
			iden := n.identifier
			e.env.AddVar(iden, n)
		}
	case *printExpr:
		{
			res := e.Eval(n.expr)
			fmt.Println(res)

		}
	case *ifExpr:
		{
			res := e.Eval(n.cmpExpr)
			//assume that 0 represents false
			if res != 0 {
				return e.Eval(n.thenStmt)
			} else {
				return e.Eval(n.elseStmt)
			}
		}
	case *binaryExpr:
		{
			var se *subExpr

			res := 0
			if len(n.subExprs) > 0 {
				res = e.Eval(n.subExprs[0])
			}

			for i := 1; i < len(n.subExprs); i++ {
				se = n.subExprs[i]

				switch op := se.Op; op {
				case PLUS:
					{
						res += e.Eval(se)
					}
				case MINUS:
					{
						res -= e.Eval(se)
					}
				case MULTIPLY:
					{
						res *= e.Eval(se)
					}
				case DIVIDE:
					{
						res /= e.Eval(se)
					}
				case ILLEGALOP:
					{
						panic("ILLEGALOP")
					}
				}
			}
			return res
		}
	case *unaryExpr:
		{
			if n.Op == MINUS {
				return -e.Eval(n.Right)
			}
			return e.Eval(n.Right)
		}
	case *callExpr:
		{
			fi, err := e.env.GetVar(n.funcName)
			if err != nil {
				panic(fmt.Sprintf("could not find function %s", n.funcName))
			}
			f := fi.(*funcStmt)

			if len(f.params) != len(n.args) {
				err := fmt.Sprintf("num params %d is not eql to num args %d", len(f.params), len(n.args))
				panic(err)
			}
			e.env.CreateFrame()

			var na int
			for i := 0; i < len(f.params); i++ {
				na = e.Eval(n.args[i])
				e.env.AddVar(f.params[i], na)
			}

			e.Eval(f.block)
		}
	case *subExpr:
		{
			return e.Eval(n.Expr)
		}
	case *number:
		{
			return n.num
		}
	case *identifier:
		{
			if val, err := e.env.GetVar(n.iden); err == nil {
				return val.(int)
			}
			panic(fmt.Sprintf("unbound indentifier %s", n.iden))
		}
	default:
		{
			panic(fmt.Sprintf("Unknown type %T", n))
		}
	}

	return 0
}
