package cli

func CheckInput() error {
	// Validate paths
	if err := IsOfExtension("xml", AudioXMLPath, ".xml"); err != nil {
		return err
	}
	if err := IsOfExtension("oin", O2RSrcPath, ".o2r"); err != nil {
		return err
	}
	if err := IsDirectory("msrc", MusicSrcPath); err != nil {
		return err
	}

	return nil
}
