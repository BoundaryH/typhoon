package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"typhoon"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "conf", "", "configure file")

	cj := &typhoon.ConfigJSON{}
	flag.IntVar(&cj.NumCPU, "cpu", 0, "the numbers of cpu used")
	flag.IntVar(&cj.NumThread, "t", 10, "the numbers of threads used")
	flag.StringVar(&cj.Duration, "d", "2s", "duration of the test e.g. 2s, 2m, 1h")

	flag.StringVar(&cj.Method, "m", "GET", "the http request method")
	flag.StringVar(&cj.Cookie, "ck", "", "Cookie of http request")
	flag.StringVar(&cj.UserAgent, "ua", "", "User-Agent of http request")
	flag.StringVar(&cj.BodyFile, "body", "", "request body file")

	flag.BoolVar(&cj.DisableCompression, "nc", false, "disable compress")
	flag.BoolVar(&cj.KeepAlive, "nk", false, "disable keep-alive")
	flag.BoolVar(&cj.SkipTLSVerify, "nt", false, "disable tls verify")

	flag.StringVar(&cj.ServerCertFile, "ca", "", "PEM encoded CA's certificate file")
	flag.StringVar(&cj.ClientCertFile, "cert", "", "PEM encoded client certificate file")
	flag.StringVar(&cj.ClientKeyFile, "key", "", "PEM encoded client private key file")
	flag.Parse()

	if configFile != "" {
		var err error
		cj, err = typhoon.ReadConfigJSON(configFile)
		if err != nil {
			log.Fatal(err)
		}
	} else if len(flag.Args()) > 0 {
		cj.Target = flag.Args()[0]
	} else {
		flag.Usage()
		return
	}

	conf, err := cj.Config()
	if err != nil {
		log.Fatal(err)
	}
	tp, err := conf.Typhoon()
	if err != nil {
		log.Fatal(err)
	}

	signChan := make(chan os.Signal, 1)
	signal.Notify(signChan, os.Interrupt)
	go func() {
		<-signChan
		tp.Stop()
	}()

	fmt.Printf("Testing : %s\n", conf.Target)
	fmt.Println(tp.Start())
}
