package main

import (
	"fmt"
	"github.com/gmyuval/libracrypt/pkg/docmgmt"
	"github.com/gmyuval/libracrypt/pkg/libracrypt"
)

func main() {
	fmt.Printf("============= Docload example =============\n")
	ld, err := docmgmt.NewLibraDoc("D:\\test_doc.docx")
	if err != nil {
		fmt.Printf("Failed docload with err:\n%v\n", err)
	}
	fmt.Printf("doc struct:\n%v\n", ld)
	fmt.Printf("============= Finished Docload example =============\n")
	fmt.Printf("============= Scramble example =============\n")
	cypher, err := libracrypt.CreateCypher("D:\\test.json")
	if err != nil {
		fmt.Printf("Failed loading cypher with err:\n%v\n", err)
	}

	fmt.Printf("loaded cypher:\n%v", cypher)
}
