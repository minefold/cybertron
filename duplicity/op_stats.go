package duplicity

import (
	"bufio"
	"bytes"
	"strconv"
	"strings"
)

type OpStats struct {
	NewFiles int
}

func DecodeOpStats(info []byte) (*OpStats, error) {
	stats := new(OpStats)
	scanner := bufio.NewScanner(bytes.NewReader(info))
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "----") {
			break
		}
	}
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "----") {
			break
		}
		parts := strings.SplitN(scanner.Text(), " ", 2)
		key := parts[0]

		switch key {
		case "NewFiles":
			val := parts[1]

			stats.NewFiles, _ = strconv.Atoi(val)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return stats, nil
}
