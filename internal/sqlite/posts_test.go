package sqlite

import (
	"database/sql"
	"github.com/gary-norman/forum/internal/models"
	"reflect"
	"testing"
)

func TestPostModel_All(t *testing.T) {
	type fields struct {
		DB *sql.DB
	}
	tests := []struct {
		name    string
		fields  fields
		want    []models.Post
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &PostModel{
				DB: tt.fields.DB,
			}
			got, err := m.All()
			if (err != nil) != tt.wantErr {
				t.Errorf("All() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("All() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostModel_GetPostsByChannel(t *testing.T) {
	type fields struct {
		DB *sql.DB
	}
	type args struct {
		channel int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.Post
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &PostModel{
				DB: tt.fields.DB,
			}
			got, err := m.GetPostsByChannel(tt.args.channel)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPostsByChannel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPostsByChannel() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostModel_Insert(t *testing.T) {
	type fields struct {
		DB *sql.DB
	}
	type args struct {
		title        string
		content      string
		images       string
		author       string
		channel      string
		authorAvatar string
		channelID    int
		authorID     int
		commentable  bool
		isFlagged    bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &PostModel{
				DB: tt.fields.DB,
			}
			if err := m.Insert(tt.args.title, tt.args.content, tt.args.images, tt.args.author, tt.args.channel, tt.args.authorAvatar, tt.args.channelID, tt.args.authorID, tt.args.commentable, tt.args.isFlagged); (err != nil) != tt.wantErr {
				t.Errorf("Insert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
