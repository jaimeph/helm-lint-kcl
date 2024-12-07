package downloader

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/jaimeph/helm-lint-kcl/internal/logger"
)

type Downloader struct {
	chart     string
	version   string
	filePaths []string
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

	chartFileContents, err := d.ReadChart(chartDirPath, d.filePaths...)
	if err != nil {
		return nil, err
	}

	return chartFileContents, nil
}

// Obtenemos el chart, si no es una ruta local, lo descargamos con pull chart (comando de helm)
func (d *Downloader) GetChart(chart, version string) (string, error) {
	newChartDirPath := chart
	logger.Debugf("get chart %s %s", chart, version)
	if _, err := os.Stat(chart); os.IsNotExist(err) {
		newChartDirPath, err = d.pullChart(chart, version)
		if err != nil {
			return "", err
		}
	}
	return newChartDirPath, nil
}

// Descargamos el chart del repositorio de helm o de oci, añadimos la versión si se especifica
func (d *Downloader) pullChart(chart, version string) (string, error) {
	vParam := ""
	if len(version) > 0 {
		vParam = "--version " + version
	}
	cmd := fmt.Sprintf("helm pull %s -d %s %s", chart, os.TempDir(), vParam)
	logger.Debugf("exec: %s", cmd)
	output, err := exec.Command("bash", "-c", cmd).Output()
	if len(output) > 0 {
		logger.Infof("output: %s", output)
	}
	if err != nil {
		return "", fmt.Errorf("fail %s: %v", cmd, err)
	}
	return fmt.Sprintf("%s/%s-%s.tgz", os.TempDir(), path.Base(chart), version), nil
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
		content, err := os.ReadFile(filepath.Join(chartDirPath, filePath))
		if err != nil {
			return nil, err
		}
		contents[filePath] = content
	}
	return contents, nil
}

// Leemos los archivos del chart tgz
func (d *Downloader) ReadChartFileContentsTgz(chartDirPath string, filePaths ...string) (map[string][]byte, error) {
	contents := make(map[string][]byte, len(filePaths))

	file, err := os.Open(chartDirPath)
	if err != nil {
		return nil, fmt.Errorf("no se pudo abrir el archivo: %w", err)
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return nil, fmt.Errorf("error creando gzip reader: %w", err)
	}
	defer gzr.Close()

	tarReader := tar.NewReader(gzr)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error leyendo el archivo tar: %w", err)
		}

		parts := strings.Split(header.Name, "/")
		name := strings.Join(parts[1:], "/")

		for _, filePath := range filePaths {
			if name == filePath {
				logger.Debugf("extracting %s", name)
				var buffer bytes.Buffer
				if _, err := io.Copy(&buffer, tarReader); err != nil {
					return nil, fmt.Errorf("error copiando contenido del archivo: %w", err)
				}
				contents[filePath] = buffer.Bytes()
			}
		}
	}
	return contents, nil
}
