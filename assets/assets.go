package assets

//GetAsset retrives []byte of file and converts it to a string
func GetAsset(path string) (string, error) {
	//todo- think through if it's worth factoring out into an interface for mocking?
	data, err := Asset(path)
	if err != nil {
		println("Error reading file ", path)
		return "", err
	}
	s := string(data[:])
	return s, nil
}
