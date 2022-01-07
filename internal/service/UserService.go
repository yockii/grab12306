package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/yockii/qscore/pkg/cache"

	"github.com/yockii/qscore/pkg/constant"
	"github.com/yockii/qscore/pkg/database"
	"github.com/yockii/qscore/pkg/domain"
	"github.com/yockii/qscore/pkg/util"
	"golang.org/x/crypto/bcrypt"
)

var UserService = new(userService)

type userService struct{}

func (s *userService) Add(instance *domain.User) (isDuplicated bool, success bool, err error) {
	if instance.Username == "" {
		return false, false, errors.New("用户名不能为空")
	}
	var c int64 = 0
	c, err = database.DB.Count(&domain.User{Username: instance.Username})
	if err != nil {
		return
	}
	if c > 0 {
		isDuplicated = true
		return
	}
	instance.Id = domain.UserIdPrefix + util.GenerateDatabaseID()
	pwd, _ := bcrypt.GenerateFromPassword([]byte(instance.Password), bcrypt.DefaultCost)
	instance.Password = string(pwd)
	_, err = database.DB.Insert(instance)
	success = err == nil
	instance.Password = ""
	return
}

func (s *userService) Remove(instance *domain.User) (bool, error) {
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

func (s *userService) Update(instance *domain.User) (bool, error) {
	if instance.Id == "" {
		return false, errors.New("ID不能为空")
	}
	// 不允许更改的字段
	if instance.Username != "" {
		instance.Username = ""
	}
	if instance.Password != "" {
		instance.Password = ""
	}

	c, err := database.DB.ID(instance.Id).Update(&domain.User{
		// 允许更改的字段
	})
	if err != nil {
		return false, err
	}
	if c == 0 {
		return false, nil
	}
	return true, nil
}

func (s *userService) Get(instance *domain.User) (*domain.User, error) {
	//if instance.Id == "" {
	//	return nil, errors.New("ID不能为空")
	//}
	has, err := database.DB.Omit("password").Get(instance)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return instance, nil
}

func (s *userService) Paginate(condition *domain.User, limit, offset int, orderBy string) (int, []*domain.User, error) {
	return s.PaginateBetweenTimes(condition, limit, offset, orderBy, nil)
}

func (s *userService) PaginateBetweenTimes(condition *domain.User, limit, offset int, orderBy string, tcList map[string]*domain.TimeCondition) (int, []*domain.User, error) {
	// 处理不允许查询的字段
	if condition.Password != "" {
		condition.Password = ""
	}
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
	if condition.Username != "" {
		session.Where("username like ?", condition.Username+"%")
		condition.Username = ""
	}
	var list []*domain.User
	total, err := session.Omit("password").FindAndCount(&list, condition)
	if err != nil {
		return 0, nil, err
	}
	return int(total), list, nil
}

func (s *userService) Login(instance *domain.User) (*domain.User, string, error) {
	if instance.Username == "" {
		return nil, "", errors.New("用户名不能为空")
	}
	u := new(domain.User)
	u.Username = instance.Username
	has, err := database.DB.Get(u)
	if err != nil {
		return nil, "", err
	}
	if has {
		err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(instance.Password))
		if err != nil {
			return nil, "", nil
		}
		token := ""
		token, err = generateToken(u.Id, u.Username, "", 72*3600)
		u.Password = ""
		return u, token, nil
	}
	return nil, "", nil
}

func generateToken(userId string, username, tenantId string, expireInSecond int) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	sid := util.GenerateDatabaseID()

	if cache.Enabled() {
		rConn := cache.Get()
		defer rConn.Close()
		_, err := rConn.Do("SETEX", cache.Prefix+":"+constant.AppSid+":"+sid, expireInSecond, userId)
		if err != nil {
			return "", err
		}
	}

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["uid"] = userId
	claims["sid"] = sid
	claims["tenantId"] = tenantId
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(constant.JWT_SECRET))
	return t, err
}

func (s *userService) UpdatePassword(instance *domain.User, newPassword string) (bool, error) {
	if instance.Id == "" {
		return false, errors.New("ID不能为空")
	}

	u := new(domain.User)
	if exist, err := database.DB.ID(instance.Id).Get(u); err != nil {
		return false, err
	} else if !exist {
		return false, nil
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(instance.Password)); err != nil {
		return false, nil
	}

	npwd, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return false, err
	}

	c, err := database.DB.ID(instance.Id).Update(&domain.User{Password: string(npwd)})
	if err != nil {
		return false, err
	}
	return c > 0, nil
}
