package firestore

import (
	"bytes"
	"context"
	"path"
	"strings"

	"cloud.google.com/go/firestore/apiv1/firestorepb"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/stats"
)

type firestoreAttrGetter struct{}

var _ stats.Handler = firestoreAttrGetter{}

func (firestoreAttrGetter) HandleRPC(ctx context.Context, sr stats.RPCStats) {
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return
	}

	switch sr := sr.(type) {
	case *stats.InPayload:
		if resp, ok := sr.Payload.(*firestorepb.RunQueryResponse); ok {
			addResponseAttrs(span, resp)
		}
	case *stats.OutPayload:
		if req, ok := sr.Payload.(*firestorepb.CommitRequest); ok {
			span.SetAttributes(attribute.Int("db.operation.batch.size", len(req.Writes)))
		}
		if query, ok := getQuery(sr.Payload); ok {
			span.SetAttributes(attribute.String("db.query.text", strings.TrimSpace(query)))
		}
	}
}

func (firestoreAttrGetter) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	return ctx
}

func (firestoreAttrGetter) TagConn(ctx context.Context, info *stats.ConnTagInfo) context.Context {
	return ctx
}

func (firestoreAttrGetter) HandleConn(ctx context.Context, sr stats.ConnStats) {}

func getQuery(payload any) (string, bool) {
	switch payload := payload.(type) {
	case *firestorepb.RunAggregationQueryRequest:
		if query, ok := payload.QueryType.(*firestorepb.RunAggregationQueryRequest_StructuredAggregationQuery); ok {
			return query.StructuredAggregationQuery.String(), true
		}
	case *firestorepb.RunQueryRequest:
		if query, ok := payload.QueryType.(*firestorepb.RunQueryRequest_StructuredQuery); ok {
			return query.StructuredQuery.String(), true
		}
	case *firestorepb.CommitRequest:
		buf := new(bytes.Buffer)
		var seen bool
		for _, w := range payload.Writes {
			if seen {
				buf.WriteRune(';')
			}
			seen = true
			buf.Write([]byte(w.String()))
		}
		return buf.String(), true
	}
	return "", false
}

func addResponseAttrs(span trace.Span, resp *firestorepb.RunQueryResponse) {
	attrs := make([]attribute.KeyValue, 1, 2)
	name := resp.GetDocument().GetName()
	attrs[0] = attribute.String("db.collection.name", path.Dir(name))
	if ns, _, ok := strings.Cut(name, "/documents/"); ok {
		attrs = append(attrs, attribute.String("db.namespace", ns))
	}
	span.SetAttributes(attrs...)
}
