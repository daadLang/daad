package interpreter

import (
	"fmt"
)

var formaters = map[string]string{
	"ر": "d", // رقم
	"ع": "f", // عدد عشري
	"ن": "s", // نص
	"م": "t", // منطقي
	"ث": "b", // ثنائي
	"و": "o", // ثماني
	"س": "x", // سداسي عشري
	"ل": "e", // كتابة علمية
}

func GetFormater(format string) string {
	if formater, ok := formaters[format]; ok {
		return formater
	}
	return ""
}

func arabicFormat(arabicFormat string) string {
	txt := []rune(arabicFormat)
	txtlen := len(txt)

	for index, char := range txt {
		if char == '٪' || char == '%' {
			if index+1 < txtlen {
				formater := GetFormater(string(txt[index+1]))
				if formater != "" {
					txt[index] = '%'
					txt[index+1] = rune(formater[0])
				}
			}
		}
	}
	return string(txt)
}

// طباعة
func Print(args FuncArgsAdapter, kwargs FuncKwargsAdapter) {

	fmt.Print(args.Args...)
}

// طباعة_منسقة
func Printf(format string, args ...interface{}) {

	fmt.Printf(arabicFormat(format), args...)
}

// صياغة_منسقة
func Sprintf(format string, args ...interface{}) string {
	return fmt.Sprintf(arabicFormat(format), args...)
}
