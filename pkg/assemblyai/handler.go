package assemblyai

// This package contains needed function to transcribe speech to text using assemblyAI

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

// constant to AssemblyAI endpoint, need to get transcribeID before transcribtion. It's not guranteed that we will id,
// because only 5 parallel transcribtions allowed in free version.
const TRANSCRIPT_URL = "https://api.assemblyai.com/v2/transcript"
const ENV_VAR = "ASSEMBLY_TOKEN"

// get TranscribeID, this function intended to work with timeout context
func getEndpoint(ctx context.Context, audioUrl string) (string, error) {
	var answer map[string]string
	assemblyToken, ok := os.LookupEnv(ENV_VAR)
	if !ok {
		return "", errors.New("no token for assemblyai passed")
	}
	content := map[string]string{"audio_url": audioUrl}
	contentJson, err := json.Marshal(&content)
	if err != nil {
		return "", fmt.Errorf("cannot marshal audiourl %s", audioUrl)
	}
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "POST", TRANSCRIPT_URL, bytes.NewBuffer(contentJson))
	if err != nil {
		panic(err)
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("authorization", assemblyToken)
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	json.NewDecoder(res.Body).Decode(&answer)
	return answer["id"], nil
}

// Get text from audio file by its url
func GetText(ctx context.Context, audioUrl string) (string, error) {
	var workUrl string
	client := http.Client{}
	transcribeID, err := getEndpoint(ctx, audioUrl)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		return "", err
	}
	workUrl = TRANSCRIPT_URL + "/" + transcribeID
	req, _ := http.NewRequestWithContext(ctx, "GET", workUrl, nil)
	req.Header.Set("content-type", "application/json")
	req.Header.Set("authorization", ENV_VAR)
	res, err := client.Do(req)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		return "", err
	}
	defer res.Body.Close()
	var result map[string]string
	json.NewDecoder(res.Body).Decode(&result)

	// Check status and print the transcribed text
	if result["status"] == "completed" {
		fmt.Println(result["text"])
		return result["text"], nil
	}
	return "", errors.New("no completed status from assemblyai")
}
