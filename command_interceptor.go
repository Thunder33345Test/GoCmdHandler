package CommandHandler

import (
	"github.com/bwmarrin/discordgo"
)

type InterceptorHandler struct {
	interceptors []Interceptor
}

func (i *InterceptorHandler) Intercept(s *discordgo.Session, m *discordgo.MessageCreate) (
	Interceptor Interceptor, command string, content string, success bool) {
	for _, interceptor := range i.interceptors {
		cmd, content, success := interceptor.Intercept(s, m)
		if success {
			return interceptor, cmd, content, success
		}
	}
	return nil, "", "", false
}

func (i *InterceptorHandler) Register(interceptor Interceptor) {
	i.interceptors = append(i.interceptors, interceptor)
}

func (i *InterceptorHandler) GetAll() []Interceptor {
	return i.interceptors
}

type Interceptor interface {
	//name of the interceptor(internal)
	Name() string
	//short prefix the user need to use inorder to trigger it(for explanation)
	Prefix(s *discordgo.Session, context CommandContext) string
	//brief description of this interceptor, ex "personal prefix"
	Description(s *discordgo.Session, context CommandContext) string
	//brief full explanation of this interceptor
	// ex "intercepts anything starting with * sent by thunder(id)"
	Explanation(s *discordgo.Session, context CommandContext) string
	//intercept test
	Intercept(s *discordgo.Session, m *discordgo.MessageCreate) (command string, content string, success bool)
	//show errors when command isn't found
	//example for bot tag prefix, explaining help command maybe?
	Error(s *discordgo.Session, m *discordgo.MessageCreate)
}
