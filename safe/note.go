package safe

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"github.com/mitchellh/go-homedir"
)

const (
	defaultNoteEditor      = "vim"
	defaultNoteTmpFileName = "pick.note.tmp"
)

type notesManager struct {
	Notes map[string]note `json:"notes"`
	safe  *Safe
}

func (n *notesManager) Edit(name string) error {
	if len(name) == 0 {
		return errors.New("Empty note name specified")
	}
	note, exists := n.Notes[name]
	if !exists {
		note = NewEmptyNote(name)
	}
	if err := note.Edit(); err != nil {
		return err
	}
	n.Notes[name] = note
	fmt.Println("Note saved")
	return n.safe.save()
}

func (n *notesManager) List() map[string]note {
	return n.Notes
}

func (n *notesManager) Remove(name string) error {
	_, exists := n.Notes[name]
	if !exists {
		return fmt.Errorf("Note not found")
	}

	delete(n.Notes, name)

	return n.safe.save()
}

func newNotesManager(safe *Safe) *notesManager {
	return &notesManager{
		Notes: make(map[string]note),
		safe:  safe,
	}
}

type note struct {
	Name       string `json:"name"`
	Text       string `json:"text"`
	CreatedOn  int64  `json:"createdOn"`
	ModifiedOn int64  `json:"modifiedOn"`
}

func (n *note) Edit() error {
	text, err := editorReadText(n.Text)
	if err != nil {
		return err
	}
	n.Text = text
	n.ModifiedOn = time.Now().Unix()
	return nil
}

func NewEmptyNote(name string) note {
	ts := time.Now().Unix()
	return note{
		Name:       name,
		Text:       "",
		CreatedOn:  ts,
		ModifiedOn: ts,
	}
}

func editorReadText(existingText string) (string, error) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = defaultNoteEditor
	}
	editorPath, err := exec.LookPath(editor)
	if err != nil {
		return "", err
	}
	homeDir, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	tmpFile, err := ioutil.TempFile(homeDir, defaultNoteTmpFileName)
	if err != nil {
		return "", err
	}
	if existingText != "" {
		if _, err := tmpFile.WriteString(existingText); err != nil {
			return "", err
		}
	}
	cmd := exec.Command(editorPath, tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return "", err
	}
	if err := cmd.Wait(); err != nil {
		return "", err
	}
	note, err := ioutil.ReadFile(tmpFile.Name())
	if err != nil {
		return "", err
	}
	if err := os.Remove(tmpFile.Name()); err != nil {
		return "", err
	}
	return string(note), nil
}
