package main

import (
	"context"
	"sort"
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
		e := e
		estateCache.Store(e.ID, &e)
	}
	return nil
}

func addEstate(e Estate) {
	estateCache.Store(e.ID, &e)
}

func searchEstatesCache(ctx context.Context, doorHeight, doorWidth, rent *Range, features []string, page, perPage int) ([]Estate, int, error) {
	defer nrsgmt(ctx, "searchEstatesCache").End()

	all := make([]Estate, 0)
	var err error

	estateCache.Range(func(_, v interface{}) bool {
		e := v.(*Estate)
		all = append(all, *e)
		return true
	})

	estates := make([]Estate, 0)
	for _, e := range all {
		if doorHeight != nil {
			if doorHeight.Min != -1 && e.DoorHeight < doorHeight.Min {
				continue
			}
			if doorHeight.Max != -1 && e.DoorHeight >= doorHeight.Max {
				continue
			}
		}

		if doorWidth != nil {
			if doorWidth.Min != -1 && e.DoorWidth < doorWidth.Min {
				continue
			}
			if doorWidth.Max != -1 && e.DoorWidth >= doorWidth.Max {
				continue
			}
		}

		if rent != nil {
			if rent.Min != -1 && e.Rent < rent.Min {
				continue
			}
			if rent.Max != -1 && e.Rent >= rent.Max {
				continue
			}
		}

		if len(features) > 0 {
			var matched bool
			for _, f := range features {
				if strings.Contains(e.Features, f) {
					matched = true
					break
				}
			}
			if !matched {
				break
			}
		}

		estates = append(estates, e)
	}

	left := page * perPage
	right := left + perPage
	total := len(estates)
	if right > total {
		right = total
	}

	sort.Slice(estates, func(i, j int) bool {
		if estates[i].Popularity == estates[j].Popularity {
			return estates[i].ID < estates[j].ID
		}
		return estates[i].Popularity > estates[j].Popularity
	})

	return estates[left:right], total, err
}
