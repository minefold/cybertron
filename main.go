package main

import (
	// "fmt"
	"fmt"
	"github.com/rcrowley/goagain"
	"github.com/whatupdave/s3/s3util"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
	"time"
)

var (
	s3url string
)

func main() {
	go memDebugger()

	var (
		err  error
		l    net.Listener
		ppid int
	)
	port := os.Getenv("PORT")
	if port == "" {
		port = "7070"
	}

	l, ppid, err = goagain.GetEnvs()

	if nil != err {
		laddr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:"+port)
		if nil != err {
			log.Fatalln(err)
		}
		log.Printf("listening on %v", laddr)
		l, err = net.ListenTCP("tcp", laddr)
		if nil != err {
			log.Fatalln(err)
		}

		go serve(l)
	} else {
		log.Printf("resuming listening on %v", l.Addr())
		go serve(l)

		if err := goagain.KillParent(ppid); nil != err {
			log.Fatalln(err)
		}
	}

	// Block the main goroutine awaiting signals.
	if err := goagain.AwaitSignals(l); nil != err {
		log.Fatalln(err)
	}

	// Do whatever's necessary to ensure a graceful exit like waiting for
	// goroutines to terminate or a channel to become closed.

	// In this case, we'll simply stop listening and wait one second.
	if err := l.Close(); nil != err {
		log.Fatalln(err)
	}
	time.Sleep(1e9)

}

func serve(l net.Listener) {
	s3url = os.Getenv("S3_URL")
	s3util.DefaultConfig.AccessKey = os.Getenv("AWS_ACCESS_KEY")
	s3util.DefaultConfig.SecretKey = os.Getenv("AWS_SECRET_KEY")

	store := NewS3Store(s3url)

	http.Serve(l, &Router{Store: store})
}

func memDebugger() {
	t := time.NewTicker(5 * time.Second)
	for _ = range t.C {
		var ms runtime.MemStats
		runtime.ReadMemStats(&ms)
		fmt.Printf("goroutines: %d  alloc: %d\n", runtime.NumGoroutine(), ms.Alloc)
	}
}
