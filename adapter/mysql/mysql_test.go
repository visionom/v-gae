package mysql

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestInit(t *testing.T) {
	type args struct {
		user   string
		passwd string
		host   string
		db     string
		port   int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"A",
			args{
				"root",
				"",
				"localhost",
				"test",
				3306,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Init(tt.args.user, tt.args.passwd, tt.args.host,
				tt.args.db, tt.args.port); (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDBMysql_Find(t *testing.T) {
	Init("root", "", "127.0.0.1", "test", 3306)
	type args struct {
		sql  string
		args []interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"A",
			args{
				sql:  "select * from pools",
				args: []interface{}{},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Query(tt.args.sql, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("DB.Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for _, rd := range got {
				t.Logf("%s\n", rd)
			}
		})
	}
}

func TestDBMysql_FindFirst(t *testing.T) {
	Init("root", "", "127.0.0.1", "test", 3306)
	type args struct {
		sql  string
		args []interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"A",
			args{
				sql:  "select * from pools",
				args: []interface{}{},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := QueryFirst(tt.args.sql, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("DB.FindFirst() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("%+v", got)
		})
	}
}
