package helpers

import "strings"

// NormalizeSortField maps common frontend camelCase sort fields to backend snake_case.
func NormalizeSortField(field string) string {
	switch strings.ToLower(field) {
	case "createdat":
		return "created_at"
	case "updatedat":
		return "updated_at"
	case "publisherid":
		return "publisher_id"
	case "authorid":
		return "author_id"
	default:
		return field
	}
}
