package porter

import "strings"

// replacer replaces the texture names with the correct names.
var replacer = strings.NewReplacer(
	"chainmail_layer_1", "chain_1",
	"chainmail_layer_2", "chain_2",

	"diamond_layer_1", "diamond_1",
	"diamond_layer_2", "diamond_2",

	"gold_layer_1", "gold_1",
	"gold_layer_2", "gold_2",

	"iron_layer_1", "iron_1",
	"iron_layer_2", "iron_2",

	"leather_layer_1", "leather_1",
	"leather_layer_2", "leather_2",

	"netherite_layer_1", "netherite_1",
	"netherite_layer_2", "netherite_2",
)
