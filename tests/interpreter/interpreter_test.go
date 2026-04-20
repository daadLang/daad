package interpreter_test

import (
	"testing"

	"github.com/daadLang/daad/internals/ast"
	"github.com/daadLang/daad/internals/interpreter"
	"github.com/daadLang/daad/internals/lexer"
)

// Helper to extract raw value from interpreter.Value for comparison
func getRawValue(v interpreter.Value) interface{} {
	switch val := v.(type) {
	case interpreter.IntValue:
		return val.V
	case interpreter.FloatValue:
		return val.V
	case interpreter.StringValue:
		return val.V
	case interpreter.BoolValue:
		return val.V
	case interpreter.CharValue:
		return val.V
	case interpreter.NoneValue:
		return nil
	default:
		return v
	}
}

// Test basic arithmetic operations
func TestBasicArithmetic(t *testing.T) {
	tests := []struct {
		name     string
		module   *ast.Module
		varName  string
		expected interface{}
	}{
		{
			name: "addition",
			module: &ast.Module{
				Body: []ast.Stmt{
					&ast.ExprStmt{
						Value: &ast.Assign{
							Target: &ast.Name{Id: "متغير"},
							Value: &ast.BinOp{
								Left:  &ast.Constant{Value: 10},
								Op:    lexer.PLUS,
								Right: &ast.Constant{Value: 5},
							},
						},
					},
				},
			},
			varName:  "متغير",
			expected: 15,
		},
		{
			name: "multiplication and subtraction",
			module: &ast.Module{
				Body: []ast.Stmt{
					&ast.ExprStmt{
						Value: &ast.Assign{
							Target: &ast.Name{Id: "متغير"},
							Value:  &ast.Constant{Value: 10},
						},
					},
					&ast.ExprStmt{
						Value: &ast.Assign{
							Target: &ast.Name{Id: "نتيجة"},
							Value: &ast.BinOp{
								Left: &ast.BinOp{
									Left:  &ast.Name{Id: "متغير"},
									Op:    lexer.MULT,
									Right: &ast.Constant{Value: 2},
								},
								Op:    lexer.MINUS,
								Right: &ast.Constant{Value: 3},
							},
						},
					},
				},
			},
			varName:  "نتيجة",
			expected: 17,
		},
		{
			name: "power",
			module: &ast.Module{
				Body: []ast.Stmt{
					&ast.ExprStmt{
						Value: &ast.Assign{
							Target: &ast.Name{Id: "قوة"},
							Value: &ast.BinOp{
								Left:  &ast.Constant{Value: 2},
								Op:    lexer.POWER,
								Right: &ast.Constant{Value: 8},
							},
						},
					},
				},
			},
			varName:  "قوة",
			expected: 256,
		},
		{
			name: "division",
			module: &ast.Module{
				Body: []ast.Stmt{
					&ast.ExprStmt{
						Value: &ast.Assign{
							Target: &ast.Name{Id: "قسمة"},
							Value: &ast.BinOp{
								Left:  &ast.Constant{Value: 100},
								Op:    lexer.DIVIDE,
								Right: &ast.Constant{Value: 4},
							},
						},
					},
				},
			},
			varName:  "قسمة",
			expected: 25.0,
		},
		{
			name: "modulo",
			module: &ast.Module{
				Body: []ast.Stmt{
					&ast.ExprStmt{
						Value: &ast.Assign{
							Target: &ast.Name{Id: "باقي"},
							Value: &ast.BinOp{
								Left:  &ast.Constant{Value: 17},
								Op:    lexer.MOD,
								Right: &ast.Constant{Value: 5},
							},
						},
					},
				},
			},
			varName:  "باقي",
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			interp := interpreter.NewInterpreter()
			interp.Run(tt.module)
			result := getRawValue(interp.GetVar(tt.varName))
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// Test comparison and boolean operations
func TestComparisonOperators(t *testing.T) {
	tests := []struct {
		name     string
		module   *ast.Module
		varName  string
		expected bool
	}{
		{
			name: "equality",
			module: &ast.Module{
				Body: []ast.Stmt{
					&ast.ExprStmt{
						Value: &ast.Assign{
							Target: &ast.Name{Id: "نتيجة1"},
							Value: &ast.Compare{
								Left:       &ast.Constant{Value: 5},
								Op:         lexer.EQ,
								Comparator: &ast.Constant{Value: 5},
							},
						},
					},
				},
			},
			varName:  "نتيجة1",
			expected: true,
		},
		{
			name: "not equal",
			module: &ast.Module{
				Body: []ast.Stmt{
					&ast.ExprStmt{
						Value: &ast.Assign{
							Target: &ast.Name{Id: "نتيجة2"},
							Value: &ast.Compare{
								Left:       &ast.Constant{Value: 10},
								Op:         lexer.NEQ,
								Comparator: &ast.Constant{Value: 5},
							},
						},
					},
				},
			},
			varName:  "نتيجة2",
			expected: true,
		},
		{
			name: "less than",
			module: &ast.Module{
				Body: []ast.Stmt{
					&ast.ExprStmt{
						Value: &ast.Assign{
							Target: &ast.Name{Id: "نتيجة3"},
							Value: &ast.Compare{
								Left:       &ast.Constant{Value: 7},
								Op:         lexer.LESS,
								Comparator: &ast.Constant{Value: 10},
							},
						},
					},
				},
			},
			varName:  "نتيجة3",
			expected: true,
		},
		{
			name: "greater than",
			module: &ast.Module{
				Body: []ast.Stmt{
					&ast.ExprStmt{
						Value: &ast.Assign{
							Target: &ast.Name{Id: "نتيجة4"},
							Value: &ast.Compare{
								Left:       &ast.Constant{Value: 15},
								Op:         lexer.GREATER,
								Comparator: &ast.Constant{Value: 10},
							},
						},
					},
				},
			},
			varName:  "نتيجة4",
			expected: true,
		},
		{
			name: "less than or equal",
			module: &ast.Module{
				Body: []ast.Stmt{
					&ast.ExprStmt{
						Value: &ast.Assign{
							Target: &ast.Name{Id: "نتيجة5"},
							Value: &ast.Compare{
								Left:       &ast.Constant{Value: 5},
								Op:         lexer.LEQ,
								Comparator: &ast.Constant{Value: 5},
							},
						},
					},
				},
			},
			varName:  "نتيجة5",
			expected: true,
		},
		{
			name: "greater than or equal",
			module: &ast.Module{
				Body: []ast.Stmt{
					&ast.ExprStmt{
						Value: &ast.Assign{
							Target: &ast.Name{Id: "نتيجة6"},
							Value: &ast.Compare{
								Left:       &ast.Constant{Value: 20},
								Op:         lexer.GEQ,
								Comparator: &ast.Constant{Value: 15},
							},
						},
					},
				},
			},
			varName:  "نتيجة6",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			interp := interpreter.NewInterpreter()
			interp.Run(tt.module)
			result := getRawValue(interp.GetVar(tt.varName))
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// Test boolean operations
func TestBooleanOperators(t *testing.T) {
	tests := []struct {
		name     string
		expr     ast.Expr
		expected bool
	}{
		{
			name: "and true false",
			expr: &ast.BoolOp{
				Left:  &ast.Constant{Value: true},
				Op:    lexer.AND,
				Right: &ast.Constant{Value: false},
			},
			expected: false,
		},
		{
			name: "or true false",
			expr: &ast.BoolOp{
				Left:  &ast.Constant{Value: true},
				Op:    lexer.OR,
				Right: &ast.Constant{Value: false},
			},
			expected: true,
		},
		{
			name: "not true",
			expr: &ast.UnaryOp{
				Op:   lexer.NOT,
				Expr: &ast.Constant{Value: true},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			interp := interpreter.NewInterpreter()
			module := &ast.Module{
				Body: []ast.Stmt{
					&ast.ExprStmt{Value: tt.expr},
				},
			}
			interp.Run(module)
			// Note: This test assumes execExpr is accessible or we store result
			// For now, we'll just run it to ensure no panic
		})
	}
}

// Test string operations
func TestStringOperations(t *testing.T) {
	tests := []struct {
		name     string
		module   *ast.Module
		varName  string
		expected string
	}{
		{
			name: "string assignment",
			module: &ast.Module{
				Body: []ast.Stmt{
					&ast.ExprStmt{
						Value: &ast.Assign{
							Target: &ast.Name{Id: "اسم"},
							Value:  &ast.Constant{Value: "محمد"},
						},
					},
				},
			},
			varName:  "اسم",
			expected: "محمد",
		},
		{
			name: "string concatenation",
			module: &ast.Module{
				Body: []ast.Stmt{
					&ast.ExprStmt{
						Value: &ast.Assign{
							Target: &ast.Name{Id: "رسالة"},
							Value: &ast.BinOp{
								Left:  &ast.Constant{Value: "مرحبا "},
								Op:    lexer.PLUS,
								Right: &ast.Constant{Value: "بك"},
							},
						},
					},
				},
			},
			varName:  "رسالة",
			expected: "مرحبا بك",
		},
		{
			name: "empty string",
			module: &ast.Module{
				Body: []ast.Stmt{
					&ast.ExprStmt{
						Value: &ast.Assign{
							Target: &ast.Name{Id: "فارغ"},
							Value:  &ast.Constant{Value: ""},
						},
					},
				},
			},
			varName:  "فارغ",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			interp := interpreter.NewInterpreter()
			interp.Run(tt.module)
			result := getRawValue(interp.GetVar(tt.varName))
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// Test control flow - if/else
func TestIfStatement(t *testing.T) {
	module := &ast.Module{
		Body: []ast.Stmt{
			&ast.ExprStmt{
				Value: &ast.Assign{
					Target: &ast.Name{Id: "العدد"},
					Value:  &ast.Constant{Value: 15},
				},
			},
			&ast.IfStmt{
				Test: &ast.Compare{
					Left:       &ast.Name{Id: "العدد"},
					Op:         lexer.GREATER,
					Comparator: &ast.Constant{Value: 10},
				},
				Body: []ast.Stmt{
					&ast.ExprStmt{
						Value: &ast.Assign{
							Target: &ast.Name{Id: "نتيجة"},
							Value:  &ast.Constant{Value: "كبير"},
						},
					},
				},
				Orelse: []ast.Stmt{
					&ast.ExprStmt{
						Value: &ast.Assign{
							Target: &ast.Name{Id: "نتيجة"},
							Value:  &ast.Constant{Value: "صغير"},
						},
					},
				},
			},
		},
	}

	interp := interpreter.NewInterpreter()
	interp.Run(module)
	result := getRawValue(interp.GetVar("نتيجة"))
	if result != "كبير" {
		t.Errorf("expected كبير, got %v", result)
	}
}

// Test augmented assignment
func TestAugmentedAssignment(t *testing.T) {
	tests := []struct {
		name     string
		op       lexer.TokenType
		initial  int
		operand  int
		expected interface{}
	}{
		{"plus assign", lexer.PLUS_ASSIGN, 10, 5, 15},
		{"minus assign", lexer.MINUS_ASSIGN, 10, 3, 7},
		{"mult assign", lexer.MULT_ASSIGN, 10, 2, 20},
		{"divide assign", lexer.DIVIDE_ASSIGN, 20, 4, 5.0},
		{"mod assign", lexer.MOD_ASSIGN, 17, 5, 2},
		{"power assign", lexer.POWER_ASSIGN, 2, 3, 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			module := &ast.Module{
				Body: []ast.Stmt{
					&ast.ExprStmt{
						Value: &ast.Assign{
							Target: &ast.Name{Id: "عدد"},
							Value:  &ast.Constant{Value: tt.initial},
						},
					},
					&ast.AugmentedAssignStmt{
						Target: &ast.Name{Id: "عدد"},
						Op:     tt.op,
						Value:  &ast.Constant{Value: tt.operand},
					},
				},
			}

			interp := interpreter.NewInterpreter()
			interp.Run(module)
			result := getRawValue(interp.GetVar("عدد"))
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestBasicOOPClassInstantiationAndMethod(t *testing.T) {
	module := &ast.Module{
		Body: []ast.Stmt{
			&ast.ClassDefStmt{
				Name: "شخص",
				Body: []ast.Stmt{
					&ast.FunctionDefStmt{
						Name: "__بناء__",
						Args: []string{"ذاتي", "اسم"},
						Body: []ast.Stmt{
							&ast.ExprStmt{
								Value: &ast.Assign{
									Target: &ast.Attribute{
										Value: &ast.Name{Id: "ذاتي"},
										Attr:  "الاسم",
									},
									Value: &ast.Name{Id: "اسم"},
								},
							},
						},
					},
					&ast.FunctionDefStmt{
						Name: "خذ_الاسم",
						Args: []string{"ذاتي"},
						Body: []ast.Stmt{
							&ast.ReturnStmt{
								Value: &ast.Attribute{
									Value: &ast.Name{Id: "ذاتي"},
									Attr:  "الاسم",
								},
							},
						},
					},
				},
			},
			&ast.ExprStmt{
				Value: &ast.Assign{
					Target: &ast.Name{Id: "م"},
					Value: &ast.Call{
						Func: &ast.Name{Id: "شخص"},
						Args: []ast.Expr{&ast.Constant{Value: "أحمد"}},
					},
				},
			},
			&ast.ExprStmt{
				Value: &ast.Assign{
					Target: &ast.Name{Id: "اسم_م"},
					Value: &ast.Call{
						Func: &ast.Attribute{
							Value: &ast.Name{Id: "م"},
							Attr:  "خذ_الاسم",
						},
					},
				},
			},
			&ast.ExprStmt{
				Value: &ast.Assign{
					Target: &ast.Name{Id: "اسم_مباشر"},
					Value: &ast.Attribute{
						Value: &ast.Name{Id: "م"},
						Attr:  "الاسم",
					},
				},
			},
		},
	}

	interp := interpreter.NewInterpreter()
	interp.Run(module)

	nameByMethod := getRawValue(interp.GetVar("اسم_م"))
	if nameByMethod != "أحمد" {
		t.Fatalf("expected method result أحمد, got %v", nameByMethod)
	}

	nameDirect := getRawValue(interp.GetVar("اسم_مباشر"))
	if nameDirect != "أحمد" {
		t.Fatalf("expected direct attribute أحمد, got %v", nameDirect)
	}
}
