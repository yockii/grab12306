package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yockii/qscore/pkg/authorization"
	"github.com/yockii/qscore/pkg/constant"
	"github.com/yockii/qscore/pkg/database"
	coreDomain "github.com/yockii/qscore/pkg/domain"
	"github.com/yockii/qscore/pkg/logger"

	"github.com/yockii/grab12306/internal/domain"
	"github.com/yockii/grab12306/internal/service"
)

var RoleController = new(roleController)

type roleController struct{}

func (c *roleController) Add(ctx *fiber.Ctx) error {
	instance := new(coreDomain.Role)
	if err := ctx.BodyParser(instance); err != nil {
		logger.Error(err)
		return ctx.JSON(&coreDomain.CommonResponse{
			Code: constant.ErrorCodeBodyParse,
			Msg:  "参数解析失败!",
		})
	}

	// 处理必填
	if instance.RoleName == "" {
		return ctx.JSON(&coreDomain.CommonResponse{
			Code: constant.ErrorCodeLackOfField,
			Msg:  "角色名必须提供",
		})
	}

	duplicated, success, err := service.RoleService.Add(instance)
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

func (c *roleController) Delete(ctx *fiber.Ctx) error {
	instance := new(coreDomain.Role)
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
	deleted, err := service.RoleService.Remove(instance)
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

func (c *roleController) Update(ctx *fiber.Ctx) error {
	instance := new(coreDomain.Role)
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
	updated, err := service.RoleService.Update(instance)
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

func (c *roleController) Paginate(ctx *fiber.Ctx) error {
	pr := new(coreDomain.RoleRequest)
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

	total, list, err := service.RoleService.PaginateBetweenTimes(&pr.Role, limit, offset, orderBy, timeRangeMap)
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

func (c *roleController) Get(ctx *fiber.Ctx) error {
	instance := new(coreDomain.Role)
	var err error
	if err = ctx.QueryParser(instance); err != nil {
		logger.Error(err)
		return ctx.JSON(&coreDomain.CommonResponse{
			Code: constant.ErrorCodeBodyParse,
			Msg:  "参数解析失败!",
		})
	}
	instance, err = service.RoleService.Get(instance)
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

func (c *roleController) GetRoleResourceIds(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.JSON(&coreDomain.CommonResponse{
			Code: constant.ErrorCodeLackOfField,
			Msg:  "ID必须提供",
		})
	}

	super, ids, err := authorization.GetSubjectResourceIds(id, "")
	if err != nil {
		logger.Error(err)
		return ctx.JSON(&coreDomain.CommonResponse{
			Code: constant.ErrorCodeService,
			Msg:  "服务出现异常",
		})
	}

	return ctx.JSON(&coreDomain.CommonResponse{Data: fiber.Map{
		"super": super,
		"ids":   ids,
	}})
}

func (c *roleController) DispatchResources(ctx *fiber.Ctx) error {
	urr := new(domain.RoleResourcesRequest)
	if err := ctx.BodyParser(urr); err != nil {
		logger.Error(err)
		return ctx.JSON(&coreDomain.CommonResponse{
			Code: constant.ErrorCodeBodyParse,
			Msg:  "参数解析失败!",
		})
	}
	if _, err := authorization.RemoveSubjectResources(urr.RoleId); err != nil {
		logger.Error(err)
		return ctx.JSON(&coreDomain.CommonResponse{
			Code: constant.ErrorCodeService,
			Msg:  "服务出现异常",
		})
	}
	for _, resourceId := range urr.ResourceIds {
		resource := new(coreDomain.Resource)
		if exist, err := database.DB.ID(resourceId).Get(resource); err != nil {
			logger.Error(err)
			return ctx.JSON(&coreDomain.CommonResponse{
				Code: constant.ErrorCodeService,
				Msg:  "服务出现异常",
			})
		} else if !exist {
			continue
		}
		_, err := authorization.AddSubjectResource(urr.RoleId, resource.ResourceContent, resource.Action, "", resourceId)
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
