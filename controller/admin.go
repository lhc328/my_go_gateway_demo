package controller

import (
	"encoding/json"
	"fmt"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/lhc328/go_gateway_demo/dao"
	"github.com/lhc328/go_gateway_demo/dto"
	"github.com/lhc328/go_gateway_demo/middleware"
	"github.com/lhc328/go_gateway_demo/public"
)

type AdminController struct{}

func AdminRegister(group *gin.RouterGroup) {
	adminLogin := &AdminController{}
	group.GET("/admin_info", adminLogin.AdminInfo)
	group.POST("/change_pwd", adminLogin.ChangePwd)
}

// AdminInfo godoc
// @Summary 管理员信息
// @Description 管理员信息
// @Tags 管理员接口
// @ID /admin/admin_info
// @Accept json
// @Produce json
// @Success 200 {object} middleware.Response{data=dto.AdminInfoOutput} "success"
// @Router /admin/admin_info [get]
func (adminInfo *AdminController) AdminInfo(c *gin.Context)  {
	sess := sessions.Default(c)
	sessInfo := sess.Get(public.AdminSessionInfoKey)
	adminSessionInfo := &dto.AdminSessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)),
		adminSessionInfo); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	out := &dto.AdminInfoOutput{
		ID: adminSessionInfo.ID,
		Name: adminSessionInfo.UserName,
		LoginTime: adminSessionInfo.LoginTime,
		Avatar:       "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
		Introduction: "I am a super administrator",
		Roles:        []string{"admin"},
	}
	middleware.ResponseSuccess(c, out)
}

// AdminInfo godoc
// @Summary 修改管理员密码
// @Description 修改管理员密码
// @Tags 管理员接口
// @ID /admin/change_pwd
// @Accept json
// @Produce json
// @Param body body dto.ChangePassInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /admin/change_pwd [post]
func (adminLogin *AdminController) ChangePwd(c *gin.Context)  {
	params := &dto.ChangePassInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	// 1 session get user info
	sess := sessions.Default(c)
	sessInfo := sess.Get(public.AdminSessionInfoKey)
	adminSessionInfo := &dto.AdminSessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)),
		adminSessionInfo); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	// 2 session.ID get mysql user info
	adminInfo := &dao.Admin{}
	adminInfo, err = adminInfo.Find(c, tx, (&dao.Admin{Id: adminSessionInfo.ID}))
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	// 3 params.password + salt --> new saltPassword
	newSaltPassword := public.GenSaltPassword(adminInfo.Salt, params.Password)
	adminInfo.Password = newSaltPassword

	// 4 save new user password
	if err = adminInfo.Save(c, tx); err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	middleware.ResponseSuccess(c, "")
}