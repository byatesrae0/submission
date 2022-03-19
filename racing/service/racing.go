package service

import (
	"git.neds.sh/matty/entain/racing/proto/racing"

	"golang.org/x/net/context"
)

// RacesRepo will be used as repository access to races.
type RacesRepo interface {
	// List should return a list of races.
	List(filter *racing.ListRacesRequestFilter) ([]*racing.Race, error)
}

type Racing interface {
	// ListRaces will return a collection of races.
	ListRaces(ctx context.Context, in *racing.ListRacesRequest) (*racing.ListRacesResponse, error)
}

// racingService implements the Racing interface.
type racingService struct {
	racesRepo RacesRepo
}

// NewRacingService instantiates and returns a new racingService.
func NewRacingService(racesRepo RacesRepo) Racing {
	return &racingService{racesRepo}
}

func (s *racingService) ListRaces(ctx context.Context, in *racing.ListRacesRequest) (*racing.ListRacesResponse, error) {
	races, err := s.racesRepo.List(in.Filter)
	if err != nil {
		return nil, err
	}

	return &racing.ListRacesResponse{Races: races}, nil
}
