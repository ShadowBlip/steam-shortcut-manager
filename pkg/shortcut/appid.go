package shortcut

import (
	"hash/crc32"
)

// https://github.com/boppreh/steamgrid/blob/9c8788db4f04613ecfb3e8fb36a0af02395e4593/games.go

/// Calculate an app id for a exe and app_name.
///
/// The app id is a 32-bit hash of the shortcut exe path and its app_name.
/// It is used to identify custom images for the shortcut.
// https://gaming.stackexchange.com/questions/386882/how-do-i-find-the-appid-for-a-non-steam-game-on-steam

func CalculateBPMID(exe, name string) uint64 {
	high32 := CalculateAppID(exe, name)
	full64 := (high32 << 32) | 0x02000000
	return uint64(full64)
}

func CalculateAppID(exe, name string) uint64 {
	combined := exe + name
	return uint64(crc32.ChecksumIEEE([]byte(combined))) | 0x80000000
}
