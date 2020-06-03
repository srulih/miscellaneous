package calculator

import (
	"fmt"
	"regexp"
	"strings"
)

func MapValues(m [][2]string) []string {
	v := make([]string, len(m))
	for i := 0; i < len(m); i++ {
		v[i] = m[i][0]
	}
	return v
}

func CreateRegType(m [][2]string) []RegType {
	regTypes := make([]RegType, len(m))

	for i := 0; i < len(m); i++ {
		regTypes[i] = RegType{Regex: m[i][0], Type: m[i][1]}
	}
	return regTypes
}

type Token struct {
	Pos   int
	Value string
	Type  string
}

type RegType struct {
	Regex string
	Type  string
}

func (t Token) String() string {
	return fmt.Sprintf("%s with value %s at position %d,", t.Type, t.Value, t.Pos)
}

type Lexer struct {
	rules           []RegType
	regex           *regexp.Regexp
	whitespaceRegex *regexp.Regexp
	skipWhitespace  bool
	buffer          string
	pos             int
}

func CreateLexer(rules [][2]string, skipWhitespace bool) *Lexer {
	rs := strings.Join(MapValues(rules), "|")
	r := regexp.MustCompile(rs)

	// newline is not included in the whitespace
	rw, _ := regexp.Compile(`[^\t\f\r ]`)
	return &Lexer{
		skipWhitespace:  skipWhitespace,
		rules:           CreateRegType(rules),
		regex:           r,
		whitespaceRegex: rw,
	}
}

func (l *Lexer) Input(buf string) {
	l.buffer = buf
}

func (l *Lexer) Token() (*Token, error) {
	if l.pos >= len(l.buffer) {
		return &Token{Pos: l.pos, Type: "EOF"}, nil
	}
	if l.skipWhitespace {
		m := l.whitespaceRegex.FindAllStringIndex(l.buffer[l.pos:], 1)
		if len(m) > 0 {
			l.pos += m[0][0]
		}
	}
	r := l.regex.FindAllStringIndex(l.buffer[l.pos:], 1)
	if len(r) > 0 && r[0][0] == 0 {
		value := l.buffer[l.pos : l.pos+r[0][1]]
		var ty string
		for _, rt := range l.rules {
			if m, _ := regexp.MatchString(rt.Regex, value); m {
				ty = rt.Type
				break
			}
		}
		t := Token{
			Value: value,
			Pos:   l.pos,
			Type:  ty,
		}
		l.pos += r[0][1]
		return &t, nil
	}
	return nil, fmt.Errorf("could not match anything at position %d", l.pos)
}

func (l *Lexer) Tokens() []*Token {
	tokens := make([]*Token, 0)
	for {
		tok, err := l.Token()
		if err != nil {
			panic(err)
		}
		if tok.Type == "EOF" {
			break
		}
		tokens = append(tokens, tok)
	}
	return tokens
}

func (l *Lexer) Reset() {
	l.pos = 0
}

func main() {
	m := [][2]string{
		{`\d+`, "NUMBER"},
		{`[a-zA-Z_]\w+`, "IDENTIFIER"},
		{`\+`, "PLUS"},
	}
	l := CreateLexer(m, true)
	l.Input("hello 23 + 34    ft")
	toks := l.Tokens()
	for _, v := range toks {
		fmt.Println(v)
	}
}
