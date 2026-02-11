package constants

type PostgresErrCode string

const (
	UNIQUE_CONSTRAINT_VIOLATION PostgresErrCode = "23505"
)
