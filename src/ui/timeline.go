package ui

import (
	"fmt"
	"image/color"
	"time"

	"todo-list-migration/src/localization"
	"todo-list-migration/src/models"
	"todo-list-migration/src/persistence"
	"todo-list-migration/src/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Timeline represents the timeline visualization widget
type Timeline struct {
	widget.BaseWidget

	dataManager *persistence.MonthlyManager
	currentDate *utils.CustomDate
	todos       []*models.TodoItem
	viewMode    models.ViewMode

	// Timeline state
	scrollPosition float32
	itemHeight     float32
	dateGroups     map[string][]*models.TodoItem
	visibleItems   []*models.TodoItem

	// Event callbacks
	onTodoSelected func(*models.TodoItem, time.Time)
}

// NewTimeline creates a new timeline widget
func NewTimeline(dataManager *persistence.MonthlyManager) *Timeline {
	t := &Timeline{
		dataManager:    dataManager,
		currentDate:    utils.GetCurrentDate(),
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
func (t *Timeline) SetDate(date *utils.CustomDate) {
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

// organizeByDate groups todos by date for display
func (t *Timeline) organizeByDate() {
	t.dateGroups = make(map[string][]*models.TodoItem)
	t.visibleItems = make([]*models.TodoItem, 0)

	// currentTime no longer needed here; filtering handled in SetTodos/FilterItems

	for _, todo := range t.todos {
		// Apply view mode filtering
		if t.viewMode == models.ViewIncomplete && todo.Done {
			continue
		}
		// Other modes are handled in mw.loadTodos() via FilterItems.

		dateKey := todo.TodoTime.Format("2006-01-02")
		t.dateGroups[dateKey] = append(t.dateGroups[dateKey], todo)
		t.visibleItems = append(t.visibleItems, todo)
	}

	// Sort todos within each date group by time (reverse chronological)
	for dateKey := range t.dateGroups {
		todos := t.dateGroups[dateKey]
		// Sort by time descending (most recent first)
		for i := 0; i < len(todos)-1; i++ {
			for j := i + 1; j < len(todos); j++ {
				if todos[i].TodoTime.Before(todos[j].TodoTime) {
					todos[i], todos[j] = todos[j], todos[i]
				}
			}
		}
	}
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
}

func (r *timelineRenderer) Layout(size fyne.Size) {
	if r.scroll == nil {
		// Initialize scroll if not already done
		content := r.createTimelineContent()
		r.scroll = container.NewScroll(content)
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
	// Store current scroll position
	var offset fyne.Position
	if r.scroll != nil {
		offset = r.scroll.Offset
	}

	// Recreate scroll container entirely to avoid layout bugs
	content := r.createTimelineContent()
	r.scroll = container.NewScroll(content)

	// Restore scroll position
	r.scroll.Offset = offset

	// Ensure proper sizing
	if size := r.timeline.Size(); size.Width > 0 && size.Height > 0 {
		r.scroll.Resize(size)
	}
}

func (r *timelineRenderer) BackgroundColor() fyne.ThemeColorName {
	return theme.ColorNameBackground
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
	var objects []fyne.CanvasObject

	// Add empty state if no todos
	if len(r.timeline.visibleItems) == 0 {
		emptyLabel := widget.NewLabel(localization.GetString("status_empty_list"))
		emptyLabel.Alignment = fyne.TextAlignCenter
		objects = append(objects, emptyLabel)
		return container.NewVBox(objects...)
	}

	// Group todos by date and create visual groups
	dateKeys := r.getSortedDateKeys()

	for _, dateKey := range dateKeys {
		todos := r.timeline.dateGroups[dateKey]

		// Add date header
		dateHeader := r.createDateHeader(dateKey)
		objects = append(objects, dateHeader)

		// Add todo items for this date
		for _, todo := range todos {
			todoItem := r.createTodoItem(todo)
			objects = append(objects, todoItem)
		}
	}

	return container.NewVBox(objects...)
}

func (r *timelineRenderer) getSortedDateKeys() []string {
	var keys []string
	for key := range r.timeline.dateGroups {
		keys = append(keys, key)
	}

	// Sort dates in descending order (most recent first)
	for i := 0; i < len(keys)-1; i++ {
		for j := i + 1; j < len(keys); j++ {
			date1, _ := time.Parse("2006-01-02", keys[i])
			date2, _ := time.Parse("2006-01-02", keys[j])
			if date1.Before(date2) {
				keys[i], keys[j] = keys[j], keys[i]
			}
		}
	}

	return keys
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

	return container.NewVBox(
		headerLabel,
		CreateSpacer(1, 10),
	)
}

func (r *timelineRenderer) createTodoItem(todo *models.TodoItem) fyne.CanvasObject {
	// Priority indicator: colored SQUARE 32x32px from mockup
	colorSquare := canvas.NewRectangle(todo.GetLevelColor())
	colorSquare.CornerRadius = 6
	colorSquare.SetMinSize(fyne.NewSize(32, 32))
	colorSquareWrap := container.NewGridWrap(fyne.NewSize(32, 32), colorSquare)

	// Custom square checkbox 20x20 per mockup
	doneCheck := newSquareCheckbox(todo.Done, func(checked bool) {
		updated := *todo
		updated.Done = checked
		_ = r.timeline.dataManager.UpdateTodo(&updated, todo.TodoTime)
		if todos, err := r.timeline.dataManager.GetTodosForMonth(r.timeline.currentDate.Year, r.timeline.currentDate.Month); err == nil {
			filtered := r.timeline.viewMode.FilterItems(todos, time.Now())
			r.timeline.SetTodos(filtered)
		}
	})

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
			if todos, err := r.timeline.dataManager.GetTodosForMonth(r.timeline.currentDate.Year, r.timeline.currentDate.Month); err == nil {
				filtered := r.timeline.viewMode.FilterItems(todos, time.Now())
				r.timeline.SetTodos(filtered)
			}
		}
	})

	// Layout: [ColorSquare] [Checkbox] [Name................] [Time] [Star] [Delete]
	leftSection := container.NewHBox(colorSquareWrap, doneCheck)
	rightSection := container.NewHBox(timeLabel, status)
	content := container.NewBorder(nil, nil, leftSection, rightSection, nameLabel)

	// Row with bottom border only (no card)
	var borderClr color.Color
	if _, ok := fyne.CurrentApp().Settings().Theme().(*LightSoftTheme); ok {
		borderClr = color.NRGBA{R: 0xF0, G: 0xF0, B: 0xF0, A: 0xFF}
	} else {
		borderClr = color.NRGBA{R: 0x3c, G: 0x38, B: 0x36, A: 0xFF}
	}
	bottomLine := canvas.NewRectangle(borderClr)
	bottomLine.SetMinSize(fyne.NewSize(10, 1))
	row := container.NewVBox(
		container.NewPadded(content),
		bottomLine,
	)

	// Make the entire item clickable
	tappable := &tappableTodo{
		todo:       todo,
		todoTime:   todo.TodoTime,
		container:  row,
		onSelected: r.timeline.onTodoSelected,
	}
	tappable.ExtendBaseWidget(tappable)

	return tappable
}

// tappableTodo makes todo items clickable
type tappableTodo struct {
	widget.BaseWidget
	todo       *models.TodoItem
	todoTime   time.Time
	container  *fyne.Container
	onSelected func(*models.TodoItem, time.Time)
}

func (tt *tappableTodo) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(tt.container)
}

func (tt *tappableTodo) Tapped(*fyne.PointEvent) {
	if tt.onSelected != nil {
		tt.onSelected(tt.todo, tt.todoTime)
	}
}

// Custom square checkbox per mockup
type squareCheckbox struct {
	widget.BaseWidget
	checked   bool
	onChanged func(bool)
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
	rect := canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 0})
	rect.StrokeColor = border
	rect.StrokeWidth = 2
	rect.CornerRadius = 3
	rect.SetMinSize(fyne.NewSize(20, 20))

	tick := canvas.NewText("", color.White)
	if c.checked {
		rect.FillColor = fillChecked
		tick.Text = "✓"
	}
	cont := container.NewMax(rect, container.NewCenter(tick))
	return widget.NewSimpleRenderer(cont)
}

func (c *squareCheckbox) MinSize() fyne.Size { return fyne.NewSize(20, 20) }

func (c *squareCheckbox) Tapped(*fyne.PointEvent) {
	c.checked = !c.checked
	if c.onChanged != nil {
		c.onChanged(c.checked)
	}
	c.Refresh()
}

// Status indicator (✓ or ★), star toggles on tap when shown
type statusIndicator struct {
	widget.BaseWidget
	todo     *models.TodoItem
	onToggle func(toggleStar bool)
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
		txt = "✓"
		if _, ok := fyne.CurrentApp().Settings().Theme().(*LightSoftTheme); ok {
			col = color.NRGBA{R: 0x4C, G: 0xAF, B: 0x50, A: 0xFF}
		} else {
			col = color.NRGBA{R: 0x98, G: 0x97, B: 0x1a, A: 0xFF}
		}
	} else if s.todo.Starred {
		txt = "★"
		if _, ok := fyne.CurrentApp().Settings().Theme().(*LightSoftTheme); ok {
			col = color.NRGBA{R: 0xFF, G: 0x8C, B: 0x42, A: 0xFF}
		} else {
			col = color.NRGBA{R: 0xFE, G: 0x80, B: 0x19, A: 0xFF}
		}
	} else {
		txt = ""
		col = theme.Color(theme.ColorNameForeground)
	}
	t := canvas.NewText(txt, col)
	t.TextSize = 20
	return widget.NewSimpleRenderer(container.NewCenter(t))
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

// ScrollToTop scrolls to the top of the timeline
func (t *Timeline) ScrollToTop() {
	if renderer := t.CreateRenderer(); renderer != nil {
		if scroll, ok := renderer.(*timelineRenderer); ok {
			scroll.scroll.ScrollToTop()
		}
	}
}
