package dbrepo

import (
	"github.com/golang/mock/gomock"
	mock_database "homework-3/internal/pkg/db/mocks"
	"homework-3/internal/pkg/repository"
	"testing"
)

type DBRepoFixture struct {
	ctrl   *gomock.Controller
	repo   repository.DatabaseRepo
	mockDb *mock_database.MockDBops
}

func setUp(t *testing.T) DBRepoFixture {
	ctrl := gomock.NewController(t)
	mockDb := mock_database.NewMockDBops(ctrl)
	repo := NewPostgresRepo(mockDb)

	return DBRepoFixture{
		ctrl:   ctrl,
		repo:   repo,
		mockDb: mockDb,
	}
}

func (a *DBRepoFixture) tearDown() {
	a.ctrl.Finish()
}
