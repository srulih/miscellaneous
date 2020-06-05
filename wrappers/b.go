package main

type B struct {
	C C
}

func (b *B) Exclaim() {
	b.C.Exclaim()
}
