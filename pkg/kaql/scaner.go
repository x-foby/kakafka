package kaql

type Scanner struct {
	src    []rune
	len    int
	ch     rune
	offset int
}

func (s *Scanner) Init(src []rune) {
	s.src = src
	s.len = len(s.src)
	s.offset = -1
	s.next()
}

func (s *Scanner) next() {
	if s.offset < s.len-1 {
		s.offset++
		s.ch = s.src[s.offset]
	} else {
		s.offset = s.len
		s.ch = -1 // eof
	}
}

func (s *Scanner) peek() rune {
	if s.offset < s.len-1 {
		return s.src[s.offset+1]
	}

	return -1 // eof
}

func (s *Scanner) skipWhitespace() {
	for s.ch == ' ' {
		s.next()
	}
}

func (s *Scanner) scanIdentifier() string {
	offs := s.offset

	// TODO: scan any identifiers
	for isLetter(s.ch) || isDigit(s.ch) || s.ch == '.' {
		s.next()
	}

	return string(s.src[offs:s.offset])
}

func (s *Scanner) scanNumber() (Token, string) {
	offs := s.offset
	tok := ILLEGAL

	var isFloat bool

	for isDigit(s.ch) || (s.ch == '.' && !isFloat) {
		if s.ch == '.' {
			isFloat = true
		}

		tok = NUMBER

		s.next()
	}

	return tok, string(s.src[offs:s.offset])
}

func (s *Scanner) scanString() (int, Token, string) {
	offs := s.offset

	for s.ch != '"' || s.offset == offs {
		if s.ch == -1 {
			return s.offset, ILLEGAL, ""
		}

		s.next()
	}

	return offs, STRING, string(s.src[offs+1 : s.offset])
}

func isLetter(ch rune) bool {
	return (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') || ch == '_'
}

func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func (s *Scanner) Scan() (pos int, tok Token, lit string) {
	tok = ILLEGAL
	s.skipWhitespace()

	pos = s.offset

	switch ch := s.ch; {
	case isLetter(ch):
		lit = s.scanIdentifier()
		tok = IDENT
	case isDigit(ch),
		ch == '.' && isDigit(s.peek()):
		tok, lit = s.scanNumber()
	default:
		switch ch {
		case -1:
			tok = EOF
		case '-':
			tok = MINUS
		case '"':
			pos, tok, lit = s.scanString()
		case '(':
			tok = LPAREN
		case ')':
			tok = RPAREN
		case '<':
			if s.peek() == '=' {
				s.next()
				tok = LEQ
			} else {
				tok = LSS
			}
		case '>':
			if s.peek() == '=' {
				s.next()
				tok = GEQ
			} else {
				tok = GTR
			}
		case '=':
			if s.peek() == '=' {
				s.next()
				tok = EQL
			}
		case '!':
			if s.peek() == '=' {
				s.next()
				tok = NEQ
			}
		case '&':
			if s.peek() == '&' {
				s.next()
				tok = AND
			}
		case '|':
			if s.peek() == '|' {
				s.next()
				tok = OR
			}
		}

		s.next()
	}

	return
}
