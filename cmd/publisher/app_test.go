package main

import (
	"context"
	"testing"
)

func TestApp_RunStream(t *testing.T) {
	type fields struct {
		writer StreamWriter
	}

	type args struct {
		ctx context.Context
		n   int
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test",
			fields: fields{
				writer: nil,
			},
			args: args{
				ctx: context.Background(),
				n:   2,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &App{
				writer: tt.fields.writer,
			}

			if err := a.RunStream(tt.args.ctx, tt.args.n); (err != nil) != tt.wantErr {
				t.Errorf("RunConsume() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
