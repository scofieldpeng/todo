package common

import (
	"github.com/labstack/echo"
	"github.com/pquerna/ffjson/ffjson"
	"io/ioutil"
)

// GetBody 获取请求的body
func GetBody(ctx echo.Context) (string, error) {
	bodyByte, err := ioutil.ReadAll(ctx.Request().Body())
	if err != nil {
		return "", err
	}

	return string(bodyByte), nil
}

// GetBodyStruct 获取请求的body体并且将其解析成struct
func GetBodyStruct(ctx echo.Context, bodyStruct interface{}) error {
	body, err := GetBody(ctx)
	if err != nil {
		return err
	}

	if body == "" {
		return nil
	}

	if err := ffjson.Unmarshal([]byte(body), bodyStruct); err != nil {
		return err
	}

	return nil
}
