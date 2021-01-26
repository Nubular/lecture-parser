package main

import(
    "fmt"
    "os"
    "encoding/json"
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/polly"
    "io"
    "io/ioutil"
)

type Config struct{
    Files []string
}

func mp3Exists(audiofile string) bool{
    config_file,_ := os.Open("config.json");
    defer config_file.Close();
    decoder := json.NewDecoder(config_file)
    configuration := Config{}
    err := decoder.Decode(&configuration)
    if err != nil {
    fmt.Println("error:", err)
    }
    for i:=0;i<len(configuration.Files);i++{
        if(configuration.Files[i]==audiofile){
            return true;
        }
    }
    return false;
    
}

func updateConfig(audiofile string){
    config_file,_ := os.Open("config.json");
    defer config_file.Close();
    decoder := json.NewDecoder(config_file)
    configuration := Config{}
    err := decoder.Decode(&configuration)
    if err != nil {
    fmt.Println("error:", err)
    }
    configuration.Files = append(configuration.Files,audiofile);
    bs,err:= json.Marshal(configuration)
    if err!= nil{
        fmt.Println("error:",err)
    }
    ioutil.WriteFile("config.json",bs,0644)
}

func main(){
    
    fname := "slide3_1" 
    s := "Please read this mp3" //Text to be converted to speech
    audiofile := fname + ".mp3"
    //Creates session based on credentials & config present in .aws directory 
    sess := session.Must(session.NewSessionWithOptions(session.Options{
        SharedConfigState: session.SharedConfigEnable,
    }))
    
    //Flag to check if mp3 has already been created
    if mp3Exists(audiofile){
        fmt.Println("MP3 for this slide has already been generated.")

    } else {   
        //Making request to Polly
        //Old code for 3000 characters - using synthesize speech
        svc := polly.New(sess)
        input := &polly.SynthesizeSpeechInput{OutputFormat: aws.String("mp3"), Text: aws.String(s), VoiceId: aws.String("Joanna")}
        output, err := svc.SynthesizeSpeech(input)
    
        if err != nil {
            fmt.Println("Got error calling SynthesizeSpeech:")
            fmt.Print(err.Error())
            os.Exit(1)
        }

        //create file
        name := audiofile
        outFile, err := os.Create(name)
        if err != nil {
            fmt.Println("Got error creating " + name + ":")
            fmt.Print(err.Error())
            os.Exit(1)
        }
        defer outFile.Close()
        _, err = io.Copy(outFile, output.AudioStream)
        if err != nil {
            fmt.Println("Got error saving MP3:")
            fmt.Print(err.Error())
            os.Exit(1)
        }
        updateConfig(audiofile);
    }    
}
    
    
    
  
