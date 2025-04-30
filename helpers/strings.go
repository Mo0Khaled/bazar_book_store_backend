package helpers

import (
	"database/sql"
	"encoding/json"
	"strconv"
)

func ToNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

func StringToNullInt32(value string) sql.NullInt32 {
	var result sql.NullInt32
	if value != "" {
		id, err := strconv.Atoi(value)
		if err == nil {
			result = sql.NullInt32{Int32: int32(id), Valid: true}
		}
	}
	return result
}

func ParseInt32JSON(value string) []int32 {
	var result []int32
	_ = json.Unmarshal([]byte(value), &result)
	return result
}
