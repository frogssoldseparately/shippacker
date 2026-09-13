package cli

func CheckInput() error {
	return IsDirectory("msrc", MusicSrcPath)
}
