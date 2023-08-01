// platform utility
package platform

import "runtime"

type Arch string

const (
	ArchLinux Arch = "x64_linux"
	ArchWin32 Arch = "x64_win32"
)

func (a Arch) ServerBinaryName() string {
	switch a {
	case ArchLinux:
		return "altv-server"
	case ArchWin32:
		return "altv-server.exe"
	default:
		return ""
	}
}

func (a Arch) SharedLibExt() string {
	switch a {
	case ArchLinux:
		return "so"
	case ArchWin32:
		return "dll"
	default:
		return ""
	}
}

func (a Arch) String() string {
	return string(a)
}

func (a Arch) OS() string {
	switch a {
	case ArchLinux:
		return "linux"
	case ArchWin32:
		return "windows"
	default:
		return runtime.GOOS
	}
}

func Platform() Arch {
	switch os := runtime.GOOS; os {
	case "windows":
		return ArchWin32
	case "linux":
		return ArchLinux
	default:
		return ""
	}
}
