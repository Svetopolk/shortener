package logging

import (
	"log"
	"runtime"
	"strconv"
	"strings"
)

func fileName(original string) string {
	i := strings.LastIndex(original, "/")
	if i == -1 {
		return original
	}
	return original[i+1:]
}

func stackRecord(depth int) string {
	function, file, line, _ := runtime.Caller(depth + 1)
	name := runtime.FuncForPC(function).Name()
	return name + "(" + fileName(file) + ":" + strconv.Itoa(line) + ")"
}

func Enter() {
	log.Println(stackRecord(1))
}

func Exit() {
	log.Println("~" + stackRecord(1))
}
