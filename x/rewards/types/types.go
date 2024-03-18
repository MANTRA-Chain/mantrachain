package types

type ClaimParams struct {
	StartClaimedSnapshotId *uint64
	EndClaimedSnapshotId   *uint64
	IsQuery                bool
	IsWithdraw             bool
}
