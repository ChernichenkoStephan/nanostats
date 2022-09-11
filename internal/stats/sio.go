package stats

import (
	"fmt"

	"bufio"
	"log"
	"os"
)

func OutputChats(names []string, counts []int, path string) error {
	parts := make([]string, 0)

	for i, name := range names {
		parts = append(parts, fmt.Sprintf("\"%s\": %d\n", name, counts[i]))
	}

	return output(parts, path)
}

// In process of fething channels are mmashed, do this function
// sorts output array to input sequence form
func SortLikeInputNames(stats []*Stats, names []string) {
	for _, n := range names {
		for i := 0; i < len(stats); i++ {
			if n != "@"+stats[i].Username {
				for j := 0; j < len(stats); j++ {
					if n == "@"+stats[j].Username {
						temp := stats[i]
						stats[i] = stats[j]
						stats[j] = temp
					}
				}
			}
		}
	}
}

func OutputStats(stats []Stats, path string) error {
	parts := make([]string, 0)

	participantsString := ""

	for _, s := range stats {
		parts = append(parts, (s.String() + "\n"))
		participantsString += fmt.Sprintf("%d ", s.Participants)
	}

	// Line for easy copying to excel
	parts = append(parts, participantsString+"\n")

	return output(parts, path)
}

func output(ss []string, path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			file, err := os.Create(path)
			if err != nil {
				return err
			}
			defer file.Close()
		}
	}

	file, err := os.OpenFile(path, os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	var bytesWritten int

	for _, s := range ss {
		newWritten, err := writer.WriteString(s)

		if err != nil {
			log.Fatal(err)
		}

		bytesWritten += newWritten
	}

	log.Printf("Bytes written: %d\n", bytesWritten)

	writer.Flush()

	return err
}
