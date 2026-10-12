package model

type Menu struct {
	BaseModel

	Name    string `json:"name"`
	Pid     int64  `json:"pid"`
	Path    string `json:"path"`
	Status  int8   `json:"status"`
	Display int8   `json:"display"`
	Sort    int8   `json:"sort"`
	Desc    string `json:"desc"`

	CommonTimestampsField
}

const Hidden = 0
const Display = 1
