package service

import "sync"

type RatingStore interface {
	Add(laptopID string, score float64) (*Rating, error)
}

type Rating struct {
	Count uint32
	Sum   float64
}

type InMemoryRatingScore struct {
	mutex  sync.RWMutex
	rating map[string]*Rating
}

func (i *InMemoryRatingScore) Add(laptopID string, score float64) (*Rating, error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	r := i.rating[laptopID]
	if r == nil {
		r = &Rating{
			Count: 1,
			Sum:   score,
		}
	} else {
		r.Count++
		r.Sum += score
	}
	i.rating[laptopID] = r
	return r, nil
}

func NewInMemoryRatingScore() *InMemoryRatingScore {
	return &InMemoryRatingScore{rating: map[string]*Rating{}}
}
