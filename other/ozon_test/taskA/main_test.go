package taskA

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	oldStdin := fillStdin([]byte("1\n2\n2\n1\n2\n3\n2"))
	defer func() { os.Stdin = oldStdin }()
	exitCode := m.Run()
	os.Exit(exitCode)
}

func fillStdin(input []byte) *os.File {
	oldStdin := os.Stdin

	tmpFile, err := ioutil.TempFile(".", "input")
	if err != nil {
		log.Fatal(err)
	}

	if _, err := tmpFile.Write(input); err != nil {
		log.Fatal(err)
	}

	if _, err := tmpFile.Seek(0, 0); err != nil {
		log.Fatal(err)
	}

	defer func() {
		_ = os.Remove(tmpFile.Name())
	}()

	os.Stdin = tmpFile
	return oldStdin
}

func TestTaskA(t *testing.T) {
	scanner := bufio.NewScanner(os.Stdin)
	uniqueInputs := make(map[string]bool)
	for scanner.Scan() {
		number := scanner.Text()
		if _, ok := uniqueInputs[number]; !ok {
			uniqueInputs[number] = true
		} else {
			delete(uniqueInputs, number)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}

	for number := range uniqueInputs {
		fmt.Print(number)
	}
}
