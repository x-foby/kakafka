package kaql

import "fmt"

type Parser struct {
	scanner Scanner
	pos     int
	tok     Token
	lit     string
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(src string) (Node, error) {
	p.scanner.Init([]rune(src))

	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	// if p.tok != EOF {
	// 	return nil, p.unexpect()
	// }

	return expr, nil
}

func (p *Parser) next() {
	p.pos, p.tok, p.lit = p.scanner.Scan()
}

func (p *Parser) unexpect() error {
	return fmt.Errorf("unexpected %v at %v", p.tok, p.pos)
}

func (p *Parser) parseIdent() (*Ident, error) {
	if p.tok != IDENT {
		return nil, p.unexpect()
	}

	return NewIdent(p.lit, p.pos), nil
}

func (p *Parser) parseExpr() (Node, error) {
	x, err := p.parseNode()
	if err != nil {
		return nil, err
	}

	p.next()

	switch p.tok {
	case EOF, RPAREN:
		return x, nil
	}

	if p.tok.IsOperator() {
		return p.parseBinaryExpr(x)
	}

	return nil, p.unexpect()
}

func (p *Parser) parseNode() (Node, error) {
	p.next()

	switch p.tok {
	case EOF:
		return nil, nil
	case IDENT:
		return NewIdent(p.lit, p.pos), nil
	case NUMBER, STRING:
		return NewConst(p.lit, p.pos, p.tok), nil
	case MINUS:
		return p.parseUnaryExpr()
	case LPAREN:
		return p.parseParen()
	default:
		return nil, p.unexpect()
	}
}

func (p *Parser) parseUnaryExpr() (*UnaryExpr, error) {
	var op Token

	switch p.tok {
	case MINUS:
		op = MINUS
	default:
		return nil, p.unexpect()
	}

	x, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	return NewUnaryExpr(op, x), nil
}

func (p *Parser) parseParen() (Node, error) {
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	if p.tok != RPAREN {
		return nil, p.unexpect()
	}

	return expr, nil
}

func (p *Parser) parseBinaryExpr(x Node) (*BinaryExpr, error) {
	if !p.tok.IsOperator() {
		return nil, p.unexpect()
	}

	op := p.tok

	y, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	expr := NewBinaryExpr(op, x, y)

	if needSwap(expr, y) {
		swap(expr, y.(*BinaryExpr))
	}

	return expr, nil
}

func needSwap(x *BinaryExpr, y Node) bool {
	yBinaryExpr, ok := y.(*BinaryExpr)
	if !ok {
		return false
	}

	return x.Op.Precedence() > yBinaryExpr.Op.Precedence()
}

func swap(x, y *BinaryExpr) {
	*x = *NewBinaryExpr(y.Op, NewBinaryExpr(x.Op, x.X, y.X), y.Y)

	xx := x.X.(*BinaryExpr)
	if needSwap(xx, xx.Y) {
		swap(xx, xx.Y.(*BinaryExpr))
	}
}
