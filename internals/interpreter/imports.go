package interpreter

import (
	"os"
	"path/filepath"
	"strings"

	ast "github.com/daadLang/daad/internals/ast"
	"github.com/daadLang/daad/internals/lexer"
	"github.com/daadLang/daad/internals/parser"
)

func (i *Interpreter) execImportStmt(stmt *ast.ImportStmt) Signal {
	for _, alias := range stmt.Names {
		module := i.loadModule(alias.Name, 0)

		targetName := importBindingName(alias)
		i.env.Set(targetName, module)
	}

	return NewNoSignal()
}

func (i *Interpreter) execFromImportStmt(stmt *ast.FromImportStmt) Signal {
	if stmt.Module == "" {
		// Python-like: from . import modA, modB
		for _, imported := range stmt.Names {
			if imported.Name == "*" {
				panic(newRuntimeError("star import requires a module name in from-import statement"))
			}

			module := i.loadModule(imported.Name, stmt.Level)
			targetName := imported.Name
			if imported.AsName != nil {
				targetName = *imported.AsName
			}
			i.env.Set(targetName, module)
		}
		return NewNoSignal()
	}

	module := i.loadModule(stmt.Module, stmt.Level)
	for _, imported := range stmt.Names {
		if imported.Name == "*" {
			for name, value := range module.Attributes {
				i.env.Set(name, value)
			}
			continue
		}

		value, ok := module.Attributes[imported.Name]
		if !ok {
			panic(newRuntimeError("module '%s' has no exported name '%s'", stmt.Module, imported.Name))
		}

		targetName := imported.Name
		if imported.AsName != nil {
			targetName = *imported.AsName
		}
		i.env.Set(targetName, value)
	}

	return NewNoSignal()
}

func importBindingName(alias ast.Alias) string {
	if alias.AsName != nil {
		return *alias.AsName
	}

	name := alias.Name
	if dot := strings.IndexRune(name, '.'); dot != -1 {
		return name[:dot]
	}
	return name
}

func (i *Interpreter) loadModule(moduleName string, level int) *ModuleValue {
	modulePath := i.resolveModulePath(moduleName, level)

	if cached, ok := i.moduleCache[modulePath]; ok {
		return cached
	}

	if i.loadingModules[modulePath] {
		panic(newRuntimeError("circular import detected for module '%s'", moduleName))
	}
	i.loadingModules[modulePath] = true
	defer delete(i.loadingModules, modulePath)

	tokens, err := lexer.Tokenize(modulePath)
	if err != nil {
		panic(newRuntimeError("failed to tokenize module '%s': %v", moduleName, err))
	}

	p := parser.NewParser(tokens)
	moduleAst := p.Parse()

	moduleEnv := NewEnv(nil)
	RegisterBuiltins(moduleEnv)

	parentEnv := i.env
	parentDir := i.currentDir
	i.env = moduleEnv
	i.currentDir = filepath.Dir(modulePath)

	sig := i.execBlock(moduleAst.Body)
	if sig.IsError() {
		i.env = parentEnv
		i.currentDir = parentDir
		panic(sig.Err)
	}
	if sig.IsReturn() || sig.IsBreak() || sig.IsContinue() {
		i.env = parentEnv
		i.currentDir = parentDir
		panic(newRuntimeError("invalid control flow signal in module '%s'", moduleName))
	}

	i.env = parentEnv
	i.currentDir = parentDir

	exports := make(map[string]Value)
	for name, value := range moduleEnv.values {
		exports[name] = value
	}

	loaded := &ModuleValue{
		Name:       moduleName,
		Path:       modulePath,
		Attributes: exports,
	}
	i.moduleCache[modulePath] = loaded
	return loaded
}

func (i *Interpreter) resolveModulePath(moduleName string, level int) string {
	relativePath := strings.ReplaceAll(moduleName, ".", string(filepath.Separator))

	bases := make([]string, 0, 3)
	if level > 0 {
		base := i.currentDir
		if base == "" {
			base = i.sourceDir
		}
		if base == "" {
			base = "."
		}

		// level=1 => current package dir, level=2 => parent, ...
		for step := 1; step < level; step++ {
			base = filepath.Dir(base)
		}
		bases = append(bases, base)
	} else {
		if i.currentDir != "" {
			bases = append(bases, i.currentDir)
		}
		if i.sourceDir != "" && i.sourceDir != i.currentDir {
			bases = append(bases, i.sourceDir)
		}
		bases = append(bases, ".")
	}

	seen := make(map[string]bool)
	for _, base := range bases {
		if base == "" {
			continue
		}
		if seen[base] {
			continue
		}
		seen[base] = true

		candidates := make([]string, 0, 4)
		if relativePath != "" {
			// ? _.daad or _.ض is equivalent to __init__.py in python
			candidates = append(candidates,
				filepath.Join(base, relativePath+".daad"),
				filepath.Join(base, relativePath+".ض"),
				filepath.Join(base, relativePath, "_.daad"),
				filepath.Join(base, relativePath, "_.ض"),
			)
		} else {
			candidates = append(candidates,
				filepath.Join(base, "_.daad"),
				filepath.Join(base, "_.ض"),
			)
		}

		for _, candidate := range candidates {
			if _, err := filepath.Abs(candidate); err != nil {
				continue
			}
			if fileExists(candidate) {
				abs, err := filepath.Abs(candidate)
				if err != nil {
					return candidate
				}
				return abs
			}
		}
	}

	if level > 0 {
		panic(newRuntimeError("cannot resolve relative module '%s' (level=%d)", moduleName, level))
	}
	panic(newRuntimeError("module '%s' not found", moduleName))
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}
