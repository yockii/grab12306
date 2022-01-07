package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yockii/qscore/pkg/authorization"
	"github.com/yockii/qscore/pkg/constant"
	coreDomain "github.com/yockii/qscore/pkg/domain"
	"github.com/yockii/qscore/pkg/logger"

	"github.com/yockii/grab12306/internal/domain"
	"github.com/yockii/grab12306/internal/service"
)

var UserController = new(userController)

type userController struct{}

func (c *userController) Login(ctx *fiber.Ctx) error {
	instance := new(coreDomain.User)
	if err := ctx.BodyParser(instance); err != nil {
		logger.Error(err)
		return ctx.JSON(&coreDomain.CommonResponse{
			Code: constant.ErrorCodeBodyParse,
			Msg:  "参数解析失败!",
		})
	}
	if instance.Username == "" || instance.Password == "" {
		return ctx.JSON(&coreDomain.CommonResponse{
			Code: constant.ErrorCodeLackOfField,
			Msg:  "用户名及密码必须提供",
		})
	}

	user, token, err := service.UserService.Login(instance)
	if err != nil {
		logger.Error(err)
		return ctx.JSON(&coreDomain.CommonResponse{
			Code: constant.ErrorCodeService,
			Msg:  "服务出现异常",
		})
	}
	if token == "" {
		return ctx.JSON(&coreDomain.CommonResponse{
			Code: constant.ErrorCodeNotFound,
			Msg:  "登录失败",
		})
	}

	// 获取权限
	super, rIds, err := authorization.GetSubjectResourceIds(user.Id, "")
	if err != nil {
		logger.Error(err)
		return ctx.JSON(&coreDomain.CommonResponse{
			Code: constant.ErrorCodeService,
			Msg:  "服务出现异常",
		})
	}
	var resourceList []*coreDomain.Resource
	if !super {
		resourceList, err = service.ResourceService.GetListByIds(rIds)
		if err != nil {
			logger.Error(err)
			return ctx.JSON(&coreDomain.CommonResponse{
				Code: constant.ErrorCodeService,
				Msg:  "服务出现异常",
			})
		}
	}

	return ctx.JSON(&coreDomain.CommonResponse{
		Data: fiber.Map{
			"token":     token,
			"user":      user,
			"isSuper":   super,
			"resources": resourceList,
		},
	})
}

func (c *userController) Add(ctx *fiber.Ctx) error {
	instance := new(coreDomain.User)
	if err := ctx.BodyParser(instance); err != nil {
		logger.Error(err)
		return ctx.JSON(&coreDomain.CommonResponse{
			Code: constant.ErrorCodeBodyParse,
			Msg:  "参数解析失败!",
		})
	}

	// 处理必填
	if instance.Username == "" || instance.Password == "" {
		return ctx.JSON(&coreDomain.CommonResponse{
			Code: constant.ErrorCodeLackOfField,
			Msg:  "用户名/密码必须提供",
		})
	}

	duplicated, success, err := service.UserService.Add(instance)
	if err != nil {
		logger.Error(err)
		return ctx.JSON(&coreDomain.CommonResponse{
			Code: constant.ErrorCodeService,
			Msg:  "服务出现异常",
		})
	}
	if duplicated {
		return ctx.JSON(&coreDomain.CommonResponse{
			Code: constant.ErrorCodeDuplicate,
			Msg:  "有重复记录",
		})
	}
	if success {
		return ctx.JSON(&coreDomain.CommonResponse{Data: instance})
	}
	return ctx.JSON(&coreDomain.CommonResponse{
		Code: constant.ErrorCodeUnknown,
		Msg:  "服务出现异常",
	})
}

func (c *userController) Delete(ctx *fiber.Ctx) error {
	instance := new(coreDomain.User)
	if err := ctx.QueryParser(instance); err != nil {
		logger.Error(err)
		return ctx.JSON(&coreDomain.CommonResponse{
			Code: constant.ErrorCodeBodyParse,
			Msg:  "参数解析失败!",
		})
	}
	if instance.Id == "" {
		return ctx.JSON(&coreDomain.CommonResponse{
			Code: constant.ErrorCodeLackOfField,
			Msg:  "ID必须提供",
		})
	}
	deleted, err := service.UserService.Remove(instance)
	if err != nil {
		logger.Error(err)
		return ctx.JSON(&coreDomain.CommonResponse{
			Code: constant.ErrorCodeService,
			Msg:  "服务出现异常",
		})
	}
	if deleted {
		return ctx.JSON(&coreDomain.CommonResponse{})
	}
	return ctx.JSON(&coreDomain.CommonResponse{
		Msg:  "无数据被删除",
		Data: false,
	})
}

func (c *userController) Update(ctx *fiber.Ctx) error {
	instance := new(coreDomain.User)
	if err := ctx.BodyParser(instance); err != nil {
		logger.Error(err)
		return ctx.JSON(&coreDomain.CommonResponse{
			Code: constant.ErrorCodeBodyParse,
			Msg:  "参数解析失败!",
		})
	}
	if instance.Id == "" {
		return ctx.JSON(&coreDomain.CommonResponse{
			Code: constant.ErrorCodeLackOfField,
			Msg:  "ID必须提供",
		})
	}
	updated, err := service.UserService.Update(instance)
	if err != nil {
		logger.Error(err)
		return ctx.JSON(&coreDomain.CommonResponse{
			Code: constant.ErrorCodeService,
			Msg:  "服务出现异常",
		})
	}
	if updated {
		return ctx.JSON(&coreDomain.CommonResponse{})
	}
	return ctx.JSON(&coreDomain.CommonResponse{
		Msg:  "无数据被更新",
		Data: false,
	})
}

func (c *userController) UpdateUserPassword(ctx *fiber.Ctx) error {
	instance := new(domain.UserPasswordRequest)
	if err := ctx.BodyParser(instance); err != nil {
		logger.Error(err)
		return ctx.JSON(&coreDomain.CommonResponse{
			Code: constant.ErrorCodeBodyParse,
			Msg:  "参数解析失败!",
		})
	}
	if instance.Id == "" {
		return ctx.JSON(&coreDomain.CommonResponse{
			Code: constant.ErrorCodeLackOfField,
			Msg:  "ID必须提供",
		})
	}
	updated, err := service.UserService.UpdatePassword(&instance.User, instance.NewPassword)
	if err != nil {
		logger.Error(err)
		return ctx.JSON(&coreDomain.CommonResponse{
			Code: constant.ErrorCodeService,
			Msg:  "服务出现异常",
		})
	}
	if updated {
		return ctx.JSON(&coreDomain.CommonResponse{})
	}
	return ctx.JSON(&coreDomain.CommonResponse{
		Msg:  "无数据被更新",
		Data: false,
	})
}

func (c *userController) Paginate(ctx *fiber.Ctx) error {
	pr := new(coreDomain.UserRequest)
	if err := ctx.QueryParser(pr); err != nil {
		logger.Error(err)
		return ctx.JSON(&coreDomain.CommonResponse{
			Code: constant.ErrorCodeBodyParse,
			Msg:  "参数解析失败!",
		})
	}
	limit, offset, orderBy, err := parsePaginationInfoFromQuery(ctx)
	if err != nil {
		logger.Error(err)
		return ctx.JSON(&coreDomain.CommonResponse{
			Code: constant.ErrorCodeBodyParse,
			Msg:  "参数解析失败!",
		})
	}

	timeRangeMap := make(map[string]*coreDomain.TimeCondition)
	if pr.CreateTimeRange != nil {
		timeRangeMap["create_time"] = &coreDomain.TimeCondition{
			Start: pr.CreateTimeRange.Start,
			End:   pr.CreateTimeRange.End,
		}
	}

	total, list, err := service.UserService.PaginateBetweenTimes(&pr.User, limit, offset, orderBy, timeRangeMap)
	if err != nil {
		logger.Error(err)
		return ctx.JSON(&coreDomain.CommonResponse{
			Code: constant.ErrorCodeService,
			Msg:  "服务出现异常",
		})
	}
	return ctx.JSON(&coreDomain.CommonResponse{Data: &coreDomain.Paginate{
		Total:  total,
		Offset: offset,
		Limit:  limit,
		Items:  list,
	}})
}

func (c *userController) Get(ctx *fiber.Ctx) error {
	instance := new(coreDomain.User)
	var err error
	if err = ctx.QueryParser(instance); err != nil {
		logger.Error(err)
		return ctx.JSON(&coreDomain.CommonResponse{
			Code: constant.ErrorCodeBodyParse,
			Msg:  "参数解析失败!",
		})
	}
	instance, err = service.UserService.Get(instance)
	if err != nil {
		logger.Error(err)
		return ctx.JSON(&coreDomain.CommonResponse{
			Code: constant.ErrorCodeService,
			Msg:  "服务出现异常",
		})
	}
	return ctx.JSON(&coreDomain.CommonResponse{
		Data: instance,
	})
}

func (c *userController) GetUserRoleIds(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.JSON(&coreDomain.CommonResponse{
			Code: constant.ErrorCodeLackOfField,
			Msg:  "ID必须提供",
		})
	}

	ids, err := authorization.GetSubjectGroupIds(id, "")
	if err != nil {
		logger.Error(err)
		return ctx.JSON(&coreDomain.CommonResponse{
			Code: constant.ErrorCodeService,
			Msg:  "服务出现异常",
		})
	}

	return ctx.JSON(&coreDomain.CommonResponse{Data: ids})
}

func (c *userController) DispatchRoles(ctx *fiber.Ctx) error {
	urr := new(domain.UserRolesRequest)
	if err := ctx.BodyParser(urr); err != nil {
		logger.Error(err)
		return ctx.JSON(&coreDomain.CommonResponse{
			Code: constant.ErrorCodeBodyParse,
			Msg:  "参数解析失败!",
		})
	}
	if _, err := authorization.RemoveSubjectGroups(urr.UserId, ""); err != nil {
		logger.Error(err)
		return ctx.JSON(&coreDomain.CommonResponse{
			Code: constant.ErrorCodeService,
			Msg:  "服务出现异常",
		})
	}
	for _, roleId := range urr.RoleIds {
		_, err := authorization.AddSubjectGroup(urr.UserId, roleId, "")
		if err != nil {
			logger.Error(err)
			return ctx.JSON(&coreDomain.CommonResponse{
				Code: constant.ErrorCodeService,
				Msg:  "服务出现异常",
			})
		}
	}
	return ctx.JSON(&coreDomain.CommonResponse{})
}
