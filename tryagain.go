package kbutils

import (
	"fmt"
	"os"
)

// TryCatch exception capture
func TryCatch(fn func() error) (err error) {
	defer func() {
		if recoverErr := recover(); recoverErr != nil {
			if e, ok := recoverErr.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("%v", recoverErr)
			}
		}
	}()
	err = fn()
	return
}

// Deprecated: please use TryCatch
// Try exception capture
func Try(fn func(), catch func(e interface{}), finally ...func()) {
	if len(finally) > 0 {
		defer func() {
			finally[0]()
		}()
	}
	defer func() {
		if err := recover(); err != nil {
			if catch != nil {
				catch(err)
			} else {
				panic(err)
			}
		}
	}()
	fn()
}

// CheckErr CheckErr
func CheckErr(err error, exit ...bool) {
	if err != nil {
		if len(exit) > 0 && exit[0] {
			fmt.Println(err)
			os.Exit(1)
			return
		}
		panic(err)
	}
}
