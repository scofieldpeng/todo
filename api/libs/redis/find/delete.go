package find

import(
    libRedis "github.com/scofieldpeng/redis-go"
    "github.com/garyburd/redigo/redis"
)

// Delete 删除某个用户的找回密码所有信息
func Delete(userid int) error {
    conn := libRedis.Pool("default").Get()
    defer conn.Close()

    tokens,err := GetUserTokens(userid)
    if err != nil {
        if err == redis.ErrNil {
            return nil
        }
        return err
    }
    tokensInterface := make([]interface{},0)
    for _,token := range tokens {
        tokensInterface = append(tokensInterface,generateRedisToken(token))
    }

    if _,err := conn.Do("DEL",tokensInterface...);err != nil {
        return err
    }

    return nil
}
