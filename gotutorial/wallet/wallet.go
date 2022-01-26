package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"gotutorial/utils"
	"math/big"
	"os"
)

const (
	privateKEy    string = "307702010104206f42f81044928170d0a1aa716cda556134d051ad3f1fb4012cd62fb1a7e89547a00a06082a8648ce3d030107a144034200043558f6ae7ae9fc2af2e8230d6aad50a07acaa98b04d909ae8330ac169b83144769009df09c33e7f859d9487d0bdb1bfe87d8364b13d32de2d8e67a3ffcb0cb81"
	signature     string = "be12ea139814e08b14eca72689a23834c3c9951ec0c93a20180df4f711b6ed201eed6e980d4ceb7b75adc0620dd8102f8cc3acbfb9aed3ec37323ca6c0646b4c"
	hashedMessage        = "da504da2fceb673eacd7e1cb1e5b6f8d5e1ea2720d8ab8e664742a56014340b2"

	walletName string = "sycoin.wallet"
)

type wallet struct {
	privateKey *ecdsa.PrivateKey
	Address    string
}

var w *wallet

func hasWalletFile() bool {
	_, err := os.Stat(walletName)
	return !os.IsNotExist(err)
}

func createPrivateKey() *ecdsa.PrivateKey {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleErr(err)
	return privateKey
}

// func restoreKey() *ecdsa.PrivateKey {
// 	bytes, err := os.ReadFile(walletName)
// 	utils.HandleErr(err)
// 	privateKey, err := x509.ParseECPrivateKey(bytes)
// 	utils.HandleErr(err)
// 	return privateKey
// }
func restoreKey() (key *ecdsa.PrivateKey) {
	bytes, err := os.ReadFile(walletName)
	utils.HandleErr(err)
	key, err = x509.ParseECPrivateKey(bytes)
	utils.HandleErr(err)
	return
}

func persistKey(privateKey *ecdsa.PrivateKey) {
	bytes, err := x509.MarshalECPrivateKey(privateKey)
	utils.HandleErr(err)
	os.WriteFile(walletName, bytes, 0644)
}

func aFromK(p *ecdsa.PrivateKey) string {
	return encodeBigInt(p.X.Bytes(), p.Y.Bytes())
}

func Sign(wallet *wallet, payload string) string {
	bytes, err := hex.DecodeString(payload)
	utils.HandleErr(err)
	r, s, err := ecdsa.Sign(rand.Reader, wallet.privateKey, bytes)
	return encodeBigInt(r.Bytes(), s.Bytes())
}

func extractBigInt(hash string) (big.Int, big.Int, error) {
	bytes, err := hex.DecodeString(hash)
	utils.HandleErr(err)

	a, b := big.Int{}, big.Int{}

	a.SetBytes(bytes[:len(bytes)/2])
	b.SetBytes(bytes[len(bytes)/2:])
	return a, b, nil
}

func encodeBigInt(a, b []byte) string {
	bytes := append(a, b...)
	return fmt.Sprintf("%x", bytes)
}

func Verify(signature string, payload string, publicKey string) bool {
	r, s, err := extractBigInt(signature)
	utils.HandleErr(err)

	payloadBytes, err := hex.DecodeString(payload)
	utils.HandleErr(err)

	x, y, err := extractBigInt(publicKey)
	utils.HandleErr(err)

	pk := ecdsa.PublicKey{Curve: elliptic.P256(), X: &x, Y: &y}

	return ecdsa.Verify(&pk, payloadBytes, &r, &s)
}

func Wallet() *wallet {
	if w == nil {
		w = &wallet{}
		var privateKey *ecdsa.PrivateKey
		if hasWalletFile() {
			privateKey = restoreKey()
			fmt.Println(privateKey)
		} else {
			privateKey = createPrivateKey()
			persistKey(privateKey)
		}
		w.privateKey = privateKey
		w.Address = aFromK(privateKey)
	}
	return w
}

func Start() {
	// privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	// pkBytes, err := x509.MarshalECPrivateKey(privateKey)
	// fmt.Printf("%x\n\n\n", pkBytes)

	// utils.HandleErr(err)
	// message := "i love sy"
	// hashedMessage := utils.Hash(message)
	// fmt.Printf("%s\n\n", hashedMessage)
	// bytes, _ := hex.DecodeString(hashedMessage)
	// r, s, err := ecdsa.Sign(rand.Reader, privateKey, bytes)
	// signature := append(r.Bytes(), s.Bytes()...)
	// fmt.Printf("%x\n\n", signature)
	pkBytes, err := hex.DecodeString(privateKEy)
	utils.HandleErr(err)
	private, err := x509.ParseECPrivateKey(pkBytes)
	utils.HandleErr(err)
	fmt.Println(private)
	signatureBytes, err := hex.DecodeString(signature)
	utils.HandleErr(err)
	mid := len(signatureBytes) / 2
	rBytes := signatureBytes[:mid]
	sBytes := signatureBytes[mid:]

	var bigR, bigS = big.Int{}, big.Int{}
	bigR.SetBytes(rBytes)
	bigS.SetBytes(sBytes)
	fmt.Println(bigR, bigS)

	dataBytes, err := hex.DecodeString(hashedMessage)
	utils.HandleErr(err)

	ok := ecdsa.Verify(&private.PublicKey, dataBytes, &bigR, &bigS)

	fmt.Println(ok)
}

// signature verification publickey, privatekey
// -> backend for wallet
// -> signature and verification with transaction

/*

1. "my data" -> Hash -> "hashed message"

2. Generate key pair(Public, Private)

3. "hashed message" + Private Key -> Signature

4. "hashed message" + Signature + Public Key -> Verification true or false

*/
