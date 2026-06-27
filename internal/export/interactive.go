package export

import (
	"fmt"
	"os"
	"time"

	"tracker/internal/client"
	"tracker/internal/input"
	"tracker/internal/ui"
)

func RunInteractive() error {
	fmt.Println()
	ui.Header("Интерактивный экспорт задач")
	fmt.Println()

	fmt.Println("1. Формат экспорта:")
	fmt.Printf("   [1] CSV\n   [2] XLSX\n   [3] JSON\n   Выбор: ")
	formatChoice := input.ReadLine()
	format := "csv"
	switch formatChoice {
	case "2":
		format = "xlsx"
	case "3":
		format = "json"
	}

	fmt.Println("\n2. Период:")
	fmt.Println("   [1] Сегодня")
	fmt.Println("   [2] Эта неделя")
	fmt.Println("   [3] Этот месяц")
	fmt.Println("   [4] Последние 7 дней")
	fmt.Println("   [5] Последние 30 дней")
	fmt.Println("   [6] Произвольный")
	fmt.Printf("   Выбор: ")
	periodChoice := input.ReadLine()

	var dateFrom, dateTo string
	switch periodChoice {
	case "1":
		dateFrom = "today"
	case "2":
		dateFrom = "this week"
	case "3":
		dateFrom = "this month"
	case "4":
		dateFrom = "last 7 days"
	case "5":
		dateFrom = "last 30 days"
	case "6":
		fmt.Print("   Дата от (например: 2026-06-01 или 'last monday'): ")
		dateFrom = input.ReadLine()
		fmt.Print("   Дата до (например: 2026-06-30 или 'today'): ")
		dateTo = input.ReadLine()
	}

	fmt.Println("\n3. Фильтры (Enter — пропустить):")
	fmt.Print("   Компания: ")
	company := input.ReadLine()
	fmt.Print("   Исполнитель: ")
	assignee := input.ReadLine()
	fmt.Print("   Статус решения: ")
	solution := input.ReadLine()
	fmt.Print("   Поиск: ")
	search := input.ReadLine()

	fmt.Println("\n4. Дополнительные опции:")
	openOnly := input.ReadBool("   Только открытые задачи?", false)
	allUsers := input.ReadBool("   Показать задачи всех пользователей?", false)

	fmt.Println("\n5. Часовой пояс:")
	fmt.Printf("   [по умолчанию: Europe/Moscow]: ")
	timezone := input.ReadLine()
	if timezone == "" {
		timezone = "Europe/Moscow"
	}

	fmt.Println("\n6. Имя выходного файла:")
	ext := format
	defaultFilename := fmt.Sprintf("tasks_%s.%s", time.Now().Format("2006-01-02"), ext)
	fmt.Printf("   [по умолчанию: %s]: ", defaultFilename)
	output := input.ReadLine()
	if output == "" {
		output = defaultFilename
	}

	resolvedFrom, resolvedTo, err := ResolveDates("", dateFrom, dateTo)
	if err != nil {
		return err
	}

	params := map[string]string{
		"format":   format,
		"timezone": timezone,
	}
	if resolvedFrom != "" {
		params["date_from"] = resolvedFrom
	}
	if resolvedTo != "" {
		params["date_to"] = resolvedTo
	}
	if company != "" {
		params["company"] = company
	}
	if assignee != "" {
		params["assignee"] = assignee
	}
	if solution != "" {
		params["solution"] = solution
	}
	if search != "" {
		params["search"] = search
	}
	if openOnly {
		params["open_only"] = "true"
	}
	if allUsers {
		params["all_users"] = "true"
	}

	fmt.Println()
	ui.Header("Сводка экспорта")
	fmt.Printf("  Формат: %s\n", format)
	if resolvedFrom != "" {
		fmt.Printf("  Период: %s", resolvedFrom)
		if resolvedTo != "" {
			fmt.Printf(" — %s", resolvedTo)
		}
		fmt.Println()
	}
	if company != "" {
		fmt.Printf("  Компания: %s\n", company)
	}
	if assignee != "" {
		fmt.Printf("  Исполнитель: %s\n", assignee)
	}
	if solution != "" {
		fmt.Printf("  Статус: %s\n", solution)
	}
	if search != "" {
		fmt.Printf("  Поиск: %s\n", search)
	}
	fmt.Printf("  Часовой пояс: %s\n", timezone)
	fmt.Printf("  Файл: %s\n", output)
	fmt.Println()

	if !input.ReadBool("Продолжить экспорт?", true) {
		fmt.Println(ui.Warning("Экспорт отменён."))
		return nil
	}

	fmt.Println()
	data, apiFilename, err := client.ExportTasks(params)
	if err != nil {
		return err
	}

	if err := os.WriteFile(output, data, 0644); err != nil {
		return err
	}

	fmt.Println(ui.Checkmark(), ui.Successf("Экспортировано в %s", ui.Bold(output)))
	if apiFilename != output {
		fmt.Println(ui.Dimf("Имя файла от сервера: %s", apiFilename))
	}

	return nil
}
