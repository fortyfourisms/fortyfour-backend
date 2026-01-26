package middleware

import "context"

// Reuse the contextKey struct & UserIDKey, Username, Role variables you already have in auth.go

// SetUserID menyimpan userID ke context menggunakan UserIDKey yang sama.
func SetUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// GetUserID membaca userID dari context.
// Mengembalikan "" jika tidak ditemukan.
func GetUserID(ctx context.Context) string {
	if v := ctx.Value(UserIDKey); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// SetUsername / GetUsername
func SetUsername(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, Username, username)
}

func GetUsername(ctx context.Context) string {
	if v := ctx.Value(Username); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// SetRole / GetRole
func SetRole(ctx context.Context, role string) context.Context {
	return context.WithValue(ctx, Role, role)
}

func GetRole(ctx context.Context) string {
	if v := ctx.Value(Role); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
