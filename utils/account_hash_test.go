package utils

import "testing"

func TestAccountHash(t *testing.T) {
	firstKey := AccountHash("01e23d200eb0f3c8a3dacc8453644e6fcf4462585a68234ebb1c3d6cc8971148c2")
	if firstKey != "14b94d33a1be1a2741ddefa7ae68a28cd1956e3801730bea617bf529d50f8aea" {
		t.Errorf("Wrong account hash for 01e23d200eb0f3c8a3dacc8453644e6fcf4462585a68234ebb1c3d6cc8971148c2, should have been : 14b94d33a1be1a2741ddefa7ae68a28cd1956e3801730bea617bf529d50f8aea instead received %s", firstKey)
	}
	secondKey := AccountHash("02035724e2530c5c8f298ba41fe1cafa28294ab7b04d4f1ade025a4a268138570b3a")
	if secondKey != "db5fd4ce31e448eafa70fa2d46d254b3bf5107322da49f2c4e457b5fbfae4f8e" {
		t.Errorf("Wrong account hash for 02035724e2530c5c8f298ba41fe1cafa28294ab7b04d4f1ade025a4a268138570b3a, should have been : db5fd4ce31e448eafa70fa2d46d254b3bf5107322da49f2c4e457b5fbfae4f8e instead received %s", secondKey)
	}
}
