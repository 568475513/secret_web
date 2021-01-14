/*
Package http implements a HTTP reporter to send spans to Zipkin V2 collectors.
*/
package zipkin

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/openzipkin/zipkin-go/model"
	"github.com/openzipkin/zipkin-go/reporter"
	"go.uber.org/zap"
)

// defaults
const (
	defaultBatchInterval = time.Second * 1 // BatchInterval in seconds
	defaultBatchSize     = 100
	defaultMaxBacklog    = 1000
)

// LogReporter will send spans to a Zipkin.
type logReporter struct {
	log           *zap.Logger
	logger        *log.Logger
	batchInterval time.Duration
	batchSize     int
	maxBacklog    int
	sendMtx       *sync.Mutex
	batchMtx      *sync.Mutex
	batch         []*model.SpanModel
	spanC         chan *model.SpanModel
	quit          chan struct{}
	shutdown      chan error
	reqCallback   RequestCallbackFn
}

// Send implements reporter
func (r *logReporter) Send(s model.SpanModel) {
	r.spanC <- &s
}

// Close implements reporter
func (r *logReporter) Close() error {
	close(r.quit)
	return <-r.shutdown
}

func (r *logReporter) loop() {
	var (
		nextSend   = time.Now().Add(r.batchInterval)
		ticker     = time.NewTicker(r.batchInterval / 10)
		tickerChan = ticker.C
	)
	defer ticker.Stop()

	for {
		select {
		case span := <-r.spanC:
			currentBatchSize := r.append(span)
			if currentBatchSize >= r.batchSize {
				nextSend = time.Now().Add(r.batchInterval)
				go func() {
					_ = r.sendBatch()
				}()
			}
		case <-tickerChan:
			if time.Now().After(nextSend) {
				nextSend = time.Now().Add(r.batchInterval)
				go func() {
					_ = r.sendBatch()
				}()
			}
		case <-r.quit:
			r.shutdown <- r.sendBatch()
			return
		}
	}
}

func (r *logReporter) append(span *model.SpanModel) (newBatchSize int) {
	r.batchMtx.Lock()

	r.batch = append(r.batch, span)
	if len(r.batch) > r.maxBacklog {
		dispose := len(r.batch) - r.maxBacklog
		r.logger.Printf("backlog too long, disposing %d spans", dispose)
		r.batch = r.batch[dispose:]
	}
	newBatchSize = len(r.batch)

	r.batchMtx.Unlock()
	return
}

func (r *logReporter) sendBatch() error {
	// in order to prevent sending the same batch twice
	r.sendMtx.Lock()
	defer r.sendMtx.Unlock()

	// Select all current spans in the batch to be sent
	r.batchMtx.Lock()
	sendBatch := r.batch[:]
	r.batchMtx.Unlock()

	if len(sendBatch) == 0 {
		return nil
	}

	body, err := json.Marshal(sendBatch)
	if err != nil {
		r.logger.Printf("failed when marshalling the spans batch: %s\n", err.Error())
		return err
	}

	// 直接写日志
	r.log.Info(string(body))

	// Remove sent spans from the batch even if they were not saved
	r.batchMtx.Lock()
	r.batch = r.batch[len(sendBatch):]
	r.batchMtx.Unlock()

	return nil
}

// RequestCallbackFn receives the initialized request from the Collector before
// sending it over the wire. This allows one to plug in additional headers or
// do other customization.
type RequestCallbackFn func(*http.Request)

// ReporterOption sets a parameter for the HTTP Reporter
type ReporterOption func(r *logReporter)

// BatchSize sets the maximum batch size, after which a collect will be
// triggered. The default batch size is 100 traces.
func BatchSize(n int) ReporterOption {
	return func(r *logReporter) { r.batchSize = n }
}

// MaxBacklog sets the maximum backlog size. When batch size reaches this
// threshold, spans from the beginning of the batch will be disposed.
func MaxBacklog(n int) ReporterOption {
	return func(r *logReporter) { r.maxBacklog = n }
}

// BatchInterval sets the maximum duration we will buffer traces before
// emitting them to the collector. The default batch interval is 1 second.
func BatchInterval(d time.Duration) ReporterOption {
	return func(r *logReporter) { r.batchInterval = d }
}

// RequestCallback registers a callback function to adjust the reporter
// *http.Request before it sends the request to Zipkin.
func RequestCallback(rc RequestCallbackFn) ReporterOption {
	return func(r *logReporter) { r.reqCallback = rc }
}

// NewReporter returns a new Log Reporter.
// url should be the endpoint to send the spans to, e.g.
func NewReporter(logObj *zap.Logger, opts ...ReporterOption) reporter.Reporter {
	r := logReporter{
		log:           logObj,
		logger:        log.New(os.Stderr, "", log.LstdFlags),
		batchInterval: defaultBatchInterval,
		batchSize:     defaultBatchSize,
		maxBacklog:    defaultMaxBacklog,
		batch:         []*model.SpanModel{},
		spanC:         make(chan *model.SpanModel),
		quit:          make(chan struct{}, 1),
		shutdown:      make(chan error, 1),
		sendMtx:       &sync.Mutex{},
		batchMtx:      &sync.Mutex{},
	}

	for _, opt := range opts {
		opt(&r)
	}

	go r.loop()

	return &r
}
