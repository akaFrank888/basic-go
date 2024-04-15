package cache

import (
	"basic-go/week2/webook/internal/domain"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"time"
)

type UserCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

func (c *UserCache) Get(ctx *gin.Context, uid int64) (domain.User, error) {
	key := c.Key(uid)
	// 假定用JSON来存储val
	val, err := c.cmd.Get(ctx, key).Result()
	if err != nil {
		return domain.User{}, err
	}
	var u domain.User
	// 反序列化
	err = json.Unmarshal([]byte(val), &u)
	return u, err
}

func (c *UserCache) Key(uid int64) string {
	// 格式化字符串
	return fmt.Sprintf("user:info:%d", uid)
}

func (c *UserCache) Set(ctx *gin.Context, du domain.User) error {
	key := c.Key(du.Id)
	// 序列化
	val, err := json.Marshal(du)
	if err != nil {
		return err
	}
	return c.cmd.Set(ctx, key, val, c.expiration).Err()
}

func NewUserCache(cmd redis.Cmdable) *UserCache {
	return &UserCache{
		cmd:        cmd,
		expiration: time.Minute * 15,
	}
}
