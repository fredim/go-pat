package p

import (
	"testing"
)

func TestNode(t *testing.T) {
	p := NewSimpleParseNode(0, "null")
	if p.Type != 0 || string(p.Value) != "null" {
		t.Error("invalid parse node p")
	}
	if p.Len() != 0 {
		t.Error("shouldn't be any children, yet")
	}
	x := NewParseNode(1, []byte("left"))
	if x.Type != 1 || string(x.Value) != "left" {
		t.Error("invalid parse node x")
	}
	p.Push(x)
	if p.At(0).Type != 1 {
		t.Error("didn't push right")
	}
	p = NewSimpleParseNode(0, "null") // reset p
	y := NewSimpleParseNode(2, "right")
	p.PushTwo(x, y)
	if p.Len() != 2 || p.At(1).Type != 2 {
		t.Error("didn't pushtwo")
	}
	p.Set(0, y)
	if p.At(0).Type != 2 {
		t.Error("didn't set right")
	}
	p.SetExpr(x)
	if p.Expr.Type != 1 {
		t.Error("expr not right")
	}
}
