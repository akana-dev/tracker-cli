package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"tracker/internal/ui"
	"tracker/internal/updater"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Проверить и установить обновление",
	Long: `Проверить наличие новой версии tracker и установить её.

Примеры:
  tracker update                  # Проверить и установить
  tracker update --check          # Только проверить (с полной информацией)
  tracker update --version 1.2.3  # Установить конкретную версию
  tracker update --pre-release    # Включая pre-release версии`,
	RunE: func(cmd *cobra.Command, args []string) error {
		checkOnly, _ := cmd.Flags().GetBool("check")
		specificVersion, _ := cmd.Flags().GetString("version")
		includePreRelease, _ := cmd.Flags().GetBool("pre-release")

		updater.ResetNotification()

		fmt.Println(ui.Dim("Проверка обновлений..."))

		ctx := context.Background()
		result, err := updater.CheckForUpdate(ctx, githubOwner, githubRepo, includePreRelease)
		if err != nil {
			return fmt.Errorf("ошибка проверки: %w", err)
		}

		if !result.HasUpdate {
			fmt.Println(ui.Checkmark(), ui.Successf("У вас актуальная версия (%s)", formatVersion(result.CurrentVersion)))
			return nil
		}

		if checkOnly {
			updater.ShowFullUpdateInfo(result)
			return nil
		}

		updater.ShowFullUpdateInfo(result)

		fmt.Println(ui.Bold("Начинаю обновление..."))
		fmt.Println()

		targetVersion := result.LatestVersion
		if specificVersion != "" {
			targetVersion = strings.TrimPrefix(specificVersion, "v")
		}

		if err := updater.DownloadAndInstall(result.ReleaseURL, targetVersion); err != nil {
			return fmt.Errorf("ошибка обновления: %w", err)
		}

		fmt.Println()
		fmt.Println(ui.Checkmark(), ui.Successf("Обновление до версии %s завершено!", formatVersion(targetVersion)))
		fmt.Println(ui.Dim("Перезапустите tracker для использования новой версии."))

		return nil
	},
}

func formatVersion(v string) string {
	if v == "dev" || v == "" {
		return "dev"
	}
	if v[0] != 'v' {
		return "v" + v
	}
	return v
}

func init() {
	updateCmd.Flags().BoolP("check", "c", false, "Только проверить наличие обновления")
	updateCmd.Flags().StringP("version", "v", "", "Установить конкретную версию")
	updateCmd.Flags().BoolP("pre-release", "p", false, "Включая pre-release версии")
}
