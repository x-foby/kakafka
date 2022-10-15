package kaql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	cases := []struct {
		name     string
		src      string
		expected Node
		err      error
	}{
		// positive
		{
			name: "empty",
			src:  "   ",
		},
		{
			name: "equal (string)",
			src:  `property == "value"`,
			expected: &BinaryExpr{
				X: &Ident{
					Name: "property",
					pos:  0,
				},
				Op: EQL,
				Y: &Const{
					Value: "value",
					pos:   12,
					tok:   STRING,
				},
			},
		},
		{
			name: "equal (int)",
			src:  "property == 1",
			expected: &BinaryExpr{
				X: &Ident{
					Name: "property",
					pos:  0,
				},
				Op: EQL,
				Y: &Const{
					Value: "1",
					pos:   12,
					tok:   NUMBER,
				},
			},
		},
		{
			name: "equal (float)",
			src:  "property == 1.1",
			expected: &BinaryExpr{
				X: &Ident{
					Name: "property",
					pos:  0,
				},
				Op: EQL,
				Y: &Const{
					Value: "1.1",
					pos:   12,
					tok:   NUMBER,
				},
			},
		},
		{
			name: "equal (float without leading symbol)",
			src:  "property == .1",
			expected: &BinaryExpr{
				X: &Ident{
					Name: "property",
					pos:  0,
				},
				Op: EQL,
				Y: &Const{
					Value: ".1",
					pos:   12,
					tok:   NUMBER,
				},
			},
		},
		{
			name: "equal (true)",
			src:  "property == true",
			expected: &BinaryExpr{
				X: &Ident{
					Name: "property",
					pos:  0,
				},
				Op: EQL,
				Y: &Ident{
					Name: "true",
					pos:  12,
				},
			},
		},
		{
			name: "equal (false)",
			src:  "property == false",
			expected: &BinaryExpr{
				X: &Ident{
					Name: "property",
					pos:  0,
				},
				Op: EQL,
				Y: &Ident{
					Name: "false",
					pos:  12,
				},
			},
		},
		{
			name: "equal (null)",
			src:  "property == null",
			expected: &BinaryExpr{
				X: &Ident{
					Name: "property",
					pos:  0,
				},
				Op: EQL,
				Y: &Ident{
					Name: "null",
					pos:  12,
				},
			},
		},
		// {
		// 	name: "equal undefined",
		// 	src:  "property == undefined",
		// 	expected: &BinaryExpr{
		// 		X: &Ident{
		// 			Name: "property",
		// 			pos:  0,
		// 		},
		// 		Op: EQL,
		// 		Y: &Ident{
		// 			Name: "undefined",
		// 			pos:  12,
		// 		},
		// 	},
		// },
		{
			name: "and",
			src:  "property == true && property == false",
			expected: &BinaryExpr{
				X: &BinaryExpr{
					X: &Ident{
						Name: "property",
						pos:  0,
					},
					Op: EQL,
					Y: &Ident{
						Name: "true",
						pos:  12,
					},
				},
				Op: AND,
				Y: &BinaryExpr{
					X: &Ident{
						Name: "property",
						pos:  20,
					},
					Op: EQL,
					Y: &Ident{
						Name: "false",
						pos:  32,
					},
				},
			},
		},
		{
			name: "or",
			src:  "property == true || property == false",
			expected: &BinaryExpr{
				X: &BinaryExpr{
					X: &Ident{
						Name: "property",
						pos:  0,
					},
					Op: EQL,
					Y: &Ident{
						Name: "true",
						pos:  12,
					},
				},
				Op: OR,
				Y: &BinaryExpr{
					X: &Ident{
						Name: "property",
						pos:  20,
					},
					Op: EQL,
					Y: &Ident{
						Name: "false",
						pos:  32,
					},
				},
			},
		},
		{
			name: "precedence",
			src:  "property == 1 || property == 2 && property == 3",
			expected: &BinaryExpr{
				X: &BinaryExpr{
					X: &Ident{
						Name: "property",
						pos:  0,
					},
					Op: EQL,
					Y: &Const{
						Value: "1",
						pos:   12,
						tok:   NUMBER,
					},
				},
				Op: OR,
				Y: &BinaryExpr{
					X: &BinaryExpr{
						X: &Ident{
							Name: "property",
							pos:  17,
						},
						Op: EQL,
						Y: &Const{
							Value: "2",
							pos:   29,
							tok:   NUMBER,
						},
					},
					Op: AND,
					Y: &BinaryExpr{
						X: &Ident{
							Name: "property",
							pos:  34,
						},
						Op: EQL,
						Y: &Const{
							Value: "3",
							pos:   46,
							tok:   NUMBER,
						},
					},
				},
			},
		},
		{
			name: "parens",
			src:  "(property == 1 || property == 2) && property == 3",
			expected: &BinaryExpr{
				X: &BinaryExpr{
					X: &BinaryExpr{
						X: &Ident{
							Name: "property",
							pos:  1,
						},
						Op: EQL,
						Y: &Const{
							Value: "1",
							pos:   13,
							tok:   NUMBER,
						},
					},
					Op: OR,
					Y: &BinaryExpr{
						X: &Ident{
							Name: "property",
							pos:  18,
						},
						Op: EQL,
						Y: &Const{
							Value: "2",
							pos:   30,
							tok:   NUMBER,
						},
					},
				},
				Op: AND,
				Y: &BinaryExpr{
					X: &Ident{
						Name: "property",
						pos:  36,
					},
					Op: EQL,
					Y: &Const{
						Value: "3",
						pos:   48,
						tok:   NUMBER,
					},
				},
			},
		},
		{
			name: "not equal",
			src:  `property != "value"`,
			expected: &BinaryExpr{
				X: &Ident{
					Name: "property",
					pos:  0,
				},
				Op: NEQ,
				Y: &Const{
					Value: "value",
					pos:   12,
					tok:   STRING,
				},
			},
		},
		{
			name: "gtr",
			src:  "property > 1",
			expected: &BinaryExpr{
				X: &Ident{
					Name: "property",
					pos:  0,
				},
				Op: GTR,
				Y: &Const{
					Value: "1",
					pos:   11,
					tok:   NUMBER,
				},
			},
		},
		{
			name: "geq",
			src:  "property >= 1",
			expected: &BinaryExpr{
				X: &Ident{
					Name: "property",
					pos:  0,
				},
				Op: GEQ,
				Y: &Const{
					Value: "1",
					pos:   12,
					tok:   NUMBER,
				},
			},
		},
		{
			name: "lss",
			src:  "property < 1",
			expected: &BinaryExpr{
				X: &Ident{
					Name: "property",
					pos:  0,
				},
				Op: LSS,
				Y: &Const{
					Value: "1",
					pos:   11,
					tok:   NUMBER,
				},
			},
		},
		{
			name: "leq",
			src:  "property <= 1",
			expected: &BinaryExpr{
				X: &Ident{
					Name: "property",
					pos:  0,
				},
				Op: LEQ,
				Y: &Const{
					Value: "1",
					pos:   12,
					tok:   NUMBER,
				},
			},
		},
		{
			name: "children property",
			src:  "property.children == 1",
			expected: &BinaryExpr{
				X: &Ident{
					Name: "property.children",
					pos:  0,
				},
				Op: EQL,
				Y: &Const{
					Value: "1",
					pos:   21,
					tok:   NUMBER,
				},
			},
		},
		// {
		// 	name:     "nested",
		// 	src:      `property{children == 1 && children == 2}`,
		// 	expected: &Query{},
		// },
		// {
		// 	name: "",
		// 	src:  "",
		// 	err: &ParsingError{
		// 		message: "qwe",
		// 		pos:     1,
		// 	},
		// },
	}

	var p Parser

	for _, cc := range cases {
		t.Run(cc.name, func(t *testing.T) {
			node, err := p.Parse(cc.src)

			if cc.err != nil {
				require.ErrorAs(t, err, &cc.err)
				require.EqualError(t, cc.err, err.Error())
				require.Empty(t, node)
			} else {
				require.NoError(t, err)
				require.Equal(t, cc.expected, node)
			}
		})
	}
}
