package app

import "runtime/debug"

var readBuildInfo = debug.ReadBuildInfo

func Version(version string) string {
	if version != "" {
		return version
	}

	if info, ok := readBuildInfo(); ok {
		if info.Main.Version != "" && info.Main.Version != "(devel)" {
			return info.Main.Version
		}
	}

	return "dev"
}
