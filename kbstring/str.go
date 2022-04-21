package kbstring

import (
	"fmt"
	"strings"
	"unicode/utf8"
	"unsafe"
)

const (
	// PadRight Right padding character
	PadRight int = iota
	// PadLeft Left padding character
	PadLeft
)

type (
	sliceX struct {
		arr unsafe.Pointer
		len int
		cap int
	}
	stringX struct {
		str unsafe.Pointer
		len int
	}
)

// Len string length (utf8)
func Len(str string) int { return utf8.RuneCountInString(str) }

// Pad String padding
func Pad(raw string, length int, padStr string, padType int) string {
	l := length - Len(raw)
	if l <= 0 {
		return raw
	}
	if padType == PadRight {
		raw = fmt.Sprintf("%s%s", raw, strings.Repeat(padStr, l))
	} else if padType == PadLeft {
		raw = fmt.Sprintf("%s%s", strings.Repeat(padStr, l), raw)
	} else {
		left := 0
		right := 0
		if l > 1 {
			left = l / 2
			right = (l / 2) + (l % 2)
		}

		raw = fmt.Sprintf("%s%s%s", strings.Repeat(padStr, left), raw, strings.Repeat(padStr, right))
	}
	return raw
}

// Substr returns part of a string
func Substr(str string, start int, length ...int) string {
	var size, ll, n, nn int
	if len(length) > 0 {
		ll = length[0] + start
	}
	lb := ll == 0
	if start < 0 {
		start = Len(str) + start
	}
	for i := 0; i < len(str); i++ {
		_, size = utf8.DecodeRuneInString(str[nn:])
		if i < start {
			n += size
		} else if lb {
			break
		}
		if !lb && i < ll {
			nn += size
		} else if lb {
			nn += size
		}
	}
	if !lb {
		return str[n:nn]
	}
	return str[n:]
}

// Buffer Buffer
func Buffer(size ...int) *strings.Builder {
	var _b strings.Builder
	if len(size) > 0 {
		_b.Grow(size[0])
	}
	return &_b
}

// Bytes2String bytes to string
func Bytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// String2Bytes string to bytes
// remark: read only, the structure of runtime changes will be affected, the role of unsafe.Pointer will be changed, and it will also be affected
func String2Bytes(s string) []byte {
	var b []byte
	str := (*stringX)(unsafe.Pointer(&s))
	pbytes := (*sliceX)(unsafe.Pointer(&b))
	pbytes.arr = str.str
	pbytes.len = str.len
	pbytes.cap = str.len
	return b
}

// CamelCaseToSnakeCase camelCase To SnakeCase helloWorld/HelloWorld => hello_world
func CamelCaseToSnakeCase(str string, delimiter ...string) string {
	if str == "" {
		return ""
	}
	sep := []byte("_")
	if len(delimiter) > 0 {
		sep = []byte(delimiter[0])
	}
	strLen := len(str)
	result := make([]byte, 0, strLen*2)
	j := false
	for i := 0; i < strLen; i++ {
		char := str[i]
		if i > 0 && char >= 'A' && char <= 'Z' && j {
			result = append(result, sep...)
		}
		if char != '_' {
			j = true
		}
		result = append(result, char)
	}
	return strings.ToLower(string(result))
}

// XSSClean clean html tag
func XSSClean(str string) string {
	str, _ = RegexReplaceFunc("<[\\S\\s]+?>", str, strings.ToLower)
	str, _ = RegexReplace("<style[\\S\\s]+?</style>", str, "")
	str, _ = RegexReplace("<script[\\S\\s]+?</script>", str, "")
	str, _ = RegexReplace("<[\\S\\s]+?>", str, "")
	str, _ = RegexReplace("\\s{2,}", str, " ")
	return strings.TrimSpace(str)
}

// RegexReplaceFunc replacing matches of the Regexp
func RegexReplaceFunc(pattern string, str string, repl func(string) string) (string, error) {
	r, err := getRegexpCompile(pattern)
	if err == nil {
		str = r.ReplaceAllStringFunc(str, repl)
	}
	return str, err
}

// RegexReplace replacing matches of the Regexp
func RegexReplace(pattern string, str, repl string) (string, error) {
	r, err := getRegexpCompile(pattern)
	if err == nil {
		str = r.ReplaceAllString(str, repl)
	}
	return str, err
}

// TrimSpace TrimSpace
func TrimSpace(s string) string {
	space := [...]uint8{127, 128, 133, 160, 194, 226, 227}
	well := func(s uint8) bool {
		for i := range space {
			if space[i] == s {
				return true
			}
		}
		return false
	}
	for len(s) > 0 {
		if (s[0] <= 31) || s[0] <= ' ' || well(s[0]) {
			s = s[1:]
			continue
		}
		break
	}
	for len(s) > 0 {
		if s[len(s)-1] <= ' ' || (s[len(s)-1] <= 31) || well(s[len(s)-1]) {
			s = s[:len(s)-1]
			continue
		}
		break
	}
	return s
}

// SnakeCaseToCamelCase snakeCase To CamelCase: hello_world => helloWorld
func SnakeCaseToCamelCase(str string, ucfirst bool, delimiter ...string) string {
	if str == "" {
		return ""
	}
	sep := "_"
	if len(delimiter) > 0 {
		sep = delimiter[0]
	}
	slice := strings.Split(str, sep)
	for i := range slice {
		if ucfirst || i > 0 {
			slice[i] = strings.Title(slice[i])
		}
	}
	return strings.Join(slice, "")
}
