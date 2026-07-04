package utils

import (
	"budget_tracket/constants"
	"context"
)

func GetUIDFromCtx(ctx context.Context) string {
	val := ctx.Value(constants.USER_ID_KEY)
	if val == nil {
		return ""
	}

	userID, ok := val.(string)
	if !ok {
		return ""
	}

	return userID
}
