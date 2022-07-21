package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	hssm "github.com/codacy/helm-ssm/internal"
	"github.com/spf13/cobra"
)

var (
	valueFiles valueFilesList
	targetDir  string
	profile    string
	verbose    bool
	dryRun     bool
)

type valueFilesList []string

func (v *valueFilesList) String() string {
	return fmt.Sprint(*v)
}

func (v *valueFilesList) Type() string {
	return "valueFilesList"
}

func (v *valueFilesList) Set(value string) error {
	for _, filePath := range strings.Split(value, ",") {
		*v = append(*v, filePath)
	}

	return nil
}

func main() {
	cmd := &cobra.Command{
		Use:   "ssm [flags]",
		Short: "",
		RunE:  run,
	}

	f := cmd.Flags()
	f.VarP(&valueFiles, "values", "f", "specify values in a YAML file (can specify multiple)")
	f.BoolVarP(&verbose, "verbose", "v", false, "show the computed YAML values file/s")
	f.BoolVarP(&dryRun, "dry-run", "d", false, "doesn't replace the file content")
	f.StringVarP(&targetDir, "target-dir", "o", "", "dir to output content")
	f.StringVarP(&profile, "profile", "p", "", "aws profile to fetch the ssm parameters")

	if err := cmd.MarkFlagRequired("values"); err != nil {
		log.Fatalf("unable to mark cobra flag required: %s", err.Error())
		os.Exit(1)
	}

	if err := cmd.Execute(); err != nil {
		log.Fatalf("unable to run cobra command: %s", err.Error())
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	funcMap := hssm.GetFuncMap(profile)
	for _, filePath := range valueFiles {
		content, err := hssm.ExecuteTemplate(filePath, funcMap, verbose)
		if err != nil {
			return fmt.Errorf("unable to execute template: %w", err)
		}

		if !dryRun {
			if err := write(filePath, targetDir, content); err != nil {
				return fmt.Errorf("unable to write output: %w", err)
			}
		}
	}

	return nil
}

func write(filePath string, targetDir string, content string) error {
	if targetDir != "" {
		fileName := filepath.Base(filePath)

		if err := hssm.WriteFileD(fileName, targetDir, content); err != nil {
			return fmt.Errorf(
				"unable to dump content of file to %s/%s: %w",
				targetDir,
				fileName,
				err,
			)
		}
	}

	if err := hssm.WriteFile(filePath, content); err != nil {
		return fmt.Errorf(
			"unable to write file %s: %w",
			filePath,
			err,
		)
	}

	return nil
}
