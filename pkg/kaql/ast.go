package kaql

type Node interface {
	Token() Token
}

type Ident struct {
	Name string
	pos  int
}

func NewIdent(name string, pos int) *Ident {
	return &Ident{Name: name, pos: pos}
}

func (i *Ident) Token() Token { return IDENT }

type Const struct {
	Value string
	tok   Token
	pos   int
}

func NewConst(value string, pos int, tok Token) *Const {
	return &Const{Value: value, pos: pos, tok: tok}
}

func (c *Const) Token() Token { return c.tok }

type BinaryExpr struct {
	X  Node
	Y  Node
	Op Token
}

func NewBinaryExpr(operator Token, x, y Node) *BinaryExpr {
	return &BinaryExpr{Op: operator, X: x, Y: y}
}

func (e *BinaryExpr) Token() Token { return e.Op }

type UnaryExpr struct {
	X  Node
	Op Token
}

func NewUnaryExpr(operator Token, x Node) *UnaryExpr {
	return &UnaryExpr{Op: operator, X: x}
}

func (e *UnaryExpr) Token() Token { return e.Op }
