package data

import (
	"github.com/Skaifai/gophers-microservice/product-service/pkg/proto"
	"math"
	"strings"
)

func sortColumn(filters *proto.Filters) string {
	for _, safeValue := range filters.SortSafeList {
		if filters.Sort == safeValue {
			return strings.TrimPrefix(filters.Sort, "-")
		}
	}
	panic("unsafe sort parameter: " + filters.Sort)
}

func sortDirection(filters *proto.Filters) string {
	if strings.HasPrefix(filters.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

func limit(filters *proto.Filters) int32 {
	return filters.PageSize
}

func offset(filters *proto.Filters) int32 {
	return (filters.Page - 1) * filters.PageSize
}

func calculateMetadata(totalRecords int32, filters *proto.Filters) *proto.Metadata {
	if totalRecords == 0 {
		return &proto.Metadata{}
	}
	return &proto.Metadata{
		CurrentPage:  filters.Page,
		PageSize:     filters.PageSize,
		FirstPage:    1,
		LastPage:     int32(math.Ceil(float64(totalRecords) / float64(filters.PageSize))),
		TotalRecords: totalRecords,
	}
}
