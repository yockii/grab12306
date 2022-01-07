package main

import (
	"github.com/yockii/qscore/pkg/authorization"
	"github.com/yockii/qscore/pkg/config"
	"github.com/yockii/qscore/pkg/constant"
	"github.com/yockii/qscore/pkg/database"
	"github.com/yockii/qscore/pkg/domain"
	"github.com/yockii/qscore/pkg/logger"
	"github.com/yockii/qscore/pkg/server"
	"github.com/yockii/qscore/pkg/util"

	"github.com/yockii/grab12306/internal/controller"
	"github.com/yockii/grab12306/internal/model"
	"github.com/yockii/grab12306/internal/service"
)

func main() {
	{
		// 初始化配置项(已引入自动初始化)
		// 使用默认即可，默认是在conf/config.toml
		logger.SetLevel(config.GetString("log.level"))
		logger.SetReportCaller(config.GetBool("log.showCode"))
		logger.SetLogDir(config.GetString("log.dir"), config.GetInt("log.rotate"))
	}
	{ // 初始化数据库，默认使用 database/ driver、host、user、password、db、port、prefix、showSql    log/ level
		database.InitSysDB()
		defer database.Close()
	}
	// 初始化数据
	authorization.Init()
	initialData()

	// 启动服务
	controller.InitRouter()
	logger.Error(server.Start(":" + config.GetString("server.port")))
}

func initialData() {
	database.DB.Sync2(domain.SyncDomains...)
	database.DB.Sync2(model.SyncModels...)

	// 检查并初始化用户角色数据
	checkInitialAuthorizationData()
	// 检查并初始化资源数据
	checkRouteResourceData()
}
func checkRouteResourceData() {
	resources := []*domain.Resource{
		// resource资源
		{ResourceName: "新增资源", ResourceContent: "/api/v1/resource", ResourceType: "route", Action: "post"},
		{ResourceName: "更新资源", ResourceContent: "/api/v1/resource", ResourceType: "route", Action: "put"},
		{ResourceName: "删除资源", ResourceContent: "/api/v1/resource", ResourceType: "route", Action: "delete"},
		{ResourceName: "获取资源详情", ResourceContent: "/api/v1/resource/instance", ResourceType: "route", Action: "get"},
		{ResourceName: "获取资源列表", ResourceContent: "/api/v1/resource/list", ResourceType: "route", Action: "get"},
		// role角色
		{ResourceName: "新增角色", ResourceContent: "/api/v1/role", ResourceType: "route", Action: "post"},
		{ResourceName: "更新角色", ResourceContent: "/api/v1/role", ResourceType: "route", Action: "put"},
		{ResourceName: "删除角色", ResourceContent: "/api/v1/role", ResourceType: "route", Action: "delete"},
		{ResourceName: "获取角色详情", ResourceContent: "/api/v1/role/instance", ResourceType: "route", Action: "get"},
		{ResourceName: "获取角色列表", ResourceContent: "/api/v1/role/list", ResourceType: "route", Action: "get"},
		{ResourceName: "获取角色的资源ID列表", ResourceContent: "/api/v1/role/resourceIds/:id", ResourceType: "route", Action: "get"},
		{ResourceName: "分配角色的资源", ResourceContent: "/api/v1/role/dispatch", ResourceType: "route", Action: "post"},
		// user用户
		{ResourceName: "新增用户", ResourceContent: "/api/v1/user", ResourceType: "route", Action: "post"},
		{ResourceName: "更新用户", ResourceContent: "/api/v1/user", ResourceType: "route", Action: "put"},
		{ResourceName: "删除用户", ResourceContent: "/api/v1/user", ResourceType: "route", Action: "delete"},
		{ResourceName: "获取用户详情", ResourceContent: "/api/v1/user/instance", ResourceType: "route", Action: "get"},
		{ResourceName: "获取用户列表", ResourceContent: "/api/v1/user/list", ResourceType: "route", Action: "get"},
		{ResourceName: "修改密码", ResourceContent: "/api/v1/user/pwd", ResourceType: "route", Action: "put"},
		{ResourceName: "获取用户角色ID列表", ResourceContent: "/api/v1/user/roleIds/:id", ResourceType: "route", Action: "get"},
		{ResourceName: "分配用户角色", ResourceContent: "/api/v1/user/dispatch", ResourceType: "route", Action: "post"},
	}
	session := database.DB.NewSession()
	defer session.Close()
	session.Begin()
	for _, r := range resources {
		c, err := database.DB.Count(r)
		if err != nil {
			logger.Panic(err)
		}
		if c == 0 {
			r.Id = util.GenerateDatabaseID()
			_, err = session.Insert(r)
			if err != nil {
				logger.Panic(err)
			}
		}
	}
	if err := session.Commit(); err != nil {
		logger.Panic(err)
	}
}

func checkInitialAuthorizationData() {
	need := false
	role := &domain.Role{RoleName: constant.DefaultRoleName}
	has, err := database.DB.Get(role)
	if err != nil {
		logger.Error(err)
		return
	}
	if !has {
		role.Id = domain.RoleIdPrefix + util.GenerateDatabaseID()
		_, err = database.DB.Insert(role)
		if err != nil {
			logger.Error(err)
			return
		}
		need = true
	}
	// 处理用户
	admin := &domain.User{Username: constant.DefaultUsername}
	has, err = database.DB.Get(admin)
	if err != nil {
		logger.Error(err)
		return
	}
	if !has {
		admin.Password = "123456"
		_, _, err = service.UserService.Add(admin)
		if err != nil {
			logger.Error(err)
			return
		}
		need = true
	}
	// 若角色或用户不存在，则处理用户角色关系
	if need {
		_, err = authorization.AddSubjectGroup(admin.Id, role.Id, "")
		if err != nil {
			logger.Error(err)
			return
		}
	}
	// 超级管理员赋权
	authorization.SetSuperAdmin(role.Id)
}
