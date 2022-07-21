package hssm

import (
	"fmt"
	"io/ioutil"
	"syscall"
	"testing"
	"text/template"

	"github.com/Masterminds/sprig"
)

func createTempFile() (string, error) {
	file, err := ioutil.TempFile("", "")
	if err != nil {
		return "", fmt.Errorf(
			"unable to create temp file: %w", err,
		)
	}

	return file.Name(), nil
}

func TestExecuteTemplate(t *testing.T) {
	t.Parallel()

	templateContent := "example: {{ and true false }}"
	expectedOutput := "example: false"
	t.Logf("Template with content: %s , should out put a file with content: %s", templateContent, expectedOutput)

	templateFilePath, err := createTempFile()
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := syscall.Unlink(templateFilePath); err != nil {
			t.Errorf("unable to unlink template file: %s", err.Error())
		}
	}()

	if err := ioutil.WriteFile(templateFilePath, []byte(templateContent), 0o600); err != nil {
		t.Fatalf("Unable to write template file at path %s: %s", templateFilePath, err.Error())
	}

	content, err := ExecuteTemplate(templateFilePath, template.FuncMap{}, false)
	if content != expectedOutput {
		t.Errorf("Expected content \"%s\". Got \"%s\"", expectedOutput, content)
	}

	if err != nil {
		t.Errorf("unexpected error executing template: %s", err.Error())
	}
}

func TestWriteFile(t *testing.T) {
	t.Parallel()

	templateContent := "write_file_example: true"
	expectedOutput := "write_file_example: true"

	t.Logf("Template with content: %s , should out put a file with content: %s", templateContent, expectedOutput)

	templateFilePath, err := createTempFile()
	if err != nil {
		panic(err)
	}

	if err := WriteFile(templateFilePath, templateContent); err != nil {
		t.Fatalf("unable write template file at path %s: %s", templateFilePath, err.Error())
	}

	fileContent, err := ioutil.ReadFile(templateFilePath)
	if err != nil {
		panic(err)
	}

	if content := string(fileContent); content != expectedOutput {
		t.Errorf("Expected file with content \"%s\". Got \"%s\"", expectedOutput, content)
	}
}

func TestFailExecuteTemplate(t *testing.T) {
	t.Parallel()

	t.Logf("Non existing template should return \"no such file or directory\" error.")

	_, err := ExecuteTemplate("", template.FuncMap{}, false)
	if err == nil {
		t.Error("Should fail with \"no such file or directory\"")
	}
}

func TestSsmFunctionExistsInFuncMap(t *testing.T) {
	t.Parallel()

	t.Logf("\"ssm\" function should exist in function map.")

	funcMap := GetFuncMap("")
	keys := make([]string, len(funcMap))

	for k := range funcMap {
		keys = append(keys, k)
	}

	if _, exists := funcMap["ssm"]; !exists {
		t.Errorf("Expected \"ssm\" function in function map. Got the following functions: %s", keys)
	}
}

func TestSprigFunctionsExistInFuncMap(t *testing.T) {
	t.Parallel()

	t.Logf("\"quote\" function (from sprig) should exist in function map.")

	funcMap := GetFuncMap("")
	keys := make([]string, len(funcMap))

	for k := range funcMap {
		keys = append(keys, k)
	}

	if _, exists := funcMap["quote"]; !exists {
		t.Errorf("Expected \"quote\" function in function map. Got the following functions: %s", keys)
	}

	t.Logf("number of functions in function map minus custom functions should match those in sprig")

	if len(funcMap)-1 != len(sprig.GenericFuncMap()) {
		t.Errorf("Expected function map to include all sprig functions. Got the following functions: %s", keys)
	}
}

func TestResolveSSMParameter(t *testing.T) {
	t.Parallel()

	t.Logf("TODO")
}

func TestHandleOptions(t *testing.T) {
	t.Parallel()

	t.Logf("TODO")
}
