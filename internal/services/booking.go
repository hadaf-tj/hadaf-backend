package services

import (
	"context"
	"errors"
	"fmt"
	"shb/internal/models"
	"shb/pkg/myerrors"
)

func (s *Service) CreateBooking(ctx context.Context, userID, needID int, quantity float64, note string) (int, error) {
	// Validate need exists and is not deleted
	need, err := s.repo.GetNeedByID(ctx, needID)
	if err != nil {
		if errors.Is(err, myerrors.ErrNotFound) {
			return 0, myerrors.NewBadRequestErr("need not found")
		}
		return 0, fmt.Errorf("get need: %w", err)
	}

	// Validate user exists and is active
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, myerrors.ErrNotFound) {
			return 0, myerrors.NewBadRequestErr("user not found")
		}
		return 0, fmt.Errorf("get user: %w", err)
	}
	if !user.IsActive {
		return 0, myerrors.NewBadRequestErr("user is not active")
	}

	// Validate quantity > 0
	if quantity <= 0 {
		return 0, myerrors.NewBadRequestErr("quantity must be greater than 0")
	}

	// Check if user already has an active booking for this need
	existingBooking, err := s.repo.GetActiveBookingByUserAndNeed(ctx, userID, needID)
	if err != nil {
		return 0, fmt.Errorf("check existing booking: %w", err)
	}
	if existingBooking != nil {
		return 0, myerrors.NewConflictErr("у вас уже есть активная заявка на помощь по этой нужде")
	}

	// Create booking with status "pending"
	booking := &models.Booking{
		UserID:   userID,
		NeedID:   needID,
		Quantity: quantity,
		Note:     note,
		Status:   models.BookingStatusPending,
	}

	bookingID, err := s.repo.CreateBooking(ctx, booking)
	if err != nil {
		return 0, fmt.Errorf("create booking: %w", err)
	}

	// Fetch institution email from need's institution
	institution, err := s.repo.GetInstitutionByID(ctx, need.InstitutionID)
	if err != nil {
		s.logger.Error().Ctx(ctx).Err(err).Int("institution_id", need.InstitutionID).Msg("failed to get institution for email")
		// Don't fail the booking creation if email fetch fails
		return bookingID, nil
	}

	// Send email notification to institution director
	if institution.Email != nil && *institution.Email != "" {
		userPhone := ""
		if user.Phone != nil {
			userPhone = *user.Phone
		}
		userFullName := ""
		if user.FullName != nil {
			userFullName = *user.FullName
		}

		subject := "Новый волонтер готов помочь"
		body := fmt.Sprintf(`Учреждение: %s
Нужда: %s
Волонтер: %s
Телефон: %s
Количество: %.2f %s
Сообщение: %s

Пожалуйста, свяжитесь с волонтером для согласования.`,
			institution.Name,
			need.Name,
			userFullName,
			userPhone,
			quantity,
			need.Unit,
			note,
		)

		if err := s.email.SendEmail(ctx, *institution.Email, subject, body); err != nil {
			s.logger.Error().Ctx(ctx).Err(err).Str("email", *institution.Email).Msg("failed to send booking notification email")
			// Don't fail the booking creation if email sending fails
		}
	}

	return bookingID, nil
}

func (s *Service) ApproveBooking(ctx context.Context, bookingID, institutionUserID int) error {
	// Validate booking exists
	booking, err := s.repo.GetBookingByID(ctx, bookingID)
	if err != nil {
		return fmt.Errorf("get booking: %w", err)
	}

	// Get need to check institution
	need, err := s.repo.GetNeedByID(ctx, booking.NeedID)
	if err != nil {
		return fmt.Errorf("get need: %w", err)
	}

	// Validate requester is employee/super_admin of the institution
	requester, err := s.repo.GetUserByID(ctx, institutionUserID)
	if err != nil {
		return fmt.Errorf("get requester user: %w", err)
	}

	if requester.Role != models.RoleSuperAdmin && requester.Role != models.RoleEmployee {
		return myerrors.NewForbiddenErr("only employees and super admins can approve bookings")
	}

	if requester.Role == models.RoleEmployee {
		if requester.InstitutionID == nil || *requester.InstitutionID != need.InstitutionID {
			return myerrors.NewForbiddenErr("you can only approve bookings for your own institution")
		}
	}

	// Update status to "approved"
	err = s.repo.UpdateBookingStatus(ctx, bookingID, models.BookingStatusApproved)
	if err != nil {
		return fmt.Errorf("update booking status: %w", err)
	}

	return nil
}

func (s *Service) RejectBooking(ctx context.Context, bookingID, institutionUserID int) error {
	// Validate booking exists
	booking, err := s.repo.GetBookingByID(ctx, bookingID)
	if err != nil {
		return fmt.Errorf("get booking: %w", err)
	}

	// Get need to check institution
	need, err := s.repo.GetNeedByID(ctx, booking.NeedID)
	if err != nil {
		return fmt.Errorf("get need: %w", err)
	}

	// Validate requester is employee/super_admin of the institution
	requester, err := s.repo.GetUserByID(ctx, institutionUserID)
	if err != nil {
		return fmt.Errorf("get requester user: %w", err)
	}

	if requester.Role != models.RoleSuperAdmin && requester.Role != models.RoleEmployee {
		return myerrors.NewForbiddenErr("only employees and super admins can reject bookings")
	}

	if requester.Role == models.RoleEmployee {
		if requester.InstitutionID == nil || *requester.InstitutionID != need.InstitutionID {
			return myerrors.NewForbiddenErr("you can only reject bookings for your own institution")
		}
	}

	// Update status to "rejected"
	err = s.repo.UpdateBookingStatus(ctx, bookingID, models.BookingStatusRejected)
	if err != nil {
		return fmt.Errorf("update booking status: %w", err)
	}

	return nil
}

func (s *Service) GetBookingsByInstitution(ctx context.Context, institutionID int) ([]*models.Booking, error) {
	bookings, err := s.repo.GetBookingsByInstitution(ctx, institutionID)
	if err != nil {
		return nil, fmt.Errorf("get bookings by institution: %w", err)
	}
	return bookings, nil
}

func (s *Service) GetBookingsByUser(ctx context.Context, userID int) ([]*models.Booking, error) {
	bookings, err := s.repo.GetBookingsByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get bookings by user: %w", err)
	}
	return bookings, nil
}

func (s *Service) CompleteBooking(ctx context.Context, bookingID, institutionUserID int) error {
	// Validate booking exists
	booking, err := s.repo.GetBookingByID(ctx, bookingID)
	if err != nil {
		return fmt.Errorf("get booking: %w", err)
	}

	// Get need to check institution
	need, err := s.repo.GetNeedByID(ctx, booking.NeedID)
	if err != nil {
		return fmt.Errorf("get need: %w", err)
	}

	// Validate requester is employee/super_admin of the institution
	requester, err := s.repo.GetUserByID(ctx, institutionUserID)
	if err != nil {
		return fmt.Errorf("get requester user: %w", err)
	}

	if requester.Role != models.RoleSuperAdmin && requester.Role != models.RoleEmployee {
		return myerrors.NewForbiddenErr("only employees and super admins can complete bookings")
	}

	if requester.Role == models.RoleEmployee {
		if requester.InstitutionID == nil || *requester.InstitutionID != need.InstitutionID {
			return myerrors.NewForbiddenErr("you can only complete bookings for your own institution")
		}
	}

	// Update status to "completed"
	err = s.repo.UpdateBookingStatus(ctx, bookingID, models.BookingStatusCompleted)
	if err != nil {
		return fmt.Errorf("update booking status: %w", err)
	}

	// Increment received_qty on the need
	err = s.repo.IncrementReceivedQty(ctx, booking.NeedID, booking.Quantity)
	if err != nil {
		return fmt.Errorf("increment received qty: %w", err)
	}

	return nil
}

func (s *Service) CancelMyBooking(ctx context.Context, bookingID int, userID int) error {
	booking, err := s.repo.GetBookingByID(ctx, bookingID)
	if err != nil {
		return fmt.Errorf("get booking: %w", err)
	}
	if booking.UserID != userID {
		return myerrors.NewForbiddenErr("you can only cancel your own bookings")
	}
	if booking.Status != models.BookingStatusPending {
		return myerrors.NewBadRequestErr("only pending bookings can be cancelled")
	}
	err = s.repo.UpdateBookingStatus(ctx, bookingID, "cancelled")
	if err != nil {
		return fmt.Errorf("update booking status: %w", err)
	}
	return nil
}

func (s *Service) UpdateMyBooking(ctx context.Context, bookingID int, userID int, qty float64) error {
	booking, err := s.repo.GetBookingByID(ctx, bookingID)
	if err != nil {
		return fmt.Errorf("get booking: %w", err)
	}
	if booking.UserID != userID {
		return myerrors.NewForbiddenErr("you can only modify your own bookings")
	}
	if booking.Status != models.BookingStatusPending {
		return myerrors.NewBadRequestErr("only pending bookings can be modified")
	}
	if qty <= 0 {
		return myerrors.NewBadRequestErr("quantity must be greater than 0")
	}
	err = s.repo.UpdateBookingQuantity(ctx, bookingID, qty)
	if err != nil {
		return fmt.Errorf("update booking quantity: %w", err)
	}
	return nil
}
