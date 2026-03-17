package jwtx

import "github.com/golang-jwt/jwt/v4"

// GetToken 生成 JWT Token
// secretKey: 密钥 (对应前端 .env 的 JWT_SECRET)
// iat: 签发时间 (seconds)
// seconds: 过期时间 (seconds)
// payload: 自定义载荷 (例如用户ID)
func GetToken(secretKey string, iat, seconds int64, payload map[string]interface{}) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat

	// 将自定义负载放入 claims
	for k, v := range payload {
		claims[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// ParseToken 验证并解析 JWT，返回 claims 的 map
func ParseToken(secretKey string, tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// 仅支持 HMAC 签名
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenUnverifiable
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrTokenInvalidClaims
}
