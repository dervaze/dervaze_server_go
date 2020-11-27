package lang

func CompareStringSlices(slice1, slice2 []string) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	for i, v := range slice1 {
		if slice2[i] != v {
			return false
		}
	}
	return true
}

// TFstring returns ifTrue or ifFalse according to condition
func TFstring(condition bool, ifTrue, ifFalse string) string {
	if condition {
		return ifTrue
	} else {
		return ifFalse
	}
}

// TFString returns ifTrue or ifFalse according to condition
func TFint(condition bool, ifTrue, ifFalse int) int {
	if condition {
		return ifTrue
	} else {
		return ifFalse
	}
}
