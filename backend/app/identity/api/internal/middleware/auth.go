package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"go-zero-learning/backend/common/jwtx"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// JWTMiddleware 会尝试从 Authorization 头解析 Token，并把 userId 写入 Context（key: "userId"）
// 这是一个轻量提取器；实际的鉴权仍建议配合 go-zero 的 rest.WithJwt 使用。
func JWTMiddleware(secret string) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" {
				// 没有 Token，直接继续（一些路由可能不需要鉴权）
				next(w, r)
				return
			}

			// 支持 "Bearer <token>"
			token := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer"))
			token = strings.TrimSpace(token)
			if token == "" {
				httpx.ErrorCtx(r.Context(), w, errors.New("invalid token"))
				return
			}

			claims, err := jwtx.ParseToken(secret, token)
			if err != nil {
				httpx.ErrorCtx(r.Context(), w, errors.New("invalid token: "+err.Error()))
				return
			}

			// 把 userId 放到 Context，供逻辑层读取
			ctx := context.WithValue(r.Context(), "userId", claims["userId"])
			next(w, r.WithContext(ctx))
		}
	}
}
