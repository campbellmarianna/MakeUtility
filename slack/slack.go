package slack

import (
	"fmt"
	"strings"
	"github.com/gocolly/colly"
	"github.com/nlopes/slack"
	"math/rand"
)

/*
   NOTE: command_arg_1 and command_arg_2 represent optional parameteras that you define
   in the Slack API UI
*/
const helpMessage = "type in '@ScienceDaily skincare to rapidly get a truthy article on skincare'"

/*
   CreateSlackClient sets up the slack RTM (real-timemessaging) client library,
   initiating the socket connection and returning the client.
*/
func CreateSlackClient(apiKey string) *slack.RTM {
	api := slack.New(apiKey)
	rtm := api.NewRTM()
	go rtm.ManageConnection() // goroutine!
	return rtm
}

/*
   RespondToEvents waits for messages on the Slack client's incomingEvents channel,
   and sends a response when it detects the bot has been tagged in a message with @<botTag>.

   EDIT THIS FUNCTION IN THE SPACE INDICATED ONLY!
*/
func RespondToEvents(slackClient *slack.RTM) {
	for msg := range slackClient.IncomingEvents {
		fmt.Println("Event Received: ", msg.Type)
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			botTagString := fmt.Sprintf("<@%s> ", slackClient.GetInfo().User.ID)
			if !strings.Contains(ev.Msg.Text, botTagString) {
				continue
			}
			message := strings.Replace(ev.Msg.Text, botTagString, "", -1)

			// TODO: Make your bot do more than respond to a help command. See notes below.
			// Make changes below this line and add additional funcs to support your bot's functionality.
			// sendHelp is provided as a simple example. Your team may want to call a free external API
			// in a function called sendResponse that you'd create below the definition of sendHelp,
			// and call in this context to ensure execution when the bot receives an event.

			// START SLACKBOT CUSTOM CODE
			// ===============================================================
			sendSkincareArticle(slackClient, message, ev.Channel)
			sendHelp(slackClient, message, ev.Channel)
			
			// ===============================================================
			// END SLACKBOT CUSTOM CODE
		default:

		}
	}
}

// sendHelp is a working help message, for reference.
func sendHelp(slackClient *slack.RTM, message, slackChannel string) {
	if strings.ToLower(message) != "help" {
		return
	}
	slackClient.SendMessage(slackClient.NewOutgoingMessage(helpMessage, slackChannel))
}

// sendSkincareArticle when given th keyword skincare scraps the Science Daily website for articles on skin care and returns an article to the slack channel
func sendSkincareArticle(slackClient *slack.RTM, message, slackChannel string) {
	command := strings.ToLower(message)
	command = strings.TrimSpace(command)
	println("[RECEIVED] sendSkincareArticle:", command)
	// if strings.ToLower(command) == "skincare" {
	articleMessage := "Article Link"
	// get a random number
	randNum := rand.Intn(30)
	// create a counter for how many print statements you do
	outputCounter := 0
	// Instantiate default collector
	c := colly.NewCollector()

	// On every a element which has href attribute call callback
	c.OnHTML("h3 > a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Only those links are visited which are in AllowedDomains
		c.Visit(e.Request.AbsoluteURL(link)) // 35 links
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		outputCounter++
		fmt.Println("Visiting", r.URL.String())
		if outputCounter == randNum {
			articleMessage = r.URL.String()
		}
	})

	// Start scraping on https://www.sciencedaily.com/news/health_medicine/skin_care/
	c.Visit("https://www.sciencedaily.com/news/health_medicine/skin_care/")
	
	
	// if it is scrap an article from the website science daily using this link https://www.sciencedaily.com/news/health_medicine/skin_care/ so that everytime the command is received we return a new article 
	// EXPECTED OUTPUT - articleMessage := "https://www.sciencedaily.com/releases/2019/08/190815140834.htm"
	
	slackClient.SendMessage(slackClient.NewOutgoingMessage(articleMessage, slackChannel))
	// }
}
