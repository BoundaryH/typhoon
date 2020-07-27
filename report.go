package typhoon

import (
	"bytes"
	"fmt"
	"sort"
	"time"
)

// Report represents the statistic report
type Report struct {
	Records []*Record

	Total     int
	TotalTime time.Duration
	ReqPerSec float64

	AvgReqTime time.Duration
	MedianReq  time.Duration
	FastestReq time.Duration
	The99thReq time.Duration
	SlowestReq time.Duration

	AvgBodySize float64
	BytePerSec  float64

	StatusErr int
	Status200 int
	Status300 int
	Status400 int
	Status500 int
}

func newReport(rs []*Record, d time.Duration) *Report {
	var bufTotal int64
	var timeTotal time.Duration
	var sErr, s200, s300, s400, s500 int
	var count int

	sortList := make([]*Record, 0, len(rs))
	for _, r := range rs {
		code := r.StatusCode
		switch {
		case code < 200:
			sErr++
		case code < 300:
			s200++
		case code < 400:
			s300++
		case code < 500:
			s400++
		default:
			s500++
		}

		if r.Err == nil {
			bufTotal += r.Length
			timeTotal += r.Duration
			sortList = append(sortList, r)
			count++
		}
	}
	sort.Slice(sortList, func(i, j int) bool {
		return sortList[i].Duration < sortList[j].Duration
	})

	var avgReqTime time.Duration
	var avgBodySize float64
	if count > 0 {
		avgReqTime = time.Duration(float64(timeTotal) / float64(count))
		avgBodySize = float64(bufTotal) / float64(count)
	}

	var medianTime, fastestTime, the99thTime, slowestTime time.Duration
	if l := len(sortList); l > 0 {
		medianTime = sortList[l/2].Duration
		fastestTime = sortList[0].Duration
		the99thTime = sortList[l*99/100].Duration
		slowestTime = sortList[l-1].Duration
	}

	return &Report{
		Records:   rs,
		Total:     len(rs),
		TotalTime: d,

		ReqPerSec:  float64(count) / d.Seconds(),
		AvgReqTime: avgReqTime,
		MedianReq:  medianTime,
		FastestReq: fastestTime,
		The99thReq: the99thTime,
		SlowestReq: slowestTime,

		AvgBodySize: avgBodySize,
		BytePerSec:  float64(bufTotal) / d.Seconds(),

		StatusErr: sErr,
		Status200: s200,
		Status300: s300,
		Status400: s400,
		Status500: s500,
	}
}

func (rp *Report) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Total calls  : %8d\n", rp.Total)
	fmt.Fprintf(&buf, "Total time   : %8.2f s\n", rp.TotalTime.Seconds())
	fmt.Fprintf(&buf, "Requests/sec : %8.2f\n", rp.ReqPerSec)
	fmt.Fprintf(&buf, "\n")

	fmt.Fprintf(&buf, "Avg Req Time : %8.2f ms\n", float64(rp.AvgReqTime.Microseconds())/1e3)
	fmt.Fprintf(&buf, "Fastest  Req : %8.2f ms\n", float64(rp.FastestReq.Microseconds())/1e3)
	fmt.Fprintf(&buf, "Median   Req : %8.2f ms\n", float64(rp.MedianReq.Microseconds())/1e3)
	fmt.Fprintf(&buf, "99%%      Req : %8.2f ms\n", float64(rp.The99thReq.Microseconds())/1e3)
	fmt.Fprintf(&buf, "Slowest  Req : %8.2f ms\n", float64(rp.SlowestReq.Microseconds())/1e3)
	fmt.Fprintf(&buf, "\n")

	fmt.Fprintf(&buf, "Avg body size: %8.2f KB\n", rp.AvgBodySize/1e3)
	fmt.Fprintf(&buf, "Transfer /sec: %8.2f MB\n", rp.BytePerSec/1e6)
	fmt.Fprintf(&buf, "\n")

	l := float64(len(rp.Records))
	fmt.Fprintf(&buf, "20X Responses: %8d (%6.2f%%)\n", rp.Status200, float64(rp.Status200*100)/l)
	fmt.Fprintf(&buf, "30X Responses: %8d (%6.2f%%)\n", rp.Status300, float64(rp.Status300*100)/l)
	fmt.Fprintf(&buf, "40X Responses: %8d (%6.2f%%)\n", rp.Status400, float64(rp.Status400*100)/l)
	fmt.Fprintf(&buf, "50X Responses: %8d (%6.2f%%)\n", rp.Status500, float64(rp.Status500*100)/l)
	fmt.Fprintf(&buf, "Err Responses: %8d (%6.2f%%)\n", rp.StatusErr, float64(rp.StatusErr*100)/l)
	return buf.String()
}
