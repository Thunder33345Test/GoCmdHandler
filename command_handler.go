package CommandHandler

import (
	"Hex/RateLimiter"
	"Hex/StateStore"
	"github.com/bwmarrin/discordgo"
	"time"
)

type CommandHandler struct {
	//responsible of intercepting commands, return cleaned content(with no prefixes/suffixes)
	Interceptor  InterceptorHandler
	CommandGroup []CommandGroup
	StateStore   *StateStore.Collection
}

func NewCommandHandler(state *StateStore.Collection) *CommandHandler {
	ch := CommandHandler{
		StateStore: state,
	}
	return &ch
}

func (c *CommandHandler) Register(cg CommandGroup) {
	c.CommandGroup = append(c.CommandGroup, cg)
}

func (c *CommandHandler) Unregister(cg CommandGroup) {
	for i, cg2 := range c.CommandGroup {
		if cg2 == cg {
			c.CommandGroup = append(c.CommandGroup[:i], c.CommandGroup[(i+1):]...)
		}
	}
}

func (c *CommandHandler) Lookup(cmd string) (Command, CommandGroup, bool) {
	for _, cg := range c.CommandGroup {
		ci, ok := cg.Lookup(cmd)
		if ok {
			return ci, cg, ok
		}
	}
	return nil, nil, false
}

//Raw command event handler
//this should be the ONLY function registered as interface between discord api
func (c *CommandHandler) HandleEvent(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID || m.Author.Bot {
		return
	}
	if c.StateStore.User(m.Author.ID).UserFlags.BlackList {
		return
	}

	interceptor, commandName, content, success := c.Interceptor.Intercept(s, m)
	if !success {
		return
	}

	cmd, group, res := c.Lookup(commandName)
	if !res {
		interceptor.Error(s, m)
		return
	}

	if oo, ok := cmd.(OwnerOnly); ok {
		if oo.OwnerOnly() && !c.StateStore.User(m.Author.ID).UserFlags.BotOwner {
			return
		}
	}

	cc := CommandContext{
		Sender:          m.Author,
		Guild:           m.GuildID,
		Channel:         m.ChannelID,
		Content:         content,
		OriginalMessage: m.Message,
		Interceptor:     interceptor,
		CommandGroup:    group,
	}
	hc := HexContext{StateStore: c.StateStore, CommandHandler: c}

	//defer func() {
	//	if r := recover(); r != nil {
	//		txt := fmt.Sprintf("Panick: %T %v\nOn:%v\n", r, r, cc)
	//		fmt.Println(txt)
	//		//_, _ = s.ChannelMessageSend(cc.Channel, fmt.Sprintf("Something went wrong\nPanick: %T %v", r, r))
	//	}
	//}()

	if !c.StateStore.User(m.Author.ID).UserFlags.IgnoreRateLimit {
		l, _ := c.StateStore.User(m.Author.ID).Limiter.GetOrCreate("global.limit", RateLimiter.NewLimiterGroupAnd(
			RateLimiter.NewLimiter(time.Second*3, 2),
			RateLimiter.NewLimiter(time.Minute, 40),
		))
		if !l.TryAddTally() {
			return
		}
	}

	if rl, ok := cmd.(RateLimit); ok {
		var l RateLimiter.ILimiter
		if c.StateStore.User(m.Author.ID).UserFlags.IgnoreRateLimit {
			l = RateLimiter.NewInfiniteLimiter()
		} else {
			l, _ = c.StateStore.User(m.Author.ID).Limiter.GetOrCreate("cmd."+cmd.Name(), rl.RateLimit())
		}
		cc.Limiter = l
	}

	cmd.Executor(s, cc, hc)

}
