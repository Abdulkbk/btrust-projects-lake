package main

func byteFromHexChar(hexChar byte) byte {
	switch {
	case hexChar >= '0' && hexChar <= '9':
		return hexChar - '0'
	case hexChar >= 'a' && hexChar <= 'f':
		return hexChar - 'a' + 10
	case hexChar >= 'A' && hexChar <= 'F':
		return hexChar - 'A' + 10
	default:
		return 0
	}
}
