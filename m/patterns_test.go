package m

import (
	"testing"
)

type Aifc interface {
	Tryme()
}

type Astr struct {
	i int
	s string
}

func (a *Astr) Tryme() {
}

type Bstr struct {
	I int
	S string
}

func (b *Bstr) Tryme() {
}

func TestLiterals(t *testing.T) {
	var matchFn func() string
	matchTester1 := BuildCases(
		`"test"`, func() string { return "string" },
		"32", func() string { return "integer" },
		"32.2", func() string { return "float" },
	)

	matchTester1.Match("test", &matchFn)
	t.Log("got a", matchFn())

	matchTester1.Match(32, &matchFn)
	t.Log("got a", matchFn())

	matchTester1.Match(32.2, &matchFn)
	t.Log("got a", matchFn())
}

func TestBasicPattern(t *testing.T) {
	var matchFn func() string
	matchTester1 := BuildCases(
		`string("test") as s`, func(s string) string { return s },
		"string() as s", func(s string) string { return s },
		"*Astr{} as a", func(a *Astr) string { return "member:" + a.s },
	)

	obj := &Astr{0, "test"}
	matchTester1.Match(obj.s, &matchFn)
	t.Log("got a", matchFn())

	matchTester1.Match(obj, &matchFn)
	t.Log("got a", matchFn())

	matchTester1.Match("generic", &matchFn)
	t.Log("got a", matchFn())

	matchTester2 := BuildCases(
		"*patterns.Astr{}", func() string { return "somethign else" },
		"_x", func() string { return "should match with no args" },
	)

	matchTester2.Match(obj, &matchFn)
	t.Log("got a", matchFn())

	matchTester2.Match(obj.s, &matchFn)
	t.Log("got a", matchFn())

	matchTester3 := BuildCases(
		// "x", func(x Aifc) string { return "got a matching interface" },
		"_", func() string { return "should match with no args" },
	)

	matchTester3.Match(obj, &matchFn)
	t.Log("got a", matchFn())

	matchTester3.Match(obj.s, &matchFn)
	t.Log("got a", matchFn())
}

func TestInnerPattern(t *testing.T) {
	var matchFn func() string
	matchTester1 := BuildCases(
		`*Astr{s:string("test")} as a`, func(a *Astr) string { return "Matched test" },
		`*Astr{s:"test2"} as a`, func(a *Astr) string { return "Matched test2" },
	)

	obj := &Astr{0, "test"}
	matchTester1.Match(obj, &matchFn)
	t.Log("got a", matchFn())

	obj = &Astr{0, "test2"}
	matchTester1.Match(obj, &matchFn)
	t.Log("got a", matchFn())

	matchTester2 := BuildCases(
		`*Bstr{S:string("test")} as b`, func(b *Bstr) string { return "Matched test" },
		"*Bstr{S:s} as b", func(b *Bstr, s string) string { return "Generic string " + s },
	)

	obj2 := &Bstr{0, "test"}
	matchTester2.Match(obj2, &matchFn)
	t.Log("got a", matchFn())

	obj2 = &Bstr{0, "something else"}
	matchTester2.Match(obj2, &matchFn)
	t.Log("got a", matchFn())
}
