package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gcinterceptor/gci-go/httphandler"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	UseGCI     bool `default:"false" envconfig:"USE_GCI"`
	Port       int  `default:"3000" envconfig:"PORT"`
	WindowSize int  `default:"1" envconfig:"WINDOW_SIZE"`
	MsgSize    int  `default:"1024" envconfig:"MSG_SIZE"`
}

var (
	msgCount = 0
	buffer   [][]byte
	mu       sync.Mutex
)

func main() {
	var c config
	err := envconfig.Process("", &c)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("Configuration: %+v\n", c)
	// Inspiration: https://making.pusher.com/golangs-real-time-gc-in-theory-and-practice/
	buffer = make([][]byte, c.WindowSize)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m := make([]byte, c.MsgSize)
		for i := range m {
			m[i] = byte(i)
		}
		mu.Lock()
		buffer[msgCount] = m
		msgCount = (msgCount + 1) % c.WindowSize
		mu.Unlock()
	})
	if c.UseGCI {
		http.Handle("/", httphandler.GCI(handler))
		fmt.Println("==< Using GCI >==")
	} else {
		http.Handle("/", handler)
	}
	http.ListenAndServe(fmt.Sprintf(":%d", c.Port), nil)
}
