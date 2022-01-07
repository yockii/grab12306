package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yockii/qscore/pkg/constant"
	"github.com/yockii/qscore/pkg/domain"
	"github.com/yockii/qscore/pkg/logger"

	"github.com/yockii/grab12306/internal/service"
)

var DictController = new(dictController)

type dictController struct{}

func (c *dictController) Add(ctx *fiber.Ctx) error {
	instance := new(domain.Dict)
	if err := ctx.BodyParser(instance); err != nil {
		logger.Error(err)
		return ctx.JSON(&domain.CommonResponse{
			Code: constant.ErrorCodeBodyParse,
			Msg:  "参数解析失败!",
		})
	}

	// 处理必填
	if instance.DictKey == "" {
		return ctx.JSON(&domain.CommonResponse{
			Code: constant.ErrorCodeLackOfField,
			Msg:  "字典键必须提供",
		})
	}

	duplicated, success, err := service.DictService.Add(instance)
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

func (c *dictController) Delete(ctx *fiber.Ctx) error {
	instance := new(domain.Dict)
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
	deleted, err := service.DictService.Remove(instance)
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

func (c *dictController) Update(ctx *fiber.Ctx) error {
	instance := new(domain.Dict)
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
	updated, err := service.DictService.Update(instance)
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

func (c *dictController) Paginate(ctx *fiber.Ctx) error {
	pr := new(domain.DictRequest)
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

	total, list, err := service.DictService.PaginateBetweenTimes(&pr.Dict, limit, offset, orderBy, timeRangeMap)
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

func (c *dictController) Get(ctx *fiber.Ctx) error {
	instance := new(domain.Dict)
	var err error
	if err = ctx.QueryParser(instance); err != nil {
		logger.Error(err)
		return ctx.JSON(&domain.CommonResponse{
			Code: constant.ErrorCodeBodyParse,
			Msg:  "参数解析失败!",
		})
	}
	instance, err = service.DictService.Get(instance)
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
