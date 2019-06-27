package rc

import (
	"net/url"
	"strconv"
	"time"
)

type ChannelList struct {
	Channels []Channel `json:"channels"`
	Offset   int64     `json:"offset"`
	Count    int64     `json:"count"`
	Total    int64     `json:"total"`
	Success  bool      `json:"success"`
}

type Channel struct {
	ID        string      `json:"_id"`
	Name      string      `json:"name"`
	Type      string      `json:"t"`
	Usernames []string    `json:"usernames"`
	Msgs      int64       `json:"msgs"`
	User      ChannelUser `json:"u"`
	Timestamp string      `json:"ts"`
	ReadOnly  bool        `json:"ro"`
	SysMes    bool        `json:"sysMes"`
	UpdatedAt string      `json:"_updatedAt"`
}

type ChannelUser struct {
	ID       string `json:"_id"`
	Username string `json:"username"`
}

type ChannelInfoEnv struct {
	Channel ChannelInfo `json:"channel"`
	Success bool        `json:"success"`
}

type ChannelInfo struct {
	ID        string   `json:"_id"`
	Timestamp string   `json:"ts"`
	Type      string   `json:"t"`
	Name      string   `json:"name"`
	Usernames []string `json:"usernames"`
	Msgs      int64    `json:"msgs"`
	Default   bool     `json:"default"`
	UpdatedAt string   `json:"_updatedAt"`
	LM        string   `json:"lm"`
}

type ChannelHistory struct {
	Messages []ChannelMessage `json:"messages"`
	Success  bool             `json:"success"`
}

type ChannelMessage struct {
	ID          string      `json:"_id"`
	RoomID      string      `json:"rid"`
	Msg         string      `json:"msg"`
	Timestamp   string      `json:"ts"`
	User        ChannelUser `json:"u"`
	UpdatedAt   string      `json:"_updatedAt"`
	ChannelType *string     `json:"t,omitempty"`
	Groupable   *bool       `json:"groupable,omitempty"`
}

type ChannelCounters struct {
	Joined       bool   `json:"joined"`
	Members      int64  `json:"members"`
	Unreads      int64  `json:"unreads"`
	UnreadsFrom  string `json:"unreadsFrom"`
	Msgs         int64  `json:"msgs"`
	Latest       string `json:"latest"`
	UserMentions int64  `json:"userMentions"`
	Success      bool   `json:"success"`
}

type ChannelMembers struct {
	Members []ChannelMember `json:"members"`
	Count   int64           `json:"count"`
	Offset  int64           `json:"offset"`
	Total   int64           `json:"total"`
	Success bool            `json:"success"`
}

type ChannelMember struct {
	ID       string `json:"_id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Status   string `json:"status"`
}

type ChannelOnline struct {
	Online  []ChannelUser `json:"online"`
	Success bool          `json:"success"`
}

type ChannelRoles struct {
	Roles   []ChannelRole `json:"roles"`
	Success bool          `json:"success"`
}

type ChannelRole struct {
	RoomID string        `json:"rid"`
	User   ChannelMember `json:"u"`
	Roles  []string      `json:"roles"`
	ID     string        `json:"_id"`
}

func (c *Client) GetChannelList() (*ChannelList, error) {
	chlist := &ChannelList{}
	if err := c.c.get("/channels.list", nil).JSON(chlist); err != nil {
		return nil, err
	}

	return chlist, nil
}

func (c *Client) GetChannelInfo(roomID string) (*ChannelInfo, error) {
	chinfo := &ChannelInfoEnv{}
	q := query("roomId", roomID)
	if err := c.c.get("/channels.info", q.Q()).JSON(chinfo); err != nil {
		return nil, err
	}

	ch := chinfo.Channel
	return &ch, nil
}

func (c *Client) GetChannelCounters(roomID string) (*ChannelCounters, error) {
	chcounters := &ChannelCounters{}
	q := query("roomId", roomID)
	if err := c.c.get("/channels.counters", q.Q()).JSON(chcounters); err != nil {
		return nil, err
	}

	return chcounters, nil
}

func (c *Client) GetChannelMembers(roomID string) (*ChannelMembers, error) {
	chmembers := &ChannelMembers{}
	q := query("roomId", roomID)
	if err := c.c.get("/channels.members", q.Q()).JSON(chmembers); err != nil {
		return nil, err
	}

	return chmembers, nil
}

func (c *Client) GetChannelOnline(roomID string) ([]ChannelUser, error) {
	chonline := &ChannelOnline{}
	vals := url.Values{}
	if roomID != "*" && roomID != "" {
		vals = query("_id", roomID).Q()
	}

	if err := c.c.get("/channels.online", vals).JSON(chonline); err != nil {
		return nil, err
	}

	o := chonline.Online
	return o, nil
}

func (c *Client) GetChannelRoles(roomID string) ([]ChannelRole, error) {
	chroles := &ChannelRoles{}
	q := query("roomId", roomID)
	if err := c.c.get("/channels.roles", q.Q()).JSON(chroles); err != nil {
		return nil, err
	}

	r := chroles.Roles
	return r, nil
}

func (c *Client) GetChannelHistory(q HistoryQuery) (*ChannelHistory, error) {
	chhistory := &ChannelHistory{}
	if err := c.c.get("/channels.history", q.Q()).JSON(chhistory); err != nil {
		return nil, err
	}

	return chhistory, nil
}

type HistoryQuery struct {
	RoomID         string
	Latest         *time.Time
	Oldest         *time.Time
	Inclusive      bool
	Offset         int
	Count          int
	IncludeUnreads bool
}

func (h *HistoryQuery) Q() url.Values {
	vals := url.Values{}
	vals["roomId"] = []string{h.RoomID}
	if h.Latest != nil {
		vals["latest"] = []string{h.Latest.Format(TimeFormat)}
	}

	if h.Oldest != nil {
		vals["oldest"] = []string{h.Oldest.Format(TimeFormat)}
	}

	if h.Inclusive {
		vals["inclusive"] = []string{"true"}
	}

	if h.Offset > 0 {
		vals["offset"] = []string{strconv.Itoa(h.Offset)}
	}

	if h.Count > 0 {
		vals["count"] = []string{strconv.Itoa(h.Count)}
	}

	if h.IncludeUnreads {
		vals["unreads"] = []string{"true"}
	}

	return vals
}
