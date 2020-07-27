package typhoon

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestHttpGet(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := io.WriteString(w, "Hello World"); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
	go func() {
		http.Serve(ln, mux)
	}()
	target := "http://" + ln.Addr().String()

	rp, err := HTTPGet(10, time.Second, target)
	if err != nil {
		t.Fatal(err)
	}
	if rp.Total == 0 || rp.Status200 < rp.Total*9/10 {
		fmt.Println("Server : ", target)
		fmt.Println(rp)
		t.Fatal("Error")
	}
}

func TestEcho(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := io.Copy(w, r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
	go func() {
		http.Serve(ln, mux)
	}()
	target := "http://" + ln.Addr().String()

	rq, err := http.NewRequest("GET", target, strings.NewReader("Hello World"))
	if err != nil {
		t.Fatal(err)
	}
	handle := func(resp *http.Response) {
		if resp.ContentLength == 0 {
			t.Fatal("Error")
		}
	}
	tp := NewTyphoon(10, time.Second, nil, rq, handle)
	rp := tp.Start()
	if rp.Total == 0 || rp.Status200 < rp.Total*9/10 {
		fmt.Println("Server : ", target)
		fmt.Println(rp)
		t.Fatal("Error")
	}
}
