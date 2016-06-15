package p

import (
	"testing"
)

func TestNewTokenizer(t *testing.T) {
	tk := NewStringTokenizer("")
	if tk.In != "" {
		t.Error("invalid tokenizer>")
	}

	if escapeDecodeMap['Z'] != 26 {
		t.Error("invalid decode map")
	}

	out := &yySymType{}
	x := tk.Lex(out)
	if x != 0 {
		t.Error("invalid response")
	}

	tk.Error("simple error")
	if tk.LastError != "Syntax Error at position 2 after token  :: \n\n" {
		t.Error("invalid error", tk.LastError)
	}

	tk = NewStringTokenizer(" \n\r\t")
	tk.skipBlank()
	if tk.lastChar != 0 {
		t.Error("skip blank failed")
	}

	letters := "@_ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	for _, l := range letters {
		if !isLetter(uint16(l)) {
			t.Error(l, "was not a letter!")
		}
	}

	numbers := "0123456789"
	for _, n := range numbers {
		if !isDigit(uint16(n)) {
			t.Error(n, "was not a number!")
		}
	}
}

func checkToken(name string, n *Node, ty int, t *testing.T) {
	if n.Type != ty {
		t.Error(name, "Failed token", n.Value, ". Was expecting ", ty, " but got ", n.Type)
	}
}

func TestTokenizerScanIdentifier(t *testing.T) {
	// out := &yySymType{}

	tk := NewStringTokenizer("as")
	tk.Next() // roll to the first char
	checkToken("TestTokenizerScanIdentifier, AS: ", tk.scanIdentifier(), AS, t)

	tk = NewStringTokenizer("a")
	tk.Next()
	checkToken("TestTokenizerScanIdentifier, Token: ", tk.scanIdentifier(), IDTOKEN, t)
}

func TestTokenizerScanNumber(t *testing.T) {
	// out := &yySymType{}

	tk := NewStringTokenizer("123")
	tk.Next() // roll to the first char
	checkToken("TestTokenizerScanNumber 123: ", tk.scanNumber(false), INTEGER, t)

	tk = NewStringTokenizer("0123")
	tk.Next()
	checkToken("TestTokenizerScanNumber 0123: ", tk.scanNumber(false), INTEGER, t)

	tk = NewStringTokenizer("0129a")
	tk.Next()
	checkToken("TestTokenizerScanNumber 0129a: ", tk.scanNumber(false), LEX_ERROR, t)

	tk = NewStringTokenizer("0x123abc")
	tk.Next()
	checkToken("TestTokenizerScanNumber 0x123abc: ", tk.scanNumber(false), INTEGER, t)

	tk = NewStringTokenizer("0x123ABC")
	tk.Next()
	checkToken("TestTokenizerScanNumber 0x123ABC: ", tk.scanNumber(false), INTEGER, t)

	tk = NewStringTokenizer("1.23")
	tk.Next()
	checkToken("TestTokenizerScanNumber 1.23: ", tk.scanNumber(false), NUMBER, t)

	tk = NewStringTokenizer("1.23.")
	tk.Next()
	checkToken("TestTokenizerScanNumber 1.23.: ", tk.scanNumber(false), NUMBER, t)

	tk = NewStringTokenizer("12e3")
	tk.Next()
	checkToken("TestTokenizerScanNumber 12e3: ", tk.scanNumber(false), INTEGER, t)

	tk = NewStringTokenizer("1.23e4")
	tk.Next()
	checkToken("TestTokenizerScanNumber 1.23: ", tk.scanNumber(false), NUMBER, t)
	// TODO: skips invalid hex chars?
	// tk = NewStringTokenizer("0x12gh")
	// tk.Next()
	// n = tk.scanNumber(false)
	// if n.Type != LEX_ERROR {
	// 	t.Error("should have caught non-hex char", n.Type, n.Value)
	// }
}

func TestTokenizerScanString(t *testing.T) {
	tk := NewStringTokenizer("123\"")
	tk.Next() // roll to the first char
	checkToken("TestTokenizerScanString 123: ", tk.scanString('"'), STRING, t)

	tk = NewStringTokenizer("12\\n\"")
	tk.Next()
	checkToken("TestTokenizerScanString 123\\n: ", tk.scanString('"'), STRING, t)

}

func TestTokenizerScan(t *testing.T) {
	// out := &yySymType{}

	tk := NewStringTokenizer("a")
	checkToken("TestTokenizerScan a:", tk.Scan(), IDTOKEN, t)

	tk = NewStringTokenizer("1")
	checkToken("TestTokenizerScan 1:", tk.Scan(), INTEGER, t)
}
