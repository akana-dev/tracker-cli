package export

import (
	"fmt"
	"os"
	"strings"

	"tracker/internal/client"
	"tracker/internal/service"
	"tracker/internal/ui"
	"tracker/pkg/table"
)

func RunPreview(params map[string]string, format string) error {
	fmt.Println()
	ui.Header("Preview режима экспорта")
	fmt.Println()

	listParams := map[string]string{}
	for k, v := range params {
		if k == "format" || k == "filename" || k == "fields" {
			continue
		}
		listParams[k] = v
	}
	listParams["limit"] = "5"

	resp, err := client.ListTasks(listParams, 5, 0)
	if err != nil {
		return fmt.Errorf("ошибка preview: %w", err)
	}

	tasks := resp.Tasks
	if len(tasks) == 0 {
		fmt.Println(ui.Warning("Задачи не найдены для экспорта."))
		return nil
	}

	tbl := table.New("Тикет", "Дата", "Задача", "Часы", "Статус")
	tbl.SetColumnWidths(map[int]int{0: 10, 1: 12, 2: 45, 3: 6, 4: 20})

	for _, t := range tasks {
		hours := fmt.Sprintf("%.1f", service.CalculateTaskHours(t))
		status := service.FormatStatus(t)
		tbl.AddRow(
			ui.Ticket(t.Ticket),
			t.StartTime.Local().Format("02.01.2006"),
			service.FormatTaskCell(t),
			hours,
			status,
		)
	}
	tbl.Render()

	fmt.Println()
	if resp.Total > 5 {
		fmt.Printf("  %s\n", ui.Dim(fmt.Sprintf("Показано 5 из ~%d задач", resp.Total)))
	} else {
		fmt.Printf("  %s\n", ui.Dim(fmt.Sprintf("Найдено задач: %d", resp.Total)))
	}
	fmt.Println()

	fmt.Print("Продолжить экспорт? [Y/n]: ")
	answer := readLine()
	answer = strings.ToLower(strings.TrimSpace(answer))
	if answer == "n" || answer == "no" {
		fmt.Println(ui.Warning("Экспорт отменён."))
		return nil
	}

	fmt.Println()
	ui.Header("Выполнение экспорта...")

	data, apiFilename, err := client.ExportTasks(params)
	if err != nil {
		return err
	}

	output := apiFilename
	if v, ok := params["filename"]; ok && v != "" {
		output = v
	}

	if err := os.WriteFile(output, data, 0644); err != nil {
		return err
	}

	fmt.Println(ui.Checkmark(), ui.Successf("Экспортировано в %s", ui.Bold(output)))
	return nil
}

func readLine() string {
	var line string
	fmt.Scanln(&line)
	return line
}
