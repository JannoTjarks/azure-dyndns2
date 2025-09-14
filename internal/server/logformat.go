package server

import (
	"fmt"
	"net/http"
	"time"
)

// https://httpd.apache.org/docs/2.4/logs.html
// LogFormat "%h %l %u %t \"%r\" %>s %b" common
// 127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif HTTP/1.0" 200 2326

// Ignore "Logging should not be vulnerable to injection attacks"
// The only string, which can be free chosen by the requestor, is the path and query.
// Even then, it will just get printed to stdout. From my understanding the severity is absolut minmal.
func formatCommonLog(req http.Request, currentTime time.Time, statusCode int) string {
	var userId string
	if req.URL.User.Username() != "" {
		userId = req.URL.User.Username()
	} else {
		userId = "-"
	}

	var query string
	if req.URL.RawQuery != "" {
		query = fmt.Sprintf("?%s", req.URL.RawQuery)
	} else {
		query = ""
	}

	return fmt.Sprintf("%s %s %s [%s] \"%s %s%s %s\" %d %d", req.RemoteAddr, "-", userId, currentTime.Format("2006-01-02T15:04:05Z07:00"), req.Method, req.URL.Path, query, req.Proto, statusCode, req.ContentLength)
}
