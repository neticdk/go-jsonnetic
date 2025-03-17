package utils

import (
	"fmt"
	"github.com/neticdk/go-jsonnetic/pkg/jsonnetic/native"
	"os"
	"path/filepath"
	"sort"
)

// PrettyFuncList returns a pretty list of all available functions.
func PrettyFuncList() string {
	funcList := ""
	for _, f := range native.Funcs() {
		funcList += fmt.Sprintf("- %s\n", f.Name)
	}
	return funcList
}

func WriteMultiOutputFiles(output map[string]string, outputDir, outputFile string, createDirs bool) (err error) {
	// If multiple file output is used, then iterate over each string from
	// the sequence of strings returned by jsonnet_evaluate_snippet_multi,
	// construct pairs of filename and content, and write each output file.

	var manifest *os.File

	if outputFile == "" {
		manifest = os.Stdout
	} else {
		manifest, err = os.Create(outputFile)
		if err != nil {
			return err
		}
		defer func() {
			if ferr := manifest.Close(); ferr != nil {
				err = ferr
			}
		}()
	}

	// Iterate through the map in order.
	keys := make([]string, 0, len(output))
	for k := range output {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		newContent := output[key]
		filename := outputDir + key

		_, err := manifest.WriteString(filename)
		if err != nil {
			return err
		}

		_, err = manifest.WriteString("\n")
		if err != nil {
			return err
		}

		if _, err := os.Stat(filename); !os.IsNotExist(err) {
			existingContent, err := os.ReadFile(filename)
			if err != nil {
				return err
			}
			if string(existingContent) == newContent {
				// Do not bump the timestamp on the file if its content is
				// the same. This may trigger other tools (e.g. make) to do
				// unnecessary work.
				continue
			}
		}
		if createDirs {
			if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
				return err
			}
		}

		err = os.WriteFile(filename, []byte(newContent), 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

// WriteOutputFile writes the output to the given file, creating directories
// if requested, and printing to stdout instead if the outputFile is "".
func WriteOutputFile(output string, outputFile string, createDirs bool) (err error) {
	if outputFile == "" {
		fmt.Print(output)
		return nil
	}

	if createDirs {
		if err := os.MkdirAll(filepath.Dir(outputFile), 0755); err != nil {
			return err
		}
	}

	f, createErr := os.Create(outputFile)
	if createErr != nil {
		return createErr
	}
	defer func() {
		if ferr := f.Close(); ferr != nil {
			err = ferr
		}
	}()

	_, err = f.WriteString(output)
	return err
}
