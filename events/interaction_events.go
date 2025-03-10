package events

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
)

// InteractionResponderFunc is a function that can be used to respond to a discord.Interaction.
type InteractionResponderFunc func(responseType discord.InteractionResponseType, data discord.InteractionResponseData, opts ...rest.RequestOpt) error

// InteractionCreate indicates that a new interaction has been created.
type InteractionCreate struct {
	*GenericEvent
	discord.Interaction
	Respond InteractionResponderFunc
}

// Guild returns the guild that the interaction happened in if it happened in a guild.
// If the interaction happened in a DM, it returns nil.
// This only returns cached guilds.
func (e *InteractionCreate) Guild() (discord.Guild, bool) {
	if e.GuildID() != nil {
		return e.Client().Caches().Guild(*e.GuildID())
	}
	return discord.Guild{}, false
}

// ApplicationCommandInteractionCreate is the base struct for all application command interaction create events.
type ApplicationCommandInteractionCreate struct {
	*GenericEvent
	discord.ApplicationCommandInteraction
	Respond InteractionResponderFunc
}

// Guild returns the guild that the interaction happened in if it happened in a guild.
// If the interaction happened in a DM, it returns nil.
// This only returns cached guilds.
func (e *ApplicationCommandInteractionCreate) Guild() (discord.Guild, bool) {
	if e.GuildID() != nil {
		return e.Client().Caches().Guild(*e.GuildID())
	}
	return discord.Guild{}, false
}

// CreateMessage responds to the interaction with a new message.
func (e *ApplicationCommandInteractionCreate) CreateMessage(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeCreateMessage, messageCreate, opts...)
}

// DeferCreateMessage responds to the interaction with a "bot is thinking..." message which should be edited later.
func (e *ApplicationCommandInteractionCreate) DeferCreateMessage(ephemeral bool, opts ...rest.RequestOpt) error {
	var data discord.InteractionResponseData
	if ephemeral {
		data = discord.MessageCreate{Flags: discord.MessageFlagEphemeral}
	}
	return e.Respond(discord.InteractionResponseTypeDeferredCreateMessage, data, opts...)
}

// CreateModal responds to the interaction with a new modal.
func (e *ApplicationCommandInteractionCreate) CreateModal(modalCreate discord.ModalCreate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeModal, modalCreate, opts...)
}

// ComponentInteractionCreate indicates that a new component interaction has been created.
type ComponentInteractionCreate struct {
	*GenericEvent
	discord.ComponentInteraction
	Respond InteractionResponderFunc
}

// Guild returns the guild that the interaction happened in if it happened in a guild.
// If the interaction happened in a DM, it returns nil.
// This only returns cached guilds.
func (e *ComponentInteractionCreate) Guild() (discord.Guild, bool) {
	if e.GuildID() != nil {
		return e.Client().Caches().Guild(*e.GuildID())
	}
	return discord.Guild{}, false
}

// CreateMessage responds to the interaction with a new message.
func (e *ComponentInteractionCreate) CreateMessage(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeCreateMessage, messageCreate, opts...)
}

// DeferCreateMessage responds to the interaction with a "bot is thinking..." message which should be edited later.
func (e *ComponentInteractionCreate) DeferCreateMessage(ephemeral bool, opts ...rest.RequestOpt) error {
	var data discord.InteractionResponseData
	if ephemeral {
		data = discord.MessageCreate{Flags: discord.MessageFlagEphemeral}
	}
	return e.Respond(discord.InteractionResponseTypeDeferredCreateMessage, data, opts...)
}

// UpdateMessage responds to the interaction with updating the message the component is from.
func (e *ComponentInteractionCreate) UpdateMessage(messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeUpdateMessage, messageUpdate, opts...)
}

// DeferUpdateMessage responds to the interaction with nothing.
func (e *ComponentInteractionCreate) DeferUpdateMessage(opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeDeferredUpdateMessage, nil, opts...)
}

// CreateModal responds to the interaction with a new modal.
func (e *ComponentInteractionCreate) CreateModal(modalCreate discord.ModalCreate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeModal, modalCreate, opts...)
}

// AutocompleteInteractionCreate indicates that a new autocomplete interaction has been created.
type AutocompleteInteractionCreate struct {
	*GenericEvent
	discord.AutocompleteInteraction
	Respond InteractionResponderFunc
}

// Guild returns the guild that the interaction happened in if it happened in a guild.
// If the interaction happened in a DM, it returns nil.
// This only returns cached guilds.
func (e *AutocompleteInteractionCreate) Guild() (discord.Guild, bool) {
	if e.GuildID() != nil {
		return e.Client().Caches().Guild(*e.GuildID())
	}
	return discord.Guild{}, false
}

// Result responds to the interaction with a slice of choices.
func (e *AutocompleteInteractionCreate) Result(choices []discord.AutocompleteChoice, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeApplicationCommandAutocompleteResult, discord.AutocompleteResult{Choices: choices}, opts...)
}

// ModalSubmitInteractionCreate indicates that a new modal submit interaction has been created.
type ModalSubmitInteractionCreate struct {
	*GenericEvent
	discord.ModalSubmitInteraction
	Respond InteractionResponderFunc
}

// Guild returns the guild that the interaction happened in if it happened in a guild.
// If the interaction happened in a DM, it returns nil.
// This only returns cached guilds.
func (e *ModalSubmitInteractionCreate) Guild() (discord.Guild, bool) {
	if e.GuildID() != nil {
		return e.Client().Caches().Guild(*e.GuildID())
	}
	return discord.Guild{}, false
}

// CreateMessage responds to the interaction with a new message.
func (e *ModalSubmitInteractionCreate) CreateMessage(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeCreateMessage, messageCreate, opts...)
}

// DeferCreateMessage responds to the interaction with a "bot is thinking..." message which should be edited later.
func (e *ModalSubmitInteractionCreate) DeferCreateMessage(ephemeral bool, opts ...rest.RequestOpt) error {
	var data discord.InteractionResponseData
	if ephemeral {
		data = discord.MessageCreate{Flags: discord.MessageFlagEphemeral}
	}
	return e.Respond(discord.InteractionResponseTypeDeferredCreateMessage, data, opts...)
}

// UpdateMessage responds to the interaction with updating the message the component is from.
func (e *ModalSubmitInteractionCreate) UpdateMessage(messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeUpdateMessage, messageUpdate, opts...)
}

// DeferUpdateMessage responds to the interaction with nothing.
func (e *ModalSubmitInteractionCreate) DeferUpdateMessage(opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeDeferredUpdateMessage, nil, opts...)
}
