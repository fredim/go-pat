//go:generate go tool yacc -o grammar.go grammar.y
package p

import (
	"errors"
)

func Parse(s string) (*Node, error) {
	tokenizer := NewStringTokenizer(s)
	if yyParse(tokenizer) != 0 {
		return nil, errors.New(tokenizer.LastError)
	}
	if tokenizer.ParseTree == nil {
		return nil, errors.New("Empty query")
	}
	return tokenizer.ParseTree, nil
}

/* simple Nodes are used for expressions */

type Node struct {
	Type   int
	Value  []byte
	Sub    []*Node
	Expr   *Node
	Export string
}

func NewSimpleParseNode(nodeType int, value string) *Node {
	return &Node{Type: nodeType, Value: []byte(value)}
}

func NewParseNode(nodeType int, value []byte) *Node {
	return &Node{Type: nodeType, Value: value}
}

func (self *Node) PushTwo(left *Node, right *Node) *Node {
	self.Push(left)
	return self.Push(right)
}

func (self *Node) Push(value *Node) *Node {
	if self.Sub == nil {
		self.Sub = make([]*Node, 0, 2)
	}
	self.Sub = append(self.Sub, value)
	return self
}

func (self *Node) Pop() *Node {
	self.Sub = self.Sub[:len(self.Sub)-1]
	return self
}

func (self *Node) At(index int) *Node {
	return self.Sub[index]
}

func (self *Node) Set(index int, val *Node) {
	self.Sub[index] = val
}

func (self *Node) SetExpr(expr *Node) *Node {
	self.Expr = expr
	return self
}

func (self *Node) SetExport(name string) *Node {
	self.Export = name
	return self
}

func (self *Node) SetType(t int) *Node {
	self.Type = t
	return self
}

func (self *Node) Len() int {
	return len(self.Sub)
}
