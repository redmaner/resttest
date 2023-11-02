package resttest

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

func ExecuteTests(logger *zap.Logger, t http.RoundTripper, baseUrl string, tests []HttpTest) {

	s := semaphore.NewWeighted(2)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	group, gCtx := errgroup.WithContext(ctx)

	for i := range tests {
		httpTest := tests[i]
		group.Go(func() error {
			runTest(gCtx, logger, t, baseUrl, s, &httpTest)
			return nil
		})
	}

	_ = group.Wait()
}

func runTest(ctx context.Context, logger *zap.Logger, t http.RoundTripper, baseUrl string, s *semaphore.Weighted, httpTest *HttpTest) {

	if err := s.Acquire(ctx, 1); err != nil {
		return
	}
	defer s.Release(1)

	var bodyBuffer io.Reader
	if len(httpTest.Body) > 0 {
		bodyBuffer = bytes.NewBufferString(httpTest.Body)
	}

	requestUrl := baseUrl + httpTest.Path
	httpRequest, err := http.NewRequest(httpTest.Method, requestUrl, bodyBuffer)
	if err != nil {
		logger.Error("encountered an error creating request", zap.Error(err), zap.String("url", requestUrl), zap.String("method", httpTest.Method))
		return
	}

	resp, err := t.RoundTrip(httpRequest)
	if err != nil {
		logger.Error("encountered an error executing request", zap.Error(err), zap.String("url", requestUrl), zap.String("method", httpTest.Method))
		return
	}
	defer resp.Body.Close()

	var testFailed bool

	if resp.StatusCode != httpTest.Expect.StatusCode {
		logger.Error("Expectation failed", zap.String("url", requestUrl), zap.String("method", httpTest.Method), zap.Int("expected_status_code", httpTest.Expect.StatusCode), zap.Int("response_status_code", resp.StatusCode))
		testFailed = true
	}

	for i := range httpTest.Expect.Headers {
		key := httpTest.Expect.Headers[i].Key
		value := httpTest.Expect.Headers[i].Value

		if respValue := resp.Header.Get(key); respValue != value {
			logger.Error("Expectation failed", zap.String("url", requestUrl), zap.String("method", httpTest.Method), zap.String("header", key), zap.String("expected_value", value), zap.String("response_value", respValue))
			testFailed = true
		}
	}

	if !testFailed {
		logger.Info("expectation passed", zap.String("url", requestUrl), zap.String("method", httpTest.Method))
	}
}
