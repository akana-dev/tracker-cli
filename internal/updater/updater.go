package updater

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"tracker/internal/ui"
)

func DownloadAndInstall(releaseURL, expectedVersion string) error {
	assetName := getAssetName()

	baseURL := strings.TrimSuffix(releaseURL, "/tag/"+expectedVersion)
	baseURL = strings.TrimSuffix(baseURL, "/tag/v"+expectedVersion)
	binaryURL := fmt.Sprintf("%s/download/v%s/%s", baseURL, expectedVersion, assetName)
	checksumsURL := fmt.Sprintf("%s/download/v%s/checksums.txt", baseURL, expectedVersion)

	fmt.Println("  Скачивание " + assetName + "...")

	tmpFile, err := downloadFile(binaryURL)
	if err != nil {
		return fmt.Errorf("ошибка скачивания бинарника: %w", err)
	}
	defer os.Remove(tmpFile)

	fmt.Println("  Проверка целостности...")
	checksumsFile, err := downloadFile(checksumsURL)
	if err != nil {
		return fmt.Errorf("ошибка скачивания checksums: %w", err)
	}
	defer os.Remove(checksumsFile)

	if err := verifyChecksum(tmpFile, checksumsFile, assetName); err != nil {
		return fmt.Errorf("ошибка проверки целостности: %w", err)
	}

	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("не удалось получить путь к бинарнику: %w", err)
	}
	execPath, err = filepath.Abs(execPath)
	if err != nil {
		return fmt.Errorf("не удалось получить абсолютный путь: %w", err)
	}

	fmt.Println("  Установка обновления...")
	if err := replaceBinary(tmpFile, execPath); err != nil {
		fmt.Printf("  Ошибка: %v\n", err)
		fmt.Println()
		fmt.Println(ui.Warning("Не удалось автоматически обновить tracker."))
		fmt.Println()
		fmt.Println("Возможные причины:")
		fmt.Println("  1. Нет прав на запись в директорию с бинарником")
		fmt.Println("  2. Антивирус или защита ОС блокируют замену")
		fmt.Println()
		fmt.Println("Решение:")
		fmt.Println("  1. Закройте все экземпляры tracker")
		fmt.Printf("  2. Выполните: %s\n", ui.Cyan("tracker update"))
		fmt.Printf("  3. Или скачайте вручную: %s\n", ui.Dim(releaseURL))
		fmt.Println()
		return fmt.Errorf("ошибка обновления: %w", err)
	}

	fmt.Println("Обновление успешно завершено!")
	return nil
}

func getAssetName() string {
	targetOS := runtime.GOOS
	targetArch := runtime.GOARCH

	var ext string
	if targetOS == "windows" {
		ext = ".exe"
	}

	return fmt.Sprintf("tracker-%s-%s%s", targetOS, targetArch, ext)
}

func downloadFile(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	tmpFile, err := os.CreateTemp("", "tracker-update-*")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		os.Remove(tmpFile.Name())
		return "", err
	}

	return tmpFile.Name(), nil
}

func verifyChecksum(filePath, checksumsFile, assetName string) error {
	data, err := os.ReadFile(checksumsFile)
	if err != nil {
		return err
	}

	var expectedChecksum string
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) == 2 && parts[1] == assetName {
			expectedChecksum = parts[0]
			break
		}
	}

	if expectedChecksum == "" {
		return fmt.Errorf("checksum для %s не найден", assetName)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return err
	}

	actualChecksum := hex.EncodeToString(hasher.Sum(nil))

	if actualChecksum != expectedChecksum {
		return fmt.Errorf("checksum не совпадает: ожидался %s, получен %s", expectedChecksum, actualChecksum)
	}

	return nil
}

func extractBinary(archivePath, assetName string) (string, error) {
	if !strings.HasSuffix(archivePath, ".tar.gz") && !strings.HasSuffix(archivePath, ".zip") {
		return archivePath, nil
	}

	tmpFile, err := os.CreateTemp("", "tracker-binary-*")
	if err != nil {
		return "", err
	}
	tmpFile.Close()

	if strings.HasSuffix(archivePath, ".tar.gz") {
		return extractTarGz(archivePath, tmpFile.Name())
	} else if strings.HasSuffix(archivePath, ".zip") {
		return extractZip(archivePath, tmpFile.Name())
	}

	return archivePath, nil
}

func extractTarGz(archivePath, outputPath string) (string, error) {
	file, err := os.Open(archivePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return "", err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		if header.Typeflag == tar.TypeReg {
			outFile, err := os.Create(outputPath)
			if err != nil {
				return "", err
			}
			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return "", err
			}
			outFile.Close()
			return outputPath, nil
		}
	}

	return "", fmt.Errorf("бинарник не найден в архиве")
}

func extractZip(archivePath, outputPath string) (string, error) {
	r, err := zip.OpenReader(archivePath)
	if err != nil {
		return "", err
	}
	defer r.Close()

	for _, f := range r.File {
		if !f.FileInfo().IsDir() {
			rc, err := f.Open()
			if err != nil {
				return "", err
			}

			outFile, err := os.Create(outputPath)
			if err != nil {
				rc.Close()
				return "", err
			}

			if _, err := io.Copy(outFile, rc); err != nil {
				outFile.Close()
				rc.Close()
				return "", err
			}

			outFile.Close()
			rc.Close()
			return outputPath, nil
		}
	}

	return "", fmt.Errorf("бинарник не найден в архиве")
}

func replaceBinary(newBinaryPath, targetPath string) error {
	oldPath := targetPath + ".old"

	if err := os.Rename(targetPath, oldPath); err != nil {
		return fmt.Errorf("не удалось создать резервную копию: %w", err)
	}

	if err := copyFile(newBinaryPath, targetPath); err != nil {
		os.Rename(oldPath, targetPath)
		return fmt.Errorf("ошибка копирования нового бинарника: %w", err)
	}

	os.Remove(oldPath)
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
