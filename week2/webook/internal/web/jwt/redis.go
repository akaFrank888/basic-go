package jwt

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

//	func NewRedisJTWHandler() RedisJTWHandler {
//		return RedisJTWHandler{
//			signMethod: jwt.SigningMethodHS512,
//			rcExpire:   time.Hour * 24 * 7,
//		}
//	}
type RedisJWTHandler struct {
	client        redis.Cmdable
	signingMethod jwt.SigningMethod
	rcExpiration  time.Duration
}

func (c *RedisJWTHandler) ExtractToken(ctx *gin.Context) string {
	authCode := ctx.GetHeader("Authorization")
	if authCode == "" {
		return authCode
	}
	segs := strings.Split(authCode, " ")
	if len(segs) != 2 {
		return ""
	}
	return segs[1]
}

func NewRedisJWTHandler(client redis.Cmdable) Handler {
	return &RedisJWTHandler{
		client:        client,
		rcExpiration:  time.Hour * 24 * 7,
		signingMethod: jwt.SigningMethodHS512,
	}
}
func (c *RedisJWTHandler) SetLoginToken(ctx *gin.Context, uid int64) error {
	ssid := uuid.New().String()
	err := c.SetRefreshToken(ctx, uid, ssid)
	if err != nil {
		return err
	}
	return c.SetJWTToken(ctx, uid, ssid)
}

func (c *RedisJWTHandler) ClearToken(ctx *gin.Context) error {
	ctx.Header("x-jwt-token", "")
	ctx.Header("x-refresh-token", "")
	uc := ctx.MustGet("user").(UserClaims)

	return c.client.Set(ctx, fmt.Sprintf("users:ssid:%s", uc.Ssid), "", c.rcExpiration).Err()
}

// SetJWTToken 登录成功后，返回长短token中的其中一个token：access-token
func (c *RedisJWTHandler) SetJWTToken(ctx *gin.Context, uid int64, Ssid string) error {
	uc := UserClaims{
		Uid:       uid,
		Ssid:      Ssid,
		UserAgent: ctx.GetHeader("User-Agent"),
		RegisteredClaims: jwt.RegisteredClaims{
			// JWT设置为1min过期
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
	}
	token := jwt.NewWithClaims(c.signingMethod, uc)
	tokenStr, err := token.SignedString(JWTKey)
	if err != nil {
		return err
	}
	ctx.Header("x-jwt-token", tokenStr)
	return nil
}

func (c *RedisJWTHandler) SetRefreshToken(ctx *gin.Context, uid int64, Ssid string) error {
	rc := RefreshClaims{
		Uid:  uid,
		Ssid: Ssid,
		RegisteredClaims: jwt.RegisteredClaims{
			// refresh-token设置为7天过期
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(c.rcExpiration)),
		},
	}
	token := jwt.NewWithClaims(c.signingMethod, rc)
	tokenStr, err := token.SignedString(RefreshKey)
	if err != nil {
		return err
	}
	ctx.Header("x-refresh-token", tokenStr)
	return nil
}

func (c *RedisJWTHandler) CheckSession(ctx *gin.Context, ssid string) error {
	cnt, err := c.client.Exists(ctx, fmt.Sprintf("users:ssid:%s", ssid)).Result()
	if err != nil {
		return err
	}
	if cnt > 0 {
		return errors.New("token无效")
	}
	return nil
}

type RefreshClaims struct {
	jwt.RegisteredClaims
	Uid  int64
	Ssid string
}
type UserClaims struct {
	jwt.RegisteredClaims
	Uid       int64
	UserAgent string
	Ssid      string
}

var JWTKey = []byte("IKD20XkWAXJus2zS7R97SH51K7XgQrLb")
var RefreshKey = []byte("IKD20XkWAXJus2zS7R97SH51K7XgQrLbA")
