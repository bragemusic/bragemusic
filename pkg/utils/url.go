package utils

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/bragemusic/core/pkg/bragerr"
	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid/v5"
)

type paramFace interface {
	~string | ~float32 | ~int | uuid.UUID
}

func GetURLParameter[T paramFace](ctx context.Context, key string) (T, error) {
	var zero T

	berr := bragerr.NewFactory("server")

	val := chi.URLParamFromCtx(ctx, key)
	if val == "" {
		return zero, errors.New("could not parse parameter")
	}

	var t T
	switch any(t).(type) {

	case string:
		return any(val).(T), nil

	case int:
		parsed, e := strconv.Atoi(val)
		if e != nil {
			return zero, berr.ParamWrongFormat(e, key, "int")
		}
		return any(parsed).(T), nil

	case float32:
		parsed, e := strconv.ParseFloat(val, 32)
		if e != nil {
			return zero, berr.ParamWrongFormat(e, key, "float")
		}
		return any(float32(parsed)).(T), nil

	case uuid.UUID:
		parsed, e := uuid.FromString(val)
		if e != nil {
			return zero, berr.ParamWrongFormat(e, key, "uuid")
		}
		return any(parsed).(T), nil

	default:
		return zero, fmt.Errorf("unsupported parameter type")
	}
}
