package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	rootHandler := wrapHandlerWithLogging(http.HandlerFunc(handleRoot))
	http.Handle("/", rootHandler)

	healthzHandler := wrapHandlerWithLogging(http.HandlerFunc(handleHealthz))
	http.Handle("/healthz", healthzHandler)

	log.Fatal(http.ListenAndServe("127.0.0.1:8888", nil))
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	printOSEnvVersion(w)
	printRequestHeaders(w, r)
}

// 当访问 localhost/healthz 时，应返回 200
func handleHealthz(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "200\n\n")
	printOSEnvVersion(w)
	printRequestHeaders(w, r)
}

// 读取当前系统的环境变量中的 VERSION 配置，并写入 response header
func printOSEnvVersion(w http.ResponseWriter) {
	var v string
	if v = os.Getenv("VERSION"); v == "" {
		v = "Unknown"
	}
	fmt.Fprintf(w, "Version: %s\n\n", v)
}

// 接收客户端 request，并将 request 中带的 header 写入 response header
func printRequestHeaders(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Request Headers (Total of %d):\n", len(r.Header))
	for k, v := range r.Header {
		fmt.Fprintf(w, "==> %s: %s\n", k, v)
	}
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
func wrapHandlerWithLogging(wrappedHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Printf("--> %s %s", req.Method, req.URL.Path)

		lrw := newLoggingResponseWriter(w)
		wrappedHandler.ServeHTTP(lrw, req)

		statusCode := lrw.statusCode
		log.Printf("<-- %d %s", statusCode, http.StatusText(statusCode))
	})
}
