package ui

import (
	"fmt"
	"image/color"
	"time"

	"todo-list-migration/src/localization"
	"todo-list-migration/src/models"
	"todo-list-migration/src/persistence"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Timeline represents the timeline visualization widget
type Timeline struct {
	widget.BaseWidget

	dataManager *persistence.MonthlyManager
	currentDate time.Time // Changed to time.Time for daily view
	todos       []*models.TodoItem
	viewMode    models.ViewMode

	// Timeline state
	scrollPosition float32
	itemHeight     float32
	dateGroups     map[string][]*models.TodoItem
	visibleItems   []*models.TodoItem

	// Event callbacks
	onTodoSelected    func(*models.TodoItem, time.Time)
	onTodoReorder     func(*models.TodoItem, int) // delta: -1 up, +1 down
	onReorderFinished func()

	// drag state
	draggingTodo *models.TodoItem
}

// NewTimeline creates a new timeline widget
func NewTimeline(dataManager *persistence.MonthlyManager) *Timeline {
	t := &Timeline{
		dataManager:    dataManager,
		currentDate:    time.Now(), // Default to now
		viewMode:       models.ViewAll,
		scrollPosition: 0,
		itemHeight:     80,
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

// SetOnTodoReorder sets callback for explicit reorder actions (up/down or DnD)
func (t *Timeline) SetOnTodoReorder(callback func(*models.TodoItem, int)) {
	t.onTodoReorder = callback
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
		r.listBox.Refresh()
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

	// Format date header like original: "2025年10月15日 星期一   +"
	weekdayNames := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	weekdayName := weekdayNames[date.Weekday()]

	headerText := fmt.Sprintf("%d/%02d/%02d %s",
		date.Year(), date.Month(), date.Day(), weekdayName)

	// Use canvas.Text to control size ~20px per mockup
	var fg color.Color
	if _, ok := fyne.CurrentApp().Settings().Theme().(*LightSoftTheme); ok {
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
	if _, ok := fyne.CurrentApp().Settings().Theme().(*LightSoftTheme); ok {
		dividerColor = color.NRGBA{R: 0xD0, G: 0xD0, B: 0xD0, A: 0xFF}
	} else {
		dividerColor = color.NRGBA{R: 0x3c, G: 0x38, B: 0x36, A: 0xFF}
	}
	divider := canvas.NewRectangle(dividerColor)
	divider.SetMinSize(fyne.NewSize(10, 1))

	return container.NewVBox(
		CreateSpacer(1, 2),
		container.NewCenter(headerLabel),
		CreateSpacer(1, 3),
		divider,
	)
}

func (r *timelineRenderer) createTodoItem(todo *models.TodoItem) fyne.CanvasObject {
	// Priority indicator: colored vertical RECTANGLE 16x32px (half width, same height)
	colorSquare := canvas.NewRectangle(todo.GetLevelColor())
	colorSquare.CornerRadius = 2 // Sharper corners (was 6)
	colorSquare.SetMinSize(fyne.NewSize(16, 32))
	colorSquareWrap := container.NewGridWrap(fyne.NewSize(16, 32), colorSquare)

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
	doneCheckCentered := container.NewCenter(doneCheck)

	// Todo name - takes the remaining space, 18px from mockup
	nameLabel := widget.NewLabel(todo.Name)
	nameLabel.Wrapping = fyne.TextWrapWord
	nameLabel.TextStyle = fyne.TextStyle{}

	// Time display - right-aligned, 18px from mockup
	var timeColor color.Color
	if _, ok := fyne.CurrentApp().Settings().Theme().(*LightSoftTheme); ok {
		timeColor = color.NRGBA{R: 0x66, G: 0x66, B: 0x66, A: 0xFF} // #666
	} else {
		timeColor = color.NRGBA{R: 0xA8, G: 0x99, B: 0x84, A: 0xFF} // #a89984
	}
	timeText := canvas.NewText(fmt.Sprintf("%02d:%02d", todo.TodoTime.Hour(), todo.TodoTime.Minute()), timeColor)
	timeText.TextSize = 18
	timeLabel := container.NewCenter(timeText)

	// Status: ✓ if done, else ★ if starred; toggle star on tap when not done
	status := newStatusIndicator(todo, func(toggleStar bool) {
		if toggleStar {
			updated := *todo
			updated.Starred = !todo.Starred
			_ = r.timeline.dataManager.UpdateTodo(&updated, todo.TodoTime)
			if todos, err := r.timeline.dataManager.GetTodosForMonth(r.timeline.currentDate.Year(), int(r.timeline.currentDate.Month())); err == nil {
				filtered := r.timeline.viewMode.FilterItems(todos, time.Now())
				r.timeline.SetTodos(filtered)
				r.timeline.Refresh()
			}
		}
	})
	// Push star down to center it vertically in the row
	statusCentered := container.NewVBox(CreateSpacer(1, 7), status)

	// Layout: [ColorSquare] [Spacer] [Checkbox] [Name................] [Time] [Star] [Delete]
	// Add spacer between color and checkbox (doubled spacing)
	leftSection := container.NewHBox(colorSquareWrap, CreateSpacer(8, 1), doneCheckCentered)
	rightSection := container.NewHBox(timeLabel, CreateSpacer(8, 1), statusCentered)
	content := container.NewBorder(nil, nil, leftSection, rightSection, nameLabel)

	// Row with bottom border only (no card)
	var borderClr color.Color
	if _, ok := fyne.CurrentApp().Settings().Theme().(*LightSoftTheme); ok {
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

	// press/drag overlay for row feedback
	pressOverlay := canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 0})
	pressOverlay.CornerRadius = 0
	pressStack := container.NewMax(row, pressOverlay)

	// If this is the currently dragged item - keep it highlighted during drag
	if r.timeline.draggingTodo == todo {
		// Slight blue tint with high transparency (~20% opacity)
		pressOverlay.FillColor = color.NRGBA{R: 0x3C, G: 0x82, B: 0xFF, A: 50}
		pressOverlay.Refresh()
	}

	// Make the entire item clickable
	tappable := &tappableTodo{
		todo:         todo,
		todoTime:     todo.TodoTime,
		container:    pressStack,
		press:        pressOverlay,
		onSelected:   r.timeline.onTodoSelected,
		onReorder:    r.timeline.onTodoReorder,
		rowThresh:    60,
		onReorderEnd: r.timeline.onReorderFinished,
		timeline:     r.timeline,
	}
	tappable.ExtendBaseWidget(tappable)

	return tappable
}

// tappableTodo makes todo items clickable
type tappableTodo struct {
	widget.BaseWidget
	todo         *models.TodoItem
	todoTime     time.Time
	container    *fyne.Container
	press        *canvas.Rectangle
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
		col := toNRGBA(theme.Color(theme.ColorNameHover))
		if col.A < 60 {
			col.A = 60
		}
		tt.press.FillColor = col
		tt.press.Refresh()
		go func(p *canvas.Rectangle) {
			time.Sleep(120 * time.Millisecond)
			p.FillColor = color.NRGBA{R: 0, G: 0, B: 0, A: 0}
			p.Refresh()
		}(tt.press)
	}
	if tt.onSelected != nil {
		tt.onSelected(tt.todo, tt.todoTime)
	}
}

// Implement fyne.Draggable
func (tt *tappableTodo) Dragged(e *fyne.DragEvent) {
	tt.dragAccumY += e.Dragged.DY
	// Move downwards
	for tt.dragAccumY > tt.rowThresh {
		if tt.onReorder != nil {
			tt.onReorder(tt.todo, 1)
		}
		tt.dragAccumY -= tt.rowThresh
	}
	// Move upwards
	for tt.dragAccumY < -tt.rowThresh {
		if tt.onReorder != nil {
			tt.onReorder(tt.todo, -1)
		}
		tt.dragAccumY += tt.rowThresh
	}
}

func (tt *tappableTodo) DragEnd() {
	tt.dragAccumY = 0
	if tt.onReorderEnd != nil {
		tt.onReorderEnd()
	}
	if tt.timeline != nil {
		tt.timeline.draggingTodo = nil
		// Force rebuild so any drag highlight applied during re-render is cleared
		tt.timeline.Refresh()
	}
	// Ensure highlight is cleared even if MouseUp was not delivered
	if tt.press != nil {
		tt.press.FillColor = color.NRGBA{R: 0, G: 0, B: 0, A: 0}
		tt.press.Refresh()
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
		tt.press.Refresh()
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
	if tt.onReorderEnd != nil {
		tt.onReorderEnd()
	}
	// Remove highlight
	if tt.press != nil {
		tt.press.FillColor = color.NRGBA{R: 0, G: 0, B: 0, A: 0}
		tt.press.Refresh()
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
	var border color.Color
	var fillChecked color.Color
	if _, ok := fyne.CurrentApp().Settings().Theme().(*LightSoftTheme); ok {
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
	c.Refresh()
}

func (c *squareCheckbox) Refresh() {
	// update visuals according to checked state and theme
	var fillChecked color.Color
	if _, ok := fyne.CurrentApp().Settings().Theme().(*LightSoftTheme); ok {
		fillChecked = color.NRGBA{R: 0x4C, G: 0xAF, B: 0x50, A: 0xFF}
	} else {
		fillChecked = color.NRGBA{R: 0x98, G: 0x97, B: 0x1a, A: 0xFF}
	}
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
	if s.todo.Done {
		// Completed: show check mark
		txt = "✓"
		if _, ok := fyne.CurrentApp().Settings().Theme().(*LightSoftTheme); ok {
			col = color.NRGBA{R: 0x4C, G: 0xAF, B: 0x50, A: 0xFF}
		} else {
			col = color.NRGBA{R: 0x98, G: 0x97, B: 0x1a, A: 0xFF}
		}
	} else {
		// Not done: always show a star
		txt = "★"
		if s.todo.Starred {
			// Starred color per theme
			if _, ok := fyne.CurrentApp().Settings().Theme().(*LightSoftTheme); ok {
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
