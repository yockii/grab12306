package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/yockii/qscore/pkg/database"
	"github.com/yockii/qscore/pkg/domain"
	"github.com/yockii/qscore/pkg/util"
)

var ResourceService = new(resourceService)

type resourceService struct{}

func (s *resourceService) Add(instance *domain.Resource) (isDuplicated bool, success bool, err error) {
	if instance.ResourceName == "" {
		return false, false, errors.New("资源名不能为空")
	}
	if instance.ResourceContent == "" {
		return false, false, errors.New("资源内容不能为空")
	}
	if instance.ResourceType == "" {
		return false, false, errors.New("资源类型不能为空")
	}
	if instance.Action == "" {
		return false, false, errors.New("资源行为不能为空")
	}
	var c int64 = 0
	c, err = database.DB.Count(&domain.Resource{
		ResourceName:    instance.ResourceName,
		ResourceContent: instance.ResourceContent,
		ResourceType:    instance.ResourceType,
		Action:          instance.Action,
	})
	if err != nil {
		return
	}
	if c > 0 {
		isDuplicated = true
		return
	}
	instance.Id = domain.ResourceIdPrefix + util.GenerateDatabaseID()
	_, err = database.DB.Insert(instance)
	success = err == nil
	return
}

func (s *resourceService) Remove(instance *domain.Resource) (bool, error) {
	if instance.Id == "" {
		return false, errors.New("id不能为空")
	}
	c, err := database.DB.Delete(instance)
	if err != nil {
		return false, err
	}
	if c == 0 {
		return false, nil
	}
	return true, nil
}

func (s *resourceService) Update(instance *domain.Resource) (bool, error) {
	if instance.Id == "" {
		return false, errors.New("ID不能为空")
	}
	// 不允许更改的字段

	c, err := database.DB.ID(instance.Id).Update(&domain.Resource{
		// 允许更改的字段
		ResourceName:    instance.ResourceName,
		ResourceContent: instance.ResourceContent,
		ResourceType:    instance.ResourceType,
		Action:          instance.Action,
	})
	if err != nil {
		return false, err
	}
	if c == 0 {
		return false, nil
	}
	return true, nil
}

func (s *resourceService) Get(instance *domain.Resource) (*domain.Resource, error) {
	if instance.Id == "" {
		return nil, errors.New("ID不能为空")
	}
	has, err := database.DB.Get(instance)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return instance, nil
}

func (s *resourceService) Paginate(condition *domain.Resource, limit, offset int, orderBy string) (int, []*domain.Resource, error) {
	return s.PaginateBetweenTimes(condition, limit, offset, orderBy, nil)
}

func (s *resourceService) PaginateBetweenTimes(condition *domain.Resource, limit, offset int, orderBy string, tcList map[string]*domain.TimeCondition) (int, []*domain.Resource, error) {
	// 处理不允许查询的字段

	// 处理sql
	session := database.DB.NewSession()
	if limit > -1 && offset > -1 {
		session.Limit(limit, offset)
	}

	if orderBy != "" {
		session.OrderBy(orderBy)
	}
	session.Desc("create_time")

	// 处理时间字段，在某段时间之间
	for tc, tr := range tcList {
		if tc != "" {
			if !tr.Start.IsZero() && !tr.End.IsZero() {
				session.Where(tc+" between ? and ?", time.Time(tr.Start), time.Time(tr.End))
			} else if tr.Start.IsZero() && !tr.End.IsZero() {
				session.Where(tc+" <= ?", time.Time(tr.End))
			} else if !tr.Start.IsZero() && tr.End.IsZero() {
				session.Where(tc+" > ?", time.Time(tr.Start))
			}
		}
	}

	// 模糊查找
	if condition.ResourceName != "" {
		session.Where("resource_name like ?", condition.ResourceName+"%")
		condition.ResourceName = ""
	}
	if condition.ResourceContent == "" {
		session.Where("resource_content like ?", condition.ResourceContent+"%")
		condition.ResourceContent = ""
	}
	if condition.ResourceType == "" {
		session.Where("resource_type like ?", condition.ResourceType+"%")
		condition.ResourceType = ""
	}
	if condition.Action == "" {
		session.Where("action like ?", condition.Action+"%")
		condition.Action = ""
	}
	var list []*domain.Resource
	total, err := session.FindAndCount(&list, condition)
	if err != nil {
		return 0, nil, err
	}
	return int(total), list, nil
}

func (s *resourceService) GetListByIds(ids []string) (rs []*domain.Resource, err error) {
	if len(ids) == 0 {
		return
	}
	err = database.DB.Where(fmt.Sprintf("id in (%s)", "'"+strings.Join(ids, "','")+"'")).Find(&rs)
	return
}
