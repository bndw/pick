package safe

import "testing"

const (
	accountName = "name2"
	initialUser = "username1"
	initialPswd = "password1"
	updatedUser = "username2"
	updatedPswd = "password2"
)

func TestSyncSameHistory(t *testing.T) {
	safe1, err := createTestSafe()
	if err != nil {
		t.Fatal(err)
	}
	safe2, err := createTestSafe()
	if err != nil {
		t.Fatal(err)
	}
	defer removeTestSafe()

	if acc1, err := safe1.Add(accountName, initialUser, initialPswd); err != nil {
		t.Fatal(err)
	} else if len(acc1.History) != 0 {
		t.Fatal("safe1 account should not have a history after creation")
	}

	if _, err := safe2.Get(accountName); err == nil {
		t.Fatal("safe2 should not have this account yet")
	}

	// Import accounts from safe1 into safe2
	if err := safe2.SyncWith(safe1); err != nil {
		t.Fatal(err)
	}

	if acc2, err := safe2.Edit(accountName, updatedUser, updatedPswd); err != nil {
		t.Fatal(err)
	} else {
		// Hack to update ModifiedOn
		acc2.ModifiedOn++
		safe2.Accounts[accountName] = *acc2
	}

	if err := safe1.SyncWith(safe2); err != nil {
		t.Fatal(err)
	}

	if acc1, err := safe1.Get(accountName); err != nil {
		t.Fatal(err)
	} else if len(acc1.History) == 0 {
		t.Fatal("safe1 account should have a history after non-empty sync")
	}
}

func TestSyncDifferentHistory(t *testing.T) {
	safe1, err := createTestSafe()
	if err != nil {
		t.Fatal(err)
	}
	safe2, err := createTestSafe()
	if err != nil {
		t.Fatal(err)
	}
	defer removeTestSafe()

	if acc1, err := safe1.Add(accountName, initialUser, initialPswd); err != nil {
		t.Fatal(err)
	} else if len(acc1.History) != 0 {
		t.Fatal("safe1 account should not have a history after creation")
	}

	if acc2, err := safe2.Add(accountName, initialUser, initialPswd); err != nil {
		t.Fatal(err)
	} else if len(acc2.History) != 0 {
		t.Fatal("safe2 account should not have a history after creation")
	} else {
		// Hack to update CreatedOn, this prevents to sync this account with acc1
		acc2.CreatedOn++
		safe2.Accounts[accountName] = *acc2
	}

	if err := safe2.SyncWith(safe1); err != nil {
		t.Fatal(err)
	}

	if acc2, err := safe2.Edit(accountName, updatedUser, updatedPswd); err != nil {
		t.Fatal(err)
	} else {
		// Hack to update ModifiedOn
		acc2.ModifiedOn++
		safe2.Accounts[accountName] = *acc2
	}

	if err := safe1.SyncWith(safe2); err != nil {
		t.Fatal(err)
	}

	if acc1, err := safe1.Get(accountName); err != nil {
		t.Fatal(err)
	} else if len(acc1.History) != 0 {
		t.Fatal("safe1 account should still not have a history, as it should not have been synced")
	}
}
