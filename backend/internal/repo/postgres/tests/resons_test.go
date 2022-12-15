package tests

import (
	"log"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func NewMock() (*sqlx.DB, sqlmock.Sqlmock) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	db := sqlx.NewDb(mockDB, "sqlmock")
	return db, mock
}

// func TestReasonRepository_Get(t *testing.T) {

// }

// func TestReasonRepository_Create(t *testing.T) {
// 	type args struct {
// 		ctx   context.Context
// 		input models.ReasonDTO
// 	}

// 	testReason := models.ReasonDTO{
// 		Id:          uuid.New(),
// 		OperationId: uuid.New(),
// 		Date:        time.Now().Format("02.01.2006 15:04"),
// 		Value:       "test value",
// 	}

// 	tests := []struct {
// 		name       string
// 		args       args
// 		beforeTest func(sqlmock.Sqlmock)
// 		want       uuid.UUID
// 		wantErr    bool
// 	}{
// 		{
// 			name: "success create reason",
// 			args: args{
// 				ctx: context.TODO(),
// 				input: models.ReasonDTO{
// 					Id:          testReason.Id,
// 					OperationId: testReason.OperationId,
// 					Date:        testReason.Date,
// 					Value:       testReason.Value,
// 				},
// 			},
// 			beforeTest: func(mockSQL sqlmock.Sqlmock) {
// 				mockSQL.ExpectQuery("INSERT INTO reasons (id, operation_id, date, value) VALUES ($1, $2, $3, $4)").
// 					WithArgs(testReason.Id, testReason.OperationId, testReason.Date, testReason.Value)
// 			},
// 			want: testReason.Id,
// 		},
// 		{
// 			name: "failed create reason",
// 			args: args{
// 				ctx:   context.TODO(),
// 				input: models.ReasonDTO{},
// 			},
// 			beforeTest: func(mockSQL sqlmock.Sqlmock) {
// 				mockSQL.
// 					ExpectQuery("INSERT INTO reasons (id, operation_id, date, value) VALUES ($1, $2, $3, $4)").
// 					WithArgs("", "", "", "").
// 					WillReturnError(fmt.Errorf("failed to execute query"))
// 			},
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			db, mock := NewMock()
// 			defer db.Close()

// 			repo := postgres.NewReasonRepo(db)

// 			if tt.beforeTest != nil {
// 				tt.beforeTest(mock)
// 			}

// 			got, err := repo.Create(tt.args.ctx, tt.args.input)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("reasonRepo.Create() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("reasonRepo.Create() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
