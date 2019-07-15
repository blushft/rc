package rc

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/gopackage/ddp"
	"github.com/mitchellh/mapstructure"
)

var (
	DefaultEventBuffer   = 10
	DefaultMessageBuffer = 50
)

type StreamSubscription struct {
	name   string
	events map[NotificationEvent]string
}

type NotificationEvent int

const (
	minMsgBuffer = 10
	minEvtBuffer = 10

	// Events used by several subs
	NotifyUpdateAvatar = iota
	NotifyRolesChange
	NotifyUpdateEmoji
	NotifyDeleteEmoji
	NotifyWebRTC
	// SubNotifyAll Events
	NotifyPublicSettingsChanged
	NotifyPermissionsChanged
	// SubNotifyLogged Events
	NotifyNameChanged
	NotifyUserDeleted
	NotifyUserStatus
	// SubNotifyUser Events
	NotifyMessage
	NotifyOTR
	NotifyNotification
	NotifyRoomsChanged
	NotifySubsChanged
	// SubNotifyRoom Events
	NotifyDeleteMessage
	NotifyTyping
)

var (
	SubNotifyAll = StreamSubscription{
		name: "stream-notify-all",
		events: map[NotificationEvent]string{
			NotifyUpdateAvatar:          "updateAvatar",
			NotifyUpdateEmoji:           "updateEmojiCustom",
			NotifyDeleteEmoji:           "deleteEmojiCustom",
			NotifyRolesChange:           "roles-change",
			NotifyPublicSettingsChanged: "public-settings-changed",
			NotifyPermissionsChanged:    "permissions-changed",
		},
	}
	SubNotifyLogged = StreamSubscription{
		name: "stream-notify-logged",
		events: map[NotificationEvent]string{
			NotifyNameChanged:  "Users:NameChanged",
			NotifyUserDeleted:  "Users:Deleted",
			NotifyUpdateAvatar: "updateAvatar",
			NotifyUpdateEmoji:  "updateEmojiCustom",
			NotifyDeleteEmoji:  "deleteEmojiCustom",
			NotifyRolesChange:  "roles-change",
			NotifyUserStatus:   "user-status",
		}}
	SubNotifyRoom = StreamSubscription{
		name: "stream-notify-room",
		events: map[NotificationEvent]string{
			NotifyDeleteMessage: "deleteMessage",
			NotifyTyping:        "typing",
		}}

	SubNotifyUser = StreamSubscription{
		name: "stream-notify-user",
		events: map[NotificationEvent]string{
			NotifyMessage:      "message",
			NotifyOTR:          "otr",
			NotifyWebRTC:       "webrtc",
			NotifyNotification: "notification",
			NotifyRoomsChanged: "rooms-changed",
			NotifySubsChanged:  "subscriptions-changed",
		}}
	SubNotifyRoomUsers = StreamSubscription{
		name: "stream-notify-room-users",
		events: map[NotificationEvent]string{
			NotifyWebRTC: "webrtc",
		}}
)

func getStreamSub(name string) *StreamSubscription {
	switch name {
	case "stream-notify-room-users":
		return &SubNotifyRoomUsers
	case "stream-notify-user":
		return &SubNotifyUser
	case "stram-notify-room":
		return &SubNotifyRoom
	case "stream-notify-logged":
		return &SubNotifyLogged
	case "stream-notify-all":
		return &SubNotifyAll
	default:
		return nil
	}
}

type StreamOption func(*streams)

func EventSubscription(sub StreamSubscription, evt NotificationEvent) StreamOption {
	return func(str *streams) {
		var ne []NotificationEvent
		if v, ok := str.subEvts[sub.name]; ok {
			ne = v
		} else {
			ne = make([]NotificationEvent, 0)
		}
		ne = append(ne, evt)
		str.subEvts[sub.name] = ne
	}
}

func RoomSubscription(roomID string) StreamOption {
	return func(str *streams) {
		str.subRooms = append(str.subRooms, roomID)
	}
}

type streams struct {
	subEvts  map[string][]NotificationEvent
	subRooms []string

	msgs []*SubChannel
	evts []*SubChannel

	allMsgs chan []RoomMessage
	allEvts chan *StreamEvent
	allErrs chan error
}

func newStreams(opts ...StreamOption) (*streams, error) {
	str := &streams{
		subEvts:  make(map[string][]NotificationEvent),
		subRooms: make([]string, 0),
		msgs:     make([]*SubChannel, 0),
		evts:     make([]*SubChannel, 0),
	}

	for _, o := range opts {
		o(str)
	}

	msgBuf := len(str.subRooms) * DefaultMessageBuffer
	if msgBuf == 0 {
		msgBuf = minMsgBuffer
	}

	evtBuf := 0
	for _, ev := range str.subEvts {
		evtBuf += len(ev) * DefaultEventBuffer
	}
	if evtBuf == 0 {
		evtBuf = minEvtBuffer
	}

	str.allMsgs = make(chan []RoomMessage, msgBuf)
	str.allEvts = make(chan *StreamEvent, evtBuf)
	str.allErrs = make(chan error, 10)

	return str, nil
}

func (str *streams) iterMsg(c *SubChannel) {
	for {
		select {
		case mm := <-c.Updates:
			m, ok := mm.([]RoomMessage)
			if !ok {
				continue
			}
			str.allMsgs <- m
		case err := <-c.Errors:
			str.allErrs <- err
		}
	}
}

func (str *streams) iterEvt(c *SubChannel) {
	for {
		select {
		case mm := <-c.Updates:
			m, ok := mm.(*StreamEvent)
			if !ok {
				continue
			}
			str.allEvts <- m
		case err := <-c.Errors:
			str.allErrs <- err
		}
	}
}

func (str *streams) runStreams(c *ddp.Client) error {
	for _, v := range str.subRooms {
		ns, err := subscribeToRoomMessages(v, c)
		if err != nil {
			return err
		}
		str.msgs = append(str.msgs, ns)
	}

	for _, msgsub := range str.msgs {
		go str.iterMsg(msgsub)
	}

	for k, vv := range str.subEvts {
		ss := getStreamSub(k)
		if ss == nil {
			continue
		}
		for _, v := range vv {
			ns, err := subscribeToEvent(*ss, v, c)
			if err != nil {
				return err
			}
			str.evts = append(str.evts, ns)
		}
	}

	for _, evtsub := range str.evts {
		go str.iterEvt(evtsub)
	}

	return nil
}

type StreamEvent struct {
	Event string
	Args  []interface{}
}

func subscribeToEvent(subName StreamSubscription, evt NotificationEvent, c *ddp.Client) (*SubChannel, error) {
	subn := subName.name
	sube, ok := subName.events[evt]
	if !ok {
		return nil, fmt.Errorf("%s does not have event %d", subn, evt)
	}
	err := c.Sub(subn, sube, true)
	if err != nil {
		return nil, err
	}

	fn := func(u ddp.Update) (interface{}, error) {
		e, ok := u["eventName"].(string)
		if !ok {
			return nil, errors.New("invalid response")
		}
		args, ok := u["args"].([]interface{})
		if !ok {
			return nil, errors.New("invalid response")
		}
		res := &StreamEvent{
			Event: e,
			Args:  args,
		}
		return res, nil
	}

	list, sub := NewUpdateListener(fn)

	c.CollectionByName(subn).
		AddUpdateListener(list)

	return sub, nil
}

func subscribeToRoomMessages(roomID string, c *ddp.Client) (*SubChannel, error) {
	err := c.Sub("stream-room-messages", roomID, true)
	if err != nil {
		return nil, err
	}

	fn := func(u ddp.Update) (interface{}, error) {
		rawmsgs := []RoomMessage{}
		m, ok := u["args"]
		if !ok {
			return nil, fmt.Errorf("unexpected args: %v", u)
		}
		decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			WeaklyTypedInput: true,
			TagName:          "json",
			Result:           &rawmsgs,
		})
		if err != nil {
			return nil, err
		}
		err = decoder.Decode(m)
		if err != nil {
			return nil, err
		}
		msgs := []RoomMessage{}
		for _, mm := range rawmsgs {
			if mm.ID != "" {
				msgs = append(msgs, mm)
			}
		}
		return msgs, nil
	}

	list, sub := NewUpdateListener(fn)

	c.CollectionByName("stream-room-messages").
		AddUpdateListener(list)

	return sub, nil
}

func newId() string {
	sid := uuid.New().String()
	return strings.Replace(sid, "-", "", -1)
}

func mapUserStatus(args []interface{}) (map[string]string, error) {
	for _, x := range args {
		status := map[string]string{}
		switch v := x.(type) {
		case []interface{}:
			if len(v) > 0 {
				status["id"] = v[0].(string)
				status["user"] = v[1].(string)
				status["status"] = userStatus(v[2].(int))
				status["message"] = v[3].(string)
				return status, nil
			}
			return nil, errors.New("not a status message")
		}
	}
	return nil, errors.New("not a status message")
}

func userStatus(i int) string {
	switch i {
	case 0:
		return "offline"
	case 1:
		return "online"
	case 2:
		return "away"
	case 3:
		return "busy"
	default:
		return "unknown"
	}
}
