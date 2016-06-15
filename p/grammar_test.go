package p

import (
	"testing"
)

func TestRoot(t *testing.T) {
	tk := NewStringTokenizer("")
	x := yyParse(tk)
	if x != 0 {
		t.Error("error parsing")
	}
}

func TestLiteral(t *testing.T) {
	tk := NewStringTokenizer("\"abc\"")
	if x := yyParse(tk); x != 0 {
		t.Error("error parsing")
	} else {
		n := tk.ParseTree
		if n.Type != STRING {
			t.Error("got the wrong thing", n)
		}
	}

	tk = NewStringTokenizer("123")
	if x := yyParse(tk); x != 0 {
		t.Error("error parsing")
	} else {
		n := tk.ParseTree
		if n.Type != INTEGER {
			t.Error("got the wrong thing", n)
		}
	}

	tk = NewStringTokenizer("123.3")
	if x := yyParse(tk); x != 0 {
		t.Error("error parsing")
	} else {
		n := tk.ParseTree
		if n.Type != NUMBER {
			t.Error("got the wrong thing", n)
		}
	}
}

func TestPrimitive(t *testing.T) {
	tk := NewStringTokenizer("a()")
	if x := yyParse(tk); x != 0 {
		t.Error("error parsing")
	} else {
		n := tk.ParseTree
		if n.Type != PRIMITIVE {
			t.Error("got the wrong thing", n)
		}
	}
}

func TestPtrPackagedName(t *testing.T) {
	tk := NewStringTokenizer("a{}")
	if x := yyParse(tk); x != 0 {
		t.Error("error parsing")
	} else {
		n := tk.ParseTree
		if n.Type != STRUCT {
			t.Error("got the wrong thing", n)
		}
	}

	tk = NewStringTokenizer("*a{}")
	if x := yyParse(tk); x != 0 {
		t.Error("error parsing")
	} else {
		n := tk.ParseTree
		if n.Type != '*' || n.At(0).Type != STRUCT {
			t.Error("got the wrong thing", n)
		}
	}

	tk = NewStringTokenizer("*b.c{}")
	if x := yyParse(tk); x != 0 {
		t.Error("error parsing", tk.LastError)
	} else {
		n := tk.ParseTree
		n1 := n.At(0)
		if n.Type != '*' || n1.Type != '.' || n1.At(0).Type != PACKAGE || n1.At(1).Type != STRUCT {
			t.Error("got the wrong thing", n)
		}
	}

	tk = NewStringTokenizer("a{b:b}")
	if x := yyParse(tk); x != 0 {
		t.Error("error parsing")
	} else {
		n := tk.ParseTree
		t.Log("inner struct", n)
		if n.Type != STRUCT || n.Expr == nil || n.Expr.Type != ARRAY {
			t.Error("got the wrong thing", n)
		}
	}

	tk = NewStringTokenizer("*a{b:b}")
	if x := yyParse(tk); x != 0 {
		t.Error("error parsing")
	} else {
		n := tk.ParseTree
		n1 := n.At(0)
		t.Log("inner struct", n)
		if n.Type != '*' || n1.Type != STRUCT {
			t.Error("got the wrong thing", n)
		}
	}
}
