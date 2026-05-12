package interpreter

// TODO: refactor this file into multiple files (math.go, random.go, time.go, os.go, path.go)
import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

type NativeModuleLoader func() *ModuleValue

func defaultNativeModules() map[string]NativeModuleLoader {
	return map[string]NativeModuleLoader{
		"رياضيات": newMathModule,
		"عشوائي":  newRandomModule,
		"وقت":     newTimeModule,
		"نظام":    newOSModule,
		"مسار":    newPathModule,
	}
}

func newNativeModule(name string, exports map[string]Value) *ModuleValue {
	return &ModuleValue{
		Name:       name,
		Path:       "native:" + name,
		Attributes: exports,
	}
}

func builtin(name string, fn BuiltinFunc, variadic bool) *BuiltinValue {
	return &BuiltinValue{Name: name, Fn: fn, Variadic: variadic}
}

func expectNumber(arg Value, name string) (float64, bool, error) {
	switch v := arg.(type) {
	case IntValue:
		return float64(v.V), true, nil
	case FloatValue:
		return v.V, false, nil
	default:
		return 0, false, fmt.Errorf("%s يتوقع رقما، حصلنا على %T", name, arg)
	}
}

func expectInt(arg Value, name string) (int, error) {
	switch v := arg.(type) {
	case IntValue:
		return v.V, nil
	case FloatValue:
		return int(v.V), nil
	default:
		return 0, fmt.Errorf("%s يتوقع عددا صحيحا، حصلنا على %T", name, arg)
	}
}

func expectString(arg Value, name string) (string, error) {
	if s, ok := arg.(StringValue); ok {
		return s.V, nil
	}
	return "", fmt.Errorf("%s يتوقع نصا، حصلنا على %T", name, arg)
}

func newMathModule() *ModuleValue {
	exports := map[string]Value{
		"باي":      FloatValue{V: math.Pi},
		"هـ":       FloatValue{V: math.E},
		"مطلق":     builtin("مطلق", mathAbs, false),
		"جذر":      builtin("جذر", mathSqrt, false),
		"جيب":      builtin("جيب", mathSin, false),
		"جيب_تمام": builtin("جيب_تمام", mathCos, false),
		"ظل":       builtin("ظل", mathTan, false),
		"أس":       builtin("أس", mathPow, false),
		"أرض":      builtin("أرض", mathFloor, false),
		"سقف":      builtin("سقف", mathCeil, false),
		"تقريب":    builtin("تقريب", mathRound, false),
	}
	return newNativeModule("رياضيات", exports)
}

func mathAbs(args []Value, _ map[string]Value) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("مطلق يتوقع وسيطا واحدا")
	}
	val, isInt, err := expectNumber(args[0], "مطلق")
	if err != nil {
		return nil, err
	}
	if isInt {
		if val < 0 {
			return IntValue{V: -int(val)}, nil
		}
		return IntValue{V: int(val)}, nil
	}
	return FloatValue{V: math.Abs(val)}, nil
}

func mathSqrt(args []Value, _ map[string]Value) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("جذر يتوقع وسيطا واحدا")
	}
	val, _, err := expectNumber(args[0], "جذر")
	if err != nil {
		return nil, err
	}
	return FloatValue{V: math.Sqrt(val)}, nil
}

func mathSin(args []Value, _ map[string]Value) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("جيب يتوقع وسيطا واحدا")
	}
	val, _, err := expectNumber(args[0], "جيب")
	if err != nil {
		return nil, err
	}
	return FloatValue{V: math.Sin(val)}, nil
}

func mathCos(args []Value, _ map[string]Value) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("جيب_تمام يتوقع وسيطا واحدا")
	}
	val, _, err := expectNumber(args[0], "جيب_تمام")
	if err != nil {
		return nil, err
	}
	return FloatValue{V: math.Cos(val)}, nil
}

func mathTan(args []Value, _ map[string]Value) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("ظل يتوقع وسيطا واحدا")
	}
	val, _, err := expectNumber(args[0], "ظل")
	if err != nil {
		return nil, err
	}
	return FloatValue{V: math.Tan(val)}, nil
}

func mathPow(args []Value, _ map[string]Value) (Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("أس يتوقع وسيطين")
	}
	left, _, err := expectNumber(args[0], "أس")
	if err != nil {
		return nil, err
	}
	right, _, err := expectNumber(args[1], "أس")
	if err != nil {
		return nil, err
	}
	return FloatValue{V: math.Pow(left, right)}, nil
}

func mathFloor(args []Value, _ map[string]Value) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("أرض يتوقع وسيطا واحدا")
	}
	val, _, err := expectNumber(args[0], "أرض")
	if err != nil {
		return nil, err
	}
	return IntValue{V: int(math.Floor(val))}, nil
}

func mathCeil(args []Value, _ map[string]Value) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("سقف يتوقع وسيطا واحدا")
	}
	val, _, err := expectNumber(args[0], "سقف")
	if err != nil {
		return nil, err
	}
	return IntValue{V: int(math.Ceil(val))}, nil
}

func mathRound(args []Value, _ map[string]Value) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("تقريب يتوقع وسيطا واحدا")
	}
	val, _, err := expectNumber(args[0], "تقريب")
	if err != nil {
		return nil, err
	}
	return IntValue{V: int(math.Round(val))}, nil
}

func newRandomModule() *ModuleValue {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	exports := map[string]Value{
		"اضبط_البذرة": builtin("اضبط_البذرة", randomSeed(rng), false),
		"عشوائي":      builtin("عشوائي", randomFloat(rng), false),
		"عدد_عشوائي":  builtin("عدد_عشوائي", randomInt(rng), false),
		"منتظم":       builtin("منتظم", randomUniform(rng), false),
		"اختر":        builtin("اختر", randomChoice(rng), false),
		"رتب_عشوائيا": builtin("رتب_عشوائيا", randomShuffle(rng), false),
	}
	return newNativeModule("عشوائي", exports)
}

func randomSeed(rng *rand.Rand) BuiltinFunc {
	return func(args []Value, _ map[string]Value) (Value, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("اضبط_البذرة يتوقع وسيطا واحدا")
		}
		seed, err := expectInt(args[0], "اضبط_البذرة")
		if err != nil {
			return nil, err
		}
		rng.Seed(int64(seed))
		return NoneValue{}, nil
	}
}

func randomFloat(rng *rand.Rand) BuiltinFunc {
	return func(args []Value, _ map[string]Value) (Value, error) {
		if len(args) != 0 {
			return nil, fmt.Errorf("عشوائي لا يتوقع وسائط")
		}
		return FloatValue{V: rng.Float64()}, nil
	}
}

func randomInt(rng *rand.Rand) BuiltinFunc {
	return func(args []Value, _ map[string]Value) (Value, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("عدد_عشوائي يتوقع وسيطين")
		}
		start, err := expectInt(args[0], "عدد_عشوائي")
		if err != nil {
			return nil, err
		}
		end, err := expectInt(args[1], "عدد_عشوائي")
		if err != nil {
			return nil, err
		}
		if start > end {
			start, end = end, start
		}
		return IntValue{V: start + rng.Intn(end-start+1)}, nil
	}
}

func randomUniform(rng *rand.Rand) BuiltinFunc {
	return func(args []Value, _ map[string]Value) (Value, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("منتظم يتوقع وسيطين")
		}
		start, _, err := expectNumber(args[0], "منتظم")
		if err != nil {
			return nil, err
		}
		end, _, err := expectNumber(args[1], "منتظم")
		if err != nil {
			return nil, err
		}
		if start > end {
			start, end = end, start
		}
		return FloatValue{V: start + rng.Float64()*(end-start)}, nil
	}
}

func randomChoice(rng *rand.Rand) BuiltinFunc {
	return func(args []Value, _ map[string]Value) (Value, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("اختر يتوقع وسيطاً واحداً")
		}
		list, ok := args[0].(ListValue)
		if !ok {
			return nil, fmt.Errorf("اختر يتوقع قائمة")
		}
		if len(list.Elements) == 0 {
			return NoneValue{}, nil
		}
		return list.Elements[rng.Intn(len(list.Elements))], nil
	}
}

func randomShuffle(rng *rand.Rand) BuiltinFunc {
	return func(args []Value, _ map[string]Value) (Value, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("رتب_عشوائيا يتوقع وسيطاً واحداً")
		}
		list, ok := args[0].(ListValue)
		if !ok {
			return nil, fmt.Errorf("رتب_عشوائيا يتوقع قائمة")
		}
		out := make([]Value, len(list.Elements))
		copy(out, list.Elements)
		rng.Shuffle(len(out), func(i, j int) { out[i], out[j] = out[j], out[i] })
		return ListValue{Elements: out}, nil
	}
}

func newTimeModule() *ModuleValue {
	exports := map[string]Value{
		"الان": builtin("الان", timeNow, false),
		"نوم":  builtin("نوم", timeSleep, false),
	}
	return newNativeModule("وقت", exports)
}

func timeNow(args []Value, _ map[string]Value) (Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("الان لا يتوقع وسائط")
	}
	return IntValue{V: int(time.Now().Unix())}, nil
}

func timeSleep(args []Value, _ map[string]Value) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("نوم يتوقع وسيطا واحدا")
	}
	seconds, _, err := expectNumber(args[0], "نوم")
	if err != nil {
		return nil, err
	}
	if seconds < 0 {
		return nil, fmt.Errorf("نوم يتوقع مدة غير سالبة")
	}
	d := time.Duration(seconds * float64(time.Second))
	time.Sleep(d)
	return NoneValue{}, nil
}

func newOSModule() *ModuleValue {
	exports := map[string]Value{
		"احصل_على_المسار_الحالي": builtin("احصل_على_المسار_الحالي", osGetCwd, false),
		"غير_المسار":             builtin("غير_المسار", osChdir, false),
		"موجود":                  builtin("موجود", osExists, false),
		"هل_ملف":                 builtin("هل_ملف", osIsFile, false),
		"هل_مجلد":                builtin("هل_مجلد", osIsDir, false),
		"قائمة_المحتويات":        builtin("قائمة_المحتويات", osListDir, false),
		"احذف":                   builtin("احذف", osRemove, false),
		"انشئ_مجلد":              builtin("انشئ_مجلد", osMkdir, false),
	}
	return newNativeModule("نظام", exports)
}

func osGetCwd(args []Value, _ map[string]Value) (Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("احصل_على_المسار_الحالي لا يتوقع وسائط")
	}
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return StringValue{V: cwd}, nil
}

func osChdir(args []Value, _ map[string]Value) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("غير_المسار يتوقع وسيطا واحدا")
	}
	path, err := expectString(args[0], "غير_المسار")
	if err != nil {
		return nil, err
	}
	return NoneValue{}, os.Chdir(path)
}

func osExists(args []Value, _ map[string]Value) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("موجود يتوقع وسيطا واحدا")
	}
	path, err := expectString(args[0], "موجود")
	if err != nil {
		return nil, err
	}
	_, err = os.Stat(path)
	return BoolValue{V: err == nil}, nil
}

func osIsFile(args []Value, _ map[string]Value) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("هل_ملف يتوقع وسيطا واحدا")
	}
	path, err := expectString(args[0], "هل_ملف")
	if err != nil {
		return nil, err
	}
	info, err := os.Stat(path)
	if err != nil {
		return BoolValue{V: false}, nil
	}
	return BoolValue{V: info.Mode().IsRegular()}, nil
}

func osIsDir(args []Value, _ map[string]Value) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("هل_مجلد يتوقع وسيطا واحدا")
	}
	path, err := expectString(args[0], "هل_مجلد")
	if err != nil {
		return nil, err
	}
	info, err := os.Stat(path)
	if err != nil {
		return BoolValue{V: false}, nil
	}
	return BoolValue{V: info.IsDir()}, nil
}

func osListDir(args []Value, _ map[string]Value) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("قائمة_المحتويات يتوقع وسيطا واحدا")
	}
	path, err := expectString(args[0], "قائمة_المحتويات")
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	items := make([]Value, 0, len(entries))
	for _, entry := range entries {
		items = append(items, StringValue{V: entry.Name()})
	}
	return ListValue{Elements: items}, nil
}

func osRemove(args []Value, _ map[string]Value) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("احذف يتوقع وسيطا واحدا")
	}
	path, err := expectString(args[0], "احذف")
	if err != nil {
		return nil, err
	}
	return NoneValue{}, os.RemoveAll(path)
}

func osMkdir(args []Value, _ map[string]Value) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("انشئ_مجلد يتوقع وسيطا واحدا")
	}
	path, err := expectString(args[0], "انشئ_مجلد")
	if err != nil {
		return nil, err
	}
	return NoneValue{}, os.MkdirAll(path, 0755)
}

func newPathModule() *ModuleValue {
	exports := map[string]Value{
		"ربط_المسار": builtin("ربط_المسار", pathJoin, true),
		"اسم_الملف":  builtin("اسم_الملف", pathBase, false),
		"اسم_المجلد": builtin("اسم_المجلد", pathDir, false),
		"الامتداد":   builtin("الامتداد", pathExt, false),
		"نظف":        builtin("نظف", pathClean, false),
		"مسار_مطلق":  builtin("مسار_مطلق", pathAbs, false),
	}
	return newNativeModule("مسار", exports)
}

func pathJoin(args []Value, _ map[string]Value) (Value, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("ربط_المسار يتوقع وسيطا واحدا على الاقل")
	}
	parts := make([]string, 0, len(args))
	for _, arg := range args {
		part, err := expectString(arg, "ربط_المسار")
		if err != nil {
			return nil, err
		}
		parts = append(parts, part)
	}
	return StringValue{V: filepath.Join(parts...)}, nil
}

func pathBase(args []Value, _ map[string]Value) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("اسم_الملف يتوقع وسيطا واحدا")
	}
	path, err := expectString(args[0], "اسم_الملف")
	if err != nil {
		return nil, err
	}
	return StringValue{V: filepath.Base(path)}, nil
}

func pathDir(args []Value, _ map[string]Value) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("اسم_المجلد يتوقع وسيطا واحدا")
	}
	path, err := expectString(args[0], "اسم_المجلد")
	if err != nil {
		return nil, err
	}
	return StringValue{V: filepath.Dir(path)}, nil
}

func pathExt(args []Value, _ map[string]Value) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("الامتداد يتوقع وسيطا واحدا")
	}
	path, err := expectString(args[0], "الامتداد")
	if err != nil {
		return nil, err
	}
	return StringValue{V: filepath.Ext(path)}, nil
}

func pathClean(args []Value, _ map[string]Value) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("نظف يتوقع وسيطا واحدا")
	}
	path, err := expectString(args[0], "نظف")
	if err != nil {
		return nil, err
	}
	return StringValue{V: filepath.Clean(path)}, nil
}

func pathAbs(args []Value, _ map[string]Value) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("مسار_مطلق يتوقع وسيطا واحدا")
	}
	path, err := expectString(args[0], "مسار_مطلق")
	if err != nil {
		return nil, err
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	return StringValue{V: abs}, nil
}
