package safe

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	homedir "github.com/mitchellh/go-homedir"
	"golang.leonklingele.de/securetemp"
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
	nt, exists := n.Notes[name]
	if !exists {
		nt = NewEmptyNote(name)
	}
	if err := nt.EditInEditor(); err != nil {
		return err
	}
	n.Notes[name] = nt
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

func (n *note) update(text string) {
	n.Text = text
	n.ModifiedOn = time.Now().Unix()
}

func (n *note) EditInEditor() error {
	text, err := editorReadText(n.Text)
	if err != nil {
		return err
	}
	n.update(text)
	return nil
}

func (n *note) SyncWith(otherNote *note, name string) (bool, error) {
	if n.CreatedOn != otherNote.CreatedOn {
		// Apparently not the same note
		// TODO(leon): Implement unique ID for a note
		return false, fmt.Errorf("Notes '%s' differ in creation date, skipping", name)
	}
	if n.ModifiedOn < otherNote.ModifiedOn {
		// Other note is newer, update ourself
		n.Text = otherNote.Text
		// Update ModifiedOn
		n.ModifiedOn = otherNote.ModifiedOn
		return true, nil
	}
	return false, nil
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

func editorPath() (string, error) {
	editor := os.Getenv("VISUAL")
	if editor == "" {
		editor = os.Getenv("EDITOR")
		if editor == "" {
			editor = defaultNoteEditor
		}
	}
	return exec.LookPath(editor)
}

func editorReadText(existingText string) (string, error) {
	editorPath, err := editorPath()
	if err != nil {
		return "", err
	}
	homeDir, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	tmpFile, cleanupFunc, err := securetemp.TempFile(8 * securetemp.SizeMB)
	if err != nil {
		tmpFile, err = ioutil.TempFile(homeDir, defaultNoteTmpFileName)
		if err != nil {
			return "", err
		}
		cleanupFunc = func() {
			os.Remove(tmpFile.Name())
		}
	}
	defer cleanupFunc()
	if existingText != "" {
		if _, err := tmpFile.WriteString(existingText); err != nil { // nolint: vetshadow
			return "", err
		}
	}
	cmd := exec.Command(editorPath, tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil { // nolint: vetshadow
		return "", err
	}
	if err := cmd.Wait(); err != nil { // nolint: vetshadow
		return "", err
	}
	note, err := ioutil.ReadFile(tmpFile.Name())
	if err != nil {
		return "", err
	}
	return string(note), nil
}
