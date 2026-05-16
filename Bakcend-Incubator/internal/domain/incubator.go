package domain

type IncubatorState struct {
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Mode        string  `json:"mode"`
	Lamp        bool    `json:"lamp"`
}

type IncubatorRepository interface {
	PublishMode(mode string) error
	PublishLamp(state bool) error
}

type IncubatorUsecase interface {
	GetState() IncubatorState
	UpdateState(state IncubatorState)
	SetMode(mode string) error
	SetLamp(on bool) error
}

type CommandMode struct {
	Mode string `json:"mode"`
}

type CommandLamp struct {
	Lamp bool `json:"lamp"`
}
