package api

import (
	"github.com/labstack/echo/v4"
	"github.com/lw396/WeComCopilot/internal/errors"
)

func (a *Api) getContactPerson(c echo.Context) (err error) {
	nickname := c.QueryParam("nickname")
	if nickname == "" {
		return errors.New(errors.CodeInvalidParam, "请输入联系人昵称")
	}

	result, err := a.service.GetContactPersonByNickname(c.Request().Context(), nickname)
	if err != nil {
		return
	}
	return OK(c, result)
}