package rc

import (
	"testing"
)

func TestWebHook_Send(t *testing.T) {

	type args struct {
		msg Message
	}
	tests := []struct {
		name    string
		h       *WebHook
		args    args
		wantErr bool
	}{
		{
			name: "hook_test",
			h:    NewWebHook("http://localhost:3000/hooks/pfGsyM2NjiMxntNK9/5wpZd8XMDN2izBeFhkoAeEPjQujgrhZfabJoHi2BA6yr4z5e"),
			args: args{
				msg: Message{
					Text:  "this is a test",
					Alias: "Admin",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.h.Send(tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("WebHook.Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
