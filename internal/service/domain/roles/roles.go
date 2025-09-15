package rolesService

import (
	"context"

	usersEntity "github.com/orangeMangoDimz/go-social/internal/entities/users"
	"github.com/orangeMangoDimz/go-social/internal/storage"
)

type RoleService struct {
	roleRepository storage.RolesRepository
	// validation
}

func NewRoleService(roleRepository storage.RolesRepository) *RoleService {
	return &RoleService{
		roleRepository: roleRepository,
		// validation
	}
}

func (s *RoleService) GetByName(ctx context.Context, roleName string) (*usersEntity.Role, error) {
	role, err := s.roleRepository.GetByName(ctx, roleName)
	return role, err
}
