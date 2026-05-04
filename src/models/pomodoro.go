package models

import "time"

const (
	defaultPomodoroWorkMinutes       = 25
	defaultPomodoroShortBreakMinutes = 5
	defaultPomodoroLongBreakMinutes  = 15
	defaultPomodoroLongBreakEvery    = 4
)

// PomodoroState stores timer phase
type PomodoroState int

const (
	PomodoroIdle PomodoroState = iota
	PomodoroWork
	PomodoroShortBreak
	PomodoroLongBreak
	PomodoroPaused
)

// PomodoroConfig stores session lengths in minutes
type PomodoroConfig struct {
	WorkDuration           int
	ShortBreakDuration     int
	LongBreakDuration      int
	SessionsUntilLongBreak int
}

// NewDefaultPomodoroConfig creates default Pomodoro config
func NewDefaultPomodoroConfig() *PomodoroConfig {
	return &PomodoroConfig{
		WorkDuration:           defaultPomodoroWorkMinutes,
		ShortBreakDuration:     defaultPomodoroShortBreakMinutes,
		LongBreakDuration:      defaultPomodoroLongBreakMinutes,
		SessionsUntilLongBreak: defaultPomodoroLongBreakEvery,
	}
}

// PomodoroTimer stores Pomodoro runtime state
type PomodoroTimer struct {
	Config            *PomodoroConfig
	State             PomodoroState
	StateBeforePause  PomodoroState
	TimeRemaining     time.Duration
	SessionsCompleted int
	StartTime         time.Time
	PausedAt          time.Time
}

// NewPomodoroTimer creates timer with safe config
func NewPomodoroTimer(config *PomodoroConfig) *PomodoroTimer {
	return &PomodoroTimer{
		Config:           normalizedPomodoroConfig(config),
		State:            PomodoroIdle,
		StateBeforePause: PomodoroIdle,
	}
}

// Start begins work session
func (pt *PomodoroTimer) Start() {
	config := pt.config()

	pt.State = PomodoroWork
	pt.StateBeforePause = PomodoroIdle
	pt.TimeRemaining = durationMinutes(config.WorkDuration)
	pt.StartTime = time.Now()
	pt.PausedAt = time.Time{}
}

// Pause pauses current session
func (pt *PomodoroTimer) Pause() {
	if !isActivePomodoroState(pt.State) {
		return
	}

	pt.StateBeforePause = pt.State
	pt.State = PomodoroPaused
	pt.PausedAt = time.Now()
}

// Resume resumes paused session
func (pt *PomodoroTimer) Resume() {
	if pt.State != PomodoroPaused {
		return
	}

	if !pt.PausedAt.IsZero() && !pt.StartTime.IsZero() {
		pauseDuration := time.Since(pt.PausedAt)
		pt.StartTime = pt.StartTime.Add(pauseDuration)
	}

	// Возвращаем именно ту фазу, где нажали паузу
	state := pt.StateBeforePause
	if !isActivePomodoroState(state) {
		state = PomodoroWork
	}

	pt.State = state
	pt.StateBeforePause = PomodoroIdle
	pt.PausedAt = time.Time{}
}

// Reset returns timer to idle state
func (pt *PomodoroTimer) Reset() {
	pt.State = PomodoroIdle
	pt.StateBeforePause = PomodoroIdle
	pt.TimeRemaining = 0
	pt.SessionsCompleted = 0
	pt.StartTime = time.Time{}
	pt.PausedAt = time.Time{}
}

// StartBreak starts next break session
func (pt *PomodoroTimer) StartBreak() {
	config := pt.config()
	pt.SessionsCompleted++

	if useLongBreak(pt.SessionsCompleted, config.SessionsUntilLongBreak) {
		pt.State = PomodoroLongBreak
		pt.TimeRemaining = durationMinutes(config.LongBreakDuration)
	} else {
		pt.State = PomodoroShortBreak
		pt.TimeRemaining = durationMinutes(config.ShortBreakDuration)
	}

	pt.StateBeforePause = PomodoroIdle
	pt.StartTime = time.Now()
	pt.PausedAt = time.Time{}
}

// Update refreshes timer state
func (pt *PomodoroTimer) Update() {
	if pt.State == PomodoroIdle || pt.State == PomodoroPaused {
		return
	}

	elapsed := time.Since(pt.StartTime)
	duration := pt.GetCurrentDuration()

	if elapsed >= duration {
		pt.TimeRemaining = 0
		pt.OnTimerComplete()

		return
	}

	pt.TimeRemaining = duration - elapsed
}

// GetCurrentDuration returns current phase duration
func (pt *PomodoroTimer) GetCurrentDuration() time.Duration {
	config := pt.config()

	return stateDuration(pt.durationState(), config)
}

// OnTimerComplete switches timer after phase end
func (pt *PomodoroTimer) OnTimerComplete() {
	switch pt.State {
	case PomodoroWork:
		pt.StartBreak()
	case PomodoroShortBreak, PomodoroLongBreak:
		pt.State = PomodoroIdle
		pt.StateBeforePause = PomodoroIdle
		pt.TimeRemaining = 0
		pt.StartTime = time.Time{}
		pt.PausedAt = time.Time{}
	}
}

// GetStateString returns visible state label
func (pt *PomodoroTimer) GetStateString() string {
	switch pt.State {
	case PomodoroIdle:
		return "Ready"
	case PomodoroWork:
		return "Work"
	case PomodoroShortBreak:
		return "Short Break"
	case PomodoroLongBreak:
		return "Long Break"
	case PomodoroPaused:
		return "Paused"
	default:
		return "Unknown"
	}
}

func (pt *PomodoroTimer) config() *PomodoroConfig {
	pt.Config = normalizedPomodoroConfig(pt.Config)

	return pt.Config
}

func (pt *PomodoroTimer) durationState() PomodoroState {
	if pt.State == PomodoroPaused {
		return pt.StateBeforePause
	}

	return pt.State
}

func normalizedPomodoroConfig(config *PomodoroConfig) *PomodoroConfig {
	if config == nil {
		return NewDefaultPomodoroConfig()
	}

	if config.WorkDuration <= 0 {
		config.WorkDuration = defaultPomodoroWorkMinutes
	}

	if config.ShortBreakDuration <= 0 {
		config.ShortBreakDuration = defaultPomodoroShortBreakMinutes
	}

	if config.LongBreakDuration <= 0 {
		config.LongBreakDuration = defaultPomodoroLongBreakMinutes
	}

	if config.SessionsUntilLongBreak <= 0 {
		config.SessionsUntilLongBreak = defaultPomodoroLongBreakEvery
	}

	return config
}

func isActivePomodoroState(state PomodoroState) bool {
	return state == PomodoroWork ||
		state == PomodoroShortBreak ||
		state == PomodoroLongBreak
}

func stateDuration(state PomodoroState, config *PomodoroConfig) time.Duration {
	switch state {
	case PomodoroWork:
		return durationMinutes(config.WorkDuration)
	case PomodoroShortBreak:
		return durationMinutes(config.ShortBreakDuration)
	case PomodoroLongBreak:
		return durationMinutes(config.LongBreakDuration)
	default:
		return 0
	}
}

func useLongBreak(sessionsCompleted, sessionsUntilLongBreak int) bool {
	return sessionsUntilLongBreak > 0 && sessionsCompleted%sessionsUntilLongBreak == 0
}

func durationMinutes(minutes int) time.Duration {
	return time.Duration(minutes) * time.Minute
}
