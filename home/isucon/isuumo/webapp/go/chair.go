package main

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
)

var chairCache sync.Map

func loadChairs(ctx context.Context) error {
	chairs := []Chair{}
	if err := db.SelectContext(ctx, &chairs, "SELECT * FROM chair"); err != nil {
		return err
	}
	chairCache = sync.Map{}
	for _, c := range chairs {
		c := c
		chairCache.Store(c.ID, &c)
	}
	return nil
}

func addChair(c Chair) {
	chairCache.Store(c.ID, &c)
}

func searchChairsCache(price, height, width, depth *Range, kind, color string, features []string, page, perPage int) ([]Chair, int, error) {
	chairs := make([]Chair, 0)
	var err error
	chairCache.Range(func(_, v interface{}) bool {
		c := v.(*Chair)

		if c.Stock == 0 {
			return true
		}

		if height != nil {
			if height.Min != -1 && c.Height < height.Min {
				return true
			}
			if height.Max != -1 && c.Height >= height.Max {
				return true
			}
		}

		if width != nil {
			if width.Min != -1 && c.Width < width.Min {
				return true
			}
			if width.Max != -1 && c.Width >= width.Max {
				return true
			}
		}

		if depth != nil {
			if depth.Min != -1 && c.Depth < depth.Min {
				return true
			}
			if depth.Max != -1 && c.Depth >= depth.Max {
				return true
			}
		}

		if kind != "" && c.Kind != kind {
			return true
		}

		if color != "" && c.Color != color {
			return true
		}

		if len(features) > 0 {
			for _, f := range features {
				if !strings.Contains(c.Features, f) {
					return true
				}
			}
		}

		chairs = append(chairs, *c)
		return true
	})
	left := page * perPage
	right := left + perPage
	total := len(chairs)
	if right > total {
		right = total
	}
	fmt.Printf("total=%v [%v:%v]\n", total, left, right)

	sort.Slice(chairs, func(i, j int) bool {
		if chairs[i].Popularity == chairs[j].Popularity {
			return chairs[i].ID < chairs[j].ID
		}
		return chairs[i].Popularity > chairs[j].Popularity
	})
	return chairs[left:right], total, err
}
