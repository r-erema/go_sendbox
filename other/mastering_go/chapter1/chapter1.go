package chapter1

import (
	"io"
	"log"
	"os"
	"strconv"
)

const StopWord = "END"

func sum() (result float64) {
	for _, arg := range os.Args {
		number, err := strconv.ParseFloat(arg, 64)
		if err != nil {
			continue
		}
		result += number
	}
	return
}

func average() (result float64) {
	var validNumsCount float64
	for _, arg := range os.Args {
		number, err := strconv.ParseFloat(arg, 64)
		if err != nil {
			continue
		}
		result += number
		validNumsCount++
	}
	result = result / validNumsCount
	return
}

func sumInt() (result int64) {
	for _, arg := range os.Args {
		if arg == StopWord {
			break
		}
		number, err := strconv.ParseInt(arg, 10, 64)
		if err != nil {
			continue
		}
		result += number
	}
	return
}

func customLog(data string, writers ...io.Writer) {
	w := io.MultiWriter(writers...)
	iLog := log.New(w, "customLogLineNumber ", log.LstdFlags)
	iLog.SetFlags(log.LstdFlags)
	iLog.Println(data)
}
