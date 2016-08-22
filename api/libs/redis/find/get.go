package find

import(
    libRedis "github.com/scofieldpeng/redis-go"
    "fmt"
    "github.com/garyburd/redigo/redis"
)

// GetToken 获取某个用户的token值列表,返回的第一个参数为token列表值数据,如果出错,第二个值为error
func GetUserTokens(userid int) ([]string,error){
    conn := libRedis.Pool("default").Get()
    defer conn.Close()

    tokens,err := redis.Strings( conn.Do("SMEMBERS",fmt.Sprintf(RedisUserTokens,userid)) );
    if err != nil {
        if err == redis.ErrNil {
            return []string{},nil
        }
        return []string{},err
    }
    return tokens,nil
}

// GetUserIDFromToken 通过userid获取token,如果出错,第二个参数返回error
func GetUserIDFromToken(token string) (int,error) {
    conn := libRedis.Pool("default").Get()
    defer conn.Close()

    userid,err := redis.Int(conn.Do("GET",fmt.Sprintf(RedisToken,token)))
    if err != nil && err != redis.ErrNil {
        return 0,err
    }

    return userid,nil
}

// GetLastTime 获取用户上次记录,如果出错,第二个参数返回error
func GetLastTime(userid int) (int64,error) {
    conn := libRedis.Pool("default").Get()
    defer conn.Close()

    lastTime,err := redis.Int64(conn.Do("GET",fmt.Sprintf(TokenLastTime,userid)))
    if err != nil && err != redis.ErrNil {
        return 0,err
    }

    return lastTime,nil
}
