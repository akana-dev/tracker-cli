package ui

import "fmt"

func SuccessAction(action, target string) string {
	return Successf("%s %s", Checkmark(), fmt.Sprintf("%s: %s", action, Bold(target)))
}

func WarningAction(action, target string) string {
	return Warningf("⚠ %s: %s", action, Bold(target))
}

func ErrorAction(action, target, reason string) string {
	if reason == "" {
		return Errorf("%s %s: %s", Cross(), action, Bold(target))
	}
	return Errorf("%s %s: %s — %s", Cross(), action, Bold(target), reason)
}

func InfoAction(action, target string) string {
	return Infof("ℹ %s: %s", action, Bold(target))
}

func ProgressMessage(action, target string) string {
	return Dimf("→ %s %s...", action, target)
}

func CompletionMessage(action, target string) string {
	return Successf("✓ %s %s завершено", action, target)
}

func CountMessage(noun string, count int) string {
	return Dimf("%s: %d", noun, count)
}

func StatusMessage(status, target string) string {
	return Infof("[%s] %s", status, target)
}

func ConfirmMessage(action, target string) string {
	return Warningf("Вы уверены, что хотите %s %s?", action, target)
}

func HelpMessage(command, description string) string {
	return fmt.Sprintf("  %s  %s", Cyan(command), Dim(description))
}

func ExampleMessage(example string) string {
	return fmt.Sprintf("  %s %s", Dim("Пример:"), Cyan(example))
}
