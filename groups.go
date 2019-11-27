package rc

import "time"

type (
	GroupList struct {
		Groups []struct {
			ID    string `json:"_id"`
			Name  string `json:"name"`
			Fname string `json:"fname"`
			T     string `json:"t"`
			Msgs  int    `json:"msgs"`
			U     struct {
				ID       string `json:"_id"`
				Username string `json:"username"`
			} `json:"u"`
			CustomFields struct {
				CompanyID string `json:"companyId"`
			} `json:"customFields"`
			Ts        time.Time `json:"ts"`
			Ro        bool      `json:"ro"`
			SysMes    bool      `json:"sysMes"`
			UpdatedAt time.Time `json:"_updatedAt"`
		} `json:"groups"`
		Offset  int  `json:"offset"`
		Count   int  `json:"count"`
		Total   int  `json:"total"`
		Success bool `json:"success"`
	}

	GroupMembers struct {
		Members []struct {
			ID        string  `json:"_id"`
			Status    string  `json:"status"`
			Name      string  `json:"name"`
			UtcOffset float64 `json:"utcOffset"`
			Username  string  `json:"username"`
		} `json:"members"`
		Count   int  `json:"count"`
		Offset  int  `json:"offset"`
		Total   int  `json:"total"`
		Success bool `json:"success"`
	}
)

func (c *Client) GetGroupList() (*GroupList, error) {
	glist := &GroupList{}
	if err := c.c.get("/groups.listAll", nil).JSON(glist); err != nil {
		return nil, err
	}

	return glist, nil
}

func (c *Client) GetGroupMembers(roomID string) (*GroupMembers, error) {
	gmembers := &GroupMembers{}
	q := query("roomId", roomID)
	if err := c.c.get("/groups.members", q.Q()).JSON(gmembers); err != nil {
		return nil, err
	}

	return gmembers, nil
}

func (c *Client) GetGroupMembersByRoomName(roomName string) (*GroupMembers, error) {
	gmembers := &GroupMembers{}
	q := query("roomName", roomName)
	if err := c.c.get("/groups.members", q.Q()).JSON(gmembers); err != nil {
		return nil, err
	}

	return gmembers, nil
}
