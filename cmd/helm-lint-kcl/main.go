package main

import (
	"fmt"
	"os"

	"github.com/jaimeph/helm-lint-kcl/internal/downloader"
	"github.com/jaimeph/helm-lint-kcl/internal/logger"
	"github.com/jaimeph/helm-lint-kcl/internal/merger"
	"github.com/jaimeph/helm-lint-kcl/internal/validator"
	"github.com/spf13/cobra"
)

const (
	valuesFilePath  = "values.yaml"
	schemasFilePath = "schemas.k"
)

var (
	debug      bool
	showValues bool
	version    string
	values     []string
	sets       []string
)

func main() {
	rootCmd := &cobra.Command{
		Use:           "helm lint-kcl [NAME] [CHART] [flags]",
		Short:         "Validate Helm values using KCL schemas",
		Long:          "A Helm plugin to validate values.yaml using KCL schemas defined in schemas.k.",
		Args:          mainArgs,
		RunE:          mainRunE,
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	rootCmd.Flags().StringVarP(&version, "version", "v", "", "chart version")
	rootCmd.Flags().BoolVar(&debug, "debug", false, "enable debug")
	rootCmd.Flags().BoolVar(&showValues, "show-values", false, "show values")
	rootCmd.Flags().StringSliceVarP(&values, "values", "f", []string{}, "specify values in a YAML file or a URL (can specify multiple)")
	rootCmd.Flags().StringSliceVar(&sets, "set", []string{}, "set values on the command line (can specify multiple or separate values with commas: key1=val1,key2=val2)")

	if err := rootCmd.Execute(); err != nil {
		logger.Errorf("Execution error: %s", err)
		os.Exit(1)
	}
}

func mainArgs(cmd *cobra.Command, args []string) error {
	logger.SetLevelDebug(debug)
	if len(args) != 2 {
		return fmt.Errorf("name and chart are required")
	}
	return cobra.OnlyValidArgs(cmd, args)
}

func mainRunE(cmd *cobra.Command, args []string) error {
	logger.SetLevelDebug(debug)

	_, chart := args[0], args[1]

	filePaths := []string{valuesFilePath, schemasFilePath}

	d := downloader.New(chart, version)
	contentFilePaths, err := d.GetFilesContents(filePaths...)
	if err != nil {
		return err
	}

	m, err := merger.New(contentFilePaths[valuesFilePath])
	if err != nil {
		return err
	}
	err = m.Values(values)
	if err != nil {
		return err
	}
	err = m.Sets(sets)
	if err != nil {
		return err
	}
	mValues, err := m.Merged()
	if err != nil {
		return err
	}

	if showValues {
		logger.Infof("Values:\n%s", mValues)
	}

	v := validator.New(mValues, contentFilePaths[schemasFilePath])
	if err := v.Validate(); err != nil {
		return err
	} else {
		logger.Info("✅ Values validation successful!")
	}

	return nil
}
