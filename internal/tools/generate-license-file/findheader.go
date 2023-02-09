package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"go.uber.org/multierr"
)

var (
	copyrightLocations = []string{
		"LICENSE",
		"license.md",
		"LICENSE.md",
		"LICENSE.txt",
		"License.txt",
		"COPYING",
		"NOTICE",
		"README",
		"README.md",
		"README.mdown",
		"README.markdown",
		"COPYRIGHT",
		"COPYRIGHT.txt",
	}
	authorLocations = []string{
		"AUTHORS",
		"AUTHORS.md",
		"CONTRIBUTORS",
	}
)

// findCopyrightNotices on a given dependency.
// The dependency is passed by its full path, e.g. github.com/DataDog/sketches-go.
func findCopyrightNotices(origin string) (copyrightHeaders []string, err error) {
	if headers, ok := globalOverrides.CopyrightNotice(origin); ok {
		return headers, nil
	}

	if strings.Contains(origin, "/") {
		parentHeaders, err := findCopyrightNotices(origin[:strings.LastIndex(origin, "/")])
		if err != nil {
			return nil, err
		}
		copyrightHeaders = append(copyrightHeaders, parentHeaders...)
	}

	pkgDir := path.Join("vendor", origin)

	for _, filename := range copyrightLocations {
		var lines []string
		lines, err = mapLines(path.Join(pkgDir, filename), getCopyrightNotice)
		if err != nil {
			return
		}
		copyrightHeaders = append(copyrightHeaders, lines...)

	}

	for _, filename := range authorLocations {
		var lines []string
		lines, err = mapLines(path.Join(pkgDir, filename), getAuthors)
		if err != nil {
			return
		}
		copyrightHeaders = append(copyrightHeaders, lines...)
	}

	return
}

// mapLines maps lines of a file given by its full path.
// The mapping function fn may return ok=false to indicate that a line should be skipped
func mapLines(fullPath string, fn func(line string) (string, bool)) ([]string, error) {
	if _, err := os.Stat(fullPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// file does not exist, nothing to do.
			return nil, nil
		} else {
			// true error. bubble up.
			return nil, fmt.Errorf("failed to stat %q: %w", fullPath, err)
		}
	}

	file, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open %q: %w", fullPath, err)
	}
	defer func() { err = multierr.Append(err, file.Close()) }()

	var notices []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if notice, ok := fn(scanner.Text()); ok {
			notices = append(notices, notice)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan %q: %w", fullPath, err)
	}

	return notices, nil
}

var (
	copyrightHeaderRegexp = regexp.MustCompile(`(?i)copyright\s+(?:©|\(c\)\s+)?(?:(?:[0-9 ,-]|present)+\s+)?(?:by\s+)?(.*)`)

	copyrightIgnoreRegexp = []*regexp.Regexp{
		regexp.MustCompile(`(?i)copyright(:? and license)?$`),
		regexp.MustCompile(`(?i)copyright (:?holder|owner|notice|license|statement)`),
		regexp.MustCompile(`Copyright & License -`),
		regexp.MustCompile(`(?i)copyright .yyyy. .name of copyright owner.`),
		regexp.MustCompile(`(?i)copyright .yyyy. .name of copyright owner.`),
	}
)

// getCopyrightNotice from a given line on LICENSE-like file.
func getCopyrightNotice(line string) (notice string, ok bool) {
	matches := copyrightHeaderRegexp.FindStringSubmatch(line)
	if len(matches) == 0 {
		return
	}

	notice = matches[0]
	var shouldIgnore bool
	for _, reg := range copyrightIgnoreRegexp {
		if reg.MatchString(notice) {
			shouldIgnore = true
			break
		}
	}
	if shouldIgnore {
		return
	}

	return strings.TrimSuffix(strings.TrimSpace(notice), "."), true
}

// getAuthors from a given line on an AUTHORS-like file.
func getAuthors(line string) (author string, ok bool) {
	line = strings.TrimSpace(line)
	if len(line) == 0 || line[0] == '#' {
		return
	}
	return line, true
}
