package model_test

import (
	"errors"
	"testing"

	"github.com/enolgor/lurra/internal/model"
)

func TestCreateTodo_Validation(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		wantErr error
	}{
		{name: "ok", title: "buy milk"},
		{name: "with spaces", title: "  buy milk  "},
		{name: "empty", title: "", wantErr: model.ErrTitleInvalid},
		{name: "only spaces", title: "  ", wantErr: model.ErrTitleInvalid},
		{name: "tab and newline", title: "\t\n", wantErr: model.ErrTitleInvalid},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createTodo := model.CreateTodo{Title: tt.title}
			gotErr := createTodo.Validation()
			if !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("Validation() error = %v, wantErr %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestUpdateTodo_Validation(t *testing.T) {
	tests := []struct {
		name    string
		title   *string
		wantErr error
	}{
		{name: "nil"},
		{name: "ok", title: new("buy milk")},
		{name: "with spaces", title: new("  buy milk  ")},
		{name: "empty", title: new(""), wantErr: model.ErrTitleInvalid},
		{name: "only spaces", title: new("  "), wantErr: model.ErrTitleInvalid},
		{name: "tab and newline", title: new("\t\n"), wantErr: model.ErrTitleInvalid},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updateTodo := model.UpdateTodo{Title: tt.title}
			gotErr := updateTodo.Validation()
			if !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("Validation() error = %v, wantErr %v", gotErr, tt.wantErr)
			}
		})
	}
}
