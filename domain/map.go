package domain

type MapNode struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Coordinates Coordinates `json:"coordinates"`
	Endpoint    bool        `json:"endpoint"`
}

type MapConnection struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Status string `json:"status"`
}

type MapData struct {
	Nodes       []MapNode       `json:"nodes"`
	Connections []MapConnection `json:"connections"`
}
