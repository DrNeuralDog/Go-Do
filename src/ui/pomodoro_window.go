package ui

import (
	"fmt"
	"image/color"
	"math"
	"time"

	"godo/src/models"
	"godo/src/ui/helpers"
	"godo/src/ui/widgets"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const (
	completionFrameDuration = 16 * time.Millisecond
	completionAnimDuration  = 2500 * time.Millisecond
	completionFlashEnd      = 0.33
	completionCheckEnd      = 0.66
)

// PomodoroWindow manages Pomodoro timer UI
type PomodoroWindow struct {
	window    fyne.Window
	timer     *models.PomodoroTimer
	config    *models.PomodoroConfig
	isGruvbox bool // legacy theme flag kept for old callers

	timerCanvas    *canvas.Text
	stateCanvas    *canvas.Text
	startBtn       *widgets.SimpleRectButton
	pauseBtn       *widgets.SimpleRectButton
	resetBtn       *widgets.SimpleRectButton
	sessionsCanvas *canvas.Text
	progressRing   *ProgressRing

	workSpinner       *widgets.NumberSpinner
	shortBreakSpinner *widgets.NumberSpinner
	longBreakSpinner  *widgets.NumberSpinner

	anim           *fyne.Animation
	lastUpdate     time.Time
	isInitializing bool                 // blocks animation during first refresh
	timerContainer *fyne.Container      // hides digits while completion animation runs
	lastState      models.PomodoroState // previous state for transition detection
}

// NewPomodoroWindow builds Pomodoro timer window
func NewPomodoroWindow(app fyne.App, isGruvbox bool) *PomodoroWindow {
	config := models.NewDefaultPomodoroConfig()
	timer := models.NewPomodoroTimer(config)

	pw := &PomodoroWindow{
		window:         app.NewWindow("Pomodoro Timer"),
		timer:          timer,
		config:         config,
		isGruvbox:      isGruvbox,
		lastUpdate:     time.Now().Add(-time.Second),
		isInitializing: true,
		lastState:      models.PomodoroIdle,
	}

	pw.setupUI()
	pw.startTicker()
	pw.tick()

	time.AfterFunc(2*time.Second, func() {
		runOnMainThread(func() {
			pw.isInitializing = false
		})
	})

	return pw
}

// setupUI rebuilds Pomodoro window content
func (pw *PomodoroWindow) setupUI() {
	var bgStart, bgEnd color.Color
	var titleColor color.Color

	currentTheme := fyne.CurrentApp().Settings().Theme()
	isLightTheme := helpers.IsLightTheme()

	if gradientTheme, ok := currentTheme.(interface {
		GetHeaderGradientColors() (color.Color, color.Color)
	}); ok {
		bgStart, bgEnd = gradientTheme.GetHeaderGradientColors()
	} else {
		bgStart = helpers.GetBackgroundColor()
		bgEnd = bgStart
	}

	titleColor = color.White

	if !isLightTheme {
		titleColor = helpers.Hex("#fabd2f")
	}

	pw.timerCanvas = canvas.NewText("25:00", titleColor)
	pw.timerCanvas.TextStyle = fyne.TextStyle{Bold: true}
	pw.timerCanvas.TextSize = 77

	pw.stateCanvas = canvas.NewText("Ready", titleColor)
	pw.stateCanvas.TextSize = 18

	pw.sessionsCanvas = canvas.NewText("Sessions: 0", titleColor)
	pw.sessionsCanvas.TextSize = 16

	var btnBg, btnFg color.Color
	if isLightTheme {
		btnBg = helpers.Hex("#ff8c42")
		btnFg = color.White
	} else {
		btnBg = helpers.Hex("#504945")
		btnFg = helpers.Hex("#fabd2f")
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

	darkText := helpers.Hex("#3c3836")
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
	wrapSpinner := func(spinner *widgets.NumberSpinner) fyne.CanvasObject {
		return container.NewVBox(
			helpers.CreateSpacer(1, spinnerVerticalOffset),
			spinner,
		)
	}

	configForm := container.NewVBox(
		container.NewCenter(cfgHeader),
		helpers.CreateSpacer(1, 7),
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

	var tickBg color.Color

	if isLightTheme {
		tickBg = color.White
	} else {
		tickBg = helpers.Hex("#504945")
	}

	pw.progressRing = NewProgressRing(tickBg)

	pw.timerContainer = container.NewCenter(pw.timerCanvas)
	ringWithDigits := container.NewMax(
		container.NewCenter(pw.progressRing),
		pw.timerContainer,
	)

	ringContainer := container.NewCenter(container.NewGridWrap(fyne.NewSize(250, 270), ringWithDigits))
	labelsBlock := container.NewVBox(
		container.NewCenter(pw.stateCanvas),
		helpers.CreateSpacer(1, 33),
		container.NewCenter(pw.sessionsCanvas),
	)

	labelHeight := labelsBlock.MinSize().Height
	liftOffset := float32(55)
	ringHeight := ringContainer.MinSize().Height

	if ringHeight < liftOffset {
		liftOffset = ringHeight / 2
	}

	spacerBeforeLabels := ringHeight - liftOffset

	// Метки подтянуты к кольцу без изменения высоты контейнера
	stackBase := container.NewVBox(
		ringContainer,
		helpers.CreateSpacer(1, labelHeight),
	)

	labelOverlay := container.NewVBox(
		helpers.CreateSpacer(1, spacerBeforeLabels),
		labelsBlock,
	)

	ringAndLabels := container.NewStack(stackBase, labelOverlay)

	timerDisplay := container.NewVBox(
		ringAndLabels,
		helpers.CreateSpacer(1, 0),
		container.NewCenter(buttonRow),
		helpers.CreateSpacer(1, 40),
	)

	content := container.NewVBox(
		timerDisplay,
		helpers.CreateFixedSeparator(),
		helpers.CreateSpacer(1, 20),
		configForm,
		helpers.CreateSpacer(1, 20),
	)

	paddedContent := container.NewBorder(
		helpers.CreateSpacer(1, 50), nil,
		helpers.CreateSpacer(ButtonPadding, 1),
		helpers.CreateSpacer(ButtonPadding, 1),
		content,
	)

	background := NewGradientRect(bgStart, bgEnd, 0)
	finalContent := container.NewMax(background, paddedContent)

	pw.window.SetContent(finalContent)
	pw.window.SetFixedSize(true)
	pw.window.Resize(fyne.NewSize(PomodoroWindowWidth, PomodoroWindowHeight))
	pw.window.CenterOnScreen()
}

// startTicker keeps timer display synced
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

// stopTicker stops timer sync
func (pw *PomodoroWindow) stopTicker() {
	if pw.anim != nil {
		pw.anim.Stop()

		pw.anim = nil
	}
}

// tick syncs timer state with visible controls
func (pw *PomodoroWindow) tick() {
	pw.timer.Update()

	pw.refreshTimerText()

	if pw.progressRing != nil {
		pw.progressRing.SetProgress(pw.currentProgress())
	}

	if pw.consumePeriodCompletion() {
		pw.playPeriodCompletion()
	}

	pw.syncButtonsState()
}

// refreshTimerText updates timer labels
func (pw *PomodoroWindow) refreshTimerText() {
	pw.timerCanvas.Text = pw.timerText()
	pw.timerCanvas.Refresh()

	pw.stateCanvas.Text = pw.timer.GetStateString()
	pw.stateCanvas.Refresh()

	pw.sessionsCanvas.Text = fmt.Sprintf("Sessions: %d", pw.timer.SessionsCompleted)
	pw.sessionsCanvas.Refresh()
}

// timerText returns display text for current timer state
func (pw *PomodoroWindow) timerText() string {
	if pw.timer.State == models.PomodoroIdle {
		return fmt.Sprintf("%02d:00", pw.config.WorkDuration)
	}

	minutes := int(pw.timer.TimeRemaining.Minutes())
	seconds := int(pw.timer.TimeRemaining.Seconds()) % 60

	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

// currentProgress returns elapsed timer progress from 0 to 1
func (pw *PomodoroWindow) currentProgress() float32 {
	total := pw.currentTotalDuration()

	if pw.timer.State == models.PomodoroIdle || total <= 0 {
		return 0
	}

	elapsed := total - pw.timer.TimeRemaining
	progress := float32(elapsed.Seconds() / total.Seconds())

	return clampProgress(progress)
}

// currentTotalDuration returns visible duration for progress ring
func (pw *PomodoroWindow) currentTotalDuration() time.Duration {
	total := pw.timer.GetCurrentDuration()

	if pw.timer.TimeRemaining <= 0 {
		return total
	}

	if pw.timer.State != models.PomodoroPaused && total != 0 {
		return total
	}

	work := time.Duration(pw.config.WorkDuration) * time.Minute
	shortBreak := time.Duration(pw.config.ShortBreakDuration) * time.Minute
	longBreak := time.Duration(pw.config.LongBreakDuration) * time.Minute

	// На паузе модель не хранит прошлую фазу, поэтому берем ближайшую длительность
	candidates := []time.Duration{work, shortBreak, longBreak}

	return nearestDuration(pw.timer.TimeRemaining, candidates, work)
}

// nearestDuration finds smallest duration that can contain remaining time
func nearestDuration(remaining time.Duration, candidates []time.Duration, fallback time.Duration) time.Duration {
	var best time.Duration

	for _, duration := range candidates {
		if duration >= remaining && (best == 0 || duration < best) {
			best = duration
		}
	}

	if best == 0 {
		return fallback
	}

	return best
}

// clampProgress keeps progress inside renderer range
func clampProgress(progress float32) float32 {
	if progress < 0 {
		return 0
	}

	if progress > 1 {
		return 1
	}

	return progress
}

// consumePeriodCompletion detects completed Pomodoro phase
func (pw *PomodoroWindow) consumePeriodCompletion() bool {
	isWorkTransition := pw.lastState == models.PomodoroWork &&
		(pw.timer.State == models.PomodoroShortBreak || pw.timer.State == models.PomodoroLongBreak)

	isBreakTransition := (pw.lastState == models.PomodoroShortBreak || pw.lastState == models.PomodoroLongBreak) &&
		pw.timer.State == models.PomodoroIdle

	pw.lastState = pw.timer.State

	return (isWorkTransition || isBreakTransition) && !pw.isInitializing
}

// playPeriodCompletion runs completion animation
func (pw *PomodoroWindow) playPeriodCompletion() {
	runOnMainThread(func() {
		if pw.timerContainer != nil {
			pw.timerContainer.Hide()
		}

		if pw.progressRing != nil {
			pw.progressRing.PlayCompletionAnimation()
		}
	})

	time.AfterFunc(completionAnimDuration, func() {
		runOnMainThread(func() {
			if pw.timerContainer == nil {
				return
			}

			pw.timerContainer.Show()
			pw.timerContainer.Refresh()
		})
	})
}

// syncButtonsState updates control buttons
func (pw *PomodoroWindow) syncButtonsState() {
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
		pw.startBtn.Disable()
		pw.pauseBtn.Enable()
		pw.resetBtn.Enable()
		pw.pauseBtn.SetText("Resume")
	}
}

// ProgressRing renders segmented timer progress
type ProgressRing struct {
	widget.BaseWidget
	Progress       float32     // elapsed progress from 0 to 1
	Segments       int         // segment count around the circle
	StartAngle     float64     // radians, 0 is right and pi is left
	StartColor     color.Color // low progress segment color
	EndColor       color.Color // high progress segment color
	BgColor        color.Color // idle segment color
	InnerRatio     float32     // inner radius ratio
	SegLength      float32     // radial segment length
	StrokeWidth    float32     // segment thickness
	IsCompleting   bool        // completion animation flag
	CompletionAnim float32     // completion animation progress from 0 to 1
}

// NewProgressRing builds segmented progress ring
func NewProgressRing(bg color.Color) *ProgressRing {
	pr := &ProgressRing{
		Progress:    0,
		Segments:    60,
		StartAngle:  math.Pi,
		StartColor:  helpers.Hex("#d65c5c"),
		EndColor:    helpers.Hex("#a4d868"),
		BgColor:     bg,
		InnerRatio:  1.2,
		SegLength:   31,
		StrokeWidth: 6,
	}

	pr.ExtendBaseWidget(pr)

	return pr
}

// SetProgress updates ring progress
func (pr *ProgressRing) SetProgress(p float32) {
	p = clampProgress(p)

	if pr.Progress == p {
		return
	}

	pr.Progress = p

	pr.Refresh()
}

// PlayCompletionAnimation starts green flash and checkmark animation
func (pr *ProgressRing) PlayCompletionAnimation() {
	runOnMainThread(func() {
		pr.IsCompleting = true
		pr.CompletionAnim = 0

		pr.Refresh()
	})

	const framesCount = int(completionAnimDuration / completionFrameDuration)

	for frame := 0; frame <= framesCount; frame++ {
		frame := frame

		time.AfterFunc(time.Duration(frame)*completionFrameDuration, func() {
			runOnMainThread(func() {
				if frame >= framesCount {
					pr.IsCompleting = false
					pr.CompletionAnim = 0
				} else {
					pr.CompletionAnim = float32(frame) / float32(framesCount)
				}

				pr.Refresh()
			})
		})
	}
}

// MinSize returns fixed ring size
func (pr *ProgressRing) MinSize() fyne.Size {
	return fyne.NewSize(200, 200)
}

// CreateRenderer builds progress ring renderer
func (pr *ProgressRing) CreateRenderer() fyne.WidgetRenderer {
	lines := make([]*canvas.Line, pr.Segments)
	objs := make([]fyne.CanvasObject, pr.Segments)

	for i := 0; i < pr.Segments; i++ {
		ln := canvas.NewLine(pr.BgColor)
		ln.StrokeWidth = pr.StrokeWidth
		lines[i] = ln
		objs[i] = ln
	}

	completionBg := canvas.NewCircle(color.NRGBA{R: 0, G: 0, B: 0, A: 0})

	checkLine1 := canvas.NewLine(completionGreen(0))
	checkLine1.StrokeWidth = 10
	checkLine2 := canvas.NewLine(completionGreen(0))
	checkLine2.StrokeWidth = 10

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

// Layout positions ring segments and completion overlay
func (r *progressRingRenderer) Layout(size fyne.Size) {
	cx := size.Width / 2
	cy := size.Height / 2
	halfMin := float32(math.Min(float64(size.Width), float64(size.Height)) / 2)
	inner := halfMin * r.ring.InnerRatio
	outer := inner + r.ring.SegLength

	r.completionBg.Move(fyne.NewPos(0, 0))
	r.completionBg.Resize(size)

	for i := 0; i < r.ring.Segments; i++ {
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

// MinSize returns renderer minimum size
func (r *progressRingRenderer) MinSize() fyne.Size {
	return r.ring.MinSize()
}

// Refresh applies current progress or completion animation state
func (r *progressRingRenderer) Refresh() {
	size := r.ring.Size()
	r.Layout(size)

	if r.ring.IsCompleting {
		r.refreshCompletion(size)
		return
	}

	r.refreshProgress()
}

// refreshCompletion renders current completion animation phase
func (r *progressRingRenderer) refreshCompletion(size fyne.Size) {
	anim := r.ring.CompletionAnim

	switch {
	case anim < completionFlashEnd:
		r.refreshFlashPhase(anim)
	case anim < completionCheckEnd:
		r.refreshCheckmarkPhase(size)
	default:
		r.refreshFadePhase(anim)
	}
}

// refreshFlashPhase renders green flash
func (r *progressRingRenderer) refreshFlashPhase(anim float32) {
	phaseLeft := 1 - anim/completionFlashEnd

	r.setSegmentsColor(completionGreen(uint8(255 * phaseLeft)))
	r.completionBg.FillColor = completionGreen(uint8(100 * phaseLeft))
	r.completionBg.Refresh()
	r.hideCheckmark()
}

// refreshCheckmarkPhase renders large completion checkmark
func (r *progressRingRenderer) refreshCheckmarkPhase(size fyne.Size) {
	r.setSegmentsColor(r.ring.BgColor)

	r.completionBg.FillColor = completionGreen(200)
	r.completionBg.Refresh()

	r.layoutCheckmark(size)
	r.checkLine1.StrokeColor = completionGreen(255)
	r.checkLine2.StrokeColor = completionGreen(255)
	r.checkLine1.Refresh()
	r.checkLine2.Refresh()
}

// refreshFadePhase fades completion checkmark out
func (r *progressRingRenderer) refreshFadePhase(anim float32) {
	phaseLeft := 1 - ((anim - completionCheckEnd) / (1 - completionCheckEnd))
	phaseLeft = clampProgress(phaseLeft)

	r.setSegmentsColor(r.ring.BgColor)

	r.completionBg.FillColor = completionGreen(uint8(200 * phaseLeft))
	r.completionBg.Refresh()

	checkAlpha := uint8(255 * phaseLeft)
	r.checkLine1.StrokeColor = completionGreen(checkAlpha)
	r.checkLine2.StrokeColor = completionGreen(checkAlpha)
	r.checkLine1.Refresh()
	r.checkLine2.Refresh()
}

// refreshProgress renders normal segmented progress
func (r *progressRingRenderer) refreshProgress() {
	filled := r.filledSegments()

	sr, sg, sb, _ := r.ring.StartColor.RGBA()
	er, eg, eb, _ := r.ring.EndColor.RGBA()

	for i := 0; i < r.ring.Segments; i++ {
		if i < filled {
			t := float32(i) / float32(r.ring.Segments)
			r.lines[i].StrokeColor = blendedSegmentColor(t, sr, sg, sb, er, eg, eb)
		} else {
			r.lines[i].StrokeColor = r.ring.BgColor
		}

		r.lines[i].Refresh()
	}

	r.hideCompletionObjects()
}

// filledSegments returns number of active progress segments
func (r *progressRingRenderer) filledSegments() int {
	if r.ring.Progress <= 0 {
		return 0
	}

	filled := int(float64(r.ring.Segments)*float64(r.ring.Progress) + 0.5)

	if filled < 0 {
		return 0
	}

	if filled > r.ring.Segments {
		return r.ring.Segments
	}

	return filled
}

// setSegmentsColor applies one color to all ring segments
func (r *progressRingRenderer) setSegmentsColor(clr color.Color) {
	for _, line := range r.lines {
		line.StrokeColor = clr
		line.Refresh()
	}
}

// layoutCheckmark positions completion checkmark
func (r *progressRingRenderer) layoutCheckmark(size fyne.Size) {
	cx := size.Width / 2
	cy := size.Height / 2
	halfMin := float32(math.Min(float64(size.Width), float64(size.Height)) / 2)

	// Галочка намеренно крупная, чтобы перекрывать кольцо
	checkSize := halfMin * 1.2

	r.checkLine1.Position1 = fyne.NewPos(cx-checkSize*0.3, cy+checkSize*0.2)
	r.checkLine1.Position2 = fyne.NewPos(cx, cy+checkSize*0.5)
	r.checkLine2.Position1 = fyne.NewPos(cx, cy+checkSize*0.5)
	r.checkLine2.Position2 = fyne.NewPos(cx+checkSize*0.5, cy-checkSize*0.3)
}

// hideCompletionObjects clears completion overlay
func (r *progressRingRenderer) hideCompletionObjects() {
	r.completionBg.FillColor = completionGreen(0)
	r.completionBg.Refresh()
	r.hideCheckmark()
}

// hideCheckmark clears checkmark lines
func (r *progressRingRenderer) hideCheckmark() {
	r.checkLine1.StrokeColor = completionGreen(0)
	r.checkLine2.StrokeColor = completionGreen(0)
	r.checkLine1.Refresh()
	r.checkLine2.Refresh()
}

// blendedSegmentColor returns gradient color for segment position
func blendedSegmentColor(t float32, sr, sg, sb, er, eg, eb uint32) color.Color {
	nr := uint8(float32(sr>>8)*(1-t) + float32(er>>8)*t)
	ng := uint8(float32(sg>>8)*(1-t) + float32(eg>>8)*t)
	nb := uint8(float32(sb>>8)*(1-t) + float32(eb>>8)*t)

	return color.RGBA{R: nr, G: ng, B: nb, A: 255}
}

// completionGreen returns completion animation color
func completionGreen(alpha uint8) color.NRGBA {
	return color.NRGBA{R: 164, G: 216, B: 104, A: alpha}
}

// BackgroundColor keeps renderer transparent
func (r *progressRingRenderer) BackgroundColor() fyne.ThemeColorName { return "" }

// Objects returns renderer canvas objects
func (r *progressRingRenderer) Objects() []fyne.CanvasObject { return r.objs }

// Destroy releases renderer resources
func (r *progressRingRenderer) Destroy() {}

// onStartClicked starts or resumes timer
func (pw *PomodoroWindow) onStartClicked() {
	if pw.timer.State == models.PomodoroPaused {
		pw.timer.Resume()
	} else {
		pw.timer.Start()
	}

	pw.tick()
}

// onPauseClicked pauses or resumes timer
func (pw *PomodoroWindow) onPauseClicked() {
	if pw.timer.State == models.PomodoroPaused {
		pw.timer.Resume()
	} else {
		pw.timer.Pause()
	}

	pw.tick()
}

// onResetClicked returns timer to idle state
func (pw *PomodoroWindow) onResetClicked() {
	pw.timer.Reset()

	pw.tick()
}

// Show displays Pomodoro window
func (pw *PomodoroWindow) Show() {
	pw.window.Show()
}

// SetOnClosed registers window close callback
func (pw *PomodoroWindow) SetOnClosed(callback func()) {
	pw.window.SetOnClosed(func() {
		pw.stopTicker()

		if callback != nil {
			callback()
		}
	})
}

// UpdateTheme rebuilds window after theme switch
func (pw *PomodoroWindow) UpdateTheme(isGruvbox bool) {
	pw.isGruvbox = isGruvbox

	pw.setupUI()

	pw.tick()
}
