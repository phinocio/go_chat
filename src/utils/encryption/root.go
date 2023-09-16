package encryption

import (
	"fmt"
	"os"
	"encoding/base64"
	"github.com/cossacklabs/themis/gothemis/keys"
	"github.com/cossacklabs/themis/gothemis/message"
)

func Gen_Keys() ( *keys.Keypair) {
	alice_keyPair, err := keys.New(keys.TypeEC)
	if nil != err {
		fmt.Println("Keypair generating error")
		os.Exit(1)
	}
	return alice_keyPair
}

func Print_Keys(key_pair *keys.Keypair) {
	fmt.Print("Private: ")
	fmt.Println(base64.StdEncoding.EncodeToString(key_pair.Private.Value))
	fmt.Print("Public: ")
	fmt.Println(base64.StdEncoding.EncodeToString(key_pair.Public.Value))
	fmt.Println()
}

func Encryptor(clear_text []byte, my_priv_key *keys.PrivateKey, peer_publ_key *keys.PublicKey) ([]byte) {
	aliceToBob := message.New(my_priv_key, peer_publ_key)
	cipher_text, err := aliceToBob.Wrap(clear_text)
	if err != nil {
		fmt.Println("message cannot be empty")
	}
	return cipher_text
}

func Decryptor(cipher_text []byte, my_priv_key *keys.PrivateKey, peer_publ_key *keys.PublicKey) ([]byte) {
	// we do not need the priv key to dec, the SecureMessage API is just misleading
	bobToAlice := message.New(my_priv_key, peer_publ_key)
	decrypted, err := bobToAlice.Unwrap(cipher_text)
	if err != nil {
		fmt.Println("decryption failure", err)
	}
	return decrypted
}

