package version

type Branch string

const (
	BranchRelease Branch = "release"
	BranchRC      Branch = "rc"
	BranchDev     Branch = "dev"
)

func (b Branch) String() string {
	return string(b)
}
