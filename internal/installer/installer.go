package installer

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"tracker/internal/ui"
)

func Install() error {
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("не удалось получить путь к бинарнику: %w", err)
	}
	execPath, err = filepath.Abs(execPath)
	if err != nil {
		return fmt.Errorf("не удалось получить абсолютный путь: %w", err)
	}

	switch runtime.GOOS {
	case "windows":
		return installWindows(execPath)
	case "linux", "darwin":
		return installUnix(execPath)
	default:
		return fmt.Errorf("неподдерживаемая ОС: %s", runtime.GOOS)
	}
}

func Uninstall() error {
	switch runtime.GOOS {
	case "windows":
		return uninstallWindows()
	case "linux", "darwin":
		return uninstallUnix()
	default:
		return fmt.Errorf("неподдерживаемая ОС: %s", runtime.GOOS)
	}
}

func installWindows(execPath string) error {
	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData == "" {
		localAppData = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local")
	}

	targetDir := filepath.Join(localAppData, "Programs", "tracker")
	targetPath := filepath.Join(targetDir, "tracker.exe")

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("не удалось создать директорию %s: %w", targetDir, err)
	}

	if err := copyFile(execPath, targetPath); err != nil {
		return fmt.Errorf("не удалось скопировать бинарник: %w", err)
	}

	fmt.Println(ui.Checkmark(), ui.Successf("Бинарник скопирован в %s", targetPath))

	if isInWindowsPath(targetDir) {
		fmt.Println(ui.Checkmark(), ui.Success("Директория уже в PATH"))
	} else {
		if err := addToWindowsPath(targetDir); err != nil {
			fmt.Println(ui.Warningf("Не удалось добавить в PATH автоматически: %v", err))
			fmt.Println()
			fmt.Println("Добавьте директорию в PATH вручную:")
			fmt.Printf("  %s\n", targetDir)
			fmt.Println()
			fmt.Println("Или выполните в PowerShell от администратора:")
			fmt.Printf("  setx PATH \"$env:PATH;%s\"\n", targetDir)
		} else {
			fmt.Println(ui.Checkmark(), ui.Successf("Директория добавлена в PATH: %s", targetDir))
			fmt.Println(ui.Warning("Перезапустите терминал для применения изменений"))
		}
	}

	fmt.Println()
	fmt.Println(ui.Success("Установка завершена!"))
	fmt.Println(ui.Dim("Теперь вы можете использовать команду 'tracker' в новом терминале"))
	return nil
}

func installUnix(execPath string) error {
	targetPath := "/usr/local/bin/tracker"

	if err := copyFile(execPath, targetPath); err != nil {
		fmt.Println(ui.Warning("Требуется sudo для установки в /usr/local/bin"))
		if err := runSudo("cp", execPath, targetPath); err != nil {
			return fmt.Errorf("не удалось скопировать бинарник: %w", err)
		}
		if err := runSudo("chmod", "+x", targetPath); err != nil {
			return fmt.Errorf("не удалось сделать бинарник исполняемым: %w", err)
		}
	} else {
		if err := os.Chmod(targetPath, 0755); err != nil {
			return fmt.Errorf("не удалось установить права: %w", err)
		}
	}

	fmt.Println(ui.Checkmark(), ui.Successf("Установлено в %s", targetPath))
	fmt.Println(ui.Success("Теперь вы можете использовать команду 'tracker'"))
	return nil
}

func uninstallWindows() error {
	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData == "" {
		localAppData = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local")
	}

	targetDir := filepath.Join(localAppData, "Programs", "tracker")

	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		fmt.Println(ui.Warning("Tracker не установлен"))
		return nil
	}

	psScript := fmt.Sprintf(`
$currentPath = [Environment]::GetEnvironmentVariable('Path', 'User')
$newPath = ($currentPath -split ';' | Where-Object { $_ -ne '%s' }) -join ';'
[Environment]::SetEnvironmentVariable('Path', $newPath, 'User')
`, targetDir)

	cmd := exec.Command("powershell", "-NoProfile", "-Command", psScript)
	cmd.Run()

	if err := os.RemoveAll(targetDir); err != nil {
		return fmt.Errorf("не удалось удалить: %w", err)
	}

	fmt.Println(ui.Checkmark(), ui.Success("Tracker удалён"))
	fmt.Println(ui.Warning("Перезапустите терминал для применения изменений"))
	return nil
}

func uninstallUnix() error {
	targetPath := "/usr/local/bin/tracker"

	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		fmt.Println(ui.Warning("Tracker не установлен"))
		return nil
	}

	if err := runSudo("rm", targetPath); err != nil {
		return fmt.Errorf("не удалось удалить: %w", err)
	}

	fmt.Println(ui.Checkmark(), ui.Success("Tracker удалён"))
	return nil
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return err
	}

	return destFile.Sync()
}

func runSudo(args ...string) error {
	cmd := exec.Command("sudo", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func isInWindowsPath(dir string) bool {
	pathEnv := os.Getenv("PATH")
	paths := strings.Split(pathEnv, ";")
	dirLower := strings.ToLower(dir)

	for _, p := range paths {
		if strings.ToLower(strings.TrimSpace(p)) == dirLower {
			return true
		}
	}
	return false
}

func addToWindowsPath(dir string) error {
	psScript := fmt.Sprintf(`
$currentPath = [Environment]::GetEnvironmentVariable('Path', 'User')
if ($currentPath -notlike '*%s*') {
    $newPath = $currentPath + ';%s'
    [Environment]::SetEnvironmentVariable('Path', $newPath, 'User')
    Write-Host 'PATH updated'
} else {
    Write-Host 'Already in PATH'
}
`, dir, dir)

	cmd := exec.Command("powershell", "-NoProfile", "-Command", psScript)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
