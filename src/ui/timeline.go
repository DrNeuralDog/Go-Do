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
	fmt.Printf("[Timeline.SetTodos] Called with %d todos\n", len(todos))
	t.todos = todos
	t.organizeByDate()
	fmt.Printf("[Timeline.SetTodos] After organize: %d visible items, %d date groups\n", len(t.visibleItems), len(t.dateGroups))
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
	fmt.Printf("[Timeline.Layout] Called with size: %v\n", size)
	if r.scroll == nil {
		// Initialize scroll if not already done
		content := r.createTimelineContent()
		r.scroll = container.NewScroll(content)
		fmt.Printf("[Timeline.Layout] Created new scroll container\n")
	}
	// Ensure scroll fills all available area
	r.scroll.Resize(size)
	fmt.Printf("[Timeline.Layout] Resized scroll to: %v\n", size)
}

func (r *timelineRenderer) MinSize() fyne.Size {
	// Always return a reasonable minimum size to prevent UI shrinking
	// The scroll container's MinSize of {32 32} is too small and causes layout issues
	minSize := fyne.NewSize(350, 200)
	fmt.Printf("[Timeline.MinSize] Returning: %v (scroll exists: %v)\n", minSize, r.scroll != nil)
	return minSize
}

func (r *timelineRenderer) Refresh() {
	fmt.Printf("[Timeline.Refresh] Starting refresh, timeline size: %v\n", r.timeline.Size())

	// Store current scroll position
	var offset fyne.Position
	if r.scroll != nil {
		offset = r.scroll.Offset
		fmt.Printf("[Timeline.Refresh] Storing scroll offset: %v\n", offset)
	}

	// Recreate scroll container entirely to avoid layout bugs
	content := r.createTimelineContent()
	fmt.Printf("[Timeline.Refresh] Created new content with %d visible items\n", len(r.timeline.visibleItems))

	r.scroll = container.NewScroll(content)
	fmt.Printf("[Timeline.Refresh] Created new scroll container, MinSize: %v\n", r.scroll.MinSize())

	// Restore scroll position
	r.scroll.Offset = offset

	// Ensure proper sizing
	if size := r.timeline.Size(); size.Width > 0 && size.Height > 0 {
		r.scroll.Resize(size)
		fmt.Printf("[Timeline.Refresh] Resized scroll to: %v\n", size)
	} else {
		fmt.Printf("[Timeline.Refresh] WARNING: Timeline has invalid size: %v\n", size)
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

	headerText := fmt.Sprintf("%d/%02d/%02d %s   +",
		date.Year(), date.Month(), date.Day(), weekdayName)

	headerLabel := widget.NewLabel(headerText)
	headerLabel.TextStyle = fyne.TextStyle{Bold: true}

	// Add separator line
	separator := canvas.NewLine(color.RGBA{R: 200, G: 200, B: 200, A: 255})
	separator.StrokeWidth = 1

	return container.NewVBox(
		container.NewHBox(headerLabel, widget.NewSeparator()),
		separator,
	)
}

func (r *timelineRenderer) createTodoItem(todo *models.TodoItem) fyne.CanvasObject {
	// Priority indicator using colored background rectangle
	priorityRect := canvas.NewRectangle(todo.GetLevelColor())
	priorityRect.SetMinSize(fyne.NewSize(20, 20))

	// Todo name
	nameLabel := widget.NewLabel(todo.Name)
	nameLabel.Wrapping = fyne.TextWrapWord

	// Time display
	timeLabel := widget.NewLabel(fmt.Sprintf("%02d:%02d", todo.TodoTime.Hour(), todo.TodoTime.Minute()))
	timeLabel.Alignment = fyne.TextAlignTrailing

	// Real completion checkbox (avoid firing OnChanged during initial SetChecked)
	doneCheck := widget.NewCheck("", nil)
	doneCheck.SetChecked(todo.Done)
	doneCheck.OnChanged = func(checked bool) {
		// Update persistence
		updated := *todo
		updated.Done = checked
		_ = r.timeline.dataManager.UpdateTodo(&updated, todo.TodoTime)
		// Reload month and re-apply filter
		if todos, err := r.timeline.dataManager.GetTodosForMonth(r.timeline.currentDate.Year, r.timeline.currentDate.Month); err == nil {
			filtered := r.timeline.viewMode.FilterItems(todos, time.Now())
			r.timeline.SetTodos(filtered)
		}
	}

	// Delete button
	delButton := widget.NewButtonWithIcon("", RedCrossIcon, func() {
		_ = r.timeline.dataManager.RemoveTodo(todo.TodoTime)
		if todos, err := r.timeline.dataManager.GetTodosForMonth(r.timeline.currentDate.Year, r.timeline.currentDate.Month); err == nil {
			filtered := r.timeline.viewMode.FilterItems(todos, time.Now())
			r.timeline.SetTodos(filtered)
		}
	})
	delButton.Importance = widget.LowImportance
	// Make delete icon smaller and vertically centered
	delWrap := container.NewGridWrap(fyne.NewSize(16, 16), delButton)

	// Star button
	// Use built-in icons: Filled vs Confirm
	// Choose star icon: outline when off, filled tinted primary when on
	var starIcon fyne.Resource
	if todo.Starred {
		// Use blue filled star explicitly
		starIcon = StarBlueIcon
	} else {
		starIcon = StarOutlineIcon
	}
	starButton := widget.NewButtonWithIcon("", starIcon, func() {
		updated := *todo
		updated.Starred = !todo.Starred
		_ = r.timeline.dataManager.UpdateTodo(&updated, todo.TodoTime)
		if todos, err := r.timeline.dataManager.GetTodosForMonth(r.timeline.currentDate.Year, r.timeline.currentDate.Month); err == nil {
			filtered := r.timeline.viewMode.FilterItems(todos, time.Now())
			r.timeline.SetTodos(filtered)
		}
	})
	starButton.Importance = widget.LowImportance
	starWrap := container.NewGridWrap(fyne.NewSize(20, 20), starButton)

	// Layout the todo item
	leftSection := container.NewHBox(priorityRect, doneCheck)
	// Use Border layout so the name takes the center and time is pinned right
	rightSection := container.NewHBox(container.NewCenter(timeLabel), container.NewCenter(starWrap), container.NewCenter(delWrap))
	content := container.NewBorder(nil, nil, leftSection, rightSection, nameLabel)
	// Add soft background behind each item
	itemBg := canvas.NewRectangle(theme.Color(theme.ColorNameInputBackground))
	item := container.NewMax(itemBg, container.NewPadded(content))

	// Make the entire item clickable
	tappable := &tappableTodo{
		todo:       todo,
		todoTime:   todo.TodoTime,
		container:  item,
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
