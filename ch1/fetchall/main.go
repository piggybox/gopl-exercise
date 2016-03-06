// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 17.
//!+

// Fetchall fetches URLs in parallel and reports their times and sizes.
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
    "gopkg.in/redis.v3"
    "strconv"
)

func main() {
	start := time.Now()
    
    client := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    })
    
    count := 0
	
    ch := make(chan string)
    
	for _, url := range os.Args[1:] {
        val, err := client.Get(url).Result()
        if err == redis.Nil {
            go fetch(url, ch, client) // start a goroutine
            count ++
        } else {
            fmt.Printf("%v  %s\n", val, url)
        }
	}
    
    for i := 1; i <= count; i++ {
        fmt.Println(<-ch)
    }
    
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func fetch(url string, ch chan<- string, client *redis.Client) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err) // send to channel ch
		return
	}

	nbytes, err := io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close() // don't leak resources
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}
	secs := time.Since(start).Seconds()
    
    err = client.Set(url, strconv.FormatFloat(secs, 'f', 1, 64), 0).Err()
    if err != nil {
        panic(err)
    }
    
	ch <- fmt.Sprintf("%.2fs  %7d  %s", secs, nbytes, url)
}

//!-
