package armor_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"testing"

	armor2 "github.com/cosmos/cosmos-sdk/crypto/armor"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/crypto/bcrypt"
	tmcrypto "github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/armor"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/crypto/xsalsa20symmetric"

	"github.com/cosmos/cosmos-sdk/codec/legacy"
	cryptoAmino "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/types"
)

func TestArmorUnarmorPrivKey(t *testing.T) {
	priv := secp256k1.GenPrivKey()
	armored := armor2.EncryptArmorPrivKey(priv, "passphrase", "")
	_, _, err := armor2.UnarmorDecryptPrivKey(armored, "wrongpassphrase")
	require.Error(t, err)
	decrypted, algo, err := armor2.UnarmorDecryptPrivKey(armored, "passphrase")
	require.NoError(t, err)
	require.Equal(t, string(hd.Secp256k1Type), algo)
	require.True(t, priv.Equals(decrypted))

	// empty string
	decrypted, algo, err = armor2.UnarmorDecryptPrivKey("", "passphrase")
	require.Error(t, err)
	require.True(t, errors.Is(io.EOF, err))
	require.Nil(t, decrypted)
	require.Empty(t, algo)

	// wrong key type
	armored = armor2.ArmorPubKeyBytes(priv.PubKey().Bytes(), "")
	_, _, err = armor2.UnarmorDecryptPrivKey(armored, "passphrase")
	require.Error(t, err)
	require.Contains(t, err.Error(), "unrecognized armor type")

	// armor key manually
	encryptPrivKeyFn := func(privKey tmcrypto.PrivKey, passphrase string) (saltBytes []byte, encBytes []byte) {
		saltBytes = tmcrypto.CRandBytes(16)
		key, err := bcrypt.GenerateFromPassword(saltBytes, []byte(passphrase), armor2.BcryptSecurityParameter)
		require.NoError(t, err)
		key = tmcrypto.Sha256(key) // get 32 bytes
		privKeyBytes := legacy.Cdc.Amino.MustMarshalBinaryBare(privKey)
		return saltBytes, xsalsa20symmetric.EncryptSymmetric(privKeyBytes, key)
	}
	saltBytes, encBytes := encryptPrivKeyFn(priv, "passphrase")

	// wrong kdf header
	headerWrongKdf := map[string]string{
		"kdf":  "wrong",
		"salt": fmt.Sprintf("%X", saltBytes),
		"type": "secp256k",
	}
	armored = armor.EncodeArmor("TENDERMINT PRIVATE KEY", headerWrongKdf, encBytes)
	_, _, err = armor2.UnarmorDecryptPrivKey(armored, "passphrase")
	require.Error(t, err)
	require.Equal(t, "unrecognized KDF type: wrong", err.Error())
}

func TestArmorUnarmorPubKey(t *testing.T) {
	// Select the encryption and storage for your cryptostore
	cstore := keyring.NewInMemory()

	// Add keys and see they return in alphabetical order
	info, _, err := cstore.NewMnemonic("Bob", keyring.English, types.FullFundraiserPath, hd.Secp256k1)
	require.NoError(t, err)
	armored := armor2.ArmorPubKeyBytes(legacy.Cdc.Amino.MustMarshalBinaryBare(info.GetPubKey()), "")
	pubBytes, algo, err := armor2.UnarmorPubKeyBytes(armored)
	require.NoError(t, err)
	pub, err := cryptoAmino.PubKeyFromBytes(pubBytes)
	require.NoError(t, err)
	require.Equal(t, string(hd.Secp256k1Type), algo)
	require.True(t, pub.Equals(info.GetPubKey()))

	armored = armor2.ArmorPubKeyBytes(legacy.Cdc.Amino.MustMarshalBinaryBare(info.GetPubKey()), "unknown")
	pubBytes, algo, err = armor2.UnarmorPubKeyBytes(armored)
	require.NoError(t, err)
	pub, err = cryptoAmino.PubKeyFromBytes(pubBytes)
	require.NoError(t, err)
	require.Equal(t, "unknown", algo)
	require.True(t, pub.Equals(info.GetPubKey()))

	armored, err = cstore.ExportPrivKeyArmor("Bob", "passphrase")
	require.NoError(t, err)
	_, _, err = armor2.UnarmorPubKeyBytes(armored)
	require.Error(t, err)
	require.Equal(t, `couldn't unarmor bytes: unrecognized armor type "TENDERMINT PRIVATE KEY", expected: "TENDERMINT PUBLIC KEY"`, err.Error())

	// armor pubkey manually
	header := map[string]string{
		"version": "0.0.0",
		"type":    "unknown",
	}
	armored = armor.EncodeArmor("TENDERMINT PUBLIC KEY", header, pubBytes)
	_, algo, err = armor2.UnarmorPubKeyBytes(armored)
	require.NoError(t, err)
	// return secp256k1 if version is 0.0.0
	require.Equal(t, "secp256k1", algo)

	// missing version header
	header = map[string]string{
		"type": "unknown",
	}
	armored = armor.EncodeArmor("TENDERMINT PUBLIC KEY", header, pubBytes)
	bz, algo, err := armor2.UnarmorPubKeyBytes(armored)
	require.Nil(t, bz)
	require.Empty(t, algo)
	require.Error(t, err)
	require.Equal(t, "header's version field is empty", err.Error())

	// unknown version header
	header = map[string]string{
		"type":    "unknown",
		"version": "unknown",
	}
	armored = armor.EncodeArmor("TENDERMINT PUBLIC KEY", header, pubBytes)
	bz, algo, err = armor2.UnarmorPubKeyBytes(armored)
	require.Nil(t, bz)
	require.Empty(t, algo)
	require.Error(t, err)
	require.Equal(t, "unrecognized version: unknown", err.Error())
}

func TestArmorInfoBytes(t *testing.T) {
	bs := []byte("test")
	armoredString := armor2.ArmorInfoBytes(bs)
	unarmoredBytes, err := armor2.UnarmorInfoBytes(armoredString)
	require.NoError(t, err)
	require.True(t, bytes.Equal(bs, unarmoredBytes))
}

func TestUnarmorInfoBytesErrors(t *testing.T) {
	unarmoredBytes, err := armor2.UnarmorInfoBytes("")
	require.Error(t, err)
	require.True(t, errors.Is(io.EOF, err))
	require.Nil(t, unarmoredBytes)

	header := map[string]string{
		"type":    "Info",
		"version": "0.0.1",
	}
	unarmoredBytes, err = armor2.UnarmorInfoBytes(armor.EncodeArmor(
		"TENDERMINT KEY INFO", header, []byte("plain-text")))
	require.Error(t, err)
	require.Equal(t, "unrecognized version: 0.0.1", err.Error())
	require.Nil(t, unarmoredBytes)
}

func BenchmarkBcryptGenerateFromPassword(b *testing.B) {
	passphrase := []byte("passphrase")
	for securityParam := 9; securityParam < 16; securityParam++ {
		param := securityParam
		b.Run(fmt.Sprintf("benchmark-security-param-%d", param), func(b *testing.B) {
			saltBytes := tmcrypto.CRandBytes(16)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := bcrypt.GenerateFromPassword(saltBytes, passphrase, param)
				require.Nil(b, err)
			}
		})
	}
}