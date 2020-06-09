package visitorpattern

import "fmt"
import "testing"

func TestVisitor(t *testing.T) {

	ceA := ConcreteElementA{A: 100}
	ceB := ConcreteElementB{B: 400}

	v1 := &ConcreteVisitor1{}

	ceA.accept(v1)
	ceB.accept(v1)

	res := 500
	if v1.Collect != res {
		fmt.Errorf("Expected %d got %d instead", res, v1.Collect)
	}
}
