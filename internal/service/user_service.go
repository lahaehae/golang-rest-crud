package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/lahaehae/crud_project/internal/models"
	"github.com/lahaehae/crud_project/internal/repository"
	"github.com/lahaehae/crud_project/internal/telemetry"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type UserService struct {	
	repo repository.UserRepo
	meter metric.Meter;
	tracer trace.Tracer;
}

func NewUserService(repo repository.UserRepo) *UserService {
	return &UserService{
		repo: repo,
		meter: otel.Meter("service"),
		tracer: otel.Tracer("service"),
	}
}

func (s *UserService) CreateUser(ctx context.Context, name, email string, balance int64 ) (*models.User, error) {
	ctx, span := s.tracer.Start(ctx, "Service.CreateUser")
	defer span.End()

	start := time.Now()

	if telemetry.RequestsCounter != nil {
		telemetry.RequestsCounter.Add(ctx, 1,
			metric.WithAttributes(
				attribute.String("method: ", "CreateUser"),
			),
		)
	}
	user, err := s.repo.CreateUser(ctx, name, email, balance)
	if err != nil {
		span.RecordError(err)
		telemetry.ErrorCounter.Add(ctx, 1, metric.WithAttributes(
			attribute.Int64("userId: ", int64(user.Id)),
			attribute.String("error.type", fmt.Sprintf("%T", err)),
			attribute.String("error.msg", err.Error()),
			))
		return nil, err	
	}

	if telemetry.LatencyRecorder != nil{
		telemetry.LatencyRecorder.Record(ctx, time.Since(start).Seconds())
	}
	return &models.User{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
		Balance: user.Balance,
	}, nil
}

func (s *UserService) GetUser(ctx context.Context, id int64) (*models.User, error) {
	ctx, span := s.tracer.Start(ctx, "Service.GetUser")
	defer span.End()

	start := time.Now()
	
	if telemetry.RequestsCounter != nil {
		telemetry.RequestsCounter.Add(ctx, 1,
			metric.WithAttributes(
				attribute.String("method: ", "GetUser"),
			),
		)
	}


	user, err := s.repo.GetUser(ctx, id)
	if user == nil {
		return nil, errors.New("user not found")
	}
	if err != nil {
		span.RecordError(err)
		telemetry.ErrorCounter.Add(ctx, 1, metric.WithAttributes(
			attribute.String("error.type", fmt.Sprintf("%T", err)),
			attribute.String("error.msg", err.Error()),
			))
		return nil, err	
	}

	if telemetry.LatencyRecorder != nil{
		telemetry.LatencyRecorder.Record(ctx, time.Since(start).Seconds())
	}
	return &models.User{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
		Balance: user.Balance,
	}, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id int64, name, email string, balance int64) (*models.User, error) {
	ctx, span := s.tracer.Start(ctx, "Service.UpdateUser")
	defer span.End()

	if telemetry.RequestsCounter != nil {
		telemetry.RequestsCounter.Add(ctx, 1,
			metric.WithAttributes(
				attribute.String("method: ", "UpdateUser"),
			),
		)
	}

	user, err := s.repo.UpdateUser(ctx, id, name, email, balance)
	if err != nil {
		span.RecordError(err)
		telemetry.RecordErrorMetric(ctx, "repo_update_user", err)
		return nil, err
	}
	return &models.User{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
		Balance: user.Balance,
	}, nil
}

func (s *UserService) TransferFunds(ctx context.Context, fromId, toId, balance int64) (*models.User, error) {
	ctx, span := s.tracer.Start(ctx, "Service.TransferFunds")
	defer span.End()
	
	if telemetry.RequestsCounter != nil {
		telemetry.RequestsCounter.Add(ctx, 1,
			metric.WithAttributes(
				attribute.String("method: ", "TransferFunds"),
			),
		)
	}
	user, err := s.repo.TransferFunds(ctx, fromId, toId, balance)
	if err != nil {
		span.RecordError(err)
		telemetry.RecordErrorMetric(ctx, "repo_transfer_funds", err)
		return nil, err
	}
	user1, err := s.repo.GetUser(ctx, toId)
	if err != nil {
		span.RecordError(err)
		telemetry.RecordErrorMetric(ctx, "repo_get_user", err)
		return nil, err
	}
	return &models.User{
		Id:      user.Id,
		Name:    user1.Name,
		Email:   user1.Email,
		Balance: user.Balance,
	}, nil
}


func (s *UserService) DeleteUser(ctx context.Context, id int64)  error {
	ctx, span := s.tracer.Start(ctx, "Service.DeleteUser")
	defer span.End()

	if telemetry.RequestsCounter != nil {
		telemetry.RequestsCounter.Add(ctx, 1,
			metric.WithAttributes(
				attribute.String("method: ", "DeleteUser"),
			),
		)
	}

	err := s.repo.DeleteUser(ctx, id)
	if err != nil {
		span.RecordError(err)
		telemetry.RecordErrorMetric(ctx, "repo_delete_user", err)
		return err
	}
	return err
}