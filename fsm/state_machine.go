package fsm

import "fmt"

const (
	NoSuchEventTransitionError = "failed to [%s] from state [%s]: no such event"
	AmbiguousTransitionError   = "failed to [%s] from state [%s]: ambiguous transitions"
	NoMatchedTransitionError   = "failed to [%s] from state [%s]: no matched transition"
	IllegalStateCodeError      = "illegal state code: %v"
)

type StatusCode int32

const (
	NilState StatusCode = 0
)

type Stateful struct {
	State StatusCode `xorm:"state" json:"state"`
}

func (status *Stateful) SetState(code StatusCode) {
	status.State = code
}

func (status *Stateful) GetState() StatusCode {
	return status.State
}

func (status *Stateful) GetStateName() string {
	return string(status.State)
}

func IllegalStateCode(code StatusCode) error {
	return fmt.Errorf(IllegalStateCodeError, code)
}

type Stater interface {
	SetState(code StatusCode)
	GetState() StatusCode
	GetStateName() string
}

type StateMachine struct {
	initialState StatusCode
	states       map[StatusCode]*State
	events       map[EventCode]*Event
}

func New() *StateMachine {
	return &StateMachine{
		states: map[StatusCode]*State{},
		events: map[EventCode]*Event{},
	}
}

func (sm *StateMachine) Initial(code StatusCode) *StateMachine {
	sm.initialState = code
	return sm
}

func (sm *StateMachine) State(code StatusCode) *State {
	state := &State{code: code}
	sm.states[code] = state
	return state
}

func (sm *StateMachine) States(states ...StatusCode) {
	for state := range states {
		sm.State(StatusCode(state))
	}
}

func (sm *StateMachine) Event(code EventCode) *Event {
	event := &Event{code: code}
	sm.events[code] = event
	return event
}

func (sm *StateMachine) Trigger(eventCode EventCode, model Stater, notes ...string) error {
	return sm.trigger(eventCode, model, false, notes...)
}

func (sm *StateMachine) TriggerOn(eventCode EventCode, model Stater, notes ...string) error {
	return sm.trigger(eventCode, model, true, notes...)
}

func (sm *StateMachine) trigger(eventCode EventCode, model Stater, continuable bool, notes ...string) error {
	if model.GetState() == NilState {
		model.SetState(sm.initialState)
	}
	var event *Event
	if event = sm.events[eventCode]; event == nil {
		return fmt.Errorf(NoSuchEventTransitionError, eventCode, model.GetStateName())
	}
	var (
		originState     StatusCode
		originStateName string
	)
	for {
		originState = model.GetState()
		originStateName = model.GetStateName()

		matchedTransitions := matchTransitions(event, originState)
		if len(matchedTransitions) == 0 {
			return fmt.Errorf(NoMatchedTransitionError, eventCode, originStateName)
		} else if len(matchedTransitions) > 1 {
			return fmt.Errorf(AmbiguousTransitionError, eventCode, originStateName)
		}

		transition := matchedTransitions[0]

		// State: exit
		if state, ok := sm.states[originState]; ok {
			for _, exit := range state.exits {
				if err := exit(model); err != nil {
					return err
				}
			}
		}

		// Transition: before
		for _, before := range transition.befores {
			if err := before(model); err != nil {
				return err
			}
		}

		model.SetState(transition.to)

		// State: enter
		if state, ok := sm.states[transition.to]; ok {
			for _, enter := range state.enters {
				if err := enter(model); err != nil {
					model.SetState(originState)
					return err
				}
			}
		}

		// Transition: after
		for _, after := range transition.afters {
			if err := after(model); err != nil {
				model.SetState(originState)
				return err
			}
		}

		if continuable && transition.continuable {
			if originState != model.GetState() {
				continue
			}
		}
		return nil
	}
}

func matchTransitions(event *Event, originState StatusCode) []*Transition {
	var matchedTransitions []*Transition
	for _, transition := range event.transitions {
		var validFrom = len(transition.froms) == 0
		if len(transition.froms) > 0 {
			for _, from := range transition.froms {
				if from == originState {
					validFrom = true
				}
			}
		}
		if validFrom {
			matchedTransitions = append(matchedTransitions, transition)
		}
	}
	return matchedTransitions
}

type State struct {
	code   StatusCode
	enters []Handler
	exits  []Handler
}

type Handler func(model interface{}) error

func (state *State) Enter(handler Handler) *State {
	state.enters = append(state.enters, handler)
	return state
}

func (state *State) Exit(handler Handler) *State {
	state.exits = append(state.exits, handler)
	return state
}

type EventCode string

type Event struct {
	code        EventCode
	transitions []*Transition
}

func (event *Event) To(code StatusCode) *Transition {
	transition := &Transition{to: code, continuable: true}
	event.transitions = append(event.transitions, transition)
	return transition
}

type Transition struct {
	to          StatusCode
	froms       []StatusCode
	befores     []Handler
	afters      []Handler
	continuable bool
}

func (transition *Transition) From(states ...StatusCode) *Transition {
	transition.froms = states
	return transition
}

func (transition *Transition) Before(handler Handler) *Transition {
	transition.befores = append(transition.befores, handler)
	return transition
}

func (transition *Transition) After(handler Handler) *Transition {
	transition.afters = append(transition.afters, handler)
	return transition
}

func (transition *Transition) Continuable() *Transition {
	transition.continuable = true
	return transition
}
