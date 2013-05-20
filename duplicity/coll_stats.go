package duplicity

import (
	"bufio"
	"bytes"
	"errors"
	"strconv"
	"strings"
)

type CollStats struct {
	Backups []string
}

func DecodeCollStats(info []byte) (*CollStats, error) {
	stats := new(CollStats)
	scanner := bufio.NewScanner(bytes.NewReader(info))

	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "Number of contained backup sets:") {
			parts := strings.Split(scanner.Text(), ": ")
			if len(parts) < 2 {
				return nil, errors.New("Invalid: " + scanner.Text())
			}
			val, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, errors.New("Invalid: " + scanner.Text())
			}

			stats.Backups = make([]string, val)
		}
	}

	return stats, nil
}
