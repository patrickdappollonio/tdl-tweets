package main

import (
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"
)

func Tweet(s *Stream, ctx context.Context) error {
	// Check if twitterapi is nil
	if twitterapi == nil {
		return fmt.Errorf("twitterapi: not set!")
	}

	// Check if the stream is nil
	if s == nil {
		return fmt.Errorf("stream: pointer is nil")
	}

	// Check if stream is valid
	if s.StreamID == 0 {
		return fmt.Errorf("stream: the stream provided is empty")
	}

	// Set the Anaconda's API transport compatible
	// with AppEngine
	twitterapi.HttpClient.Transport = &urlfetch.Transport{Context: ctx}

	// Convert image to base64-encoded string
	b64, err := ConvertImage(ctx, s.PreviewURL)

	// Check if it was possible
	if err != nil {
		return err
	}

	// Try uploading the image to the Twitter API
	image, err := twitterapi.UploadMedia(b64)

	// Check for existence of an error
	if err != nil {
		return err
	}

	// Create a value placeholder
	// and pass the media_id acquired in the upload
	values := url.Values{}
	values.Set("media_ids", strconv.FormatInt(image.MediaID, 10))

	// Send the tweet to twitter
	_, err = twitterapi.PostTweet(createMessage(s), values)

	// Check if it was possible
	if err != nil {
		return err
	}

	return nil
}

// Message formats that'll be updated with their proper values
var messages = []string{
	`#TheDivision se juega en el canal de {user}: {url}`,
	`{user} está jugando #TheDivision: {url}`,
	`Ganas de ver #TheDivision? Ve al canal de {user}: {url}`,
	`{user} está en vivo jugando #TheDivision: {url}`,
	`En vivo: {user} juega #TheDivision: {url}`,
	`Te invitamos al canal de {user} a ver #TheDivision: {url}`,
	`Hype de #TheDivision en el canal de {user}: {url}`,
	`Y @TheDivisionGame se juega en el Canal de {user}: {url}`,
	`Estamos todos en el Canal de {user} viendo #TheDivision: {url}`,
	`Únete al stream de {user} que juega @TheDivisionGame: {url}`,
	`Pasa un buen momento en el canal de {user} que juega #TheDivision: {url}`,
}

// For known people, we convert Twitch usernames to @names
var conversions = map[string]string{
	"patrickdap":       "marlex",
	"hawk12fps":        "hawk12fps",
	"mrjutsu":          "mrjutsu",
	"mrprobeta":        "mrprobeta",
	"zeromexico":       "ZeroMexico",
	"boga_xp":          "BoGA_xP",
	"monstergmer":      "MonsterGmer",
	"victorzcre":       "zucre_",
	"vmt85":            "Mantis_V8",
	"glenfy":           "GlenfyGames",
	"ardashe":          "Ardashe",
	"elsuperwtf88":     "ElsupeRWTF88",
	"lrockyhd":         "lRockyHD",
	"orihalcon_tsuyoi": "CarlosOrihalcon",
	"lyberion":         "lyberion",
	"dopiko":           "DopikoS",
	"balrickps4":       "Balrick_Twitch",
	"duendepablo":      "DuendeGaming",
	"soxx18":           "BysoXx",
	"dm_daedalus":      "Dm_Daedalus",
	"lupesparza":       "lupesparzatv",
	"redreichelgames":  "RedReichelGames",
	"mariyolo1":        "Cenandoconlobos",
	"xboxmexico":       "XboxMexico",
	"enemykitty":       "enemykitty",
	"imburundi":        "imburundi",
	"carlosbanano":     "_CarlosBanano",
	"arrasatorr":       "javierarrasa",
	"riimpo":           "riimpo",
	"neisrosver":       "NeisRosver",
	"dk1nghd":          "danik1ng",
	"tavinhokun":       "tavinho_kun",
	"danstertv":        "danstter",
	"kraviuz":          "kraviuz",
	"cordonesdesata2":  "hdcordoneshd",
	"dopeyferes":       "DopeyFeres",
	"leduend":          "duend",
	"alexander_cavali": "AlexanderCavali",
	"kurojesterhead":   "TanyaGelotte",
	"javidemon":        "JaViDeMoN_Co",
	"zenixrevenge":     "zenixrevenge",
}

func convertToTwitterHandler(channel string) string {
	channelLower := strings.ToLower(channel)

	if handler, ok := conversions[channelLower]; ok {
		return fmt.Sprintf("@%v", handler)
	}

	return channel
}

func getRandomMessageFormat() string {
	rand.Seed(int64(time.Now().Nanosecond()))
	return messages[rand.Intn(len(messages))]
}

func createMessage(s *Stream) string {
	rep := strings.NewReplacer(
		`{user}`, convertToTwitterHandler(s.DisplayName),
		`{url}`, s.URL,
	)

	tweetstr := rep.Replace(getRandomMessageFormat())

	if strings.HasPrefix(tweetstr, "@") {
		return fmt.Sprintf(".%s", tweetstr)
	}

	return tweetstr
}

func GetFollowerIDs(ctx context.Context) ([]int64, error) {
	// Check if twitterapi is nil
	if twitterapi == nil {
		return []int64{}, fmt.Errorf("twitterapi: not set!")
	}

	// Pass default AppEngine transport
	twitterapi.HttpClient.Transport = &urlfetch.Transport{Context: ctx}

	// Holder for users
	var userIDs []int64
	var err error

	// Fetch a list of all followers
	pages := twitterapi.GetFollowersIdsAll(nil)

	// Iterate over the channel retrieving the values
	for page := range pages {
		// Check if there was an error
		if page.Error != nil {
			err = page.Error
			break
		}

		// If no error, append the new list
		userIDs = append(userIDs, page.Ids...)
	}

	return userIDs, err
}
