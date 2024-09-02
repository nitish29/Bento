package evil

import "main/bot"

// silly thing to verify interfaces as implimented
var _ bot.SpokeMessageCreateHandler = &Evil{}
var _ bot.SpokeMessageReactionAddHandler = &Evil{}
