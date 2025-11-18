# Текущий CR: Cursor IDE конфиги для запуска Go/Fyne (Debug/Release)

## Цель
Обеспечить запуск и отладку приложения на Go (Fyne) в Cursor по F5, с профилями Debug и Release.

## Статус задач
- [x] Создать задачи сборки (.vscode/tasks.json) для Debug и Release
- [x] Создать конфигурации запуска (.vscode/launch.json) для Debug и Release
- [x] Добавить базовые IDE-настройки (.vscode/settings.json)
- [x] Зафиксировать изменения в журналах (docs/WorkflowLogs/*.md)
- [ ] Проверка: вручную запустить F5 (Debug) и профиль Release (без отладки)

## Примечания
- Debug использует предварительную сборку с флагами -gcflags "all=-N -l" и запуск через Delve (режим exec).
- Release собирает бинарник в bin/release с -ldflags "-s -w -H=windowsgui" и запускается без отладчика.


