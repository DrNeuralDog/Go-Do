package ui

import (
	"fmt"
	"image/color"
	"time"

	"godo/src/models"
	"godo/src/ui/helpers"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// createTodoItem builds one todo row with controls and drag layers
func (r *timelineRenderer) createTodoItem(todo *models.TodoItem) fyne.CanvasObject {
	isLightTheme := helpers.IsLightTheme()

	colorSquare := canvas.NewRectangle(todo.GetLevelColor())
	colorSquare.CornerRadius = 2

	colorSquare.SetMinSize(fyne.NewSize(16, 32))

	colorSquareWrap := container.NewGridWrap(fyne.NewSize(16, 32), colorSquare)
	colorSquareAligned := verticallyCenterCompact(colorSquareWrap)

	doneCheck := newSquareCheckbox(todo.Done, func(checked bool) {
		updated := *todo
		updated.Done = checked

		if err := r.timeline.dataManager.UpdateTodo(&updated, todo.TodoTime); err != nil {
			r.timeline.showError(err)

			return
		}

		r.timeline.notifyTodosChanged()
	})

	doneCheckCentered := verticallyCenterCompact(doneCheck)

	nameLabel := widget.NewLabel(todo.Name)
	nameLabel.Wrapping = fyne.TextWrapWord

	var timeColor color.Color

	if isLightTheme {
		timeColor = color.NRGBA{R: 0x66, G: 0x66, B: 0x66, A: 0xFF}
	} else {
		timeColor = color.NRGBA{R: 0xA8, G: 0x99, B: 0x84, A: 0xFF}
	}

	timeText := canvas.NewText(fmt.Sprintf("%02d:%02d", todo.TodoTime.Hour(), todo.TodoTime.Minute()), timeColor)
	timeText.TextSize = 18
	timeLabel := verticallyCenterCompact(timeText)

	status := newStatusIndicator(todo, func(toggleStar bool) {
		if !toggleStar {
			return
		}

		updated := *todo
		updated.Starred = !todo.Starred

		if err := r.timeline.dataManager.UpdateTodo(&updated, todo.TodoTime); err != nil {
			r.timeline.showError(err)

			return
		}

		r.timeline.notifyTodosChanged()
	})

	statusCentered := verticallyCenterCompact(status)

	deleteBtn := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		r.timeline.confirmDelete(todo)
	})
	deleteBtn.Importance = widget.LowImportance
	deleteBtnCentered := verticallyCenterCompact(deleteBtn)

	leftSection := container.NewHBox(colorSquareAligned, helpers.CreateSpacer(8, 1), doneCheckCentered)
	rightSection := container.NewHBox(timeLabel, helpers.CreateSpacer(8, 1), statusCentered, helpers.CreateSpacer(4, 1), deleteBtnCentered)
	content := container.NewBorder(nil, nil, leftSection, rightSection, verticallyCenterWide(nameLabel))

	var borderColor color.Color

	if isLightTheme {
		borderColor = color.NRGBA{R: 0xD0, G: 0xD0, B: 0xD0, A: 0xFF}
	} else {
		borderColor = color.NRGBA{R: 0x3c, G: 0x38, B: 0x36, A: 0xFF}
	}

	bottomLine := canvas.NewRectangle(borderColor)
	bottomLine.SetMinSize(fyne.NewSize(10, 1))

	row := container.NewVBox(
		container.NewPadded(content),
		bottomLine,
	)

	shadowRect := canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 0})
	shadowRect.CornerRadius = 4

	borderRect := canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 0})
	borderRect.CornerRadius = 4

	pressOverlay := canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 0})
	pressOverlay.CornerRadius = 4

	pressStack := container.NewMax(shadowRect, borderRect, row, pressOverlay)
	isDragging := r.timeline.draggingTodo == todo

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

	if isDragging {
		tappable.applyDragVisuals(80)
	}

	tappable.ExtendBaseWidget(tappable)

	return tappable
}

// verticallyCenterCompact centers fixed-size controls
func verticallyCenterCompact(obj fyne.CanvasObject) fyne.CanvasObject {
	return container.NewVBox(layout.NewSpacer(), container.NewCenter(obj), layout.NewSpacer())
}

// verticallyCenterWide centers stretchable row content
func verticallyCenterWide(obj fyne.CanvasObject) fyne.CanvasObject {
	return container.NewVBox(layout.NewSpacer(), obj, layout.NewSpacer())
}

// tappableTodo handles row selection and drag reorder
type tappableTodo struct {
	widget.BaseWidget

	todo         *models.TodoItem
	todoTime     time.Time
	container    *fyne.Container
	press        *canvas.Rectangle
	shadow       *canvas.Rectangle
	border       *canvas.Rectangle
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

func (tt *tappableTodo) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

// Tapped flashes row and opens todo editor
func (tt *tappableTodo) Tapped(*fyne.PointEvent) {
	if tt.press != nil {
		col := helpers.ToNRGBA(theme.Color(theme.ColorNameHover))

		if col.A < 60 {
			col.A = 60
		}

		tt.press.FillColor = col

		runOnMainThread(func() {
			tt.press.Refresh()
		})

		go func(press *canvas.Rectangle) {
			time.Sleep(120 * time.Millisecond)

			runOnMainThread(func() {
				press.FillColor = color.NRGBA{R: 0, G: 0, B: 0, A: 0}

				press.Refresh()
			})
		}(tt.press)
	}

	if tt.onSelected != nil {
		tt.onSelected(tt.todo, tt.todoTime)
	}
}

// Dragged handles touch drag reorder
func (tt *tappableTodo) Dragged(e *fyne.DragEvent) {
	if !tt.dragging {
		tt.dragging = true
	}

	tt.applyDragVisuals(80)
	tt.processDrag(e.Dragged.DY)
}

// DragEnd finishes touch drag
func (tt *tappableTodo) DragEnd() {
	tt.finishDrag()
}

// MouseDown starts desktop drag tracking
func (tt *tappableTodo) MouseDown(e *desktop.MouseEvent) {
	tt.dragging = true
	tt.lastMouseY = e.Position.Y

	tt.applyDragVisuals(50)
}

// MouseMoved converts mouse movement to reorder steps
func (tt *tappableTodo) MouseMoved(e *desktop.MouseEvent) {
	if !tt.dragging {
		return
	}

	dy := e.Position.Y - tt.lastMouseY
	tt.lastMouseY = e.Position.Y

	tt.processDrag(dy)
}

// MouseUp finishes desktop drag
func (tt *tappableTodo) MouseUp(*desktop.MouseEvent) {
	tt.finishDrag()
}

// applyDragVisuals highlights active drag row
func (tt *tappableTodo) applyDragVisuals(pressAlpha uint8) {
	runOnMainThread(func() {
		if tt.press != nil {
			tt.press.FillColor = color.NRGBA{R: 0x3C, G: 0x82, B: 0xFF, A: pressAlpha}

			tt.press.Refresh()
		}

		if tt.border != nil {
			tt.border.FillColor = color.Transparent
			tt.border.StrokeColor = color.NRGBA{R: 0x3C, G: 0x82, B: 0xFF, A: 200}
			tt.border.StrokeWidth = 2

			tt.border.Refresh()
		}

		if tt.shadow != nil {
			tt.shadow.FillColor = color.NRGBA{R: 0, G: 0, B: 0, A: 30}

			tt.shadow.Refresh()
		}
	})

	if tt.timeline != nil {
		tt.timeline.draggingTodo = tt.todo
	}
}

// clearDragVisuals resets drag highlight layers
func (tt *tappableTodo) clearDragVisuals() {
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
}

// processDrag converts accumulated Y movement to reorder callbacks
func (tt *tappableTodo) processDrag(deltaY float32) {
	tt.dragAccumY += deltaY

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

// finishDrag resets drag state and refreshes timeline
func (tt *tappableTodo) finishDrag() {
	tt.dragging = false
	tt.dragAccumY = 0

	tt.clearDragVisuals()

	if tt.onReorderEnd != nil {
		tt.onReorderEnd()
	}

	if tt.timeline != nil {
		tt.timeline.draggingTodo = nil
		tt.timeline.Refresh()
	}
}
