# Daad Language Grammar Specification

## Overview

This document defines the formal grammar for the **Daad (ض)** programming language.  
Daad is a Python-style, dynamically typed language with Arabic keywords.

**Excluded features:**
- Object-Oriented Programming (no classes, inheritance, methods)
- Type annotations and type checking
- Lambda expressions

---

## Notation

- `|` denotes alternatives
- `*` means zero or more
- `+` means one or more
- `?` means optional (zero or one)
- `()` for grouping
- `UPPERCASE` for tokens/terminals
- `lowercase` for non-terminals (grammar rules)

---

## 1. Program Structure

```ebnf
module          ::= statement*

statement       ::= simple_stmt NEWLINE
                  | compound_stmt

simple_stmt     ::= expr_stmt
                  | assignment_stmt
                  | augmented_assign_stmt
                  | return_stmt
                  | break_stmt
                  | continue_stmt

compound_stmt   ::= if_stmt
                  | while_stmt
                  | for_stmt
                  | repeat_stmt
                  | function_def
```

---

## 2. Simple Statements

### 2.1 Expression Statement
```ebnf
expr_stmt       ::= expression
```

### 2.2 Assignment Statement
```ebnf
assignment_stmt ::= target "=" expression
                  | target "=" assignment_stmt    # chained assignment

target          ::= NAME
                  | subscript

subscript       ::= primary "[" expression "]"
```

### 2.3 Augmented Assignment
```ebnf
augmented_assign_stmt ::= target augmented_op expression

augmented_op    ::= "+=" | "-=" | "*=" | "/=" | "%=" | "**="
```

### 2.4 Return Statement
```ebnf
return_stmt     ::= "ارجع" expression?
                  | "أرجع" expression?
```

### 2.5 Break Statement
```ebnf
break_stmt      ::= "اخرج"
```

### 2.6 Continue Statement
```ebnf
continue_stmt   ::= "تابع"
```

---

## 3. Compound Statements

### 3.1 If Statement
```ebnf
if_stmt         ::= "إذا" expression ":" suite
                    ("وإذا" expression ":" suite)*
                    ("وإلا" ":" suite)?

# Alternative spellings
                  | "اذا" expression ":" suite
                    ("واذا" expression ":" suite)*
                    ("والا" ":" suite)?
```

### 3.2 While Statement
```ebnf
while_stmt      ::= "طالما" expression ":" suite
```

### 3.3 For Statement
```ebnf
for_stmt        ::= "لكل" NAME "في" expression ":" suite
```

### 3.4 Repeat Statement (Syntactic Sugar)
```ebnf
repeat_stmt     ::= "كرر" expression ("مرات" | "مرة") ":" suite

# Equivalent to: لكل _ في مدى(expression): suite
```

### 3.5 Function Definition
```ebnf
function_def    ::= "دالة" NAME "(" parameters? ")" ":" suite

parameters      ::= param ("," param)*

param           ::= NAME
                  | NAME "=" expression    # default value
```

### 3.6 Suite (Indented Block)
```ebnf
suite           ::= NEWLINE INDENT statement+ DEDENT
                  | simple_stmt NEWLINE
```

---

## 4. Expressions

### 4.1 Expression Hierarchy (Precedence: low to high)
```ebnf
expression      ::= or_expr

or_expr         ::= and_expr ("أو" and_expr)*

and_expr        ::= not_expr ("و" not_expr)*

not_expr        ::= "ليس" not_expr
                  | "لا" not_expr
                  | comparison

comparison      ::= bitor_expr (comp_op bitor_expr)*

comp_op         ::= "==" | "!=" | "<" | ">" | "<=" | ">="
                  | "في"                    # membership test (in)
                  | "ليس" "في"              # not in
```

### 4.2 Bitwise Expressions
```ebnf
bitor_expr      ::= bitxor_expr ("|" bitxor_expr)*

bitxor_expr     ::= bitand_expr ("^" bitand_expr)*

bitand_expr     ::= shift_expr ("&" shift_expr)*

shift_expr      ::= arith_expr (("<<" | ">>") arith_expr)*
```

### 4.3 Arithmetic Expressions
```ebnf
arith_expr      ::= term (("+"|"-") term)*

term            ::= factor (("*"|"/"|"//"|"%") factor)*

factor          ::= ("+"|"-"|"~") factor
                  | power

power           ::= primary ("**" factor)?
```

### 4.4 Primary Expressions
```ebnf
primary         ::= atom
                  | primary "(" arguments? ")"     # function call
                  | primary "[" expression "]"      # subscript
                  | primary "." NAME                # attribute (for dicts)

atom            ::= NAME
                  | NUMBER
                  | STRING
                  | "صحيح"                          # True
                  | "خطأ"                           # False
                  | "عدم"                           # None
                  | "(" expression ")"              # parenthesized
                  | list_expr
                  | dict_expr
                  | tuple_expr
```

### 4.5 Function Call Arguments
```ebnf
arguments       ::= expression ("," expression)*
```

---

## 5. Collection Literals

### 5.1 List
```ebnf
list_expr       ::= "[" list_items? "]"

list_items      ::= expression ("," expression)* ","?
```

### 5.2 Dictionary
```ebnf
dict_expr       ::= "{" dict_items? "}"

dict_items      ::= key_value ("," key_value)* ","?

key_value       ::= expression ":" expression
```

### 5.3 Tuple
```ebnf
tuple_expr      ::= "(" ")"
                  | "(" expression "," ")"
                  | "(" expression ("," expression)+ ","? ")"
```

---

## 6. Lexical Elements

### 6.1 Keywords (Arabic)
```
إذا / اذا      - if
وإذا / واذا    - elif  
وإلا / والا    - else
طالما          - while
لكل            - for
في             - in
دالة           - function (def)
ارجع / أرجع    - return
اخرج           - break
تابع           - continue
كرر            - repeat
مرات           - times
و              - and
أو / او        - or
ليس / لا       - not
صحيح           - True
خطأ            - False
عدم            - None
```

### 6.2 Operators
```
# Arithmetic
+       addition
-       subtraction (and unary minus)
*       multiplication
/       division
//      floor division
%       modulo
**      power

# Assignment
=       assignment
+=      add and assign
-=      subtract and assign
*=      multiply and assign
/=      divide and assign
%=      modulo and assign
**=     power and assign

# Comparison
==      equal
!=      not equal
<       less than
>       greater than
<=      less or equal
>=      greater or equal

# Bitwise
&       bitwise AND
|       bitwise OR
^       bitwise XOR
~       bitwise NOT
<<      left shift
>>      right shift
```

### 6.3 Delimiters
```
(       left parenthesis
)       right parenthesis
[       left bracket
]       right bracket
{       left brace
}       right brace
,       comma
:       colon
.       dot
```

### 6.4 Literals

#### Number
```ebnf
NUMBER          ::= INTEGER | FLOAT

INTEGER         ::= DIGIT+
                  | "0x" HEXDIGIT+          # hexadecimal
                  | "0b" BINDIGIT+          # binary
                  | "0o" OCTDIGIT+          # octal

FLOAT           ::= DIGIT+ "." DIGIT*
                  | DIGIT* "." DIGIT+
                  | (DIGIT+ | FLOAT) ("e"|"E") ("+"|"-")? DIGIT+

DIGIT           ::= [0-9]
HEXDIGIT        ::= [0-9a-fA-F]
BINDIGIT        ::= [0-1]
OCTDIGIT        ::= [0-7]
```

#### String
```ebnf
STRING          ::= SHORT_STRING | LONG_STRING

SHORT_STRING    ::= '"' STRING_CHAR* '"'
                  | "'" STRING_CHAR* "'"

LONG_STRING     ::= '"""' LONG_CHAR* '"""'
                  | "'''" LONG_CHAR* "'''"

STRING_CHAR     ::= any character except newline, quote, or backslash
                  | ESCAPE_SEQUENCE

LONG_CHAR       ::= any character except backslash
                  | ESCAPE_SEQUENCE

ESCAPE_SEQUENCE ::= "\\" | "\'" | '\"' | "\n" | "\t" | "\r" | "\0"
```

#### Identifier (Name)
```ebnf
NAME            ::= (LETTER | "_") (LETTER | DIGIT | "_")*

LETTER          ::= [a-zA-Z]
                  | ARABIC_LETTER

ARABIC_LETTER   ::= [\u0600-\u06FF]         # Arabic Unicode block
```

### 6.5 Comments
```ebnf
COMMENT         ::= "#" (any character except newline)* NEWLINE
```

### 6.6 Indentation
```
INDENT          - increase in indentation level
DEDENT          - decrease in indentation level
NEWLINE         - end of logical line
```

---

## 7. Built-in Functions

```
اطبع(...)       - print (output to console)
مدى(n)          - range (generates sequence 0 to n-1)
مدى(start, end) - range (generates sequence start to end-1)
طول(x)          - len (length of collection)
نوع(x)          - type (returns type name as string)
عدد(x)          - int (convert to integer)
عشري(x)         - float (convert to float)
نص(x)           - str (convert to string)
قائمة(x)        - list (convert to list)
ادخل(prompt)    - input (read from console)
```

---

## 8. Grammar Examples

### Variable Assignment
<pre dir="rtl"><code>س = 10
ص = س + 5
</code></pre>

### Conditional
<pre dir="rtl"><code>إذا س > 10:
    اطبع("كبير")
وإذا س > 5:
    اطبع("متوسط")
وإلا:
    اطبع("صغير")
</code></pre>

### While Loop
<pre dir="rtl"><code>عداد = 0
طالما عداد < 10:
    اطبع(عداد)
    عداد += 1
</code></pre>

### For Loop
<pre dir="rtl"><code>أرقام = [1, 2, 3, 4, 5]
لكل رقم في أرقام:
    اطبع(رقم)
</code></pre>

### Repeat Loop
<pre dir="rtl"><code>كرر 5 مرات:
    اطبع("مرحبا")
</code></pre>

### Function Definition
<pre dir="rtl"><code>دالة مجموع(أ, ب):
    ارجع أ + ب

نتيجة = مجموع(10, 20)
اطبع(نتيجة)
</code></pre>

### Function with Default Parameter
<pre dir="rtl"><code>دالة تحية(اسم = "ضيف"):
    اطبع("مرحبا " + اسم)

تحية()           # prints: مرحبا ضيف
تحية("أحمد")     # prints: مرحبا أحمد
</code></pre>

### Collections
<pre dir="rtl"><code># List
قائمة = [1, 2, 3]
قائمة[0] = 10

# Dictionary
قاموس = {"اسم": "أحمد", "عمر": 25}
اطبع(قاموس["اسم"])

# Tuple (immutable)
ثنائي = (1, 2)
</code></pre>

### Logical Expressions
<pre dir="rtl"><code>إذا س > 0 و س < 100:
    اطبع("في المدى")

إذا س == 0 أو ص == 0:
    اطبع("أحدهما صفر")

إذا ليس متصل:
    اطبع("غير متصل")
</code></pre>

### Bitwise Operations
<pre dir="rtl"><code>قيمة = 5 & 3      # bitwise AND
قيمة = 5 | 3      # bitwise OR
قيمة = 5 ^ 3      # bitwise XOR
قيمة = ~5         # bitwise NOT
قيمة = 2 << 3     # left shift
قيمة = 16 >> 2    # right shift
</code></pre>
---

## 9. Operator Precedence (Highest to Lowest)

| Precedence | Operators                          | Description           |
|------------|------------------------------------|-----------------------|
| 1          | `()`                               | Parentheses           |
| 2          | `f()`, `x[i]`, `x.attr`            | Call, subscript, attr |
| 3          | `**`                               | Exponentiation        |
| 4          | `+x`, `-x`, `~x`                   | Unary operators       |
| 5          | `*`, `/`, `//`, `%`                | Multiplicative        |
| 6          | `+`, `-`                           | Additive              |
| 7          | `<<`, `>>`                         | Shift                 |
| 8          | `&`                                | Bitwise AND           |
| 9          | `^`                                | Bitwise XOR           |
| 10         | `\|`                               | Bitwise OR            |
| 11         | `==`, `!=`, `<`, `>`, `<=`, `>=`   | Comparison            |
| 12         | `ليس` / `لا`                       | Logical NOT           |
| 13         | `و`                                | Logical AND           |
| 14         | `أو`                               | Logical OR            |
| 15         | `=`, `+=`, `-=`, etc.              | Assignment            |

---

## 10. Reserved for Future Versions

The following features are intentionally **NOT** included in this grammar and may be added in future versions:

- **Classes and OOP**: `صنف` (class), inheritance, `ذاتي` (self)
- **Type annotations**: `-> نوع`, parameter types
- **Lambda expressions**: anonymous functions
- **Decorators**: `@decorator` syntax
- **Generators**: `أعط` (yield)
- **Context managers**: `مع` (with)
- **Exception handling**: `حاول` / `امسك` / `أخيرا` (try/except/finally)
- **Import system**: `استورد` (import)
- **Async/Await**: asynchronous programming

---

## Appendix: Token Reference Table

| Token Type     | Arabic              | Symbol/Keyword      |
|----------------|---------------------|---------------------|
| IF             | إذا / اذا           | if                  |
| ELIF           | وإذا / واذا         | elif                |
| ELSE           | وإلا / والا         | else                |
| WHILE          | طالما               | while               |
| FOR            | لكل                 | for                 |
| IN             | في                  | in                  |
| FUNC           | دالة                | def/function        |
| RETURN         | ارجع / أرجع         | return              |
| BREAK          | اخرج                | break               |
| CONTINUE       | تابع                | continue            |
| REPEAT         | كرر                 | repeat              |
| TIMES          | مرات                | times               |
| AND            | و                   | and                 |
| OR             | أو / او             | or                  |
| NOT            | ليس / لا            | not                 |
| TRUE           | صحيح                | True                |
| FALSE          | خطأ                 | False               |
| NONE           | عدم                 | None                |
