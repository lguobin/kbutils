package kbutils

import (
	"bytes"
	"sync"
)

const bufTypeNum = len(bufTypes)

var (
	bufPools [bufTypeNum]sync.Pool
	BuffSize = bufTypes[3]
	bufTypes = [...]int{
		0, 16, 32, 64, 128, 256, 512, 1024, 2048, 5120, 1048576, 5242880, 10485760, 52428800, 104857600,
	}
)

func init() {
	// 必须初始化 bytes.NewBuffer 否则会 panic
	for i := 0; i < bufTypeNum; i++ {
		l := bufTypes[i]
		bufPools[i].New = func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, l))
		}
	}
}

func GetBuff(ss ...int) *bytes.Buffer {
	size := BuffSize
	if len(ss) > 0 {
		size = ss[0]
	}
	if size > 0 {
		if size <= bufTypes[bufTypeNum-1] {
			for i := 0; i < bufTypeNum; i++ {
				if size <= bufTypes[i] {
					return bufPools[i].Get().(*bytes.Buffer)
				}
			}
		}
		return bytes.NewBuffer(make([]byte, 0, size))
	}

	return bufPools[0].Get().(*bytes.Buffer)
}

func PutBuff(buffer *bytes.Buffer) {
	size := buffer.Cap()
	buffer.Reset()
	if size > bufTypes[bufTypeNum-1] {
		bufPools[0].Put(buffer)
		return
	}
	for i := 1; i < bufTypeNum; i++ {
		if size <= bufTypes[i] {
			if size == bufTypes[i] {
				bufPools[i].Put(buffer)
			} else {
				bufPools[i-1].Put(buffer)
			}
			return
		}
	}
}
