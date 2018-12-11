package pwdhelper

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)

/**
 * 生成rsa公私钥
 */
func CreateRsaKey() map[string]string {
	var key_map map[string]string
	key_map = make(map[string]string)
	// 私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return nil
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	data := pem.EncodeToMemory(block)
	if data == nil {
		return nil
	}
	key_map["private_key"] = string(data)

	// 公钥
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	data_pu := pem.EncodeToMemory(block)
	if data_pu == nil {
		return nil
	}
	var pu_key string
	pu_key = string(data_pu)
	//pu_key = pu_key[strings.IndexAny(pu_key, "Y")+7:]
	//pu_key = pu_key[:strings.LastIndex(pu_key, "N")-7]
	key_map["public_key"] = string(pu_key)
	return key_map
}

/**
 * rsa加密
 */
func RsaEncrypt(origData []byte, publicKey string) (string, error) {
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return "", errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}
	pub := pubInterface.(*rsa.PublicKey)

	partLen := pub.N.BitLen()/8 - 11
	chunks := rsa_split(origData, partLen)

	buffer := bytes.NewBufferString("")
	for _, chunk := range chunks {
		bytes, err := rsa.EncryptPKCS1v15(rand.Reader, pub, chunk)
		if err != nil {
			return "", err
		}
		buffer.Write(bytes)
	}
	return base64.StdEncoding.EncodeToString(buffer.Bytes()), nil
}

/**
 * rsa解密
 */
func RsaDecrypt(ciphertext string, publicKey string, privateKey string) (string, error) {
	pub_block, _ := pem.Decode([]byte(publicKey))
	if pub_block == nil {
		return "", errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(pub_block.Bytes)
	if err != nil {
		return "", err
	}
	pub := pubInterface.(*rsa.PublicKey)

	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return "", errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	partLen := pub.N.BitLen() / 8
	raw, err := base64.StdEncoding.DecodeString(ciphertext)
	chunks := rsa_split([]byte(raw), partLen)

	buffer := bytes.NewBufferString("")
	for _, chunk := range chunks {
		decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, priv, chunk)
		if err != nil {
			return "", err
		}
		buffer.Write(decrypted)
	}
	return buffer.String(), err
}

/**
 * rsa分段处理
 */
func rsa_split(buf []byte, lim int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:len(buf)])
	}
	return chunks
}
