package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/lejianwen/rustdesk-api/v2/global"
	"github.com/lejianwen/rustdesk-api/v2/http/response"
	"github.com/lejianwen/rustdesk-api/v2/model"
	"github.com/lejianwen/rustdesk-api/v2/service"
	"os"
	"strings"
)

type Config struct {
}

// ServerConfig RUSTDESK服务配置
// @Tags ADMIN
// @Summary RUSTDESK服务配置
// @Description 服务配置,给webclient提供api-server
// @Accept  json
// @Produce  json
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/config/server [get]
// @Security token
func (co *Config) ServerConfig(c *gin.Context) {
	cf := &response.ServerConfigResponse{
		IdServer:    global.Config.Rustdesk.IdServer,
		Key:         global.Config.Rustdesk.Key,
		RelayServer: global.Config.Rustdesk.RelayServer,
		ApiServer:   global.Config.Rustdesk.ApiServer,
	}
	response.Success(c, cf)
}

// AppConfig APP服务配置
// @Tags ADMIN
// @Summary APP服务配置
// @Description APP服务配置
// @Accept  json
// @Produce  json
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/config/app [get]
// @Security token
func (co *Config) AppConfig(c *gin.Context) {
	response.Success(c, &gin.H{
		"web_client": global.Config.App.WebClient,
	})
}

// AdminConfig ADMIN服务配置
// @Tags ADMIN
// @Summary ADMIN服务配置
// @Description ADMIN服务配置
// @Accept  json
// @Produce  json
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/config/admin [get]
// @Security token
func (co *Config) AdminConfig(c *gin.Context) {

	u := &model.User{}
	token := c.GetHeader("api-token")
	if token != "" {
		u, _ = service.AllService.UserService.InfoByAccessToken(token)
		if !service.AllService.UserService.CheckUserEnable(u) {
			u.Id = 0
		}
	}

	if u.Id == 0 {
		response.Success(c, &gin.H{
			"title": global.Config.Admin.Title,
		})
		return
	}

	hello := global.Config.Admin.Hello
	if hello == "" {
		helloFile := global.Config.Admin.HelloFile
		if helloFile != "" {
			b, err := os.ReadFile(helloFile)
			if err == nil && len(b) > 0 {
				hello = string(b)
			}
		}
	}

	//replace {{username}} to username
	hello = strings.Replace(hello, "{{username}}", u.Username, -1)
	response.Success(c, &gin.H{
		"title": global.Config.Admin.Title,
		"hello": hello,
	})
}

// UpdateAdminConfig 更新ADMIN服务配置
// @Tags ADMIN
// @Summary 更新ADMIN服务配置
// @Description 更新标题和欢迎语
// @Accept  json
// @Produce  json
// @Param title body string false "页面标题"
// @Param hello body string false "欢迎语"
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/config/admin [put]
// @Security token
func (co *Config) UpdateAdminConfig(c *gin.Context) {
	var req struct {
		Title string `json:"title"`
		Hello string `json:"hello"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "参数错误")
		return
	}

	// 更新内存中的配置
	if req.Title != "" {
		global.Config.Admin.Title = req.Title
		global.Viper.Set("admin.title", req.Title)
	}
	if req.Hello != "" {
		global.Config.Admin.Hello = req.Hello
		global.Viper.Set("admin.hello", req.Hello)
		// 清空 hello-file，优先使用 hello 字段
		global.Viper.Set("admin.hello-file", "")
	}

	// 写入配置文件
	if err := global.Viper.WriteConfig(); err != nil {
		response.Error(c, "写入配置文件失败")
		return
	}

	response.Success(c, &gin.H{
		"title": global.Config.Admin.Title,
		"hello": global.Config.Admin.Hello,
	})
}
