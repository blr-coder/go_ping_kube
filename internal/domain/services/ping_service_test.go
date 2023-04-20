package services

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/suite"
	"go_ping_kube/internal/domain/errs"
	"go_ping_kube/internal/domain/models"
	mock_repository "go_ping_kube/internal/infrastructure/repository/mock"
	"testing"
	"time"
)

func NewMockRepoData(t *testing.T) *models.PingData {
	t.Helper()
	f := faker.New()
	return &models.PingData{
		Id:         uuid.New(),
		CreatedAt:  f.Time().Time(time.Now()).UTC(),
		IP:         f.Internet().Ipv4(),
		RequestURI: f.Internet().URL(),
		UserAgent:  f.UserAgent().UserAgent(),
		Headers:    nil,
		Client:     "",
		Device: f.RandomStringElement([]string{
			"Mobile", "Tablet", "Desktop", "Bot",
		}),
		OS:      "",
		Lang:    []string{f.Language().Language()},
		Country: f.Address().Country(),
	}
}

type PingServiceTestsSuite struct {
	suite.Suite
	mockController          *gomock.Controller
	mockMockIPingRepository *mock_repository.MockIPingRepository
	pingService             *PingService
}

func (ts *PingServiceTestsSuite) SetupTest() {
	ts.mockController = gomock.NewController(ts.T())
	ts.mockMockIPingRepository = mock_repository.NewMockIPingRepository(ts.mockController)
	ts.pingService = &PingService{
		storage: ts.mockMockIPingRepository,
	}
}

func (ts *PingServiceTestsSuite) clear() {
	ts.mockController.Finish()
}

func TestPing(t *testing.T) {
	suite.Run(t, new(PingServiceTestsSuite))
}

func (ts *PingServiceTestsSuite) TestGetOK() {
	defer ts.clear()
	ctx := context.Background()

	repoData := NewMockRepoData(ts.T())

	ts.mockMockIPingRepository.
		EXPECT().
		Get(ctx, repoData.Id).
		Return(repoData, nil).
		Times(1)

	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	tests := []struct {
		name    string
		args    args
		want    *models.PingData
		wantErr error
	}{
		{
			name: "OK",
			args: args{
				ctx: ctx,
				id:  repoData.Id,
			},
			want:    repoData,
			wantErr: nil,
		},
	}

	for _, tc := range tests {
		ts.Run(tc.name, func() {
			actual, err := ts.pingService.Get(tc.args.ctx, tc.args.id)
			ts.Require().Equal(tc.want, actual)
			ts.Require().Equal(tc.wantErr, err)
		})
	}
}

func (ts *PingServiceTestsSuite) TestGetError() {
	defer ts.clear()
	ctx := context.Background()

	randUuid := uuid.New()

	ts.mockMockIPingRepository.
		EXPECT().
		Get(ctx, randUuid).
		Return(nil, errs.NewDomainNotFoundError(randUuid.String())).
		Times(1)

	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	tests := []struct {
		name    string
		args    args
		want    *models.PingData
		wantErr error
	}{
		{
			name: "RECORD NOT FOUND",
			args: args{
				ctx: ctx,
				id:  randUuid,
			},
			want:    nil,
			wantErr: errs.NewDomainNotFoundError(randUuid.String()),
		},
	}

	for _, tc := range tests {
		ts.Run(tc.name, func() {
			actual, err := ts.pingService.Get(tc.args.ctx, tc.args.id)
			ts.Require().Equal(tc.want, actual)
			ts.Require().Equal(tc.wantErr, err)
		})
	}
}
