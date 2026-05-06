package routes

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"path"
	"reflect"
	"strconv"
	"strings"

	"github.com/bragemusic/core/pkg/auth"
	"github.com/bragemusic/core/pkg/bragerr"
	"github.com/bragemusic/core/pkg/types"
	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid/v5"
	"github.com/swaggest/openapi-go"
	"github.com/swaggest/openapi-go/openapi31"
)

var uuidType = reflect.TypeOf(uuid.UUID{})

type RouteFunc[Req Validator, Resp any] func(ctx context.Context, req Req, user types.UserDetails, w http.ResponseWriter, r *http.Request) (resp types.Response[Resp], err error)

type Validator interface {
	Validate() (validationMessages string, err error)
}

type RouteHandler interface {
	Method() string
	Path() string
	Roles() []types.UserRole
	Handler(log, errLog *slog.Logger, berr *bragerr.BragErrFactory) http.Handler
	Docs(refl *openapi31.Reflector, basePath string) error
	AddMiddleware(RouteMiddleware) RouteHandler
	Middlewares() []RouteMiddleware
}

type RouteMeta struct {
	Summary             string
	Description         string
	ExpectedDescription string
	Tags                []string
	Errors              []RouteErrorMeta
	ExpectedStatus      int
	Deprecated          bool
}

type RouteErrorMeta struct {
	Description string
	Status      int
}

type RouteMiddleware struct {
	Name        string
	Description string
	Func        func(next http.Handler) http.Handler
}

type RouteObject[Req Validator, Resp any] struct {
	handlerFunc RouteFunc[Req, Resp]
	method      string
	path        string
	roles       []types.UserRole
	meta        RouteMeta
	middlewares []RouteMiddleware
}

func (ro RouteObject[Req, Resp]) AddMiddleware(rm RouteMiddleware) RouteHandler {
	ro.middlewares = append(ro.middlewares, rm)
	return ro
}

func (ro RouteObject[Req, Resp]) Middlewares() []RouteMiddleware {
	return ro.middlewares
}

func (ro RouteObject[Req, Resp]) Roles() []types.UserRole {
	return ro.roles
}

func (ro RouteObject[Req, Resp]) Handler(log, errLog *slog.Logger, berr *bragerr.BragErrFactory) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var jsonResp bool

		var zero Resp
		switch any(zero).(type) {
		case types.NoResponse:
			jsonResp = false
		default:
			jsonResp = true
		}

		if jsonResp {
			w.Header().Set("Content-Type", "application/json")
		}

		var reqStruct Req

		err := ParsePaths(ctx, berr, &reqStruct)
		if err != nil {
			bragerr.HandleHttpResponse(ctx, err, w, errLog)
			return
		}

		err = ParseQueries(r.URL.Query(), &reqStruct)
		if err != nil {
			bragerr.HandleHttpResponse(ctx, err, w, errLog)
			return
		}

		ct := r.Header.Get("Content-Type")
		if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
			if !strings.Contains(ct, "multipart/form-data") {
				if err = json.NewDecoder(r.Body).Decode(&reqStruct); err != nil {
					if err != io.EOF {
						bragerr.HandleHttpResponse(ctx, err, w, errLog)
						return
					}
				}
			}
		}

		vErrs, err := reqStruct.Validate()
		if err != nil {
			bragerr.HandleHttpResponse(ctx, err, w, errLog)
			return
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, vErrs)
			slog.ErrorContext(ctx, "validation of incoming json data failed")
			return
		}

		user, err := auth.UserFromContext(ctx)
		if err != nil {
			bragerr.HandleHttpResponse(ctx, err, w, errLog)
			return
		}

		resp, err := ro.handlerFunc(ctx, reqStruct, user, w, r)
		if err != nil {
			bragerr.HandleHttpResponse(ctx, err, w, errLog)
			return
		}

		w.WriteHeader(resp.Status)
		if jsonResp {
			if err := json.NewEncoder(w).Encode(resp.Payload); err != nil {
				bragerr.HandleHttpResponse(ctx, err, w, errLog)
				return
			}
		}
	})
}

func (ro RouteObject[Req, Resp]) Docs(refl *openapi31.Reflector, basePath string) error {
	op, err := refl.NewOperationContext(ro.Method(), path.Join(basePath, ro.Path()))
	if err != nil {
		return err
	}
	op.SetTags(ro.meta.Tags...)
	op.SetSummary(ro.meta.Summary)

	if len(ro.roles) > 0 {
		var rFormatted strings.Builder
		for _, r := range ro.roles {
			fmt.Fprintf(&rFormatted, "\n - `%s`", r)
		}
		desc := fmt.Sprintf("%s<br><br>Requires at least one of the following user roles:%s", ro.meta.Description, rFormatted.String())
		op.SetDescription(desc)
	} else {
		op.SetDescription(ro.meta.Description)
	}

	if ro.meta.Deprecated {
		op.SetIsDeprecated(true)
	}

	op.AddReqStructure(new(Req))
	op.AddRespStructure(new(Resp), func(cu *openapi.ContentUnit) {
		cu.HTTPStatus = ro.meta.ExpectedStatus
		cu.Description = ro.meta.ExpectedDescription
	})

	for _, e := range ro.meta.Errors {
		_ = e
		op.AddRespStructure(new(bragerr.BragErr), func(cu *openapi.ContentUnit) {
			cu.HTTPStatus = e.Status
			cu.Description = e.Description
		})
	}

	op.AddRespStructure(new(bragerr.BragErr), func(cu *openapi.ContentUnit) { cu.HTTPStatus = http.StatusInternalServerError })

	return refl.AddOperation(op)
}

func (ro RouteObject[Req, Resp]) Method() string {
	return ro.method
}

func (ro RouteObject[Req, Resp]) Path() string {
	return ro.path
}

func New[Req Validator, T any](method, path string, f RouteFunc[Req, T], r []types.UserRole, m RouteMeta) RouteObject[Req, T] {
	return RouteObject[Req, T]{
		handlerFunc: f,
		method:      method,
		path:        path,
		roles:       r,
		meta:        m,
	}
}

func ParsePaths[T Validator](ctx context.Context, berr *bragerr.BragErrFactory, v *T) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return fmt.Errorf("must be a pointer, not %s", rv.Type())
	}
	rv = reflect.Indirect(rv)

	rctx := chi.RouteContext(ctx)
	if rctx == nil {
		return errors.New("not a route context")
	}

	tags := map[string]int{}

	for i := 0; i < rv.NumField(); i++ {
		f := rv.Type().Field(i)
		if f.Tag.Get("path") != "" {
			tags[f.Tag.Get("path")] = i
		}
	}

	for i := 0; i < len(rctx.URLParams.Keys); i++ {
		k := rctx.URLParams.Keys[i]
		v := rctx.URLParams.Values[i]
		if k == "*" {
			continue
		}

		fidx, ok := tags[k]
		if !ok {
			return fmt.Errorf("'%s' does not have a '%s' path tag", rv.Type(), k)
		}

		field := rv.Field(fidx)
		fieldType := field.Type()

		switch rv.Field(fidx).Kind() {
		case reflect.String:
			rv.Field(fidx).SetString(v)
		case reflect.Int:
			intVal, err := strconv.Atoi(v)
			if err != nil {
				return berr.ParamWrongFormat(err, k, "int")
			}
			rv.Field(fidx).SetInt(int64(intVal))
		default:
			if fieldType == uuidType {
				parsed, err := uuid.FromString(v)
				if err != nil {
					return berr.ParamWrongFormat(err, k, "uuid")
				}
				field.Set(reflect.ValueOf(parsed))
			} else {
				return fmt.Errorf("cannot parse path '%s' of type '%s'. No support.", k, rv.Field(fidx).Kind())
			}
		}
	}

	return nil
}

func ParseQueries[T Validator](q url.Values, v *T) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return fmt.Errorf("must be a pointer, not %s", rv.Type())
	}
	rv = reflect.Indirect(rv)

	tags := map[string]int{}

	for i := 0; i < rv.NumField(); i++ {
		f := rv.Type().Field(i)
		if f.Tag.Get("query") != "" {
			tags[f.Tag.Get("query")] = i
		}
	}

	for k, vs := range q {

		fidx, ok := tags[k]
		if !ok {
			return fmt.Errorf("'%s' does not have a '%s' query tag", rv.Type(), k)
		}

		field := rv.Field(fidx)
		fieldType := field.Type()

		switch rv.Field(fidx).Kind() {
		case reflect.String:
			if len(vs) > 0 {
				rv.Field(fidx).SetString(vs[0])
			}
		case reflect.Int:
			if len(vs) > 0 {
				intVal, err := strconv.Atoi(vs[0])
				if err != nil {
					return err
				}
				rv.Field(fidx).SetInt(int64(intVal))
			}

		case reflect.Bool:
			if len(vs) > 0 {
				rv.Field(fidx).SetBool(vs[0] == "true")
			}

		default:
			if fieldType == uuidType {
				if len(vs) > 0 {
					parsed, err := uuid.FromString(vs[0])
					if err != nil {
						return err
					}
					field.Set(reflect.ValueOf(parsed))
					return nil
				}
			}
			return fmt.Errorf("cannot parse path '%s' of type '%s'. No support.", k, rv.Field(fidx).Kind())
		}
	}

	return nil
}
