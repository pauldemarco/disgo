package bot

import (
	"fmt"

	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/httpserver"
	"github.com/disgoorg/disgo/internal/tokenhelper"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/disgo/sharding"
	"github.com/disgoorg/log"
	"github.com/pkg/errors"
)

func DefaultConfig(gatewayHandlers map[discord.GatewayEventType]GatewayEventHandler, httpHandler HTTPServerEventHandler) *Config {
	return &Config{
		Logger:                 log.Default(),
		EventManagerConfigOpts: []EventManagerConfigOpt{WithGatewayHandlers(gatewayHandlers), WithHTTPServerHandler(httpHandler)},
	}
}

// Config lets you configure your Client instance
// Config is the core.Client config used to configure everything
type Config struct {
	Logger log.Logger

	RestClient           rest.Client
	RestClientConfigOpts []rest.ConfigOpt
	RestServices         rest.Services

	EventManager           EventManager
	EventManagerConfigOpts []EventManagerConfigOpt

	Gateway           gateway.Gateway
	GatewayConfigOpts []gateway.ConfigOpt

	ShardManager           sharding.ShardManager
	ShardManagerConfigOpts []sharding.ConfigOpt

	HTTPServer           httpserver.Server
	HTTPServerConfigOpts []httpserver.ConfigOpt

	Caches          cache.Caches
	CacheConfigOpts []cache.ConfigOpt

	MemberChunkingManager MemberChunkingManager
	MemberChunkingFilter  MemberChunkingFilter
}

type ConfigOpt func(config *Config)

func (c *Config) Apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

// WithLogger lets you inject your own logger implementing log.Logger
//goland:noinspection GoUnusedExportedFunction
func WithLogger(logger log.Logger) ConfigOpt {
	return func(config *Config) {
		config.Logger = logger
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithRestClient(restClient rest.Client) ConfigOpt {
	return func(config *Config) {
		config.RestClient = restClient
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithRestClientConfigOpts(opts ...rest.ConfigOpt) ConfigOpt {
	return func(config *Config) {
		config.RestClientConfigOpts = append(config.RestClientConfigOpts, opts...)
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithRestServices(restServices rest.Services) ConfigOpt {
	return func(config *Config) {
		config.RestServices = restServices
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithEventManager(eventManager EventManager) ConfigOpt {
	return func(config *Config) {
		config.EventManager = eventManager
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithEventManagerConfigOpts(opts ...EventManagerConfigOpt) ConfigOpt {
	return func(config *Config) {
		config.EventManagerConfigOpts = append(config.EventManagerConfigOpts, opts...)
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithEventListeners(eventListeners ...EventListener) ConfigOpt {
	return func(config *Config) {
		config.EventManagerConfigOpts = append(config.EventManagerConfigOpts, WithListeners(eventListeners...))
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithGateway(gateway gateway.Gateway) ConfigOpt {
	return func(config *Config) {
		config.Gateway = gateway
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithGatewayConfigOpts(opts ...gateway.ConfigOpt) ConfigOpt {
	return func(config *Config) {
		config.GatewayConfigOpts = append(config.GatewayConfigOpts, opts...)
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithShardManager(shardManager sharding.ShardManager) ConfigOpt {
	return func(config *Config) {
		config.ShardManager = shardManager
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithShardManagerConfigOpts(opts ...sharding.ConfigOpt) ConfigOpt {
	return func(config *Config) {
		config.ShardManagerConfigOpts = append(config.ShardManagerConfigOpts, opts...)
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithHTTPServer(httpServer httpserver.Server) ConfigOpt {
	return func(config *Config) {
		config.HTTPServer = httpServer
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithHTTPServerConfigOpts(opts ...httpserver.ConfigOpt) ConfigOpt {
	return func(config *Config) {
		config.HTTPServerConfigOpts = append(config.HTTPServerConfigOpts, opts...)
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithCaches(caches cache.Caches) ConfigOpt {
	return func(config *Config) {
		config.Caches = caches
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithCacheConfigOpts(opts ...cache.ConfigOpt) ConfigOpt {
	return func(config *Config) {
		config.CacheConfigOpts = append(config.CacheConfigOpts, opts...)
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithMemberChunkingManager(memberChunkingManager MemberChunkingManager) ConfigOpt {
	return func(config *Config) {
		config.MemberChunkingManager = memberChunkingManager
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithMemberChunkingFilter(memberChunkingFilter MemberChunkingFilter) ConfigOpt {
	return func(config *Config) {
		config.MemberChunkingFilter = memberChunkingFilter
	}
}

func BuildClient(token string, config Config, gatewayEventHandlerFunc func(client Client) gateway.EventHandlerFunc, httpServerEventHandlerFunc func(client Client) httpserver.EventHandlerFunc, os string, name string, github string, version string) (Client, error) {
	if token == "" {
		return nil, discord.ErrNoBotToken
	}
	id, err := tokenhelper.IDFromToken(token)
	if err != nil {
		return nil, errors.Wrap(err, "error while getting application id from token")
	}
	client := &clientImpl{
		token:  token,
		logger: config.Logger,
	}

	// TODO: figure out how we handle different application & client ids
	client.applicationID = *id
	client.clientID = *id

	if config.RestClient == nil {
		// prepend standard user-agent. this can be overridden as it's appended to the front of the slice
		config.RestClientConfigOpts = append([]rest.ConfigOpt{
			rest.WithUserAgent(fmt.Sprintf("DiscordBot (%s, %s)", github, version)),
			rest.WithLogger(client.logger),
		}, config.RestClientConfigOpts...)

		config.RestClient = rest.NewClient(client.token, config.RestClientConfigOpts...)
	}

	if config.RestServices == nil {
		config.RestServices = rest.NewServices(config.RestClient)
	}
	client.restServices = config.RestServices

	if config.EventManager == nil {
		config.EventManager = NewEventManager(client, config.EventManagerConfigOpts...)
	}
	client.eventManager = config.EventManager

	if config.Gateway == nil && config.GatewayConfigOpts != nil {
		var gatewayRs *discord.Gateway
		gatewayRs, err = client.restServices.GatewayService().GetGateway()
		if err != nil {
			return nil, err
		}

		config.GatewayConfigOpts = append([]gateway.ConfigOpt{
			gateway.WithGatewayURL(gatewayRs.URL),
			gateway.WithLogger(client.logger),
			gateway.WithOS(os),
			gateway.WithBrowser(name),
			gateway.WithDevice(name),
		}, config.GatewayConfigOpts...)

		config.Gateway = gateway.New(token, gatewayEventHandlerFunc(client), config.GatewayConfigOpts...)
	}
	client.gateway = config.Gateway

	if config.ShardManager == nil && config.ShardManagerConfigOpts != nil {
		var gatewayBotRs *discord.GatewayBot
		gatewayBotRs, err = client.restServices.GatewayService().GetGatewayBot()
		if err != nil {
			return nil, err
		}

		shardIDs := make([]int, gatewayBotRs.Shards-1)
		for i := 0; i < gatewayBotRs.Shards-1; i++ {
			shardIDs[i] = i
		}

		config.ShardManagerConfigOpts = append([]sharding.ConfigOpt{
			sharding.WithShardCount(gatewayBotRs.Shards),
			sharding.WithShards(shardIDs...),
			sharding.WithGatewayConfigOpts(
				gateway.WithGatewayURL(gatewayBotRs.URL),
				gateway.WithLogger(client.logger),
			),
			sharding.WithLogger(client.logger),
		}, config.ShardManagerConfigOpts...)

		config.ShardManager = sharding.New(token, gatewayEventHandlerFunc(client), config.ShardManagerConfigOpts...)
	}
	client.shardManager = config.ShardManager

	if config.HTTPServer == nil && config.HTTPServerConfigOpts != nil {
		config.HTTPServerConfigOpts = append([]httpserver.ConfigOpt{
			httpserver.WithLogger(client.logger),
		}, config.HTTPServerConfigOpts...)

		config.HTTPServer = httpserver.New(httpServerEventHandlerFunc(client), config.HTTPServerConfigOpts...)
	}
	client.httpServer = config.HTTPServer

	if config.MemberChunkingManager == nil {
		config.MemberChunkingManager = NewMemberChunkingManager(client, config.MemberChunkingFilter)
	}
	client.memberChunkingManager = config.MemberChunkingManager

	if config.Caches == nil {
		config.Caches = cache.NewCaches(config.CacheConfigOpts...)
	}
	client.caches = config.Caches

	return client, nil
}
