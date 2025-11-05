package ui

import (
    "fmt"
    "image/color"
    "math"
    "time"

    "todo-list-migration/src/models"

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/canvas"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
)

// PomodoroWindow represents the Pomodoro timer window
type PomodoroWindow struct {
	window         fyne.Window
	timer          *models.PomodoroTimer
	config         *models.PomodoroConfig
	isGruvbox      bool

	// UI components
    timerCanvas    *canvas.Text
    stateCanvas    *canvas.Text
    startBtn       *SimpleRectButton
    pauseBtn       *SimpleRectButton
    resetBtn       *SimpleRectButton
    sessionsCanvas *canvas.Text
    progressRing   *ProgressRing

	// Configuration inputs
    workSpinner       *NumberSpinner
    shortBreakSpinner *NumberSpinner
    longBreakSpinner  *NumberSpinner

	// Timer animation
	anim           *fyne.Animation
	lastUpdate     time.Time
}

// NewPomodoroWindow creates a new Pomodoro timer window
func NewPomodoroWindow(app fyne.App, isGruvbox bool) *PomodoroWindow {
	config := models.NewDefaultPomodoroConfig()
	timer := models.NewPomodoroTimer(config)

	pw := &PomodoroWindow{
		window:    app.NewWindow("Pomodoro Timer"),
		timer:     timer,
		config:    config,
		isGruvbox: isGruvbox,
		lastUpdate: time.Now().Add(-time.Second),
	}

	pw.setupUI()
	pw.startTicker()
	pw.tick() // initial update

	return pw
}

// setupUI initializes the user interface
func (pw *PomodoroWindow) setupUI() {
	// Window properties
	pw.window.Resize(fyne.NewSize(400, 600))
	pw.window.CenterOnScreen()

    // Get colors based on theme (match main window)
    var bgStart, bgEnd color.Color
    var titleColor color.Color
    currentTheme := fyne.CurrentApp().Settings().Theme()
    if lightTheme, ok := currentTheme.(*LightSoftTheme); ok {
        bgStart, bgEnd = lightTheme.GetHeaderGradientColors()
        titleColor = color.White
    } else if gruv, ok := currentTheme.(*GruvboxBlackTheme); ok {
        bgStart, bgEnd = gruv.GetHeaderGradientColors()
        titleColor = hex("#fabd2f")
    } else {
        bgStart = bgEnd // fallback
        titleColor = color.White
    }

    // Timer display - big digits, themed color
    pw.timerCanvas = canvas.NewText("25:00", titleColor)
    pw.timerCanvas.TextStyle = fyne.TextStyle{Bold: true}
    pw.timerCanvas.TextSize = 77
    
    // State label
    pw.stateCanvas = canvas.NewText("Ready", titleColor)
    pw.stateCanvas.TextSize = 18
    
    // Sessions completed label
    pw.sessionsCanvas = canvas.NewText("Sessions: 0", titleColor)
    pw.sessionsCanvas.TextSize = 16

    // Control buttons (match theme button styles from main window)
    var btnBg, btnFg color.Color
    if _, ok := fyne.CurrentApp().Settings().Theme().(*LightSoftTheme); ok {
        btnBg = hex("#ff8c42")
        btnFg = color.White
    } else {
        btnBg = hex("#504945")
        btnFg = hex("#fabd2f")
    }
    pw.startBtn = NewSimpleRectButton("Start", btnBg, btnFg, fyne.NewSize(90, 36), 8, pw.onStartClicked)
    pw.pauseBtn = NewSimpleRectButton("Pause", btnBg, btnFg, fyne.NewSize(90, 36), 8, pw.onPauseClicked)
    pw.pauseBtn.Disable()
    pw.resetBtn = NewSimpleRectButton("Reset", btnBg, btnFg, fyne.NewSize(90, 36), 8, pw.onResetClicked)

	buttonRow := container.NewHBox(
		pw.startBtn,
		pw.pauseBtn,
		pw.resetBtn,
	)

    // Configuration section (labels themed, inputs as white-back spinners)
    cfgHeader := canvas.NewText("Configuration", titleColor)
    cfgHeader.TextStyle = fyne.TextStyle{Bold: true}
    cfgHeader.TextSize = 18

    labelWork := canvas.NewText("Work time (min):", titleColor)
    labelWork.TextSize = 16
    labelShort := canvas.NewText("Short break (min):", titleColor)
    labelShort.TextSize = 16
    labelLong := canvas.NewText("Long break (min):", titleColor)
    labelLong.TextSize = 16

    darkText := hex("#3c3836")
    whiteBg := color.White

    pw.workSpinner = NewNumberSpinner(pw.window, pw.config.WorkDuration, 1, 120, 1, darkText, whiteBg, func(v int) {
        pw.config.WorkDuration = v
        pw.tick()
    })
    pw.shortBreakSpinner = NewNumberSpinner(pw.window, pw.config.ShortBreakDuration, 1, 60, 1, darkText, whiteBg, func(v int) {
        pw.config.ShortBreakDuration = v
        pw.tick()
    })
    pw.longBreakSpinner = NewNumberSpinner(pw.window, pw.config.LongBreakDuration, 1, 120, 1, darkText, whiteBg, func(v int) {
        pw.config.LongBreakDuration = v
        pw.tick()
    })

    configForm := container.NewVBox(
        cfgHeader,
        widget.NewSeparator(),
        container.NewGridWithColumns(2,
            labelWork,
            pw.workSpinner,
        ),
        container.NewGridWithColumns(2,
            labelShort,
            pw.shortBreakSpinner,
        ),
        container.NewGridWithColumns(2,
            labelLong,
            pw.longBreakSpinner,
        ),
    )

    // Progress ring colors
    var tickBg color.Color
    if _, ok := currentTheme.(*LightSoftTheme); ok {
        tickBg = hex("#e0e0e0")
    } else {
        tickBg = hex("#504945")
    }
    pw.progressRing = NewProgressRing(hex("#fb4934"), tickBg)

    // Layout
    ringWithDigits := container.NewMax(
        container.NewCenter(pw.progressRing),
        container.NewCenter(pw.timerCanvas),
    )
    timerDisplay := container.NewVBox(
        CreateSpacer(1, 40),
        container.NewCenter(container.NewGridWrap(fyne.NewSize(240, 240), ringWithDigits)),
        CreateSpacer(1, 10),
        container.NewCenter(pw.stateCanvas),
        CreateSpacer(1, 20),
        container.NewCenter(pw.sessionsCanvas),
        CreateSpacer(1, 30),
        container.NewCenter(buttonRow),
        CreateSpacer(1, 40),
    )

	content := container.NewVBox(
		timerDisplay,
		widget.NewSeparator(),
		CreateSpacer(1, 20),
		configForm,
		CreateSpacer(1, 20),
	)

	// Wrap in padding
	paddedContent := container.NewBorder(
		nil, nil,
		CreateSpacer(24, 1),
		CreateSpacer(24, 1),
		content,
	)

	// Background gradient
	background := NewGradientRect(bgStart, bgEnd, 0)
	finalContent := container.NewMax(background, paddedContent)

    pw.window.SetContent(finalContent)
}

// startTicker starts the timer update ticker
func (pw *PomodoroWindow) startTicker() {
    pw.anim = fyne.NewAnimation(time.Second, func(_ float32) {
		now := time.Now()
		if now.Sub(pw.lastUpdate) >= time.Second {
			pw.lastUpdate = now
			pw.tick()
		}
	})
	pw.anim.RepeatCount = fyne.AnimationRepeatForever
	pw.anim.Start()
}

// stopTicker stops the timer update ticker
func (pw *PomodoroWindow) stopTicker() {
	if pw.anim != nil {
		pw.anim.Stop()
		pw.anim = nil
	}
}

func (pw *PomodoroWindow) tick() {
	pw.timer.Update()

	minutes := int(pw.timer.TimeRemaining.Minutes())
	seconds := int(pw.timer.TimeRemaining.Seconds()) % 60

	timerText := ""
	if pw.timer.State == models.PomodoroIdle {
		timerText = fmt.Sprintf("%02d:00", pw.config.WorkDuration)
	} else {
		timerText = fmt.Sprintf("%02d:%02d", minutes, seconds)
	}

    stateText := pw.timer.GetStateString()
    sessionsText := fmt.Sprintf("Sessions: %d", pw.timer.SessionsCompleted)

    pw.timerCanvas.Text = timerText
    pw.timerCanvas.Refresh()
    pw.stateCanvas.Text = stateText
    pw.stateCanvas.Refresh()
    pw.sessionsCanvas.Text = sessionsText
    pw.sessionsCanvas.Refresh()

    // Update progress ring (remaining fraction, clockwise from left)
    var prog float32 = 1
    total := pw.timer.GetCurrentDuration()
    // Handle paused/unknown state by estimating total from config
    if (pw.timer.State == models.PomodoroPaused || total == 0) && pw.timer.TimeRemaining > 0 {
        wrk := time.Duration(pw.config.WorkDuration) * time.Minute
        sbr := time.Duration(pw.config.ShortBreakDuration) * time.Minute
        lbr := time.Duration(pw.config.LongBreakDuration) * time.Minute
        // choose the smallest duration that is >= remaining as a plausible total
        candidates := []time.Duration{wrk, sbr, lbr}
        var best time.Duration
        for _, d := range candidates {
            if d >= pw.timer.TimeRemaining && (best == 0 || d < best) {
                best = d
            }
        }
        if best == 0 {
            best = wrk
        }
        total = best
    }
    if pw.timer.State != models.PomodoroIdle && total > 0 {
        prog = float32(pw.timer.TimeRemaining.Seconds() / total.Seconds())
        if prog < 0 {
            prog = 0
        }
        if prog > 1 {
            prog = 1
        }
    }
    if pw.progressRing != nil {
        pw.progressRing.SetProgress(prog)
    }

	// Update button states
    switch pw.timer.State {
	case models.PomodoroIdle:
        pw.startBtn.Enable()
        pw.pauseBtn.Disable()
        pw.resetBtn.Disable()
        pw.startBtn.SetText("Start")
	case models.PomodoroWork, models.PomodoroShortBreak, models.PomodoroLongBreak:
        pw.startBtn.Disable()
        pw.pauseBtn.Enable()
        pw.resetBtn.Enable()
        pw.pauseBtn.SetText("Pause")
	case models.PomodoroPaused:
        pw.startBtn.Enable()
        pw.pauseBtn.Enable()
        pw.resetBtn.Enable()
        pw.startBtn.SetText("Resume")
        pw.pauseBtn.SetText("Resume")
	}
}

// ProgressRing renders a circular segmented progress indicator.
type ProgressRing struct {
    widget.BaseWidget
    Progress    float32      // 0..1 of remaining progress
    Segments    int          // number of segments around the circle
    StartAngle  float64      // radians; 0 is right, pi is left
    FillColor   color.Color  // color for filled segments
    BgColor     color.Color  // color for background segments
    InnerRatio  float32      // inner radius ratio relative to half of min(size)
    SegLength   float32      // length of each radial segment in px
    StrokeWidth float32      // thickness of each segment
}

func NewProgressRing(fill, bg color.Color) *ProgressRing {
    pr := &ProgressRing{
        Progress:    1,
        Segments:    90,
        StartAngle:  math.Pi, // start from left side
        FillColor:   fill,
        BgColor:     bg,
        InnerRatio:  0.68,
        SegLength:   14,
        StrokeWidth: 4,
    }
    pr.ExtendBaseWidget(pr)
    return pr
}

func (pr *ProgressRing) SetProgress(p float32) {
    if p < 0 {
        p = 0
    }
    if p > 1 {
        p = 1
    }
    if pr.Progress == p {
        return
    }
    pr.Progress = p
    pr.Refresh()
}

func (pr *ProgressRing) MinSize() fyne.Size {
    return fyne.NewSize(200, 200)
}

func (pr *ProgressRing) CreateRenderer() fyne.WidgetRenderer {
    lines := make([]*canvas.Line, pr.Segments)
    objs := make([]fyne.CanvasObject, pr.Segments)
    for i := 0; i < pr.Segments; i++ {
        ln := canvas.NewLine(pr.BgColor)
        ln.StrokeWidth = pr.StrokeWidth
        lines[i] = ln
        objs[i] = ln
    }
    return &progressRingRenderer{
        ring:  pr,
        lines: lines,
        objs:  objs,
    }
}

type progressRingRenderer struct {
    ring  *ProgressRing
    lines []*canvas.Line
    objs  []fyne.CanvasObject
}

func (r *progressRingRenderer) Layout(size fyne.Size) {
    cx := size.Width / 2
    cy := size.Height / 2
    halfMin := float32(math.Min(float64(size.Width), float64(size.Height)) / 2)
    inner := halfMin * r.ring.InnerRatio
    outer := inner + r.ring.SegLength

    for i := 0; i < r.ring.Segments; i++ {
        // clockwise angle from start
        ang := r.ring.StartAngle + 2*math.Pi*float64(i)/float64(r.ring.Segments)
        cos := float32(math.Cos(ang))
        sin := float32(math.Sin(ang))
        x1 := cx + inner*cos
        y1 := cy + inner*sin
        x2 := cx + outer*cos
        y2 := cy + outer*sin
        ln := r.lines[i]
        ln.Position1 = fyne.NewPos(x1, y1)
        ln.Position2 = fyne.NewPos(x2, y2)
    }
}

func (r *progressRingRenderer) MinSize() fyne.Size {
    return r.ring.MinSize()
}

func (r *progressRingRenderer) Refresh() {
    // recolor based on current progress
    filled := int(float64(r.ring.Segments) * float64(r.ring.Progress) + 0.5)
    if filled < 0 {
        filled = 0
    }
    if filled > r.ring.Segments {
        filled = r.ring.Segments
    }
    for i := 0; i < r.ring.Segments; i++ {
        if i < filled {
            r.lines[i].StrokeColor = r.ring.FillColor
        } else {
            r.lines[i].StrokeColor = r.ring.BgColor
        }
        r.lines[i].Refresh()
    }
    // ensure geometry is correct if resized
    r.Layout(r.ring.Size())
}

func (r *progressRingRenderer) BackgroundColor() fyne.ThemeColorName { return "" }
func (r *progressRingRenderer) Objects() []fyne.CanvasObject         { return r.objs }
func (r *progressRingRenderer) Destroy()                             {}

// Event handlers
func (pw *PomodoroWindow) onStartClicked() {
	if pw.timer.State == models.PomodoroPaused {
		pw.timer.Resume()
	} else {
		pw.timer.Start()
	}
	pw.tick()
}

func (pw *PomodoroWindow) onPauseClicked() {
	if pw.timer.State == models.PomodoroPaused {
		pw.timer.Resume()
	} else {
		pw.timer.Pause()
	}
	pw.tick()
}

func (pw *PomodoroWindow) onResetClicked() {
	pw.timer.Reset()
	pw.tick()
}

// (config change handled by spinners' callbacks)

// Show displays the window
func (pw *PomodoroWindow) Show() {
	pw.window.Show()
}

// SetOnClosed sets the callback for when the window is closed
func (pw *PomodoroWindow) SetOnClosed(callback func()) {
	pw.window.SetOnClosed(func() {
		pw.stopTicker()
		if callback != nil {
			callback()
		}
	})
}

// UpdateTheme updates the theme of the pomodoro window
func (pw *PomodoroWindow) UpdateTheme(isGruvbox bool) {
	pw.isGruvbox = isGruvbox

	// Update the app theme first (should already be done by MainWindow)
	// Then recreate the UI with new colors
	pw.setupUI()

	// Preserve timer state and refresh display
	pw.tick()
}
