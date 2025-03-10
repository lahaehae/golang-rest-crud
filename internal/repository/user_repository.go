package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	pb "github.com/lahaehae/crud_project/internal/pb"
	"github.com/lahaehae/crud_project/internal/telemetry"

	//"github.com/lahaehae/crud_project/internal/telemetry"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type UserRepo interface {
	CreateUser(ctx context.Context, name, email string) (*pb.UserResponse, error)
	GetUser(ctx context.Context, id int64) (*pb.UserResponse, error)
	UpdateUser(ctx context.Context, id int64) (*pb.UserResponse, error)
	DeleteUser(ctx context.Context, id int64) error
	TransferFunds(ctx context.Context, fromId, toId, balance int64) (*pb.UserResponse, error)
}

type UserRepository struct {
	db     *pgxpool.Pool
	meter  metric.Meter
	tracer trace.Tracer
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db:     db,
		meter:  otel.Meter("repository"),
		tracer: otel.Tracer("repository"),
	}
}

// method CreateUser without transaction
func (r *UserRepository) CreateUser(ctx context.Context, name, email string, balance int64) (*pb.UserResponse, error) {
	ctx, span := r.tracer.Start(ctx, "Repository.CreateUser")
	defer span.End()

	start := time.Now()

	query := "INSERT INTO users (name, email, balance) VALUES ($1, $2, $3) RETURNING id"
	var id int64
	err := r.db.QueryRow(ctx, query, name, email, balance).Scan(&id)
	if err != nil {
		span.RecordError(err)
		telemetry.ErrorCounter.Add(ctx, 1, metric.WithAttributes(
			attribute.Int64("userId: ", id),
			attribute.String("error.type", fmt.Sprintf("%T", err)),
			attribute.String("error.msg", err.Error()),
			attribute.String("query", query),
		))
		return nil, err
	}
	duration := time.Since(start).Milliseconds()
	span.SetAttributes(
		attribute.Int64("db_query.time_ms", duration),
		attribute.Int64("db_query.user_id", id),
	)

	if telemetry.RepoLatencyRecorder != nil {
		telemetry.RepoLatencyRecorder.Record(ctx, time.Since(start).Seconds())
	}
	return &pb.UserResponse{
		Id:      id,
		Name:    name,
		Email:   email,
		Balance: balance,
	}, nil
}

// method GetUser without transaction
func (r *UserRepository) GetUser(ctx context.Context, id int64) (*pb.UserResponse, error) {
	ctx, span := r.tracer.Start(ctx, "Repository.GetUser")
	defer span.End()

	start := time.Now()

	var user pb.UserResponse
	query := "SELECT id, name, email, balance FROM users WHERE id = $1"
	err := r.db.QueryRow(ctx, query, id).Scan(&user.Id, &user.Name, &user.Email, &user.Balance)
	if err != nil {
		span.RecordError(err)
		telemetry.ErrorCounter.Add(ctx, 1, metric.WithAttributes(
			attribute.String("method:", "GetUser"),
			attribute.String("error.type", fmt.Sprintf("%T", err)),
			attribute.String("error.msg", err.Error()),
			attribute.String("query", query),
		))
		return nil, err
	}

	duration := time.Since(start).Milliseconds()
	span.SetAttributes(
		attribute.Int64("db_query.time_ms", duration),
		attribute.Int64("db_query.user_id", int64(id)),
	)

	if telemetry.RepoLatencyRecorder != nil {
		telemetry.RepoLatencyRecorder.Record(ctx, time.Since(start).Seconds())
	}
	return &user, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, id int64, name, email string, balance int64) (*pb.UserResponse, error) {
	ctx, span := r.tracer.Start(ctx, "Repository.UpdateUser")
	defer span.End()

	start := time.Now()

	query := "UPDATE users SET name = $1, email = $2, balance = $3 WHERE id = $4"
	_, err := r.db.Exec(ctx, query, name, email, balance, id)
	if err != nil {
		span.RecordError(err)
		telemetry.ErrorCounter.Add(ctx, 1, metric.WithAttributes(
			attribute.String("method:", "UpdateUser"),
			attribute.String("error.type", fmt.Sprintf("%T", err)),
			attribute.String("error.msg", err.Error()),
			attribute.String("query", query),
		))
		return nil, err
	}

	duration := time.Since(start).Milliseconds()
	span.SetAttributes(
		attribute.Int64("db_query.time_ms", duration),
		attribute.Int64("db_query.user_id", id),
	)

	if telemetry.RepoLatencyRecorder != nil {
		telemetry.RepoLatencyRecorder.Record(ctx, time.Since(start).Seconds())
	}
	return &pb.UserResponse{
		Id:      id,
		Name:    name,
		Email:   email,
		Balance: balance,
	}, nil
}

func (r *UserRepository) TransferFunds(ctx context.Context, fromId, toId, balance int64) (*pb.UserResponse, error) {
	ctx, span := r.tracer.Start(ctx, "Repository.TransferFunds")
	defer span.End()

	start:= time.Now()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		span.RecordError(err)
		telemetry.RecordErrorMetric(ctx, "begin_transaction", err)
		return nil, err
	}
	defer tx.Rollback(ctx)

	query1 := "UPDATE users SET balance = balance - $1 WHERE id = $2"
	_, err = tx.Exec(ctx, query1, balance, fromId)
	if err != nil {
		span.RecordError(err)
		telemetry.RecordErrorMetric(ctx, "update_balance_from", err)	
		return nil, err
	}
	query2 := "UPDATE users SET balance = balance + $1 WHERE id = $2"
	_, err = tx.Exec(ctx, query2, balance, toId)
	if err != nil {
		span.RecordError(err)
		telemetry.RecordErrorMetric(ctx, "update_balance_to", err)
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		span.RecordError(err)
		telemetry.RecordErrorMetric(ctx, "commit_transaction", err)
		return nil, err
	}

	var newBalance int64
	err = r.db.QueryRow(ctx, "SELECT balance FROM users WHERE id = $1", toId).Scan(&newBalance)
	if err != nil {
		span.RecordError(err)
		telemetry.RecordErrorMetric(ctx, "select_balance", err)
		return nil, err
	}

	duration := time.Since(start).Milliseconds()
	span.SetAttributes(
		attribute.Int64("db_query.time_ms", duration),
		attribute.Int64("db_query.user_fromId", fromId),
		attribute.Int64("db_query.user_toId", toId),
	)

	if telemetry.RepoLatencyRecorder != nil {
		telemetry.RepoLatencyRecorder.Record(ctx, time.Since(start).Seconds())
	}

	return &pb.UserResponse{
		Id:      toId,
		Balance: newBalance,
	}, nil

}

func (r *UserRepository) DeleteUser(ctx context.Context, id int64) error {
	ctx, span := r.tracer.Start(ctx, "Repository.DeleteUser")
	defer span.End()

	query := "DELETE FROM users WHERE id = $1"

	start := time.Now()

	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		span.RecordError(err)
		telemetry.ErrorCounter.Add(ctx, 1, metric.WithAttributes(
			attribute.String("method:", "DeleteUser"),
			attribute.String("error.type", fmt.Sprintf("%T", err)),
			attribute.String("error.msg", err.Error()),
			attribute.String("query", query),
		))
		return err
	}

	duration := time.Since(start).Milliseconds()
	span.SetAttributes(
		attribute.Int64("db_query.time_ms", duration),
		attribute.Int64("db_query.user_id", id),
	)

	if telemetry.RepoLatencyRecorder != nil {
		telemetry.RepoLatencyRecorder.Record(ctx, time.Since(start).Seconds())
	}

	return err
}

