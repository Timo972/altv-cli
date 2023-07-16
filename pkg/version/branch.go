package version

type Branch string

const (
	BranchRelease Branch = "release"
	BranchRC      Branch = "rc"
	BranchDev     Branch = "dev"
)
