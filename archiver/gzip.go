package archiver

import (
	"bufio"
	"compress/gzip"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
)

func (a *Archiver) gzipFile(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Errorf("File not found: %s Error: %s", filePath, err.Error())
		return
	}
	defer file.Close()

	log.Infof("Archiving: %s", filePath)

	reader := bufio.NewReader(file)
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Errorf("Failure reading file: %s Error: %s", filePath, err.Error())
		return
	}

	gzipFile, err := os.Create(filePath + ".gz")
	if err != nil {
		log.Errorf("Failure creating file: %s.gz Error: %s", filePath, err.Error())
		return
	}
	defer gzipFile.Close()

	w := gzip.NewWriter(gzipFile)
	_, err = w.Write(content)
	if err != nil {
		log.Errorf("Failure writing content in file: %s.gz Error: %s", filePath, err.Error())
	}
	w.Close()

	err = os.Remove(filePath)
	if err != nil {
		log.Errorf("Failure deleting file: %s Error: %s", filePath, err.Error())
	}
}
