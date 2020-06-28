package dao

import (
	"errors"
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
	"github.com/lhc328/go_gateway_demo/dto"
	"github.com/lhc328/go_gateway_demo/public"
	"time"
)

type Admin struct {
	Id int `json:"id" gorm:"primary_key" description:"自增主键"`
	UserName string `json:"username" gorm:"column:user_name" description:"姓名"`
	Salt string `json:"salt" gorm:"column:salt" description:"盐"`
	Password string	`json:"password" gorm:"column:password" description:"密码"`
	CreatedAt time.Time	`json:"create_at" gorm:"column:create_at" description:"创建时间"`
	UpdatedAt time.Time	`json:"update_at" gorm:"column:update_at" description:"修改时间"`
	IsDelete int	`json:"is_delete" gorm:"column:is_delete" description:"是否删除"`
}

func (t *Admin) TableName() string {
	return "gateway_admin"
}

func (t *Admin) Find(c *gin.Context, tx *gorm.DB, search *Admin) (*Admin, error) {
	out := &Admin{}
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (t *Admin) Save(c *gin.Context, tx *gorm.DB) error {
	return tx.SetCtx(public.GetGinTraceContext(c)).Save(t).Error
}

func (t *Admin) LoginCheck(c *gin.Context, tx *gorm.DB, search *dto.AdminLoginInput) (*Admin, error) {
	adminInfo, err := t.Find(c, tx, (&Admin{UserName: search.UserName, IsDelete: 0}))
	if err != nil {
		return nil, errors.New("用户信息不存在")
	}
	saltPassword := public.GenSaltPassword(adminInfo.Salt, search.Password)
	if saltPassword == adminInfo.Password {
		return adminInfo, nil
	}
	return nil, errors.New("姓名或密码错误")
}