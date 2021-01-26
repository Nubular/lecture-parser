package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
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

func main() {
	err := createMP3(true)
	if err != nil {
		fmt.Println(err)
	}
}

func createMP3(cacheFiles bool) error {

	fname := "slide3_1"
	s := "Please read this mp3" //Text to be converted to speech
	audiofile := fname + ".mp3"
	//Creates session based on credentials & config present in .aws directory
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	//Flag to check if mp3 has already been created
	if cacheFiles {
		log.Println("Operating in dev mode. Using Cached files.")

	} else {
		//Making request to Polly
		//Old code for 3000 characters - using synthesize speech
		svc := polly.New(sess)
		input := &polly.SynthesizeSpeechInput{OutputFormat: aws.String("mp3"), Text: aws.String(s), VoiceId: aws.String("Joanna")}
		output, err := svc.SynthesizeSpeech(input)

		if err != nil {
			log.Panic("Got error calling SynthesizeSpeech:", err)
			return err
		}

		//create file
		name := audiofile
		outFile, err := os.Create(name)
		if err != nil {
			log.Panic("Got error creating " + name + ":")
			return err
		}
		defer outFile.Close()
		_, err = io.Copy(outFile, output.AudioStream)
		if err != nil {
			log.Panic("Got error saving MP3:")
			return err
		}
	}
	return nil
}
