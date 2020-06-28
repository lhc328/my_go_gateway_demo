package controller

import (
	"encoding/json"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/lhc328/go_gateway_demo/dao"
	"github.com/lhc328/go_gateway_demo/dto"
	"github.com/lhc328/go_gateway_demo/middleware"
	"github.com/lhc328/go_gateway_demo/public"
	"time"
)

type AdminLoginController struct {}

func AdminLoginRegister(group *gin.RouterGroup) {
	adminLogin := &AdminLoginController{}
	group.POST("/login", adminLogin.AdminLogin)
	group.GET("/logout", adminLogin.AdminLoginOut)
}

// AdminLogin godoc
// @Summary 管理员登录
// @Description 管理员登录
// @Tags 管理员接口
// @ID /admin_login/login
// @Accept json
// @Produce json
// @Param body body dto.AdminLoginInput true "body"
// @Success 200 {object} middleware.Response{data.DemoInput} "success"
// @Router /admin_login/login [post]
func (adminLogin *AdminLoginController) AdminLogin(c *gin.Context)  {
	params := &dto.AdminLoginInput{}
	if err := params.BindValidParam(c); err!=nil {
		middleware.ResponseError(c, 1001, err)
		return
	}
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	// 1 get admin info
	// 2 salt + admin.password
	// 3 dao.password == sha256(admin.passowrd + salt)
	// 4 gen token
	admin := &dao.Admin{}
	admin, err = admin.LoginCheck(c, tx, params)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	sessionInfo := &dto.AdminSessionInfo{
		ID: admin.Id,
		UserName: admin.UserName,
		LoginTime: time.Now(),
	}
	sessBts, err := json.Marshal(sessionInfo)
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}

	session := sessions.Default(c)
	session.Set(public.AdminSessionInfoKey, string(sessBts))
	session.Save()

	out := &dto.AdminLoginOutput{Token: admin.UserName}
	middleware.ResponseSuccess(c, out)
}

// AdminLogin godoc
// @Summary 管理员退出
// @Description 管理员退出
// @Tags 管理员接口
// @ID /admin_login/logout
// @Accept json
// @Produce json
// @Success 200 {object} middleware.Response{string} "success"
// @Router /admin_login/logout [get]
func (adminLogin *AdminLoginController) AdminLoginOut(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete(public.AdminSessionInfoKey)
	session.Save()

	middleware.ResponseSuccess(c, "")
}