package find

import(
    libRedis "github.com/scofieldpeng/redis-go"
    "fmt"
    "time"
)

// SetToken 设置用户的token
func SetToken(userid int,token string) error {
    conn := libRedis.Pool("default").Get()
    defer conn.Close()

    // 设置token所属的用户
    if err := conn.Send("SETEX",fmt.Sprintf(RedisToken,token),RedisTokenExpire,userid);err != nil {
        return err
    }
    // 设置用户最近一次请求的时间
    if err := conn.Send("SETEX",fmt.Sprintf(TokenLastTime,userid),TokenTimeExpire,time.Now().Unix());err != nil {
        return err
    }
    // 设置该用户的token列表
    if err := conn.Send("SADD",fmt.Sprintf(RedisUserTokens,userid),token);err != nil {
        return err
    }

    return nil
}
