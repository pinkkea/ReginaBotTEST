package main

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

var errInvalidFormat = errors.New("invalid format")
var attendeeRole = "1136878110646743170"

// https://yourbasic.org/golang/time-change-convert-location-timezone/
// TimeIn returns the time in UTC if the name is "" or "UTC".
// It returns the local time if the name is "Local".
// Otherwise, the name is taken to be a location name in
// the IANA Time Zone database, such as "Africa/Lagos".
func TimeIn(t time.Time, name string) (time.Time, error) {
	loc, err := time.LoadLocation(name)
	if err == nil {
		t = t.In(loc)
	}
	return t, err
}

func main() {
	discord, err := discordgo.New("Bot <TOKEN>")

	if err != nil {
		log.Fatal(err)
	}

	// Add event handler
	discord.AddHandler(newMessage)

	// Open session
	discord.Open()
	defer discord.Close()

	// Run until code is terminated
	fmt.Println("Bot running...")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func helpMessage() (s string) {
	var message = `
$color <hex code>: I'll set your color to <hex code> (example: #FFFFFF).
$help: Make me repeat this message for some ungodly reason
$rsvp: RSVP to WWD24!

$date: Get today's date
$dogepoint: <:dogekek:1135750882584182925> 👉
$dumphim: <:dumphim:1136426428565569567>
$horse: 🐴
$horses: Several horses
$limit: Talk about the limit
$mathletes: Mathletes opinion
$skillissue: I'll say "skill issue"
$uck: Sometimes people get up to this I guess
$tacobell: I'll express my feelings about Taco Bell
$talkshit: I'll say "post fit"
$waluigi <@user>: Get the waluigi of someone else (or get your own waluigi by not tagging anyone)
$wednesday: If it's Wednesday I'll let you know
$white: <:white:1136047355997720747>
`
	return message
}

func wednesdayMessage() (s string) {
	// NYC supremacy
	t, _ := TimeIn(time.Now(), "America/New_York")
	weekday := t.Weekday()
	if int(weekday) == 3 {
		return "https://giphy.com/gifs/filmeditor-mean-girls-movie-3otPozZKy1ALqGLoVG"
	} else {
		return "It's not Wednesday numbnuts it is " + weekday.String()
	}
}

func dateMessage() (s string) {
	t, _ := TimeIn(time.Now(), "America/New_York")
	dateString := t.Format("January 2")
	if dateString != "October 3" {
		return dateString
	}
	// It's October 3rd.
	return "https://tenor.com/view/crush-diary-october3rd-mean-girls-lindsay-lohan-gif-9906172"
}

func checkForStrings(discord *discordgo.Session, message *discordgo.MessageCreate) {
	// OK we don't want to spam so let's kind of sort these by funniest highest to lowest.
	// and send only one per message (even if has two.)
	msg := strings.ToLower(message.Content)
	if strings.Contains(msg, "fetch") {
		// Stop trying to make fetch happen
		discord.ChannelMessageSend(message.ChannelID,
			"https://tenor.com/view/fetch-mean-girls-gif-19691105")
	} else {
		marxStrings := []string{"marx", "capital", "landlord", "rich", "big natural",
			"communis", "comrade", "commie"}
		for _, str := range marxStrings {
			if strings.Contains(msg, str) {
				discord.MessageReactionAdd(message.ChannelID, message.ID,
					"marx:1159218841189101668")
				return
			}
		}

	}
}

func newMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {

	// Ignore bot messaage
	if message.Author.ID == discord.State.User.ID {
		return
	}

	if len(message.Content) == 0 {
		return
	}

	// Don't bother with all the below switch logic if there isn't a $ in the front.
	// Conversely, don't check for strings if there's maybe a command.
	if message.Content[0:1] != "$" {
		checkForStrings(discord, message)
		return
	}

	tokens := strings.Fields(message.Content)

	// Respond to messages
	switch tokens[0] {
	// Somewhat more involved commands
	case "$color":
		createColorRole(message, discord, tokens)
	case "$help":
		discord.ChannelMessageSend(message.ChannelID, helpMessage())
	case "$rsvp":
		discord.GuildMemberRoleAdd(message.GuildID, message.Author.ID, attendeeRole)
		discord.ChannelMessageSend(message.ChannelID,
			"Thanks, "+message.Author.Username+" for RSVPing to WWD24! Hope to see you there!")

	// Memes
	case "$date":
		discord.ChannelMessageSend(message.ChannelID, dateMessage())
	case "$dogepoint":
		discord.ChannelMessageSend(message.ChannelID, "<:dogekek:1135750882584182925> 👉")
	case "$dumphim":
		discord.ChannelMessageSend(message.ChannelID, "<:dumphim:1136426428565569567>")
	case "$horse":
		discord.ChannelMessageSend(message.ChannelID, "🐴")
	case "$horses":
		discord.ChannelMessageSend(message.ChannelID,
			"🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴")
	case "$limit":
		// the limit does not exist
		discord.ChannelMessageSend(message.ChannelID,
			"https://tenor.com/view/mean-girls-karen-gif-9300840")
	case "$mathletes":
		// you can't join mathletes. it's social suicide.
		discord.ChannelMessageSend(message.ChannelID,
			"https://y.yarn.co/ea1dd776-80ed-43fb-a53d-9e77520bf781_text.gif")
	case "$skillissue":
		discord.ChannelMessageSend(message.ChannelID, "skill issue")
	case "$uck":
		// sucking dick and cock
		discord.ChannelMessageSend(message.ChannelID,
			"https://cdn.discordapp.com/attachments/1136094668136915066/1143333955555303568/image0.gif")
	case "$tacobell":
		discord.ChannelMessageSend(message.ChannelID, "https://i.imgur.com/TkZbs3J.png")
	case "$talkshit":
		discord.ChannelMessageSend(message.ChannelID, "post fit")
	case "$wednesday":
		discord.ChannelMessageSend(message.ChannelID, wednesdayMessage())
	case "$waluigi":
		waluigi(message, discord, tokens)
	case "$white":
		discord.ChannelMessageSend(message.ChannelID, "<:white:1136047355997720747>")
	default:
		// There wasn't a command so let's just check for strings.
		checkForStrings(discord, message)
		return
	}
}

// so far all this does is print the profile picture of either the message author
// (if no one mentioned) or the user mentioned. I can't figure out how to invert stuff.
func waluigi(message *discordgo.MessageCreate, discord *discordgo.Session, tokens []string) {
	if len(tokens) != 1 {
		// convert e.g. <@425544164365565962> to 425544164365565962
		if len(tokens[1]) < 3 {
			discord.ChannelMessageSend(message.ChannelID,
				`Literally who are you even talking about? Type "$waluigi <@user> to get the waluigi of that user."`)
			return
		}
		userId := tokens[1][2 : len(tokens[1])-1]
		mentionedUser, err := discord.User(userId)
		if err != nil {
			discord.ChannelMessageSend(message.ChannelID,
				`Literally who are you even talking about? Type "$waluigi <@user> to get the waluigi of that user."`)
			return
		}
		img, e := discord.UserAvatarDecode(mentionedUser)
		if e != nil {
			discord.ChannelMessageSend(message.ChannelID,
				"idk why but i can't get that user's pfp")
			return
		}
		wah := Invert(img)
		f, err := os.Create("img.png")
		if err != nil {
			fmt.Println("image save failed")
			return
		}
		defer f.Close()
		if err = png.Encode(f, wah); err != nil {
			fmt.Printf("failed to encode: %v", err)
			return
		}
		file, errr := os.Open("img.png")
		if errr != nil {
			fmt.Println("image save failed")
			return
		}
		discord.ChannelFileSend(message.ChannelID, "waluigi.png", file)
		os.Remove("img.png")
	}
}

func createColorRole(message *discordgo.MessageCreate, discord *discordgo.Session, tokens []string) {
	if len(tokens) == 1 {
		discord.ChannelMessageSend(message.ChannelID,
			`Are you fucking with me? You need to type "$color <hex code>" (e.g. "$color #FFFFFF")`)
		return
	} else if !strings.Contains(tokens[1], "#") {
		discord.ChannelMessageSend(message.ChannelID,
			`Are you fucking with me? You need to type "$color <hex code>" (e.g. "$color #FFFFFF"). YOU NEED THE # IN THERE TOO.`)
		return
	}

	c := tokens[1]
	fmt.Println(c)
	_, e := ParseHexColorFast(c)

	if e != nil { //if the color is a bad hex code throw this error
		discord.ChannelMessageSend(message.ChannelID,
			`I couldn't find that color. Try again using a hex code (e.g. "$color #FFFFFF")`)
		fmt.Println("fail 1")
		return
	}
	//the role params struct needs a decimal INT for some reason, so we need to convert it here.
	cInt, convErr := strconv.ParseInt(c[1:], 16, 64)
	cIntPoi := int(cInt)
	if convErr != nil { //if the int parser fails for some godforsaken reason
		discord.ChannelMessageSend(message.ChannelID,
			"I had issues creating that color. Is it even real?? Try again or contact @synanasthesia")
		fmt.Println("fail 2")
		return
	}
	newRole := discordgo.RoleParams{Name: c, Color: &cIntPoi} //this creates the role parameters - currently, the name is set to just the color string and color is obvious.
	role, er := discord.GuildRoleCreate(message.GuildID, &newRole)
	if er != nil {
		discord.ChannelMessageSend(message.ChannelID,
			"Sorry, I couldn't create a role with that color. Please try again or contact @synanasthesia")
		fmt.Println("fail 3")
		return
	}
	//somewhere here we either need to remove old roles or reorder the roles.
	discord.GuildMemberRoleAdd(message.GuildID, message.Author.ID, role.ID)
	discord.ChannelMessageSend(message.ChannelID, "Done! How's that?")
}

// shout out to https://stackoverflow.com/questions/54197913/parse-hex-string-to-image-color
func ParseHexColorFast(s string) (c color.RGBA, err error) {
	c.A = 0xff

	if s[0] != '#' {
		return c, errInvalidFormat
	}

	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10
		case b >= 'A' && b <= 'F':
			return b - 'A' + 10
		}
		err = errInvalidFormat
		return 0
	}

	switch len(s) {
	case 7:
		c.R = hexToByte(s[1])<<4 + hexToByte(s[2])
		c.G = hexToByte(s[3])<<4 + hexToByte(s[4])
		c.B = hexToByte(s[5])<<4 + hexToByte(s[6])
	case 4:
		c.R = hexToByte(s[1]) * 17
		c.G = hexToByte(s[2]) * 17
		c.B = hexToByte(s[3]) * 17
	default:
		err = errInvalidFormat
	}
	return
}

func Invert(img image.Image) *image.NRGBA {
	dst := image.NewNRGBA(image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()))
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c := img.At(x, y)
			r, g, b, a := c.RGBA()
			wC := color.RGBA{R: uint8(255 - r), G: uint8(255 - g), B: uint8(255 - b), A: uint8(a)}
			dst.Set(x, y, wC)
		}
	}
	return dst
}
