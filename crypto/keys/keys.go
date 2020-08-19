package keys

import (
	"github.com/cosmos/cosmos-sdk/crypto"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
)

var (
	_, _, _ crypto.PubKey  = &Secp256K1PubKey{}, &Ed25519PubKey{}, &Sr25519PubKey{}
	_, _, _ crypto.PrivKey = &Secp256K1PrivKey{}, &Ed25519PrivKey{}, &Sr25519PrivKey{}
)

func (m *Secp256K1PubKey) Address() crypto.Address {
	return m.Key.Address()
}

func (m *Secp256K1PubKey) Bytes() []byte {
	return m.Key.Bytes()
}

func (m *Secp256K1PubKey) VerifySignature(msg []byte, sig []byte) bool {
	return m.Key.VerifySignature(msg, sig)
}

func (m *Secp256K1PubKey) Equals(key crypto.PubKey) bool {
	return m.Key.Equals(key)
}

func (m *Secp256K1PubKey) Type() string {
	return m.Key.Type()
}

func (m *Ed25519PubKey) Address() crypto.Address {
	return m.Key.Address()
}

func (m *Ed25519PubKey) Bytes() []byte {
	return m.Key.Bytes()
}

func (m *Ed25519PubKey) VerifySignature(msg []byte, sig []byte) bool {
	return m.Key.VerifySignature(msg, sig)
}

func (m *Ed25519PubKey) Equals(key crypto.PubKey) bool {
	return m.Key.Equals(key)
}

func (m *Ed25519PubKey) Type() string {
	return m.Key.Type()
}

func (m *Sr25519PubKey) Address() crypto.Address {
	return m.Key.Address()
}

func (m *Sr25519PubKey) Bytes() []byte {
	return m.Key.Bytes()
}

func (m *Sr25519PubKey) VerifySignature(msg []byte, sig []byte) bool {
	return m.Key.VerifySignature(msg, sig)
}

func (m *Sr25519PubKey) Equals(key crypto.PubKey) bool {
	return m.Key.Equals(key)
}

func (m *Sr25519PubKey) Type() string {
	return m.Key.Type()
}

func (m *Secp256K1PrivKey) Bytes() []byte {
	return m.Bytes()
}

func (m *Secp256K1PrivKey) Sign(msg []byte) ([]byte, error) {
	return m.Sign(msg)
}

func (m *Secp256K1PrivKey) PubKey() crypto.PubKey {
	return &Secp256K1PubKey{Key: m.PubKey().(secp256k1.PubKey)}
}

func (m *Secp256K1PrivKey) Equals(key crypto.PrivKey) bool {
	return m.Key.Equals(key)
}

func (m *Secp256K1PrivKey) Type() string {
	return m.Key.Type()
}
