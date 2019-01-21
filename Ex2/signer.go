package main

import (
	"fmt"
	"log"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
)

const TH_AMOUNT = 6

type multiHashElement struct {
	id   int
	hash string
}

func main() {
	cancelCh := make(chan struct{})

	// inputData := []int{0, 1, 1, 2, 3, 5, 8}
	inputData := []int{0, 1}
	in := make(chan interface{}, 2)
	out := make(chan interface{}, 2)
	out2 := make(chan interface{}, 2)
	out3 := make(chan interface{}, 2)

	go SingleHash(in, out)
	go MultiHash(out, out2)
	go CombineResults(out2, out3)
	go reader(out3, cancelCh, in)
	for _, fibNum := range inputData {
		in <- fibNum
	}
	// runtime.Gosched()
	time.Sleep(5 * time.Second)
	cancelCh <- struct{}{}

}

func ExecutePipeline(hashignJobs ...job) {
	wg := &sync.WaitGroup{}
	in := make(chan interface{})

	for _, jobItem := range hashignJobs {
		wg.Add(1)
		out := make(chan interface{})
		go func(jobFunc job, in chan interface{}, out chan interface{}, wg *sync.WaitGroup) {
			defer wg.Done()
			defer close(out)
			jobFunc(in, out)
		}(jobItem, in, out, wg)
		in = out
	}

	defer wg.Wait()
}

func reader(out chan interface{}, cancelCh chan struct{}, in chan interface{}) {
	for {
		select {
		case v1 := <-out:
			fmt.Println("Reader: out val", v1)
		case <-cancelCh:
			fmt.Println("Reader: stop channel")
			close(in)
			return
		}
		runtime.Gosched()
	}

}

// crc32(data)+"~"+crc32(md5(data))
func SingleHash(in chan interface{}, out chan interface{}) {
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	for data := range in {
		stringData := fmt.Sprintf("%v", data)

		md5 := DataSignerMd5(stringData)
		wg.Add(1)
		go func(stringData string, md5 string) {
			defer wg.Done()
			log.Print("SingleHash: start work with ", stringData)
			left := <-genCrc32(stringData)
			// runtime.Gosched()
			right := DataSignerCrc32(md5)
			out <- left + "~" + right
		}(stringData, md5)
	}
}

func MultiHash(in chan interface{}, out chan interface{}) {
	// wg := &sync.WaitGroup{}
	// infoLog := log.New(io.Writer,
	// 	"INFO: ",
	// 	log.Ldate|log.Ltime|log.Lshortfile)
	for data := range in {
		temp := make(chan string)
		wgTH := &sync.WaitGroup{}
		inputCh := make(chan multiHashElement)
		inputLine := data.(string)
		// fmt.Println("MultiHash: got line =", inputLine)
		// infoLog.Println("MultiHash: got line =", inputLine)
		log.Print("MultiHash: got line =", inputLine)

		wgTH.Add(TH_AMOUNT)
		for i := 0; i < TH_AMOUNT; i++ {
			go genMultiHash(i, inputLine, inputCh, wgTH)
		}
		go func(wg *sync.WaitGroup, ch chan multiHashElement) {
			defer close(ch)
			wg.Wait()
		}(wgTH, inputCh)
		go sortMultiHash(inputCh, temp)
		for hash := range temp {
			out <- hash
		}
	}
}

func CombineResults(in, out chan interface{}) {
	var lines []string
	var result string

	for line := range in {
		fmt.Println("CombineResults: get from input ", line)
		lines = append(lines, (line).(string))
	}

	sort.Strings(lines)
	result = strings.Join(lines, "_")
	fmt.Println("CombineResults: send to output ", result)

	out <- result
}

func genCrc32(data string) chan string {
	result := make(chan string, 1)
	go func(out chan<- string) {
		out <- DataSignerCrc32(data)
	}(result)

	return result
}

func genMultiHash(i int, inputLine interface{}, resultCh chan multiHashElement, wg *sync.WaitGroup) {
	defer wg.Done()
	resultCh <- multiHashElement{id: i, hash: DataSignerCrc32(fmt.Sprintf("%v%v", i, inputLine))}
}

func sortMultiHash(hashElements chan multiHashElement, out chan string) {
	result := map[int]string{}
	var data []int

	for hash := range hashElements {
		fmt.Println("sortMultiHash: ", hash)
		result[hash.id] = hash.hash
		data = append(data, hash.id)
	}
	sort.Ints(data)

	var results []string
	for i := range data {
		results = append(results, result[i])
	}
	fmt.Println("sortMultiHash: result ", result)
	out <- strings.Join(results, "")
	close(out)
}
