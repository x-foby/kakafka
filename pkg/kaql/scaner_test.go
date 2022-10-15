package kaql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestScan(t *testing.T) {
	cases := []struct {
		name string
		src  string
		pos  int
		tok  Token
		lit  string
	}{
		{
			name: "EOF",
			src:  "",
			pos:  0,
			tok:  EOF,
		},
		{
			name: "identifier",
			src:  " ident ",
			pos:  1,
			tok:  IDENT,
			lit:  "ident",
		},
		{
			name: "integer",
			src:  "10 ",
			pos:  0,
			tok:  NUMBER,
			lit:  "10",
		},
		{
			name: "negative integer",
			src:  "-10 ",
			pos:  0,
			tok:  NUMBER,
			lit:  "-10",
		},
		{
			name: "float",
			src:  " 1.1",
			pos:  1,
			tok:  NUMBER,
			lit:  "1.1",
		},
		{
			name: "negative float",
			src:  " -1.1",
			pos:  1,
			tok:  NUMBER,
			lit:  "-1.1",
		},
		{
			name: "float without leading symbol",
			src:  ".1",
			pos:  0,
			tok:  NUMBER,
			lit:  ".1",
		},
		{
			name: "negative float without leading symbol",
			src:  "-.1",
			pos:  0,
			tok:  NUMBER,
			lit:  "-.1",
		},
		{
			name: "string",
			src:  `"string"`,
			pos:  0,
			tok:  STRING,
			lit:  "string",
		},
		{
			name: "and",
			src:  " && ",
			pos:  1,
			tok:  AND,
		},
		{
			name: "or",
			src:  "||",
			pos:  0,
			tok:  OR,
		},
		{
			name: "equal",
			src:  "==",
			pos:  0,
			tok:  EQL,
		},
		{
			name: "not equal",
			src:  "!=",
			pos:  0,
			tok:  NEQ,
		},
		{
			name: "lss",
			src:  "<",
			pos:  0,
			tok:  LSS,
		},
		{
			name: "leq",
			src:  "<=",
			pos:  0,
			tok:  LEQ,
		},
		{
			name: "gtr",
			src:  ">",
			pos:  0,
			tok:  GTR,
		},
		{
			name: "geq",
			src:  ">=",
			pos:  0,
			tok:  GEQ,
		},
		{
			name: "l paren",
			src:  "(",
			pos:  0,
			tok:  LPAREN,
		},
		{
			name: "r paren",
			src:  ")",
			pos:  0,
			tok:  RPAREN,
		},
	}

	var s Scanner

	for _, cc := range cases {
		t.Run(cc.name, func(t *testing.T) {
			s.Init([]rune(cc.src))

			pos, tok, lit := s.Scan()

			require.Equal(t, cc.pos, pos)
			require.Equal(t, cc.tok, tok)
			require.Equal(t, cc.lit, lit)
		})
	}
}

func TestPeek(t *testing.T) {
	cases := []struct {
		name string
		src  string
		ch   rune
	}{
		{
			name: "illegal",
			src:  "",
			ch:   -1,
		},
		{
			name: "illegal",
			src:  " ",
			ch:   -1,
		},
		{
			name: "normal",
			src:  " a",
			ch:   'a',
		},
	}

	var s Scanner

	for _, cc := range cases {
		t.Run(cc.name, func(t *testing.T) {
			s.Init([]rune(cc.src))

			ch := s.peek()

			require.Equal(t, cc.ch, ch)
		})
	}
}
