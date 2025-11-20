package ui

import (
	"fmt"
	"image/color"
	"time"

	"godo/src/localization"
	"godo/src/models"
	"godo/src/persistence"
	"godo/src/ui/helpers"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Timeline represents the timeline visualization widget
type Timeline struct {
	widget.BaseWidget

	dataManager persistence.TodoRepository
	currentDate time.Time // Changed to time.Time for daily view
	todos       []*models.TodoItem
	viewMode    models.ViewMode
	window      fyne.Window

	// Timeline state
	scrollPosition float32
	itemHeight     float32
	dateGroups     map[string][]*models.TodoItem
	visibleItems   []*models.TodoItem

	// Event callbacks
	onTodoSelected    func(*models.TodoItem, time.Time)
	onTodoReorder     func(*models.TodoItem, int) // delta: -1 up, +1 down
	onReorderFinished func()
	onTodosChanged    func()

	// drag state
	draggingTodo *models.TodoItem
}

// NewTimeline creates a new timeline widget
func NewTimeline(dataManager persistence.TodoRepository) *Timeline {
	t := &Timeline{
		dataManager:    dataManager,
		currentDate:    time.Now(), // Default to now
		viewMode:       models.ViewAll,
		scrollPosition: 0,
		itemHeight:     TimelineItemHeight,
		dateGroups:     make(map[string][]*models.TodoItem),
		visibleItems:   make([]*models.TodoItem, 0),
	}

	t.ExtendBaseWidget(t)
	return t
}

// SetDate sets the current viewing date
func (t *Timeline) SetDate(date time.Time) {
	t.currentDate = date
	// Don't auto-refresh - let caller control when to refresh
}

// SetViewMode sets the current view mode
func (t *Timeline) SetViewMode(mode models.ViewMode) {
	t.viewMode = mode
	// Don't auto-refresh - let caller control when to refresh
}

// SetTodos updates the todos and refreshes the display
func (t *Timeline) SetTodos(todos []*models.TodoItem) {
	t.todos = todos
	t.organizeByDate()
	// Don't auto-refresh - let caller control when to refresh
}

// SetWindow sets the parent window reference for dialogs.
func (t *Timeline) SetWindow(win fyne.Window) {
	t.window = win
}

// SetOnTodoReorder sets callback for explicit reorder actions (up/down or DnD)
func (t *Timeline) SetOnTodoReorder(callback func(*models.TodoItem, int)) {
	t.onTodoReorder = callback
}

// SetOnTodosChanged registers callback invoked when timeline mutates todo data.
func (t *Timeline) SetOnTodosChanged(callback func()) {
	t.onTodosChanged = callback
}

// organizeByDate groups todos by date for display
func (t *Timeline) organizeByDate() {
	t.dateGroups = make(map[string][]*models.TodoItem)
	t.visibleItems = make([]*models.TodoItem, 0)

	dateKey := t.currentDate.Format("2006-01-02")
	t.dateGroups[dateKey] = t.todos // Todos are already filtered by day
	t.visibleItems = t.todos
}

// CreateRenderer creates the widget renderer
func (t *Timeline) CreateRenderer() fyne.WidgetRenderer {
	t.ExtendBaseWidget(t)
	return &timelineRenderer{timeline: t}
}

// timelineRenderer handles the rendering of the timeline
type timelineRenderer struct {
	timeline *Timeline
	scroll   *container.Scroll
	listBox  *fyne.Container
}

func (r *timelineRenderer) Layout(size fyne.Size) {
	if r.scroll == nil {
		// Initialize scroll with persistent content container
		r.listBox = r.createTimelineContent()
		r.scroll = container.NewScroll(r.listBox)
	}
	// Ensure scroll fills all available area
	r.scroll.Resize(size)
}

func (r *timelineRenderer) MinSize() fyne.Size {
	// Always return a reasonable minimum size to prevent UI shrinking
	// The scroll container's MinSize of {32 32} is too small and causes layout issues
	return fyne.NewSize(350, 200)
}

func (r *timelineRenderer) Refresh() {
	// Update content in-place for smoother UI
	if r.scroll == nil || r.listBox == nil {
		r.listBox = r.createTimelineContent()
		r.scroll = container.NewScroll(r.listBox)
	} else {
		// Preserve scroll offset
		offset := r.scroll.Offset
		// Rebuild children objects
		objects := r.buildTimelineObjects()
		r.listBox.Objects = objects
		runOnMainThread(func() {
			r.listBox.Refresh()
		})
		r.scroll.Offset = offset
		if size := r.timeline.Size(); size.Width > 0 && size.Height > 0 {
			r.scroll.Resize(size)
		}
	}
}

func (r *timelineRenderer) BackgroundColor() fyne.ThemeColorName {
	// Transparent so the full-window background (gradient) is visible outside the card
	return ""
}

func (r *timelineRenderer) Objects() []fyne.CanvasObject {
	if r.scroll == nil {
		// Create the scrollable content
		content := r.createTimelineContent()
		r.scroll = container.NewScroll(content)
	}
	return []fyne.CanvasObject{r.scroll}
}

func (r *timelineRenderer) Destroy() {
	// Clean up resources if needed
}

func (r *timelineRenderer) createTimelineContent() *fyne.Container {
	objects := r.buildTimelineObjects()
	return container.NewVBox(objects...)
}

func (r *timelineRenderer) buildTimelineObjects() []fyne.CanvasObject {
	var objects []fyne.CanvasObject

	// Single date header for current day
	dateKey := r.timeline.currentDate.Format("2006-01-02")
	dateHeader := r.createDateHeader(dateKey)
	objects = append(objects, dateHeader)

	if len(r.timeline.visibleItems) == 0 {
		emptyLabel := widget.NewLabel(localization.GetString("status_empty_list"))
		emptyLabel.Alignment = fyne.TextAlignCenter
		objects = append(objects, emptyLabel)
		return objects
	}

	// Add todo items
	for _, todo := range r.timeline.todos {
		todoItem := r.createTodoItem(todo)
		objects = append(objects, todoItem)
	}
	return objects
}

func (r *timelineRenderer) getSortedDateKeys() []string {
	return []string{r.timeline.currentDate.Format("2006-01-02")}
}

func (r *timelineRenderer) createDateHeader(dateKey string) fyne.CanvasObject {
	date, err := time.Parse("2006-01-02", dateKey)
	if err != nil {
		return widget.NewLabel("Invalid Date")
	}
	isLightTheme := helpers.IsLightTheme()

	// Format date header like original: "2025年10月15日 星期一   +"
	weekdayNames := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	weekdayName := weekdayNames[date.Weekday()]

	headerText := fmt.Sprintf("%d/%02d/%02d %s",
		date.Year(), date.Month(), date.Day(), weekdayName)

	// Use canvas.Text to control size ~20px per mockup
	var fg color.Color
	if isLightTheme {
		fg = color.NRGBA{R: 0x3c, G: 0x38, B: 0x36, A: 0xFF} // #3c3836
	} else {
		fg = color.NRGBA{R: 0xeb, G: 0xdb, B: 0xb2, A: 0xFF} // #ebdbb2
	}
	headerLabel := canvas.NewText(headerText, fg)
	headerLabel.TextStyle = fyne.TextStyle{Bold: true}
	headerLabel.TextSize = 20

	// Center the date header within the tasks window
	headerLabel.Alignment = fyne.TextAlignCenter

	// Divider line separating header from tasks
	var dividerColor color.Color
	if isLightTheme {
		dividerColor = color.NRGBA{R: 0xD0, G: 0xD0, B: 0xD0, A: 0xFF}
	} else {
		dividerColor = color.NRGBA{R: 0x3c, G: 0x38, B: 0x36, A: 0xFF}
	}
	divider := canvas.NewRectangle(dividerColor)
	divider.SetMinSize(fyne.NewSize(10, 1))

	return container.NewVBox(
		helpers.CreateSpacer(1, 2),
		container.NewCenter(headerLabel),
		helpers.CreateSpacer(1, 3),
		divider,
	)
}

func (r *timelineRenderer) createTodoItem(todo *models.TodoItem) fyne.CanvasObject {
	isLightTheme := helpers.IsLightTheme()
	// Priority indicator: colored vertical RECTANGLE 16x32px (half width, same height)
	colorSquare := canvas.NewRectangle(todo.GetLevelColor())
	colorSquare.CornerRadius = 2 // Sharper corners (was 6)
	colorSquare.SetMinSize(fyne.NewSize(16, 32))
	colorSquareWrap := container.NewGridWrap(fyne.NewSize(16, 32), colorSquare)
	colorSquareAligned := verticallyCenterCompact(colorSquareWrap)

	// Custom square checkbox 20x20 per mockup, centered vertically
	doneCheck := newSquareCheckbox(todo.Done, func(checked bool) {
		updated := *todo
		updated.Done = checked
		_ = r.timeline.dataManager.UpdateTodo(&updated, todo.TodoTime)
		if todos, err := r.timeline.dataManager.GetTodosForMonth(r.timeline.currentDate.Year(), int(r.timeline.currentDate.Month())); err == nil {
			filtered := r.timeline.viewMode.FilterItems(todos, time.Now())
			r.timeline.SetTodos(filtered)
			r.timeline.Refresh()
		}
	})
	// Wrap checkbox in container for vertical centering
	doneCheckCentered := verticallyCenterCompact(doneCheck)

	// Todo name - takes the remaining space, 18px from mockup
	nameLabel := widget.NewLabel(todo.Name)
	nameLabel.Wrapping = fyne.TextWrapWord
	nameLabel.TextStyle = fyne.TextStyle{}

	// Time display - right-aligned, 18px from mockup
	var timeColor color.Color
	if isLightTheme {
		timeColor = color.NRGBA{R: 0x66, G: 0x66, B: 0x66, A: 0xFF} // #666
	} else {
		timeColor = color.NRGBA{R: 0xA8, G: 0x99, B: 0x84, A: 0xFF} // #a89984
	}
	timeText := canvas.NewText(fmt.Sprintf("%02d:%02d", todo.TodoTime.Hour(), todo.TodoTime.Minute()), timeColor)
	timeText.TextSize = 18
	timeLabel := verticallyCenterCompact(timeText)

	//Status indicator
	status := newStatusIndicator(todo, func(toggleStar bool) {
		if toggleStar {
			updated := *todo
			updated.Starred = !todo.Starred
			if err := r.timeline.dataManager.UpdateTodo(&updated, todo.TodoTime); err != nil {
				r.timeline.showError(err)
				return
			}
			r.timeline.notifyTodosChanged()
		}
	})
	// Keep status aligned with the schedule time for a cleaner row
	statusCentered := verticallyCenterCompact(status)

	// Delete action: show trash icon right of status
	deleteBtn := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		r.timeline.confirmDelete(todo)
	})
	deleteBtn.Importance = widget.LowImportance
	deleteBtnCentered := verticallyCenterCompact(deleteBtn)

	// Layout: [ColorSquare] [Spacer] [Checkbox] [Name................] [Time] [Star] [Delete]
	// Add spacer between color and checkbox (doubled spacing)
	leftSection := container.NewHBox(colorSquareAligned, helpers.CreateSpacer(8, 1), doneCheckCentered)
	rightSection := container.NewHBox(timeLabel, helpers.CreateSpacer(8, 1), statusCentered, helpers.CreateSpacer(4, 1), deleteBtnCentered)
	content := container.NewBorder(nil, nil, leftSection, rightSection, verticallyCenterWide(nameLabel))

	// Row with bottom border only (no card)
	var borderClr color.Color
	if isLightTheme {
		borderClr = color.NRGBA{R: 0xD0, G: 0xD0, B: 0xD0, A: 0xFF} // #d0d0d0 - darker gray for visibility
	} else {
		borderClr = color.NRGBA{R: 0x3c, G: 0x38, B: 0x36, A: 0xFF}
	}
	bottomLine := canvas.NewRectangle(borderClr)
	bottomLine.SetMinSize(fyne.NewSize(10, 1))
	row := container.NewVBox(
		container.NewPadded(content),
		bottomLine,
	)

	// Shadow layer (hidden by default, shown during drag)
	shadowRect := canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 0})
	shadowRect.CornerRadius = 4

	// Border layer (hidden by default, shown during drag)
	borderRect := canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 0})
	borderRect.CornerRadius = 4

	// Press/drag overlay for row feedback
	pressOverlay := canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 0})
	pressOverlay.CornerRadius = 4

	// Stack all layers: shadow (bottom) -> border -> row -> press overlay (top)
	pressStack := container.NewMax(shadowRect, borderRect, row, pressOverlay)

	// If this is the currently dragged item - apply drag effects
	isDragging := r.timeline.draggingTodo == todo
	if isDragging {
		// Enhanced highlight with higher opacity and border
		pressOverlay.FillColor = color.NRGBA{R: 0x3C, G: 0x82, B: 0xFF, A: 80}

		// Border effect - blue outline
		borderRect.FillColor = color.Transparent
		borderRect.StrokeColor = color.NRGBA{R: 0x3C, G: 0x82, B: 0xFF, A: 200}
		borderRect.StrokeWidth = 2

		// Shadow effect - dark gradient
		shadowRect.FillColor = color.NRGBA{R: 0, G: 0, B: 0, A: 30}

		runOnMainThread(func() {
			pressOverlay.Refresh()
			borderRect.Refresh()
			shadowRect.Refresh()
		})
	}

	// Make the entire item clickable
	tappable := &tappableTodo{
		todo:         todo,
		todoTime:     todo.TodoTime,
		container:    pressStack,
		press:        pressOverlay,
		shadow:       shadowRect,
		border:       borderRect,
		onSelected:   r.timeline.onTodoSelected,
		onReorder:    r.timeline.onTodoReorder,
		rowThresh:    60,
		onReorderEnd: r.timeline.onReorderFinished,
		timeline:     r.timeline,
	}
	tappable.ExtendBaseWidget(tappable)

	return tappable
}

// verticallyCenterCompact keeps the child at its natural size while centering it vertically.
func verticallyCenterCompact(obj fyne.CanvasObject) fyne.CanvasObject {
	return container.NewVBox(layout.NewSpacer(), container.NewCenter(obj), layout.NewSpacer())
}

// verticallyCenterWide centers the child vertically but allows it to stretch horizontally.
func verticallyCenterWide(obj fyne.CanvasObject) fyne.CanvasObject {
	return container.NewVBox(layout.NewSpacer(), obj, layout.NewSpacer())
}

// tappableTodo makes todo items clickable
type tappableTodo struct {
	widget.BaseWidget
	todo         *models.TodoItem
	todoTime     time.Time
	container    *fyne.Container
	press        *canvas.Rectangle
	shadow       *canvas.Rectangle // Shadow effect during drag
	border       *canvas.Rectangle // Border effect during drag
	originalSize fyne.Size         // Original size before scaling
	onSelected   func(*models.TodoItem, time.Time)
	onReorder    func(*models.TodoItem, int)
	onReorderEnd func()
	dragAccumY   float32
	rowThresh    float32
	dragging     bool
	lastMouseY   float32
	timeline     *Timeline
}

func (tt *tappableTodo) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(tt.container)
}

// Cursor returns a hand/move cursor to indicate draggable item
func (tt *tappableTodo) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

func (tt *tappableTodo) Tapped(*fyne.PointEvent) {
	// flash overlay
	if tt.press != nil {
		col := helpers.ToNRGBA(theme.Color(theme.ColorNameHover))
		if col.A < 60 {
			col.A = 60
		}
		tt.press.FillColor = col
		runOnMainThread(func() {
			tt.press.Refresh()
		})
		go func(p *canvas.Rectangle) {
			time.Sleep(120 * time.Millisecond)
			runOnMainThread(func() {
				p.FillColor = color.NRGBA{R: 0, G: 0, B: 0, A: 0}
				p.Refresh()
			})
		}(tt.press)
	}
	if tt.onSelected != nil {
		tt.onSelected(tt.todo, tt.todoTime)
	}
}

// Implement fyne.Draggable
func (tt *tappableTodo) Dragged(e *fyne.DragEvent) {
	// Mark as dragging and apply visual effects on first drag event
	if !tt.dragging {
		tt.dragging = true

		// Enhanced highlight with higher opacity (80 instead of 50)
		if tt.press != nil {
			tt.press.FillColor = color.NRGBA{R: 0x3C, G: 0x82, B: 0xFF, A: 80}
			runOnMainThread(func() {
				tt.press.Refresh()
			})
		}

		// Border effect - blue outline
		if tt.border != nil {
			tt.border.FillColor = color.Transparent
			tt.border.StrokeColor = color.NRGBA{R: 0x3C, G: 0x82, B: 0xFF, A: 200}
			tt.border.StrokeWidth = 2
			runOnMainThread(func() {
				tt.border.Refresh()
			})
		}

		// Shadow effect - dark semi-transparent
		if tt.shadow != nil {
			tt.shadow.FillColor = color.NRGBA{R: 0, G: 0, B: 0, A: 30}
			runOnMainThread(func() {
				tt.shadow.Refresh()
			})
		}

		if tt.timeline != nil {
			tt.timeline.draggingTodo = tt.todo
		}
	}

	// Process drag movement
	tt.dragAccumY += e.Dragged.DY
	for tt.dragAccumY > tt.rowThresh {
		if tt.onReorder != nil {
			tt.onReorder(tt.todo, 1)
		}
		tt.dragAccumY -= tt.rowThresh
	}
	for tt.dragAccumY < -tt.rowThresh {
		if tt.onReorder != nil {
			tt.onReorder(tt.todo, -1)
		}
		tt.dragAccumY += tt.rowThresh
	}
}

func (t *Timeline) confirmDelete(todo *models.TodoItem) {
	if todo == nil {
		return
	}
	if t.window == nil {
		t.deleteTodo(todo)
		return
	}

	var restoreTheme func()
	currentTheme := fyne.CurrentApp().Settings().Theme()
	isLightTheme := helpers.IsLightTheme()
	messageColor := helpers.ToNRGBA(theme.Color(theme.ColorNameForeground))
	overrideTheme := &foregroundOverrideTheme{
		base:         currentTheme,
		headingDelta: 5,
		textDelta:    5,
	}
	if isLightTheme {
		messageColor = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
		overrideTheme.overrideForeground = true
		overrideTheme.overrideColor = messageColor
	}
	fyne.CurrentApp().Settings().SetTheme(overrideTheme)
	restoreTheme = func() {
		fyne.CurrentApp().Settings().SetTheme(currentTheme)
	}
	message := canvas.NewText(localization.GetString("confirm_delete_message"), messageColor)
	message.Alignment = fyne.TextAlignLeading
	message.TextSize = 16
	content := container.NewVBox(
		container.NewHBox(message, layout.NewSpacer()),
		helpers.CreateSpacer(1, 20),
	)
	conf := dialog.NewCustomConfirm(
		localization.GetString("confirm_delete_title"),
		localization.GetString("shortcut_delete"),
		localization.GetString("form_button_cancel"),
		content,
		func(confirm bool) {
			if confirm {
				t.deleteTodo(todo)
			}
		},
		t.window,
	)
	conf.SetOnClosed(func() {
		if restoreTheme != nil {
			restoreTheme()
		}
	})
	conf.Show()
}

func (t *Timeline) deleteTodo(todo *models.TodoItem) {
	if todo == nil {
		return
	}
	if err := t.dataManager.RemoveTodo(todo.TodoTime); err != nil {
		t.showError(err)
		return
	}
	t.notifyTodosChanged()
}

func (t *Timeline) notifyTodosChanged() {
	if t.onTodosChanged != nil {
		t.onTodosChanged()
		return
	}
	t.reloadVisibleTodos()
}

func (t *Timeline) reloadVisibleTodos() {
	year, month := t.currentDate.Year(), int(t.currentDate.Month())
	todos, err := t.dataManager.GetTodosForMonth(year, month)
	if err != nil {
		t.showError(err)
		return
	}
	filtered := t.viewMode.FilterItems(todos, time.Now())
	t.SetTodos(filtered)
	t.Refresh()
}

func (t *Timeline) showError(err error) {
	if err == nil {
		return
	}
	if t.window != nil {
		dialog.NewError(err, t.window).Show()
		return
	}
	fmt.Println("timeline error:", err)
}

type foregroundOverrideTheme struct {
	base               fyne.Theme
	overrideColor      color.Color
	overrideForeground bool
	headingDelta       float32
	textDelta          float32
}

func (t *foregroundOverrideTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameForeground && t.overrideForeground {
		return t.overrideColor
	}
	if t.base != nil {
		return t.base.Color(name, variant)
	}
	return theme.DefaultTheme().Color(name, variant)
}

func (t *foregroundOverrideTheme) Font(style fyne.TextStyle) fyne.Resource {
	if t.base != nil {
		return t.base.Font(style)
	}
	return theme.DefaultTheme().Font(style)
}

func (t *foregroundOverrideTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	if t.base != nil {
		return t.base.Icon(name)
	}
	return theme.DefaultTheme().Icon(name)
}

func (t *foregroundOverrideTheme) Size(name fyne.ThemeSizeName) float32 {
	var size float32
	if t.base != nil {
		size = t.base.Size(name)
	} else {
		size = theme.DefaultTheme().Size(name)
	}
	if name == theme.SizeNameHeadingText {
		size += t.headingDelta
	}
	if name == theme.SizeNameText {
		size += t.textDelta
	}
	return size
}

func (tt *tappableTodo) DragEnd() {
	tt.dragging = false
	tt.dragAccumY = 0

	// Clear all drag visual effects
	runOnMainThread(func() {
		if tt.press != nil {
			tt.press.FillColor = color.NRGBA{R: 0, G: 0, B: 0, A: 0}
			tt.press.Refresh()
		}
		if tt.border != nil {
			tt.border.FillColor = color.NRGBA{R: 0, G: 0, B: 0, A: 0}
			tt.border.StrokeColor = color.NRGBA{R: 0, G: 0, B: 0, A: 0}
			tt.border.StrokeWidth = 0
			tt.border.Refresh()
		}
		if tt.shadow != nil {
			tt.shadow.FillColor = color.NRGBA{R: 0, G: 0, B: 0, A: 0}
			tt.shadow.Refresh()
		}
	})

	if tt.onReorderEnd != nil {
		tt.onReorderEnd()
	}
	if tt.timeline != nil {
		tt.timeline.draggingTodo = nil
		// Force rebuild so any drag highlight applied during re-render is cleared
		tt.timeline.Refresh()
	}
}

// Implement desktop.Mouseable to ensure drag works well inside Scroll
func (tt *tappableTodo) MouseDown(e *desktop.MouseEvent) {
	tt.dragging = true
	tt.lastMouseY = e.Position.Y
	// Highlight row while dragging
	if tt.press != nil {
		// Blue tint with ~20% opacity for both themes
		tt.press.FillColor = color.NRGBA{R: 0x3C, G: 0x82, B: 0xFF, A: 50}
		runOnMainThread(func() {
			tt.press.Refresh()
		})
	}
	if tt.timeline != nil {
		tt.timeline.draggingTodo = tt.todo
	}
}

func (tt *tappableTodo) MouseMoved(e *desktop.MouseEvent) {
	if !tt.dragging {
		return
	}
	dy := e.Position.Y - tt.lastMouseY
	tt.lastMouseY = e.Position.Y
	tt.dragAccumY += dy
	for tt.dragAccumY > tt.rowThresh {
		if tt.onReorder != nil {
			tt.onReorder(tt.todo, 1)
		}
		tt.dragAccumY -= tt.rowThresh
	}
	for tt.dragAccumY < -tt.rowThresh {
		if tt.onReorder != nil {
			tt.onReorder(tt.todo, -1)
		}
		tt.dragAccumY += tt.rowThresh
	}
}

func (tt *tappableTodo) MouseUp(*desktop.MouseEvent) {
	tt.dragging = false
	tt.dragAccumY = 0

	// Clear all drag visual effects
	runOnMainThread(func() {
		if tt.press != nil {
			tt.press.FillColor = color.NRGBA{R: 0, G: 0, B: 0, A: 0}
			tt.press.Refresh()
		}
		if tt.border != nil {
			tt.border.FillColor = color.NRGBA{R: 0, G: 0, B: 0, A: 0}
			tt.border.StrokeColor = color.NRGBA{R: 0, G: 0, B: 0, A: 0}
			tt.border.StrokeWidth = 0
			tt.border.Refresh()
		}
		if tt.shadow != nil {
			tt.shadow.FillColor = color.NRGBA{R: 0, G: 0, B: 0, A: 0}
			tt.shadow.Refresh()
		}
	})

	if tt.onReorderEnd != nil {
		tt.onReorderEnd()
	}
	if tt.timeline != nil {
		tt.timeline.draggingTodo = nil
		// Force rebuild so any drag highlight applied during re-render is cleared
		tt.timeline.Refresh()
	}
}

// Custom square checkbox per mockup
type squareCheckbox struct {
	widget.BaseWidget
	checked   bool
	onChanged func(bool)
	rect      *canvas.Rectangle
	tick      *canvas.Text
	overlay   *canvas.Rectangle
	cont      *fyne.Container
}

func newSquareCheckbox(initial bool, onChanged func(bool)) *squareCheckbox {
	c := &squareCheckbox{checked: initial, onChanged: onChanged}
	c.ExtendBaseWidget(c)
	return c
}

func (c *squareCheckbox) CreateRenderer() fyne.WidgetRenderer {
	// colors
	isLightTheme := helpers.IsLightTheme()
	var border color.Color
	var fillChecked color.Color
	if isLightTheme {
		border = color.NRGBA{R: 0xD0, G: 0xD0, B: 0xD0, A: 0xFF}
		fillChecked = color.NRGBA{R: 0x4C, G: 0xAF, B: 0x50, A: 0xFF} // #4caf50
	} else {
		border = color.NRGBA{R: 0x66, G: 0x5c, B: 0x54, A: 0xFF}
		fillChecked = color.NRGBA{R: 0x98, G: 0x97, B: 0x1a, A: 0xFF} // #98971a
	}
	c.rect = canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 0})
	c.rect.StrokeColor = border
	c.rect.StrokeWidth = 2
	c.rect.CornerRadius = 3
	c.rect.SetMinSize(fyne.NewSize(20, 20))

	c.tick = canvas.NewText("", color.White)
	if c.checked {
		c.rect.FillColor = fillChecked
		c.tick.Text = "✓"
	}
	c.overlay = canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 0})
	c.cont = container.NewGridWrap(
		fyne.NewSize(20, 20),
		container.NewMax(c.rect, container.NewCenter(c.tick), c.overlay),
	)
	return widget.NewSimpleRenderer(c.cont)
}

func (c *squareCheckbox) MinSize() fyne.Size { return fyne.NewSize(20, 20) }

func (c *squareCheckbox) Tapped(*fyne.PointEvent) {
	c.checked = !c.checked
	if c.onChanged != nil {
		c.onChanged(c.checked)
	}
	runOnMainThread(func() {
		c.Refresh()
	})
}

func (c *squareCheckbox) Refresh() {
	// update visuals according to checked state and theme
	isLightTheme := helpers.IsLightTheme()
	var fillChecked color.Color
	if isLightTheme {
		fillChecked = color.NRGBA{R: 0x4C, G: 0xAF, B: 0x50, A: 0xFF}
	} else {
		fillChecked = color.NRGBA{R: 0x98, G: 0x97, B: 0x1a, A: 0xFF}
	}
	runOnMainThread(func() {
		if c.rect != nil && c.tick != nil {
			if c.checked {
				c.rect.FillColor = fillChecked
				c.tick.Text = "✓"
			} else {
				c.rect.FillColor = color.NRGBA{R: 0, G: 0, B: 0, A: 0}
				c.tick.Text = ""
			}
			c.rect.Refresh()
			c.tick.Refresh()
		}
		c.BaseWidget.Refresh()
	})
}

// Status indicator (✓ or ★), star toggles on tap when shown
type statusIndicator struct {
	widget.BaseWidget
	todo     *models.TodoItem
	onToggle func(toggleStar bool)
	overlay  *canvas.Rectangle
	cont     *fyne.Container
}

func newStatusIndicator(todo *models.TodoItem, onToggle func(bool)) *statusIndicator {
	s := &statusIndicator{todo: todo, onToggle: onToggle}
	s.ExtendBaseWidget(s)
	return s
}

func (s *statusIndicator) CreateRenderer() fyne.WidgetRenderer {
	var txt string
	var col color.Color
	isLightTheme := helpers.IsLightTheme()
	if s.todo.Done {
		// Completed: show check mark
		txt = "✓"
		if isLightTheme {
			col = color.NRGBA{R: 0x4C, G: 0xAF, B: 0x50, A: 0xFF}
		} else {
			col = color.NRGBA{R: 0x98, G: 0x97, B: 0x1a, A: 0xFF}
		}
	} else {
		// Not done: always show a star
		txt = "★"
		if s.todo.Starred {
			// Starred color per theme
			if isLightTheme {
				// Light theme: blue star
				col = color.NRGBA{R: 0x3C, G: 0x82, B: 0xFF, A: 0xFF} // #3c82ff
			} else {
				// Dark theme: orange star
				col = color.NRGBA{R: 0xFE, G: 0x80, B: 0x19, A: 0xFF} // #fe8019
			}
		} else {
			// Default grey star
			col = color.NRGBA{R: 0x9E, G: 0x9E, B: 0x9E, A: 0xFF} // #9e9e9e
		}
	}
	t := canvas.NewText(txt, col)
	t.TextSize = 20
	s.overlay = canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 0})
	s.cont = container.NewMax(container.NewCenter(t), s.overlay)
	return widget.NewSimpleRenderer(s.cont)
}

func (s *statusIndicator) Tapped(*fyne.PointEvent) {
	if !s.todo.Done && s.onToggle != nil {
		s.onToggle(true)
	}
}

func (s *statusIndicator) MinSize() fyne.Size { return fyne.NewSize(20, 20) }

// SetOnTodoSelected sets the callback for when a todo is selected
func (t *Timeline) SetOnTodoSelected(callback func(*models.TodoItem, time.Time)) {
	t.onTodoSelected = callback
}

// SetOnReorderFinished sets callback when drag/reorder ends
func (t *Timeline) SetOnReorderFinished(callback func()) {
	t.onReorderFinished = callback
}

// ScrollToTop scrolls to the top of the timeline
func (t *Timeline) ScrollToTop() {
	if renderer := t.CreateRenderer(); renderer != nil {
		if scroll, ok := renderer.(*timelineRenderer); ok {
			scroll.scroll.ScrollToTop()
		}
	}
}
