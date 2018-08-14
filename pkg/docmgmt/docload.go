package docmgmt

import (
	"baliance.com/gooxml/document"
	log "github.com/sirupsen/logrus"
	"strings"

	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

type LibraDoc struct {
	SrcFileName string
	SrcFileSHA1 string
	SrcText     []string
}

func NewLibraDoc(srcFile string) (*LibraDoc, error) {
	sha, err := calcsha1(srcFile)
	if err != nil {
		return nil, err
	}
	text, err := getText(srcFile)
	if err != nil {
		return nil, err
	}
	return &LibraDoc{
		SrcFileName: srcFile,
		SrcFileSHA1: sha,
		SrcText:     text,
	}, nil
}

func (ld LibraDoc) String() string {
	return fmt.Sprintf("Libradoc {\n    SrcFileName: %s\n    SHA1: %s\n    nParagraphs: %d\n}",
		ld.SrcFileName, ld.SrcFileSHA1, len(ld.SrcText))
}

func getText(srcfile string) ([]string, error) {
	doc, err := document.Open(srcfile)
	if err != nil {
		log.WithField("Error", err).Fatal("Failed opening %s", srcfile)
		return nil, err
	}

	var paragraphs []document.Paragraph
	for _, p := range doc.Paragraphs() {
		paragraphs = append(paragraphs, p)
	}

	var text []string
	for _, p := range paragraphs {
		paragraph := ""
		for _, r := range p.Runs() {
			if strings.TrimSpace(r.Text()) != "" {
				paragraph += r.Text()
			}
		}
		if strings.TrimSpace(paragraph) != "" {
			text = append(text, paragraph)
		} else {
			text = append(text, "\n")
		}
	}
	return text, nil
}

func calcsha1(filepath string) (string, error) {
	var shaString string
	file, err := os.Open(filepath)
	if err != nil {
		log.WithFields(log.Fields{
			"file":  filepath,
			"error": err,
		}).Error("Failed opening file")
		return shaString, err
	}
	defer file.Close()

	hash := sha1.New()
	if _, err := io.Copy(hash, file); err != nil {
		log.WithField("Error", err).Fatal("Failed hash calculation")
		return shaString, err
	}
	shaString = hex.EncodeToString(hash.Sum(nil))
	return shaString, nil
}
