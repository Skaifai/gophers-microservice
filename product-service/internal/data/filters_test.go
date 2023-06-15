package data

import (
	"github.com/Skaifai/gophers-microservice/product-service/pkg/proto"
	"testing"
)

var testcaseFilterByNameDesc = &proto.Filters{
	Page:         5,
	PageSize:     10,
	Sort:         "-name",
	SortSafeList: []string{"id", "-id", "name", "-name"},
}

var testcaseFilterByIdAsc = &proto.Filters{
	Page:         2,
	PageSize:     5,
	Sort:         "id",
	SortSafeList: []string{"id", "-id", "name", "-name"},
}

func TestSortColumn(t *testing.T) {
	tests := []struct {
		name     string
		filter   *proto.Filters
		expected string
	}{
		{
			name:     "Filter by Name (Descending)",
			filter:   testcaseFilterByNameDesc,
			expected: "name",
		},
		{
			name:     "Filter by ID (Ascending)",
			filter:   testcaseFilterByIdAsc,
			expected: "id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			column := sortColumn(tt.filter)
			if column != tt.expected {
				t.Errorf("sortColumn() returned %s, expected %s", column, tt.expected)
			}
		})
	}
}

func TestSortDirection(t *testing.T) {
	tests := []struct {
		name     string
		filter   *proto.Filters
		expected string
	}{
		{
			name:     "Filter by Name (Descending)",
			filter:   testcaseFilterByNameDesc,
			expected: "DESC",
		},
		{
			name:     "Filter by ID (Ascending)",
			filter:   testcaseFilterByIdAsc,
			expected: "ASC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			direction := sortDirection(tt.filter)
			if direction != tt.expected {
				t.Errorf("sortDirection() returned %s, expected %s", direction, tt.expected)
			}
		})
	}
}

func TestLimit(t *testing.T) {
	tests := []struct {
		name     string
		filter   *proto.Filters
		expected int32
	}{
		{
			name:     "Filter by Name (Descending)",
			filter:   testcaseFilterByNameDesc,
			expected: 10,
		},
		{
			name:     "Filter by ID (Ascending)",
			filter:   testcaseFilterByIdAsc,
			expected: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			limit := limit(tt.filter)
			if limit != tt.expected {
				t.Errorf("limit() returned %d, expected %d", limit, tt.expected)
			}
		})
	}
}

func TestOffSet(t *testing.T) {
	tests := []struct {
		name     string
		filter   *proto.Filters
		expected int32
	}{
		{
			name:     "Filter by Name (Descending)",
			filter:   testcaseFilterByNameDesc,
			expected: 40,
		},
		{
			name:     "Filter by ID (Ascending)",
			filter:   testcaseFilterByIdAsc,
			expected: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			offset := offset(tt.filter)
			if offset != tt.expected {
				t.Errorf("offset() returned %d, expected %d", offset, tt.expected)
			}
		})
	}
}
