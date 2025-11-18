# CR: Pomodoro Module Improvements
**Version:** 1.0
**Status:** In Progress
**Priority:** High
**Author:** Development Team
**Date:** 2025-11-17

---

## Executive Summary

This Change Request details a comprehensive overhaul of the Pomodoro timer module to fix multiple UI/UX issues. The module currently has a window that is too large, incorrect progress ring color gradients, spinner value display synchronization problems, button state inconsistencies, and missing animations and styling improvements on spinner controls.

---

## Issues Summary

| # | Issue | Current Behavior | Desired Behavior | Priority | Complexity |
|---|-------|------------------|------------------|----------|-----------|
| 1 | Window size | 400x600px (too large) | Reduce to main window size (~320-350px width) | High | Low |
| 2 | Progress ring colors | Green→Red gradient (wrong direction) | Yellow-Orange→Green gradient (left to right) | High | Medium |
| 3 | Spinner value sync | Arrow clicks update circle display but not spinner display | Both circle and spinner field update together | High | High |
| 4 | Pause/Resume button behavior | Both Start and Pause change to "Resume" | Only Pause changes to "Resume", Start becomes disabled (gray) | High | Medium |
| 5 | Spinner arrow animations | No visual feedback on click | Add click press animation with visual feedback | Medium | Medium |
| 6 | Spinner arrow hover state | No hover feedback | Add hover animation that disappears on mouse out | Medium | Medium |
| 7 | Spinner arrow vertical alignment | Arrows cramped at bottom (large top padding) | Center arrows vertically in spinner | Medium | Low |
| 8 | Configuration label position | "Configuration" text left-aligned | Center-align "Configuration" text | Low | Low |

---

## Detailed Requirements

### Issue 1: Window Size Reduction

**Current State:**
- Window size: 400x600px (line 66 in pomodoro_window.go)
- Too large relative to main application window

**Requirements:**
- Reduce window width to ~320-350px
- Adjust height proportionally to maintain layout
- Ensure all UI elements remain properly positioned and readable
- Ring display should still be clearly visible and appropriately sized
- Configuration section should remain fully visible without excessive scrolling

**Implementation Notes:**
- Update `pw.window.Resize()` call in `setupUI()`
- May need to adjust ring dimensions in `ringWithDigits` container
- Retest all UI elements for proper alignment and spacing
- Ensure no content is hidden or truncated

---

### Issue 2: Progress Ring Color Gradient

**Current State:**
```go
StartColor:  hex("#b8bb26"), // Gruvbox green - start
EndColor:    hex("#fb4934"), // Gruvbox red - end
```
- Gradient progresses from green (start) to red (end)
- Direction: left-to-right, but color semantics are backwards

**Requirements:**
- Change gradient to Yellow-Orange → Green
- Yellow-Orange at the beginning (0% progress)
- Green at the end (100% progress)
- Semantics: "danger" (yellow) → "success" (green)
- Should feel natural: warming up → ready/complete

**Color Palette Options:**
```
Option A (Gruvbox-based):
- Start (0%):   #fabd2f (Yellow)      [Gruvbox yellow]
- Mid (~50%):   #fe8019 (Orange)      [Gruvbox orange]
- End (100%):   #b8bb26 (Green)       [Gruvbox green]

Option B (Alternative):
- Start (0%):   #ffb347 (Light Orange)
- End (100%):   #98c379 (Light Green)
```

**Implementation Details:**
- Update `NewProgressRing()` constructor
- Reverse or remap color interpolation in `progressRingRenderer.Refresh()` (lines 409-435)
- Test visual appearance at 0%, 25%, 50%, 75%, 100% progress
- Ensure colors are distinguishable on both light and dark themes

---

### Issue 3: Spinner Value Display Synchronization

**Current Problem:**
- When arrow buttons (up/down) are clicked, `NumberSpinner.increment()`/`decrement()` is called
- These methods call `SetValue()` which updates internal `Value` field
- `SetValue()` calls `Refresh()` to update the display
- However, the text displayed in the spinner widget is not being updated
- The progress ring circle updates correctly (via `pw.tick()` callback)
- The spinner's own text field shows stale value

**Root Cause Analysis:**
- `NumberSpinner.CreateRenderer()` (lines 711-739) creates renderer with static text
- When `Refresh()` is called, it doesn't recreate the renderer
- The text canvas object is not being updated with new value

**Solution Approach:**

**Option A: Store Renderer Reference (Complex)**
- Keep reference to renderer in NumberSpinner
- On Refresh(), directly update text in stored renderer
- Risk: Violates Fyne architecture patterns

**Option B: Implement Custom Widget Properly (Recommended)**
- Override `Refresh()` method in NumberSpinner
- Make text update part of refresh cycle
- Store text reference or use canvas refresh

**Implementation Steps:**
1. Add text field reference to NumberSpinner struct:
   ```go
   type NumberSpinner struct {
       // ... existing fields ...
       displayText *canvas.Text  // NEW
       renderer    fyne.WidgetRenderer // NEW (optional)
   }
   ```

2. Create/update a custom `numberSpinnerRenderer` struct:
   ```go
   type numberSpinnerRenderer struct {
       spinner *NumberSpinner
       bg      *canvas.Rectangle
       txt     *canvas.Text
       // ... other fields ...
   }
   ```

3. Implement proper refresh chain:
   - In NumberSpinner.SetValue(): call refresh
   - In numberSpinnerRenderer.Refresh(): update `txt.Text = fmt.Sprintf("%d", ns.Value)`
   - Store renderer reference for direct access

4. Alternative quick fix (if above too invasive):
   - Override CreateRenderer() to return custom renderer that stores text reference
   - Make text update happen in a separate refresh handler
   - Use widget.Refresh() pattern properly

**Testing:**
- Click up arrow → value in spinner text should change immediately
- Click down arrow → value in spinner text should change immediately
- Value in progress ring should also update (via callback)
- Configuration changes should persist

---

### Issue 4: Pause/Resume Button Behavior

**Current State (lines 293-310):**
```go
case models.PomodoroPaused:
    pw.startBtn.Enable()
    pw.pauseBtn.Enable()
    pw.startBtn.SetText("Resume")
    pw.pauseBtn.SetText("Resume")  // ← PROBLEM
```

**Problem:**
- Both buttons show "Resume" when paused
- User can click either button to resume (confusing UX)
- No clear indication which button is "active"

**Requirements:**
- When paused:
  - **Pause button** → "Resume" (enabled, can be clicked)
  - **Start button** → Disabled (grayed out, cannot be clicked)
  - The Pause button becomes the primary action

- When running:
  - **Start button** → Disabled (grayed out, cannot be clicked)
  - **Pause button** → "Pause" (enabled, can be clicked)

- When idle (reset):
  - **Start button** → "Start" (enabled, can be clicked)
  - **Pause button** → Disabled (grayed out)

**Implementation (lines 293-310):**
```go
case models.PomodoroPaused:
    pw.startBtn.Disable()        // ← START IS DISABLED
    pw.pauseBtn.Enable()
    pw.pauseBtn.SetText("Resume") // ← ONLY PAUSE CHANGES
    pw.resetBtn.Enable()
```

**Event Handler Update (lines 454-461):**
- `onPauseClicked()` should handle both Pause and Resume since button shows appropriate text
- `onStartClicked()` should only handle Start (Resume handled by Pause button when paused)
- Simplify logic:
  ```go
  func (pw *PomodoroWindow) onPauseClicked() {
      if pw.pauseBtn.Text == "Resume" {  // or check timer state
          pw.timer.Resume()
      } else {
          pw.timer.Pause()
      }
      pw.tick()
  }
  ```

---

### Issue 5: Spinner Arrow Click Animation

**Current State:**
- Arrow buttons (TinyIconButton from line 727-728) have no visual feedback
- Clicking arrows provides no tactile/visual confirmation
- Users unsure if click registered

**Requirements:**
- Add press animation when arrow is clicked
- Animation duration: 120ms (consistent with button style)
- Visual feedback: Brief color change or scale change
- Must not interfere with spinner value update

**Implementation:**
- Enhance `TinyIconButton` or create `AnimatedTinyIconButton`
- Add press handler similar to SimpleRectButton (lines 355-371):
  ```go
  func (b *TinyIconButton) Tapped(*fyne.PointEvent) {
      // Store original state
      originalStyle := b.getCurrentStyle()

      // Apply pressed state (darken/scale)
      b.setPressed(true)
      b.Refresh()

      // Execute callback
      if b.OnTapped != nil {
          b.OnTapped()
      }

      // Animate back to normal after 120ms
      go func() {
          time.Sleep(120 * time.Millisecond)
          b.setPressed(false)
          b.Refresh()
      }()
  }
  ```

**Options:**
1. Darken icon color (like SimpleRectButton)
2. Scale icon down slightly (0.85x) and back
3. Combination of both

**Testing:**
- Click up arrow multiple times → see brief visual feedback each time
- Click down arrow multiple times → see brief visual feedback each time
- Value updates correctly while animation plays

---

### Issue 6: Spinner Arrow Hover Animation

**Current State:**
- Arrow buttons have no hover feedback
- No visual indication that they're interactive

**Requirements:**
- Show hover effect when mouse enters arrow button area
- Effect should smoothly appear on hover
- Effect should disappear when mouse leaves
- Should work on both light and dark themes
- Must not interfere with click animation

**Implementation:**
- Add hover state tracking to arrow buttons
- Implement `MouseIn()`, `MouseMoved()`, `MouseOut()` handlers (like SimpleRectButton lines 395-407)
- Visual feedback options:
  1. Brighten/lighten icon color
  2. Add subtle background circle behind icon
  3. Scale icon up slightly (1.1x)

**Recommended approach:**
- Add `hovered bool` field to TinyIconButton or wrapper
- In Refresh(), apply color change if hovered:
  ```go
  if b.hovered {
      iconColor = lighten(baseColor, 0.15)
  }
  ```
- Quick implementation: Use existing icon color and brighten on hover

**Testing:**
- Hover over up arrow → icon brightens/becomes visible
- Move mouse away → icon returns to normal
- Hover over down arrow → same behavior
- Works on light and dark themes

---

### Issue 7: Spinner Arrow Vertical Alignment

**Current State (lines 729-732):**
```go
upDown := container.NewVBox(
    container.NewGridWrap(fyne.NewSize(24, 24), upBtn),
    container.NewGridWrap(fyne.NewSize(24, 24), downBtn),
)
```
- Arrows are cramped at bottom with large top padding
- Appears misaligned relative to number display
- Visual hierarchy is off

**Root Cause:**
- TinyIconButton.MinSize() returns 16x16 (line 619)
- Wrapped in 24x24 GridWrap containers
- Vertical spacing not optimized

**Requirements:**
- Center arrows vertically within the spinner height (48px)
- Balance space above and below arrow group
- Arrows should be vertically centered with number text
- Total spinner height should remain 48px

**Implementation:**
```go
// Option A: Use centered container with spacing
upDown := container.NewVBox(
    container.NewCenter(upBtn),      // 24x24 with auto-centering
    container.NewCenter(downBtn),    // 24x24 with auto-centering
)

// Option B: Use GridWrap with auto-sizing
upDown := container.NewGridWrap(
    fyne.NewSize(24, 24),
    container.NewVBox(
        container.NewCenter(upBtn),
        container.NewCenter(downBtn),
    ),
)

// Option C: Custom padding/margin
upDown := container.NewVBox(
    CreateSpacer(24, 2),              // Top padding
    upBtn,
    downBtn,
    CreateSpacer(24, 2),              // Bottom padding
)
```

**Testing:**
- Arrows should be horizontally centered in their space
- Arrows should be vertically centered in the 48px height
- Text and arrows appear balanced
- Works at different window sizes

---

### Issue 8: Configuration Label Center Alignment

**Current State (line 119):**
```go
cfgHeader := canvas.NewText("Configuration", titleColor)
cfgHeader.TextStyle = fyne.TextStyle{Bold: true}
cfgHeader.TextSize = 18
```
- Added directly to VBox (line 147)
- Left-aligned by default
- Not centered in the UI

**Requirements:**
- "Configuration" label should be center-aligned
- Visual hierarchy should be clear (prominent header)
- Should align with other centered elements in window

**Implementation:**
- Wrap cfgHeader in a centering container:
```go
cfgHeader := canvas.NewText("Configuration", titleColor)
cfgHeader.TextStyle = fyne.TextStyle{Bold: true}
cfgHeader.TextSize = 18

cfgHeaderCenter := container.NewCenter(cfgHeader)

configForm := container.NewVBox(
    cfgHeaderCenter,  // ← CENTERED
    widget.NewSeparator(),
    // ... rest of form ...
)
```

**Testing:**
- "Configuration" text is centered horizontally
- Alignment is consistent with timer display above
- Text is readable and properly styled

---

## Implementation Plan

### Phase 1: Essential Fixes (Week 1)
**Priority: Critical fixes that impact usability**

1. **Issue 4 - Button State Logic** (Issue Priority: HIGH, Complexity: MEDIUM)
   - Estimated time: 1-2 hours
   - Impact: Fixes confusing UX immediately
   - Files to modify: `pomodoro_window.go`
   - Changes: Update tick() method lines 293-310

2. **Issue 1 - Window Size** (Issue Priority: HIGH, Complexity: LOW)
   - Estimated time: 1-2 hours
   - Impact: Improves overall aesthetics
   - Files to modify: `pomodoro_window.go`
   - Changes: Line 66, potentially adjust container sizes

3. **Issue 8 - Configuration Label** (Issue Priority: LOW, Complexity: LOW)
   - Estimated time: 30 min
   - Impact: Minor visual improvement
   - Files to modify: `pomodoro_window.go`
   - Changes: Lines 119-147

### Phase 2: Core Functionality (Week 1-2)
**Priority: Major fixes affecting core functionality**

1. **Issue 3 - Spinner Value Sync** (Issue Priority: HIGH, Complexity: HIGH)
   - Estimated time: 3-4 hours
   - Impact: Fixes critical display bug
   - Files to modify: `style_helpers.go` (NumberSpinner), `pomodoro_window.go`
   - Refactoring needed: Implement custom renderer properly
   - Testing: Extensive UI testing

2. **Issue 2 - Progress Ring Colors** (Issue Priority: HIGH, Complexity: MEDIUM)
   - Estimated time: 1-2 hours
   - Impact: Better UX feedback
   - Files to modify: `pomodoro_window.go` (ProgressRing)
   - Changes: Lines 327-340, 409-435

### Phase 3: Polish & Animation (Week 2)
**Priority: Enhancement fixes for improved UX**

1. **Issue 5 - Click Animation** (Issue Priority: MEDIUM, Complexity: MEDIUM)
   - Estimated time: 1-2 hours
   - Impact: Better tactile feedback
   - Files to modify: `style_helpers.go` (TinyIconButton)
   - Changes: Add Tapped() method, press state handling

2. **Issue 6 - Hover Animation** (Issue Priority: MEDIUM, Complexity: MEDIUM)
   - Estimated time: 1-2 hours
   - Impact: Better visual feedback
   - Files to modify: `style_helpers.go` (TinyIconButton)
   - Changes: Add MouseIn/MouseOut handlers

3. **Issue 7 - Arrow Alignment** (Issue Priority: MEDIUM, Complexity: LOW)
   - Estimated time: 30 min - 1 hour
   - Impact: Visual polish
   - Files to modify: `pomodoro_window.go`
   - Changes: Lines 729-732 container setup

---

## Technical Approach

### NumberSpinner Redesign (Issue 3)

**Current Architecture Problem:**
- NumberSpinner inherits from BaseWidget
- CreateRenderer() creates anonymous renderer (via SimpleRenderer)
- No persistent reference to renderer components
- Refresh() doesn't reach text element

**Proposed Solution:**

Create a dedicated renderer type:
```go
type numberSpinnerRenderer struct {
    spinner *NumberSpinner
    bg      *canvas.Rectangle
    txt     *canvas.Text
    upBtn   *TinyIconButton
    downBtn *TinyIconButton
    cont    *fyne.Container
}

func (r *numberSpinnerRenderer) Refresh() {
    // Update text display
    r.txt.Text = fmt.Sprintf("%d", r.spinner.Value)
    r.txt.Refresh()
    r.cont.Refresh()
}

func (ns *NumberSpinner) CreateRenderer() fyne.WidgetRenderer {
    // Create all components including text
    txt := canvas.NewText(fmt.Sprintf("%d", ns.Value), ...)
    // Store text in renderer so Refresh can access it
    renderer := &numberSpinnerRenderer{
        spinner: ns,
        txt:    txt,
        // ... other components
    }
    return renderer
}
```

---

## Testing Strategy

### Unit Tests

1. **Spinner Value Sync:**
   - Test increment() updates Value field
   - Test increment() calls OnChanged callback
   - Test decrement() updates Value field
   - Test SetValue() clamps to min/max

2. **Timer Button States:**
   - Test idle state enables Start, disables Pause/Reset
   - Test running state disables Start, enables Pause/Reset
   - Test paused state disables Start, enables Resume/Reset

3. **Progress Ring:**
   - Test gradient colors at different progress values
   - Test color interpolation logic

### Integration Tests

1. **Spinner Display:**
   - Create spinner, click arrow buttons
   - Verify displayed value matches internal Value
   - Verify OnChanged callback fires

2. **Timer Controls:**
   - Start timer, verify Pause button shows "Pause"
   - Pause timer, verify Resume button appears
   - Click Resume, verify timer continues
   - Reset, verify back to Start state

3. **Full Window:**
   - Open Pomodoro window
   - Verify size is correct (~320-350px width)
   - Verify all elements properly spaced
   - Verify Configuration header is centered
   - Click arrows multiple times with animations
   - Hover over arrows to verify hover effects

### Visual/Manual Tests

1. **Theme Testing:**
   - Test on light theme
   - Test on dark (Gruvbox) theme
   - Verify colors are appropriate for both

2. **Animation Testing:**
   - Click arrows rapidly → see animation each time
   - Hover then click → both effects visible
   - Move mouse away → hover effect disappears

3. **Layout Testing:**
   - Arrows should be centered in spinner
   - Configuration text should be centered
   - Window should be compact (~320-350px)
   - No content should overflow or clip

---

## Files to Modify

| File | Location | Changes | Complexity |
|------|----------|---------|-----------|
| `pomodoro_window.go` | Lines 66, 119-147, 179, 293-310, 327-340, 409-435, 729-732 | Window size, button logic, colors, container layout | High |
| `style_helpers.go` | Lines 597-626, 711-785 | TinyIconButton enhancements, NumberSpinner renderer redesign | High |

---

## Rollback Plan

If issues arise during implementation:

1. **For window size changes:** Revert line 66 to `pw.window.Resize(fyne.NewSize(400, 600))`
2. **For button logic:** Revert lines 293-310 to original state
3. **For NumberSpinner:** Create feature branch, test extensively before merge
4. **For animations:** Disable by removing animation logic, keep structural changes

---

## Success Criteria

- [ ] Window size is reduced to 320-350px width
- [ ] Progress ring gradient is Yellow-Orange → Green
- [ ] Spinner arrows update displayed value immediately
- [ ] Pause button shows "Resume" when paused, Start button is disabled
- [ ] Clicking arrow buttons shows brief visual feedback
- [ ] Hovering over arrow buttons shows highlight that disappears on mouseout
- [ ] Arrow buttons are vertically centered in spinner
- [ ] "Configuration" header text is center-aligned
- [ ] All changes work on both light and dark themes
- [ ] No regressions in existing functionality
- [ ] All tests pass

---

## Risk Assessment

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|-----------|
| NumberSpinner refactoring breaks display | Medium | High | Extensive testing, use feature branch |
| Window resize breaks layout | Low | Medium | Test all screen sizes |
| Button state logic confusion | Low | Low | Clear comment, validate with team |
| Animation performance issues | Low | Low | Monitor frame rate during testing |
| Theme color contrast issues | Low | Medium | Test on both light/dark themes |

---

## Documentation Updates

After implementation, update:

1. **CLAUDE.md** - Add section about PomodoroWindow customization
2. **Code comments** - Document button state logic clearly
3. **README** - If Pomodoro feature is documented there
4. **Examples** - If any example code uses Pomodoro window

---

## Timeline

- **Phase 1 (Critical Fixes):** 2-3 hours
- **Phase 2 (Core Functionality):** 4-5 hours
- **Phase 3 (Polish):** 3-4 hours
- **Testing & Documentation:** 2-3 hours
- **Total Estimated:** 11-15 hours (1.5-2 days of development)

---

## Approval & Sign-Off

- [ ] Requirements reviewed and approved
- [ ] Technical approach approved
- [ ] Timeline agreed
- [ ] Testing strategy validated
- [ ] Ready to begin implementation

---

## Changelog

**v1.0 - 2025-11-17**
- Initial CR document created
- All 8 issues identified and documented
- Implementation plan established
- Testing strategy defined
