package safe

import (
	"fmt"
)

// SyncWith syncs the current safe with another safe if they're branching off the same
// safe, i.e. are not completely different safes. It simply imports non-existing
// accounts & notes into the current safe and updates existing accounts & notes if the
// other safe has a more recent version of an account / note. A backup is made when
// existing accounts / notes are updated.
func (s *Safe) SyncWith(otherSafe *Safe) error {
	syncAccounts(s.Accounts, otherSafe.Accounts)
	syncNotes(s.Notes.Notes, otherSafe.Notes.Notes)

	// TODO(leon): Should we also sync in the reverse direction, i.e. local -> remote?
	return s.save()
}

func syncAccounts(localAccounts, remoteAccounts map[string]Account) {
	for name, remoteAccount := range remoteAccounts {
		// If we don't have a remote account yet, simply add it
		if localAccount, ok := localAccounts[name]; !ok {
			fmt.Printf("Importing new account '%s'..\n", name)
			localAccounts[name] = remoteAccount
		} else {
			// We already have this account, now sync current state & history
			fmt.Printf("Syncing '%s'.. ", name)
			synced, err := localAccount.SyncWith(&remoteAccount, name)
			if err != nil {
				fmt.Printf("error: %s\n", err)
				continue
			}
			if !synced {
				fmt.Printf("account already up-to-date\n")
			} else {
				fmt.Printf("done\n")
				// Save it
				localAccounts[name] = localAccount
			}
		}
	}
}

func syncNotes(localNotes, remoteNotes map[string]note) {
	for name, remoteNote := range remoteNotes {
		// If we don't have a remote note yet, simply add it
		if localNote, ok := localNotes[name]; !ok {
			fmt.Printf("Importing new note '%s'..\n", name)
			localNotes[name] = remoteNote
		} else {
			// We already have this note, now sync current state & history
			fmt.Printf("Syncing '%s'.. ", name)
			synced, err := localNote.SyncWith(&remoteNote, name)
			if err != nil {
				fmt.Printf("error: %s\n", err)
				continue
			}
			if !synced {
				fmt.Printf("note already up-to-date\n")
			} else {
				fmt.Printf("done\n")
				// Save it
				localNotes[name] = localNote
			}
		}
	}
}
