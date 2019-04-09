package go_schnorr

import (
	"errors"
	"github.com/btcsuite/btcd/btcec"
	"math/big"
)

var (
	secp256k1 = btcec.S256()
)

func TrySign(privateKey []byte, publicKey []byte, message []byte, k []byte) ([]byte, []byte, error) {

	priKey := new(big.Int).SetBytes(privateKey)

	// 1a. check if private key is 0
	if priKey.Cmp(new(big.Int).SetInt64(0)) <= 0 {
		return nil, nil, errors.New("private key must be > 0")
	}

	// 1b. check if private key is less than curve order, i.e., within [1...n-1]
	if priKey.Cmp(secp256k1.N) >= 0 {
		return nil, nil, errors.New("private key cannot be greater than curve order")
	}

	// 2. Compute commitment Q = kG, where G is the base point
	Qx, Qy := secp256k1.ScalarBaseMult(k)

	Q := Compress(secp256k1, Qx, Qy)

	// 3. Compute the challenge r = H(Q || pubKey || msg)
	// mod reduce r by the order of secp256k1, n
	r := new(big.Int).SetBytes(hash(Q, publicKey, message[:]))
	r = r.Mod(r, secp256k1.N)

	if r.Cmp(new(big.Int).SetInt64(0)) == 0 {
		return nil, nil, errors.New("invalid r")
	}

	//4. Compute s = k - r * prv
	// 4a. Compute r * prv
	_r := *r
	s := new(big.Int).Mod(_r.Sub(new(big.Int).SetBytes(k), _r.Mul(&_r, priKey)), secp256k1.N)

	if s.Cmp(big.NewInt(0)) == 0 {
		return nil, nil, errors.New("invalid s")
	}

	return r.Bytes(), s.Bytes(), nil
}
