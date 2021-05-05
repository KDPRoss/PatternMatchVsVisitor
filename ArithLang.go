package main

import "fmt"

// To hell with the normal Golang 'everything returns an error' thing; we'll
// just `panic`.

type Type interface{}

type Typable interface {
	Type() Type
}

type Evalable interface {
	Eval() Exp
}

type Exp interface {
	Typable
	Evalable
}

type Num struct {
	val int
}

type TNum struct{}

func (Num) Type() Type {
	return TNum{}
}

func (e Num) Eval() Exp {
	return e
}

type Bool struct {
	val bool
}

type TBool struct{}

func (Bool) Type() Type {
	return TBool{}
}

func (e Bool) Eval() Exp {
	return e
}

type Add struct {
	e1 Exp
	e2 Exp
}

func (e Add) Type() Type {
	switch e.e1.Type().(type) {
	case TNum:
		switch e.e2.Type().(type) {
		case TNum:
			return TNum{}
		default:
			panic("Fail")
		}
	default:
		panic("Fail")
	}
}

func (e Add) Eval() Exp {
	v1 := e.e1.Eval()
	v2 := e.e2.Eval()
	return Num{v1.(Num).val + v2.(Num).val}
}

type IsZ struct {
	e Exp
}

func (e IsZ) Type() Type {
	switch e.e.Type().(type) {
	case TNum:
		return TBool{}
	default:
		panic("Fail")
	}
}

func (e IsZ) Eval() Exp {
	return Bool{e.e.Eval().(Num).val == 0}
}

type If struct {
	e1 Exp
	e2 Exp
	e3 Exp
}

func (e If) Type() Type {
	switch e.e1.Type().(type) {
	case TBool:
		t2 := e.e2.Type()
		t3 := e.e3.Type()
		if t2 == t3 {
			return t2
		}
		panic("Fail")
	default:
		panic("Fail")
	}
}

func (e If) Eval() Exp {
	if e.e1.Eval().(Bool).val {
		return e.e2.Eval()
	}
	return e.e3.Eval()
}

type ShittyMaybe struct {
	valid bool
	e     Exp
}

type TestCase struct {
	in  Exp
	out ShittyMaybe
}

func Just(e Exp) ShittyMaybe {
	return ShittyMaybe{true, e}
}

var Nothing = ShittyMaybe{valid: false}

var Examples = []TestCase{
	{Num{3}, Just(Num{3})},
	{Bool{true}, Just(Bool{true})},
	{Add{Num{2}, Num{3}}, Just(Num{5})},
	{Add{Num{2}, Bool{true}}, Nothing},
	{IsZ{Num{0}}, Just(Bool{true})},
	{IsZ{Num{3}}, Just(Bool{false})},
	{If{Bool{true}, Num{0}, Num{1}}, Just(Num{0})},
	{If{IsZ{Num{0}}, Add{Num{1}, Num{1}}, Num{0}}, Just(Num{2})},
	{If{Num{0}, Num{1}, Num{2}}, Nothing},
}

func RunOne(c TestCase) {
	defer func() {
		if r := recover(); r != nil {
			if c.out.valid {
				panic("Bad test case!")
			}
		}
	}()

	e := c.in
	te := e.Type()
	v := e.Eval()
	if !c.out.valid || c.out.e != v {
		panic("Evaluation fail!")
	}
	tv := v.Type()
	if te != tv {
		panic("Type-preservation fail!")
	}
}

func main() {
	for _, e := range Examples {
		RunOne(e)
	}

	fmt.Println("The coast is clear!")
}
