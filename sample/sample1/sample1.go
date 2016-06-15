package sample1

import (
	"fmt"
	"strconv"

	"github.com/fredim/go-pat/m"
)

type Expr interface {
	isExpr()
}

type ExprBase struct{}

func (e *ExprBase) isExpr() {}

type ExprInt interface {
	Expr
	isIntExpr()
}

type ExprIntBase struct {
	ExprBase
}

func (e *ExprIntBase) isIntExpr() {}

type IntLiteral struct {
	ExprIntBase
	Val int
}

type OpAddInt struct {
	ExprIntBase
	Lhs, Rhs ExprInt
}

type OpNeg struct {
	ExprIntBase
	Val ExprInt
}

func transInt(n ExprInt) string {
	var matchFn func() string
	m.Match(n, m.BuildCases(
		"*IntLiteral{Val:i}", func(i int) string {
			return strconv.Itoa(i)
		},
		"*OpAddInt{Rhs:*OpNeg{Val:r}} as oai", func(oai *OpAddInt, r ExprInt) string {
			return transInt(oai.Lhs) + "-" + transInt(r)
		},
		"*OpAddInt{Lhs:l, Rhs:r}", func(l ExprInt, r ExprInt) string {
			return transInt(l) + "+" + transInt(r)
		},
		"_", func() string { return "failed match" },
	), &matchFn)
	return matchFn()
}

func Sample() {
	lit1 := &IntLiteral{Val: 1}
	lit2 := &IntLiteral{Val: 2}
	add1 := &OpAddInt{Lhs: lit1, Rhs: lit2}
	neg1 := &OpAddInt{Lhs: lit1, Rhs: &OpNeg{Val: lit2}}
	fmt.Println("lit1", transInt(lit1))
	fmt.Println("lit2", transInt(lit2))
	fmt.Println("add1", transInt(add1))
	fmt.Println("neg1", transInt(neg1))

}
