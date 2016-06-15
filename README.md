# go-pat
pattern matcher for go

## import the match interface

	import "github.com/fredim/go-pat/m"

## build cases

	matchCases := m.BuildCases(pat1, func1, pat2, func2, ...)

## the arguments to the function depend on the exported values from the patterns

	"string() as s", func(s string) string { return "got a string!" },

## primitives are simple tokens with values in parantheses

	"int(123)" matches an integer 123

## structs are tokens with members in braces

	"MyStruct{Member:value}"

## Useful Members need to be Exported (Capitalized) to be able to use in funcs

	"MyStruct{Member:val}", func(val MemberType) string { ... }

## value can be a sub structure, as well

	"MyStruct{Member:InnerStruct{}}"

## the as keyword assigns a argument name to the match

	"MyStruct{Member:InnerStruct{}} as ms", func(ms MyStruct) string { .. }

## all funcs need to have the same return value

	var matchFn func() string

## Finally, execute the match, and provide a func pointer

	matchCases.Match(obj, &matchFn)

## Match will reassign the matchFn to the return value, so it can be invoked without arguments

	return matchFn()

## done!
