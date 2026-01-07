
# Jeem Language Notes

## Overview

- Jeem is an **interpreted Arabic programming language**.
- It will be **Python-style**, dynamically typed, and beginner-friendly.
- Version 1 (v1) will be **simple and minimal**, focusing on core concepts.

---

## v1 Design

### Implementation

- The interpreter will consist of:
  - **Tokenizer (lexer)**
  - **Parser → AST**
  - **AST executor** (execute directly from AST, no bytecode in the first version)
- The **AST should be compatible with Python's AST** structure to make it easier to integrate or reuse in Python code in future versions (v2).
- **For simplicity, the first version will be implemented in Go**.
  - Later, performance-critical parts can be rewritten in C/C++ or even assembly for optimization.
- Goal: **Keep it simple and understandable**.

### Input simplification

- Typing should be easy:
  - Letters like `أ`, `ا`, `ء`, `آ`, `إ` can be considered equivalent.

### Features in v1

- Variables
- Numbers and strings
- Expressions
- Conditions (`if / elif / else`)
- Loops (`while`, `for`)
- Basic error handling

### File extension

- `.jeem` , `.jm` or `.ج`
 
---

## Example Code

### Function & While Loop

```jeem
دالة تحقق(عدد):
    إذا عدد % 2 == 0:
        أرجع "زوجي"
    وإلا:
        أرجع "فردي"

س = 1
طالما س <= 5:
    اطبع(تحقق(س))
    س = س + 1
````

### Conditional Example

```jeem
إذا درجة >= 90:
    اطبع("ممتاز")
وإذا درجة >= 75:
    اطبع("جيد جداً")
وإذا درجة >= 60:
    اطبع("جيد")
وإلا:
    اطبع("راسب")
```

### For Loop Example

```jeem
أسماء = ["علي", "أحمد"]

لكل اسم في أسماء:
    اطبع(اسم)
```

### Repeat Loop Example (syntactic sugar)

```jeem
كرر 5 مرات:
    اطبع("مرحبا")
```

Equivalent to:

```jeem
لكل _ في مدى(5):
    اطبع("مرحبا")
```

---

## Jeem Keyword Table

### Control Flow

| Concept  | Jeem (Arabic) |
| -------- | ------------- |
| if       | إذا           |
| elif     | وإذا          |
| else     | وإلا          |
| while    | طالما         |
| for      | لكل           |
| in       | في            |
| break    | اخرج          |
| continue | تابع          |
| return   | أرجع          |

### Functions

| Concept  | Jeem |
| -------- | ---- |
| function | دالة |

### Boolean Values

| Concept | Jeem |
| ------- | ---- |
| true    | صحيح |
| false   | خطأ  |

### Built-in Operations (Initial)

| Concept | Jeem |
| ------- | ---- |
| print   | اطبع |
| range   | مدى  |

### Optional / Future Keywords

| Concept | Jeem   | Notes              |
| ------- | ------ | ------------------ |
| and     | و      | Logical AND        |
| or      | أو     | Logical OR         |
| not     | ليس    | Logical NOT        |
| null    | لاشيء  | None / null        |
| import  | استورد | For later versions |

### Exception Handling

| English | Jeem (Arabic) | Notes                                     |
| ------- | ------------- | ----------------------------------------- |
| try     | حاول          | Simple and intuitive                      |
| except  | باستثناء      | Can optionally capture exception variable |
| finally | أخيراً        | Optional                                  |

---

## Example Using Many Keywords

```jeem
دالة تحقق(عدد):
    إذا عدد > 0:
        أرجع صحيح
    وإذا عدد < 0:
        أرجع خطأ
    وإلا:
        أرجع لاشيء

لكل س في مدى(1, 5):
    إذا تحقق(س):
        اطبع(س)
    وإلا:
        تابع
```

---

## Additional Notes

* **VSCode Extensions / Setup:**

  * RTL support is important
  * Syntax highlighting

* **Design Philosophy:**

  * Keep v1 simple
  * AST-based execution
  * Arabic-friendly and readable

