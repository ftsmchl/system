package helper

import ()

//checks if a string is a valid hex string
func IsHex(str string) bool {

	bytes := []byte(str)

	if bytes[0] == 48 && bytes[1] == 120 {
		//start with 0x, do nothing
	} else {
		return false
	}

	for i := 2; i < len(bytes); i++ {
		//checks if character is 0-9 or ABCDEF or abcdef
		if (48 <= bytes[i] && bytes[i] <= 57) || (65 <= bytes[i] && bytes[i] <= 70) || (97 <= bytes[i] && bytes[i] <= 102) {
			//char is hexademical do nothing
		} else {
			return false
		}
	}

	return true
}
