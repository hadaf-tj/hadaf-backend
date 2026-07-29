// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

package services_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"shb/internal/models"
	"shb/internal/repositories/filters"
	"shb/pkg/external/sms/smsProvider"
	"shb/pkg/myerrors"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_GetAllInstitutions_GetInstitutionByID_CreateInstitution(t *testing.T) {
	ctx := context.Background()
	svc, d := newTestService(t)

	q := models.InstitutionListQuery{Limit: 10, Offset: 0}
	page := &models.InstitutionPage{Total: 1, Limit: 10}
	d.Repo.On("GetAllInstitutions", ctx, q).Return(page, nil)
	out, err := svc.GetAllInstitutions(ctx, q)
	require.NoError(t, err)
	require.Equal(t, page, out)

	inst := &models.Institution{ID: 1, Name: "X"}
	d.Repo.On("GetInstitutionByID", ctx, 1).Return(inst, nil)
	g, err := svc.GetInstitutionByID(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, inst, g)

	d.Repo.On("CreateInstitution", ctx, inst).Return(2, nil)
	id, err := svc.CreateInstitution(ctx, inst)
	require.NoError(t, err)
	require.Equal(t, 2, id)
}

func TestService_GetAllEvents_GetEventByID_GetEventDetail(t *testing.T) {
	ctx := context.Background()
	svc, d := newTestService(t)

	eq := models.EventListQuery{UserID: 0, Limit: 20, Offset: 0}
	ep := &models.EventPage{Total: 0, Limit: 20}
	d.Repo.On("GetAllEvents", ctx, eq).Return(ep, nil)
	out, err := svc.GetAllEvents(ctx, eq)
	require.NoError(t, err)
	require.Equal(t, ep, out)

	ev := &models.Event{ID: 1}
	d.Repo.On("GetEventByID", ctx, 1).Return(ev, nil)
	g, err := svc.GetEventByID(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, ev, g)

	detail := models.EventDetailQuery{EventID: 1, ViewerUserID: 0}
	er := &models.EventResponse{ID: 1}
	d.Repo.On("GetEventDetail", ctx, detail).Return(er, nil)
	r, err := svc.GetEventDetail(ctx, detail)
	require.NoError(t, err)
	require.Equal(t, er, r)
}

func TestService_GetEventDetail_notFound(t *testing.T) {
	ctx := context.Background()
	svc, d := newTestService(t)
	q := models.EventDetailQuery{EventID: 99, ViewerUserID: 0}
	d.Repo.On("GetEventDetail", ctx, q).Return(nil, pgx.ErrNoRows)
	_, err := svc.GetEventDetail(ctx, q)
	require.Error(t, err)
	require.ErrorIs(t, err, myerrors.ErrNotFound)
}

func TestService_GetPublicStats(t *testing.T) {
	ctx := context.Background()
	svc, d := newTestService(t)
	stats := map[string]int{"users": 5}
	d.Repo.On("GetPublicStats", ctx).Return(stats, nil)
	out, err := svc.GetPublicStats(ctx)
	require.NoError(t, err)
	require.Equal(t, stats, out)
}

func TestService_GetPublicStats_error(t *testing.T) {
	ctx := context.Background()
	svc, d := newTestService(t)
	d.Repo.On("GetPublicStats", ctx).Return(nil, errors.New("db"))
	_, err := svc.GetPublicStats(ctx)
	require.Error(t, err)
}

func TestService_CheckSMSBalance(t *testing.T) {
	ctx := context.Background()
	svc, d := newTestService(t)
	bal := &smsProvider.BalanceResult{}
	d.SMS.On("CheckBalance", ctx).Return(bal, nil)
	res, err := svc.CheckSMSBalance(ctx)
	require.NoError(t, err)
	require.Equal(t, bal, res)
}

func TestService_NeedsCRUD(t *testing.T) {
	ctx := context.Background()
	svc, d := newTestService(t)

	n := &models.Need{InstitutionID: 1, Name: "n", Unit: "kg", RequiredQty: 1, Urgency: "low"}
	d.Repo.On("CreateNeed", ctx, n).Return(10, nil)
	id, err := svc.CreateNeed(ctx, n)
	require.NoError(t, err)
	require.Equal(t, 10, id)

	old := &models.Need{ID: 1, InstitutionID: 1, Name: "n", Unit: "u", RequiredQty: 1, ReceivedQty: 0, Urgency: "low"}
	upd := *old
	upd.Name = "new"
	d.Repo.On("GetNeedByID", ctx, 1).Return(old, nil)
	d.Repo.On("UpdateNeed", ctx, mock.MatchedBy(func(x *models.Need) bool { return x.Name == "new" })).Return(nil)
	d.Repo.On("CreateNeedHistory", ctx, mock.AnythingOfType("*models.NeedsHistory")).Return(nil)
	require.NoError(t, svc.UpdateNeed(ctx, &upd))

	d.Repo.On("DeleteNeed", ctx, 2).Return(nil)
	d.Repo.On("CreateNeedHistory", ctx, mock.AnythingOfType("*models.NeedsHistory")).Return(nil)
	require.NoError(t, svc.DeleteNeed(ctx, 2))

	f := filters.NeedsFilter{}
	d.Repo.On("GetNeedsByInstitution", ctx, f, 3).Return([]*models.Need{}, nil)
	list, err := svc.GetNeedsByInstitution(ctx, f, 3)
	require.NoError(t, err)
	require.Empty(t, list)

	d.Repo.On("GetNeedByID", ctx, 5).Return(old, nil)
	got, err := svc.GetNeedByID(ctx, 5)
	require.NoError(t, err)
	require.Equal(t, old, got)
}

func TestService_CreateEvent_JoinLeave(t *testing.T) {
	ctx := context.Background()

	t.Run("past date rejected", func(t *testing.T) {
		svc, _ := newTestService(t)
		_, err := svc.CreateEvent(ctx, &models.Event{EventDate: time.Now().Add(-time.Hour), InstitutionID: 1})
		require.Error(t, err)
		var br myerrors.BadRequestErr
		require.ErrorAs(t, err, &br)
	})

	t.Run("create success", func(t *testing.T) {
		svc, d := newTestService(t)
		future := time.Now().Add(48 * time.Hour)
		d.Repo.On("GetInstitutionByID", ctx, 5).Return(&models.Institution{ID: 5}, nil)
		d.Repo.On("CreateEvent", ctx, mock.MatchedBy(func(e *models.Event) bool {
			return e.InstitutionID == 5 && e.Title == "T"
		})).Return(100, nil)
		id, err := svc.CreateEvent(ctx, &models.Event{
			Title: "T", EventDate: future, InstitutionID: 5, CreatorID: 1,
		})
		require.NoError(t, err)
		require.Equal(t, 100, id)
	})

	t.Run("join past event", func(t *testing.T) {
		svc, d := newTestService(t)
		d.Repo.On("GetEventByID", ctx, 1).Return(&models.Event{
			ID: 1, EventDate: time.Now().Add(-time.Hour),
		}, nil)
		err := svc.JoinEvent(ctx, 1, 2)
		require.Error(t, err)
		var br myerrors.BadRequestErr
		require.ErrorAs(t, err, &br)
	})

	t.Run("leave calls repo", func(t *testing.T) {
		svc, d := newTestService(t)
		d.Repo.On("GetEventByID", ctx, 1).Return(&models.Event{ID: 1, EventDate: time.Now().Add(time.Hour)}, nil)
		d.Repo.On("LeaveEvent", ctx, 1, 2).Return(nil)
		require.NoError(t, svc.LeaveEvent(ctx, 1, 2))
	})

	t.Run("join future event", func(t *testing.T) {
		svc, d := newTestService(t)
		d.Repo.On("GetEventByID", ctx, 3).Return(&models.Event{ID: 3, EventDate: time.Now().Add(48 * time.Hour)}, nil)
		d.Repo.On("JoinEvent", ctx, 3, 9).Return(nil)
		require.NoError(t, svc.JoinEvent(ctx, 3, 9))
	})
}
