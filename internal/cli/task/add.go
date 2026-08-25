package task

import (
	"fmt"
	"strings"
	"time"

	"tracker/internal/app"
	"tracker/internal/config"
	"tracker/internal/service"
	"tracker/internal/ui"
	"tracker/pkg/timeparse"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var AddCmd = &cobra.Command{
	Use:   "add [название]",
	Short: "Создать новую задачу",
	Args:  cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		interactive, _ := cmd.Flags().GetBool("interactive")

		if len(args) == 0 && !interactive && ui.IsInteractive() {
			interactive = true
		}

		if interactive && ui.IsInteractive() {
			return addInteractive()
		}

		if len(args) == 0 {
			return fmt.Errorf("укажите название задачи или используйте -i для интерактивного режима")
		}

		return addFromArgs(cmd, args)
	},
}

func addFromArgs(cmd *cobra.Command, args []string) error {
	title := strings.Join(args, " ")
	if err := service.ValidateTitle(title); err != nil {
		return err
	}

	start, _ := cmd.Flags().GetString("start")
	end, _ := cmd.Flags().GetString("end")
	company, _ := cmd.Flags().GetString("company")
	assignee, _ := cmd.Flags().GetString("assignee")
	solution, _ := cmd.Flags().GetString("solution")
	comment, _ := cmd.Flags().GetString("comment")
	tagNames, _ := cmd.Flags().GetStringSlice("tag")

	if err := service.ValidateComment(comment); err != nil {
		return err
	}
	if err := service.ValidateSolution(solution); err != nil {
		return err
	}

	if company == "" {
		if server, err := config.GetCurrentServer(); err == nil && server.DefaultCompany != "" {
			company = server.DefaultCompany
		}
	}

	startTime, err := timeparse.Parse(start)
	if err != nil {
		return fmt.Errorf("ошибка в start: %w", err)
	}

	payload := map[string]interface{}{
		"title":        title,
		"company_name": company,
		"start_time":   startTime.UTC().Format(time.RFC3339),
	}

	if end != "" {
		endTime, err := timeparse.Parse(end)
		if err != nil {
			return fmt.Errorf("ошибка в end: %w", err)
		}
		payload["end_time"] = endTime.UTC().Format(time.RFC3339)
	}

	if assignee != "" {
		payload["assignee_username"] = assignee
	}
	if solution != "" {
		payload["solution"] = solution
	}
	if comment != "" {
		payload["comment"] = comment
	}

	if len(tagNames) > 0 {
		tagIDs, err := resolveTagNamesToIDs(tagNames)
		if err != nil {
			return err
		}
		payload["tag_ids"] = tagIDs
	}

	c := app.GetClient()
	task, err := c.CreateTask(payload)
	if err != nil {
		return err
	}

	fmt.Println(ui.Checkmark(), ui.Successf("Создана задача %s: %s",
		ui.Ticket(task.Ticket), ui.Bold(task.Title)))
	fmt.Println(ui.Dimf("Исполнитель: %s", task.GetAssigneeDisplay()))
	return nil
}

func addInteractive() error {
	fmt.Println()
	ui.Header("Создание задачи")
	fmt.Println()

	var title string
	var company string
	var assignee string
	var solution string
	var comment string
	var startStr string
	var selectedTags []string

	defaultCompany := ""
	if server, err := config.GetCurrentServer(); err == nil && server.DefaultCompany != "" {
		defaultCompany = server.DefaultCompany
	}

	companies, _ := loadCompanies()
	tags, _ := loadTags()

	companyOptions := make([]huh.Option[string], 0, len(companies)+1)
	if defaultCompany != "" {
		companyOptions = append(companyOptions, huh.NewOption(defaultCompany+" (по умолчанию)", defaultCompany))
	}
	for _, c := range companies {
		if c != defaultCompany {
			companyOptions = append(companyOptions, huh.NewOption(c, c))
		}
	}

	tagOptions := make([]huh.Option[string], 0, len(tags))
	for _, t := range tags {
		tagOptions = append(tagOptions, huh.NewOption(t.Name, t.Name))
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Название задачи:").
				Value(&title).
				Validate(func(s string) error {
					return service.ValidateTitle(s)
				}),

			huh.NewSelect[string]().
				Title("Компания:").
				Options(companyOptions...).
				Value(&company),

			huh.NewInput().
				Title("Исполнитель (опционально):").
				Value(&assignee),

			huh.NewInput().
				Title("Статус решения (опционально):").
				Value(&solution).
				Validate(func(s string) error {
					if s == "" {
						return nil
					}
					return service.ValidateSolution(s)
				}),
		),

		huh.NewGroup(
			huh.NewInput().
				Title("Начало (Enter = now):").
				Value(&startStr).
				Placeholder("now").
				Validate(func(s string) error {
					if s == "" || s == "now" {
						return nil
					}
					_, err := timeparse.Parse(s)
					return err
				}),

			huh.NewMultiSelect[string]().
				Title("Теги (пробел — выбрать, enter — подтвердить):").
				Options(tagOptions...).
				Value(&selectedTags),
		),

		huh.NewGroup(
			huh.NewText().
				Title("Комментарий (опционально):").
				Description("Поддерживается Markdown. Ctrl+D — завершить").
				Value(&comment),
		),
	)

	if err := form.Run(); err != nil {
		return fmt.Errorf("форма отменена: %w", err)
	}

	if startStr == "" || startStr == "now" {
		startStr = "now"
	}

	startTime, err := timeparse.Parse(startStr)
	if err != nil {
		return fmt.Errorf("ошибка в start: %w", err)
	}

	payload := map[string]interface{}{
		"title":        title,
		"company_name": company,
		"start_time":   startTime.UTC().Format(time.RFC3339),
	}

	if assignee != "" {
		payload["assignee_username"] = assignee
	}
	if solution != "" {
		payload["solution"] = solution
	}
	if comment != "" {
		payload["comment"] = comment
	}

	if len(selectedTags) > 0 {
		tagIDs, err := resolveTagNamesToIDs(selectedTags)
		if err != nil {
			return err
		}
		payload["tag_ids"] = tagIDs
	}

	c := app.GetClient()
	task, err := c.CreateTask(payload)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println(ui.SuccessAction("Создана задача", task.Ticket))
	fmt.Println(ui.Dimf("  Название: %s", task.Title))
	if len(selectedTags) > 0 {
		fmt.Println(ui.Dimf("Теги: %s", strings.Join(selectedTags, ", ")))
	}
	fmt.Println()

	return nil
}

func loadCompanies() ([]string, error) {
	c := app.GetClient()
	resp, err := c.ListCompanies(0, 0)
	if err != nil {
		return nil, err
	}
	companies := make([]string, 0, len(resp.Companies))
	for _, c := range resp.Companies {
		companies = append(companies, c.Name)
	}
	return companies, nil
}

func loadTags() ([]struct{ Name string }, error) {
	tags, err := getTagsFromCache()
	if err != nil {
		return nil, err
	}
	result := make([]struct{ Name string }, 0, len(tags))
	for _, t := range tags {
		result = append(result, struct{ Name string }{Name: t.Name})
	}
	return result, nil
}

func init() {
	AddCmd.Flags().BoolP("interactive", "i", false, "Интерактивный режим")
	AddCmd.Flags().StringP("start", "s", "now", "Начало")
	AddCmd.Flags().StringP("end", "e", "", "Конец")
	AddCmd.Flags().StringP("company", "q", "", "Компания (по умолчанию — из конфига)")
	AddCmd.Flags().StringP("assignee", "a", "", "Исполнитель")
	AddCmd.Flags().StringP("solution", "S", "", "Статус")
	AddCmd.Flags().StringP("comment", "C", "", "Комментарий")
	AddCmd.Flags().StringSliceP("tag", "T", nil, "Теги задачи (можно указать несколько)")
}
