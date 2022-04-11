package generators

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"image"
	"image/jpeg"
	"liokoredu/pkg/constants"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

/* Converts dataURL to file and saves it. Returns file path. Only jpg and png supported
Usage example:

path, err := dataURLToFile("wolchara", newData.AvatarURL, 500)
if err != nil {
	fmt.Println(err.Error())
} else {
	fmt.Println(path) // wolchara.jpg
}
*/

func DataURLToFile(path string, dataURL string, maxSizeKB int) (string, error) {
	if dataURL == "" {
		return "", nil
	}

	splittedURL := strings.Split(dataURL, ",")
	if len(splittedURL) != 2 {
		return "", errors.New("incorrect data url")
	}

	meta := splittedURL[0]
	if !strings.Contains(meta, "image/jpeg") && !strings.Contains(meta, "image/png") {
		return "", errors.New("forbidden data format")
	}

	base64Data := splittedURL[1]
	decoded, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", errors.New("unable to base64 decode")
	}

	if len(decoded) > maxSizeKB*1024 {
		return "", errors.New("image is too big")
	}

	img, format, err := image.Decode(bytes.NewReader(decoded))
	if err != nil {
		return "", err
	}

	var ext string
	if (format == "jpeg") || (format == "png") {
		ext = "jpg" // because we convert both jpg and png to jpg
	} else {
		return "", errors.New("forbidden data format")
	}

	path += "." + ext
	f, err := os.Create(path)
	if err != nil {
		return "", errors.New("unable to save file to " + path)
	}
	defer f.Close()
	err = jpeg.Encode(f, img, nil)
	if err != nil {
		return "", errors.New("unable to encode jpeg")
	}

	return path, nil
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n uint8) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func HashPassword(oldPassword string) string {
	hash := sha256.New()
	salt := RandStringRunes(constants.SaltLength)
	_, _ = hash.Write([]byte(salt + oldPassword))
	return salt + base64.URLEncoding.EncodeToString(hash.Sum(nil))
}

func CheckHashedPassword(databasePassword string, gotPassword string) bool {
	salt := databasePassword[:8]
	hash := sha256.New()
	_, _ = hash.Write([]byte(salt + gotPassword))
	gotPassword = base64.URLEncoding.EncodeToString(hash.Sum(nil))

	return gotPassword == databasePassword[8:]
}

func CreateCookieValue(n uint8) string {
	key := RandStringRunes(n)
	return key
}

func CreateCookieWithValue(value string) *http.Cookie {
	newCookie := &http.Cookie{
		Name:     constants.SessionCookieName,
		Value:    value,
		Expires:  time.Now().Add(24 * 7 * time.Hour),
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}

	return newCookie
}
func CreateCookie(n uint8) *http.Cookie {
	key := RandStringRunes(n)

	newCookie := &http.Cookie{
		Name:     constants.SessionCookieName,
		Value:    key,
		Expires:  time.Now().Add(24 * 7 * time.Hour),
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}

	return newCookie
}
