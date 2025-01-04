package main

//
//import (
//	"bytes"
//	"encoding/json"
//	"log"
//	"math/rand"
//	"net/http"
//	"time"
//
//	"github.com/prometheus/client_golang/prometheus"
//	"github.com/prometheus/client_golang/prometheus/collectors"
//	"github.com/prometheus/client_golang/prometheus/promhttp"
//)
//
//const (
//	testAddr = "localhost:9123"
//)
//
//var httpStatusCodeCounter = prometheus.NewCounterVec(
//	prometheus.CounterOpts{
//		Name: "http_status_code_counter",
//		Help: "Count http status code",
//	},
//	[]string{"status_code"},
//)
//
//func main() {
//	go produceData()
//	reg := prometheus.NewRegistry()
//	prometheus.WrapRegistererWith(prometheus.Labels{"serviceName": "demo-service"}, reg).MustRegister(
//		collectors.NewGoCollector(),
//		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
//		httpStatusCodeCounter,
//	)
//	// localhost:9123/metrics
//	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
//	http.HandleFunc("/", sendMetricsHandler)
//	log.Fatal(http.ListenAndServe(testAddr, nil))
//}
//
//func sendMetricsHandler(w http.ResponseWriter, r *http.Request) {
//	var req request
//	defer func() {
//		httpStatusCodeCounter.WithLabelValues(req.StatusCode).Inc()
//		log.Printf("add 1 to %s", req.StatusCode)
//	}()
//	_ = json.NewDecoder(r.Body).Decode(&req)
//	log.Printf("receive req:%+v", req)
//	_, _ = w.Write([]byte(req.StatusCode))
//}
//
//type request struct {
//	StatusCode string
//}
//
//func produceData() {
//	codes := []string{"503", "404", "400", "200", "304", "500"}
//	for {
//		body, _ := json.Marshal(request{
//			StatusCode: codes[rand.Intn(len(codes))],
//		})
//		requestBody := bytes.NewBuffer(body)
//		http.Post("http://"+testAddr, "application/json", requestBody)
//		log.Printf("send request=%s to %s", requestBody.String(), testAddr)
//		time.Sleep(2 * time.Second)
//	}
//}
