package service

import (
	"context"
	"fmt"
	"time"

	"github.com/lahaehae/crud_project/internal/pb"
	"github.com/lahaehae/crud_project/internal/repository"
	"github.com/lahaehae/crud_project/internal/telemetry"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// UserService реализует pb.UserServiceServer
type UserService struct {	
	pb.UnimplementedUserServiceServer
	repo repository.UserRepository
	meter metric.Meter;
	tracer trace.Tracer;
}

// NewUserService создаёт новый экземпляр UserService
func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
		meter: otel.Meter("service"),
		tracer: otel.Tracer("service"),
	}
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
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
	user, err := s.repo.CreateUser(ctx, req.Name, req.Email, req.Balance)
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
	return &pb.UserResponse{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
		Balance: user.Balance,
	}, nil
}

func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
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


	user, err := s.repo.GetUser(ctx, req.Id)
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
	return &pb.UserResponse{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
		Balance: user.Balance,
	}, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	ctx, span := s.tracer.Start(ctx, "Service.UpdateUser")
	defer span.End()

	if telemetry.RequestsCounter != nil {
		telemetry.RequestsCounter.Add(ctx, 1,
			metric.WithAttributes(
				attribute.String("method: ", "UpdateUser"),
			),
		)
	}

	user, err := s.repo.UpdateUser(ctx, req.Id, req.Name, req.Email, req.Balance)
	if err != nil {
		span.RecordError(err)
		telemetry.RecordErrorMetric(ctx, "repo_update_user", err)
		return nil, err
	}
	return &pb.UserResponse{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
		Balance: user.Balance,
	}, nil
}

func (s *UserService) TransferFunds(ctx context.Context, req *pb.TransferFundsUserRequset) (*pb.UserResponse, error) {
	ctx, span := s.tracer.Start(ctx, "Service.TransferFunds")
	defer span.End()
	
	if telemetry.RequestsCounter != nil {
		telemetry.RequestsCounter.Add(ctx, 1,
			metric.WithAttributes(
				attribute.String("method: ", "TransferFunds"),
			),
		)
	}
	user, err := s.repo.TransferFunds(ctx, req.Fromid, req.Toid, req.Balance)
	if err != nil {
		span.RecordError(err)
		telemetry.RecordErrorMetric(ctx, "repo_transfer_funds", err)
		return nil, err
	}
	user1, err := s.repo.GetUser(ctx, req.Toid)
	if err != nil {
		span.RecordError(err)
		telemetry.RecordErrorMetric(ctx, "repo_get_user", err)
		return nil, err
	}
	return &pb.UserResponse{
		Id: user.Id,
		Name: user1.Name,
		Email: user1.Email,
		Balance: user.Balance,
	}, nil
}


func (s *UserService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.Empty, error) {
	ctx, span := s.tracer.Start(ctx, "Service.DeleteUser")
	defer span.End()

	if telemetry.RequestsCounter != nil {
		telemetry.RequestsCounter.Add(ctx, 1,
			metric.WithAttributes(
				attribute.String("method: ", "DeleteUser"),
			),
		)
	}

	err := s.repo.DeleteUser(ctx, req.Id)
	if err != nil {
		span.RecordError(err)
		telemetry.RecordErrorMetric(ctx, "repo_delete_user", err)
		return nil, err
	}
	return &pb.Empty{}, nil
}
