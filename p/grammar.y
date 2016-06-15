%{
package p

func SetParseTree(yylex interface{}, root *Node) {
	tn := yylex.(*Tokenizer)
	tn.ParseTree = root
}
%}


// fields inside this union end up as the fields in a structure known
// as ${PREFIX}SymType, of which a reference is passed to the lexer.
%union{
	node *Node
	nodes []*Node
}

// any non-terminal which returns a value needs a type, which is
// really a field name in the above union struct

// same for terminals
%token <node> STRING INTEGER NUMBER TRUE FALSE
%token <node> IDTOKEN TERMINATOR PACKAGE STRUCT ARRAY
%token <node> MEMBER PRIMITIVE

%right <node> AS '*'
%right <node> '(' '[' '{'
%left <node> ')' ']' '}'
%left <node> RANGE '.' DOUBLECOLON ','

%start Root

// Fake Tokens
%token <node> COMMENT LEX_ERROR

%type <node> Query Object Literal
%type <node> PackagedName PtrPackagedName Primitive
%type <node> Member MemberList

%%

Root
	: /* empty, do something? */
	| Query
		{ SetParseTree(yylex, $1) }

Query
	: Literal
		{ $$ = $1 }
	| Object
		{ $$ = $1 }
	| Object AS IDTOKEN
		{ $$ = $1.SetExport(string($3.Value)) }

Literal
	: IDTOKEN
		{ $$ = $1 }
	| STRING
		{ $$ = $1 }
	| INTEGER
		{ $$ = $1 }
	| NUMBER
		{ $$ = $1 }

Object
	: Primitive '(' ')'
		{ $$ = $1 }
	| Primitive '(' Literal ')'
		{ $$ = $1.SetExpr($3) }
	| PtrPackagedName '{' '}'
		{ $$ = $1 }
	| PtrPackagedName '{' MemberList '}'
		{ $$ = $1.SetExpr($3) }

MemberList
	: Member
		{ $$ = NewSimpleParseNode(ARRAY, "MemberList").Push($1) }
	| MemberList ',' Member
		{ $$ = $1.Push($3) }

Member
	: IDTOKEN ':' Query
		{ $$ = $1.SetType(MEMBER).SetExpr($3) }

Primitive
	: IDTOKEN
		{ $$ = $1.SetType(PRIMITIVE) }

PtrPackagedName
	: PackagedName
		{ $$ = $1 }
	| '*' PackagedName
		{ $$ = $1.Push($2) }

PackagedName
	: IDTOKEN
		{ $$ = $1.SetType(STRUCT) }
	| IDTOKEN '.' IDTOKEN
		{ $$ = $2.PushTwo($1.SetType(PACKAGE), $3.SetType(STRUCT)) }
