package main

import (
	"context"
	"strings"
	"sync"
)

var estateCache sync.Map

func loadEstates(ctx context.Context) error {
	estates := []Estate{}
	if err := db.SelectContext(ctx, &estates, "SELECT * FROM estate"); err != nil {
		return err
	}
	estateCache = sync.Map{}
	for _, e := range estates {
		estateCache.Store(e.ID, e)
	}
	return nil
}

func addEstate(e Estate) {
	estateCache.Store(e.ID, e)
}

func searchEstatesCache(doorHeight, doorWidth, rent *Range, features []string, page, perPage int) ([]Estate, error) {
	estates := make([]Estate, 0)
	var err error
	estateCache.Range(func(_, v interface{}) bool {
		e := v.(Estate)

		if doorHeight != nil {
			if doorHeight.Min == -1 || e.DoorHeight < doorHeight.Min {
				return true
			}
			if doorHeight.Max == -1 || e.DoorHeight >= doorHeight.Max {
				return true
			}
		}

		if doorWidth != nil {
			if doorWidth.Min == -1 || e.DoorWidth < doorWidth.Min {
				return true
			}
			if doorWidth.Max == -1 || e.DoorWidth >= doorWidth.Max {
				return true
			}
		}

		if rent != nil {
			if rent.Min == -1 || e.Rent < rent.Min {
				return true
			}
			if rent.Max == -1 || e.Rent >= rent.Max {
				return true
			}
		}

		var matched bool
		for _, f := range features {
			if strings.Contains(e.Features, f) {
				matched = true
				break
			}
		}
		if !matched {
			return true
		}

		estates = append(estates, e)
		return true
	})
	left := page * perPage
	right := left + perPage
	total := len(estates)
	if right > total {
		right = total
	}
	return estates[left:right], err
}
