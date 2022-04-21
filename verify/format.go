package verify

import (
	"strings"

	"github.com/lguobin/kbutils/kbstring"
	"golang.org/x/crypto/bcrypt"
)

// Trim remove leading and trailing spaces
func (v Engine) Trim() Engine {
	return pushQueue(v, func(v Engine) Engine {
		if notEmpty(&v) {
			v.value = kbstring.TrimSpace(v.value)
		}

		return v
	})
}

// RemoveSpace remove all spaces
func (v Engine) RemoveSpace() Engine {
	return pushQueue(v, func(v Engine) Engine {
		if notEmpty(&v) {
			v.value = strings.Replace(v.value, " ", "", -1)
		}
		return v
	})
}

// Replace replace text
func (v Engine) Replace(old, new string, n int) Engine {
	return pushQueue(v, func(v Engine) Engine {
		if notEmpty(&v) {
			v.value = strings.Replace(v.value, old, new, n)
		}
		return v
	})
}

// ReplaceAll replace all text
func (v Engine) ReplaceAll(old, new string) Engine {
	return v.Replace(old, new, -1)
}

// XSSClean clean html tag
func (v Engine) XSSClean() Engine {
	return pushQueue(v, func(v Engine) Engine {
		if notEmpty(&v) {
			v.value = kbstring.XSSClean(v.value)
		}
		return v
	})
}

// SnakeCaseToCamelCase snakeCase To CamelCase: hello_world => helloWorld
func (v Engine) SnakeCaseToCamelCase(ucfirst bool, delimiter ...string) Engine {
	return pushQueue(v, func(v Engine) Engine {
		if notEmpty(&v) {
			v.value = kbstring.SnakeCaseToCamelCase(v.value, ucfirst, delimiter...)
		}
		return v
	})
}

// CamelCaseToSnakeCase camelCase To SnakeCase helloWorld/HelloWorld => hello_world
func (v Engine) CamelCaseToSnakeCase(delimiter ...string) Engine {
	return pushQueue(v, func(v Engine) Engine {
		if notEmpty(&v) {
			v.value = kbstring.CamelCaseToSnakeCase(v.value, delimiter...)
		}
		return v
	})
}

// EncryptPassword encrypt the password
func (v Engine) EncryptPassword(cost ...int) Engine {
	return pushQueue(v, func(v Engine) Engine {
		if notEmpty(&v) {
			bcost := bcrypt.DefaultCost
			if len(cost) > 0 {
				bcost = cost[0]
			}
			if bytes, err := bcrypt.GenerateFromPassword(kbstring.String2Bytes(v.value), bcost); err == nil {
				v.value = kbstring.Bytes2String(bytes)
			} else {
				v.err = err
			}
		}
		return v
	})
}
