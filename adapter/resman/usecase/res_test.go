package usecase

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/visionom/v-gae/adapter/mysql"
	"github.com/visionom/v-gae/adapter/mysql/config"
	"github.com/visionom/v-gae/adapter/resman/domain"
)

func init() {
	mysql.Init(&config.DBConfig{
		Host:   "127.0.0.1",
		Port:   3306,
		Name:   "test_resman",
		User:   "root",
		Passwd: "",
	})
	mysql.Update("truncate res")
	mysql.Update("truncate tags")
	c := NewComRepo()
	c.NewRes(genRes(-1000, -2, true, 100, true))
}

func genRes(start, end int, hasTag bool, tagNum int, hasInfo bool) []domain.Res {
	var reses []domain.Res
	for i := start; i < end; i++ {
		id := fmt.Sprintf("test_id_%d", i)
		name := fmt.Sprintf("test_name_%d", i)
		var info json.RawMessage
		if hasInfo {
			info = json.RawMessage{'{', '}'}
		}
		res := domain.NewRes(id, name, info)
		if hasTag {
			for j := 0; j < tagNum; j++ {
				tag := fmt.Sprintf("test_tag_%d", j)
				res.AddTags(domain.NewTag(tag, tag))
			}
		}
		reses = append(reses, res)
	}
	return reses
}

func genIDs(start, end int) []string {
	var resIDs []string
	for i := start; i < end; i++ {
		resIDs = append(resIDs, fmt.Sprintf("test_id_%d", i))
	}
	return resIDs
}

func TestComRepoImpl_NewRes(t *testing.T) {
	type args struct {
		reses []domain.Res
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"insert nil", args{nil}, false},
		{"insert empty", args{genRes(0, 0, false, 0, false)}, false},
		{"insert one", args{genRes(0, 1, false, 0, false)}, false},
		{"insert one with empty id", args{[]domain.Res{domain.NewRes("", "test", json.RawMessage{'{', '}'})}}, false},
		{"insert two", args{genRes(2, 3, false, 0, false)}, false},
		{"insert 100", args{genRes(4, 104, false, 0, false)}, false},
		{"insert one with info", args{genRes(1000, 1001, false, 0, true)},
			false},
		{"insert two with info", args{genRes(1002, 1003, false, 0, true)},
			false},
		{"insert 100 with info", args{genRes(1004, 1104, false, 0, true)},
			false},
		{"insert one with tags", args{genRes(2000, 2001, true, 100, true)},
			false},
		{"insert two with tags", args{genRes(2002, 2003, true, 100, true)},
			false},
		{"insert 100 with tags", args{genRes(2004, 2104, true, 100, true)},
			false},
		{"insert one repeat", args{genRes(0, 1, false, 0, true)},
			true},
		{"insert two and one repeat", args{genRes(-1, 1, false, 0, true)},
			true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ComRepoImpl{}
			if err := c.NewRes(tt.args.reses); (err != nil) != tt.wantErr {
				t.Errorf("ComRepoImpl.NewRes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestComRepoImpl_DelRes(t *testing.T) {
	type args struct {
		ids []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"delete nil", args{nil}, false},
		{"delete empty", args{genIDs(0, 0)}, false},
		{"delete one", args{genIDs(-1000, -999)}, false},
		{"delete 100", args{genIDs(-998, -898)}, false},
		{"delete 100 not exist", args{genIDs(-1000, -900)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ComRepoImpl{}
			if err := c.DelRes(tt.args.ids); (err != nil) != tt.wantErr {
				t.Errorf("ComRepoImpl.DelRes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestComRepoImpl_ModRes(t *testing.T) {
	type args struct {
		res domain.Res
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"modify", args{
				domain.NewRes(
					fmt.Sprintf("test_id_%d", 0),
					fmt.Sprintf("new_test_name_%d", 0),
					json.RawMessage{'{', '}'},
				)},
			false,
		},
		{
			"modify with empty id", args{
				domain.NewRes(
					"",
					fmt.Sprintf("new_test_name_%d", 0),
					json.RawMessage{'{', '}'},
				)},
			true,
		}, {
			"modify with empty name", args{
				domain.NewRes(
					fmt.Sprintf("test_id_%d", 0),
					"",
					json.RawMessage{'{', '}'},
				)},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ComRepoImpl{}
			if err := c.ModRes(tt.args.res); (err != nil) != tt.wantErr {
				t.Errorf("ComRepoImpl.ModRes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestComRepoImpl_AddTag(t *testing.T) {
	type args struct {
		resID string
		tags  domain.Tags
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"add tag empty",
			args{
				"test",
				domain.Tags{},
			},
			false,
		},
		{
			"add tag ",
			args{
				"test",
				domain.Tags{
					domain.NewTag("tag1", "v1"),
					domain.NewTag("tag2", "v2"),
				},
			},
			false,
		},
		{
			"add tag repeat",
			args{
				"test",
				domain.Tags{
					domain.NewTag("tag1", "v1"),
					domain.NewTag("tag2", "v2"),
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ComRepoImpl{}
			if err := c.AddTag(tt.args.resID, tt.args.tags); (err != nil) != tt.wantErr {
				t.Errorf("ComRepoImpl.AddTag() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestComRepoImpl_RmTag(t *testing.T) {
	type args struct {
		resID string
		tags  domain.Tags
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"del tag empty",
			args{
				"test",
				domain.Tags{},
			},
			false,
		},
		{
			"del tag",
			args{
				"test",
				domain.Tags{
					domain.NewTag("tag1", "v1"),
					domain.NewTag("tag2", "v2"),
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ComRepoImpl{}
			_, err := c.RmTag(tt.args.resID, tt.args.tags)
			if (err != nil) != tt.wantErr {
				t.Errorf("ComRepoImpl.RmTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestComRepoImpl_ChangeTag(t *testing.T) {
	type args struct {
		resID string
		key   string
		old   string
		new   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"change empty", args{}, false},
		{"change empty", args{"test_id_-100", "test_tag_10", "test_tag_10", "new_test_tag_10"},
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ComRepoImpl{}
			if err := c.ChangeTag(tt.args.resID, tt.args.key, tt.args.old, tt.args.new); (err != nil) != tt.wantErr {
				t.Errorf("ComRepoImpl.ChangeTags() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestComRepoImpl_FindRes(t *testing.T) {
	type args struct {
		tags domain.Tags
		page int
		size int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"find empty", args{}, false},
		{"find one", args{domain.Tags{domain.NewTag("test_tag_1", "test_tag_1")}, 0, 10}, false},
		{"find three", args{domain.Tags{domain.NewTag("test_tag_1", "test_tag_1"),
			domain.NewTag("test_tag_2", "test_tag_2"),
			domain.NewTag("test_tag_3", "test_tag_3"),
		},
			10, 100},
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ComRepoImpl{}
			gotReses, err := c.FindRes(tt.args.tags, tt.args.page, tt.args.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("ComRepoImpl.FindRes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("%v", gotReses)
		})
	}
}

func Test_findResById(t *testing.T) {
	type args struct {
		resIDs []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"find empty", args{}, false},
		{"find empty", args{resIDs: genIDs(0, 100)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotReses, err := findResById(tt.args.resIDs)
			if (err != nil) != tt.wantErr {
				t.Errorf("findResById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("%v", gotReses)
		})
	}
}

func TestComRepoImpl_FindResByID(t *testing.T) {
	type args struct {
		ids []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"find empty", args{}, false},
		{"find empty", args{genIDs(0, 100)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ComRepoImpl{}
			gotReses, err := c.FindResByID(tt.args.ids)
			if (err != nil) != tt.wantErr {
				t.Errorf("ComRepoImpl.FindResByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("%v", gotReses)
		})
	}
}

func TestComRepoImpl_Count(t *testing.T) {
	type args struct {
		tags domain.Tags
	}
	tests := []struct {
		name    string
		args    args
		wantC   int
		wantErr bool
	}{
		{"find empty", args{}, 0, false},
		{"find one", args{domain.Tags{domain.NewTag("test_tag_1", "test_tag_1")}}, 998, false},
		{"find three", args{domain.Tags{
			domain.NewTag("test_tag_1", "test_tag_1"),
			domain.NewTag("test_tag_2", "test_tag_2"),
			domain.NewTag("test_tag_3", "test_tag_3")}}, 998, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ComRepoImpl{}
			gotC, err := c.Count(tt.args.tags)
			if (err != nil) != tt.wantErr {
				t.Errorf("ComRepoImpl.Count() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotC != tt.wantC {
				t.Errorf("ComRepoImpl.Count() = %v, want %v", gotC, tt.wantC)
			}
		})
	}
}
