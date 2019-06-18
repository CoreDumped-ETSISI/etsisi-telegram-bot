package tts

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"

	"github.com/CoreDumped-ETSISI/etsisi-telegram-bot/state"

	tb "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/guad/commander"
)

func TtsCmd(ctx commander.Context) error {
	update := ctx.Arg("update").(state.Update)
	bot := update.State.Bot()

	msg := update.Message

	if msg.ReplyToMessage == nil || msg.ReplyToMessage.Text == "" {
		return nil
	}

	options := map[string]interface{}{
		"input": map[string]interface{}{
			"text": msg.ReplyToMessage.Text,
		},
		"voice": map[string]interface{}{
			"languageCode": "es-ES",
			"name":         "es-ES-Standard-A",
		},
		"audioConfig": map[string]interface{}{
			"audioEncoding": "MP3",
			"pitch":         15.2,
			"speakingRate":  1,
		},
	}

	b, _ := json.Marshal(options)
	buf := bytes.NewBuffer(b)

	resp, err := http.Post("https://texttospeech.googleapis.com/v1beta1/text:synthesize?key="+os.Getenv("TTS_API_KEY"), "application/json", buf)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var content map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&content)

	if err != nil {
		return err
	}

	encoded := content["audioContent"].(string)
	audio, err := base64.StdEncoding.DecodeString(encoded)

	if err != nil {
		return err
	}

	m := tb.NewVoiceUpload(msg.Chat.ID, tb.FileBytes{
		Name:  "Mensaje.mp3",
		Bytes: audio,
	})

	_, err = bot.Send(m)

	return err
}
