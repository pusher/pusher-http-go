package pusher

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	// "fmt"
	"strings"

	"golang.org/x/crypto/nacl/secretbox"
)

func hmacSignature(toSign, secret string) string {
	return hex.EncodeToString(hmacBytes([]byte(toSign), []byte(secret)))
}

func hmacBytes(toSign, secret []byte) []byte {
	_authSignature := hmac.New(sha256.New, secret)
	_authSignature.Write(toSign)
	return _authSignature.Sum(nil)
}

func checkSignature(result, secret string, body []byte) bool {
	expected := hmacBytes(body, []byte(secret))
	resultBytes, err := hex.DecodeString(result)

	if err != nil {
		return false
	}
	return hmac.Equal(expected, resultBytes)
}

func createAuthMap(key, secret, stringToSign string, sharedSecret string) map[string]string {
	authSignature := hmacSignature(stringToSign, secret)
	authString := strings.Join([]string{key, authSignature}, ":")
	if sharedSecret != "" {
		return map[string]string{"auth": authString, "shared_secret": sharedSecret}
	}
	return map[string]string{"auth": authString}
}

func md5Signature(body []byte) string {
	_bodyMD5 := md5.New()
	_bodyMD5.Write([]byte(body))
	return hex.EncodeToString(_bodyMD5.Sum(nil))
}

func encrypt(channel string, databytes []byte, encryptionKey string) string {
	sharedSecret := generateSharedSecret(channel, encryptionKey)
	nonce := generateNonce()
	nonceB64 := base64.StdEncoding.EncodeToString(nonce[:])
	cipherText := secretbox.Seal([]byte{}, databytes, &nonce, &sharedSecret)
	cipherTextB64 := base64.StdEncoding.EncodeToString(cipherText)
	return fmt.Sprintf("encrypted_data:%s:%s", nonceB64, cipherTextB64)
}

func generateNonce() [24]byte {
	var nonce [24]byte
	//Trick ReadFull into thinking nonce is a slice
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		panic(err)
	}
	return nonce
}

func generateSharedSecret(channel string, encryptionKey string) [32]byte {
	return sha256.Sum256([]byte(channel + encryptionKey))
}

func decryptEvents(webhookData Webhook, encryptionKey string) (*Webhook, error) {
	// TODO: This is nasty as heck, could probably do with some tidying
	decryptedWebhooks := &Webhook{}
	decryptedWebhooks.TimeMs = webhookData.TimeMs
	for _, event := range webhookData.Events {

		if isEncryptedPayload(event.Data) {
			payload := event.Data
			splitPayload := strings.Split(payload, ":")

			payloadbytes, err := base64.StdEncoding.DecodeString(splitPayload[2])
			if err != nil {
				return decryptedWebhooks, err
			}
			box := []byte(payloadbytes)

			var nonce [24]byte
			nonceBytes, err2 := base64.StdEncoding.DecodeString(splitPayload[1])
			if err2 != nil {
				return decryptedWebhooks, err
			}
			copy(nonce[:], []byte(nonceBytes[:]))

			sharedSecret := generateSharedSecret(event.Channel, encryptionKey)
			decryptedBox, ok := secretbox.Open([]byte{}, box, &nonce, &sharedSecret)
			if !ok {
				return decryptedWebhooks, errors.New("Failed to decrypt event, possibly wrong key")
			}
			event.Data = string(decryptedBox)
		}
		decryptedWebhooks.Events = append(decryptedWebhooks.Events, event)
	}
	return decryptedWebhooks, nil
}
