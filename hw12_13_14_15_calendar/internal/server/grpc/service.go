package internalgrpc

import (
	"database/sql"

	"github.com/director74/homeworks/hw12_13_14_15_calendar/internal/app"
	"github.com/director74/homeworks/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"github.com/director74/homeworks/hw12_13_14_15_calendar/internal/storage"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/timestamp"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	pb.UnimplementedCalendarServer
	storage app.Storage
}

func NewService(storage app.Storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) Add(ctx context.Context, event *pb.Event) (*pb.AddResponse, error) {
	id, err := s.storage.Add(s.convertRequestEvent(event))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &pb.AddResponse{ID: id}, status.Error(codes.OK, "Success")
}

func (s *Service) Edit(ctx context.Context, event *pb.EditEvent) (*empty.Empty, error) {
	foundedEvent, err := s.storage.GetByID(event.GetID())
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if err := s.storage.Edit(event.ID, s.editMergeEvents(foundedEvent, event)); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &empty.Empty{}, status.Error(codes.OK, "Success")
}

func (s *Service) Delete(ctx context.Context, request *pb.DeleteRequest) (*empty.Empty, error) {
	if err := s.storage.Delete(request.ID); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &empty.Empty{}, status.Error(codes.OK, "Success")
}

func (s *Service) ListEventsDay(
	ctx context.Context,
	date *pb.ListEventsDayRequest,
) (*pb.ListEventsResponse, error) {
	result := make([]*pb.Event, 0)
	events, err := s.storage.ListEventsDay(date.GetDate())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	for _, event := range events {
		result = append(result, s.convertResponseEvent(event))
	}
	return &pb.ListEventsResponse{Events: result}, status.Error(codes.OK, "Success")
}

func (s *Service) ListEventsWeek(
	ctx context.Context,
	date *pb.ListEventsWeekRequest,
) (*pb.ListEventsResponse, error) {
	result := make([]*pb.Event, 0)
	events, err := s.storage.ListEventsWeek(date.GetWeekBeginDate())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	for _, event := range events {
		result = append(result, s.convertResponseEvent(event))
	}
	return &pb.ListEventsResponse{Events: result}, status.Error(codes.OK, "Success")
}

func (s *Service) ListEventsMonth(
	ctx context.Context,
	date *pb.ListEventsMonthRequest,
) (*pb.ListEventsResponse, error) {
	result := make([]*pb.Event, 0)
	events, err := s.storage.ListEventsMonth(date.GetMonthBeginDate())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	for _, event := range events {
		result = append(result, s.convertResponseEvent(event))
	}
	return &pb.ListEventsResponse{Events: result}, status.Error(codes.OK, "Success")
}

func (s *Service) convertRequestEvent(event *pb.Event) storage.Event {
	return storage.Event{
		ID:                   event.ID,
		Title:                event.Title,
		DateStart:            event.DateStart.AsTime(),
		DateEnd:              event.DateEnd.AsTime(),
		Description:          sql.NullString{String: event.GetDescription(), Valid: true},
		UserID:               event.UserID,
		NotificationInterval: sql.NullInt32{Int32: event.GetNotificationInterval(), Valid: true},
	}
}

func (s *Service) editMergeEvents(storageEvent storage.Event, pbEvent *pb.EditEvent) storage.Event {
	result := storageEvent

	if pbEvent.GetTitle() != nil {
		result.Title = pbEvent.GetTitle().Value
	}

	if pbEvent.GetDateStart() != nil {
		result.DateStart = pbEvent.GetDateStart().AsTime()
	}

	if pbEvent.GetDateEnd() != nil {
		result.DateEnd = pbEvent.GetDateEnd().AsTime()
	}

	if pbEvent.GetUserID() != nil {
		result.UserID = pbEvent.GetUserID().Value
	}

	if pbEvent.GetDescription() != nil {
		result.Description = sql.NullString{String: pbEvent.GetDescription().Value, Valid: true}
	}

	if pbEvent.GetNotificationInterval() != nil {
		result.NotificationInterval = sql.NullInt32{Int32: pbEvent.GetNotificationInterval().Value, Valid: true}
	}

	return result
}

func (s *Service) convertResponseEvent(event storage.Event) *pb.Event {
	return &pb.Event{
		ID:                   event.ID,
		Title:                event.Title,
		DateStart:            &timestamp.Timestamp{Seconds: event.DateStart.Unix()},
		DateEnd:              &timestamp.Timestamp{Seconds: event.DateEnd.Unix()},
		UserID:               event.UserID,
		Description:          event.Description.String,
		NotificationInterval: event.NotificationInterval.Int32,
	}
}
