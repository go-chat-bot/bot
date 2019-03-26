package treta

import (
	"math/rand"
	"strings"

	"github.com/go-chat-bot/bot"
)

const (
	msgInvalidAmountOfParams = "Invalid amount of parameters"
	msgInvalidParam          = "Invalid parameter"
)

var (
	quotes = map[string][]string{
		"APPLE": {
			"Apple Inc. is not just a computer/portable device company, but at its inner core a philosophy. It's a philosophy of life, of living, of being alive, of stayin' alive, and of livin' la vida loca. It is a way of thinking and consuming overpriced monochrome technology that's designed with elegance.",
			"Apple decided to discontinue the Macbook Air, though there are rumours of its reintroduction, as several enthusiastic customers have threatened suicide if the line is permanently discontinued.",
			"Every once in a while, some brave soul dares to ask: \"Why do I keep buying Apple products? They're shit and extremely overpriced, but I just keep buying! WHY?!\" Apple's response is always some along the lines of: \"We own you.\"",
			"The iPhone is the culmination of several years of research and development at Apple, in to how they could further extort money from customers while maintaining an almost Big Brother style control over their device.",
			"Safari is worse than the fucking Internet Explorer Version 6.",
			"OS X in some ways is actually worse than Windows to program for. Their file system is complete and utter crap, which is scary.",
		},
		"DELPHI": {
			"Delphi. Now there's a name I haven't heard in a long time.",
			"Access Violation at address 69696969 in module 'Treta.exe'. Read of address 00000666.",
			"Delphi supports the best way to develop iOS applications",
			"It’s not difficult to read and listen about the wonders of Embarcadero DataSnap technology around the world.",
			"Delphi was, and remains, vastly superior to anyone developer tools, in that users can actually produce working programs with it.",
		},
		"GO": {
			"Go is the official programming language of the eXtreme Go Horse",
			"If you're looking for a language optimized for your problem domain, Golang is not the language for you.",
			"I don't know that Golang is a great language.",
			"Go don't have classes/constructors, but we have to reinvent them... with much worse practices.",
			"Oh Go! You so crazy!",
			"Good, Good... Let the Golang flow through you.",
		},
		"LINUX": {
			"The Linux philosophy is 'Laugh in the face of danger'. Oops. Wrong One. 'Do it yourself'. Yes, that's it.",
			"Software is like sex: it's better when it's free.",
			"My name is Linus, and I am your God.",
			"I think the OpenBSD crowd is a bunch of masturbating monkeys, in that they make such a big deal about concentrating on security to the point where they pretty much admit that nothing else matters to them.",
			"Nvidia, fuck you!",
		},
		"JAVA": {
			"You're using Java? Well there's your problem.",
			"I had a problem so I thought to use Java. Now I have a ProblemFactory.",
			"Many individual Java programmers claim that it is the very best technology available, particularly when they don't know anything else.",
			"Java Performance? You must be joking!",
			"It is said that Java was an idea of God to show to Humans how stupid they were",
		},
		"JAVASCRIPT": {
			"Javascript is not funny at all",
			"JavaScript, why don't you work?",
			"Brace yourself. A new Javascript framework is coming.",
			"JavaScript... Whoops! Maybe you were looking for Java?",
			"JavaScript is a computer language for writing ineffectual computer viruses (interruptions to web surfing that will annoy the user without completely ruining his computer)",
		},
		"PYTHON": {
			"We'll can do cool things... even with Python",
			"No one has been able to live programming with Python",
			"Python is the best programming language in the world... for kids to play and have fun.",
		},
		"RUBY": {
			"Can Rails Scale? NOOOOO!",
			"Why is Ruby so slow?",
			"I hate managing inventory and the game drops more weapon than the rails can handle the requests",
			"Ruby on Rails? Pleaaase. Do you even code, bro?",
			"The classic Hello, world! program is really easy with Ruby. You just need to know the name of the gem you want to install.",
			"Python is known for its clear, readable, and regular syntax. Ruby code is vandalism!",
			"Python is better than Ruby",
			"even PHP is better than Ruby",
			"Ruby may do something completely useless and have infinite ways of doing something completely useless.",
			"I've hit this a few times in Ruby and it bugged me like crazy. But then I grew up, learned Python, and dealt with it.",
			"Do your best to program, not just uses Ruby.",
		},
		"VIM": {
			"Emacs > VIM",
			"Sublime Text > VIM",
			"even Notepad > VIM",
			"VIM... Why can't I quit you?!",
			"Vim Is Too Mainstream. I'm Switching To Emacs",
		},
		"WINDOWS": {
			"Why I love Windows: Keyboard not responding. Press any key to continue.",
			"Why I love Windows: A system call that should never fail has failed.",
			"Why I love Windows: Bluescreen has performed an illegal operation. Bluescreen must be closed.",
			"Why I love Windows: An error occurred whilst trying to load the previous error.",
			"Help and Support Error: Windows cannot open Help and Support because a system service is not running. To fix this problems, start the service named Help and Support",
			"Windows is the collective name for a series of failures that began in 1983 as a means of reversing the stagnation of the computer hardware market.",
			"I mean, it's obvious, isn't it? Windows seems perfectly clear and simple to use, but it crashes with the slightest pressure, or sometimes breaks inexplicably.",
			"Windows was officially confirmed to work correctly on i386, X86-64, IA64, ARM - it crashes on all of them. Undesired productivity boost when run under VirtualBox on Ubuntu.",
			"Microsoft isn't evil, they just make really crappy operating systems.",
			"Hoping the problem magically goes away by ignoring it is the “microsoft approach to programming” and should never be allowed.",
		},
	}
)

func treta(command *bot.Cmd) (string, error) {
	var key string
	switch len(command.Args) {
	case 0:
		key = randKey()
	case 1:
		key = strings.ToUpper(command.Args[0])
	default:
		return msgInvalidAmountOfParams, nil
	}

	q, found := quotes[key]
	if !found {
		return msgInvalidParam, nil
	}
	return q[rand.Intn(len(q))], nil
}

func randKey() string {
	keys := make([]string, 0, len(quotes))
	for k := range quotes {
		keys = append(keys, k)
	}
	return keys[rand.Intn(len(keys))]
}

func init() {
	bot.RegisterCommand(
		"treta",
		"sowing discord",
		"",
		treta)
}
