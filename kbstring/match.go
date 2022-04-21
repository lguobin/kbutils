package kbstring

import (
	"regexp"
	"sync"
	"time"
	"unicode/utf8"
)

var (
	rk         sync.RWMutex
	regexCache = map[string]*regexMapStruct{}
)

type regexMapStruct struct {
	Value *regexp.Regexp
	Time  int64
	sync.RWMutex
}

// RegexMatch check for match
func RegexMatch(pattern string, str string) bool {
	if r, err := getRegexpCompile(pattern); err == nil {
		return r.Match(String2Bytes(str))
	}
	return false
}

// RegexExtract extract matching text
func RegexExtract(pattern string, str string) ([]string, error) {
	r, err := getRegexpCompile(pattern)
	if err == nil {
		return r.FindStringSubmatch(str), nil
	}
	return nil, err
}

func getRegexpCompile(pattern string) (r *regexp.Regexp, err error) {
	rk.RLock()
	var data *regexMapStruct
	var ok bool
	data, ok = regexCache[pattern]
	rk.RUnlock()
	if ok {
		r = data.Value
		return
	}
	r, err = regexp.Compile(pattern)
	if err != nil {
		return
	}
	rk.Lock()
	regexCache[pattern] = &regexMapStruct{Value: r, Time: time.Now().Unix()}
	rk.Unlock()
	return
}

//字符串匹配
func Match(str, pattern string) bool {
	if pattern == "*" {
		return true
	}
	return deepMatch(str, pattern)
}

func deepMatch(str, pattern string) bool {
	for len(pattern) > 0 {
		if pattern[0] > 0x7f {
			return deepMatchRune(str, pattern)
		}
		switch pattern[0] {
		default:
			if len(str) == 0 {
				return false
			}
			if str[0] > 0x7f {
				return deepMatchRune(str, pattern)
			}
			if str[0] != pattern[0] {
				return false
			}
		case '?':
			if len(str) == 0 {
				return false
			}
		case '*':
			return deepMatch(str, pattern[1:]) ||
				(len(str) > 0 && deepMatch(str[1:], pattern))
		}
		str = str[1:]
		pattern = pattern[1:]
	}
	return len(str) == 0 && len(pattern) == 0
}

func deepMatchRune(str, pattern string) bool {
	var sr, pr rune
	var srsz, prsz int

	x7f := func(isStr bool) (r rune, p int) {
		var s uint8
		if isStr {
			s = str[0]
		} else {
			s = pattern[0]
		}
		if str[0] > 0x7f {
			r, p = utf8.DecodeRuneInString(str)
		} else {
			r, p = rune(s), 1
		}
		return
	}

	if len(str) > 0 {
		sr, srsz = x7f(true)
	} else {
		sr, srsz = utf8.RuneError, 0
	}
	if len(pattern) > 0 {
		pr, prsz = x7f(false)
	} else {
		pr, prsz = utf8.RuneError, 0
	}
	for pr != utf8.RuneError {
		switch pr {
		default:
			if srsz == utf8.RuneError {
				return false
			}
			if sr != pr {
				return false
			}
		case '?':
			if srsz == utf8.RuneError {
				return false
			}
		case '*':
			return deepMatchRune(str, pattern[prsz:]) ||
				(srsz > 0 && deepMatchRune(str[srsz:], pattern))
		}
		str = str[srsz:]
		pattern = pattern[prsz:]
		if len(str) > 0 {
			sr, srsz = x7f(true)
		} else {
			sr, srsz = utf8.RuneError, 0
		}
		if len(pattern) > 0 {
			pr, prsz = x7f(false)
		} else {
			pr, prsz = utf8.RuneError, 0
		}
	}

	return srsz == 0 && prsz == 0
}
