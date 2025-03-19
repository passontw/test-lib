package base

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"testing"
	"time"
)

func ExampleNewCBCEncrypter() {
	// Load your secret key from a safe place and reuse it across multiple
	// NewCipher calls. (Obviously don't use this example key for anything
	// real.) If you want to convert a passphrase to a key, use a suitable
	// package like bcrypt or scrypt.
	key, _ := hex.DecodeString("6368616e676520746869732070617373")
	plaintext := []byte("exampleplaintext")

	// CBC mode works on blocks so plaintexts may need to be padded to the
	// next whole block. For an example of such padding, see
	// https://tools.ietf.org/html/rfc5246#section-6.2.3.2. Here we'll
	// assume that the plaintext is already of the correct length.
	if len(plaintext)%aes.BlockSize != 0 {
		panic("plaintext is not a multiple of the block size")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	// It's important to remember that ciphertexts must be authenticated
	// (i.e. by using crypto/hmac) as well as being encrypted in order to
	// be secure.

	fmt.Printf("%x\n", ciphertext)
}

func ExampleNewCFBDecrypter() {
	// Load your secret key from a safe place and reuse it across multiple
	// NewCipher calls. (Obviously don't use this example key for anything
	// real.) If you want to convert a passphrase to a key, use a suitable
	// package like bcrypt or scrypt.

	//text := "123456"
	//strkey:="e10adc3949ba59abbe56e057f20f883e"
	//text:c5fdd1b4b5925b3815179123b3e71e66

	key, _ := hex.DecodeString("6368616e676520746869732070617373")
	ciphertext, _ := hex.DecodeString("7dd015f06bec7f1b8f6559dad89f4131da62261786845100056b353194ad")

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)
	fmt.Printf("%s", ciphertext)
	// Output: some plaintext
}

func TestHttpEncTest(t *testing.T) {
	authPwd := "xdr56yhn"
	keyMd5 := GetMd5(authPwd)
	//paramm := `{"timestamp":"1111","userName":"mina1234561"}`
	paramm := fmt.Sprintf(`{"serverId":"A01","serverPwd":"iv66","timestamp":"%d"}`, time.Now().UnixNano())
	keyByte, _ := hex.DecodeString(keyMd5)
	encrypt, error := AesEcbEncrypt([]byte(paramm), keyByte, PKCS5Padding)
	fmt.Println(hex.EncodeToString(encrypt), error)

	//decrypt
	//str := hex.EncodeToString(encrypt)
	//str := "7bcdd32297b8d2e196b2b379f1cd83fd2985e2c088c6b497ee0ca1b42cb8536b318cfbb3ef649d7bbbd177cc24acf783ed65f000fcf8cd048d9ca299bd922c21acee24a6184ae4bcf2148ec884ba44fb"
	str := "27186c4bfae27ba4b207279757828ea4f11cd8169236e47524b943eb5a652136ee0123b28dc00545377a976bd3e478c5189c2d066f248e0f999c5deb9b8b8bc89146fd2ba051829131e33f6241afaad70d21588f3512986065aec30aedcc9d5ff8f23b955d2f31769d41dc808f825fde271ee2231606cf9a9fbc1f1b6e50ece12034a52c590489af34110d4130736fff33a48641859da378dd29d978686a9dd377ea1ba9fadfd02e480f4cb251f3c48402b8a57b181b26b6aad5e999b86cd4274b9c5f3d816d788b8b687249565ba833233f07cc63948abee6a380dac13da7871c5cf18f4f7b4c2f7463e1d8b2b2268fdda3b883c82732942334bfef96aff21fc161ff4e4c1321a2b600c6b35311c016a7df5b63998db5323f364a12e568c0f0811eb60813f1b5635e452680af7da2db"
	enByte, _ := hex.DecodeString(str)
	decrypt, err := AesEcbDecrypt(enByte, keyByte, PKCS5UnPadding)
	fmt.Println(string(decrypt), err)
}

func TestHtonl(t *testing.T) {
	//fmt.Println(unsafe.Sizeof(PacketHeader{}))
	fmt.Println("des")
	//   Des des("a0fa6db8", 8, Cipher::ENCODING::BASE64, Cipher::MODE::ECB, Cipher::PADDING::PKCS7);
	// syQNvjGylUm/nz3e7lfi7Q==;
	// len(key) % 8 == 0, len(iv) == 8

	//dbconn := "Provider=OraOLEDB.Oracle;Data Source=MYTEST150;User id=newAg;Password=syQNvjGylUm/nz3e7lfi7Q==;"

	//fmt.Println(GetDSN(dbconn))

	//str := ""

	//src, _ := DesDecrypt(str, "a0fa6db8", PKCS7UnPadding)
	//fmt.Println(src)

	/*
		var text = "ag123qwe"
		str, _ := DesEncrypt(text, "a0fa6db8", PKCS7Padding)
		fmt.Println(str)

		src, _ := DesDecrypt(str, "a0fa6db8", PKCS7UnPadding)
		fmt.Println(src)
	*/

	//text := "123456"
	//strkey:="e10adc3949ba59abbe56e057f20f883e"
	//text:c5fdd1b4b5925b3815179123b3e71e66
	//AesEncrypt([]byte(text), []byte(strkey), "", PKCS7Padding)

	//text := "123456"
	//strkey:="e10adc3949ba59abbe56e057f20f883e"
	//text:c5fdd1b4b5925b3815179123b3e71e66

	// Load your secret key from a safe place and reuse it across multiple
	// NewCipher calls. (Obviously don't use this example key for anything
	// real.) If you want to convert a passphrase to a key, use a suitable
	// package like bcrypt or scrypt.
	/*
		key, _ := hex.DecodeString("e10adc3949ba59abbe56e057f20f883e")
		//key := "e10adc3949ba59abbe56e057f20f883e"
		plaintext := []byte("123456")

		// CBC mode works on blocks so plaintexts may need to be padded to the
		// next whole block. For an example of such padding, see
		// https://tools.ietf.org/html/rfc5246#section-6.2.3.2. Here we'll
		// assume that the plaintext is already of the correct length.
		if len(plaintext)%aes.BlockSize != 0 {
			panic("plaintext is not a multiple of the block size")
		}

		block, err := aes.NewCipher(key)
		if err != nil {
			panic(err)
		}

		// The IV needs to be unique, but not secure. Therefore it's common to
		// include it at the beginning of the ciphertext.
		ciphertext := make([]byte, aes.BlockSize+len(plaintext))
		iv := ciphertext[:aes.BlockSize]
		if _, err := io.ReadFull(rand.Reader, iv); err != nil {
			panic(err)
		}

		mode := cipher.NewCBCEncrypter(block, iv)
		mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

		// It's important to remember that ciphertexts must be authenticated
		// (i.e. by using crypto/hmac) as well as being encrypted in order to
		// be secure.

		fmt.Printf("%x\n", ciphertext)
	*/
	now := time.Now().Unix()
	now += 100
	fmt.Println(ConvertSecsToStr(now))

	a := IpToInt("127.0.0.1")
	fmt.Println(a)
	v := IpToStr(a + 1)
	fmt.Println(v)

	fmt.Println(ToHex([]byte{1, 1, 1}))

	//var mapdata map[string]string
	//fmt.Println(mapdata)
	//mapdata["aa"] = "hello"
	//fmt.Println(mapdata)

}
