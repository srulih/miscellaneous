package visitorpattern

type Element interface {
	accept(Visitor)
}

type ConcreteElementA struct {
	A int
}

func (ce *ConcreteElementA) accept(v Visitor) {
	v.visitA(ce)
}

type ConcreteElementB struct {
	B int
}

func (ce *ConcreteElementB) accept(v Visitor) {
	v.visitB(ce)
}

type Visitor interface {
	visitA(*ConcreteElementA)
	visitB(*ConcreteElementB)
}

type ConcreteVisitor1 struct {
	Collect int
}

func (c *ConcreteVisitor1) visitA(ce *ConcreteElementA) {
	c.Collect += ce.A
}

func (c *ConcreteVisitor1) visitB(ce *ConcreteElementB) {
	c.Collect += ce.B
}
