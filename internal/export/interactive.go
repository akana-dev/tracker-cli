package export

import (
	"fmt"
	"os"
	"time"

	"tracker/internal/client"
	"tracker/internal/ui"

	"github.com/charmbracelet/huh"
)

func RunInteractive() error {
	fmt.Println()
	ui.Header("Интерактивный экспорт задач")
	fmt.Println()

	var format string
	var periodChoice int
	var dateFrom, dateTo string
	var company, assignee, solution, search string
	var openOnly, allUsers bool
	var timezone, output string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Формат экспорта:").
				Options(
					huh.NewOption("CSV", "csv"),
					huh.NewOption("XLSX (Excel)", "xlsx"),
					huh.NewOption("JSON", "json"),
				).
				Value(&format),

			huh.NewSelect[int]().
				Title("Период:").
				Options(
					huh.NewOption("Сегодня", 1),
					huh.NewOption("Эта неделя", 2),
					huh.NewOption("Этот месяц", 3),
					huh.NewOption("Последние 7 дней", 4),
					huh.NewOption("Последние 30 дней", 5),
					huh.NewOption("Произвольный", 6),
				).
				Value(&periodChoice),
		),

		huh.NewGroup(
			huh.NewInput().
				Title("Компания (Enter — пропустить):").
				Value(&company),

			huh.NewInput().
				Title("Исполнитель (Enter — пропустить):").
				Value(&assignee),

			huh.NewInput().
				Title("Статус решения (Enter — пропустить):").
				Value(&solution),

			huh.NewInput().
				Title("Поиск (Enter — пропустить):").
				Value(&search),
		),

		huh.NewGroup(
			huh.NewConfirm().
				Title("Только открытые задачи?").
				Value(&openOnly),

			huh.NewConfirm().
				Title("Показать задачи всех пользователей?").
				Value(&allUsers),

			huh.NewInput().
				Title("Часовой пояс:").
				Value(&timezone).
				Placeholder("Europe/Moscow"),

			huh.NewInput().
				Title("Имя выходного файла:").
				Value(&output).
				Placeholder(fmt.Sprintf("tasks_%s.csv", time.Now().Format("2006-01-02"))),
		),
	)

	if err := form.Run(); err != nil {
		return fmt.Errorf("экспорт отменён: %w", err)
	}

	switch periodChoice {
	case 1:
		dateFrom = "today"
	case 2:
		dateFrom = "this week"
	case 3:
		dateFrom = "this month"
	case 4:
		dateFrom = "last 7 days"
	case 5:
		dateFrom = "last 30 days"
	case 6:
		customForm := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Дата от (например: 2026-06-01 или 'last monday'):").
					Value(&dateFrom),

				huh.NewInput().
					Title("Дата до (например: 2026-06-30 или 'today'):").
					Value(&dateTo),
			),
		)
		if err := customForm.Run(); err != nil {
			return fmt.Errorf("экспорт отменён: %w", err)
		}
	}

	if timezone == "" {
		timezone = "Europe/Moscow"
	}
	if output == "" {
		output = fmt.Sprintf("tasks_%s.%s", time.Now().Format("2006-01-02"), format)
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

	var confirmed bool
	confirmForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Продолжить экспорт?").
				Affirmative("Да").
				Negative("Нет").
				Value(&confirmed),
		),
	)

	if err := confirmForm.Run(); err != nil || !confirmed {
		fmt.Println(ui.Warning("Экспорт отменён."))
		return nil
	}

	fmt.Println()
	data, apiFilename, err := ui.WithSpinnerExport("Экспорт задач", func() ([]byte, string, error) {
		return client.ExportTasks(params)
	})
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
