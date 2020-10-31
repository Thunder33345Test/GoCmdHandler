package CommandHandler

import (
	"Hex/RateLimiter"
	"Hex/StateStore"
	"github.com/bwmarrin/discordgo"
)

type Command interface {
	Name() string
	Alias() []string
	Executor(ses *discordgo.Session, context CommandContext, hexContext HexContext)
}

type RateLimit interface {
	RateLimit() RateLimiter.ILimiter
}

type PermissionMetadata struct {
	//todo overwrite this with internal permission management in a latter date
	BotPermission  int
	UserPermission int
}

type HasPermission interface {
	Permission() PermissionMetadata
}

type OwnerOnly interface {
	OwnerOnly() bool
}

type Visible interface {
	Visible() bool
}

type HasHelp interface {
	Help() string
	//Arguments   func() string //generate the arguments help maybe
}

type CommandContext struct {
	Sender  *discordgo.User
	Guild   string
	Channel string
	//maybe also interceptor used
	//intercepted content, which should have necessary arguments, and excluded prefixes
	Content string
	//the original content should not be relied upon
	//for all we know, it could be a base64-ed message that needs to be intercepted first
	OriginalMessage *discordgo.Message
	//the interceptor used
	Interceptor Interceptor
	//the command group of the command
	CommandGroup CommandGroup
	//nullable limiter
	Limiter RateLimiter.ILimiter
}

type HexContext struct {
	StateStore     *StateStore.Collection
	CommandHandler *CommandHandler
}

func (h *HexContext) User(id string) *StateStore.UserStore {
	return h.StateStore.User(id)
}

//should be moved to state store
/*type UserFlags struct {
	//should tbe user be allowed to execute bot owner commands
	BotOwner bool
	//should the user bypass user permission check
	RootMode bool
	//should the user be free from rate limit
	IgnoreRateLimit bool
	//should tbe user be ignored
	BlackList bool
}*/

//future explaining faults
//panic or return idk
/*type Explainable interface {
	Error() string
	Explain() string
}*/
