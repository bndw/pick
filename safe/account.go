package safe

import (
	"fmt"
	"sort"
	"time"
)

type Account struct {
	Username   string         `json:"username"`
	Password   string         `json:"password"`
	CreatedOn  int64          `json:"createdOn"`
	ModifiedOn int64          `json:"modifiedOn"`
	History    accountHistory `json:"history,omitempty"`
}

func (acc *Account) Update(cb func(*Account)) {
	cb(acc)
	acc.ModifiedOn = time.Now().Unix()
}

func (acc *Account) SyncWith(otherAccount *Account, name string) (bool, error) {
	if acc.CreatedOn != otherAccount.CreatedOn {
		// Apparently not the same account
		// TODO(leon): Implement unique ID for an account
		return false, fmt.Errorf("Accounts '%s' differ in creation date, skipping", name)
	}
	if acc.ModifiedOn < otherAccount.ModifiedOn {
		// Other account is newer, update ourself
		// Defer backup creation
		defer func(acc Account) {
			acc.History = append(acc.History, acc)
		}(*acc)
		acc.Username = otherAccount.Username
		acc.Password = otherAccount.Password
		// Sync history
		acc.syncHistory(otherAccount.History)
		// Update ModifiedOn
		acc.ModifiedOn = otherAccount.ModifiedOn
		return true, nil
	}
	return false, nil
}

type accountHistory []Account

func (a accountHistory) Len() int {
	return len(a)
}

func (a accountHistory) Less(i, j int) bool {
	return a[i].ModifiedOn < a[j].ModifiedOn
}

func (a accountHistory) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (acc *Account) syncHistory(otherHistory accountHistory) {
	for _, otherAccount := range otherHistory {
		found := false
		for _, localAccount := range acc.History {
			if localAccount.ModifiedOn == otherAccount.ModifiedOn {
				// We already have this account in our history
				found = true
				break
			}
		}
		if !found {
			acc.History = append(acc.History, otherAccount)
		}
	}
	// Sort our history to preserve order
	historySorted := make(accountHistory, 0, len(acc.History))
	for _, account := range acc.History {
		historySorted = append(historySorted, account)
	}
	sort.Sort(historySorted)
	acc.History = historySorted
}

func NewAccount(name, username, password string) *Account {
	ts := time.Now().Unix()
	return &Account{
		Username:   username,
		Password:   password,
		CreatedOn:  ts,
		ModifiedOn: ts,
		History:    make(accountHistory, 0),
	}
}
