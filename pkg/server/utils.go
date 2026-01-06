package server

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid/v5"
)

type paramFace interface {
	~string | ~float32 | ~int | uuid.UUID
}

func getParameter[T paramFace](ctx context.Context, key string) (T, error) {
	var zero T

	val := chi.URLParamFromCtx(ctx, key)
	if val == "" {
		return zero, ErrParamNotDefined{
			idKey: key,
		}
	}

	var t T
	switch any(t).(type) {

	case string:
		return any(val).(T), nil

	case int:
		parsed, e := strconv.Atoi(val)
		if e != nil {
			return zero, ErrBadParameter{
				idKey: key,
				t:     "int",
				err:   e,
			}
		}
		return any(parsed).(T), nil

	case float32:
		parsed, e := strconv.ParseFloat(val, 32)
		if e != nil {
			return zero, ErrBadParameter{
				idKey: key,
				t:     "float",
				err:   e,
			}
		}
		return any(float32(parsed)).(T), nil

	case uuid.UUID:
		parsed, e := uuid.FromString(val)
		if e != nil {
			return zero, ErrBadParameter{
				idKey: key,
				t:     "uuid",
				err:   e,
			}
		}
		return any(parsed).(T), nil

	default:
		return zero, fmt.Errorf("unsupported parameter type")
	}
}
