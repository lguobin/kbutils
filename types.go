package kbutils

import (
	"encoding/json"
	"strconv"

	"github.com/lguobin/kbutils/kbstring"
)

// IfVal Simulate ternary calculations, pay attention to handling no variables or indexing problems
func IfVal(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}

///定义接口
type StringVariable interface{ String() string }

func ToString(i interface{}) string {
	if i == nil {
		return ""
	}
	if f, ok := i.(StringVariable); ok {
		return f.String()
	}
	switch value := i.(type) {
	case int:
		return strconv.Itoa(value)
	case int8:
		return strconv.Itoa(int(value))
	case int16:
		return strconv.Itoa(int(value))
	case int32:
		return strconv.Itoa(int(value))
	case int64:
		return strconv.Itoa(int(value))
	case uint:
		return strconv.FormatUint(uint64(value), 10)
	case uint8:
		return strconv.FormatUint(uint64(value), 10)
	case uint16:
		return strconv.FormatUint(uint64(value), 10)
	case uint32:
		return strconv.FormatUint(uint64(value), 10)
	case uint64:
		return strconv.FormatUint(value, 10)
	case float32:
		return strconv.FormatFloat(float64(value), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(value)
	case string:
		return value
	case []byte:
		return kbstring.Bytes2String(value)
	default:
		if f, ok := value.(StringVariable); ok {
			return f.String()
		}
		jsonContent, _ := json.Marshal(value)
		return string(jsonContent)
	}
}
