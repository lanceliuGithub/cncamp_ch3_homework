package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

type MyConfig struct {
	Server  Server `json:"server"`
	Log     Log    `json:"log"`
}
type Server struct {
	Host string `json:"host"`
	Port string `json:"port"`
}
type Log struct {
	Enable              bool `json:"enable"`
	EnableRequestHeader bool `json:"request_header"`
}

var myConf *MyConfig

func main() {
	// Load Config
	myConf = new(MyConfig)
	loadConfig("config.json")

	// Register http handlers
	rootHandler := wrapHandlerWithLogging(http.HandlerFunc(handleRoot))
	http.Handle("/", rootHandler)

	healthzHandler := wrapHandlerWithLogging(http.HandlerFunc(handleHealthz))
	http.Handle("/healthz", healthzHandler)

	// Start HTTP server
	go func() {
		printOSEnvVersion()

		serverHost := myConf.Server.Host
		serverPort := myConf.Server.Port
		log.Print("HTTP Server starting, listening on " + serverHost + ":" + serverPort)
		if err := http.ListenAndServe(serverHost + ":" + serverPort, nil); err != http.ErrServerClosed {
			log.Fatalf("HTTP server crashed: %v", err)
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(
		signalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
	)

	<-signalChan
	log.Print("OS Interrupt -- HTTP server shutting down")

	go func() {
		<-signalChan
		log.Fatal("OS kill -- Termination")
	}()

	defer os.Exit(0)
	return
}

func loadConfig(confFilename string) {
	ExecPath, _ := os.Executable()
	confFilepath := filepath.Dir(ExecPath) + "/" + confFilename
	configFile, err := os.Open(confFilepath)
	if err != nil {
		log.Fatalf("Config file <%s> NOT found: %v", confFilename, err)
	}
	defer configFile.Close()

	configBytes, err := ioutil.ReadAll(configFile)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(configBytes, &myConf);
	if err != nil {
		panic(err)
	}
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	printRequestHeaders(w, r)
}

// 当访问 /healthz 时，状态码200，返回success字样
func handleHealthz(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "success\n")
	printRequestHeaders(w, r)
}

// 读取当前系统的环境变量中的 VERSION 配置
func printOSEnvVersion() {
	var v string
	if v = os.Getenv("VERSION"); v == "" {
		v = "Unknown"
	}
	log.Printf("MyHTTPServer Version: %s\n", v)
}

// 接收客户端 request，并将 request 中带的 header 写入 response header
func printRequestHeaders(w http.ResponseWriter, r *http.Request) {
	logEnabled := myConf.Log.Enable
	logRequestHeaderEnabled := myConf.Log.EnableRequestHeader
	//log.Printf("logEnabled: " + strconv.FormatBool(logEnabled))

	if ! logEnabled || ! logRequestHeaderEnabled {
		return
	}

	log.Printf("Request Headers (Total of %d):\n", len(r.Header))
	for k, v := range r.Header {
		log.Printf("==> %s: %s\n", k, v)
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
		logEnabled := myConf.Log.Enable

		if logEnabled {
			log.Printf("--> %s %s", req.Method, req.URL.Path)
		}

		lrw := newLoggingResponseWriter(w)
		wrappedHandler.ServeHTTP(lrw, req)

		if logEnabled {
			statusCode := lrw.statusCode
			log.Printf("<-- %d %s", statusCode, http.StatusText(statusCode))
		}
	})
}
