package helpers

func ContainsInSlice(hook string, slice []string) bool{
	for _, element :=range slice {
		if  element == hook {
			return true
		}
	}
	
	return false
}