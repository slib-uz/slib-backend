package middlewares

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"runtime"
	"strings"
	"time"

	"slib.uz/src/core/application/response"
	"slib.uz/src/infrastructure/logger"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

const maxBodyBytes = 100 * 1024 * 1024 // 100MB for base64 file uploads

type RecoveryAndLoggerMiddleware struct {
	logger *logger.HttpLogger
}

// @inject
func NewRecoveryAndLoggerMiddleware(logger *logger.HttpLogger) *RecoveryAndLoggerMiddleware {
	return &RecoveryAndLoggerMiddleware{logger: logger}
}

func (this *RecoveryAndLoggerMiddleware) Wrap(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		lg := this.initLogger(c, c.Request())

		defer this.recoverPanic(c, lg, c.Request())

		req := c.Request()
		if req.Method != http.MethodPost && req.Method != http.MethodPut && req.Method != http.MethodPatch && req.Method != http.MethodDelete {
			return next(c)
		}

		return this.handleWithLog(c, next, lg, req)
	}
}

func (this *RecoveryAndLoggerMiddleware) initLogger(c echo.Context, req *http.Request) *logger.HttpLogger {
	reqID := req.Header.Get("X-Request-ID")
	if reqID == "" {
		reqID = uuid.NewString()
	}

	zapLg := this.logger.With(zap.String("request_id", reqID))

	ctx := logger.WithLogger(req.Context(), zapLg)
	req = req.WithContext(ctx)
	c.SetRequest(req)

	return &logger.HttpLogger{Logger: zapLg}
}

func (this *RecoveryAndLoggerMiddleware) recoverPanic(c echo.Context, lg *logger.HttpLogger, req *http.Request) {
	r := recover()
	if r == nil {
		return
	}

	stack := make([]byte, 4<<10)
	length := runtime.Stack(stack, false)

	lg.Error("panic recovered",
		zap.String("method", req.Method),
		zap.String("path", req.URL.Path),
		zap.String("ip", c.RealIP()),
		zap.Any("panic", r),
		zap.ByteString("stack", stack[:length]),
	)

	c.Error(echo.NewHTTPError(http.StatusInternalServerError, "internal server error"))
}

func (this *RecoveryAndLoggerMiddleware) handleWithLog(c echo.Context, next echo.HandlerFunc, lg *logger.HttpLogger, req *http.Request) error {
	reqBody := readRequestBody(req)

	rec := &responseBodyRecorder{
		ResponseWriter: c.Response().Writer,
		limit:          maxBodyBytes,
	}
	c.Response().Writer = rec

	start := time.Now()
	err := next(c)
	latency := time.Since(start)

	this.logResponse(c, lg, req, reqBody, rec, latency, err)

	if err != nil {
		return err
	}
	return nil
}

func (this *RecoveryAndLoggerMiddleware) logResponse(
	c echo.Context,
	lg *logger.HttpLogger,
	req *http.Request,
	reqBody any,
	rec *responseBodyRecorder,
	latency time.Duration,
	err error,
) {
	status := c.Response().Status

	if err != nil {
		var he *echo.HTTPError
		var resp *response.Response
		if errors.As(err, &he) {
			status = he.Code
		} else if errors.As(err, &resp) {
			status = resp.Status
		} else {
			status = http.StatusInternalServerError
		}
	}

	resBody := parseJSON(c.Response().Header().Get(echo.HeaderContentType), rec.body.Bytes())

	fields := []zap.Field{
		zap.String("method", req.Method),
		zap.String("path", req.URL.Path),
		zap.String("query", req.URL.RawQuery),
		zap.String("ip", c.RealIP()),
		zap.String("user_agent", req.UserAgent()),
		zap.Any("scope", "http"),
		zap.Int("status", status),
		zap.Duration("latency", latency),
		zap.Any("req_body", reqBody),
		zap.Any("res_body", resBody),
		zap.Any("req_header", req.Header),
	}

	if err != nil {
		lg.Error("http error", append(fields, zap.Error(err))...)
		return
	}

	switch {
	case status >= 500:
		lg.Error("http", fields...)
	case status >= 400:
		lg.Warn("http", fields...)
	default:
		lg.Info("http", fields...)
	}
}

func readRequestBody(req *http.Request) any {
	if req.Body == nil || req.Body == http.NoBody {
		return nil
	}
	b, _ := io.ReadAll(io.LimitReader(req.Body, maxBodyBytes))
	body := parseJSON(req.Header.Get(echo.HeaderContentType), b)
	req.Body = io.NopCloser(bytes.NewReader(b))
	return body
}

type responseBodyRecorder struct {
	http.ResponseWriter
	body  bytes.Buffer
	limit int64
}

func (r *responseBodyRecorder) Write(b []byte) (int, error) {
	if int64(r.body.Len()) < r.limit {
		remain := r.limit - int64(r.body.Len())
		if int64(len(b)) > remain {
			r.body.Write(b[:remain])
		} else {
			r.body.Write(b)
		}
	}
	return r.ResponseWriter.Write(b)
}

func (r *responseBodyRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := r.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, http.ErrNotSupported
	}
	return h.Hijack()
}

func parseJSON(contentType string, raw []byte) any {
	if len(raw) == 0 {
		return nil
	}
	if !strings.Contains(strings.ToLower(contentType), "application/json") {
		return truncate(string(raw), 2000)
	}
	var v any
	if err := json.Unmarshal(raw, &v); err != nil {
		return truncate(string(raw), 2000)
	}
	return v
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "...(truncated)"
}
