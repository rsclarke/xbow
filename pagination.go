package xbow

import (
	"context"
	"fmt"
	"iter"
)

// ListOptions specifies pagination options for list operations.
type ListOptions struct {
	Limit int
	After string
}

// PageInfo contains pagination metadata.
type PageInfo struct {
	NextCursor *string
	HasMore    bool
}

// Page represents a paginated response.
type Page[T any] struct {
	Items    []T
	PageInfo PageInfo
}

// listFunc is a function that fetches a page of items.
type listFunc[T any] func(ctx context.Context, opts *ListOptions) (*Page[T], error)

// paginate creates an iterator that automatically handles pagination.
func paginate[T any](ctx context.Context, opts *ListOptions, fetch listFunc[T]) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		var zero T

		cursor := ""
		if opts != nil {
			cursor = opts.After
		}

		limit := 0
		if opts != nil {
			limit = opts.Limit
		}

		for {
			pageOpts := &ListOptions{
				Limit: limit,
				After: cursor,
			}

			page, err := fetch(ctx, pageOpts)
			if err != nil {
				yield(zero, err)
				return
			}

			for _, item := range page.Items {
				if !yield(item, nil) {
					return
				}
			}

			if !page.PageInfo.HasMore {
				return
			}

			if page.PageInfo.NextCursor == nil || *page.PageInfo.NextCursor == "" {
				yield(zero, fmt.Errorf("xbow: server indicated more pages but returned no cursor"))
				return
			}
			if *page.PageInfo.NextCursor == cursor {
				yield(zero, fmt.Errorf("xbow: server returned same cursor, stopping to prevent infinite loop"))
				return
			}
			cursor = *page.PageInfo.NextCursor
		}
	}
}

// Collect gathers all items from an iterator into a slice.
func Collect[T any](seq iter.Seq2[T, error]) ([]T, error) {
	var items []T
	for item, err := range seq {
		if err != nil {
			return items, err
		}
		items = append(items, item)
	}
	return items, nil
}
