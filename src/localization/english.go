package localization

import "fmt"

// English provides English language strings for the application
var English = map[string]string{
	// Window and Navigation Elements
	"window_title":         "Go Do",
	"previous_month":       "<",
	"next_month":           ">",
	"view_mode_all":        "All",
	"view_mode_incomplete": "Incomplete",
	"view_mode_reminders":  "Reminders",

	// Form Elements
	"form_title_add":     "Add Todo",
	"form_title_edit":    "Edit Todo",
	"form_button_add":    "Add",
	"form_button_save":   "Save",
	"form_button_cancel": "Cancel",

	// Field Labels and Placeholders
	"field_name":                 "Name:",
	"field_name_placeholder":     "Enter todo item",
	"field_content":              "Content:",
	"field_content_placeholder":  "Content:",
	"field_location":             "Location:",
	"field_location_placeholder": "Location:",
	"field_label":                "Label:",
	"field_label_placeholder":    "Label:",
	"field_datetime":             "Date/Time:",
	"field_datetime_placeholder": "Date/Time (DD.MM.YYYY HH:MM)",
	"field_type":                 "Type:",
	"field_priority":             "Priority:",
	"field_reminder":             "Reminder:",
	"select_datetime":            "Select Date/Time",

	// Priority Levels
	"priority_0": "Not Important - Not Urgent",
	"priority_1": "Not Important - Urgent",
	"priority_2": "Important - Not Urgent",
	"priority_3": "Important - Urgent",

	// Types
	"type_event": "Event",
	"type_task":  "Task",

	// Reminder Messages
	"reminder_none":   "No reminder",
	"reminder_format": "Remind %s before",

	// Status Messages
	"status_empty_list":    "No todos yet. Click + to add your first todo!",
	"status_loading_error": "Error loading todos: %s",

	// Error Messages
	"error_name_required":    "Name is required",
	"error_invalid_datetime": "Invalid date/time format. Use DD.MM.YYYY HH:MM",
	"error_save_failed":      "Failed to save todo: %s",
	"error_load_failed":      "Failed to load todos: %s",

	// Success Messages
	"success_todo_saved":   "Todo saved successfully",
	"success_todo_deleted": "Todo deleted successfully",

	// Time Units
	"time_days":    "days",
	"time_hours":   "hours",
	"time_minutes": "minutes",
	"time_day":     "day",
	"time_hour":    "hour",
	"time_minute":  "minute",

	// Confirmation Dialogs
	"confirm_delete_title":       "Delete Todo",
	"confirm_delete_message":     "Are you sure you want to delete this todo?",
	"confirm_delete_all_title":   "Delete All Todos",
	"confirm_delete_all_message": "Are you sure you want to delete all todos?",

	// Menu Items
	"menu_file":  "File",
	"menu_edit":  "Edit",
	"menu_view":  "View",
	"menu_help":  "Help",
	"menu_exit":  "Exit",
	"menu_about": "About",

	// Tooltips and Help
	"tooltip_add_todo":        "Add new todo (Ctrl+N)",
	"tooltip_edit_todo":       "Edit selected todo (Ctrl+E)",
	"tooltip_delete_todo":     "Delete selected todo (Delete)",
	"tooltip_toggle_complete": "Mark as complete/incomplete",
	"tooltip_view_mode":       "Change view mode",
	"tooltip_navigation":      "Navigate between months",

	// Keyboard Shortcuts
	"shortcut_new":    "Ctrl+N",
	"shortcut_edit":   "Ctrl+E",
	"shortcut_delete": "Delete",
	"shortcut_save":   "Ctrl+S",
	"shortcut_cancel": "Esc",

	// Application Info
	"app_name":        "Go Do Todo List",
	"app_version":     "1.0.0",
	"app_description": "A simple and elegant todo list application",
	"app_author":      "Migration Team",
	"app_website":     "https://example.com",

	// Color Names for Priority Levels
	"color_green":  "Green",
	"color_blue":   "Blue",
	"color_orange": "Orange",
	"color_red":    "Red",
}

// GetString retrieves a localized string by key
func GetString(key string) string {
	if str, exists := English[key]; exists {
		return str
	}
	return key // Return key as fallback if not found
}

// GetStringWithArgs retrieves a localized string and formats it with arguments
func GetStringWithArgs(key string, args ...interface{}) string {
	format := GetString(key)
	return fmt.Sprintf(format, args...)
}
