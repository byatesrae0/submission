package service

import (
	"fmt"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"git.neds.sh/matty/entain/sports/proto/sports"
)

// SportsService can be used to look up sorting events.
type SportsService struct {
	db *simpleDB
}

// NewSportsService instantiates and returns a new SportsService.
func NewSportsService() *SportsService {
	return &SportsService{
		db: &simpleDB{
			events: []*sports.Event{
				{
					Id:                  1,
					Name:                "Event 1",
					AdvertisedStartTime: timestamppb.New(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
				{
					Id:                  2,
					Name:                "Event 2",
					AdvertisedStartTime: timestamppb.New(time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC)),
				},
				{
					Id:                  3,
					Name:                "Event 3",
					AdvertisedStartTime: timestamppb.New(time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC)),
				},
			},
		},
	}
}

// ListEvents returns a collection of events.
func (s *SportsService) ListEvents(ctx context.Context, in *sports.ListEventsRequest) (*sports.ListEventsResponse, error) {
	return &sports.ListEventsResponse{Events: s.db.list()}, nil
}

// GetEvent returns a single event.
func (s *SportsService) GetEvent(ctx context.Context, in *sports.GetEventRequest) (*sports.Event, error) {
	event := s.db.get(in.Id)

	if event == nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Event with ID %v does not exist.", in.Id))
	}

	return event, nil
}
