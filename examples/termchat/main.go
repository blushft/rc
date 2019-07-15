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
		rc.StreamOptions(
			rc.RoomSubscription("__my_messages__"),
			rc.EventSubscription(rc.SubNotifyLogged, rc.NotifyUserStatus),
		),
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

	for {
		select {
		case e := <-client.EventStream():
			log.Printf("Event: %s\n", e.Event)
			log.Printf("Args: %#v", e.Args)
		case m := <-client.MessageStream():
			log.Printf("%d messages\n", len(m))
			for _, msg := range m {
				log.Printf("From: %s - %s\n", msg.User.Username, msg.Msg)
			}

		case serr := <-client.StreamErrors():
			log.Fatal(serr)
		default:
		}
	}
}
