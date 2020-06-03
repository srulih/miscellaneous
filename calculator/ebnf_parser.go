package calculator

/*

EBNF:

<program>     : <stmt> {<seperator> <stmt>}

<stmt>        : <assign_stmt>
			  | <func_stmt>
              | <if_stmt>
              | <cmp_expr>
			  | <block>
			  | <print>

<assign_stmt> : set <id> = <cmp_expr>

<func_stmt>   : func <function>

<function>    : <id> ( [<parameters>] ) block

<parameters>  : <id> { , <id> }


Note 'else' binds to the innermost 'if', like in C

<if_stmt>     : if <cmp_expr> then <stmt> [else <stmt>]

<block>       : "{" <program> "}"

<print>       : print(<cmp_expr>)

//TODO add in all expressions until arithExpr

<cmp_expr>    : <bitor_expr> [== <bitor_expr>]
              | <bitor_expr> [!= <bitor_expr>]
              | <bitor_expr> [> <bitor_expr>]
              | <bitor_expr> [< <bitor_expr>]
              | <bitor_expr> [>= <bitor_expr>]
              | <bitor_expr> [<= <bitor_expr>]

<bitor_expr>  | <bitxor_expr> {| <bitxor_expr>}

<bitxor_expr> | <bitand_expr> {^ <bitand_expr>}

<bitand_expr> | <shift_expr> {& <shift_expr>}

<shift_expr>  | <arith_expr> {<< <arith_expr>}
              : <arith_expr> {>> <arith_expr>}

<arith_expr>  : <term> {+ <term>}
              | <term> {- <term>}

<term>        : <power> {* <power>}
              | <power> {/ <power>}

<power>       : <factor> ** <power>
              | <factor>

<factor>      : <id>
              | <number>
              | - <factor>
              | ( <cmp_expr> )

<unary>       : - <call>

<call>        : primary { ( {<arguments>} )}

<arguments>   : <cmp_expr> { , <cmp_expr> }

<primary>     : <id> | <number> |

<id>          : [a-zA-Z_]\w+
<number>      : \d+
<seperator>   : \n
			  | ;
*/

import "fmt"
import "strconv"

type Op int

const (
	ILLEGALOP Op = iota
	PLUS
	MINUS
	MULTIPLY
	DIVIDE
	POWER
)

type Node interface {
	isNode()
}

type programStmt struct {
	declarations []Node
}

type assignStmt struct {
	identifier string
	expr       Node
}

type funcStmt struct {
	identifier string
	params     []string
	block      Node
}

type binaryExpr struct {
	Op       Op
	subExprs []*subExpr
}

func (a *binaryExpr) accept(v Visitor) {
	v.binaryExpr(a)
}

//TODO should if be a stmt or expr?
type ifExpr struct {
	cmpExpr  Node
	thenStmt Node
	elseStmt Node
}

type blockStmt struct {
	program Node
}

type printExpr struct {
	expr Node
}

type identifier struct {
	iden string
}

func (a *identifier) accept(v Visitor) {
	v.identifier(a)
}

type number struct {
	num int
}

func (a *number) accept(v Visitor) {
	v.number(a)
}

type unaryExpr struct {
	Op    Op
	Right Node
}

func (a *unaryExpr) accept(v Visitor) {
	v.unaryExpr(a)
}

type callExpr struct {
	funcName string
	args     []Node
}

type subExpr struct {
	Op   Op
	Expr Node
}

func (a *subExpr) accept(v Visitor) {
	v.subExpr(a)
}

func (e *programStmt) isNode() {}
func (e *assignStmt) isNode()  {}
func (e *binaryExpr) isNode()  {}
func (e *funcStmt) isNode()    {}
func (e *identifier) isNode()  {}
func (e *number) isNode()      {}
func (e *unaryExpr) isNode()   {}
func (e *subExpr) isNode()     {}
func (e *ifExpr) isNode()      {}
func (e *blockStmt) isNode()   {}
func (e *printExpr) isNode()   {}
func (e *callExpr) isNode()    {}

type Parser struct {
	Lexer        *Lexer
	CurrentToken *Token
}

func BuildParser() *Parser {
	regexMap := [][2]string{
		{`set`, "SET"},
		{`if`, "IF"},
		{`then`, "THEN"},
		{`else`, "ELSE"},
		{`func`, "FUNC"},
		{`print`, "PRINT"},
		{`\n`, "NEWLINE"},
		{`\d+`, "NUMBER"},
		{`[a-zA-Z_]\w*`, "IDENTIFIER"},
		{`\*\*`, "**"},
		{`!=`, "!="},
		{`==`, "=="},
		{`>=`, ">="},
		{`<=`, "<="},
		{`>>`, ">>"},
		{`<<`, "<<"},
		{`&`, "&"},
		{`\^`, "^"},
		{`\|`, "|"},
		{`<`, "<"},
		{`>`, ">"},
		{`\+`, "+"},
		{`\-`, "-"},
		{`\*`, "*"},
		{`\/`, "/"},
		{`\(`, "("},
		{`\)`, ")"},
		{`\{`, "{"},
		{`\}`, "}"},
		{`=`, "="},
		{`;`, ";"},
		{`,`, ","},
	}

	p := &Parser{}
	p.Lexer = CreateLexer(regexMap, true)
	return p
}

func (p *Parser) getNextToken() {
	var err error
	p.CurrentToken, err = p.Lexer.Token()
	if err != nil {
		panic(err)
	}
}

func (p *Parser) matchMultipleTokens(ts ...string) string {
	for _, t := range ts {
		if t == p.CurrentToken.Type {
			v := p.CurrentToken.Value
			p.getNextToken()
			return v
		}
	}
	err := fmt.Errorf("expcted one of the types %+v got type %s", ts, p.CurrentToken.Type)
	panic(err)
	return ""
}

func (p *Parser) matchToken(t string) string {
	if t == p.CurrentToken.Type {
		v := p.CurrentToken.Value
		p.getNextToken()
		return v
	}

	err := fmt.Errorf("expcted type %s got type %s", t, p.CurrentToken.Type)
	panic(err)
	return ""
}

func (p *Parser) Parse(program string) Node {

	p.Lexer.Input(program)
	p.Lexer.Reset()
	p.CurrentToken, _ = p.Lexer.Token()

	return p.parseProgram()
}

func (p *Parser) parseProgram() Node {
	stmts := make([]Node, 0)
	var n Node
	for p.CurrentToken.Type != "EOF" {
		n = p.parseDeclaration()
		stmts = append(stmts, n)
		if p.CurrentToken.Type != "EOF" {
			p.matchMultipleTokens(";", "NEWLINE")
		}
	}

	if p.CurrentToken.Type != "EOF" {
		panic(fmt.Sprintf("Expectd type EOF got %s instead", p.CurrentToken.Type))
	}
	return &programStmt{
		declarations: stmts,
	}
}

func (p *Parser) parseDeclaration() Node {
	c := p.CurrentToken.Type
	var node Node
	if c == "SET" {
		node = p.parseAssignStmt()
	} else if c == "IF" {
		node = p.parseIfExpr()
	} else if c == "{" {
		node = p.parseBlockExpr()
	} else if c == "PRINT" {
		node = p.parsePrintStmt()
	} else if c == "FUNC" {
		node = p.parseFuncStmt()
	} else {
		node = p.parseCmpExpr()
	}
	return node
}

func (p *Parser) parseAssignStmt() Node {
	p.matchToken("SET")
	iden := p.matchToken("IDENTIFIER")
	p.matchToken("=")
	cmpExpr := p.parseCmpExpr()
	return &assignStmt{
		identifier: iden,
		expr:       cmpExpr,
	}
}

func (p *Parser) parseIfExpr() Node {
	p.matchToken("IF")
	ce := p.parseCmpExpr()
	p.matchToken("THEN")
	te := p.parseDeclaration()
	var ee Node
	if p.CurrentToken.Type == "ELSE" {
		p.matchToken("ELSE")
		ee = p.parseDeclaration()
	}
	return &ifExpr{
		cmpExpr:  ce,
		thenStmt: te,
		elseStmt: ee,
	}
}

//TODO code in parseProgram is quite similar
func (p *Parser) parseBlockExpr() Node {
	p.matchToken("{")
	stmts := make([]Node, 0)
	var n Node
	for p.CurrentToken.Type != "}" {
		n = p.parseDeclaration()
		stmts = append(stmts, n)
		if p.CurrentToken.Type != "}" {
			p.matchMultipleTokens(";", "NEWLINE")
		}
	}

	if p.CurrentToken.Type != "}" {
		panic(fmt.Sprintf("Expectd type } got %s instead", p.CurrentToken.Type))
	}
	p.matchToken("}")

	return &blockStmt{
		program: &programStmt{
			declarations: stmts,
		},
	}
}

func (p *Parser) parsePrintStmt() Node {
	p.matchToken("PRINT")
	p.matchToken("(")
	expr := p.parseCmpExpr()
	p.matchToken(")")

	return &printExpr{
		expr: expr,
	}
}

func (p *Parser) parseFuncStmt() Node {
	p.matchToken("FUNC")
	funcIden := p.matchToken("IDENTIFIER")
	p.matchToken("(")

	params := make([]string, 0)
	var par string
	for p.CurrentToken.Type != ")" {
		par = p.matchToken("IDENTIFIER")
		params = append(params, par)
		if p.CurrentToken.Type != ")" {
			p.matchToken(",")
		}
	}

	p.matchToken(")")

	block := p.parseBlockExpr()

	return &funcStmt{
		identifier: funcIden,
		params:     params,
		block:      block,
	}
}

func (p *Parser) parseCmpExpr() Node {
	return p.parseArithExpr()
}

func (p *Parser) parseArithExpr() Node {
	subExprs := make([]*subExpr, 0)
	var op Op

	lhs := p.parseTerm()
	subExprs = append(subExprs, &subExpr{Op: ILLEGALOP, Expr: lhs})
	for {
		if p.CurrentToken.Type == "+" {
			p.matchToken("+")
			op = PLUS
		} else if p.CurrentToken.Type == "-" {
			p.matchToken("-")
			op = MINUS
		} else {
			break
		}
		rhs := p.parseTerm()
		se := &subExpr{
			Op:   op,
			Expr: rhs,
		}
		subExprs = append(subExprs, se)
	}

	return &binaryExpr{subExprs: subExprs}
}

func (p *Parser) parseTerm() Node {
	subExprs := make([]*subExpr, 0)

	lhs := p.parsePower()
	subExprs = append(subExprs, &subExpr{Op: ILLEGALOP, Expr: lhs})

	var op Op
	for {
		if p.CurrentToken.Type == "*" {
			p.matchToken("*")
			op = MULTIPLY
		} else if p.CurrentToken.Type == "/" {
			p.matchToken("/")
			op = DIVIDE
		} else {
			break
		}

		rhs := p.parsePower()
		se := &subExpr{
			Op:   op,
			Expr: rhs,
		}
		subExprs = append(subExprs, se)
	}

	return &binaryExpr{subExprs: subExprs}
}

// power is right associative so the factor element has ILLEGALOP
func (p *Parser) parsePower() Node {
	subExprs := make([]*subExpr, 0)

	factor := p.parseUnary()

	for {
		if p.CurrentToken.Type == "**" {
			p.matchToken("**")
		} else {
			break
		}

		subExprs = append(subExprs, &subExpr{Op: POWER, Expr: factor})
		factor = p.parseUnary()
	}

	subExprs = append(subExprs, &subExpr{Op: ILLEGALOP, Expr: factor})

	return &binaryExpr{
		subExprs: subExprs,
	}
}

func (p *Parser) parseUnary() Node {
	if p.CurrentToken.Type == "-" {
		p.matchToken("-")
		right := p.parseUnary()
		return &unaryExpr{Op: MINUS, Right: right}
	}

	return p.parseCall()
}

func (p *Parser) parseCall() Node {
	primary := p.parsePrimary()

	if p.CurrentToken.Type == "(" {
		return p.finishCall(primary)
	}
	return primary
}

func (p *Parser) finishCall(expr Node) Node {
	p.matchToken("(")

	tExpr := expr.(*identifier)

	if p.CurrentToken.Type == ")" {
		p.matchToken(")")
		return &callExpr{funcName: tExpr.iden}
	}
	args := make([]Node, 0)

	a1 := p.parseCmpExpr()
	args = append(args, a1)
	for p.CurrentToken.Type == "," {
		p.matchToken(",")
		a1 = p.parseCmpExpr()
		args = append(args, a1)
	}
	p.matchToken(")")

	return &callExpr{funcName: tExpr.iden, args: args}
}

func (p *Parser) parsePrimary() Node {
	c := p.CurrentToken.Type
	if c == "NUMBER" {
		num, err := strconv.Atoi(p.matchToken("NUMBER"))
		if err != nil {
			panic(err)
		}
		return &number{num: num}
	} else if c == "IDENTIFIER" {
		iden := p.matchToken("IDENTIFIER")
		return &identifier{iden: iden}
	} else {
		p := fmt.Sprintf("Unknown factor type %s", c)
		panic(p)
	}
	return nil
}
