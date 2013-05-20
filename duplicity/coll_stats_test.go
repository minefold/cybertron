package duplicity

import (
	"testing"
)

func TestNoBackups(t *testing.T) {
	input := `Local and Remote metadata are synchronized, no sync needed.
Last full backup date: none
Collection Status
-----------------
Connecting with backend: BotoBackend
Archive dir: /Users/dave/.cache/duplicity/0187ce5469b85ab00178266045db34f3

Found 0 secondary backup chains.
No backup chains with active signatures found
No orphaned or incomplete backup sets found.`

	stats, _ := DecodeCollStats([]byte(input))
	if len(stats.Backups) != 0 {
		t.Fatalf("Expected 0 was %d\n", len(stats.Backups))
	}
}


func TestWithBackups(t *testing.T) {
	input := `Local and Remote metadata are synchronized, no sync needed.
Last full backup date: Thu May 16 15:50:21 2013
Collection Status
-----------------
Connecting with backend: BotoBackend
Archive dir: /Users/dave/.cache/duplicity/e1c6d10a9cc684be3228c481bad6e83f

Found 0 secondary backup chains.

Found primary backup chain with matching signature chain:
-------------------------
Chain start time: Thu May 16 15:50:21 2013
Chain end time: Thu May 16 15:50:42 2013
Number of contained backup sets: 3
Total number of contained volumes: 3
 Type of backup set:                            Time:      Num volumes:
                Full         Thu May 16 15:50:21 2013                 1
         Incremental         Thu May 16 15:50:33 2013                 1
         Incremental         Thu May 16 15:50:42 2013                 1
-------------------------
No orphaned or incomplete backup sets found.`

	stats, _ := DecodeCollStats([]byte(input))
	if len(stats.Backups) != 3 {
		t.Fatalf("Expected 3 was %d\n", len(stats.Backups))
	}
}
