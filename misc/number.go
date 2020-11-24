package misc

import "crypto/rand"

// GenerateNumberFixedLength menggenerate angka dengan panjang yang tetap
// Kode ini di adopsi dari stackoverflow https://stackoverflow.com/a/61600241/6541319
func GenerateNumberFixedLength(length int) (string, error) {
	otpChars := "1234567890"
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	otpCharsLength := len(otpChars)
	for i := 0; i < length; i++ {
		buffer[i] = otpChars[int(buffer[i])%otpCharsLength]
	}

	return string(buffer), nil
}

// CountDigits menghitung pangjang suatu digit/nomor
func CountDigits(i int) (count int) {
	for i != 0 {
		i /= 10
		count = count + 1
	}
	return count
}
