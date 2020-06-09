package main

import (
	"fmt"
	"go/format"
	"log"
	"os"
	"strings"
)

type astGenerator struct {
	sb *strings.Builder
}

func (ag *astGenerator) defineAst(baseName string, types []string) {
	ag.sb.WriteString("package typedcalculator\n\n")

	ag.defineNodeInterface()

	ag.OpGenerator()
	visitorTypes := make([]string, 0)

	for _, ty := range types {
		rs := strings.Split(ty, ":")
		r := strings.TrimSpace(rs[0])
		visitorTypes = append(visitorTypes, r)

		fts := make([]string, 0)
		fs := strings.Split(rs[1], ",")
		for i, _ := range fs {
			if fs[i] == " " {
				continue
			}
			fs[i] = strings.TrimSpace(fs[i])
			ev := strings.Split(fs[i], " ")
			for j, _ := range ev {
				ev[j] = strings.TrimSpace(ev[j])
			}
			fts = append(fts, ev...)
		}
		ag.defineType(baseName, r, fts)
		ag.defineVisitor(baseName, r)
	}

	ag.generateNodeType(visitorTypes)
	ag.generateVisitorInterface(baseName, visitorTypes)

}

func (ag *astGenerator) defineType(baseType string, typeName string, fields []string) {
	ag.sb.WriteString("type ")
	ag.sb.WriteString(typeName)
	ag.sb.WriteString(" struct {\n")

	// write the fields

	for i := 0; i < len(fields); i += 2 {
		ag.sb.WriteString(fields[i])
		ag.sb.WriteString("\t")
		ag.sb.WriteString(fields[i+1])
		ag.sb.WriteString("\n")
	}

	ag.sb.WriteString("\n}\n\n")
}

func (ag *astGenerator) OpGenerator() {
	ops := []string{"ILLEGALOP", "NOOP", "PLUS", "MINUS", "MULTIPLY", "DIVIDE", "POWER"}
	ag.generateConstEnum("Op", ops)

	types := []string{"NOTYPE", "FLOAT", "DOUBLE", "INT", "LONG"}
	ag.generateConstEnum("Type", types)
}

func (ag *astGenerator) generateConstEnum(c string, vals []string) {
	ag.sb.WriteString(fmt.Sprintf("type %s int\n\n", c))
	ag.sb.WriteString("const (\n")
	ag.sb.WriteString(fmt.Sprintf("%s %s = iota\n", vals[0], c))
	for i := 1; i < len(vals); i++ {
		ag.sb.WriteString(fmt.Sprintf("%s\n", vals[i]))
	}
	ag.sb.WriteString(")\n\n")

	//generate map of enum to string
	ag.sb.WriteString(fmt.Sprintf("var %sStringMap = map[%s]string {\n", c, c))
	for i := 1; i < len(vals); i++ {
		ag.sb.WriteString(fmt.Sprintf("%s: \"%s\",\n", vals[i], vals[i]))
	}
	ag.sb.WriteString("}\n\n")

	//generate map of string to enum
	ag.sb.WriteString(fmt.Sprintf("var String%sMap = map[string]%s {\n", c, c))
	for i := 1; i < len(vals); i++ {
		ag.sb.WriteString(fmt.Sprintf("\"%s\":%s,\n", vals[i], vals[i]))
	}
	ag.sb.WriteString("}\n\n")
}

func (ag *astGenerator) defineNodeInterface() {
	ag.sb.WriteString("type Node interface {\n")
	ag.sb.WriteString("isNode()\n")
	ag.sb.WriteString("accept(Visitor)}\n\n")
}

func (ag *astGenerator) defineVisitor(baseType string, typeName string) {
	ag.sb.WriteString(fmt.Sprintf("func (f *%s) accept(v Visitor) {\n", typeName))
	ag.sb.WriteString(fmt.Sprintf("\tv.visit%s(f)\n}\n\n", typeName+baseType))
}

func (ag *astGenerator) generateNodeType(types []string) {
	for _, ty := range types {
		ag.sb.WriteString(fmt.Sprintf("func (f *%s) isNode() {}\n", ty))
	}
	ag.sb.WriteString("\n")
}

func (ag *astGenerator) generateVisitorInterface(baseType string, types []string) {
	ag.sb.WriteString("type Visitor interface {\n")
	for _, ty := range types {
		ag.sb.WriteString(fmt.Sprintf("visit%s(f *%s) \n", ty+baseType, ty))
	}
	ag.sb.WriteString("}\n\n")
}

func (ag astGenerator) String() string {
	return ag.sb.String()
}

func (ag astGenerator) FormatFile() []byte {
	sb, err := format.Source([]byte(ag.String()))
	if err != nil {
		log.Fatal(err)
	}
	return sb
}

func (ag astGenerator) WriteToFile(file string) {
	f, err := os.Create(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = f.Write(ag.FormatFile())
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	ag := astGenerator{sb: &strings.Builder{}}
	ag.defineAst("Stmt", []string{
		"Program   : Lines []*Line",
		"Line     : Stmt Node",
		"Assignment : Type Type , Identifier string, Expr Node",
		"Print : Expr Node",
		"Reset : ",
		"Binary2 : Type Type, Op Op, Lhs Node, Rhs Node",
		"Identifier : Val string",
		"Number : Type Type, Fixed bool, Num int, Flt float64",
	})

	ag.WriteToFile("../ast_tree.go")

}
