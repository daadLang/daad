# OOP Internals in Daad (v1)

This document explains how Object-Oriented Programming (OOP) works in this project from parsing to runtime execution.

Scope is the current minimal OOP implementation (v1):
- class definition
- object instantiation
- instance and class attributes
- methods
- constructor via `__بناء__`

Not implemented in v1:
- inheritance
- access modifiers
- overriding rules

## 1. High-Level Model

Daad executes directly from AST (no bytecode).

OOP flow:
1. Lexer emits tokens like `CLASS`, `DOT`, `NAME`, `ASSIGN`.
2. Parser builds AST nodes such as `ClassDefStmt`, `Attribute`, `Call`, `Assign`.
3. Interpreter evaluates statements and expressions.
4. Runtime values are stored in `Env` and represented by `ClassValue` and `ObjectValue`.

## 2. AST Nodes Used by OOP

Defined in `internals/ast/nodes.go`:

- `ClassDefStmt`: class declaration
  - fields: `Name`, `Body`
- `Attribute`: attribute access expression
  - models `obj.attr`
  - fields: `Value` (object/class expression), `Attr` (attribute name)
- `Assign`: assignment expression
  - target can be `Name` or `Attribute`
- `Call`: function or method invocation
  - used for `ClassName(...)` and `obj.method(...)`

Important detail: method calls are represented as a call whose function expression is an attribute:
- `obj.method()` becomes `Call{Func: Attribute{Value: Name("obj"), Attr: "method"}}`

## 3. Parsing OOP Syntax

Implemented in `internals/parser/parser.go`.

### 3.1 Class definition

`parseClassDef` parses:

```daad
صنف شخص:
    ...
```

into:
- `ClassDefStmt{Name: "شخص", Body: classSuite}`

Class body supports (v1):
- method declarations (`FunctionDefStmt`)
- class-level assignments (as `ExprStmt` containing `Assign`)

### 3.2 Attribute access and calls

`parsePrimary` repeatedly extends atoms with postfix operations:
- `(... )` for calls
- `[...]` for subscripts
- `.name` for attributes

That means chained forms are naturally supported by AST construction order, such as:
- `a.b.c()`
- `obj.method(arg)`

## 4. Runtime OOP Types

Defined in `internals/interpreter/types.go`.

### 4.1 `ClassValue`

Represents a class at runtime:
- `Name string`
- `Attributes map[string]Value`

The attributes map stores both:
- class fields
- methods (stored as `*FunctionValue`)

### 4.2 `ObjectValue`

Represents an object instance:
- `Class *ClassValue`
- `Attributes map[string]Value`

The instance has a pointer to its class plus its own attribute map.

## 5. Environment and Scope

Environment is implemented in `internals/interpreter/env.go`:
- `Env.values` stores name -> value
- `Env.parent` forms lexical scope chain

OOP-related scope behavior:
- When class methods are created in `execClassDefStmt`, each method `FunctionValue` captures current interpreter env (`Env: i.env`).
- When any function/method is called, `callFunction` creates a child env from the function closure env:
  - `funcEnv := NewEnv(fn.Env)`
  - parameters are bound inside `funcEnv`

So methods use lexical scoping like normal functions.

## 6. Class Definition Execution

Implemented in `execClassDefStmt` in `internals/interpreter/execStmt.go`.

Execution steps:
1. Create empty `classAttributes` map.
2. Iterate class body statements.
3. For each `FunctionDefStmt`:
   - create `FunctionValue`
   - store in class attributes by method name
4. For class-level assignment expressions:
   - only `Name` assignment target is accepted
   - evaluate value and store in class attributes
5. Register class in current env:
   - `env[className] = &ClassValue{...}`

## 7. Instantiation and Constructor

Implemented in `execCallExpr` in `internals/interpreter/execExpr.go` when callee is `*ClassValue`.

When `ClassName(args...)` is executed:
1. Create object:
   - `obj := &ObjectValue{Class: class, Attributes: map[string]Value{}}`
2. Copy non-function class attributes into object attributes.
3. If class has `__بناء__` and it is a function:
   - call it automatically
   - pass object as first argument (`ذاتي`)
   - append user-provided constructor args after it
4. Return object.

Constructor convention in v1:
- Name must be exactly `__بناء__`.

## 8. Attribute Resolution and Method Binding

Implemented in `execAttributeExpr` in `internals/interpreter/execExpr.go`.

### 8.1 Access on object instance

Lookup order for `obj.attr`:
1. `obj.Attributes[attr]` (instance first)
2. `obj.Class.Attributes[attr]` (class fallback)

If class attribute is a function (`*FunctionValue`), interpreter returns a bound callable:
- wraps method in `BuiltinValue`
- prepends instance (`obj`) to call args
- then calls `callFunction`

This is how `obj.method(x)` automatically receives `ذاتي` as first parameter.

### 8.2 Access on class

For `ClassName.attr`, lookup is direct in class attributes map.

## 9. Attribute Assignment

Implemented in `execAssignExpr` in `internals/interpreter/execExpr.go`.

Supported targets:
- `Name` -> normal variable assignment in env
- `Attribute` -> attribute assignment

For attribute assignment:
- if container evaluates to `*ObjectValue`, write to instance attributes
- if container evaluates to `*ClassValue`, write to class attributes
- otherwise runtime error

So both are valid:

```daad
ذاتي.الاسم = اسم
شخص.نوع = "إنسان"
```

## 10. Example End-to-End (from tests/examples/oop_class.daad)

Source:

```daad
صنف شخص:
    ن = "يسيسي"
    دالة __بناء__(ذاتي, اسم):
        ذاتي.الاسم = اسم

    دالة عرف_نفسك(ذاتي):
        اطبع("أنا " + ذاتي.الاسم)

م = شخص("ssss")
م.عرف_نفسك()
اطبع(م.ن)
```

Runtime behavior:
1. Class `شخص` is created with:
   - class field `ن`
   - methods `__بناء__` and `عرف_نفسك`
2. `شخص("ssss")` creates object `م`, copies field `ن`, then calls constructor.
3. Constructor writes `الاسم` into instance attributes.
4. `م.عرف_نفسك()` resolves method from class, binds `م` as first arg, executes method.
5. `م.ن` succeeds via instance/class attribute resolution.

## 11. Current Limits and Notes

- No inheritance hierarchy.
- No explicit visibility (`private/public`).
- No metaclass behavior.
- Constructor invocation is convention-based (`__بناء__`).
- Instance receives a copy of non-function class attributes during instantiation.

## 12. Related Source Files

- AST nodes: `internals/ast/nodes.go`
- Parser: `internals/parser/parser.go`
- Statement execution: `internals/interpreter/execStmt.go`
- Expression execution: `internals/interpreter/execExpr.go`
- Runtime types: `internals/interpreter/types.go`
- Environment: `internals/interpreter/env.go`
- OOP AST test: `tests/ast/oop_ast_test.go`
- OOP interpreter test: `tests/interpreter/interpreter_test.go`
- OOP example: `tests/examples/oop_class.daad`
