package initcontainers

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func Init() {
	systemContainersPath := "/var/lib/contos/system-containers"
	userContainersPath := os.ExpandEnv("$HOME/contos/user-containers")

	createDirIfNotExists(systemContainersPath)
	createDirIfNotExists(userContainersPath)

	startContainers(systemContainersPath)
	startContainers(userContainersPath)
}

func createDirIfNotExists(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			log.Fatalf("Ошибка при создании каталога %s: %v", path, err)
		}
	}
}

func startContainers(path string) {
	files, err := os.ReadDir(path)
	if err != nil {
		log.Printf("Ошибка при чтении каталога контейнеров: %v", err)
		return
	}

	for _, file := range files {
		containerName := file.Name()
		containerArchivePath := filepath.Join(path, containerName)

		if !isTarArchive(containerArchivePath) {
			log.Printf("Файл %s не является архивом .tar, пропускаем.", containerArchivePath)
			continue
		}

		log.Printf("Извлечение и запуск контейнера: %s", containerName)

		containerDir := filepath.Join(path, containerName+"_extracted")
		err := extractTarArchive(containerArchivePath, containerDir)
		if err != nil {
			log.Printf("Ошибка при извлечении архива %s: %v", containerArchivePath, err)
			continue
		}

		startContainer(containerName, containerDir)
	}
}

func isTarArchive(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		log.Printf("Не удалось открыть файл %s: %v", path, err)
		return false
	}
	defer file.Close()

	// Проверяем первые 512 байт на сигнатуру формата .tar
	buf := make([]byte, 512)
	_, err = file.Read(buf)
	if err != nil {
		log.Printf("Ошибка при чтении файла %s: %v", path, err)
		return false
	}

	// Проверяем сигнатуру .tar
	return bytes.Compare(buf[257:262], []byte("ustar")) == 0
}

func extractTarArchive(archivePath, targetDir string) error {
	file, err := os.Open(archivePath)
	if err != nil {
		return fmt.Errorf("не удалось открыть архив %s: %v", archivePath, err)
	}
	defer file.Close()

	tarReader := tar.NewReader(file)

	err = os.MkdirAll(targetDir, 0755)
	if err != nil {
		return fmt.Errorf("не удалось создать каталог для извлечения: %v", err)
	}

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("ошибка при извлечении архива: %v", err)
		}

		targetPath := filepath.Join(targetDir, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			err := os.MkdirAll(targetPath, 0755)
			if err != nil {
				return fmt.Errorf("не удалось создать каталог %s: %v", targetPath, err)
			}
		case tar.TypeReg:
			targetFile, err := os.Create(targetPath)
			if err != nil {
				return fmt.Errorf("не удалось создать файл %s: %v", targetPath, err)
			}
			defer targetFile.Close()

			_, err = io.Copy(targetFile, tarReader)
			if err != nil {
				return fmt.Errorf("ошибка при извлечении файла %s: %v", targetPath, err)
			}
		}
	}

	return nil
}

func startContainer(containerName, containerDir string) {
	log.Printf("Запуск контейнера %s из каталога %s", containerName, containerDir)

	// Команда для запуска контейнера с использованием containerd
	cmdStart := exec.Command("ctr", "tasks", "start", containerName)
	cmdStart.Dir = containerDir
	cmdStart.Stdout = os.Stdout
	cmdStart.Stderr = os.Stderr

	if err := cmdStart.Run(); err != nil {
		log.Printf("Ошибка при запуске контейнера %s: %v", containerName, err)
	}
}
