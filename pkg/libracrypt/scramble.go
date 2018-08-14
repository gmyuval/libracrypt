package libracrypt

import (
	"github.com/gmyuval/libracrypt/pkg/docmgmt"
	log "github.com/sirupsen/logrus"

	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"
	"unicode/utf8"
)

type ScrambleCypher struct {
	Cypher        map[rune]rune
	SrcDoc        docmgmt.LibraDoc
	DstDoc        docmgmt.LibraDoc
	EncryptedText []string
}

func NewScrambleCypher(srcFile string, cypherJson string) (*ScrambleCypher, error) {
	doc, err := docmgmt.NewLibraDoc(srcFile)
	if err != nil {
		return nil, errors.New("failed creating libradoc")
	}
	cypher, err := CreateCypher(cypherJson)
	if err != nil {
		return nil, errors.New("failed creating cypher")
	}
	return &ScrambleCypher{
		Cypher: cypher,
		SrcDoc: *doc,
	}, nil
}

func CreateCypher(cypherPath string) (map[rune]rune, error) {
	jsonCypher, err := ioutil.ReadFile(cypherPath)
	if err != nil {
		log.WithFields(log.Fields{
			"file":  cypherPath,
			"error": err,
		}).Fatal("Failed loading json")
		return nil, err
	}
	//defer jsonCypher.Close()
	var cypher map[string]string
	var runeCypher map[rune]rune
	json.Unmarshal([]byte(jsonCypher), &cypher)
	runeCypher = make(map[rune]rune)
	for key, value := range cypher {
		if utf8.RuneCountInString(key) != 1 || utf8.RuneCountInString(value) != 1 {
			log.WithFields(log.Fields{
				"key":          key,
				"length key":   len(key),
				"value":        value,
				"length value": len(value),
			}).Fatal("Key or value have more than single character")
			return nil, errors.New("scramble json contains non single char keys/values")
		}
		runeCypher[[]rune(key)[0]] = []rune(value)[0]
	}
	return runeCypher, nil
}

func Scramble(sCypher *ScrambleCypher, overwrite bool) error {
	if sCypher.EncryptedText != nil {
		if overwrite {
			sCypher.EncryptedText = nil
		} else {
			return errors.New("encrypted text exists")
		}
	}
	for _, line := range sCypher.SrcDoc.SrcText {
		lineArray := strings.Fields(line)
		var encLine []string
		for _, word := range lineArray {
			runeWord := []rune(word)
			var encWord []rune
			for _, r := range runeWord {
				if v, found := sCypher.Cypher[r]; found {
					encWord = append(encWord, v)
				} else {
					encWord = append(encWord, r)
				}
			}
			encLine = append(encLine, string(encWord))
		}
		sCypher.EncryptedText = append(sCypher.EncryptedText, strings.Join(encLine, " "))
	}
	return nil
}
