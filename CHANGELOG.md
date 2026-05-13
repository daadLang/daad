# 0.2.0 (2026-05-13)


### Bug Fixes

* actions ([485eb0d](https://github.com/daadLang/daad/commit/485eb0d04e841e5d8a40e5b39c30c8ea3c59030e))
* actions ([0433412](https://github.com/daadLang/daad/commit/0433412aaf555bf35dff859f87c50a6b9d277be2))
* revert Go version to 1.25 ([e03892a](https://github.com/daadLang/daad/commit/e03892a11c339cbfb938a3a158297bc086189618))
* standardize AST node labels to lowercase for consistency ([2e81238](https://github.com/daadLang/daad/commit/2e8123809d4558ca0afe914a4d51a4dfd9872601))
* update build configuration for correct binary and main path ([3c104fa](https://github.com/daadLang/daad/commit/3c104fab2489856674e2534fe3bb28e098ea47dc))


### Features

* **interpreter:** add some native modules support (os, random,math,time,path) ([355a15e](https://github.com/daadLang/daad/commit/355a15e7547a36c53067d9eed2f56eeeb7149b41))
* **interpreter:** empliment simple import & FromImport ([3265589](https://github.com/daadLang/daad/commit/3265589708bd1d3be487b981b1d55b7a8b93ff74))
* **parser:** add import & importFrom ([f238356](https://github.com/daadLang/daad/commit/f238356e00d90a0eeaf68fd10d86445bc0c68093))

---

# 0.1.0 (2026-05-06)

### Features

* **interpreter:** add basic object-oriented programming support
* **interpreter:** add class definitions, attribute assignments, method calls, and single inheritance
* **docs:** add OOP internals documentation and related examples
* **examples:** add class instantiation and inheritance examples

### Internal

* expand interpreter and AST to support class-based execution

---

# 0.0.0 (2026-02-01)

### Features

* **runtime:** initialize Daad language runtime foundation
* **lexer:** add indentation-aware tokenization
* **parser:** implement recursive-descent parser and AST structure
* **interpreter:** add variables, expressions, functions, collections, and control flow support
* **stdlib:** add initial built-in operations and standard behavior
* **docs:** add project documentation, examples, and test suite
