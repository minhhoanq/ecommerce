package service_test

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/minhhoanq/ecommerce/user_service/internal/dataaccess/database"
	mockdb "github.com/minhhoanq/ecommerce/user_service/internal/dataaccess/mock"
	"github.com/minhhoanq/ecommerce/user_service/internal/service"
	"github.com/minhhoanq/ecommerce/user_service/internal/util"
	"github.com/minhhoanq/ecommerce/user_service/internal/worker"
	mockwk "github.com/minhhoanq/ecommerce/user_service/internal/worker/mock"
	"github.com/stretchr/testify/require"
)

type eqCreateUserTxParamsMatcher struct {
	arg      database.CreateUserTxParams
	password string
	user     database.User
}

func (expected eqCreateUserTxParamsMatcher) Matches(x interface{}) bool {
	actualArg, ok := x.(database.CreateUserTxParams)
	if !ok {
		return false
	}

	err := util.ComparePassword(expected.password, actualArg.Password)
	if err != nil {
		return false
	}

	expected.arg.Password = actualArg.Password
	if !reflect.DeepEqual(expected.arg.CreateUserParams, actualArg.CreateUserParams) {
		return false
	}

	err = actualArg.AfterCreate(&expected.user)

	return err == nil
}

func (e eqCreateUserTxParamsMatcher) String() string {
	return fmt.Sprintf("mactches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserTxParams(arg database.CreateUserTxParams, password string, user database.User) gomock.Matcher {
	return eqCreateUserTxParamsMatcher{arg, password, user}
}

func TestGetUserByID(t *testing.T) {
	user, _ := randomUser(t)

	testCases := []struct {
		name          string
		body          uuid.UUID
		buildStubs    func(store *mockdb.MockUserDataAccessor)
		checkResponse func(t *testing.T, user *database.User, err error)
	}{
		{
			name: "OK",
			body: user.ID,
			buildStubs: func(store *mockdb.MockUserDataAccessor) {
				store.EXPECT().
					GetUserByID(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(&user, nil)
			},
			checkResponse: func(t *testing.T, res *database.User, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Equal(t, user, *res) // Dereference the actual result
			},
		},
		{
			name: "Not Found",
			body: uuid.Nil,
			buildStubs: func(store *mockdb.MockUserDataAccessor) {
				store.EXPECT().
					GetUserByID(gomock.Any(), gomock.Eq(uuid.Nil)).
					Times(1).
					Return(nil, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, res *database.User, err error) {
				require.Error(t, err)
				require.Nil(t, res)
				require.Equal(t, err, sql.ErrNoRows)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			store := mockdb.NewMockUserDataAccessor(mockCtrl)

			tc.buildStubs(store)
			service := newTestServer(t, store, nil)
			res, err := service.GetUserByID(context.Background(), tc.body)
			tc.checkResponse(t, res, err)
		})
	}
}

func TestCreateUser(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		body          service.CreateUserParams
		buildStubs    func(store *mockdb.MockUserDataAccessor, taskDistributor *mockwk.MockTaskDistributor)
		checkResponse func(t *testing.T, res *database.User, err error)
	}{
		{
			name: "OK",
			body: service.CreateUserParams{
				Username: user.Username,
				Password: password,
				Email:    user.Email,
				RoleId:   user.RoleId,
			},
			buildStubs: func(store *mockdb.MockUserDataAccessor, taskDistributor *mockwk.MockTaskDistributor) {
				arg := database.CreateUserTxParams{
					CreateUserParams: database.CreateUserParams{
						Username: user.Username,
						Email:    user.Email,
						Password: user.Password,
						RoleId:   user.RoleId,
					},
				}

				store.EXPECT().
					CreateUserTx(gomock.Any(), EqCreateUserTxParams(arg, password, user)).
					Times(1).
					Return(database.CreateUserTxResult{
						User: user,
					}, nil)
				taskPayload := &worker.PayloadSendVerifyEmail{
					UserId: user.ID,
				}
				taskDistributor.EXPECT().
					DistributeTaskSendVerifyEmail(gomock.Any(), taskPayload, gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, res *database.User, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Equal(t, user.Email, res.Email)
				require.Equal(t, user.Username, res.Username)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			store := mockdb.NewMockUserDataAccessor(mockCtrl)

			taskCtrl := gomock.NewController(t)
			defer taskCtrl.Finish()
			taskDistributor := mockwk.NewMockTaskDistributor(taskCtrl)

			tc.buildStubs(store, taskDistributor)
			service := newTestServer(t, store, taskDistributor)
			res, err := service.CreateUser(context.Background(), tc.body)

			tc.checkResponse(t, res, err)
		})
	}
}

func randomUser(t *testing.T) (user database.User, password string) {
	password = util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	return database.User{
		Username: util.RandomUsername(),
		Password: hashedPassword,
		Email:    util.RandomEmail(),
		RoleId:   1,
	}, password
}
