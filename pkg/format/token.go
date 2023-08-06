package format

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidToken = errors.New("invalid token")

func ParseUserID(token string) (int, error) {
	parts := strings.Split(token, ":")
	if len(parts) != 2 {
		return 0, ErrInvalidToken
	}
	return strconv.Atoi(parts[0])
}
