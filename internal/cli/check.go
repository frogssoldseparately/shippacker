package cli

func CheckInput() error {
	if err := IsDirectory("msrc", MusicSrcPath); err != nil {
		return err
	}
	return nil
}
