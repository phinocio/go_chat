package encryption

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"github.com/cossacklabs/themis/gothemis/keys"
	"github.com/cossacklabs/themis/gothemis/message"
)

type Self_Config struct {
	Name     string `json:"name"`
	Priv_key string `json:"priv_key"`
	Publ_key string `json:"publ_key"`
}
type Peer_Config struct {
	Name     string `json:"name"`
	Publ_key string `json:"publ_key"`
}
type Config_Pack struct {
	Self_config Self_Config `json:"self"`
	Peer_config Peer_Config `json:"peer"`
}

// type CONFIG_PACK interface {
// 	debug_print()
// }

// func (self Config_Pack) debug_print() {
// 	fmt.Println("SELF")
// 	fmt.Println(self.Self_config.Name)
// 	fmt.Println(self.Self_config.Priv_key)
// 	fmt.Println(self.Self_config.Publ_key)
// 	fmt.Println("")
// 	fmt.Println("PEER")
// 	fmt.Println(self.Peer_config.Name)
// 	fmt.Println(self.Peer_config.Publ_key)
// }

func load_config_file(filename string) Config_Pack {
	var result Config_Pack

	b, err := os.ReadFile(filename) // just pass the file name
	if err != nil {
		fmt.Print(err)
	}

	json.Unmarshal(b, &result) // unmarshal means convert to struct

	return result
}

// This will eventually read the whole json and return an object matching the json
// representation, allowing the client to then use clientKeys.PubKey, clientKeys.PeerPubKey
// etc.
// Limitation of thing outlined above is that it only supports a single peer, but #TODO :P
func Load_Keys(src string) (*keys.PublicKey, *keys.PrivateKey) {
	var conf = load_config_file(src + ".json")

	b64Pub, _ := base64.StdEncoding.DecodeString(conf.Self_config.Publ_key)
	b64Prv, _ := base64.StdEncoding.DecodeString(conf.Self_config.Priv_key)

	var publicKey = &keys.PublicKey{
		Value: b64Pub,
	}

	var privateKey = &keys.PrivateKey{
		Value: b64Prv,
	}

	return publicKey, privateKey
}

func Gen_Keys() *keys.Keypair {
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

func Encryptor(clear_text []byte, my_priv_key *keys.PrivateKey, peer_publ_key *keys.PublicKey) []byte {
	aliceToBob := message.New(my_priv_key, peer_publ_key)
	cipher_text, err := aliceToBob.Wrap(clear_text)
	if err != nil {
		fmt.Println("message cannot be empty")
	}
	return cipher_text
}

func Decryptor(cipher_text []byte, my_priv_key *keys.PrivateKey, peer_publ_key *keys.PublicKey) []byte {
	// we do not need the priv key to dec, the SecureMessage API is just misleading
	bobToAlice := message.New(my_priv_key, peer_publ_key)
	decrypted, err := bobToAlice.Unwrap(cipher_text)
	if err != nil {
		fmt.Println("decryption failure")
		fmt.Println(err)
	}
	return decrypted
}
