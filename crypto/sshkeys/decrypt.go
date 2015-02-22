package sshkeys

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"strings"

	"github.com/Scalingo/cli/term"
	"gopkg.in/errgo.v1"
)

func (p *PrivateKey) Decrypt() error {
	var decryptedKey []byte
	iv, err := p.iv()
	if err != nil {
		return errgo.Mask(err)
	}
	cipher := p.cipher()

	if !p.IsCipherImplemented(cipher) {
		return fmt.Errorf(
			"Failed to read key '%s': password protected keys with '%s' are\n"+
				"not supported. Please decrypt the key prior to use.", p.Path, cipher)
	}

	if p.PasswordMethod == nil {
		p.PasswordMethod = term.Password
	}

	password, err := p.PasswordMethod("Encrypted SSH Key, password: ")
	if err != nil {
		return errgo.Mask(err)
	}

	switch cipher {
	case "DES-EDE3-CBC":
		key := genDES3Key(password, iv)
		decryptedKey, err = decryptKey(p.Block.Bytes, iv, key, des.NewTripleDESCipher)
		if err != nil {
			return errgo.Newf("Key is tagged DES-ECE3-CBC, but is not: %v", err)
		}
	case "AES-128-CBC":
		key := genAESKey(password, iv)
		decryptedKey, err = decryptKey(p.Block.Bytes, iv, key, aes.NewCipher)
		if err != nil {
			return errgo.Newf("Key is tagged AES-128-CBC, but is not: %v", err)
		}
	}
	decryptedBlock := &pem.Block{}
	decryptedBlock.Type = p.Block.Type
	decryptedBlock.Bytes = decryptedKey
	p.Block = decryptedBlock
	return nil
}

func (p *PrivateKey) cipher() string {
	splitCipher := strings.Split(p.Block.Headers["DEK-Info"], ",")
	return splitCipher[0]
}

func (p *PrivateKey) iv() ([]byte, error) {
	splitCipher := strings.Split(p.Block.Headers["DEK-Info"], ",")
	ivStr := strings.TrimSpace(splitCipher[1])
	iv, err := hex.DecodeString(ivStr)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return iv, nil
}

func genDES3Key(passphrase string, iv []byte) []byte {
	key := append([]byte(passphrase), iv[0:8]...)
	keyHash := md5.New()
	keyHash.Write(key)
	d1 := keyHash.Sum(nil)
	key = append(d1, []byte(passphrase)...)
	key = append(key, iv[0:8]...)
	keyHash = md5.New()
	keyHash.Write(key)
	return append(d1, keyHash.Sum(nil)[0:8]...)
}

func genAESKey(passphrase string, iv []byte) []byte {
	key := append([]byte(passphrase), iv[0:8]...)
	keyHash := md5.New()
	keyHash.Write(key)
	return keyHash.Sum(nil)
}

func decryptKey(payload []byte, iv []byte, key []byte, newCypherFunc func([]byte) (cipher.Block, error)) ([]byte, error) {
	decryptedPayload := make([]byte, len(payload))
	block, err := newCypherFunc(key)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	decrypter := cipher.NewCBCDecrypter(block, iv)
	decrypter.CryptBlocks(decryptedPayload, payload)
	decryptedPayload = bytes.TrimRight(decryptedPayload, "\x02\x08\x09\x0a")
	return decryptedPayload, nil
}
