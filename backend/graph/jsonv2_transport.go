package graph

import (
	"bufio"
	"context"
	jsonv1 "encoding/json"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"strings"
	"sync"

	"careco/backend/log/attribute"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/errcode"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// POSTTransport returns an [graphql.Transport] that expects HTTP POST request
// but using [encoding/json/v2] with strict options to marshal/unmarshal JSON.
func POSTTransport() graphql.Transport {
	return postJSONV2Transport{}
}

var graphqlParamsPool = sync.Pool{New: func() any { return &graphql.RawParams{} }}

type postJSONV2Transport struct{}

var _ graphql.Transport = postJSONV2Transport{}

func (postJSONV2Transport) Supports(r *http.Request) bool {
	if r.Header.Get("Upgrade") != "" {
		return false
	}
	mt, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		return false
	}
	return r.Method == http.MethodPost && mt == "application/json"
}

func (postJSONV2Transport) Do(w http.ResponseWriter, r *http.Request, exec graphql.GraphExecutor) {
	respType := determineResponseType(r)
	w.Header().Set("Content-Type", respType.mediaType())

	params := graphqlParamsPool.Get().(*graphql.RawParams) //nolint:forcetypeassert,errcheck
	defer func() {
		params.Headers = nil
		params.ReadTime = graphql.TraceTiming{}
		params.Extensions = nil
		params.OperationName = ""
		params.Query = ""
		params.Variables = nil

		graphqlParamsPool.Put(params)
	}()

	start := graphql.Now()

	bufReader := bufio.NewReader(r.Body)
	if err := json.UnmarshalRead(bufReader, params, unmarshalOptions); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		gqlerr := gqlerror.Errorf("JSON request body could no be decoded: %+v", err)
		resp := exec.DispatchError(r.Context(), gqlerror.List{gqlerr})
		emitJSON(r.Context(), w, resp)
		return
	}

	params.Headers = r.Header
	params.ReadTime = graphql.TraceTiming{
		Start: start,
		End:   graphql.Now(),
	}

	rc, opErr := exec.CreateOperationContext(r.Context(), params)
	if opErr != nil {
		w.WriteHeader(determineStatusCode(respType, opErr))
		resp := exec.DispatchError(graphql.WithOperationContext(r.Context(), rc), opErr)
		emitJSON(r.Context(), w, resp)
		return
	}

	respHandler, ctx := exec.DispatchOperation(r.Context(), rc)
	emitJSON(ctx, w, respHandler(ctx))
}

func emitJSON(ctx context.Context, w io.Writer, resp *graphql.Response) {
	if err := json.MarshalWrite(w, resp, json.Deterministic(true)); err != nil {
		slog.ErrorContext(ctx, "failed to emit JSON", attribute.Error(err))
		panic("failed to emit JSON")
	}
}

func determineStatusCode(rt responseType, err gqlerror.List) int {
	switch errcode.GetErrorKind(err) {
	case errcode.KindProtocol:
		if rt == responseTypeGraphql {
			return http.StatusBadRequest
		} else {
			return http.StatusUnprocessableEntity
		}
	default:
		return http.StatusOK
	}
}

// unmarshalOptions decodes JSON numbers into [any] as [jsonv1.Number] to be compatible with
// [github.com/99designs/gqlgen/graphql/handler/transport.POST] that uses [jsonv1.Decoder.UseNumber].
var unmarshalOptions = json.JoinOptions(
	json.DefaultOptionsV2(),
	json.WithUnmarshalers(json.UnmarshalFromFunc(unmarshalNumberAsJSONNumber)),
)

func unmarshalNumberAsJSONNumber(dec *jsontext.Decoder, v *any) error {
	if dec.PeekKind() != jsontext.KindNumber {
		return errors.ErrUnsupported
	}
	val, err := dec.ReadValue()
	if err != nil {
		return err
	}
	*v = jsonv1.Number(val)
	return nil
}

func determineResponseType(r *http.Request) responseType {
	accept := r.Header.Get("Accept")
	if accept == "" {
		return responseTypeGraphql
	}

	for spec := range strings.SplitSeq(accept, ",") {
		mt, _, err := mime.ParseMediaType(strings.TrimSpace(spec))
		if err != nil {
			continue
		}

		switch mt {
		case "*/*", "application/*", "application/graphql-response+json":
			return responseTypeGraphql
		case "application/json":
			return responseTypeJson
		}
	}

	return responseTypeJson
}

type responseType int

func (rt responseType) mediaType() string {
	switch rt {
	case responseTypeJson:
		return "application/json"
	case responseTypeGraphql:
		return "application/graphql-response+json"
	default:
		return "application/octet-stream"
	}
}

const (
	responseTypeJson responseType = iota
	responseTypeGraphql
)
