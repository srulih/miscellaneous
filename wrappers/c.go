package main

import (
	"os"
	"fmt"
	"runtime"
	"text/template"
)

const tmpl = 

`{{range . }}
Line Number {{ .LineNum }} in File {{ .File }} was called
{{ end}}
`

type FileInfo struct {
	File string
	LineNum int
}

type C struct {
	S string
	L int
	StackCalls []FileInfo
}


func (c *C) Exclaim() {
	ok := true
	var file string
	var line int
	for i := 0; ok; i++ {
		_, file, line, ok = runtime.Caller(i)
		if ok {
			c.StackCalls = append(c.StackCalls, FileInfo{ File: file, LineNum: line })
		}
	}
}

func (c *C) FormatStackCalls() {
	tmpl, err := template.New("template").Parse(tmpl)
	if err != nil {
		return
	}
	err = tmpl.Execute(os.Stdout, c.StackCalls)
	if err != nil {
		fmt.Println(err)
	}
}
