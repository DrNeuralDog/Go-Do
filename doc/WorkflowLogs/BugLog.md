[2025-10-02 11:36] Layout still shifts on first two filter switches - Attempt 2: moved Scroll inside Timeline renderer only and removed external wrapper; added forced resize on refresh.
[2025-10-02 11:28] Layout shrinks on view switch (extra margins right/bottom) - Attempt 1: removed fixed MinSize in Timeline and forced container refresh after SetTodos. Issue persists partially after restart.
[2025-10-02 11:05] Checkbox toggle triggers infinite save loop (recreating YAML temp files) - Fixed by deferring OnChanged binding until after SetChecked
# Bug Log for Todo List Application Migration

## Overview
This log tracks bugs, errors, and debugging activities during the Todo List application migration from C++ Qt to Go Fyne.

## Log Format
```
[Timestamp] Brief description of error/debug - Result
```

## Bug Reports and Resolutions

[2025-10-02 15:05:00 UTC] Нет зафиксированных ошибок - Проект в начальной стадии разработки
[2025-10-02 20:30:00 UTC] При первом переключении фильтра (All/Important/Complete) контент сжимается, огромные поля справа и снизу - Исправлено: в `timelineRenderer.Refresh()` создаётся Scroll при nil и вызывается `Layout(r.timeline.Size())` для корректного лэйаута
