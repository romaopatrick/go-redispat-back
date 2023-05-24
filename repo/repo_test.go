package repo

import (
	"fmt"
	"log"
	"os"
	"testing"
)

const fname = "default"
const testdir = "../%s.json"

func TestStoreConnection(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			log.Print(r)
			t.FailNow()
		}
	}()

	path := fmt.Sprintf(testdir, fname)

	t.Cleanup(func() {
		if err := os.Remove(path); err != nil {
			log.Print(err)
		}
	})

	c := NewConnection(fname, "",
		"aws-memorydb-prd-0001-001.aws-memorydb-prd.cv3tti.memorydb.sa-east-1.amazonaws.com:6379",
		"aws-memorydb-prd-0001-002.aws-memorydb-prd.cv3tti.memorydb.sa-east-1.amazonaws.com:6379")

	StoreConnection(c)

	if _, err := os.Open(path); err != nil {
		log.Panicf("err: %s", err.Error())
	}
}
