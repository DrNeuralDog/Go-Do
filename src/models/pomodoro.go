package models

import "time"

// PomodoroState represents the current state of the pomodoro timer
type PomodoroState int

const (
	PomodoroIdle PomodoroState = iota
	PomodoroWork
	PomodoroShortBreak
	PomodoroLongBreak
	PomodoroPaused
)

// PomodoroConfig holds the configuration for pomodoro sessions
type PomodoroConfig struct {
	WorkDuration      int // in minutes
	ShortBreakDuration int // in minutes
	LongBreakDuration  int // in minutes
	SessionsUntilLongBreak int // number of work sessions before long break
}

// NewDefaultPomodoroConfig creates a default configuration
func NewDefaultPomodoroConfig() *PomodoroConfig {
	return &PomodoroConfig{
		WorkDuration:           25,
		ShortBreakDuration:     5,
		LongBreakDuration:      15,
		SessionsUntilLongBreak: 4,
	}
}

// PomodoroTimer manages the pomodoro timer state 
type PomodoroTimer struct {
	Config           *PomodoroConfig
	State            PomodoroState
	TimeRemaining    time.Duration
	SessionsCompleted int
	StartTime        time.Time
	PausedAt         time.Time
}

// NewPomodoroTimer creates a new pomodoro timer
func NewPomodoroTimer(config *PomodoroConfig) *PomodoroTimer {
	return &PomodoroTimer{
		Config:           config,
		State:            PomodoroIdle,
		TimeRemaining:    0,
		SessionsCompleted: 0,
	}
}

// Start begins a new pomodoro work session
func (pt *PomodoroTimer) Start() {
	pt.State = PomodoroWork
	pt.TimeRemaining = time.Duration(pt.Config.WorkDuration) * time.Minute
	pt.StartTime = time.Now()
}

// Pause pauses the current timer
func (pt *PomodoroTimer) Pause() {
	if pt.State != PomodoroIdle && pt.State != PomodoroPaused {
		pt.State = PomodoroPaused
		pt.PausedAt = time.Now()
	}
}

// Resume resumes a paused timer
func (pt *PomodoroTimer) Resume() {
	if pt.State == PomodoroPaused {
		// Calculate pause duration and adjust start time
		pauseDuration := time.Since(pt.PausedAt)
		pt.StartTime = pt.StartTime.Add(pauseDuration)
		pt.State = PomodoroWork // Resume to work state
	}
}

// Reset resets the timer to idle state
func (pt *PomodoroTimer) Reset() {
	pt.State = PomodoroIdle
	pt.TimeRemaining = 0
	pt.SessionsCompleted = 0
}

// StartBreak starts a break session
func (pt *PomodoroTimer) StartBreak() {
	pt.SessionsCompleted++

	if pt.SessionsCompleted % pt.Config.SessionsUntilLongBreak == 0 {
		// Long break
		pt.State = PomodoroLongBreak
		pt.TimeRemaining = time.Duration(pt.Config.LongBreakDuration) * time.Minute
	} else {
		// Short break
		pt.State = PomodoroShortBreak
		pt.TimeRemaining = time.Duration(pt.Config.ShortBreakDuration) * time.Minute
	}
	pt.StartTime = time.Now()
}

// Update updates the timer state
func (pt *PomodoroTimer) Update() {
	if pt.State == PomodoroIdle || pt.State == PomodoroPaused {
		return
	}

	elapsed := time.Since(pt.StartTime)
	duration := pt.GetCurrentDuration()

	if elapsed >= duration {
		// Timer completed
		pt.OnTimerComplete()
	} else {
		pt.TimeRemaining = duration - elapsed
	}
}

// GetCurrentDuration returns the total duration for current state
func (pt *PomodoroTimer) GetCurrentDuration() time.Duration {
	switch pt.State {
	case PomodoroWork:
		return time.Duration(pt.Config.WorkDuration) * time.Minute
	case PomodoroShortBreak:
		return time.Duration(pt.Config.ShortBreakDuration) * time.Minute
	case PomodoroLongBreak:
		return time.Duration(pt.Config.LongBreakDuration) * time.Minute
	default:
		return 0
	}
}

// OnTimerComplete handles timer completion
func (pt *PomodoroTimer) OnTimerComplete() {
	switch pt.State {
	case PomodoroWork:
		pt.StartBreak()
	case PomodoroShortBreak, PomodoroLongBreak:
		pt.State = PomodoroIdle
		pt.TimeRemaining = 0
	}
}

// GetStateString returns a human-readable state string
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
