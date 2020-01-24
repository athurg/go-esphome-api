package esphome

import (
	"image/color"
	"math"

	"maze.io/esphome/api"
)

type Entity struct {
	Name     string
	ObjectID string
	UniqueID string
	Key      uint32
	client   *Client
}

// Entities is a high level map of a device's entities.
type Entities struct {
	BinarySensor map[string]BinarySensor
	Camera       map[string]Camera
	Climate      map[string]Climate
	Cover        map[string]Cover
	Fan          map[string]Fan
	Light        map[string]Light
	Sensor       map[string]Sensor
	Switch       map[string]Switch
	TextSensor   map[string]TextSensor
}

type BinarySensor struct {
	Entity
	DeviceClass string
	State       bool
}

func newBinarySensor(client *Client, entity *api.ListEntitiesBinarySensorResponse) *BinarySensor {
	return &BinarySensor{
		Entity: Entity{
			Name:     entity.Name,
			ObjectID: entity.ObjectId,
			UniqueID: entity.UniqueId,
			Key:      entity.Key,
			client:   client,
		},
		DeviceClass: entity.DeviceClass,
	}
}

type Climate struct {
	Entity

	Capabilities ClimateCapabilities
}

type (
	ClimateCapabilities struct {
		CurrentTemperature        bool
		TwoPointTargetTemperature bool
		Modes                     []ClimateMode
		VisualMinTemperature      float32
		VisualMaxTemperature      float32
		VisualTemperatureStep     float32
		Away                      bool
		Action                    bool
		FanModes                  []ClimateFanMode
		SwingModes                []ClimateSwingMode
	}
	ClimateMode      int32
	ClimateFanMode   int32
	ClimateSwingMode int32
)

const (
	ClimateModeOff ClimateMode = iota
	ClimateModeAuto
	ClimateModeCool
	ClimateModeHeat
	ClimateModeFanOnly
	ClimateModeDry
)

const (
	ClimateFanModeOn ClimateFanMode = iota
	ClimateFanModeOff
	ClimateFanModeAuto
	ClimateFanModeLow
	ClimateFanModeMedium
	ClimateFanModeHigh
	ClimateFanModeMiddle
	ClimateFanModeFocus
	ClimateFanModeDiffuse
)

const (
	ClimateSwingModeOff ClimateSwingMode = iota
	ClimateSwingModeBoth
	ClimateSwingModeVertical
	ClimateSwingModeHorizontal
)

func newClimate(client *Client, entity *api.ListEntitiesClimateResponse) *Climate {
	var (
		modes      = make([]ClimateMode, len(entity.SupportedModes))
		fanModes   = make([]ClimateFanMode, len(entity.SupportedFanModes))
		swingModes = make([]ClimateSwingMode, len(entity.SupportedSwingModes))
	)
	for i, v := range entity.SupportedModes {
		modes[i] = ClimateMode(v)
	}
	for i, v := range entity.SupportedFanModes {
		fanModes[i] = ClimateFanMode(v)
	}
	for i, v := range entity.SupportedSwingModes {
		swingModes[i] = ClimateSwingMode(v)
	}
	return &Climate{
		Entity: Entity{
			Name:     entity.Name,
			ObjectID: entity.ObjectId,
			UniqueID: entity.UniqueId,
			Key:      entity.Key,
			client:   client,
		},
		Capabilities: ClimateCapabilities{
			CurrentTemperature:        entity.SupportsCurrentTemperature,
			TwoPointTargetTemperature: entity.SupportsTwoPointTargetTemperature,
			Modes:                     modes,
			VisualMinTemperature:      entity.VisualMinTemperature,
			VisualMaxTemperature:      entity.VisualMaxTemperature,
			VisualTemperatureStep:     entity.VisualTemperatureStep,
			Away:                      entity.SupportsAway,
			Action:                    entity.SupportsAction,
			FanModes:                  fanModes,
			SwingModes:                swingModes,
		},
	}
}

type Cover struct {
	Entity

	// TODO(maze): Finish me
}

func newCover(client *Client, entity *api.ListEntitiesCoverResponse) *Cover {
	return &Cover{
		Entity: Entity{
			Name:     entity.Name,
			ObjectID: entity.ObjectId,
			UniqueID: entity.UniqueId,
			Key:      entity.Key,
			client:   client,
		},
	}
}

type Fan struct {
	Entity

	// TODO(maze): Finish me
}

func newFan(client *Client, entity *api.ListEntitiesFanResponse) *Fan {
	return &Fan{
		Entity: Entity{
			Name:     entity.Name,
			ObjectID: entity.ObjectId,
			UniqueID: entity.UniqueId,
			Key:      entity.Key,
			client:   client,
		},
	}
}

type Light struct {
	Entity

	Capabilities LightCapabilities
	Effects      []string

	State        LightState
	StateIsValid bool

	HandleState            func(on bool)
	HandleBrightness       func(float32)
	HandleColor            func(r, g, b, w float32)
	HandleColorTemperature func(float32)
	HandleEffect           func(string)
}

type LightCapabilities struct {
	Brightness       bool
	RGB              bool
	WhiteValue       bool
	ColorTemperature bool
	MinMired         float32
	MaxMired         float32
}

type LightState struct {
	On                      bool
	Brightness              float32
	Red, Green, Blue, White float32
	ColorTemperature        float32
	Effect                  string
}

func newLight(client *Client, entity *api.ListEntitiesLightResponse) *Light {
	effects := make([]string, len(entity.Effects))
	copy(effects, entity.Effects)
	return &Light{
		Entity: Entity{
			Name:     entity.Name,
			ObjectID: entity.ObjectId,
			UniqueID: entity.UniqueId,
			Key:      entity.Key,
			client:   client,
		},
		Capabilities: LightCapabilities{
			Brightness:       entity.SupportsBrightness,
			RGB:              entity.SupportsRgb,
			WhiteValue:       entity.SupportsWhiteValue,
			ColorTemperature: entity.SupportsColorTemperature,
			MinMired:         entity.MinMireds,
			MaxMired:         entity.MaxMireds,
		},
		Effects: effects,
	}
}

func (entity Light) commandRequest() *api.LightCommandRequest {
	return &api.LightCommandRequest{
		Key:                 entity.Key,
		State:               entity.State.On,
		HasBrightness:       entity.Capabilities.Brightness,
		Brightness:          entity.State.Brightness,
		HasRgb:              entity.Capabilities.RGB,
		Red:                 entity.State.Red,
		Green:               entity.State.Green,
		Blue:                entity.State.Blue,
		HasWhite:            entity.Capabilities.WhiteValue,
		White:               entity.State.White,
		HasColorTemperature: entity.Capabilities.ColorTemperature,
		ColorTemperature:    entity.State.ColorTemperature,
		HasTransitionLength: false,
		TransitionLength:    0,
		HasFlashLength:      false,
		FlashLength:         0,
		HasEffect:           false,
		Effect:              entity.State.Effect,
	}
}

func (entity *Light) update(state *api.LightStateResponse) {
	if entity.HandleState != nil {
		if state.State != entity.State.On {
			entity.HandleState(true)
		}
	}
	if entity.HandleColor != nil {
		if !entity.StateIsValid {
			entity.HandleColor(state.Red, state.Green, state.Blue, state.White)
		} else if !equal(entity.State.Red, state.Red) ||
			!equal(entity.State.Green, state.Green) ||
			!equal(entity.State.Blue, state.Blue) ||
			!equal(entity.State.White, state.White) {
			entity.HandleColor(state.Red, state.Green, state.Blue, state.White)
		}
	}
	if entity.HandleColorTemperature != nil {
		if !entity.StateIsValid || !equal(entity.State.ColorTemperature, state.ColorTemperature) {
			entity.HandleColorTemperature(state.ColorTemperature)
		}
	}
	if entity.HandleEffect != nil {
		if !entity.StateIsValid || entity.State.Effect != state.Effect {
			entity.HandleEffect(state.Effect)
		}
	}

	entity.State.On = state.State
	entity.State.Brightness = state.Brightness
	entity.State.Red = state.Red
	entity.State.Green = state.Green
	entity.State.Blue = state.Blue
	entity.State.White = state.White
	entity.State.ColorTemperature = state.ColorTemperature
	entity.State.Effect = state.Effect
	entity.StateIsValid = true
}

func (entity Light) SetBrightness(value float32) error {
	request := entity.commandRequest()
	request.Brightness = value
	return entity.client.sendTimeout(request, entity.client.Timeout)
}

func (entity Light) SetColor(value color.Color) error {
	r, g, b, _ := value.RGBA()
	request := entity.commandRequest()
	request.Red = float32(r>>4) / 256.0
	request.Green = float32(g>>4) / 256.0
	request.Blue = float32(b>>4) / 256.0
	return entity.client.sendTimeout(request, entity.client.Timeout)
}

func (entity Light) SetWhite(value float32) error {
	request := entity.commandRequest()
	request.White = value
	return entity.client.sendTimeout(request, entity.client.Timeout)
}

func (entity Light) SetEffect(effect string) error {
	request := entity.commandRequest()
	request.Effect = effect
	return entity.client.sendTimeout(request, entity.client.Timeout)
}

func (entity Light) SetState(on bool) error {
	request := entity.commandRequest()
	request.State = on
	return entity.client.sendTimeout(request, entity.client.Timeout)
}

type Sensor struct {
	Entity
	Icon              string
	UnitOfMeasurement string
	AccuracyDecimals  int32
	ForceUpdate       bool

	State        float32
	StateIsValid bool

	HandleState func(float32)
}

func newSensor(client *Client, entity *api.ListEntitiesSensorResponse) *Sensor {
	return &Sensor{
		Entity: Entity{
			Name:     entity.Name,
			ObjectID: entity.ObjectId,
			UniqueID: entity.UniqueId,
			Key:      entity.Key,
			client:   client,
		},
	}
}

func (entity *Sensor) update(state *api.SensorStateResponse) {
	if !state.MissingState && entity.HandleState != nil && !equal(entity.State, state.State) {
		entity.HandleState(state.State)
	}

	entity.State = state.State
	entity.StateIsValid = !state.MissingState
}

type Switch struct {
	Entity
	Icon         string
	AssumedState bool

	// State of the switch.
	State bool

	HandleState func(bool)
}

func newSwitch(client *Client, entity *api.ListEntitiesSwitchResponse) *Switch {
	return &Switch{
		Entity: Entity{
			Name:     entity.Name,
			ObjectID: entity.ObjectId,
			UniqueID: entity.UniqueId,
			Key:      entity.Key,
			client:   client,
		},
	}
}

func (entity *Switch) update(state *api.SwitchStateResponse) {
	if entity.HandleState != nil {
		entity.HandleState(state.State)
	}

	entity.State = state.State
	entity.AssumedState = false
}

func (entity Switch) commandRequest() *api.SwitchCommandRequest {
	return &api.SwitchCommandRequest{
		Key:   entity.Key,
		State: entity.State,
	}
}

func (entity Switch) SetState(on bool) error {
	request := entity.commandRequest()
	request.State = on
	return entity.client.sendTimeout(request, entity.client.Timeout)
}

type TextSensor struct {
	Entity
	State        string
	StateIsValid bool
}

func newTextSensor(client *Client, entity *api.ListEntitiesTextSensorResponse) *TextSensor {
	return &TextSensor{
		Entity: Entity{
			Name:     entity.Name,
			ObjectID: entity.ObjectId,
			UniqueID: entity.UniqueId,
			Key:      entity.Key,
			client:   client,
		},
	}
}

func equal(a, b float32) bool {
	const ε = 1e-6
	return math.Abs(float64(a)-float64(b)) <= ε
}
