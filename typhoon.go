package typhoon

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

// RespHandle handle the http response
type RespHandle func(*http.Response)

// HTTPGet start a new Typhoon with 'GET' Method, and returns a testing report
func HTTPGet(numThread int, timeout time.Duration, target string) (*Report, error) {
	req, err := http.NewRequest("Get", target, nil)
	if err != nil {
		return nil, err
	}
	tp := NewTyphoon(numThread, timeout, nil, req, nil)
	return tp.Start(), nil
}

// Typhoon respresent a HTTP benchmarking tester
type Typhoon struct {
	stopCh chan bool

	numThread int
	timeout   time.Duration
	cli       *http.Client
	req       *http.Request
	handle    RespHandle
}

// NewTyphoon returns a new Typhoon
// If cli is nil, http.DefaultClient be used
// If handle is nil, respHandle would be skiped
func NewTyphoon(
	numThread int,
	timeout time.Duration,
	cli *http.Client,
	req *http.Request,
	handle RespHandle) *Typhoon {

	if cli == nil {
		cli = http.DefaultClient
	}
	if numThread <= 0 {
		numThread = 1
	}
	return &Typhoon{
		stopCh: make(chan bool, 1),

		numThread: numThread,
		timeout:   timeout,
		cli:       cli,
		req:       req,
		handle:    handle,
	}
}

// Stop stop the Typhoon
func (tp *Typhoon) Stop() {
	select {
	case tp.stopCh <- true:
		// send signal
	default:
		// nothing
	}
}

// Start start a Typhoon and return the testing report
func (tp *Typhoon) Start() *Report {
	select {
	case <-tp.stopCh:
		// clear signal
	default:
		// nothing
	}
	var wg sync.WaitGroup
	recordCh := make(chan *Record, tp.numThread*2)

	now := time.Now()
	deadline := now.Add(tp.timeout)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	for i := 0; i < tp.numThread; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for time.Now().Before(deadline) {
				rec, resp := tp.do(ctx)
				if rec.Err == nil && tp.handle != nil {
					go tp.handle(resp)
				}
				if rec.Err == nil || !errors.Is(rec.Err, ctx.Err()) {
					recordCh <- rec
				}
			}
		}()
	}
	go func() {
		wg.Wait()
		tp.stopCh <- true
	}()

	rs := make([]*Record, 0)
	for running := true; running; {
		select {
		case <-tp.stopCh:
			running = false
		case r := <-recordCh:
			rs = append(rs, r)
		}
	}
	return newReport(rs, time.Now().Sub(now))
}

func (tp *Typhoon) do(ctx context.Context) (*Record, *http.Response) {
	now := time.Now()
	resp, err := tp.cli.Do(tp.req.Clone(ctx))
	if err != nil {
		return NewErrorRecord(err), nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return NewErrorRecord(err), nil
	}
	resp.Body = ioutil.NopCloser(bytes.NewReader(body))
	return NewRecord(time.Now().Sub(now), resp.StatusCode, int64(len(body))), resp
}
