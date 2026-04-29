package stats

import (
	"testing"
)

func makePagerSnapshot(keys ...int) *Snapshot {
	buckets := make([]BucketStats, len(keys))
	for i, k := range keys {
		buckets[i] = BucketStats{Name: fmt.Sprintf("bucket%d", i), KeyCount: k}
	}
	return NewSnapshot(buckets)
}

func TestPaginateSnapshot_Nil(t *testing.T) {
	p := PaginateSnapshot(nil, DefaultPageOptions())
	if p.Total != 0 || len(p.Items) != 0 {
		t.Fatal("expected empty page for nil snapshot")
	}
}

func TestPaginateSnapshot_ZeroPageSize(t *testing.T) {
	snap := makePagerSnapshot(1, 2, 3)
	p := PaginateSnapshot(snap, PageOptions{Page: 1, PageSize: 0})
	if p.Total != 0 {
		t.Fatal("expected empty page for zero page size")
	}
}

func TestPaginateSnapshot_FirstPage(t *testing.T) {
	snap := makePagerSnapshot(1, 2, 3, 4, 5)
	p := PaginateSnapshot(snap, PageOptions{Page: 1, PageSize: 2})
	if len(p.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(p.Items))
	}
	if p.Total != 5 {
		t.Fatalf("expected total 5, got %d", p.Total)
	}
}

func TestPaginateSnapshot_LastPage(t *testing.T) {
	snap := makePagerSnapshot(1, 2, 3, 4, 5)
	p := PaginateSnapshot(snap, PageOptions{Page: 3, PageSize: 2})
	if len(p.Items) != 1 {
		t.Fatalf("expected 1 item on last page, got %d", len(p.Items))
	}
}

func TestPaginateSnapshot_PageBeyondTotal(t *testing.T) {
	snap := makePagerSnapshot(1, 2)
	p := PaginateSnapshot(snap, PageOptions{Page: 99, PageSize: 5})
	if len(p.Items) != 0 {
		t.Fatalf("expected 0 items for out-of-range page, got %d", len(p.Items))
	}
	if p.Total != 2 {
		t.Fatalf("expected total 2, got %d", p.Total)
	}
}

func TestPage_TotalPages(t *testing.T) {
	p := Page{Total: 10, PageSize: 3}
	if p.TotalPages() != 4 {
		t.Fatalf("expected 4 total pages, got %d", p.TotalPages())
	}
}

func TestPage_HasNextAndHasPrev(t *testing.T) {
	p := Page{Page: 2, PageSize: 2, Total: 6}
	if !p.HasNext() {
		t.Fatal("expected HasNext true")
	}
	if !p.HasPrev() {
		t.Fatal("expected HasPrev true")
	}

	first := Page{Page: 1, PageSize: 2, Total: 6}
	if first.HasPrev() {
		t.Fatal("first page should not have prev")
	}
}

func TestDefaultPageOptions(t *testing.T) {
	opts := DefaultPageOptions()
	if opts.Page != 1 || opts.PageSize != 10 {
		t.Fatalf("unexpected defaults: %+v", opts)
	}
}
