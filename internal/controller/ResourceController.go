package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yockii/qscore/pkg/constant"
	"github.com/yockii/qscore/pkg/domain"
	"github.com/yockii/qscore/pkg/logger"

	"github.com/yockii/grab12306/internal/service"
)

var ResourceController = new(resourceController)

type resourceController struct{}

func (c *resourceController) Add(ctx *fiber.Ctx) error {
	instance := new(domain.Resource)
	if err := ctx.BodyParser(instance); err != nil {
		logger.Error(err)
		return ctx.JSON(&domain.CommonResponse{
			Code: constant.ErrorCodeBodyParse,
			Msg:  "参数解析失败!",
		})
	}

	// 处理必填
	if instance.ResourceName == "" || instance.ResourceType == "" || instance.ResourceContent == "" || instance.Action == "" {
		return ctx.JSON(&domain.CommonResponse{
			Code: constant.ErrorCodeLackOfField,
			Msg:  "资源名/资源类型/资源内容/资源操作必须提供",
		})
	}

	duplicated, success, err := service.ResourceService.Add(instance)
	if err != nil {
		logger.Error(err)
		return ctx.JSON(&domain.CommonResponse{
			Code: constant.ErrorCodeService,
			Msg:  "服务出现异常",
		})
	}
	if duplicated {
		return ctx.JSON(&domain.CommonResponse{
			Code: constant.ErrorCodeDuplicate,
			Msg:  "有重复记录",
		})
	}
	if success {
		return ctx.JSON(&domain.CommonResponse{Data: instance})
	}
	return ctx.JSON(&domain.CommonResponse{
		Code: constant.ErrorCodeUnknown,
		Msg:  "服务出现异常",
	})
}

func (c *resourceController) Delete(ctx *fiber.Ctx) error {
	instance := new(domain.Resource)
	if err := ctx.QueryParser(instance); err != nil {
		logger.Error(err)
		return ctx.JSON(&domain.CommonResponse{
			Code: constant.ErrorCodeBodyParse,
			Msg:  "参数解析失败!",
		})
	}
	if instance.Id == "" {
		return ctx.JSON(&domain.CommonResponse{
			Code: constant.ErrorCodeLackOfField,
			Msg:  "ID必须提供",
		})
	}
	deleted, err := service.ResourceService.Remove(instance)
	if err != nil {
		logger.Error(err)
		return ctx.JSON(&domain.CommonResponse{
			Code: constant.ErrorCodeService,
			Msg:  "服务出现异常",
		})
	}
	if deleted {
		return ctx.JSON(&domain.CommonResponse{})
	}
	return ctx.JSON(&domain.CommonResponse{
		Msg:  "无数据被删除",
		Data: false,
	})
}

func (c *resourceController) Update(ctx *fiber.Ctx) error {
	instance := new(domain.Resource)
	if err := ctx.BodyParser(instance); err != nil {
		logger.Error(err)
		return ctx.JSON(&domain.CommonResponse{
			Code: constant.ErrorCodeBodyParse,
			Msg:  "参数解析失败!",
		})
	}
	if instance.Id == "" {
		return ctx.JSON(&domain.CommonResponse{
			Code: constant.ErrorCodeLackOfField,
			Msg:  "ID必须提供",
		})
	}
	updated, err := service.ResourceService.Update(instance)
	if err != nil {
		logger.Error(err)
		return ctx.JSON(&domain.CommonResponse{
			Code: constant.ErrorCodeService,
			Msg:  "服务出现异常",
		})
	}
	if updated {
		return ctx.JSON(&domain.CommonResponse{})
	}
	return ctx.JSON(&domain.CommonResponse{
		Msg:  "无数据被更新",
		Data: false,
	})
}

func (c *resourceController) Paginate(ctx *fiber.Ctx) error {
	pr := new(domain.ResourceRequest)
	if err := ctx.QueryParser(pr); err != nil {
		logger.Error(err)
		return ctx.JSON(&domain.CommonResponse{
			Code: constant.ErrorCodeBodyParse,
			Msg:  "参数解析失败!",
		})
	}
	limit, offset, orderBy, err := parsePaginationInfoFromQuery(ctx)
	if err != nil {
		logger.Error(err)
		return ctx.JSON(&domain.CommonResponse{
			Code: constant.ErrorCodeBodyParse,
			Msg:  "参数解析失败!",
		})
	}

	timeRangeMap := make(map[string]*domain.TimeCondition)
	if pr.CreateTimeRange != nil {
		timeRangeMap["create_time"] = &domain.TimeCondition{
			Start: pr.CreateTimeRange.Start,
			End:   pr.CreateTimeRange.End,
		}
	}

	total, list, err := service.ResourceService.PaginateBetweenTimes(&pr.Resource, limit, offset, orderBy, timeRangeMap)
	if err != nil {
		logger.Error(err)
		return ctx.JSON(&domain.CommonResponse{
			Code: constant.ErrorCodeService,
			Msg:  "服务出现异常",
		})
	}
	return ctx.JSON(&domain.CommonResponse{Data: &domain.Paginate{
		Total:  total,
		Offset: offset,
		Limit:  limit,
		Items:  list,
	}})
}

func (c *resourceController) Get(ctx *fiber.Ctx) error {
	instance := new(domain.Resource)
	var err error
	if err = ctx.QueryParser(instance); err != nil {
		logger.Error(err)
		return ctx.JSON(&domain.CommonResponse{
			Code: constant.ErrorCodeBodyParse,
			Msg:  "参数解析失败!",
		})
	}
	instance, err = service.ResourceService.Get(instance)
	if err != nil {
		logger.Error(err)
		return ctx.JSON(&domain.CommonResponse{
			Code: constant.ErrorCodeService,
			Msg:  "服务出现异常",
		})
	}
	return ctx.JSON(&domain.CommonResponse{
		Data: instance,
	})
}
