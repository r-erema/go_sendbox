package signer

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	strings2 "strings"
	"sync"
	"time"
)

func ExecutePipeline(jobs ...job) {

	defer Duration(time.Now(), "ExecutePipeline")

	var cachedIn, newOut chan interface{}
	for _, j := range jobs {
		newIn := make(chan interface{})
		newOut = make(chan interface{})

		go pipe(newIn, newOut)
		go worker(j, cachedIn, newOut)

		cachedIn = newIn
	}

	if newOut != nil {
		fmt.Println(<-newOut)
	}
}

func pipe(in, out chan interface{}) {
	for data := range out {
		in <- data
	}
	close(in)
}

func worker(j job, in, out chan interface{}) {
	defer close(out)
	j(in, out)
}

func SingleHash(in, out chan interface{}) {

	defer Duration(time.Now(), "SingleHash")

	wg := &sync.WaitGroup{}
	for data := range in {
		data := fmt.Sprintf("%v", data)
		wg.Add(1)
		md5 := DataSignerMd5(data)

		go func(wg *sync.WaitGroup, data, md5 string) {
			ch := make(chan string)
			ch2 := make(chan string)

			go Crc32(data, ch)
			go Crc32(md5, ch2)

			go func(wg *sync.WaitGroup, ch, ch2 chan string) {
				out <- <-ch + "~" + <-ch2
				wg.Done()
			}(wg, ch, ch2)

		}(wg, data, md5)
	}
	wg.Wait()
}

func Crc32(data string, ch chan string) {
	ch <- DataSignerCrc32(data)
}

func MultiHash(in, out chan interface{}) {

	defer Duration(time.Now(), "MultiHash")

	hashesCount := 6

	wg := &sync.WaitGroup{}
	for data := range in {

		wg.Add(1)
		wg2 := &sync.WaitGroup{}
		data := data.(string)

		go func(wg, wg2 *sync.WaitGroup, data string) {
			strings := make([]string, hashesCount)
			m := &sync.Mutex{}
			for i := 0; i < hashesCount; i++ {
				wg2.Add(1)
				go func(i int, data string, strings *[]string, m *sync.Mutex) {
					result := DataSignerCrc32(strconv.Itoa(i) + data)
					m.Lock()
					defer m.Unlock()
					(*strings)[i] = result
					wg2.Done()
				}(i, data, &strings, m)
			}
			wg2.Wait()
			out <- strings2.Join(strings, "")
			wg.Done()
		}(wg, wg2, data)
	}
	wg.Wait()
}

func CombineResults(in, out chan interface{}) {

	defer Duration(time.Now(), "CombineResults")

	var strings []string
	for str := range in {
		strings = append(strings, str.(string))
	}
	sort.Strings(strings)
	result := strings2.Join(strings, "_")
	out <- result
}

func Duration(invocation time.Time, name string) {
	log.Printf("%s execution time: %s", name, time.Since(invocation))
}
