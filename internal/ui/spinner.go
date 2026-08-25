package ui

import (
	"fmt"
	"sync"

	"github.com/pterm/pterm"
)

var spinnerMu sync.Mutex

type ExportResult struct {
	Data     []byte
	Filename string
}

func WithSpinner(message string, fn func() error) error {
	if !IsInteractive() {
		fmt.Println(message + "...")
		return fn()
	}

	spinnerMu.Lock()
	defer spinnerMu.Unlock()

	spinner, _ := pterm.DefaultSpinner.
		WithShowTimer(true).
		WithTimerStyle(pterm.NewStyle(pterm.FgCyan)).
		WithStyle(pterm.NewStyle(pterm.FgCyan)).
		Start(message)

	err := fn()

	if err != nil {
		spinner.Fail(fmt.Sprintf("Ошибка: %v", err))
	} else {
		spinner.Success("Готово")
	}

	return err
}

func WithSpinnerResult[T any](message string, fn func() (T, error)) (T, error) {
	var zero T

	if !IsInteractive() {
		fmt.Println(message + "...")
		return fn()
	}

	spinnerMu.Lock()
	defer spinnerMu.Unlock()

	spinner, _ := pterm.DefaultSpinner.
		WithShowTimer(true).
		WithTimerStyle(pterm.NewStyle(pterm.FgCyan)).
		WithStyle(pterm.NewStyle(pterm.FgCyan)).
		Start(message)

	result, err := fn()

	if err != nil {
		spinner.Fail(fmt.Sprintf("Ошибка: %v", err))
		return zero, err
	}

	spinner.Success("Готово")
	return result, nil
}

func WithSpinnerExport(message string, fn func() ([]byte, string, error)) ([]byte, string, error) {
	if !IsInteractive() {
		fmt.Println(message + "...")
		return fn()
	}

	spinnerMu.Lock()
	defer spinnerMu.Unlock()

	spinner, _ := pterm.DefaultSpinner.
		WithShowTimer(true).
		WithTimerStyle(pterm.NewStyle(pterm.FgCyan)).
		WithStyle(pterm.NewStyle(pterm.FgCyan)).
		Start(message)

	data, filename, err := fn()

	if err != nil {
		spinner.Fail(fmt.Sprintf("Ошибка: %v", err))
		return nil, "", err
	}

	spinner.Success("Готово")
	return data, filename, nil
}

type ProgressWriter struct {
	Total   int64
	Current int64
	Bar     *pterm.ProgressbarPrinter
	started bool
	mu      sync.Mutex
}

func NewProgressWriter(total int64, message string) *ProgressWriter {
	pw := &ProgressWriter{
		Total: total,
	}

	if IsInteractive() && total > 0 {
		pw.Bar = pterm.DefaultProgressbar.
			WithTitle(message).
			WithTotal(int(total)).
			WithShowCount(true).
			WithShowPercentage(true)
	}

	return pw
}

func (pw *ProgressWriter) Write(p []byte) (n int, err error) {
	n = len(p)
	pw.mu.Lock()
	pw.Current += int64(n)

	if pw.Bar != nil {
		if !pw.started {
			pw.Bar.Start()
			pw.started = true
		}
		pw.Bar.Current = int(pw.Current)
	}
	pw.mu.Unlock()
	return n, nil
}

func (pw *ProgressWriter) Finish() {
	pw.mu.Lock()
	defer pw.mu.Unlock()

	if pw.Bar != nil && pw.started {
		pw.Bar.Stop()
	}
}

func PrintProgress(current, total int64, message string) {
	if total == 0 {
		return
	}
	percent := (current * 100) / total
	fmt.Printf("\r%s %d%% (%d/%d)", message, percent, current, total)
	if current == total {
		fmt.Println()
	}
}
