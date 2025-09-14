package server

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"
)

// https://httpd.apache.org/docs/2.4/logs.html
// LogFormat "%h %l %u %t \"%r\" %>s %b" common
// 127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif HTTP/1.0" 200 2326
func TestCommonLogFormat(t *testing.T) {
	var req http.Request
	req.URL = &url.URL{}

	req.RemoteAddr = "127.0.0.1"
	req.Method = "GET"
	req.URL = &url.URL{}
	req.URL.Path = "/apache_pb.gif"
	req.Proto = "HTTP/1.0"
	req.ContentLength = 2326

	currentTime := time.Now()
	statusCode := 200

	currentTime.Format("2006-01-02T15:04:05Z07:00")

	want := fmt.Sprintf("127.0.0.1 - - [%s] \"GET /apache_pb.gif HTTP/1.0\" 200 2326", currentTime.Format("2006-01-02T15:04:05Z07:00"))

	logMsg := formatCommonLog(req, currentTime, statusCode)

	if logMsg != want {
		t.Errorf(`formatCommonLog(req, currentTime, statusCode) = %q, want match for %#q, nil`, logMsg, want)
	}
}

func TestCommonLogFormatWithQuery(t *testing.T) {
	var req http.Request
	req.URL = &url.URL{}

	req.RemoteAddr = "127.0.0.1"
	req.Method = "GET"
	req.URL = &url.URL{}
	req.URL.Path = "/apache_pb.gif"
	req.URL.RawQuery = "demo=test"
	req.Proto = "HTTP/1.0"
	req.ContentLength = 2400

	currentTime := time.Now()
	statusCode := 200

	currentTime.Format("2006-01-02T15:04:05Z07:00")

	want := fmt.Sprintf("127.0.0.1 - - [%s] \"GET /apache_pb.gif?demo=test HTTP/1.0\" 200 2400", currentTime.Format("2006-01-02T15:04:05Z07:00"))

	logMsg := formatCommonLog(req, currentTime, statusCode)

	if logMsg != want {
		t.Errorf(`formatCommonLog(req, currentTime, statusCode) = %q, want match for %#q, nil`, logMsg, want)
	}
}
