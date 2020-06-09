package typedcalculator

type Node interface {
	isNode()
	accept(Visitor)
}

type Op int

const (
	ILLEGALOP Op = iota
	NOOP
	PLUS
	MINUS
	MULTIPLY
	DIVIDE
	POWER
)

var OpStringMap = map[Op]string{
	NOOP:     "NOOP",
	PLUS:     "PLUS",
	MINUS:    "MINUS",
	MULTIPLY: "MULTIPLY",
	DIVIDE:   "DIVIDE",
	POWER:    "POWER",
}

var StringOpMap = map[string]Op{
	"NOOP":     NOOP,
	"PLUS":     PLUS,
	"MINUS":    MINUS,
	"MULTIPLY": MULTIPLY,
	"DIVIDE":   DIVIDE,
	"POWER":    POWER,
}

type Type int

const (
	NOTYPE Type = iota
	FLOAT
	DOUBLE
	INT
	LONG
)

var TypeStringMap = map[Type]string{
	FLOAT:  "FLOAT",
	DOUBLE: "DOUBLE",
	INT:    "INT",
	LONG:   "LONG",
}

var StringTypeMap = map[string]Type{
	"FLOAT":  FLOAT,
	"DOUBLE": DOUBLE,
	"INT":    INT,
	"LONG":   LONG,
}

type Program struct {
	Lines []*Line
}

func (f *Program) accept(v Visitor) {
	v.visitProgramStmt(f)
}

type Line struct {
	Stmt Node
}

func (f *Line) accept(v Visitor) {
	v.visitLineStmt(f)
}

type Assignment struct {
	Type       Type
	Identifier string
	Expr       Node
}

func (f *Assignment) accept(v Visitor) {
	v.visitAssignmentStmt(f)
}

type Print struct {
	Expr Node
}

func (f *Print) accept(v Visitor) {
	v.visitPrintStmt(f)
}

type Reset struct {
}

func (f *Reset) accept(v Visitor) {
	v.visitResetStmt(f)
}

type Binary2 struct {
	Type Type
	Op   Op
	Lhs  Node
	Rhs  Node
}

func (f *Binary2) accept(v Visitor) {
	v.visitBinary2Stmt(f)
}

type Identifier struct {
	Val string
}

func (f *Identifier) accept(v Visitor) {
	v.visitIdentifierStmt(f)
}

type Number struct {
	Type  Type
	Fixed bool
	Num   int
	Flt   float64
}

func (f *Number) accept(v Visitor) {
	v.visitNumberStmt(f)
}

func (f *Program) isNode()    {}
func (f *Line) isNode()       {}
func (f *Assignment) isNode() {}
func (f *Print) isNode()      {}
func (f *Reset) isNode()      {}
func (f *Binary2) isNode()    {}
func (f *Identifier) isNode() {}
func (f *Number) isNode()     {}

type Visitor interface {
	visitProgramStmt(f *Program)
	visitLineStmt(f *Line)
	visitAssignmentStmt(f *Assignment)
	visitPrintStmt(f *Print)
	visitResetStmt(f *Reset)
	visitBinary2Stmt(f *Binary2)
	visitIdentifierStmt(f *Identifier)
	visitNumberStmt(f *Number)
}
