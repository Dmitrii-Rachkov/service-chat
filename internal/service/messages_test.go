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

func TestMessageService_AddMessage(t *testing.T) {
	// Структура для последующей реализации поведения мока
	type mockBehaviour func(s *mockRepo.MockMessage, dataDB entity.MessageAdd)

	//Инициализируем контролер для мока сервиса
	ctrl := gomock.NewController(t)
	// Завершаем работу контролера после выполнения каждого теста
	defer ctrl.Finish()

	// Создаём моки базы данных сообщений
	mockMessage := mockRepo.NewMockMessage(ctrl)

	// Создаём объект базы данных в который передадим наш мок сообщений
	repository := &db.DB{Message: mockMessage}

	// Создаём экземпляр сервиса сообщений
	serviceChat := NewMessageService(repository)

	tests := []struct {
		name      string
		inMessage dto.MessageAdd
		dataDB    entity.MessageAdd
		mock      mockBehaviour
		want      int
		wantErr   error
	}{
		{
			name: "Success",
			inMessage: dto.MessageAdd{
				ChatID: 1,
				UserID: 1,
				Text:   "test",
			},
			dataDB: entity.MessageAdd{
				ChatID: 1,
				UserID: 1,
				Text:   "test",
			},
			mock: func(s *mockRepo.MockMessage, dataDB entity.MessageAdd) {
				s.EXPECT().AddMessage(dataDB).Return(1, nil)
			},
			want:    1,
			wantErr: nil,
		},
		{
			name: "Empty chat_id",
			inMessage: dto.MessageAdd{
				ChatID: 0,
				UserID: 1,
				Text:   "test",
			},
			mock:    func(s *mockRepo.MockMessage, dataDB entity.MessageAdd) {},
			want:    0,
			wantErr: errors.New("empty chat_id or user_id"),
		},
		{
			name: "Empty user_id",
			inMessage: dto.MessageAdd{
				ChatID: 1,
				UserID: 0,
				Text:   "test",
			},
			mock:    func(s *mockRepo.MockMessage, dataDB entity.MessageAdd) {},
			want:    0,
			wantErr: errors.New("empty chat_id or user_id"),
		},
		{
			name: "Empty text",
			inMessage: dto.MessageAdd{
				ChatID: 1,
				UserID: 1,
				Text:   "",
			},
			mock:    func(s *mockRepo.MockMessage, dataDB entity.MessageAdd) {},
			want:    0,
			wantErr: errors.New("empty text"),
		},
		{
			name: "Other error",
			inMessage: dto.MessageAdd{
				ChatID: 1,
				UserID: 1,
				Text:   "test",
			},
			dataDB: entity.MessageAdd{
				ChatID: 1,
				UserID: 1,
				Text:   "test",
			},
			mock: func(s *mockRepo.MockMessage, dataDB entity.MessageAdd) {
				s.EXPECT().AddMessage(dataDB).Return(0, errors.New("other error"))
			},
			want:    0,
			wantErr: errors.New("other error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Передаём структуру сообщения
			tt.mock(mockMessage, tt.dataDB)

			// Проверяем ожидаемый и актуальный результат
			acMsg, acErr := serviceChat.AddMessage(tt.inMessage)
			assert.Equal(t, tt.want, acMsg)
			assert.Equal(t, tt.wantErr, acErr)
		})
	}
}

func TestMessageService_UpdateMessage(t *testing.T) {
	// Структура для последующей реализации поведения мока
	type mockBehaviour func(s *mockRepo.MockMessage, dataDB entity.MessageUpdate)

	//Инициализируем контролер для мока сервиса
	ctrl := gomock.NewController(t)
	// Завершаем работу контролера после выполнения каждого теста
	defer ctrl.Finish()

	// Создаём моки базы данных сообщений
	mockMessage := mockRepo.NewMockMessage(ctrl)

	// Создаём объект базы данных в который передадим наш мок сообщений
	repository := &db.DB{Message: mockMessage}

	// Создаём экземпляр сервиса сообщений
	serviceChat := NewMessageService(repository)

	tests := []struct {
		name      string
		inMessage dto.MessageUpdate
		dataDB    entity.MessageUpdate
		mock      mockBehaviour
		want      int
		wantErr   error
	}{
		{
			name: "Success",
			inMessage: dto.MessageUpdate{
				MessageID: 1,
				UserID:    1,
				NewText:   "new_text",
			},
			dataDB: entity.MessageUpdate{
				MessageID: 1,
				UserID:    1,
				NewText:   "new_text",
			},
			mock: func(s *mockRepo.MockMessage, dataDB entity.MessageUpdate) {
				s.EXPECT().UpdateMessage(dataDB).Return(1, nil)
			},
			want:    1,
			wantErr: nil,
		},
		{
			name: "Empty message_id",
			inMessage: dto.MessageUpdate{
				MessageID: 0,
				UserID:    1,
				NewText:   "new_text",
			},
			mock:    func(s *mockRepo.MockMessage, dataDB entity.MessageUpdate) {},
			want:    0,
			wantErr: errors.New("empty chat_id or user_id"),
		},
		{
			name: "Empty user_id",
			inMessage: dto.MessageUpdate{
				MessageID: 1,
				UserID:    0,
				NewText:   "new_text",
			},
			mock:    func(s *mockRepo.MockMessage, dataDB entity.MessageUpdate) {},
			want:    0,
			wantErr: errors.New("empty chat_id or user_id"),
		},
		{
			name: "Empty new_text",
			inMessage: dto.MessageUpdate{
				MessageID: 1,
				UserID:    1,
				NewText:   "",
			},
			mock:    func(s *mockRepo.MockMessage, dataDB entity.MessageUpdate) {},
			want:    0,
			wantErr: errors.New("empty new_text"),
		},
		{
			name: "Other error",
			inMessage: dto.MessageUpdate{
				MessageID: 1,
				UserID:    1,
				NewText:   "new_text",
			},
			dataDB: entity.MessageUpdate{
				MessageID: 1,
				UserID:    1,
				NewText:   "new_text",
			},
			mock: func(s *mockRepo.MockMessage, dataDB entity.MessageUpdate) {
				s.EXPECT().UpdateMessage(dataDB).Return(0, errors.New("other error"))
			},
			want:    0,
			wantErr: errors.New("other error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Передаём структуру сообщения
			tt.mock(mockMessage, tt.dataDB)

			// Проверяем ожидаемый и актуальный результат
			acMsg, acErr := serviceChat.UpdateMessage(tt.inMessage)
			assert.Equal(t, tt.want, acMsg)
			assert.Equal(t, tt.wantErr, acErr)
		})
	}
}

func TestMessageService_GetMessage(t *testing.T) {
	// Структура для последующей реализации поведения мока
	type mockBehaviour func(s *mockRepo.MockMessage, dataDB entity.MessageGet)
	limit := int64(10)
	offset := int64(0)

	//Инициализируем контролер для мока сервиса
	ctrl := gomock.NewController(t)
	// Завершаем работу контролера после выполнения каждого теста
	defer ctrl.Finish()

	// Создаём моки базы данных сообщений
	mockMessage := mockRepo.NewMockMessage(ctrl)

	// Создаём объект базы данных в который передадим наш мок сообщений
	repository := &db.DB{Message: mockMessage}

	// Создаём экземпляр сервиса сообщений
	serviceChat := NewMessageService(repository)

	tests := []struct {
		name      string
		inMessage dto.MessageGet
		dataDB    entity.MessageGet
		mock      mockBehaviour
		want      []entity.Message
		wantErr   error
	}{
		{
			name: "Success one message",
			inMessage: dto.MessageGet{
				ChatID: 1,
				Limit:  &limit,
				Offset: &offset,
			},
			dataDB: entity.MessageGet{
				ChatID: 1,
				Limit:  limit,
				Offset: offset,
				UserID: 1,
			},
			mock: func(s *mockRepo.MockMessage, dataDB entity.MessageGet) {
				s.EXPECT().GetMessage(dataDB).Return([]entity.Message{
					{
						Id:        1,
						Text:      "text_1",
						UserID:    1,
						CreatedAt: "2024-09-19T18:26:13.239627Z",
						IsDeleted: false,
					},
				}, nil)
			},
			want: []entity.Message{
				{
					Id:        1,
					Text:      "text_1",
					UserID:    1,
					CreatedAt: "2024-09-19T18:26:13.239627Z",
					IsDeleted: false,
				},
			},
			wantErr: nil,
		},
		{
			name: "Success many message",
			inMessage: dto.MessageGet{
				ChatID: 1,
				Limit:  &limit,
				Offset: &offset,
			},
			dataDB: entity.MessageGet{
				ChatID: 1,
				Limit:  limit,
				Offset: offset,
				UserID: 1,
			},
			mock: func(s *mockRepo.MockMessage, dataDB entity.MessageGet) {
				s.EXPECT().GetMessage(dataDB).Return([]entity.Message{
					{
						Id:        1,
						Text:      "text_1",
						UserID:    1,
						CreatedAt: "2024-09-19T18:26:13.239627Z",
						IsDeleted: false,
					},
					{
						Id:        2,
						Text:      "text_2",
						UserID:    1,
						CreatedAt: "2024-09-18T18:26:13.239627Z",
						IsDeleted: false,
					},
					{
						Id:        3,
						Text:      "text_3",
						UserID:    1,
						CreatedAt: "2024-09-19T17:26:13.239627Z",
						IsDeleted: false,
					},
				}, nil)
			},
			want: []entity.Message{
				{
					Id:        1,
					Text:      "text_1",
					UserID:    1,
					CreatedAt: "2024-09-19T18:26:13.239627Z",
					IsDeleted: false,
				},
				{
					Id:        2,
					Text:      "text_2",
					UserID:    1,
					CreatedAt: "2024-09-18T18:26:13.239627Z",
					IsDeleted: false,
				},
				{
					Id:        3,
					Text:      "text_3",
					UserID:    1,
					CreatedAt: "2024-09-19T17:26:13.239627Z",
					IsDeleted: false,
				},
			},
			wantErr: nil,
		},
		{
			name: "Empty chat_id",
			inMessage: dto.MessageGet{
				ChatID: 0,
				Limit:  &limit,
				Offset: &offset,
			},
			mock:    func(s *mockRepo.MockMessage, dataDB entity.MessageGet) {},
			want:    nil,
			wantErr: errors.New("empty chat_id"),
		},
		{
			name: "Empty user_id",
			inMessage: dto.MessageGet{
				ChatID: 1,
				Limit:  &limit,
				Offset: &offset,
			},
			dataDB: entity.MessageGet{
				UserID: 0,
			},
			mock:    func(s *mockRepo.MockMessage, dataDB entity.MessageGet) {},
			want:    nil,
			wantErr: errors.New("empty user_id"),
		},
		{
			name: "Other errors",
			inMessage: dto.MessageGet{
				ChatID: 1,
				Limit:  &limit,
				Offset: &offset,
			},
			dataDB: entity.MessageGet{
				ChatID: 1,
				Limit:  limit,
				Offset: offset,
				UserID: 1,
			},
			mock: func(s *mockRepo.MockMessage, dataDB entity.MessageGet) {
				s.EXPECT().GetMessage(dataDB).Return(nil, errors.New("other error"))
			},
			want:    nil,
			wantErr: errors.New("other error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Передаём структуру сообщения
			tt.mock(mockMessage, tt.dataDB)

			// Проверяем ожидаемый и актуальный результат
			acMsg, acErr := serviceChat.GetMessage(tt.inMessage, tt.dataDB.UserID)
			assert.Equal(t, tt.want, acMsg)
			assert.Equal(t, tt.wantErr, acErr)
		})
	}
}

func TestMessageService_DeleteMessage(t *testing.T) {
	// Структура для последующей реализации поведения мока
	type mockBehaviour func(s *mockRepo.MockMessage, dataDB entity.MessageDel)

	//Инициализируем контролер для мока сервиса
	ctrl := gomock.NewController(t)
	// Завершаем работу контролера после выполнения каждого теста
	defer ctrl.Finish()

	// Создаём моки базы данных сообщений
	mockMessage := mockRepo.NewMockMessage(ctrl)

	// Создаём объект базы данных в который передадим наш мок сообщений
	repository := &db.DB{Message: mockMessage}

	// Создаём экземпляр сервиса сообщений
	serviceChat := NewMessageService(repository)

	tests := []struct {
		name      string
		inMessage dto.MessageDelete
		dataDB    entity.MessageDel
		mock      mockBehaviour
		want      []entity.DelMsg
		wantErr   error
	}{
		{
			name: "Success one message",
			inMessage: dto.MessageDelete{
				MessageIds: &[]int64{1},
			},
			dataDB: entity.MessageDel{
				MsgIds: []int64{1},
				UserID: 1,
			},
			mock: func(s *mockRepo.MockMessage, dataDB entity.MessageDel) {
				s.EXPECT().DeleteMessage(dataDB).Return([]entity.DelMsg{
					{
						MessageID: 1,
						Result:    "Message successfully deleted",
					},
				}, nil)
			},
			want: []entity.DelMsg{
				{
					MessageID: 1,
					Result:    "Message successfully deleted",
				},
			},
			wantErr: nil,
		},
		{
			name: "Success many message",
			inMessage: dto.MessageDelete{
				MessageIds: &[]int64{1},
			},
			dataDB: entity.MessageDel{
				MsgIds: []int64{1},
				UserID: 1,
			},
			mock: func(s *mockRepo.MockMessage, dataDB entity.MessageDel) {
				s.EXPECT().DeleteMessage(dataDB).Return([]entity.DelMsg{
					{
						MessageID: 1,
						Result:    "Message successfully deleted",
					},
					{
						MessageID: 2,
						Result:    "Message successfully deleted",
					},
					{
						MessageID: 3,
						Result:    "Message does not exist or has already been deleted",
					},
				}, nil)
			},
			want: []entity.DelMsg{
				{
					MessageID: 1,
					Result:    "Message successfully deleted",
				},
				{
					MessageID: 2,
					Result:    "Message successfully deleted",
				},
				{
					MessageID: 3,
					Result:    "Message does not exist or has already been deleted",
				},
			},
			wantErr: nil,
		},
		{
			name: "Empty messages_ids",
			inMessage: dto.MessageDelete{
				MessageIds: &[]int64{},
			},
			mock:    func(s *mockRepo.MockMessage, dataDB entity.MessageDel) {},
			want:    nil,
			wantErr: errors.New("message_ids is empty"),
		},
		{
			name: "Empty user_id",
			inMessage: dto.MessageDelete{
				MessageIds: &[]int64{1},
			},
			dataDB: entity.MessageDel{
				UserID: 0,
			},
			mock:    func(s *mockRepo.MockMessage, dataDB entity.MessageDel) {},
			want:    nil,
			wantErr: errors.New("user_id is empty"),
		},
		{
			name: "Other errors",
			inMessage: dto.MessageDelete{
				MessageIds: &[]int64{1},
			},
			dataDB: entity.MessageDel{
				MsgIds: []int64{1},
				UserID: 1,
			},
			mock: func(s *mockRepo.MockMessage, dataDB entity.MessageDel) {
				s.EXPECT().DeleteMessage(dataDB).Return(nil, errors.New("other error"))
			},
			want:    nil,
			wantErr: errors.New("other error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Передаём структуру сообщения
			tt.mock(mockMessage, tt.dataDB)

			// Проверяем ожидаемый и актуальный результат
			acMsg, acErr := serviceChat.DeleteMessage(tt.inMessage, tt.dataDB.UserID)
			assert.Equal(t, tt.want, acMsg)
			assert.Equal(t, tt.wantErr, acErr)
		})
	}
}
