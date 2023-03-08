package converter

// uses ffmpeg to convert audio from * to wav. Of course required to have ffmpeg to be installed in path

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type converter struct {
	path    string
	filein  string
	fileout string
	sep     string
}

// Base function to start convertion. Creates new converter.
// File to convert must be located in workdirectory in folder data.
func New(filein string, fileout string) (*converter, error) {
	path, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	sep := string(os.PathSeparator)
	path = path + sep + "data"
	return &converter{path: path, filein: filein, fileout: fileout, sep: sep}, nil
}

// Convert file specified in New method of format ogg into Wav format and output it into outFile.
// outFile path is relative and goes from workdir/data
func (c converter) convertOggToWav() error {
	cmd := exec.Command("ffmpeg", "-i", c.path+c.sep+c.filein, c.path+c.sep+c.fileout)
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	return nil
}

// Convert file specified in New method of format mp3 into Wav format and output it into outFile.
// outFile path is relative and goes from workdir/data
func (c converter) convertMp3ToWav() error {
	cmd := exec.Command("ffmpeg", "-i", c.path+c.sep+c.filein, "-acodec",
		"pcm_s16le", "-ac", "1", "ar", "16000", c.path+c.sep+c.fileout)
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	return nil
}

// convert specified file to Wav format with mono sound
func (c converter) ConvertToWav() error {
	ext := filepath.Ext(c.filein)
	ext = strings.ToLower(ext)
	switch ext {
	case ".ogg", ".oga":
		c.convertOggToWav()
	case ".mp3":
		c.convertMp3ToWav()
	}
	return nil
}

// delete files before convertion or delete both if deleteAll flag set
func (c converter) Purge(deleteAll bool) {
	os.Remove(c.path + c.sep + c.filein)
	if deleteAll {
		os.Remove(c.path + c.sep + c.fileout)
	}
}
