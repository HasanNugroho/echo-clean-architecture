package service

import (
	"context"
	"fmt"
	"time"

	"github.com/HasanNugroho/golang-starter/internal/errs"
	"github.com/HasanNugroho/golang-starter/internal/helper"
	"github.com/HasanNugroho/golang-starter/internal/model"
	"github.com/HasanNugroho/golang-starter/internal/repository"
	"github.com/rs/zerolog"
)

type RoleService struct {
	repo       repository.IRoleRepository
	logger     *zerolog.Logger
	permMaster map[string]struct{}
}

func NewRoleService(repo repository.IRoleRepository, logger *zerolog.Logger) (*RoleService, error) {
	perm, err := helper.LoadStringListFromYAML("./internal/constant/data.yaml", "permission")
	if err != nil {
		logger.Error().Err(err).Msg("failed to load permission data")
		return nil, err
	}
	return &RoleService{
		repo:       repo,
		logger:     logger,
		permMaster: perm,
	}, nil
}

func (r *RoleService) Create(ctx context.Context, role *model.CreateRoleRequest) error {
	var invalid []string
	for _, p := range role.Permissions {
		if _, ok := r.permMaster[p]; !ok {
			invalid = append(invalid, p)
		}
	}
	if len(invalid) > 0 {
		return errs.BadRequest("invalid permission", fmt.Errorf("invalid permissions: %v", invalid))
	}

	payload := model.Role{
		Name:        role.Name,
		Permissions: role.Permissions,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := r.repo.Create(ctx, &payload); err != nil {
		r.logger.Error().Err(err).Fields(payload).Msg("failed to create data")
		return err
	}

	return nil
}

func (r *RoleService) FindById(ctx context.Context, id string) (*model.Role, error) {
	role, err := r.repo.FindById(ctx, id)
	if err != nil {
		r.logger.Error().Err(err).Str("roleID", id).Msg("error from repo")
		return &model.Role{}, err
	}
	return role, err
}

func (r *RoleService) FindAll(ctx context.Context, filter *model.PaginationFilter) (*[]model.Role, int64, error) {
	roles, totalItems, err := r.repo.FindAll(ctx, filter)
	if err != nil {
		r.logger.Error().Err(err).
			Str("search", filter.Search).
			Int("page", filter.Page).
			Int("limit", filter.Limit).
			Msg("error from repo")

		return &[]model.Role{}, 0, err
	}

	return roles, int64(totalItems), nil
}

func (r *RoleService) Update(ctx context.Context, id string, role *model.UpdateRoleRequest) error {
	currentRole, err := r.repo.FindById(ctx, id)
	if err != nil {
		r.logger.Error().Err(err).Str("role", id).Msg("failed to find role for update")
		return err
	}

	if role.Name != "" {
		currentRole.Name = role.Name
	}

	if role.Permissions != nil {
		var invalid []string
		for _, p := range role.Permissions {
			if _, ok := r.permMaster[p]; !ok {
				invalid = append(invalid, p)
			}
		}
		if len(invalid) > 0 {
			return errs.BadRequest("invalid permission", fmt.Errorf("invalid permissions: %v", invalid))
		}

		currentRole.Permissions = role.Permissions
	}

	return r.repo.Update(ctx, id, currentRole)
}

func (r *RoleService) Delete(ctx context.Context, id string) error {
	err := r.repo.Delete(ctx, id)
	if err != nil {
		r.logger.Error().Err(err).Str("role", id).Msg("failed to delete data")
	}
	return err
}

func (r *RoleService) AssignUser(ctx context.Context, payload *model.AssignRoleModel) error {
	err := r.repo.AssignUser(ctx, payload.UserID, payload.RoleID)
	if err != nil {
		r.logger.Error().Err(err).Fields(payload).Msg("failed to assign user")
	}
	return err
}

func (r *RoleService) UnassignUser(ctx context.Context, payload *model.AssignRoleModel) error {
	err := r.repo.UnassignUser(ctx, payload.UserID, payload.RoleID)
	if err != nil {
		r.logger.Error().Err(err).Fields(payload).Msg("failed to unassign user")
	}
	return err
}
