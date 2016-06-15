package p

import (
	"bytes"
	"errors"
	"fmt"
)

// Special marker, safe because input is read byte by byte
const EOFCHAR = 0x100

const contextLen = 14

type Tokenizer struct {
	In            string
	AllowComments bool
	lastChar      uint16
	position      int
	lastToken     *Node
	LastError     string
	ParseTree     *Node
}

func NewStringTokenizer(s string) *Tokenizer {
	return &Tokenizer{In: s, position: -1}
}

var keywords = map[string]int{
	"as": AS,
}

// escapEncodeMap specifies how to escape certain binary data with '\'
// complies to http://dev.mysql.com/doc/refman/5.1/en/string-syntax.html
var escapeEncodeMap = map[byte]byte{
	'\x00': '0',
	'\'':   '\'',
	'"':    '"',
	'\b':   'b',
	'\n':   'n',
	'\r':   'r',
	'\t':   't',
	26:     'Z', // ctl-Z
	'\\':   '\\',
}

// escapeDecodeMap is the reverse of excapeEncodeMap
var escapeDecodeMap map[byte]byte

func init() {
	escapeDecodeMap = make(map[byte]byte)
	for k, v := range escapeEncodeMap {
		escapeDecodeMap[v] = k
	}
}

func (self *Tokenizer) Lex(lval *yySymType) int {
	parseNode := self.Scan()
	for parseNode.Type == COMMENT {
		if self.AllowComments {
			break
		}
		parseNode = self.Scan()
	}
	self.lastToken = parseNode
	lval.node = parseNode
	return parseNode.Type
}

func (self *Tokenizer) Error(err string) {
	// errors.New("Error %c %d", self.lastChar, self.position)
	context := ""
	start := 0
	if self.position-contextLen > 0 {
		start = self.position - contextLen
	}
	end := len(self.In)
	if self.position+contextLen < len(self.In) {
		end = self.position + contextLen
	}
	for i, v := range self.In[start:end] {
		if i == contextLen-1 {
			context += "⇾"
			context += string(v)
			context += "⇽"
		} else {
			context += string(v)
		}
	}
	// increment position as want position of failure, not last good position
	self.LastError = fmt.Sprintf("Syntax Error at position %v after token %s :: \n%s\n", self.position+1, string(self.lastToken.Value), context)
	// errors.New("LastError %s", self.LastError)
}

func (self *Tokenizer) Scan() (parseNode *Node) {
	defer func() {
		if x := recover(); x != nil {
			err := x.(error)
			parseNode = NewSimpleParseNode(LEX_ERROR, err.Error())
		}
	}()

	if self.lastChar == 0 {
		self.Next()
	}
	self.skipBlank()
	switch ch := self.lastChar; {
	case isLetter(ch):
		return self.scanIdentifier()
	case isDigit(ch):
		return self.scanNumber(false)
	default:
		self.Next()
		switch ch {
		case EOFCHAR: // TODO: Que es?
			return NewSimpleParseNode(0, "")
		case '*', '.', ',', ':', '(', ')', '[', ']', '{', '}':
			return NewSimpleParseNode(int(ch), string(ch))
		// case '.':
		// 	if isDigit(self.lastChar) {
		// 		return self.scanNumber(true)
		// 	} else if self.lastChar == '.' {
		// 		self.Next()
		// 		return NewSimpleParseNode(RANGE, "..")
		// 	} else {
		// 		return NewSimpleParseNode(int(ch), ".")
		// 	}
		// case '\'':
		// 	tok := self.scanString(ch)
		// 	tok.Type = IDTOKEN
		// 	return tok
		case '"':
			return self.scanString(ch)
		default:
			return NewSimpleParseNode(LEX_ERROR, fmt.Sprintf("Unexpected character '%c'", ch))
		}
	}
}

func (self *Tokenizer) skipBlank() {
	ch := self.lastChar
	for ch == ' ' || ch == '\n' || ch == '\r' || ch == '\t' {
		self.Next()
		ch = self.lastChar
	}
}

func (self *Tokenizer) scanIdentifier() *Node {
	buffer := bytes.NewBuffer(make([]byte, 0, 8))
	buffer.WriteByte(byte(self.lastChar))
	for self.Next(); isLetter(self.lastChar) || isDigit(self.lastChar); self.Next() {
		buffer.WriteByte(byte(self.lastChar))
	}
	id := buffer.String()
	if keywordId, found := keywords[id]; found {
		return NewParseNode(keywordId, buffer.Bytes())
	}
	return NewParseNode(IDTOKEN, buffer.Bytes())
}

func (self *Tokenizer) scanMantissa(base int, buffer *bytes.Buffer) {
	for digitVal(self.lastChar) < base {
		self.ConsumeNext(buffer)
	}
}

func (self *Tokenizer) scanNumber(seenDecimalPoint bool) *Node {
	buffer := bytes.NewBuffer(make([]byte, 0, 8))
	nodeType := INTEGER
	if seenDecimalPoint {
		nodeType = NUMBER
		buffer.WriteByte(byte('.'))
		self.scanMantissa(10, buffer)
		goto exponent
	}

	if self.lastChar == '0' {
		// int or float
		self.ConsumeNext(buffer)
		if self.lastChar == 'x' || self.lastChar == 'X' {
			// hexadecimal int
			self.ConsumeNext(buffer)
			self.scanMantissa(16, buffer)
		} else {
			// octal int or float
			seenDecimalDigit := false
			self.scanMantissa(8, buffer)
			if self.lastChar == '8' || self.lastChar == '9' {
				// illegal octal int or float
				seenDecimalDigit = true
				self.scanMantissa(10, buffer)
			}
			if self.lastChar == '.' || self.lastChar == 'e' || self.lastChar == 'E' {
				goto fraction
			}
			// octal int
			if seenDecimalDigit {
				return NewParseNode(LEX_ERROR, buffer.Bytes())
			}
		}
		goto exit
	}

	// decimal int or float
	self.scanMantissa(10, buffer)

fraction:
	if self.lastChar == '.' {
		nextChar := self.Peek(1)
		if nextChar == '.' {
			// Could be a RANGE, return INTEGER number and keep parsing
			goto exit
		}
		if !isDigit(nextChar) {
			// We don't support trailing dot in number literals
			goto exit
		}
		nodeType = NUMBER
		self.ConsumeNext(buffer)
		self.scanMantissa(10, buffer)
	}

exponent:
	if self.lastChar == 'e' || self.lastChar == 'E' {
		self.ConsumeNext(buffer)
		if self.lastChar == '+' || self.lastChar == '-' {
			self.ConsumeNext(buffer)
		}
		self.scanMantissa(10, buffer)
	}

exit:
	return NewParseNode(nodeType, buffer.Bytes())
}

func (self *Tokenizer) scanString(delim uint16) *Node {
	buffer := bytes.NewBuffer(make([]byte, 0, 8))
	for {
		ch := self.lastChar
		self.Next()
		if ch == delim {
			if self.lastChar == delim {
				self.Next()
			} else {
				break
			}
		} else if ch == '\\' {
			if self.lastChar == EOFCHAR {
				return NewParseNode(LEX_ERROR, buffer.Bytes())
			}
			if decodedChar, ok := escapeDecodeMap[byte(self.lastChar)]; ok {
				ch = uint16(decodedChar)
			} else {
				ch = self.lastChar
			}
			self.Next()
		}
		if ch == EOFCHAR {
			return NewParseNode(LEX_ERROR, buffer.Bytes())
		}
		buffer.WriteByte(byte(ch))
	}
	return NewParseNode(STRING, buffer.Bytes())
}

func (self *Tokenizer) ConsumeNext(buffer *bytes.Buffer) {
	// Never consume an EOF
	if self.lastChar == EOFCHAR {
		panic(errors.New("Unexpected EOF"))
	}
	buffer.WriteByte(byte(self.lastChar))
	self.Next()
}

func (self *Tokenizer) Peek(offset int) uint16 {
	if self.position+offset >= len(self.In) {
		return EOFCHAR
	}
	return uint16(self.In[self.position+offset])
}

func (self *Tokenizer) Next() {
	self.lastChar = self.Peek(1)
	self.position++
}

func isLetter(ch uint16) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch == '@'
}

func digitVal(ch uint16) int {
	switch {
	case '0' <= ch && ch <= '9':
		return int(ch) - '0'
	case 'a' <= ch && ch <= 'f':
		return int(ch) - 'a' + 10
	case 'A' <= ch && ch <= 'F':
		return int(ch) - 'A' + 10
	}
	return 16 // larger than any legal digit val
}

func isDigit(ch uint16) bool {
	return '0' <= ch && ch <= '9'
}
