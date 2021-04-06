package api

import log "github.com/sirupsen/logrus"

// DisgoBuilder allows you to create a Disgo client through a series of methods
type DisgoBuilder interface {
	SetLogLevel(level log.Level) DisgoBuilder
	SetToken(token string) DisgoBuilder
	SetIntents(intents Intents) DisgoBuilder
	SetVoiceDispatchInterceptor(VoiceDispatchInterceptor) DisgoBuilder
	SetEventManager(eventManager EventManager) DisgoBuilder
	AddEventListeners(eventsListeners ...EventListener) DisgoBuilder
	SetWebhookServer(webhookServer WebhookServer) DisgoBuilder
	SetWebhookServerProperties(listenURL string, listenPort int, publicKey string) DisgoBuilder
	SetRestClient(restClient RestClient) DisgoBuilder
	SetCache(cache Cache) DisgoBuilder
	SetMemberCachePolicy(memberCachePolicy MemberCachePolicy) DisgoBuilder
	SetMessageCachePolicy(messageCachePolicy MessageCachePolicy) DisgoBuilder
	SetCacheFlags(cacheFlags CacheFlags) DisgoBuilder
	SetGateway(gateway Gateway) DisgoBuilder
	Build() (Disgo, error)
}
