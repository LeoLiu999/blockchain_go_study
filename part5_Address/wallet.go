package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"

	"golang.org/x/crypto/ripemd160"
)

const version = byte(0x00)   //版本
const addressChecksumLen = 4 //校验和是结果哈希的前四个字节。

//Wallet 钱包
type Wallet struct {
	//一个钱包只存储一个密钥对
	//ecdsa 比特币使用椭圆曲线来产生私钥 确保生成真正随机的字节
	//比特币使用的是 ECDSA（Elliptic Curve Digital Signature Algorithm）算法来对交易进行签名
	//在比特币中使用的曲线可以随机选取在 0 与 2 ^ 2 ^ 56（大概是 10^77, 而整个可见的宇宙中，原子数在 10^78 到 10^82 之间） 的一个数。
	//有如此高的一个上限，意味着几乎不可能发生有两次生成同一个私钥的事情。
	//这里也使用该算法
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

//GetAddress get address 生成地址
func (w Wallet) GetAddress() []byte {
	//1:先将publicKey用sha256 hash 再用ripemd160 hash ripemd160(sha256(pubkey))
	pubKeyHash := HashPubKey(w.PublicKey)
	//2:给哈希加上地址生成算法版本的前缀
	versionPayload := append([]byte{version}, pubKeyHash...)
	//3:使用 SHA256(SHA256(payload)) 再哈希，计算校验和。校验和是结果哈希的前四个字节。
	checksum := checksum(versionPayload)
	//4:将校验和加入{version, pubkeyhash}中 最终的地址为{version, pubkeyhash, checksum}
	fullPayload := append(versionPayload, checksum...)
	//5:将地址base58 为人类可理解的地址明文
	address := Base58Encode(fullPayload)

	return address
}

//NewWallet create Wallet
func NewWallet() *Wallet {
	private, public := newKeyPair()

	wallet := Wallet{private, public}

	return &wallet
}

//newKeyPair 生成一对key
func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256() //elliptic包实现了几条覆盖素数有限域的标准椭圆曲线。
	//ECDSA基于椭圆曲线，所以我们需要一个椭圆曲线。接下来，使用椭圆生成一个私钥，然后再从私钥生成一个公钥。
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader) //rand.Reader强伪随机生成器
	if err != nil {
		log.Panic(err)
	}
	//有一点需要注意：在基于椭圆曲线的算法中，公钥是曲线上的点。因此，公钥是 X，Y 坐标的组合。在比特币中，这些坐标会被连接起来，然后形成一个公钥。
	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)

	return *privateKey, publicKey
}

//HashPubKey HashPubKey
//先用sha256 hash 再用ripemd160 hash ripemd160(sha256(pubkey))
func HashPubKey(pubKey []byte) []byte {
	pubSha256 := sha256.Sum256(pubKey)

	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(pubSha256[:])
	if err != nil {
		log.Panic(err)
	}

	publicRipemd160 := RIPEMD160Hasher.Sum(nil)
	return publicRipemd160
}

//checksum 两次sha256计算校验和。校验和是结果哈希的前四个字节。
func checksum(versionPayload []byte) []byte {
	sha1 := sha256.Sum256(versionPayload)
	sha2 := sha256.Sum256(sha1[:])

	return sha2[:addressChecksumLen]
}
