package rc

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func clientOpts() []ClientOption {
	return []ClientOption{
		ServerURL("https://chat.dev.rpxapps.com"),
		CredFromJson("./token.json"),
		//Debug(true),
		Realtime(true),
	}
}

func createTestClient() *Client {
	opts := clientOpts()
	return New(opts...)
}

func TestNew(t *testing.T) {

	client := createTestClient()
	if err := client.Connect(); err != nil {
		t.Error(err)
	}

	tests := []struct {
		name string
	}{
		{
			name: "testrc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			result := client.c.getInfo()
			if result.Error() != nil {
				t.Errorf("info request failed: %v", result.Error())
			}
			defer client.log.Sync()
			client.log.Infow("string result: "+result.String(),
				"code", result.StatusCode(),
			)
			v := make(map[string]interface{})
			if err := result.JSON(&v); err != nil {
				t.Errorf("error getting result json: %v", err)
			}

			rooms, err := client.GetSubscriptions()
			if err != nil {
				t.Errorf("error getting rooms: %v", err)
			}

			for _, r := range rooms {
				if r.Name == "" {
					spew.Dump(r)
				} else {
					client.log.Infow(r.Name, "id", r.RoomID, "type", r.Type)
				}

			}
			/*
				_, err = client.GetRoomByName("general")
				if err != nil {
					t.Errorf("error getting room general: %v", err)
				}

				//spew.Dump(room)

				users, err := client.GetUsers()
				if err != nil {
					t.Errorf("error getting users: %v", err)
				}

				spew.Dump(users) */
		})
	}
}
