package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgconn"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type tracedDB struct {
	ContextExecutor
}

func Trace(db ContextExecutor) ContextExecutor {
	return tracedDB{ContextExecutor: db}
}

func (t tracedDB) PrepareContext(ctx context.Context, query string) (stmt *sql.Stmt, err error) {
	span := trace.SpanFromContext(ctx)
	defer func(started time.Time) {
		span.AddEvent("PrepareContext", trace.WithAttributes(
			attribute.String("Query", query),
			attribute.Float64("Took", time.Since(started).Seconds()),
		))
		t.recordError(span, err)
	}(time.Now())

	return t.ContextExecutor.PrepareContext(ctx, query)
}

func (t tracedDB) ExecContext(ctx context.Context, query string, args ...interface{}) (rs sql.Result, err error) {
	span := trace.SpanFromContext(ctx)
	defer func(started time.Time) {
		span.AddEvent("ExecContext", trace.WithAttributes(
			attribute.String("Query", query),
			attribute.Float64("Took", time.Since(started).Seconds()),
		))
		t.recordError(span, err)
	}(time.Now())

	return t.ContextExecutor.ExecContext(ctx, query, args...)
}

func (t tracedDB) QueryContext(ctx context.Context, query string, args ...interface{}) (row *sql.Rows, err error) {
	span := trace.SpanFromContext(ctx)
	defer func(started time.Time) {
		span.AddEvent("QueryContext", trace.WithAttributes(
			attribute.String("Query", query),
			attribute.Float64("Took", time.Since(started).Seconds()),
		))
		t.recordError(span, err)
	}(time.Now())

	return t.ContextExecutor.QueryContext(ctx, query, args...)
}

func (t tracedDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) (row *sql.Row) {
	span := trace.SpanFromContext(ctx)
	defer func(started time.Time) {
		span.AddEvent("QueryRowContext", trace.WithAttributes(
			attribute.String("Query", query),
			attribute.Float64("Took", time.Since(started).Seconds()),
		))
		t.recordError(span, row.Err())
	}(time.Now())

	return t.ContextExecutor.QueryRowContext(ctx, query, args...)
}

func (t tracedDB) recordError(span trace.Span, err error) {
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			span.AddEvent("Database Error", trace.WithAttributes(
				attribute.String("Error", err.Error()),
				attribute.String("Code", pgErr.Code),
				attribute.String("Severity", pgErr.Severity),
				attribute.String("Message", pgErr.Message),
				attribute.String("Detail", pgErr.Detail),
			))
		} else {
			span.AddEvent("Database Error", trace.WithAttributes(
				attribute.String("Error", err.Error()),
			))
		}
	}
}
