package converter

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
)

func Debug(s string, args ...interface{}) {
	if debugMode {
		log.Printf(s, args...)
	}
}

func Exit(v interface{}, args ...interface{}) {
	switch v := v.(type) {
	case error:
		err := v
		pc := make([]uintptr, 1)
		runtime.Callers(2, pc)
		funcName := runtime.FuncForPC(pc[0]).Name()
		funcName = LastPart(funcName, path.Dirname()+"/")

		_, file, line, ok := runtime.Caller(1)
		if !ok {
			log.Fatal(fmt.Errorf("get file & line failed"))
		}
		file = LastPart(file, path.Dirname()+"/")
		fileLine := fmt.Sprintf("%s:%d", file, line)

		newErr := fmt.Errorf("\n - line: %s\n - func: %s \n - error: %w \n%v",
			fileLine, funcName, err, string(debug.Stack()))

		log.Fatal(newErr)
	case string:
		format := v
		pc := make([]uintptr, 1)
		runtime.Callers(2, pc)
		funcName := runtime.FuncForPC(pc[0]).Name()
		funcName = LastPart(funcName, path.Dirname()+"/")

		_, file, line, ok := runtime.Caller(1)
		if !ok {
			log.Fatal(fmt.Errorf("get file & line failed"))
		}
		file = LastPart(file, path.Dirname()+"/")
		fileLine := fmt.Sprintf("%s:%d", file, line)

		err := fmt.Errorf(format, args...)
		newErr := fmt.Errorf("\n - line: %s\n - func: %s \n - error: %w \n%v",
			fileLine, funcName, err, string(debug.Stack()))

		log.Fatal(newErr)
	}
}

func ToSlice(arr interface{}) []interface{} {
	ret := make([]interface{}, 0)
	v := reflect.ValueOf(arr)
	if v.Kind() != reflect.Slice {
		ret = append(ret, arr)
		return ret
	}
	l := v.Len()
	for i := 0; i < l; i++ {
		ret = append(ret, v.Index(i).Interface())
	}
	return ret
}

func LastPart(s string, sep string) string {
	lastIndex := strings.LastIndex(s, sep)
	if lastIndex < 0 {
		return s
	}
	return s[lastIndex+len(sep):]
}

// MatchParentDir returns target's directory's full path,
// returning error if `dir`'s parent dir names don't match `target`
func MatchParentDir(dir string, target string) (string, error) {
	var currentDir string
	var file string
	for {
		currentDir = filepath.Dir(dir)
		file = filepath.Base(dir)

		// Match target directory
		if file == target {
			return dir, nil
		}

		// Reach the top of directory
		if currentDir == dir {
			return "", fmt.Errorf(
				"diretory `%s` doesn't match `%s`", dir, target)
		}

		dir = currentDir
	}
}

// Exists returns if exists a dir or file
func Exists(s string) bool {
	_, err := os.Stat(s) // os.Stat获取文件信息
	if err != nil {
		if os.IsNotExist(err) {
			return false
		} else {
			panic(err)
		}
	}
	return true
}

// float2String converts a scientific notation string to a decimal notation string
// It also removes trailing zeros in the fractional part
func float2String(f string) string {
	// Parse the string to float64
	floatVal, err := strconv.ParseFloat(f, 64)
	if err != nil {
		fmt.Println("Error parsing float:", err)
		return ""
	}
	// Use fmt.Sprintf to format the float64 as a decimal string
	decimalStr := fmt.Sprintf("%f", floatVal)

	// Remove trailing zeros in the fractional part
	decimalStr = strings.TrimRight(decimalStr, "0")

	// Remove the decimal point if it is at the end
	decimalStr = strings.TrimRight(decimalStr, ".")

	return decimalStr
}
