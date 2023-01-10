package apis

func FirstNElementsOfSlice[Element any](slice []Element, n int) []Element {
	if len(slice) < n {
		return slice
	}
	return slice[:n]
}
