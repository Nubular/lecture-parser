package util

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/nubular/lecture-parser/parser"
)

type Config struct {
	Files []string
}

func mp3Exists(audiofile string) bool {
	config_file, _ := os.Open("config.json")
	defer config_file.Close()
	decoder := json.NewDecoder(config_file)
	configuration := Config{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	for i := 0; i < len(configuration.Files); i++ {
		if configuration.Files[i] == audiofile {
			return true
		}
	}
	return false

}

func updateConfig(audiofile string) {
	config_file, _ := os.Open("config.json")
	defer config_file.Close()
	decoder := json.NewDecoder(config_file)
	configuration := Config{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	configuration.Files = append(configuration.Files, audiofile)
	bs, err := json.Marshal(configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	ioutil.WriteFile("config.json", bs, 0644)
}

func ttsPolly(ssml string, outPath string, svc *polly.Polly) error {
	input := &polly.SynthesizeSpeechInput{OutputFormat: aws.String("mp3"), Text: aws.String(ssml), TextType: aws.String("ssml"), VoiceId: aws.String("Joanna")}
	output, err := svc.SynthesizeSpeech(input)
	// _ = input
	// output, err := []byte{65, 66, 67, 226, 130, 172}, nil
	if err != nil {
		log.Println("Got error calling SynthesizeSpeech:", err)
		return err
	}

	//create file
	outFile, err := os.Create(outPath)
	if err != nil {
		log.Println("Got error creating " + outPath + ":")
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, output.AudioStream)
	if err != nil {
		log.Println("Got error saving MP3:")
		return err
	}
	return nil
}

// CreateMP3 creates the tts mp3
func CreateMP3(outPath string, ssmlFrames []parser.Section, cacheFiles bool) error {

	if _, err := os.Stat(outPath); os.IsNotExist(err) {
		log.Println("Output dir not found. Creating at ", outPath)
		os.Mkdir(outPath, os.ModePerm)
	}

	//Creates session based on credentials & config present in .aws directory
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	//Flag to check if mp3 has already been created
	if cacheFiles {
		log.Println("Operating in dev mode. Using Cached files.")
		return nil
	}

	svc := polly.New(sess)
	// _ = svc
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for _, frame := range ssmlFrames {
		ssmlOutPath := filepath.Join(outPath, frame.FrameSrc.AudioSrc)
		wg.Add(1)
		go func(ssml string, audioOutPath string) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			default:
			}
			err := ttsPolly(ssml, audioOutPath, svc)

			// log.Println("Done Writing to ", audioOutPath)
			if err != nil {
				log.Println(err)
				cancel()
				return
			}
		}(frame.SSML, ssmlOutPath)

	}
	wg.Wait()
	return ctx.Err()
}
