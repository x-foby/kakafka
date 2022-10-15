package kaql

import "fmt"

type Token int

const (
	ILLEGAL Token = iota
	EOF

	literalbeg
	IDENT  // x
	NUMBER // 1.23
	STRING // "abc"
	literalend

	operatorsbeg
	AND // &&
	OR  // ||

	EQL // ==
	NEQ // !=
	LSS // <
	LEQ // <=
	GTR // >
	GEQ // >=
	operatorsend

	LPAREN // (
	RPAREN // )

	MINUS // -
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",

	IDENT:  "IDENT",
	NUMBER: "NUMBER",
	STRING: "STRING",

	AND: "&",
	OR:  "|",

	EQL: "=",
	NEQ: "!=",
	LSS: "<",
	LEQ: "<=",
	GTR: ">",
	GEQ: ">=",

	LPAREN: "(",
	RPAREN: ")",

	MINUS: "-",
}

func (t Token) String() string {
	if t < ILLEGAL || t >= Token(len(tokens)) {
		return fmt.Sprintf("undefined (%v)", int(t))
	}

	return tokens[t]
}

func (t Token) Precedence() int {
	switch t {
	case OR:
		return 1
	case AND:
		return 2
	case EQL, NEQ, LSS, LEQ, GTR, GEQ:
		return 3

	default:
		return 0
	}
}

func (t Token) IsLiteral() bool { return literalbeg < t && t < literalend }

func (t Token) IsOperator() bool { return operatorsbeg < t && t < operatorsend }
