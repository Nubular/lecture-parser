package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

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

	absPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	outPath := filepath.Join(absPath, "output")

	err = createMP3(`<speak>
You can also write custom ssml which works with specific platforms!
You say, <phoneme alphabet="ipa" ph="pɪˈkɑːn">pecan</phoneme>. 
I say, <phoneme alphabet="ipa" ph="ˈpi.kæn">pecan</phoneme>.
Sometimes it can be useful to <prosody volume="loud">increase the volume 
for a specific speech.</prosody>
 
</speak>`, filepath.Join(outPath, "test2.mp3"), true)
	if err != nil {
		fmt.Println(err)
	}
}

func createMP3(ssml string, outPath string, cacheFiles bool) error {

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
		input := &polly.SynthesizeSpeechInput{OutputFormat: aws.String("mp3"), Text: aws.String(ssml), TextType: aws.String("ssml"), VoiceId: aws.String("Joanna")}
		output, err := svc.SynthesizeSpeech(input)

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
	}
	return nil
}
