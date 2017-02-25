package safe

import (
	"fmt"
)

// SyncWith syncs the current safe with another safe if they're branching off the same
// safe, i.e. are not completely different safes. It simply imports non-existing
// accounts into the current safe and updates existing accounts if the other safe has
// a more recent version of the account. A backup is made when updating existing accounts
func (s *Safe) SyncWith(otherSafe *Safe) error {
	localAccounts := s.Accounts
	remoteAccounts := otherSafe.Accounts

	for name, remoteAccount := range remoteAccounts {
		// If we don't have an remote account yet, simply add it
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
	// TODO(leon): Should we also sync in the reverse direction, i.e. local -> remote?
	return s.save()
}
