package domain

import "github.com/yockii/qscore/pkg/domain"

type UserPasswordRequest struct {
	domain.User
	NewPassword string `json:"newPassword,omitempty"`
}

type UserRolesRequest struct {
	UserId  string   `json:"userId,omitempty"`
	RoleIds []string `json:"roleIds,omitempty"`
}

type RoleResourcesRequest struct {
	RoleId      string   `json:"roleId,omitempty"`
	ResourceIds []string `json:"resourceIds,omitempty"`
}
