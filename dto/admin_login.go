package dto

import (
	"github.com/gin-gonic/gin"
	"github.com/lhc328/go_gateway_demo/public"
	"time"
)

type AdminLoginInput struct {
	UserName string `json:"username" form:"username" comment:"姓名" example:"admin"
validate: "required,is_valid_username"` //管理员姓名
	Password string `json:"password" form:"password" comment:"密码" example:"123456" 
validate: "required"`	//密码
}

type AdminSessionInfo struct {
	ID int `json:"id"`
	UserName string `json:"username"`
	LoginTime time.Time `json:"login_time"`
}

// 校验
func (param *AdminLoginInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}

type AdminLoginOutput struct {
	Token string `json: "token" form: "token" comment: "token" example: "token" validate:""`
}
