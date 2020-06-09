package typedcalculator

import (
	"fmt"
)

//match type of params for binary ops to type of result
var typeTable = map[[2]Type]Type{
	[2]Type{INT, INT}:     INT,
	[2]Type{FLOAT, INT}:   FLOAT,
	[2]Type{INT, FLOAT}:   FLOAT,
	[2]Type{FLOAT, FLOAT}: FLOAT,
}

type TypeChecker struct {
	Ty Type
	Nm Number
}

func (t *TypeChecker) Run(root Node) {
	root.accept(t)
}

func (t *TypeChecker) visitProgramStmt(f *Program) {
	for _, l := range f.Lines {
		l.accept(t)
	}
}

func (t *TypeChecker) visitLineStmt(f *Line) {
	f.Stmt.accept(t)
}

func (t *TypeChecker) visitAssignmentStmt(f *Assignment) {}
func (t *TypeChecker) visitPrintStmt(f *Print) {
	f.Expr.accept(t)
}

func (t *TypeChecker) visitResetStmt(f *Reset) {}
func (t *TypeChecker) visitBinary2Stmt(f *Binary2) {
	f.Lhs.accept(t)
	lhs := t.Nm

	var rhsType Number
	if f.Op != NOOP {
		f.Rhs.accept(t)
		rhsType = t.Nm
	}

	t.Nm
	f.Type = t.Ty
}

func (t *TypeChecker) visitIdentifierStmt(f *Identifier) {
}

func (t *TypeChecker) visitNumberStmt(f *Number) {
	t.Nm = f
}

func inferType(n1 Number, n2 Number) Type {
	var err error
	var t Type
	//TODO should we allow a type to change if it fixed
	if (n1.Fixed && n2.Fixed) && (n1.Type != n2.Type) {
		panic(fmt.Sprintf("Unmatched types %d and %d", n1.Type, n2.Type))
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
		return nil
	}

	//TODO a float can be converted to a int if it has no decimal part eg 3.0 -> 3

	return ILLEGALOP, fmt.Errorf("cannot change %d to %d", t, n.Type)
}

func maxType(t1 Type, t2 Type) Type {
	if t1 >= t2 {
		return t1
	}
	return t2
}
