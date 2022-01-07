package service

import (
	"errors"
	"time"

	"github.com/yockii/qscore/pkg/database"
	"github.com/yockii/qscore/pkg/domain"
	"github.com/yockii/qscore/pkg/util"
)

var RoleService = new(roleService)

type roleService struct{}

func (s *roleService) Add(instance *domain.Role) (isDuplicated bool, success bool, err error) {
	if instance.RoleName == "" {
		return false, false, errors.New("角色名不能为空")
	}
	var c int64 = 0
	c, err = database.DB.Count(&domain.Role{RoleName: instance.RoleName})
	if err != nil {
		return
	}
	if c > 0 {
		isDuplicated = true
		return
	}
	instance.Id = domain.RoleIdPrefix + util.GenerateDatabaseID()
	_, err = database.DB.Insert(instance)
	success = err == nil
	return
}

func (s *roleService) Remove(instance *domain.Role) (bool, error) {
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

func (s *roleService) Update(instance *domain.Role) (bool, error) {
	if instance.Id == "" {
		return false, errors.New("ID不能为空")
	}
	// 不允许更改的字段

	c, err := database.DB.ID(instance.Id).Update(&domain.Role{
		// 允许更改的字段
		RoleName: instance.RoleName,
		RoleDesc: instance.RoleDesc,
	})
	if err != nil {
		return false, err
	}
	if c == 0 {
		return false, nil
	}
	return true, nil
}

func (s *roleService) Get(instance *domain.Role) (*domain.Role, error) {
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

func (s *roleService) Paginate(condition *domain.Role, limit, offset int, orderBy string) (int, []*domain.Role, error) {
	return s.PaginateBetweenTimes(condition, limit, offset, orderBy, nil)
}

func (s *roleService) PaginateBetweenTimes(condition *domain.Role, limit, offset int, orderBy string, tcList map[string]*domain.TimeCondition) (int, []*domain.Role, error) {
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
	if condition.RoleName != "" {
		session.Where("role_name like ?", condition.RoleName+"%")
		condition.RoleName = ""
	}
	if condition.RoleDesc != "" {
		session.Where("role_desc like ?", condition.RoleDesc+"%")
		condition.RoleDesc = ""
	}
	var list []*domain.Role
	total, err := session.FindAndCount(&list, condition)
	if err != nil {
		return 0, nil, err
	}
	return int(total), list, nil
}
