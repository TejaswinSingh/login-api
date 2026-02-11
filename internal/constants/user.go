package constants

const (
	MIN_USER_PASSWORD_LEN_BYTES = 8
	MAX_USER_PASSWORD_LEN_BYTES = 72 // bcrypt can only hash password upto this length
	MIN_USER_USERNAME_LEN       = 8
	MAX_USER_USERNAME_LEN       = 128
)
