package microservice

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

var RE_LINE_EMPTY = regexp.MustCompile(`^\s*(?:#.*)?$`)
var RE_LINE_REMOVE_COMMENT = regexp.MustCompile(`#.*$`)
var RE_LINE_SPLIT = regexp.MustCompile(`^(?P<user>\S+)\s+(?P<password>.*)`)

type UserEntry struct {
	username     string
	passwordhash []byte
}

func InitUserEntry(e *UserEntry, username string, password string) {
	e.username = username

	if strings.HasPrefix(password, "$") {
		e.passwordhash = []byte(password)
		return
	}

	passwordhash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		log.Panic(fmt.Printf("Could not hash password, error was '%s'!", err.Error()))
	}

	log.Printf("XXX [%s][%s][%s] XXX\n", username, password, string(passwordhash))
	e.passwordhash = passwordhash
}

func UserEntryFromString(line string) *UserEntry {

	if RE_LINE_EMPTY.MatchString(line) {
		return nil
	}

	line = RE_LINE_REMOVE_COMMENT.ReplaceAllLiteralString(line, "")
	split := RE_LINE_SPLIT.FindStringSubmatch(line)

	if len(split) != 3 {
		log.Panic(fmt.Printf("Could not split password line! [%s] %d", line, len(split)))
	}

	var e UserEntry

	InitUserEntry(&e, split[1], split[2])
	return &e
}

func (e *UserEntry) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword(e.passwordhash, []byte(password))
	return err == nil
}

func UserEntriesFromFile(filename string) map[string]UserEntry {
	file, err := os.Open(filename)

	if err != nil {
		log.Panic("Failed to open file '%s': %s", filename, err)
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	var line string
	var userEntries map[string]UserEntry
	userEntries = make(map[string]UserEntry)

	for {
		line, err = reader.ReadString('\n')

		if err != nil && err != io.EOF {
			break
		}
		entry := UserEntryFromString(line)

		if entry != nil {
			userEntries[entry.username] = *entry
		}

		if err != nil {
			break
		}
	}

	if err != io.EOF {
		log.Panic("Failed to read file '%s': %s", filename, err)
	}

	return userEntries
}
