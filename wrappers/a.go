package main

type A struct {
	B B
}

func (a *A) Exclaim() {
	a.B.Exclaim()
}


func main() {
	a := A{B: B{C: C{S: "hello", StackCalls: make([]FileInfo, 0)}}}
	a.Exclaim()
	a.B.C.FormatStackCalls()

	var fi []FileInfo
	for i:=0; i < 3; i++ {
		fi = append(fi, FileInfo{})
	}


}
