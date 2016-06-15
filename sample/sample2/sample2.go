package sample2

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

var transIntMatches *m.MatchInfo

func init() {
	transIntMatches = m.BuildCases(
		"*IntLiteral{} as il", func(il *IntLiteral) string {
			return strconv.Itoa(il.Val)
		},
		"*OpAddInt{Rhs:*OpNeg{Val:r}} as oai", func(oai *OpAddInt, r ExprInt) string {
			return transInt(oai.Lhs) + "-" + transInt(r)
		},
		"*OpAddInt{} as oai", func(oai *OpAddInt) string {
			return transInt(oai.Lhs) + "+" + transInt(oai.Rhs)
		},
		"_", func() string { return "failed match" },
	)
}

func transInt(n ExprInt) string {
	var matchFn func() string
	transIntMatches.Match(n, &matchFn)
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
