package core

type RoleRepository interface {
	RemoveMember(roleID uint64, userID uint64) error
	AddMember(roleID uint64, userID uint64) error
}

type RoleService interface {
	RemoveMember(roleID, userID uint64) error
	AddMember(roleID uint64, userID uint64) error
}

type roleService struct {
	repo RoleRepository
}

func NewRoleService(re RoleRepository) RoleService {
	return &roleService{
		repo: re,
	}
}

func (s *roleService) RemoveMember(roleID uint64, userID uint64) error {
	return s.repo.RemoveMember(roleID, userID)
}

func (s *roleService) AddMember(roleID uint64, userID uint64) error {
	return s.repo.AddMember(roleID, userID)
}
