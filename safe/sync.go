package safe

import (
	"fmt"
)

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
