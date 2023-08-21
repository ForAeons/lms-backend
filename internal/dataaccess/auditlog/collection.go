package auditlog

import (
	collection "lms-backend/pkg/collectionquery"
)

func Filters() collection.FilterMap {
	return map[string]collection.Filter{
		"action": collection.StringLikeFilter("action"),
	}
}

func Sorters() collection.SortMap {
	return map[string]collection.Sorter{
		"action":     collection.SortBy("action"),
		"created_at": collection.SortBy("created_at"),
	}
}
