package stats

// Page represents a single page of bucket stats.
type Page struct {
	Items    []BucketStats
	Page     int
	PageSize int
	Total    int
}

// PageOptions configures pagination behaviour.
type PageOptions struct {
	Page     int // 1-based
	PageSize int
}

// DefaultPageOptions returns sensible pagination defaults.
func DefaultPageOptions() PageOptions {
	return PageOptions{
		Page:     1,
		PageSize: 10,
	}
}

// TotalPages returns the number of pages required to display all items.
func (p Page) TotalPages() int {
	if p.PageSize <= 0 {
		return 0
	}
	if p.Total == 0 {
		return 0
	}
	total := p.Total / p.PageSize
	if p.Total%p.PageSize != 0 {
		total++
	}
	return total
}

// HasNext reports whether there is a next page.
func (p Page) HasNext() bool {
	return p.Page < p.TotalPages()
}

// HasPrev reports whether there is a previous page.
func (p Page) HasPrev() bool {
	return p.Page > 1
}

// PaginateSnapshot returns a single page of bucket stats from the snapshot.
// If snap is nil or opts.PageSize <= 0 an empty Page is returned.
func PaginateSnapshot(snap *Snapshot, opts PageOptions) Page {
	if snap == nil || opts.PageSize <= 0 {
		return Page{}
	}
	if opts.Page < 1 {
		opts.Page = 1
	}

	buckets := snap.Buckets()
	total := len(buckets)

	start := (opts.Page - 1) * opts.PageSize
	if start >= total {
		return Page{
			Items:    []BucketStats{},
			Page:     opts.Page,
			PageSize: opts.PageSize,
			Total:    total,
		}
	}

	end := start + opts.PageSize
	if end > total {
		end = total
	}

	return Page{
		Items:    buckets[start:end],
		Page:     opts.Page,
		PageSize: opts.PageSize,
		Total:    total,
	}
}
