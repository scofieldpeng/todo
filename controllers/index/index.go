package index

import (
    "github.com/labstack/echo"
    "net/http"
)

// Index index页面
func Index(ctx echo.Context) error {
    return ctx.String(http.StatusOK,"hello world")
}
