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
	window    fyne.Window
	timer     *models.PomodoroTimer
	config    *models.PomodoroConfig
	isGruvbox bool

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
	isInitializing bool                 // prevent animation on startup
	timerContainer *fyne.Container      // for hiding/showing digits during animation
	lastState      models.PomodoroState // track previous state to detect transitions
}

// NewPomodoroWindow creates a new Pomodoro timer window
func NewPomodoroWindow(app fyne.App, isGruvbox bool) *PomodoroWindow {
	config := models.NewDefaultPomodoroConfig()
	timer := models.NewPomodoroTimer(config)

	pw := &PomodoroWindow{
		window:         app.NewWindow("Pomodoro Timer"),
		timer:          timer,
		config:         config,
		isGruvbox:      isGruvbox,
		lastUpdate:     time.Now().Add(-time.Second),
		isInitializing: true,                // prevent animation on first tick
		lastState:      models.PomodoroIdle, // track previous state
	}

	pw.setupUI()
	pw.startTicker()
	pw.tick() // initial update

	// Reset initializing flag after 2 seconds to prevent spurious animations
	time.AfterFunc(2*time.Second, func() {
		pw.isInitializing = false
	})

	return pw
}

// setupUI initializes the user interface
func (pw *PomodoroWindow) setupUI() {
	// Window properties
	pw.window.Resize(fyne.NewSize(200, 350))
	pw.window.SetFixedSize(true)
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
	cfgHeader.TextSize = 30

	labelWork := canvas.NewText("Work time (min):", titleColor)
	labelWork.TextSize = 16
	labelWork.TextStyle = fyne.TextStyle{Bold: true}
	labelShort := canvas.NewText("Short break (min):", titleColor)
	labelShort.TextSize = 16
	labelShort.TextStyle = fyne.TextStyle{Bold: true}
	labelLong := canvas.NewText("Long break (min):", titleColor)
	labelLong.TextSize = 16
	labelLong.TextStyle = fyne.TextStyle{Bold: true}

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
	spinnerVerticalOffset := pw.workSpinner.MinSize().Height * 0.2
	wrapSpinner := func(spinner *NumberSpinner) fyne.CanvasObject {
		return container.NewVBox(
			CreateSpacer(1, spinnerVerticalOffset),
			spinner,
		)
	}
	configForm := container.NewVBox(
		container.NewCenter(cfgHeader),
		CreateSpacer(1, 7),
		container.NewGridWithColumns(2,
			labelWork,
			wrapSpinner(pw.workSpinner),
		),
		container.NewGridWithColumns(2,
			labelShort,
			wrapSpinner(pw.shortBreakSpinner),
		),
		container.NewGridWithColumns(2,
			labelLong,
			wrapSpinner(pw.longBreakSpinner),
		),
	)

	// Progress ring colors
	var tickBg color.Color
	if _, ok := currentTheme.(*LightSoftTheme); ok {
		tickBg = color.White
	} else {
		tickBg = hex("#504945")
	}
	pw.progressRing = NewProgressRing(tickBg)

	// Layout
	// Wrap timer canvas so we can hide it during animation
	pw.timerContainer = container.NewCenter(pw.timerCanvas)
	ringWithDigits := container.NewMax(
		container.NewCenter(pw.progressRing),
		pw.timerContainer,
	)
	ringContainer := container.NewCenter(container.NewGridWrap(fyne.NewSize(250, 270), ringWithDigits))
	labelsBlock := container.NewVBox(
		container.NewCenter(pw.stateCanvas),
		CreateSpacer(1, 33),
		container.NewCenter(pw.sessionsCanvas),
	)
	labelHeight := labelsBlock.MinSize().Height
	liftOffset := float32(55)
	ringHeight := ringContainer.MinSize().Height
	if ringHeight < liftOffset {
		liftOffset = ringHeight / 2 // avoid negative spacer
	}
	spacerBeforeLabels := ringHeight - liftOffset
	// Overlay labels so they sit 30px closer to the ring without altering container height.
	stackBase := container.NewVBox(
		ringContainer,
		CreateSpacer(1, labelHeight),
	)
	labelOverlay := container.NewVBox(
		CreateSpacer(1, spacerBeforeLabels),
		labelsBlock,
	)
	ringAndLabels := container.NewStack(stackBase, labelOverlay)

	timerDisplay := container.NewVBox(
		ringAndLabels,
		CreateSpacer(1, 0), // ещё на ~20px ближе кнопки к индикатору
		container.NewCenter(buttonRow),
		CreateSpacer(1, 40),
	)

	content := container.NewVBox(
		timerDisplay,
		CreateFixedSeparator(),
		CreateSpacer(1, 20),
		configForm,
		CreateSpacer(1, 20),
	)

	// Wrap in padding (24px left/right, 10px top)
	paddedContent := container.NewBorder(
		CreateSpacer(1, 50), nil,
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

	// Update progress ring (elapsed fraction, fills clockwise from left)
	var prog float32 = 0
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
		// Calculate elapsed fraction instead of remaining
		elapsed := total - pw.timer.TimeRemaining
		prog = float32(elapsed.Seconds() / total.Seconds())
		if prog < 0 {
			prog = 0
		}
		if prog > 1 {
			prog = 1
		}
	}
	if pw.progressRing != nil {
		pw.progressRing.SetProgress(prog)

		// Detect transitions that mean a period just finished:
		// - Work -> Short/Long break
		// - Short/Long break -> Idle
		isWorkTransition := pw.lastState == models.PomodoroWork &&
			(pw.timer.State == models.PomodoroShortBreak || pw.timer.State == models.PomodoroLongBreak)
		isBreakTransition := (pw.lastState == models.PomodoroShortBreak || pw.lastState == models.PomodoroLongBreak) &&
			pw.timer.State == models.PomodoroIdle
		periodCompleted := (isWorkTransition || isBreakTransition)

		// Update lastState for next tick
		pw.lastState = pw.timer.State

		// Trigger animation when period completes
		if periodCompleted && !pw.isInitializing {
			// Play completion animation (on main thread)
			// Hide timer digits during animation
			pw.timerContainer.Hide()
			pw.progressRing.PlayCompletionAnimation()
			// Show digits again after animation (now 2.5 seconds for longer animation)
			time.AfterFunc(2500*time.Millisecond, func() {
				pw.timerContainer.Show()
				pw.timerContainer.Refresh()
			})
		}
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
		// When paused: only Pause button becomes Resume, Start is disabled
		pw.startBtn.Disable()
		pw.pauseBtn.Enable()
		pw.resetBtn.Enable()
		pw.pauseBtn.SetText("Resume")
	}
}

// ProgressRing renders a circular segmented progress indicator.
type ProgressRing struct {
	widget.BaseWidget
	Progress       float32     // 0..1 of elapsed progress
	Segments       int         // number of segments around the circle
	StartAngle     float64     // radians; 0 is right, pi is left
	StartColor     color.Color // color at start (green)
	EndColor       color.Color // color at end (red)
	BgColor        color.Color // color for background segments
	InnerRatio     float32     // inner radius ratio relative to half of min(size)
	SegLength      float32     // length of each radial segment in px
	StrokeWidth    float32     // thickness of each segment
	IsCompleting   bool        // true when showing completion animation
	CompletionAnim float32     // 0..1 animation progress for completion
}

func NewProgressRing(bg color.Color) *ProgressRing {
	pr := &ProgressRing{
		Progress:    0,
		Segments:    60,
		StartAngle:  math.Pi,        // start from left side
		StartColor:  hex("#d65c5c"), // Red-orange - start (0% progress, warning state)
		EndColor:    hex("#a4d868"), // Bright vibrant green - end (100% progress, complete)
		BgColor:     bg,
		InnerRatio:  1.2,
		SegLength:   31,
		StrokeWidth: 6,
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

// PlayCompletionAnimation starts the green flash + checkmark animation
func (pr *ProgressRing) PlayCompletionAnimation() {
	pr.IsCompleting = true
	pr.CompletionAnim = 0
	pr.Refresh()

	// Animate over 2500ms using sequential timeouts
	const frameDuration = 16 * time.Millisecond
	const animDuration = 2500 * time.Millisecond
	const framesCount = int(animDuration / frameDuration)

	for frame := 0; frame <= framesCount; frame++ {
		frame := frame // capture for closure
		time.AfterFunc(time.Duration(frame)*frameDuration, func() {
			if frame >= framesCount {
				// Animation complete
				pr.IsCompleting = false
				pr.CompletionAnim = 0
			} else {
				// Calculate progress 0..1
				pr.CompletionAnim = float32(frame) / float32(framesCount)
			}
			pr.Refresh()
		})
	}
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
	// Add background circle for completion animation
	completionBg := canvas.NewCircle(color.NRGBA{R: 0, G: 0, B: 0, A: 0})
	// Add checkmark lines (enlarged, 3x scale)
	checkLine1 := canvas.NewLine(color.NRGBA{R: 164, G: 216, B: 104, A: 255})
	checkLine1.StrokeWidth = 10 // increased from 4
	checkLine2 := canvas.NewLine(color.NRGBA{R: 164, G: 216, B: 104, A: 255})
	checkLine2.StrokeWidth = 10 // increased from 4

	objs = append(objs, completionBg, checkLine1, checkLine2)

	return &progressRingRenderer{
		ring:         pr,
		lines:        lines,
		objs:         objs,
		completionBg: completionBg,
		checkLine1:   checkLine1,
		checkLine2:   checkLine2,
	}
}

type progressRingRenderer struct {
	ring         *ProgressRing
	lines        []*canvas.Line
	objs         []fyne.CanvasObject
	completionBg *canvas.Circle
	checkLine1   *canvas.Line
	checkLine2   *canvas.Line
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
	size := r.ring.Size()
	cx := size.Width / 2
	cy := size.Height / 2

	if r.ring.IsCompleting {
		// Completion animation phases:
		// 0.0-0.33: Green flash - all segments turn bright green
		// 0.33-0.66: Show checkmark
		// 0.66-1.0: Fade out
		anim := r.ring.CompletionAnim

		if anim < 0.33 {
			// Flash phase - make all segments bright green
			flashAlpha := uint8(255 * (1 - (anim / 0.33)))
			for i := 0; i < r.ring.Segments; i++ {
				r.lines[i].StrokeColor = color.NRGBA{R: 164, G: 216, B: 104, A: flashAlpha}
				r.lines[i].Refresh()
			}
			// Green background circle
			r.completionBg.FillColor = color.NRGBA{R: 164, G: 216, B: 104, A: uint8(100 * (1 - (anim / 0.33)))}
		} else if anim < 0.66 {
			// Checkmark phase - hide segments, show checkmark
			for i := 0; i < r.ring.Segments; i++ {
				r.lines[i].StrokeColor = r.ring.BgColor
				r.lines[i].Refresh()
			}
			r.completionBg.FillColor = color.NRGBA{R: 164, G: 216, B: 104, A: 200}

			// Draw checkmark lines (enlarged by 3x)
			halfMin := float32(math.Min(float64(size.Width), float64(size.Height)) / 2)
			checkSize := halfMin * 0.4 * 3 // 3x larger

			// Checkmark starts at 30% down from center, goes to 60% down, then up-right
			// Left diagonal: from (cx-checkSize*0.3, cy+checkSize*0.2) to (cx, cy+checkSize*0.5)
			r.checkLine1.Position1 = fyne.NewPos(cx-checkSize*0.3, cy+checkSize*0.2)
			r.checkLine1.Position2 = fyne.NewPos(cx, cy+checkSize*0.5)

			// Right diagonal: from (cx, cy+checkSize*0.5) to (cx+checkSize*0.5, cy-checkSize*0.3)
			r.checkLine2.Position1 = fyne.NewPos(cx, cy+checkSize*0.5)
			r.checkLine2.Position2 = fyne.NewPos(cx+checkSize*0.5, cy-checkSize*0.3)

			r.checkLine1.Refresh()
			r.checkLine2.Refresh()
		} else {
			// Fade out phase
			fadeAlpha := uint8(200 * (1 - ((anim - 0.66) / 0.34)))
			r.completionBg.FillColor = color.NRGBA{R: 164, G: 216, B: 104, A: fadeAlpha}
			r.completionBg.Refresh()

			// Fade out checkmark
			checkAlpha := uint8(255 * (1 - ((anim - 0.66) / 0.34)))
			r.checkLine1.StrokeColor = color.NRGBA{R: 164, G: 216, B: 104, A: checkAlpha}
			r.checkLine2.StrokeColor = color.NRGBA{R: 164, G: 216, B: 104, A: checkAlpha}
			r.checkLine1.Refresh()
			r.checkLine2.Refresh()
		}
	} else {
		// Normal progress display
		// Only show segments if progress is greater than 0
		var filled int
		if r.ring.Progress > 0 {
			filled = int(float64(r.ring.Segments)*float64(r.ring.Progress) + 0.5)
			if filled < 0 {
				filled = 0
			}
			if filled > r.ring.Segments {
				filled = r.ring.Segments
			}
		}

		// Extract RGB components from start and end colors
		sr, sg, sb, _ := r.ring.StartColor.RGBA()
		er, eg, eb, _ := r.ring.EndColor.RGBA()

		for i := 0; i < r.ring.Segments; i++ {
			if i < filled {
				// Calculate gradient color for this segment
				t := float32(i) / float32(r.ring.Segments)
				nr := uint8((float32(sr>>8)*(1-t) + float32(er>>8)*t))
				ng := uint8((float32(sg>>8)*(1-t) + float32(eg>>8)*t))
				nb := uint8((float32(sb>>8)*(1-t) + float32(eb>>8)*t))
				r.lines[i].StrokeColor = color.RGBA{R: nr, G: ng, B: nb, A: 255}
			} else {
				r.lines[i].StrokeColor = r.ring.BgColor
			}
			r.lines[i].Refresh()
		}

		// Hide completion elements
		r.completionBg.FillColor = color.NRGBA{R: 0, G: 0, B: 0, A: 0}
		r.completionBg.Refresh()
		r.checkLine1.StrokeColor = color.NRGBA{R: 0, G: 0, B: 0, A: 0}
		r.checkLine2.StrokeColor = color.NRGBA{R: 0, G: 0, B: 0, A: 0}
		r.checkLine1.Refresh()
		r.checkLine2.Refresh()
	}

	// ensure geometry is correct if resized
	r.Layout(size)
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
