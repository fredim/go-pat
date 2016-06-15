package m

import (
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	"github.com/fredim/go-pat/p"
)

var _ = fmt.Println // dummy to keep from having to remove fmt

type MapVar map[string]interface{}

type MatchInfo struct {
	callerPackage, callerFunction string
	patterns                      []*PatResult
}

type PatResult struct {
	pat  string
	node *p.Node
	fptr interface{} // function type
}

// Cases take inputs in pairs
// match: action
// input[0]: input[1]
func BuildCases(inputs ...interface{}) *MatchInfo {
	callerPkg, callerName := getCallerName()
	ret := []*PatResult{}
	l := len(inputs)
	for i := 1; i < l; i += 2 {
		pp := inputs[i-1].(string)
		node, err := p.Parse(pp)
		if err != nil {
			panic(err)
		}
		ret = append(ret, &PatResult{pp, node, inputs[i]})
	}
	return &MatchInfo{
		callerPkg, callerName, ret,
	}
}

var (
	PrimitiveMap = map[string]reflect.Kind{
		"string": reflect.String,
		"ptr":    reflect.Ptr,
	}
	StructMap = map[string]reflect.Kind{
		"struct": reflect.Struct,
	}
)

func (mi *MatchInfo) nodeMatch(n *p.Node, obj reflect.Value, objType reflect.Type, inArgs []reflect.Value) ([]reflect.Value, bool) {
	args := inArgs
	if n.Export != "" {
		args = append(inArgs, obj)
	}
	// fmt.Println("matching", n.Type, objType.Kind())
	switch n.Type {
	case p.STRING:
		if objType.Kind() == reflect.String {
			nval := string(n.Value)
			if obj.String() == nval {
				return args, true
			}
		}
	case p.INTEGER:
		if objType.Kind() == reflect.Int {
			ival, _ := strconv.ParseInt(string(n.Value), 10, 64)
			if ival == obj.Int() {
				return args, true
			}
		}
	case p.NUMBER:
		if objType.Kind() == reflect.Float64 {
			fval, _ := strconv.ParseFloat(string(n.Value), 64)
			if fval == obj.Float() {
				return args, true
			}
		}
	case p.IDTOKEN:
		// this means stick the obj into the arg list
		// unless it starts with an _
		if n.Value[0] != '_' {
			args = append(args, obj)
		}
		return args, true
	case p.PRIMITIVE:
		nval := string(n.Value)
		// fmt.Println("Matching a PRIMITIVE: ", nval)
		if k, ok := PrimitiveMap[nval]; ok {
			if k == objType.Kind() {
				if n.Expr != nil {
					return mi.nodeMatch(n.Expr, obj, objType, args)
				}
				return args, true
			}
		}
	case p.STRUCT:
		nval := string(n.Value)
		// fmt.Println("Matching a STRUCT: ", nval, objType.Kind())
		switch objType.Kind() {
		case reflect.Struct:
			if nval != objType.Name() {
				return args, false
			}
			// fmt.Println("what's in the struct?", n)
			if n.Expr == nil {
				return args, true
			}
			// fmt.Println("struct members: ", n.Expr.Type)
			memberMatch := true
			for _, member := range n.Expr.Sub {
				// fmt.Println("struct members: ", member)
				if member.Type == p.MEMBER {
					fn := string(member.Value)
					field := obj.FieldByName(fn)
					if field.Kind() == reflect.Interface {
						field = field.Elem()
					}
					args, memberMatch = mi.nodeMatch(member.Expr, field, field.Type(), args)
					if !memberMatch {
						return args, false
					}
				} else {
					return args, false
				}
			}
			return args, memberMatch
		case reflect.Ptr:
			return mi.nodeMatch(n, obj.Elem(), objType.Elem(), inArgs)
		}
	case '*':
		switch objType.Kind() {
		case reflect.Ptr:
			ptrTo := n.At(0)
			if n.Expr != nil {
				ptrTo.SetExpr(n.Expr)
			}
			return mi.nodeMatch(ptrTo, obj.Elem(), objType.Elem(), args)
		case reflect.Struct:
			ptrTo := n.At(0)
			return mi.nodeMatch(ptrTo, obj, objType, args)
		}
	case '.':
		// TODO: match package name?
		pack := n.At(0)
		if pack.Type != p.PACKAGE {
			return args, false
		}
		pval := string(pack.Value)
		if objType.PkgPath() != pval {
			return args, false
		}
		return mi.nodeMatch(n.At(1), obj, objType, args)
	}
	return args, false
}

func (mi *MatchInfo) Match(obj interface{}, fptr interface{}) {
	objType := reflect.TypeOf(obj)
	// fmt.Println("---------------starting match")
	for _, pr := range mi.patterns {
		// fmt.Println(pr.pat)
		if pr.node == nil {
			return
		}
		args := []reflect.Value{}
		if args, ok := mi.nodeMatch(pr.node, reflect.ValueOf(obj), objType, args); ok {
			// fmt.Println("creating a custom fn", args)
			matchFn := func(in []reflect.Value) []reflect.Value {
				fn := reflect.ValueOf(pr.fptr)
				return fn.Call(args)
			}
			fn := reflect.ValueOf(fptr).Elem()
			// Make a function of the right type.
			v := reflect.MakeFunc(fn.Type(), matchFn)
			// Assign it to the value fn represents.
			fn.Set(v)
			return
			// } else {
			// 	fmt.Println("failed match", pr.pat, args, obj)
		}
	}
}

func Match(obj interface{}, mi *MatchInfo, fptr interface{}) {
	mi.Match(obj, fptr)
}

func getCallerName() (string, string) {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(3, pc)
	f := runtime.FuncForPC(pc[0])
	names := strings.Split(f.Name(), ".")
	l := len(names) - 1
	return strings.Join(names[0:l], "."), names[l]
}
