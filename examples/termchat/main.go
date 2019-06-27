package main

import (
	"log"

	"github.com/blushft/rc"
)

func clientOpts() []rc.ClientOption {
	return []rc.ClientOption{
		rc.ServerURL("http://localhost:3000"),
		rc.CredFromJson("../../token.json"),
		rc.Debug(true),
		rc.Realtime(true),
	}
}

func createTestClient() *rc.Client {
	opts := clientOpts()
	return rc.New(opts...)
}

func main() {
	client := createTestClient()
	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}

	subs, err := client.SubscribeToRoomMessages("__my_messages__")
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case m := <-subs.Updates:
			msgs, ok := m.([]rc.RoomMessage)
			if !ok {
				log.Fatalf("this isn't messages? %v", m)
			}
			log.Printf("%d messages\n", len(msgs))
			for _, msg := range msgs {
				log.Printf("From: %s - %s\n", msg.User.Username, msg.Msg)
			}
		case serr := <-subs.Errors:
			log.Fatal(serr)
		}
	}
}
