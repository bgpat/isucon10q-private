package main

import (
	"context"
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

func searchEstatesCache(q map[string]string) []Estate {
	estates := make([]Estate, 0)
	estateCache.Range(func(_, v interface{}) bool {
		e := v.(Estate)
		estates = append(estates, e)
		return true
	})
	return estates
}
