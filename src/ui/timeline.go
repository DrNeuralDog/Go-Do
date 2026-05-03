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
	"fyne.io/fyne/v2/widget"
)

// Timeline shows todos for the selected day
type Timeline struct {
	widget.BaseWidget

	dataManager persistence.TodoRepository
	currentDate time.Time
	todos       []*models.TodoItem
	viewMode    models.ViewMode
	window      fyne.Window
	renderer    *timelineRenderer

	onTodoSelected    func(*models.TodoItem, time.Time)
	onTodoReorder     func(*models.TodoItem, int)
	onReorderFinished func()
	onTodosChanged    func()

	draggingTodo *models.TodoItem
}

// NewTimeline creates timeline widget
func NewTimeline(dataManager persistence.TodoRepository) *Timeline {
	t := &Timeline{
		dataManager: dataManager,
		currentDate: time.Now(),
		viewMode:    models.ViewAll,
	}

	t.ExtendBaseWidget(t)

	return t
}

// SetDate changes selected day
func (t *Timeline) SetDate(date time.Time) {
	t.currentDate = date
}

// SetViewMode changes current filter mode
func (t *Timeline) SetViewMode(mode models.ViewMode) {
	t.viewMode = mode
}

// SetTodos stores visible todos for rendering
func (t *Timeline) SetTodos(todos []*models.TodoItem) {
	t.todos = todos
}

// SetWindow stores parent window for dialogs
func (t *Timeline) SetWindow(win fyne.Window) {
	t.window = win
}

// SetOnTodoSelected registers item selection callback
func (t *Timeline) SetOnTodoSelected(callback func(*models.TodoItem, time.Time)) {
	t.onTodoSelected = callback
}

// SetOnTodoReorder registers reorder callback
func (t *Timeline) SetOnTodoReorder(callback func(*models.TodoItem, int)) {
	t.onTodoReorder = callback
}

// SetOnReorderFinished registers drag completion callback
func (t *Timeline) SetOnReorderFinished(callback func()) {
	t.onReorderFinished = callback
}

// SetOnTodosChanged registers data change callback
func (t *Timeline) SetOnTodosChanged(callback func()) {
	t.onTodosChanged = callback
}

// CreateRenderer creates timeline renderer
func (t *Timeline) CreateRenderer() fyne.WidgetRenderer {
	renderer := &timelineRenderer{timeline: t}
	t.renderer = renderer

	return renderer
}

// ScrollToTop moves current timeline renderer to top
func (t *Timeline) ScrollToTop() {
	renderer := t.renderer

	if renderer == nil {
		return
	}

	runOnMainThread(func() {
		if renderer.scroll == nil {
			return
		}

		renderer.scroll.ScrollToTop()
	})
}

// timelineRenderer renders timeline content
type timelineRenderer struct {
	timeline *Timeline
	scroll   *container.Scroll
	listBox  *fyne.Container
}

// ensureScroll creates persistent scroll objects once
func (r *timelineRenderer) ensureScroll() {
	if r.listBox == nil {
		r.listBox = r.createTimelineContent()
	}

	if r.scroll == nil {
		r.scroll = container.NewScroll(r.listBox)
	}
}

func (r *timelineRenderer) Layout(size fyne.Size) {
	r.ensureScroll()
	r.scroll.Resize(size)
}

func (r *timelineRenderer) MinSize() fyne.Size {
	return fyne.NewSize(350, 200)
}

// Refresh rebuilds rows and keeps current scroll offset
func (r *timelineRenderer) Refresh() {
	r.ensureScroll()

	offset := r.scroll.Offset
	r.listBox.Objects = r.buildTimelineObjects()

	runOnMainThread(func() {
		r.listBox.Refresh()
	})

	r.scroll.Offset = offset

	if size := r.timeline.Size(); size.Width > 0 && size.Height > 0 {
		r.scroll.Resize(size)
	}
}

// BackgroundColor keeps parent task container visible
func (r *timelineRenderer) BackgroundColor() fyne.ThemeColorName {
	return ""
}

func (r *timelineRenderer) Objects() []fyne.CanvasObject {
	r.ensureScroll()

	return []fyne.CanvasObject{r.scroll}
}

// Destroy clears cached renderer reference
func (r *timelineRenderer) Destroy() {
	if r.timeline != nil && r.timeline.renderer == r {
		r.timeline.renderer = nil
	}
}

// createTimelineContent wraps rendered rows into scroll content
func (r *timelineRenderer) createTimelineContent() *fyne.Container {
	objects := r.buildTimelineObjects()

	return container.NewVBox(objects...)
}

// buildTimelineObjects builds date header, empty state, and todo rows
func (r *timelineRenderer) buildTimelineObjects() []fyne.CanvasObject {
	var objects []fyne.CanvasObject

	dateKey := r.timeline.currentDate.Format("2006-01-02")
	dateHeader := r.createDateHeader(dateKey)
	objects = append(objects, dateHeader)

	if len(r.timeline.todos) == 0 {
		emptyLabel := widget.NewLabel(localization.GetString("status_empty_list"))
		emptyLabel.Alignment = fyne.TextAlignCenter
		objects = append(objects, emptyLabel)

		return objects
	}

	for _, todo := range r.timeline.todos {
		todoItem := r.createTodoItem(todo)
		objects = append(objects, todoItem)
	}

	return objects
}

// createDateHeader builds selected day title and divider
func (r *timelineRenderer) createDateHeader(dateKey string) fyne.CanvasObject {
	date, err := time.Parse("2006-01-02", dateKey)

	if err != nil {
		return widget.NewLabel("Invalid Date")
	}

	isLightTheme := helpers.IsLightTheme()
	weekdayNames := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	weekdayName := weekdayNames[date.Weekday()]

	headerText := fmt.Sprintf("%d/%02d/%02d %s",
		date.Year(), date.Month(), date.Day(), weekdayName)

	var fg color.Color

	if isLightTheme {
		fg = color.NRGBA{R: 0x3c, G: 0x38, B: 0x36, A: 0xFF}
	} else {
		fg = color.NRGBA{R: 0xeb, G: 0xdb, B: 0xb2, A: 0xFF}
	}

	headerLabel := canvas.NewText(headerText, fg)
	headerLabel.TextStyle = fyne.TextStyle{Bold: true}
	headerLabel.TextSize = 20
	headerLabel.Alignment = fyne.TextAlignCenter

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
