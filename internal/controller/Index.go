package controller

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/yockii/qscore/pkg/server"
	"github.com/yockii/qscore/pkg/util"
)

func InitRouter() {
	server.Post("/login", UserController.Login)

	// Dict
	server.StandardRouter(
		"/dict",
		DictController.Add,
		DictController.Update,
		DictController.Delete,
		DictController.Get,
		DictController.Paginate,
	)
	// Resource
	server.StandardRouter(
		"/resource",
		ResourceController.Add,
		ResourceController.Update,
		ResourceController.Delete,
		ResourceController.Get,
		ResourceController.Paginate,
	)
	// Role
	rg := server.StandardRouter(
		"/role",
		RoleController.Add,
		RoleController.Update,
		RoleController.Delete,
		RoleController.Get,
		RoleController.Paginate,
	)
	rg.Get("/resourceIds/:id", RoleController.GetRoleResourceIds)
	rg.Post("/dispatch", RoleController.DispatchResources)

	// User
	ug := server.StandardRouter(
		"/user",
		UserController.Add,
		UserController.Update,
		UserController.Delete,
		UserController.Get,
		UserController.Paginate,
	)
	ug.Put("/pwd", UserController.UpdateUserPassword)
	ug.Get("/roleIds/:id", UserController.GetUserRoleIds)
	ug.Post("/dispatch", UserController.DispatchRoles)

}

func parsePaginationInfoFromQuery(ctx *fiber.Ctx) (size, offset int, orderBy string, err error) {
	sizeStr := ctx.Query("limit", "10")
	offsetStr := ctx.Query("offset", "0")
	size, err = strconv.Atoi(sizeStr)
	if err != nil {
		return
	}
	offset, err = strconv.Atoi(offsetStr)
	if err != nil {
		return
	}
	if size < -1 || size > 1000 {
		size = 10
	}
	if offset < -1 {
		offset = 0
	}
	orderBy = ctx.Query("orderBy") // orderBy=xxx-desc,yyy-asc,zzz
	if orderBy != "" {
		obs := strings.Split(orderBy, ",")
		ob := ""
		for _, s := range obs {
			kds := strings.Split(s, "-")
			ob += ", " + util.SnakeString(strings.TrimSpace(kds[0]))
			if len(kds) == 2 {
				d := strings.ToLower(kds[1])
				if d == "desc" {
					ob += " DESC"
				}
			}
		}
		orderBy = ob
	}
	return
}
