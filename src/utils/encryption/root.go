package encryption

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"github.com/cossacklabs/themis/gothemis/keys"
	"github.com/cossacklabs/themis/gothemis/message"
)

type Config struct{}

type Self_Config struct {
	Name     string `json:"name"`
	Priv_key string `json:"priv_key"`
	Publ_key string `json:"publ_key"`
}
type Peer_Config struct {
	Name     string `json:"name"`
	Publ_key string `json:"publ_key"`
}

type peer struct {
	Name     string
	Publ_key *keys.PublicKey
}

type H8go struct {
	Name     string
	Priv_key *keys.PrivateKey
	Publ_key *keys.PublicKey
	Peers    peer
}

type Config_Pack struct {
	Name        string
	Priv_key    string `json:"priv_key"`
	Publ_key    string `json:"publ_key"`
	Peer_config Peer_Config
}

type CONFIG_PACK interface {
	debug_print()
}

func (self H8go) debug_print() {
	fmt.Println("SELF")
	fmt.Println(self.Name)
	fmt.Println(self.Priv_key)
	fmt.Println(self.Publ_key)
	fmt.Println("")
	fmt.Println("PEER")
	fmt.Println(self.Peers.Name)
	fmt.Println(self.Peers.Publ_key)
}

func load_config_file(filename string) H8go {
	var tmp Config_Pack
	var result H8go

	b, err := os.ReadFile(filename) // just pass the file name
	if err != nil {
		fmt.Print(err)
	}

	json.Unmarshal(b, &tmp) // unmarshal means convert to struct

	// Cursed af :puke:
	result.Name = tmp.Name

	meow, _ := base64.StdEncoding.DecodeString(tmp.Publ_key)
	result.Publ_key = &keys.PublicKey{
		Value: meow,
	}

	meow2, _ := base64.StdEncoding.DecodeString(tmp.Priv_key)
	result.Priv_key = &keys.PrivateKey{
		Value: meow2,
	}

	result.Peers.Name = tmp.Peer_config.Name

	meow3, _ := base64.StdEncoding.DecodeString(tmp.Peer_config.Publ_key)
	result.Peers.Publ_key = &keys.PublicKey{
		Value: meow3,
	}

	return result
}

// This will eventually read the whole json and return an object matching the json
// representation, allowing the client to then use clientKeys.PubKey, clientKeys.PeerPubKey
// etc.
// Limitation of thing outlined above is that it only supports a single peer, but #TODO :P
func Load_Keys(src string) H8go {
	var conf = load_config_file(src + ".json")
	conf.debug_print()

	return conf
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
