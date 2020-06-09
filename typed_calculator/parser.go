package typedcalculator

import "fmt"
import "strconv"

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
		{`reset`, "RESET"},
		{`int`, "TYPE"},
		{`float`, "TYPE"},
		{`\n`, "NEWLINE"},
		{`\d*\.\d+`, "DECIMAL"},
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
	lines := make([]*Line, 0)
	var n Node
	for p.CurrentToken.Type != "EOF" {
		n = p.parseLine()
		if p.CurrentToken.Type != "EOF" {
			p.matchToken(";")
		}
		lines = append(lines, &Line{Stmt: n})
	}
	return &Program{
		Lines: lines,
	}
}

func (p *Parser) parseLine() Node {
	switch t := p.CurrentToken.Type; t {
	case "RESET":
		{
			p.matchToken("RESET")
			return &Reset{}
		}

	case "PRINT":
		{
			p.matchToken("PRINT")
			n := p.parseExpression2()
			return &Print{
				Expr: n,
			}
		}
	case "TYPE":
		{

			ty := StringTypeMap[p.matchToken("TYPE")]
			iden := p.matchToken("IDENTIFIER")
			p.matchToken("=")
			n := p.parseExpression2()
			return &Assignment{
				Type:       ty,
				Identifier: iden,
				Expr:       n,
			}
		}
	}
	panic(fmt.Sprintf("Unknown type %s", p.CurrentToken.Type))
	return nil
}

func (p *Parser) parseExpression2() Node {
	var op Op
	lhs := p.parseTerm2()
	be := &Binary2{
		Op:  NOOP,
		Lhs: lhs,
	}
	var be2 *Binary2
	for p.CurrentToken.Type == "+" || p.CurrentToken.Type == "-" {
		if p.CurrentToken.Type == "+" {
			p.matchToken("+")
			op = PLUS
		} else {
			p.matchToken("-")
			op = MINUS
		}

		rhs := p.parseTerm2()
		be.Op = op
		be.Rhs = rhs

		if p.CurrentToken.Type == "+" || p.CurrentToken.Type == "-" {
			be2 = &Binary2{Op: NOOP}
			be2.Lhs = be
			be = be2
		}
	}
	return be
}

func (p *Parser) parseTerm2() Node {
	var op Op
	lhs := p.parsePower()
	be := &Binary2{
		Op:  NOOP,
		Lhs: lhs,
	}
	var be2 *Binary2
	for p.CurrentToken.Type == "*" || p.CurrentToken.Type == "/" {
		if p.CurrentToken.Type == "*" {
			p.matchToken("*")
			op = MULTIPLY
		} else {
			p.matchToken("/")
			op = DIVIDE
		}

		rhs := p.parsePower()
		be.Op = op
		be.Rhs = rhs

		if p.CurrentToken.Type == "*" || p.CurrentToken.Type == "/" {
			be2 = &Binary2{Op: NOOP}
			be2.Lhs = be
			be = be2
		}

	}
	return be
}

func (p *Parser) parsePower() Node {

	lhs := p.parseFactor()

	se := &Binary2{
		Op:  NOOP,
		Lhs: lhs,
	}
	for p.CurrentToken.Type == "**" {
		p.matchToken("**")

		rhs := p.parsePower()
		se.Op = POWER
		se.Rhs = rhs
	}

	return se
}

func (p *Parser) parseFactor() Node {
	switch c := p.CurrentToken.Type; c {
	case "IDENTIFIER":
		{
			val := p.matchToken("IDENTIFIER")
			return &Identifier{
				Val: val,
			}
		}
	case "DECIMAL":
		{
			flt, err := strconv.ParseFloat(p.matchToken("DECIMAL"), 64)
			if err != nil {
				panic(err)
			}
			return &Number{
				Type: FLOAT,
				Flt:  flt,
			}
		}
	case "NUMBER":
		{
			num, err := strconv.Atoi(p.matchToken("NUMBER"))
			if err != nil {
				panic(err)
			}
			return &Number{
				Type: INT,
				Num:  num,
			}
		}
	case "(":
		{
			p.matchToken("(")
			p.parseExpression2()
			p.matchToken(")")
		}
	}

	panic(fmt.Sprintf("Unkown type %s", p.CurrentToken.Type))
	return nil
}
