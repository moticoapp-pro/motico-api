package context

import (
	"context"
)

type contextKey string

const TenantIDKey contextKey = "tenant_id"
const UserIDKey contextKey = "user_id"

func GetTenantID(ctx context.Context) string {
	if tenantID, ok := ctx.Value(TenantIDKey).(string); ok {
		return tenantID
	}
	return ""
}

func GetUserID(ctx context.Context) string {
	if userID, ok := ctx.Value(UserIDKey).(string); ok {
		return userID
	}
	return ""
}
