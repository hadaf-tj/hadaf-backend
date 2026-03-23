package services

import (
	"context"
	"shb/internal/models"
)

// GetAllTeamMembers получает всех активных участников команды
func (s *Service) GetAllTeamMembers(ctx context.Context) ([]*models.TeamMember, error) {
	return s.repo.GetAllTeamMembers(ctx)
}

// GetTeamMemberByID получает участника по ID
func (s *Service) GetTeamMemberByID(ctx context.Context, id int) (*models.TeamMember, error) {
	return s.repo.GetTeamMemberByID(ctx, id)
}
