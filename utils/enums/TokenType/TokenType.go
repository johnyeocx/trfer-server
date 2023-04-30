package TokenType

import "github.com/johnyeocx/usual/server2/utils/enums"

const (
	User 	enums.TokenType = "user"
)

func StrToTokenType(tokType string) (enums.TokenType) {
	if (tokType == "user") {
		return User
	}

	return User
}