package db

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"service-chat/internal/db/entity"
)

func TestAuthPostgres_CreateUser(t *testing.T) {
	// Создаём мок объекта базы данных
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() { _ = db.Close() }()

	r := NewAuthPostgres(db)

	// Входные данные для метода CreateUser
	type args struct {
		user entity.User
	}

	// Функция для определения поведения мока базы данных
	type mockBehavior func(args args, id int)

	tests := []struct {
		name    string
		input   args
		mock    mockBehavior
		wantID  int
		wantErr error
	}{
		{
			name: "Success",
			input: args{
				user: entity.User{
					Username: "Andrey",
					Password: "CiRA9gEG",
				},
			},
			mock: func(input args, wantID int) {
				// Мок sql запроса
				rowMock := sqlmock.NewRows([]string{"id"}).AddRow(wantID)
				mock.
					ExpectPrepare(`INSERT INTO "user" (username, password_hash) VALUES ($1, $2) RETURNING id`).
					ExpectQuery().WithArgs(input.user.Username, input.user.Password).WillReturnRows(rowMock)
			},
			wantID: 1,
		},
		{
			name: "Empty username",
			input: args{
				user: entity.User{
					Username: "",
					Password: "CiRA9gEG",
				},
			},
			mock:    func(input args, wantID int) {},
			wantID:  0,
			wantErr: errors.New("error path: db.CreateUser, error: empty username or password"),
		},
		{
			name: "Empty password",
			input: args{
				user: entity.User{
					Username: "Andrey",
					Password: "",
				},
			},
			mock:    func(input args, wantID int) {},
			wantID:  0,
			wantErr: errors.New("error path: db.CreateUser, error: empty username or password"),
		},
		{
			name:    "Nil request",
			mock:    func(input args, wantID int) {},
			wantID:  0,
			wantErr: errors.New("error path: db.CreateUser, error: empty username or password"),
		},
		{
			name: "Unique violation",
			input: args{
				user: entity.User{
					Username: "Andrey",
					Password: "CiRA9gEG",
				},
			},
			mock: func(input args, wantID int) {
				// Мок sql запроса
				rowMock := sqlmock.NewRows([]string{"id"}).AddRow(wantID).RowError(0, errors.New("unique_violation"))
				mock.
					ExpectPrepare(`INSERT INTO "user" (username, password_hash) VALUES ($1, $2) RETURNING id`).
					ExpectQuery().WithArgs(input.user.Username, input.user.Password).WillReturnRows(rowMock)
			},
			wantID:  0,
			wantErr: errors.New("error path: db.CreateUser, error: unique_violation"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.input, tt.wantID)
			acID, acErr := r.CreateUser(tt.input.user)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr.Error(), acErr.Error())
				assert.Equal(t, tt.wantID, acID)
			} else {
				assert.NoError(t, acErr)
				assert.Equal(t, tt.wantID, acID)
			}
			// Проверяем все ли моки выполнены
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAuthPostgres_GetUser(t *testing.T) {
	// Создаём мок объекта базы данных
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() { _ = db.Close() }()

	r := NewAuthPostgres(db)

	// Входные данные для метода CreateUser
	type args struct {
		user entity.User
	}

	// Функция для определения поведения мока базы данных
	type mockBehavior func(args args)

	tests := []struct {
		name     string
		input    args
		mock     mockBehavior
		wantUser *entity.User
		wantErr  error
	}{
		{
			name: "Success",
			input: args{
				user: entity.User{
					Username: "Andrey",
					Password: "adgui*",
				},
			},
			mock: func(input args) {
				// Мок sql запроса
				rowMock := sqlmock.NewRows([]string{"id", "username", "password_hash"}).
					AddRow(1, "Andrey", "CiRA9gEG")
				mock.
					ExpectPrepare(`SELECT id, username, password_hash FROM "user" WHERE username = $1`).
					ExpectQuery().WithArgs(input.user.Username).WillReturnRows(rowMock)
			},
			wantUser: &entity.User{
				Id:       1,
				Username: "Andrey",
				Password: "CiRA9gEG",
			},
		},
		{
			name: "Empty username",
			input: args{
				user: entity.User{
					Username: "",
					Password: "adgui*",
				},
			},
			mock:     func(input args) {},
			wantUser: nil,
			wantErr:  errors.New("error path: db.GetUser, error: empty username"),
		},
		{
			name: "Other error",
			input: args{
				user: entity.User{
					Username: "Andrey",
					Password: "adgui*",
				},
			},
			mock: func(input args) {
				// Мок sql запроса
				rowMock := sqlmock.NewRows([]string{"id", "username", "password_hash"}).
					AddRow(1, "Andrey", "CiRA9gEG").RowError(0, errors.New("other error"))
				mock.
					ExpectPrepare(`SELECT id, username, password_hash FROM "user" WHERE username = $1`).
					ExpectQuery().WithArgs(input.user.Username).WillReturnRows(rowMock)
			},
			wantUser: nil,
			wantErr:  errors.New("error path: db.GetUser, error: other error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.input)
			acUser, acErr := r.GetUser(tt.input.user)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr.Error(), acErr.Error())
				assert.Equal(t, tt.wantUser, acUser)
			} else {
				assert.NoError(t, acErr)
				assert.Equal(t, tt.wantUser, acUser)
			}
			// Проверяем все ли моки выполнены
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

//stmt, errS := db.Prepare(`INSERT INTO "user" (username, password_hash) VALUES ($1, $2) RETURNING id`)
//
//// Если ошибка подготовки скелета запроса
//if errS != nil {
//	t.Fatalf("error '%s' was not expected while creating a prepared statement", errS)
//}
//
//// Если подготовленный запрос nil
//if stmt == nil {
//	t.Fatalf("stmt is nil")
//}
//
//// Получаем запись из мокированной базы данных
//row := stmt.QueryRow(input.user.Username, input.user.Password)
//
//// Получаем id записи
//var id int
//if err = row.Scan(&id); err != nil {
//	t.Fatalf("error '%s' was not expected while scanning the row", err)
//}
//
//// Проверяем, что все моки отработали
//if errE := mock.ExpectationsWereMet(); errE != nil {
//	t.Fatalf("there were unfulfilled expectations: %s", errE)
//}

// новое

// Мок скелета sql запроса
//stmt, errS := db.Prepare(`INSERT INTO "user" (username, password_hash) VALUES ($1, $2) RETURNING id`)
//
//// Если ошибка подготовки скелета запроса
//if errS != nil {
//	t.Fatalf("error '%s' was not expected while creating a prepared statement", errS)
//}
//
//// Если подготовленный запрос nil
//if stmt == nil {
//	t.Fatalf("stmt is nil")
//}

// Мок непосредственно самого запроса
//rowMock := sqlmock.NewRows([]string{"id"}).AddRow(wantID)
//mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "user" (username, password_hash) VALUES ($1, $2) RETURNING id`)).
//	WithArgs(input.user.Username, input.user.Password).WillReturnRows(rowMock)
//row := stmt.QueryRow(input.user.Username, input.user.Password)

// Получаем id записи
//var id int
//if err = row.Scan(&id); err != nil {
//	t.Fatalf("error '%s' was not expected while scanning the row", err)
//}
