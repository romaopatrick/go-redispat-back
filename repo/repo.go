package repo

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

const dbtemplate = "./db/%s.json"

func StoreConnection(c *Connection) {
	f, err := os.Create(fmt.Sprintf(dbtemplate, c.Name))
	if err != nil {
		log.Printf("error: %s", err.Error())
		panic(err)
	}

	defer f.Close()

	b, err := json.Marshal(c)
	if err != nil {
		log.Printf("error: %s", err.Error())
		panic(err)
	}

	f.Write(b)
}

func GetConnection(cname string) *Connection {
	f, err := os.Open(fmt.Sprintf(dbtemplate, cname))
	if err != nil {
		log.Printf("error: %s", err.Error())
		return nil
	}

	c := &Connection{}

	b, err := io.ReadAll(f)
	if err != nil {
		log.Printf("error: %s", err.Error())
		return nil
	}

	err = json.Unmarshal(b, &c)
	if err != nil {
		log.Printf("error: %s", err.Error())
		return nil
	}

	return c
}
