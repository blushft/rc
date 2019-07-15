package rc

type IntegrationEvent string

const (
	EvtSendMessage  IntegrationEvent = "sendMessage"
	EvtFileUploaded IntegrationEvent = "fileUploaded"
	EvtRoomArchived IntegrationEvent = "roomArchived"
	EvtRoomCreated  IntegrationEvent = "roomCreated"
	EvtRoomJoined   IntegrationEvent = "roomJoined"
	EvtRoomLeft     IntegrationEvent = "roomLeft"
	EvtUserCreated  IntegrationEvent = "userCreated"
)

type Integration struct {
	Type          string           `json:"type,omitempty"`
	Name          string           `json:"name,omitempty"`
	Event         IntegrationEvent `json:"event,omitempty"`
	Enabled       bool             `json:"enabled,omitempty"`
	Username      string           `json:"username,omitempty"`
	Urls          []string         `json:"urls,omitempty"`
	ScriptEnabled bool             `json:"scriptEnabled,omitempty"`
	Channel       string           `json:"channel,omitempty"`
	TriggerWords  string           `json:"trigger_words,omitempty"`
	Alias         string           `json:"alias,omitempty"`
	Avatar        string           `json:"avatar,omitempty"`
	Emoji         string           `json:"emoji,omitempty"`
	Token         string           `json:"token,omitempty"`
	Script        string           `json:"script,omitempty"`
}

type IntegrationOption func(*Integration)

func IntegrationScript(script string) IntegrationOption {
	return func(i *Integration) {
		i.ScriptEnabled = true
		i.Script = script
	}
}

func IntegrationAvatarURL(u string) IntegrationOption {
	return func(i *Integration) {
		i.Avatar = u
	}
}

func IntegrationAvatarEmoji(e string) IntegrationOption {
	return func(i *Integration) {
		i.Emoji = e
	}
}

func IntegrationTriggers(t string) IntegrationOption {
	return func(i *Integration) {
		i.TriggerWords = t
	}
}

func IntegrationToken(t string) IntegrationOption {
	return func(i *Integration) {
		i.Token = t
	}
}

func NewIncomingIntegration(name, user, channel string, enabled bool, opts ...IntegrationOption) *Integration {
	i := &Integration{
		Type:          "webhook-incoming",
		Name:          name,
		Username:      user,
		Enabled:       enabled,
		Channel:       channel,
		ScriptEnabled: false,
	}

	for _, o := range opts {
		o(i)
	}

	return i
}

func NewOutgoingIntegration(name, user, channel string, event IntegrationEvent, urls []string, enabled bool, opts ...IntegrationOption) *Integration {
	i := &Integration{
		Type:          "webhook-outgoing",
		Name:          name,
		Event:         event,
		Username:      user,
		Enabled:       enabled,
		Channel:       channel,
		ScriptEnabled: false,
		Urls:          urls,
	}

	for _, o := range opts {
		o(i)
	}

	return i
}

type IntegrationResponse struct {
	Integration IntegrationInfo `json:"integration"`
	Success     bool            `json:"success"`
}

type IntegrationList struct {
	Integrations []IntegrationInfo `json:"integrations,omitempty"`
	Success      bool              `json:"success,omitempty"`
	Offset       int               `json:"offset,omitempty"`
	Items        int               `json:"items,omitempty"`
	Total        int               `json:"total,omitempty"`
}

type IntegrationInfo struct {
	Type          string               `json:"type"`
	Name          string               `json:"name"`
	Enabled       bool                 `json:"enabled"`
	Username      string               `json:"username"`
	Event         string               `json:"event"`
	Urls          []string             `json:"urls"`
	ScriptEnabled bool                 `json:"scriptEnabled"`
	UserID        string               `json:"userId"`
	Channel       []interface{}        `json:"channel"`
	CreatedAt     string               `json:"_createdAt"`
	CreatedBy     IntegrationCreatedBy `json:"_createdBy"`
	UpdatedAt     string               `json:"_updatedAt"`
	ID            string               `json:"_id"`
}

type IntegrationCreatedBy struct {
	Username string `json:"username"`
	ID       string `json:"_id"`
}

func (c *Client) CreateIntegration(i *Integration) (*IntegrationInfo, error) {
	result := &IntegrationResponse{}
	if err := c.c.postJSON("/integrations.create", i).JSON(result); err != nil {
		return nil, err
	}
	info := result.Integration
	return &info, nil
}

func (c *Client) GetIntegrations() ([]IntegrationInfo, error) {
	is := &IntegrationList{}
	if err := c.c.get("/integrations.list", nil).JSON(is); err != nil {
		return nil, err
	}
	result := is.Integrations
	return result, nil
}
