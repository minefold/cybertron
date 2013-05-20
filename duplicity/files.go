package duplicity

import (
	"bufio"
	"bytes"
	"regexp"
	"strings"
)

type Files []string

func DecodeFiles(info []byte) (Files, error) {
	scanner := bufio.NewScanner(bytes.NewReader(info))
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "Last full backup date") {
			break
		}
	}
	files := make(Files, 0)
	r := regexp.MustCompile(`\d{4} (.+)`)
	for scanner.Scan() {
    if matches := r.FindAllStringSubmatch(scanner.Text(), -1); matches != nil {
      files = append(files, matches[0][1])
    }
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return files, nil
}
