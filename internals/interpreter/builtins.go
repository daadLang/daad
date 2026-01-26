package interpreter

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Arabic format specifiers mapping to Go format specifiers
var formatters = map[string]string{
	"ر": "d", // رقم (integer)
	"ع": "f", // عدد عشري (float)
	"ن": "s", // نص (string)
	"م": "t", // منطقي (boolean)
	"ث": "b", // ثنائي (binary)
	"و": "o", // ثماني (octal)
	"س": "x", // سداسي عشري (hex)
	"ل": "e", // كتابة علمية (scientific)
}

// RegisterBuiltins adds all builtin functions to the environment
func RegisterBuiltins(env *Env) {
	// i/o functions
	env.Set("اطبع", &BuiltinValue{Name: "اطبع", Fn: builtinPrint, Variadic: true})
	env.Set("ادخل", &BuiltinValue{Name: "ادخل", Fn: builtinInput, Variadic: true})

	// formatting functions
	env.Set("نسق", &BuiltinValue{Name: "نسق", Fn: builtinFormat, Variadic: true})

	// array options
	env.Set("طول", &BuiltinValue{Name: "طول", Fn: builtinLen, Variadic: false})

	// list manipulation builtins
	// English: append
	env.Set("اضف", &BuiltinValue{Name: "اضف", Fn: builtinAppend, Variadic: false})
	// English: push
	env.Set("ادفع", &BuiltinValue{Name: "ادفع", Fn: builtinPush, Variadic: false})
	// English: pop
	env.Set("ازل", &BuiltinValue{Name: "ازل", Fn: builtinPop, Variadic: false})
	// English: copy
	env.Set("انسخ", &BuiltinValue{Name: "انسخ", Fn: builtinCopy, Variadic: false})
	// English: clear
	env.Set("افرغ", &BuiltinValue{Name: "افرغ", Fn: builtinClear, Variadic: false})

	// type conversion functions
	env.Set("نوع", &BuiltinValue{Name: "نوع", Fn: builtinType, Variadic: false})

	env.Set("نطاق", &BuiltinValue{Name: "نطاق", Fn: builtinRange, Variadic: true})
	env.Set("صحيح", &BuiltinValue{Name: "صحيح", Fn: builtinInt, Variadic: false})
	env.Set("عشري", &BuiltinValue{Name: "عشري", Fn: builtinFloat, Variadic: false})
	env.Set("نص", &BuiltinValue{Name: "نص", Fn: builtinStr, Variadic: false})
}

// Print function: اطبع
// Supports: اطبع("مرحبا %ن، عمرك %ر سنة", "أحمد", 25)
// kwargs: فاصل (sep), نهاية (end)
func builtinPrint(args []Value, kwargs map[string]Value) (Value, error) {
	// Get separator (default is space)
	sep := " "
	if sepVal, ok := kwargs["فاصل"]; ok {
		if s, ok := sepVal.(StringValue); ok {
			sep = s.V
		}
	}

	// Get end character (default is newline)
	end := "\n"
	if endVal, ok := kwargs["نهاية"]; ok {
		if e, ok := endVal.(StringValue); ok {
			end = e.V
		}
	}

	if len(args) == 0 {
		fmt.Print(end)
		return NoneValue{}, nil
	}

	// Check if first arg is a format string
	if formatStr, ok := args[0].(StringValue); ok && len(args) > 1 {
		// Check if it contains format specifiers
		if containsFormatSpecifier(formatStr.V) {
			result, err := formatString(formatStr.V, args[1:])
			if err != nil {
				return nil, err
			}
			fmt.Print(result + end)
			return NoneValue{}, nil
		}
	}

	// No format string, just print all args with separator
	parts := make([]string, len(args))
	for i, arg := range args {
		parts[i] = valueToString(arg)
	}
	fmt.Print(strings.Join(parts, sep) + end)

	return NoneValue{}, nil
}

// format function: نسق
// Returns formatted string without printing
// نسق("مرحبا %ن، عمرك %ر سنة", "أحمد", 25)

func builtinFormat(args []Value, kwargs map[string]Value) (Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("نسق() تتطلب على الأقل معامل واحد")
	}

	formatStr, ok := args[0].(StringValue)
	if !ok {
		return nil, fmt.Errorf("المعامل الأول لـ نسق() يجب أن يكون نصاً")
	}

	result, err := formatString(formatStr.V, args[1:])
	if err != nil {
		return nil, err
	}

	return StringValue{V: result}, nil
}

// =============================================================================
// Helper functions for formatting
// =============================================================================

// containsFormatSpecifier checks if string contains Arabic format specifiers
func containsFormatSpecifier(s string) bool {
	// Match %ر, %ع, %ن, etc. with optional width/precision
	pattern := `%[-+0-9.]*[رعنمثوسل]`
	matched, _ := regexp.MatchString(pattern, s)
	return matched
}

// formatString processes Arabic format specifiers
func formatString(format string, args []Value) (string, error) {
	// regex to find arabic format specifiers with optional flags
	// Matches: %ر, %-10ن, %.2ع, %+5ر, etc.
	pattern := regexp.MustCompile(`%([-+]?)(\d*)(?:\.(\d+))?([رعنمثوسل])`)

	argIndex := 0
	result := pattern.ReplaceAllStringFunc(format, func(match string) string {
		if argIndex >= len(args) {
			return match // Not enough args, leave as-is
		}

		// parse the format specifier
		submatches := pattern.FindStringSubmatch(match)
		flags := submatches[1]
		width := submatches[2]
		precision := submatches[3]
		specifier := submatches[4]

		// get the corresponding Go format specifier
		goSpec, ok := formatters[specifier]
		if !ok {
			return match
		}

		// build Go format string
		var goFormat strings.Builder
		goFormat.WriteString("%")
		if flags != "" {
			goFormat.WriteString(flags)
		}
		if width != "" {
			goFormat.WriteString(width)
		}
		if precision != "" {
			goFormat.WriteString(".")
			goFormat.WriteString(precision)
		}
		goFormat.WriteString(goSpec)

		// format the args
		arg := args[argIndex]
		argIndex++

		return formatValue(goFormat.String(), arg, goSpec)
	})

	return result, nil
}

// formatValue formats a single value with a Go format specifier
func formatValue(format string, val Value, spec string) string {
	switch spec {
	case "d": // Integer
		if v, ok := val.(IntValue); ok {
			return fmt.Sprintf(format, v.V)
		}
		if v, ok := val.(FloatValue); ok {
			return fmt.Sprintf(format, int(v.V))
		}
		return fmt.Sprintf("%v", getRawValue(val))

	case "f", "e": // Float / Scientific
		if v, ok := val.(FloatValue); ok {
			return fmt.Sprintf(format, v.V)
		}
		if v, ok := val.(IntValue); ok {
			return fmt.Sprintf(format, float64(v.V))
		}
		return fmt.Sprintf("%v", getRawValue(val))

	case "s": // String
		return fmt.Sprintf(format, valueToString(val))

	case "t": // Boolean
		if v, ok := val.(BoolValue); ok {
			if v.V {
				return "صحيح"
			}
			return "خطأ"
		}
		return fmt.Sprintf("%v", getRawValue(val))

	case "b": // Binary
		if v, ok := val.(IntValue); ok {
			return fmt.Sprintf(format, v.V)
		}
		return fmt.Sprintf("%v", getRawValue(val))

	case "o": // Octal
		if v, ok := val.(IntValue); ok {
			return fmt.Sprintf(format, v.V)
		}
		return fmt.Sprintf("%v", getRawValue(val))

	case "x": // Hex
		if v, ok := val.(IntValue); ok {
			return fmt.Sprintf(format, v.V)
		}
		return fmt.Sprintf("%v", getRawValue(val))

	default:
		return fmt.Sprintf("%v", getRawValue(val))
	}
}

// valueToString converts any Value to its string representation
func valueToString(val Value) string {
	switch v := val.(type) {
	case IntValue:
		return strconv.Itoa(v.V)
	case FloatValue:
		return strconv.FormatFloat(v.V, 'f', -1, 64)
	case StringValue:
		return v.V
	case CharValue:
		return string(v.V)
	case BoolValue:
		if v.V {
			return "صحيح"
		}
		return "خطأ"
	case NoneValue:
		return "عدم"
	case ListValue:
		parts := make([]string, len(v.Elements))
		for i, elem := range v.Elements {
			parts[i] = valueToRepr(elem)
		}
		return "[" + strings.Join(parts, "، ") + "]"
	case TupleValue:
		parts := make([]string, len(v.Elements))
		for i, elem := range v.Elements {
			parts[i] = valueToRepr(elem)
		}
		return "(" + strings.Join(parts, "، ") + ")"
	case DictValue:
		parts := make([]string, 0, len(v.Entries))
		for k, val := range v.Entries {
			parts = append(parts, fmt.Sprintf("%v: %s", k, valueToRepr(val)))
		}
		return "{" + strings.Join(parts, "، ") + "}"
	case *FunctionValue:
		return fmt.Sprintf("<دالة %s>", v.Name)
	case *BuiltinValue:
		return fmt.Sprintf("<دالة مدمجة %s>", v.Name)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// valueToRepr returns repr-style string (with quotes for strings)
func valueToRepr(val Value) string {
	switch v := val.(type) {
	case StringValue:
		return fmt.Sprintf("\"%s\"", v.V)
	case CharValue:
		return fmt.Sprintf("'%c'", v.V)
	default:
		return valueToString(val)
	}
}

func getRawValue(val Value) interface{} {
	switch v := val.(type) {
	case IntValue:
		return v.V
	case FloatValue:
		return v.V
	case StringValue:
		return v.V
	case CharValue:
		return v.V
	case BoolValue:
		return v.V
	case NoneValue:
		return nil
	default:
		return val
	}
}

// ? =============================================================================
// ? Other builtin functions
// ? =============================================================================

func builtinLen(args []Value, kwargs map[string]Value) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("طول() تتطلب بالضبط معامل واحد (%d معطى)", len(args))
	}

	switch v := args[0].(type) {
	case StringValue:
		return IntValue{V: len([]rune(v.V))}, nil
	case ListValue:
		return IntValue{V: len(v.Elements)}, nil
	case TupleValue:
		return IntValue{V: len(v.Elements)}, nil
	case DictValue:
		return IntValue{V: len(v.Entries)}, nil
	default:
		return nil, fmt.Errorf("معامل طول() يجب أن يكون تسلسلاً (نص/قائمة/صف/قاموس)، لكن حصلنا على %T", args[0])
	}
}

// append(list, elem) -> returns a new list with elem appended
func builtinAppend(args []Value, kwargs map[string]Value) (Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("اضف() تتطلب بالضبط معاملين (%d معطى)", len(args))
	}

	lst, ok := args[0].(ListValue)
	if !ok {
		return nil, fmt.Errorf("المعامل الأول لـ اضف() يجب أن يكون قائمة، حصلنا على %T", args[0])
	}

	newElems := make([]Value, len(lst.Elements)+1)
	copy(newElems, lst.Elements)
	newElems[len(lst.Elements)] = args[1]

	return ListValue{Elements: newElems}, nil
}

// push alias for append
func builtinPush(args []Value, kwargs map[string]Value) (Value, error) {
	return builtinAppend(args, kwargs)
}

// ازل(list) -> ترجع القائمة نفسها بدون العنصر الأخير
func builtinPop(args []Value, kwargs map[string]Value) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("ازل() تتطلب بالضبط معامل واحد (%d معطى)", len(args))
	}

	lst, ok := args[0].(ListValue)
	if !ok {
		return nil, fmt.Errorf("المعامل لـ ازل() يجب أن يكون قائمة، حصلنا على %T", args[0])
	}

	if len(lst.Elements) == 0 {
		return nil, fmt.Errorf("محاولة ازالة عنصر من قائمة فارغة")
	}

	lastIdx := len(lst.Elements) - 1
	newElems := make([]Value, lastIdx)
	copy(newElems, lst.Elements[:lastIdx])

	return ListValue{Elements: newElems}, nil
}

// copy(list) -> shallow copy of list
func builtinCopy(args []Value, kwargs map[string]Value) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("انسخ() تتطلب بالضبط معامل واحد (%d معطى)", len(args))
	}

	lst, ok := args[0].(ListValue)
	if !ok {
		return nil, fmt.Errorf("المعامل لـ انسخ() يجب أن يكون قائمة، حصلنا على %T", args[0])
	}

	newElems := make([]Value, len(lst.Elements))
	copy(newElems, lst.Elements)
	return ListValue{Elements: newElems}, nil
}

// clear(list) -> returns an empty list
func builtinClear(args []Value, kwargs map[string]Value) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("افرغ() تتطلب بالضبط معامل واحد (%d معطى)", len(args))
	}

	_, ok := args[0].(ListValue)
	if !ok {
		return nil, fmt.Errorf("المعامل لـ افرغ() يجب أن يكون قائمة، حصلنا على %T", args[0])
	}

	return ListValue{Elements: []Value{}}, nil
}

// نوع - type function
func builtinType(args []Value, kwargs map[string]Value) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("نوع() تتطلب بالضبط معامل واحد (%d معطى)", len(args))
	}

	var typeName string
	switch args[0].(type) {
	case IntValue:
		typeName = "صحيح"
	case FloatValue:
		typeName = "عشري"
	case StringValue:
		typeName = "نص"
	case CharValue:
		typeName = "حرف"
	case BoolValue:
		typeName = "منطقي"
	case ListValue:
		typeName = "قائمة"
	case TupleValue:
		typeName = "صف"
	case DictValue:
		typeName = "قاموس"
	case *FunctionValue:
		typeName = "دالة"
	case *BuiltinValue:
		typeName = "دالة_مدمجة"
	case NoneValue:
		typeName = "عدم"
	default:
		typeName = "مجهول"
	}

	return StringValue{V: typeName}, nil
}

// نطاق - range function
func builtinRange(args []Value, kwargs map[string]Value) (Value, error) {
	var start, stop, step int

	switch len(args) {
	case 1:
		if v, ok := args[0].(IntValue); ok {
			start, stop, step = 0, v.V, 1
		} else {
			return nil, fmt.Errorf("المعامل لـ نطاق() يجب أن يكون عددًا صحيحًا")
		}

	case 2:
		v1, ok1 := args[0].(IntValue)
		v2, ok2 := args[1].(IntValue)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("معاملات نطاق() يجب أن تكون أعدادًا صحيحة")
		}
		start, stop, step = v1.V, v2.V, 1
	case 3:
		v1, ok1 := args[0].(IntValue)
		v2, ok2 := args[1].(IntValue)
		v3, ok3 := args[2].(IntValue)
		if !ok1 || !ok2 || !ok3 {
			return nil, fmt.Errorf("معاملات نطاق() يجب أن تكون أعدادًا صحيحة")
		}
		start, stop, step = v1.V, v2.V, v3.V
		if step == 0 {
			return nil, fmt.Errorf("قيمة الخطوة في نطاق() لا يمكن أن تكون صفرًا")
		}
	default:
		return nil, fmt.Errorf("نطاق() تتطلب من 1 إلى 3 معاملات (%d معطى)", len(args))
	}

	var elements []Value
	if step > 0 {
		for i := start; i < stop; i += step {
			elements = append(elements, IntValue{V: i})
		}
	} else {
		for i := start; i > stop; i += step {
			elements = append(elements, IntValue{V: i})
		}
	}

	return ListValue{Elements: elements}, nil
}

func builtinInt(args []Value, kwargs map[string]Value) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("صحيح() تتطلب بالضبط معامل واحد (%d معطى)", len(args))
	}

	switch v := args[0].(type) {
	case IntValue:
		return v, nil
	case FloatValue:
		return IntValue{V: int(v.V)}, nil
	case StringValue:
		i, err := strconv.Atoi(strings.TrimSpace(v.V))
		if err != nil {
			return nil, fmt.Errorf("صحيح() لا يمكن تحويل '%s' إلى عدد صحيح", v.V)
		}
		return IntValue{V: i}, nil
	case BoolValue:
		if v.V {
			return IntValue{V: 1}, nil
		}
		return IntValue{V: 0}, nil
	default:
		return nil, fmt.Errorf("صحيح() لا يمكن تحويل %T إلى عدد صحيح", args[0])
	}
}

func builtinFloat(args []Value, kwargs map[string]Value) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("عشري() تتطلب بالضبط معامل واحد (%d معطى)", len(args))
	}

	switch v := args[0].(type) {
	case FloatValue:
		return v, nil
	case IntValue:
		return FloatValue{V: float64(v.V)}, nil
	case StringValue:
		f, err := strconv.ParseFloat(strings.TrimSpace(v.V), 64)
		if err != nil {
			return nil, fmt.Errorf("عشري() لا يمكن تحويل '%s' إلى عدد عشري", v.V)
		}
		return FloatValue{V: f}, nil
	default:
		return nil, fmt.Errorf("عشري() لا يمكن تحويل %T إلى عدد عشري", args[0])
	}
}

func builtinStr(args []Value, kwargs map[string]Value) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("نص() تتطلب بالضبط معامل واحد (%d معطى)", len(args))
	}

	return StringValue{V: valueToString(args[0])}, nil
}

func builtinInput(args []Value, kwargs map[string]Value) (Value, error) {
	// Print prompt if provided
	if len(args) > 0 {
		fmt.Print(valueToString(args[0]))
	}

	var input string
	fmt.Scanln(&input)

	return StringValue{V: input}, nil
}
