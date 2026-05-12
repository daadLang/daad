package interpreter

import (
	"os"
	"path/filepath"

	ast "github.com/daadLang/daad/internals/ast"
)

type Interpreter struct {
	env            *Env
	moduleCache    map[string]*ModuleValue
	loadingModules map[string]bool // detect circular imports
	sourceDir      string
	currentDir     string
	nativeModules  map[string]NativeModuleLoader
}

func NewInterpreter() *Interpreter {
	env := NewEnv(nil)
	RegisterBuiltins(env)

	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}

	return &Interpreter{
		env:            env,
		moduleCache:    make(map[string]*ModuleValue),
		loadingModules: make(map[string]bool),
		sourceDir:      cwd,
		currentDir:     cwd,
		nativeModules:  defaultNativeModules(),
	}
}

func (i *Interpreter) RegisterNativeModule(name string, loader NativeModuleLoader) {
	if i.nativeModules == nil {
		i.nativeModules = make(map[string]NativeModuleLoader)
	}
	i.nativeModules[name] = loader
}

func (i *Interpreter) SetSourcePath(filePath string) {
	abs, err := filepath.Abs(filePath)
	if err != nil {
		return
	}
	base := filepath.Dir(abs)
	i.sourceDir = base
	i.currentDir = base
}

func (i *Interpreter) SetVar(name string, value Value) {
	i.env.Set(name, value)
}

func (i *Interpreter) GetVar(name string) Value {
	return i.env.Get(name)
}

func (i *Interpreter) Run(m *ast.Module) {
	for _, stmt := range m.Body {
		i.execStmt(stmt)
	}
}

func (i *Interpreter) execStmt(stmt ast.Stmt) Signal {
	switch e := stmt.(type) {
	case *ast.ExprStmt:
		i.execExpr(e.Value)
		return NewNoSignal()
	case *ast.ImportStmt:
		return i.execImportStmt(e)
	case *ast.FromImportStmt:
		return i.execFromImportStmt(e)
	case *ast.IfStmt:
		return i.execIfStmt(e)
	case *ast.ForStmt:
		return i.execForStmt(e)
	case *ast.RepeatStmt:
		return i.execRepeatStmt(e)
	case *ast.WhileStmt:
		return i.execWhileStmt(e)
	case *ast.AssignStmt:
		value := i.execExpr(e.Value)
		i.env.Set(e.Target.Id, value)
		return NewNoSignal()
	case *ast.AugmentedAssignStmt:
		return i.execAugmentedAssignStmt(e)
	case *ast.FunctionDefStmt:
		return i.execFunctionDefStmt(e)
	case *ast.ClassDefStmt:
		return i.execClassDefStmt(e)
	case *ast.ReturnStmt:
		return i.execReturnStmt(e)
	case *ast.BreakStmt:
		return NewBreakSignal()
	case *ast.ContinueStmt:
		return NewContinueSignal()

	default:
		panic(newRuntimeError("unknown statement: %T", stmt))
	}
}

func (i *Interpreter) execExpr(expr ast.Expr) Value {
	switch e := expr.(type) {
	case *ast.Constant:
		return i.execConstExpr(e)

	case *ast.Name:
		return i.env.Get(e.Id)

	case *ast.UnaryOp:
		return i.execUnaryOpExpr(e)

	case *ast.BinOp:
		return i.execBinOpExpr(e)

	case *ast.BoolOp:
		return i.execBoolOpExpr(e)

	case *ast.Compare:
		return i.execCompareExpr(e)

	case *ast.Assign:
		return i.execAssignExpr(e)

	case *ast.Call:
		return i.execCallExpr(e)

	case *ast.Subscript:
		return i.execSubscriptExpr(e)

	case *ast.Attribute:
		return i.execAttributeExpr(e)

	case *ast.List:
		return i.execListExpr(e)

	case *ast.Dict:
		return i.execDictExpr(e)

	case *ast.Tuple:
		return i.execTupleExpr(e)
	}

	panic(newRuntimeError("unknown expression: %T", expr))
}
