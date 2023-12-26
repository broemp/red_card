package util

import "database/sql"

func StringToSQLString(value string) sql.NullString {
	if value != "" {
		return sql.NullString{
			String: value,
			Valid:  true,
		}
	}
	return sql.NullString{
		Valid: false,
	}
}
