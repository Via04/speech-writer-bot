package witai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

const WITAI_ENDPOINT = "https://api.wit.ai/speech"
const ENV_VAR = "WIT_TOKEN"

// MessageResponse - https://wit.ai/docs/http/20200513/#get__message_link
type MessageResponse struct {
	ID       string                     `json:"msg_id"`
	Text     string                     `json:"text"`
	Intents  []MessageIntent            `json:"intents"`
	Entities map[string][]MessageEntity `json:"entities"`
	Traits   map[string][]MessageTrait  `json:"traits"`
	IsFinal  bool                       `json:"is_final"`
}

// MessageEntity - https://wit.ai/docs/http/20200513/#get__message_link
type MessageEntity struct {
	ID         string                 `json:"id"`
	Name       string                 `json:"name"`
	Role       string                 `json:"role"`
	Start      int                    `json:"start"`
	End        int                    `json:"end"`
	Body       string                 `json:"body"`
	Value      string                 `json:"value"`
	Confidence float64                `json:"confidence"`
	Entities   []MessageEntity        `json:"entities"`
	Extra      map[string]interface{} `json:"-"`
}

// MessageTrait - https://wit.ai/docs/http/20200513/#get__message_link
type MessageTrait struct {
	ID         string                 `json:"id"`
	Value      string                 `json:"value"`
	Confidence float64                `json:"confidence"`
	Extra      map[string]interface{} `json:"-"`
}

// MessageIntent - https://wit.ai/docs/http/20200513/#get__message_link
type MessageIntent struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Confidence float64 `json:"confidence"`
}

func GetText(ctx context.Context, path string) (string, error) {
	token, ok := os.LookupEnv(ENV_VAR)
	if !ok {
		return "", errors.New("no token for witai")
	}
	contn, err := os.Open(path)
	if err != nil {
		return "", err
	}
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "POST", WITAI_ENDPOINT, contn)
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "audio/wav")
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	for {
		var answer = new(MessageResponse)

		err := decoder.Decode(answer)
		if answer.IsFinal {
			return answer.Text, nil
		}
		if err == io.EOF {
			// all done
			break
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}
	return "", errors.New("BAD REQUEST. NO ANSWER CHECK WAV")
}
