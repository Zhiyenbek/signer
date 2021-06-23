package main

import (
	"sort"
	"strconv"
	"strings"
	"sync"
)

//var result string = DataSignerCrc32(data) + "~" + DataSignerMd5(data)
func ExecutePipeline(funcs ...job) {
	wg := &sync.WaitGroup{}
	in := make(chan interface{})
	for _, jobFunc := range funcs {
		wg.Add(1)
		out := make(chan interface{})
		go func(jobFunc job, in, out chan interface{}) {
			defer wg.Done()
			defer close(out)
			jobFunc(in, out)
		}(jobFunc, in, out)
		in = out
	}
	wg.Wait()
}
func SingleHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	for i := range in {
		wg.Add(1)
		data, ok := i.(string)
		if !ok {
			val := i.(int)
			data = strconv.Itoa(val)
		}

		md5 := DataSignerMd5(data)
		go func(data string) {
			defer wg.Done()
			inWg := &sync.WaitGroup{}
			inWg.Add(2)
			var begin, end string
			go func(data string) {
				defer inWg.Done()
				begin = DataSignerCrc32(data)
			}(data)
			go func(data string) {
				defer inWg.Done()
				end = DataSignerCrc32(md5)
			}(data)
			inWg.Wait()
			result := begin + "~" + end
			out <- result

		}(data)
	}
	wg.Wait()
}

func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	for i := range in {
		wg.Add(1)
		data, ok := i.(string)
		if !ok {
			val := i.(int)
			data = strconv.Itoa(val)
		}
		go func(data string) {
			defer wg.Done()
			inWg := &sync.WaitGroup{}
			inWg.Add(6)
			result := make([]string, 6)
			for i := '0'; i < '6'; i++ {
				go func(data string, j int, k string) {
					defer inWg.Done()
					resultStr := DataSignerCrc32(k + data)
					result[j-48] = resultStr
				}(data, int(i), string(i))
			}
			inWg.Wait()
			out <- strings.Join(result[:], "")
		}(data)
	}
	wg.Wait()
}
func CombineResults(in, out chan interface{}) {
	var result []string

	for i := range in {
		result = append(result, i.(string))
	}

	sort.Strings(result)
	combineResults := strings.Join(result, "_")
	out <- combineResults
}
