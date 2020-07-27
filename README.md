# Typhoon

> HTTP benchmarking tool and library

Typhoon is a small benchmarking tool similar to [adjust/go-wrk](https://github.com/adjust/go-wrk). 

It is also a benchmarking **library**, which can be used in other project written in Go.

## Build
```
cd cmd
go build -o typhoon
```

## Example
```
typhoon -t 100 -d 20s http://localhost:8080
```
This runs a benchmark for 20 seconds using 100 goroutines.

## Command Line Options
```
  -t int
    	the numbers of threads used (default 10)
  -d string
    	duration of the test e.g. 2s, 2m, 1h (default "2s")
  -conf string
    	configure file
  -cpu int
    	the numbers of cpu used
  -ca string
    	PEM encoded CA's certificate file
  -cert string
    	PEM encoded client certificate file
  -key string
    	PEM encoded client private key file
  -m string
    	the http request method (default "GET")
  -ua string
    	User-Agent of http request
  -body string
    	request body file
  -ck string
    	Cookie of http request
  -nc
    	disable compress
  -nk
    	disable keep-alive
  -nt
    	disable tls verify
```

## Configuration Example
Configuration file is written in JSON.
```
{
    "Target": "https://localhost:8080",
    "NumCPU": 0,
    "NumThread": 10,
    "Duration": "20s",
    "Method": "GET",
    "Header": {
        "Referer": "https://www.google.com/"
    },
    "Cookie": "UID=s0YAAQspMRh; SID=AJi4QfGep3ZMNueWz5buMc4cFPd",
    "UserAgent": "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:47.0) Gecko/20100101 Firefox/47.0",
    "BodyFile": "body.txt",
    "DisableCompression": false,
    "KeepAlive": true,
    "SkipTLSVerify": false,
    "ServerCertFile": "server.crt",
    "ClientCertFile": "client.crt",
    "ClientKeyFile": "client.key"
  }
```


## References

* [wg/wrk](https://github.com/wg/wrk)
* [adjust/go-wrk](https://github.com/adjust/go-wrk)
* [tsliwowicz/go-wrk](https://github.com/tsliwowicz/go-wrk)

## License

This software is licensed under the MIT License.