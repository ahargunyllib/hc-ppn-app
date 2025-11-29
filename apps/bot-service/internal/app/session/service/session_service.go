package service

import (
	"context"
	"time"

	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/entity"
	"github.com/ahargunyllib/hc-ppn-app/apps/bot-service/domain/errx"
	"github.com/google/uuid"
)

func (s *SessionService) GetOrCreateSession(ctx context.Context, phoneNumber string) (*entity.ConversationSession, error) {
	session, err := s.sessionRepo.FindActiveByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		return nil, err
	}

	if session != nil {
		return session, nil
	}

	id, err := s.uuidPkg.NewV7()
	if err != nil {
		return nil, errx.ErrInternalServer.WithLocation("SessionService.GetOrCreateSession").WithError(err)
	}

	newSession := &entity.ConversationSession{
		ID:            id,
		PhoneNumber:   phoneNumber,
		Status:        entity.SessionStatusActive,
		LastMessageAt: time.Now(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.sessionRepo.Create(ctx, newSession); err != nil {
		return nil, err
	}

	return newSession, nil
}

func (s *SessionService) UpdateSessionActivity(ctx context.Context, sessionID uuid.UUID) error {
	session, err := s.sessionRepo.FindByID(ctx, sessionID)
	if err != nil {
		return err
	}

	session.LastMessageAt = time.Now()
	session.UpdatedAt = time.Now()

	if err := s.sessionRepo.Update(ctx, session); err != nil {
		return err
	}

	return nil
}

func (s *SessionService) MarkSessionWaitingFeedback(ctx context.Context, sessionID uuid.UUID) error {
	session, err := s.sessionRepo.FindByID(ctx, sessionID)
	if err != nil {
		return err
	}

	now := time.Now()
	session.Status = entity.SessionStatusWaitingFeedback
	session.FeedbackPromptSentAt = &now
	session.UpdatedAt = now

	if err := s.sessionRepo.Update(ctx, session); err != nil {
		return err
	}

	return nil
}

func (s *SessionService) CloseSession(ctx context.Context, sessionID uuid.UUID) error {
	session, err := s.sessionRepo.FindByID(ctx, sessionID)
	if err != nil {
		return err
	}

	session.Status = entity.SessionStatusClosed
	session.UpdatedAt = time.Now()

	if err := s.sessionRepo.Update(ctx, session); err != nil {
		return err
	}

	return nil
}

func (s *SessionService) ProcessExpiredSessions(ctx context.Context) error {
	expiredSessions, err := s.sessionRepo.FindExpiredSessions(ctx, 5)
	if err != nil {
		return err
	}

	for i := range expiredSessions {
		session := &expiredSessions[i]

		now := time.Now()
		session.Status = entity.SessionStatusWaitingFeedback
		session.FeedbackPromptSentAt = &now
		session.UpdatedAt = now

		if err := s.sessionRepo.Update(ctx, session); err != nil {
			s.logger.Error().Err(err).Str("session_id", session.ID.String()).Msg("Failed to mark session as waiting for feedback")
			continue
		}

		s.logger.Info().Str("session_id", session.ID.String()).Str("phone_number", session.PhoneNumber).Msg("Session marked as waiting for feedback")
	}

	waitingSessions, err := s.sessionRepo.List(ctx, &entity.GetSessionsFilter{
		Offset: 0,
		Limit:  100,
		Status: func() *entity.ConversationSessionStatus { status := entity.SessionStatusWaitingFeedback; return &status }(),
	})
	if err != nil {
		return err
	}

	for i := range waitingSessions {
		session := &waitingSessions[i]
		if session.FeedbackPromptSentAt == nil {
			continue
		}

		if time.Since(*session.FeedbackPromptSentAt) > 5*time.Minute {
			session.Status = entity.SessionStatusClosed
			session.UpdatedAt = time.Now()

			if err := s.sessionRepo.Update(ctx, session); err != nil {
				s.logger.Error().Err(err).Str("session_id", session.ID.String()).Msg("Failed to close session")
				continue
			}

			s.logger.Info().Str("session_id", session.ID.String()).Str("phone_number", session.PhoneNumber).Msg("Session auto-closed due to no feedback")
		}
	}

	return nil
}
