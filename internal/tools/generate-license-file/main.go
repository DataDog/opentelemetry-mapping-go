package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"

	"go.uber.org/multierr"
)

const licensesCSV = "LICENSE-3rdparty.csv"

var (
	// modules on this repository.
	modules = []string{
		"pkg/quantile",
		"pkg/otlp/attributes",
		"pkg/otlp/metrics",
		"pkg/internal/sketchtest",
		"pkg/inframetadata",
		"pkg/inframetadata/gohai/internal/gohaitest",
	}
)

func main() {
	f, err := os.Create(licensesCSV)
	if err != nil {
		log.Fatalf("Failed to open %q file: %s\n", licensesCSV, err)
	}

	w := csv.NewWriter(f)
	if err := w.Write([]string{"Component", "Origin", "License", "Copyright"}); err != nil {
		log.Fatalln("error writing record to csv:", err)
	}

	for _, module := range modules {
		deps, err := findDependenciesOf(module)
		if err != nil {
			log.Fatalln(err)
		}

		for _, p := range deps {
			if err := w.Write(p.Record()); err != nil {
				log.Fatalln("error writing record to csv:", err)
			}
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}

type Package struct {
	Component        string
	Origin           string
	License          string
	CopyrightNotices []string
}

func (p Package) Record() []string {
	return []string{p.Component, p.Origin, p.License, strings.Join(p.CopyrightNotices, " | ")}
}

var _ sort.Interface = (*Packages)(nil)

type Packages []Package

func (p Packages) Len() int           { return len(p) }
func (p Packages) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p Packages) Less(i, j int) bool { return p[i].Origin < p[j].Origin }

// findDependenciesOf a given module given by its folder path.
func findDependenciesOf(module string) (packages []Package, err error) {
	cwd, err := os.Getwd()
	if err != nil {
		return
	}

	if err = os.Chdir(module); err != nil {
		return nil, fmt.Errorf("failed to changed directory to %q: %w", module, err)
	}
	// restore directory after exit
	defer func() { err = multierr.Append(err, os.Chdir(cwd)) }()

	// wwhrd needs vendored dependencies.
	cmd := exec.Command("go", "mod", "vendor")
	if err = cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to run 'go mod vendor': %w", err)
	}
	// remove vendored dependencies after exit
	defer func() { err = multierr.Append(err, os.RemoveAll("vendor/")) }()

	cmd = exec.Command("wwhrd", "list", "--no-color")
	var out bytes.Buffer
	cmd.Stderr = &out
	if err = cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to run 'wwhrd list': %w", err)
	}

	// Parse wwhrd output
	scanner := bufio.NewScanner(&out)
	const foundLicense = "msg=\"Found License\""
	for scanner.Scan() {
		line := scanner.Text()
		index := strings.Index(line, foundLicense)
		if index == -1 {
			continue
		}

		var pkg Package
		pkg.Component = module
		parts := strings.Split(line[index+len(foundLicense):], " ")
		for _, part := range parts {
			switch {
			case strings.HasPrefix(part, "license="):
				pkg.License = part[len("license="):]
			case strings.HasPrefix(part, "package="):
				pkg.Origin = part[len("package="):]
			}
		}

		if pkg.CopyrightNotices, err = findCopyrightNotices(pkg.Origin); len(pkg.CopyrightNotices) == 0 {
			if err == nil {
				err = fmt.Errorf("could not find copyright notice for %q", pkg.Origin)
			}
			return
		}

		packages = append(packages, pkg)
	}

	sort.Sort(Packages(packages))
	return
}
