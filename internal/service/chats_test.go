package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"service-chat/internal/db"
	"service-chat/internal/db/entity"
	mockRepo "service-chat/internal/db/mocks"
	"service-chat/internal/dto"
)

func TestChatService_CreateChat(t *testing.T) {
	// Структура для последующей реализации поведения мока
	type mockBehaviour func(s *mockRepo.MockChat, dataDB entity.ChatAdd)

	tests := []struct {
		name    string
		inChat  dto.ChatAdd
		dataDB  entity.ChatAdd
		mock    mockBehaviour
		wantID  int
		wantErr error
	}{
		{
			name: "Success two user",
			inChat: dto.ChatAdd{
				ChatName: "chat_1",
				Users:    []int64{1, 2},
			},
			dataDB: entity.ChatAdd{
				ChatName: "chat_1",
				Users:    []int64{1, 2},
			},
			mock: func(s *mockRepo.MockChat, dataDB entity.ChatAdd) {
				s.EXPECT().CreateChat(dataDB).Return(1, nil)
			},
			wantID:  1,
			wantErr: nil,
		},
		{
			name: "Success many user",
			inChat: dto.ChatAdd{
				ChatName: "chat_1",
				Users:    []int64{1, 2, 3, 4},
			},
			dataDB: entity.ChatAdd{
				ChatName: "chat_1",
				Users:    []int64{1, 2, 3, 4},
			},
			mock: func(s *mockRepo.MockChat, dataDB entity.ChatAdd) {
				s.EXPECT().CreateChat(dataDB).Return(1, nil)
			},
			wantID:  1,
			wantErr: nil,
		},
		{
			name: "Empty request",
			inChat: dto.ChatAdd{
				ChatName: "",
				Users:    []int64{},
			},
			dataDB: entity.ChatAdd{
				ChatName: "",
				Users:    []int64{},
			},
			mock:    func(s *mockRepo.MockChat, dataDB entity.ChatAdd) {},
			wantID:  0,
			wantErr: errors.New("chat_name or users is empty"),
		},
		{
			name: "Empty chat_name",
			inChat: dto.ChatAdd{
				ChatName: "",
				Users:    []int64{1, 2},
			},
			dataDB: entity.ChatAdd{
				ChatName: "",
				Users:    []int64{1, 2},
			},
			mock:    func(s *mockRepo.MockChat, dataDB entity.ChatAdd) {},
			wantID:  0,
			wantErr: errors.New("chat_name or users is empty"),
		},
		{
			name:    "Nil request",
			mock:    func(s *mockRepo.MockChat, dataDB entity.ChatAdd) {},
			wantID:  0,
			wantErr: errors.New("chat_name or users is empty"),
		},
		{
			name: "Other error",
			inChat: dto.ChatAdd{
				ChatName: "chat_1",
				Users:    []int64{1, 2},
			},
			dataDB: entity.ChatAdd{
				ChatName: "chat_1",
				Users:    []int64{1, 2},
			},
			mock: func(s *mockRepo.MockChat, dataDB entity.ChatAdd) {
				s.EXPECT().CreateChat(dataDB).Return(0, errors.New("some error"))
			},
			wantID:  0,
			wantErr: errors.New("some error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//Инициализируем контролер для мока сервиса
			ctrl := gomock.NewController(t)
			// Завершаем работу контролера после выполнения каждого теста
			defer ctrl.Finish()

			// Создаём моки базы данных авторизации
			mockChat := mockRepo.NewMockChat(ctrl)

			// Передаём структуру пользователя
			tt.mock(mockChat, tt.dataDB)

			// Создаём объект базы данных в который передадим наш мок авторизации
			repository := &db.DB{Chat: mockChat}

			// Создаём экземпляр сервиса авторизации
			serviceChat := NewChatService(repository)

			// Проверяем ожидаемый и актуальный результат
			acID, acErr := serviceChat.CreateChat(tt.inChat)
			assert.Equal(t, tt.wantID, acID)
			assert.Equal(t, tt.wantErr, acErr)
		})
	}
}

func TestChatService_GetChat(t *testing.T) {
	// Структура для последующей реализации поведения мока
	type mockBehaviour func(s *mockRepo.MockChat, dataDB entity.ChatGet)
	userID := int64(1)
	userEmpty := int64(0)

	tests := []struct {
		name    string
		inChat  dto.ChatGet
		dataDB  entity.ChatGet
		mock    mockBehaviour
		want    []entity.Chat
		wantErr error
	}{
		{
			name: "Success one chat",
			inChat: dto.ChatGet{
				UserID: &userID,
			},
			dataDB: entity.ChatGet{
				UserID: userID,
			},
			mock: func(s *mockRepo.MockChat, dataDB entity.ChatGet) {
				s.EXPECT().GetChat(dataDB).Return([]entity.Chat{
					{
						Id:        userID,
						Name:      "chat_1",
						CreatedAt: "2024-09-20T18:26:13.239627Z",
						IsDeleted: false,
					},
				}, nil)
			},
			want: []entity.Chat{
				{
					Id:        userID,
					Name:      "chat_1",
					CreatedAt: "2024-09-20T18:26:13.239627Z",
					IsDeleted: false,
				},
			},
			wantErr: nil,
		},
		{
			name: "Success many chats",
			inChat: dto.ChatGet{
				UserID: &userID,
			},
			dataDB: entity.ChatGet{
				UserID: userID,
			},
			mock: func(s *mockRepo.MockChat, dataDB entity.ChatGet) {
				s.EXPECT().GetChat(dataDB).Return([]entity.Chat{
					{
						Id:        userID,
						Name:      "chat_1",
						CreatedAt: "2024-09-20T18:26:13.239627Z",
						IsDeleted: false,
					},
					{
						Id:        userID,
						Name:      "chat_2",
						CreatedAt: "2024-09-19T18:26:13.239627Z",
						IsDeleted: false,
					},
					{
						Id:        userID,
						Name:      "chat_3",
						CreatedAt: "2024-09-18T18:26:13.239627Z",
						IsDeleted: false,
					},
				}, nil)
			},
			want: []entity.Chat{
				{
					Id:        userID,
					Name:      "chat_1",
					CreatedAt: "2024-09-20T18:26:13.239627Z",
					IsDeleted: false,
				},
				{
					Id:        userID,
					Name:      "chat_2",
					CreatedAt: "2024-09-19T18:26:13.239627Z",
					IsDeleted: false,
				},
				{
					Id:        userID,
					Name:      "chat_3",
					CreatedAt: "2024-09-18T18:26:13.239627Z",
					IsDeleted: false,
				},
			},
			wantErr: nil,
		},
		{
			name: "Empty user",
			inChat: dto.ChatGet{
				UserID: &userEmpty,
			},
			mock:    func(s *mockRepo.MockChat, dataDB entity.ChatGet) {},
			want:    nil,
			wantErr: errors.New("user_id is empty"),
		},
		{
			name:    "Nil user",
			mock:    func(s *mockRepo.MockChat, dataDB entity.ChatGet) {},
			want:    nil,
			wantErr: errors.New("user_id is empty"),
		},
		{
			name: "Other error",
			inChat: dto.ChatGet{
				UserID: &userID,
			},
			dataDB: entity.ChatGet{
				UserID: userID,
			},
			mock: func(s *mockRepo.MockChat, dataDB entity.ChatGet) {
				s.EXPECT().GetChat(dataDB).Return(nil, errors.New("some error"))
			},
			want:    nil,
			wantErr: errors.New("some error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//Инициализируем контролер для мока сервиса
			ctrl := gomock.NewController(t)
			// Завершаем работу контролера после выполнения каждого теста
			defer ctrl.Finish()

			// Создаём моки базы данных авторизации
			mockChat := mockRepo.NewMockChat(ctrl)

			// Передаём структуру пользователя
			tt.mock(mockChat, tt.dataDB)

			// Создаём объект базы данных в который передадим наш мок авторизации
			repository := &db.DB{Chat: mockChat}

			// Создаём экземпляр сервиса авторизации
			serviceChat := NewChatService(repository)

			// Проверяем ожидаемый и актуальный результат
			acChats, acErr := serviceChat.GetChat(tt.inChat)
			assert.Equal(t, tt.want, acChats)
			assert.Equal(t, tt.wantErr, acErr)
		})
	}
}

func TestChatService_DeleteChat(t *testing.T) {
	// Структура для последующей реализации поведения мока
	type mockBehaviour func(s *mockRepo.MockChat, dataDB entity.ChatDelete)

	tests := []struct {
		name    string
		inChat  dto.ChatDelete
		dataDB  entity.ChatDelete
		mock    mockBehaviour
		want    []entity.DeletedChats
		wantErr error
	}{
		{
			name: "Success one chat",
			inChat: dto.ChatDelete{
				ChatIds: &[]int64{1},
			},
			dataDB: entity.ChatDelete{
				ChatIds: []int64{1},
				UserID:  1,
			},
			mock: func(s *mockRepo.MockChat, dataDB entity.ChatDelete) {
				s.EXPECT().DeleteChat(dataDB).Return([]entity.DeletedChats{
					{
						ChatID: 1,
						Result: "Chat successfully deleted",
					},
				}, nil)
			},
			want: []entity.DeletedChats{
				{
					ChatID: 1,
					Result: "Chat successfully deleted",
				},
			},
			wantErr: nil,
		},
		{
			name: "Success many chat",
			inChat: dto.ChatDelete{
				ChatIds: &[]int64{1, 2, 3},
			},
			dataDB: entity.ChatDelete{
				ChatIds: []int64{1, 2, 3},
				UserID:  1,
			},
			mock: func(s *mockRepo.MockChat, dataDB entity.ChatDelete) {
				s.EXPECT().DeleteChat(dataDB).Return([]entity.DeletedChats{
					{
						ChatID: 1,
						Result: "Chat successfully deleted",
					},
					{
						ChatID: 2,
						Result: "Chat does not exist or has already been deleted",
					},
					{
						ChatID: 3,
						Result: "Chat successfully deleted",
					},
				}, nil)
			},
			want: []entity.DeletedChats{
				{
					ChatID: 1,
					Result: "Chat successfully deleted",
				},
				{
					ChatID: 2,
					Result: "Chat does not exist or has already been deleted",
				},
				{
					ChatID: 3,
					Result: "Chat successfully deleted",
				},
			},
			wantErr: nil,
		},
		{
			name: "Empty chat_ids",
			inChat: dto.ChatDelete{
				ChatIds: &[]int64{},
			},
			mock: func(s *mockRepo.MockChat, dataDB entity.ChatDelete) {
			},
			want:    nil,
			wantErr: errors.New("chat_ids is empty"),
		},
		{
			name: "Empty userID",
			inChat: dto.ChatDelete{
				ChatIds: &[]int64{1, 2, 3},
			},
			dataDB: entity.ChatDelete{
				UserID: 0,
			},
			mock:    func(s *mockRepo.MockChat, dataDB entity.ChatDelete) {},
			want:    nil,
			wantErr: errors.New("user_id is empty"),
		},
		{
			name: "Nil chat_ids",
			mock: func(s *mockRepo.MockChat, dataDB entity.ChatDelete) {
			},
			want:    nil,
			wantErr: errors.New("chat_ids is empty"),
		},
		{
			name: "Other error",
			inChat: dto.ChatDelete{
				ChatIds: &[]int64{1},
			},
			dataDB: entity.ChatDelete{
				ChatIds: []int64{1},
				UserID:  1,
			},
			mock: func(s *mockRepo.MockChat, dataDB entity.ChatDelete) {
				s.EXPECT().DeleteChat(dataDB).Return(nil, errors.New("some error"))
			},
			want:    nil,
			wantErr: errors.New("some error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//Инициализируем контролер для мока сервиса
			ctrl := gomock.NewController(t)
			// Завершаем работу контролера после выполнения каждого теста
			defer ctrl.Finish()

			// Создаём моки базы данных авторизации
			mockChat := mockRepo.NewMockChat(ctrl)

			// Передаём структуру пользователя
			tt.mock(mockChat, tt.dataDB)

			// Создаём объект базы данных в который передадим наш мок авторизации
			repository := &db.DB{Chat: mockChat}

			// Создаём экземпляр сервиса авторизации
			serviceChat := NewChatService(repository)

			// Проверяем ожидаемый и актуальный результат
			acChats, acErr := serviceChat.DeleteChat(tt.inChat, tt.dataDB.UserID)
			assert.Equal(t, tt.want, acChats)
			assert.Equal(t, tt.wantErr, acErr)
		})
	}
}
