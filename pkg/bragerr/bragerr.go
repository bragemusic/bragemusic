package bragerr

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

type Response struct {
	Error *BragErr `json:"error"`
}

type BragErr struct {
	Code    string `json:"code"`
	Title   string `json:"title"`
	Message string `json:"message"`

	Status   int         `json:"-"`
	Service  string      `json:"-"`
	Err      error       `json:"-"`
	LogAttrs []slog.Attr `json:"-"`
}

func (e *BragErr) Error() string {
	if e.Err != nil {
		return e.Code + ": " + e.Err.Error()
	}
	return e.Code
}

func (e *BragErr) Unwrap() error {
	return e.Err
}

func (e *BragErr) With(key string, val any) *BragErr {
	e.LogAttrs = append(e.LogAttrs, slog.Attr{
		Key:   key,
		Value: slog.AnyValue(val),
	})
	return e
}

type BragErrFactory struct {
	service string
}

func NewFactory(service string) BragErrFactory {
	return BragErrFactory{
		service: service,
	}
}

func HandleHttpResponse(ctx context.Context, err error, w http.ResponseWriter, logger *slog.Logger) {
	if err == nil {
		return
	}

	berr, ok := err.(*BragErr)
	if !ok {
		berr = &BragErr{
			Code:    "UNKNOWN",
			Title:   "Internal Server Error",
			Message: "An unknown error occured",
			Status:  http.StatusInternalServerError,
		}
		logger.ErrorContext(ctx, "unknown sever error", "error", err)
	} else {
		attrs := []slog.Attr{
			slog.Any("code", berr.Code),
			slog.Any("service", berr.Service),
		}
		if berr.Err != nil {
			attrs = append(attrs, slog.Any("error", berr.Err.Error()))
		}
		attrs = append(attrs, berr.LogAttrs...)
		logger.LogAttrs(ctx, slog.LevelError, berr.Message, attrs...)
	}

	w.WriteHeader(berr.Status)
	w.Header().Add("Content-Type", "application/json")

	if ierr := json.NewEncoder(w).Encode(Response{Error: berr}); ierr != nil {
		logger.ErrorContext(ctx, "error when encoding error", "error", ierr)
	}
}
