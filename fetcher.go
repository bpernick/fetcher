package main

import (
	"fmt"
	"flag"
	"net"
	"net/url"
	"log"
	"io/ioutil"
	"os"
	"time"
	"math"
	"sort"
	"regexp"
	"strconv"
)
func http (str string) (string, error, int, int) {
	address, err := url.Parse(str)
	if err != nil {
		return "", err, 0, 0
	}
	host := string(address.Host)
	end  := string(host[len(host)-3])
	
	if end != "/" {
		host += "/"
	}
	
	connection, err := net.Dial("tcp", address.Host+":80")
	if err != nil {
		return "", err, 0, 0
	}
	
	request := fmt.Sprintf("GET %v HTTP/1.1\r\nHost: %v\r\nConnection: close\r\n\r\n", address.Path, address.Host)
	if err != nil {
		return "", err, 0, 0
	}
	
	_, err = connection.Write([]byte(request))
		if err != nil {
			return "", err, 0, 0
		}
		
	response, err := ioutil.ReadAll(connection)
    if err != nil {
			return "", err, 0, 0
		}

		data := string(response)
		bytes := len(response)
		re := regexp.MustCompile(`(HTTP/[0-2]\.[0-9]\s)(.*?)(\s)`)
		status := re.FindStringSubmatch(data)

		if len(status) >= 3 {
			statusInt, err := strconv.Atoi(status[2])
			if err != nil {
				return "", err, 0, 0
			}
			return data, nil, statusInt, bytes
		}
    return data, nil, 0, bytes
}


func main () {

	stra := flag.String("url", "", "url")
	helpa := flag.Bool("help", false, "help")
	profilea := flag.Int("profile", 0, "profile")

	flag.Parse()

	str := *stra
	help := *helpa
	profile := *profilea

	if help == true {
		fmt.Println("\n\rFetcher is a tool for making API requests to a given endpoint\nUsage: fetcher [options] [arguments]\n Options:\n--help: see list of commands and usage\n Arguments:\n--url<string>: makes an http request to the given endpoint and prints the result\n--profile<integer>: makes the given number of requests to the endpoint and prints diagnostic information\n\r")
	}

	if str == "" {
		if profile != 0 {
			log.Fatal("In order to use --profile, you must provide a value for --url")
		}
		os.Exit(0)
	}

	if profile != 0 {

		times := make ([] int, 0, profile)
		lengths := make ([] int, 0, profile)
		errors := make(map[string]int)

		for i := 0; i < profile; i++ {
			start := time.Now()
			_, err, status, bytes := http(str)
			if err != nil {
				log.Fatal(err)
			}
			times = append(times, int(time.Since(start)/1000000))
			lengths = append(lengths, bytes)

			if status > 399 {
				errors[strconv.Itoa(status)] += 1
			}
		}
		sort.Slice(times, func(i, j int) bool { return times[i] < times[j] })
		sort.Slice(lengths, func(i, j int) bool { return lengths[i] < lengths[j] })

		fastest := times[0]
		slowest := times[profile - 1]
		smallest := lengths[0]
		largest := lengths[profile - 1]
		var median int
		var mean float64

		half := float64(profile)/2.0
		intHalf := int(math.Floor(half))

		if profile % 2 == 0 {
			median = (times[intHalf-1] + times[intHalf]) / 2
		} else {
			median = times[intHalf]
		}
		for i := 0; i < profile; i++ {
			mean += float64(times[i])
		}
		mean = mean/float64(profile)
		fmt.Printf("Fastest request: %v ms\nSlowest request: %v ms\nMean: %v ms\nMedian: %v ms\nErrors: %v\nSmallest response: %v bytes\nLargest resonse: %v bytes\n\r\n", fastest, slowest, mean, median, errors, smallest, largest)
		os.Exit(0)
	}
	response, err, _, _ := http(str)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response)
}