package main

import(
    "fmt"
    "os"
    "encoding/json"
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/polly"
    "github.com/aws/aws-sdk-go/service/s3"
    "github.com/aws/aws-sdk-go/service/s3/s3manager"
    "time"
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
    
    fname := "slide2_1" 
    s := "Please read this mp3" //Text to be converted to speech
    checkfile := fname + ".mp3"
    //Creates session based on credentials & config present in .aws directory 
    sess := session.Must(session.NewSessionWithOptions(session.Options{
        SharedConfigState: session.SharedConfigEnable,
    }))
    
    //Flag to check if mp3 has already been created
    if mp3Exists(checkfile){
        fmt.Println("MP3 for this slide has already been generated.")

    } else {   
        //Making request to Polly
        svc := polly.New(sess)
        input := &polly.StartSpeechSynthesisTaskInput{OutputFormat: aws.String("mp3"), OutputS3BucketName:aws.String("se-project-mp3-files") ,Text: aws.String(s), VoiceId: aws.String("Joanna")}
        response,err := svc.StartSpeechSynthesisTask(input)
        if err != nil{
            fmt.Println("Here")
            fmt.Println(err.Error());
            os.Exit(1);
        }
        fmt.Println(response.SynthesisTask)
        fmt.Printf("%T",response);
        taskID := *response.SynthesisTask.TaskId +".mp3"
        item :=  fname + ".mp3"
        //Polly sends output to S3 bucket - filename in S3 bucket is same as taskID
        time.Sleep(20*time.Second);
        //Downloading from S3 -  - audio file named after the slide (fname)
        downloader := s3manager.NewDownloader(sess);
        file, err := os.Create(item)

        if err != nil {
            fmt.Println("Error HERE")
            fmt.Println(err);
            os.Exit(1);
        }
        //Downloading the file onto system
        defer file.Close()
        numBytes, err := downloader.Download(file,
            &s3.GetObjectInput{
                Bucket: aws.String("se-project-mp3-files"),
                Key:    aws.String(taskID),
            })
        if err != nil {
            fmt.Println("This is the error");
            fmt.Println(err);
            os.Exit(1);
        }

        fmt.Println("Downloaded", file.Name(), numBytes, "bytes")
        updateConfig(item);
    }    
}
    
    
    
    //Old code for 3000 characters - using synthesize speech
    // svc := polly.New(sess)
    // input := &polly.SynthesizeSpeechInput{OutputFormat: aws.String("mp3"), Text: aws.String(s), VoiceId: aws.String("Joanna")}
    // output, err := svc.SynthesizeSpeech(input)
    
    // if err != nil {
    //     fmt.Println("Got error calling SynthesizeSpeech:")
    //     fmt.Print(err.Error())
    //     os.Exit(1)
    // }

    // //create file
    // name := "random_audio.mp3"
    // outFile, err := os.Create(name)
    // if err != nil {
    //     fmt.Println("Got error creating " + name + ":")
    //     fmt.Print(err.Error())
    //     os.Exit(1)
    // }
    // defer outFile.Close()
    // _, err = io.Copy(outFile, output.AudioStream)
    // if err != nil {
    //     fmt.Println("Got error saving MP3:")
    //     fmt.Print(err.Error())
    //     os.Exit(1)
    // }

