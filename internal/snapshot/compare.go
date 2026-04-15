package snapshot

// ChangeKind describes how a key changed between two snapshots.
type ChangeKind string

const (
	ChangeAdded   ChangeKind = "added"
	ChangeRemoved ChangeKind = "removed"
	ChangeUpdated ChangeKind = "updated"
)

// Change represents a single key-level difference between snapshots.
type Change struct {
	Key  string
	Kind ChangeKind
}

// Diff compares a previous snapshot to the current one and returns a list
// of changes. If prev is nil every key in current is considered added.
func Diff(prev *Snapshot, current Snapshot) []Change {
	var changes []Change

	prevMap := make(map[string]string)
	if prev != nil {
		for _, e := range prev.Entries {
			prevMap[e.Key] = e.ValueHash
		}
	}

	currMap := make(map[string]string)
	for _, e := range current.Entries {
		currMap[e.Key] = e.ValueHash
	}

	for k, currHash := range currMap {
		if prevHash, ok := prevMap[k]; !ok {
			changes = append(changes, Change{Key: k, Kind: ChangeAdded})
		} else if prevHash != currHash {
			changes = append(changes, Change{Key: k, Kind: ChangeUpdated})
		}
	}

	for k := range prevMap {
		if _, ok := currMap[k]; !ok {
			changes = append(changes, Change{Key: k, Kind: ChangeRemoved})
		}
	}

	return changes
}
