package downloader

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jaimeph/helm-lint-kcl/internal/logger"
)

type Downloader struct {
	chart   string
	version string
}

func New(chart, version string) *Downloader {
	return &Downloader{
		chart:   chart,
		version: version,
	}
}

// Obtenemos del chart los archivos especificados
func (d *Downloader) GetFilesContents(filePaths ...string) (map[string][]byte, error) {
	chartDirPath, err := d.GetChart(d.chart, d.version)
	if err != nil {
		return nil, err
	}

	logger.Debugf("locate chart %s", chartDirPath)

	chartFileContents, err := d.ReadChart(chartDirPath, filePaths...)
	if err != nil {
		return nil, err
	}

	return chartFileContents, nil
}

// Obtenemos el chart, si no es una ruta local, lo descargamos con pull chart
func (d *Downloader) GetChart(chart, version string) (string, error) {
	newChartDirPath := chart
	logger.Debugf("get chart %s %s", chart, version)
	if _, err := os.Stat(chart); os.IsNotExist(err) {
		logger.Debugf("pull chart %s %s", chart, version)
		newChartDirPath, err = d.pullChart(chart, version)
		if err != nil {
			return "", err
		}
	}
	return newChartDirPath, nil
}

// Descargamos el chart del repositorio de helm o de oci, añadimos la versión si se especifica
func (d *Downloader) pullChart(chart, version string) (string, error) {
	cmdVersion := ""
	if len(version) > 0 {
		cmdVersion = "--version " + version
	}

	tempDir, err := os.MkdirTemp("", "helm-lint-kcl-*")
	if err != nil {
		return "", fmt.Errorf("the temporary directory could not be created: %w", err)
	}

	cmd := fmt.Sprintf("helm pull %s -d %s %s", chart, tempDir, cmdVersion)
	logger.Debugf("exec: %s", cmd)
	output, err := exec.Command("bash", "-c", cmd).CombinedOutput()

	if len(output) > 0 {
		logger.Debug(string(output))
	}
	if err != nil {
		return "", fmt.Errorf("helm pull %v", err)
	}

	chartTempPath, err := findTgz(tempDir)
	if err != nil {
		return "", fmt.Errorf("findTgz %v", err)
	}
	logger.Debugf("downloaded chart %s", chartTempPath)
	return chartTempPath, err
}

// Leemos los archivos del chart, si es un tgz, descomprimimos los archivos y los leemos
func (d *Downloader) ReadChart(chartDirPath string, filePath ...string) (map[string][]byte, error) {
	if strings.EqualFold(filepath.Ext(chartDirPath), ".tgz") {
		return d.ReadChartFileContentsTgz(chartDirPath, filePath...)
	} else {
		return d.ReadChartFileContentsLocal(chartDirPath, filePath...)
	}
}

// Leemos los archivos del chart local
func (d *Downloader) ReadChartFileContentsLocal(chartDirPath string, filePaths ...string) (map[string][]byte, error) {
	contents := make(map[string][]byte, len(filePaths))
	for _, filePath := range filePaths {
		fullPath := filepath.Join(chartDirPath, filePath)
		if _, err := os.Stat(fullPath); err == nil {
			content, err := os.ReadFile(fullPath)
			if err != nil {
				return nil, err
			}
			logger.Debugf("reading %s", filePath)
			contents[filePath] = content
		} else if os.IsNotExist(err) {
			return nil, fmt.Errorf("file %s does not exist into chart", filePath)
		} else {
			return nil, fmt.Errorf("error reading file %s: %w", filePath, err)
		}
	}
	return contents, nil
}

// Leemos los archivos del chart tgz
func (d *Downloader) ReadChartFileContentsTgz(chartDirPath string, filePaths ...string) (map[string][]byte, error) {
	contents := make(map[string][]byte, len(filePaths))
	for _, filePath := range filePaths {
		contents[filePath] = nil
	}

	file, err := os.Open(chartDirPath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return nil, fmt.Errorf("error creating gzip reader: %w", err)
	}
	defer gzr.Close()

	tarReader := tar.NewReader(gzr)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading tar header: %w", err)
		}

		parts := strings.Split(header.Name, "/")
		name := strings.Join(parts[1:], "/")

		for _, filePath := range filePaths {
			if name == filePath {
				logger.Debugf("extracting %s", name)
				var buffer bytes.Buffer
				if _, err := io.Copy(&buffer, tarReader); err != nil {
					return nil, fmt.Errorf("error copying file contents: %w", err)
				}
				contents[filePath] = buffer.Bytes()
			}
		}
	}

	for filePath, content := range contents {
		if content == nil {
			return nil, fmt.Errorf("file %s does not exist into chart", filePath)
		}
	}

	return contents, nil
}

func findTgz(dirPath string) (string, error) {
	var tgzFile string
	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(path) == ".tgz" {
			tgzFile = path
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("no .tgz file was found in the directory: %s", dirPath)
	}
	if tgzFile == "" {
		return "", fmt.Errorf("no .tgz file was found in the directory: %s", dirPath)
	}
	return tgzFile, nil
}
