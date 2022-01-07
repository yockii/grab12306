package service

import (
	"errors"
	"time"

	"github.com/yockii/qscore/pkg/database"
	"github.com/yockii/qscore/pkg/domain"
	"github.com/yockii/qscore/pkg/util"
)

var DictService = new(dictService)

type dictService struct{}

func (s *dictService) Add(instance *domain.Dict) (isDuplicated bool, success bool, err error) {
	if instance.DictKey == "" {
		return false, false, errors.New("字典键名不能为空")
	}
	var c int64 = 0
	c, err = database.DB.Count(&domain.Dict{
		DictKey: instance.DictKey,
	})
	if err != nil {
		return
	}
	if c > 0 {
		isDuplicated = true
		return
	}
	instance.Id = domain.DictIdPrefix + util.GenerateDatabaseID()
	_, err = database.DB.Insert(instance)
	success = err == nil
	return
}

func (s *dictService) Remove(instance *domain.Dict) (bool, error) {
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

func (s *dictService) Update(instance *domain.Dict) (bool, error) {
	if instance.Id == "" {
		return false, errors.New("ID不能为空")
	}
	// 不允许更改的字段

	c, err := database.DB.ID(instance.Id).Update(&domain.Dict{
		// 允许更改的字段
		DictKey:   instance.DictKey,
		DictValue: instance.DictValue,
		DictExt:   instance.DictExt,
		ParentId:  instance.ParentId,
	})
	if err != nil {
		return false, err
	}
	if c == 0 {
		return false, nil
	}
	return true, nil
}

func (s *dictService) Get(instance *domain.Dict) (*domain.Dict, error) {
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

func (s *dictService) Paginate(condition *domain.Dict, limit, offset int, orderBy string) (int, []*domain.Dict, error) {
	return s.PaginateBetweenTimes(condition, limit, offset, orderBy, nil)
}

func (s *dictService) PaginateBetweenTimes(condition *domain.Dict, limit, offset int, orderBy string, tcList map[string]*domain.TimeCondition) (int, []*domain.Dict, error) {
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
	if condition.DictKey != "" {
		session.Where("dict_key like ?", condition.DictKey+"%")
		condition.DictKey = ""
	}
	if condition.DictValue != "" {
		session.Where("dict_value like ?", condition.DictValue+"%")
		condition.DictValue = ""
	}
	if condition.DictExt != "" {
		session.Where("dict_ext like ?", condition.DictExt+"%")
		condition.DictExt = ""
	}
	var list []*domain.Dict
	total, err := session.FindAndCount(&list, condition)
	if err != nil {
		return 0, nil, err
	}
	return int(total), list, nil
}
