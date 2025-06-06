package fault

import "errors"

type Tag string

const (
	UNTAGGED              Tag = "UNTAGGED"
	BAD_REQUEST           Tag = "BAD_REQUEST_ERROR"
	NOT_FOUND             Tag = "NOT_FOUND_ERROR"
	INTERNAL_SERVER_ERROR Tag = "INTERNAL_SERVER_ERROR"
	UNAUTHORIZED          Tag = "UNAUTHORIZED_ERROR"
	FORBIDDEN             Tag = "FORBIDDEN_ERROR"
	CONFLICT              Tag = "CONFLICT_ERROR"
	TOO_MANY_REQUESTS     Tag = "TOO_MANY_REQUESTS_ERROR"
	UNPROCESSABLE_ENTITY  Tag = "UNPROCESSABLE_ENTITY_ERROR"
	DISABLED_USER         Tag = "DISABLED_USER_ERROR"
	DB_RESOURCE_NOT_FOUND Tag = "DB_RESOURCE_NOT_FOUND_ERROR"
	INVALID_ENTITY        Tag = "INVALID_ENTITY_ERROR"
	MAILER_ERROR          Tag = "MAILER_ERROR"
	EXPIRED               Tag = "EXPIRED_ERROR"
	CACHE_MISS            Tag = "CACHE_MISS_KEY_ERROR"
	DB_TRANSACTION        Tag = "DB_TRANSACTION_ERROR"
)

func GetTag(err error) Tag {
	if err == nil {
		return UNTAGGED
	}

	for err != nil {
		e, ok := err.(*Fault)
		if ok && e.Tag != "" {
			return e.Tag
		}
		err = errors.Unwrap(err)
	}

	return UNTAGGED
}
