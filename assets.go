package main

import (
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

var spritesheet = LoadImageFromPath("assets/spritesheet.png")
var test_bg = LoadImageFromPath("assets/backgrounds/test_bg.png")
var sewer_left = LoadImageFromPath("assets/backgrounds/sewer_left.png")
var sewer_entrance = LoadImageFromPath("assets/backgrounds/sewer_entrance.png")
var sewer_entrance_overlay = LoadImageFromPath("assets/backgrounds/sewer_entrance_overlay.png")
var sewer_middle = LoadImageFromPath("assets/backgrounds/sewer_middle.png")
var sewer_middle_with_top_connection = LoadImageFromPath("assets/backgrounds/sewer_middle_with_top_connection.png")
var sewer_going_up = LoadImageFromPath("assets/backgrounds/sewer_going_up.png")
var sewer_top = LoadImageFromPath("assets/backgrounds/sewer_top.png")
var factory_base = LoadImageFromPath("assets/backgrounds/factory_base.png")
var control_room = LoadImageFromPath("assets/backgrounds/control_room.png")
var control_room_overlay = LoadImageFromPath("assets/backgrounds/control_room_overlay.png")
var storage_b = LoadImageFromPath("assets/backgrounds/storage_b.png")
var storage_b_overlay = LoadImageFromPath("assets/backgrounds/storage_b_overlay.png")

var electrical_corridor = LoadImageFromPath("assets/backgrounds/electrical_corridor.png")
var electrical = LoadImageFromPath("assets/backgrounds/electrical.png")
var electrical_overlay = LoadImageFromPath("assets/backgrounds/electrical_overlay.png")

var final_production = LoadImageFromPath("assets/backgrounds/final_production.png")
var final_production_overlay = LoadImageFromPath("assets/backgrounds/final_production_overlay.png")

var pcb_manufacture = LoadImageFromPath("assets/backgrounds/pcb_manufacture.png")
var pcb_manufacture_overlay = LoadImageFromPath("assets/backgrounds/pcb_manufacture_overlay.png")

var metal_room = LoadImageFromPath("assets/backgrounds/metal_room.png")
var metal_room_overlay = LoadImageFromPath("assets/backgrounds/metal_room_overlay.png")

var pick_and_place = LoadImageFromPath("assets/backgrounds/pick_and_place.png")
var pick_and_place_overlay = LoadImageFromPath("assets/backgrounds/pick_and_place_overlay.png")

var vents = LoadImageFromPath("assets/backgrounds/vents.png")

var vignette = LoadImageFromPath("assets/vignette.png")
var vignette_red = LoadImageFromPath("assets/vignette_red.png")
var vignette_mild = LoadImageFromPath("assets/vignette_mild.png")
var target = LoadImageFromPath("assets/target.png")
var target_green = LoadImageFromPath("assets/target_green.png")
var target_red = LoadImageFromPath("assets/target_red.png")

var debug_wall = LoadImageFromPath("assets/debug_wall.png")

var hud_button_nine_slice = LoadImageFromPath("assets/ui/hud_button_nine_slice.png")
var hud_button_nine_slice_inverted = LoadImageFromPath("assets/ui/hud_button_nine_slice_inverted.png")
var button_nine_slice = LoadImageFromPath("assets/ui/button_nine_slice.png")
var button_nine_slice_inverted = LoadImageFromPath("assets/ui/button_nine_slice_inverted.png")
var button_nine_slice_disabled = LoadImageFromPath("assets/ui/button_nine_slice_disabled.png")
var box_nine_slice = LoadImageFromPath("assets/ui/box_nine_slice.png")

var hammer = LoadImageFromPath("assets/hammer.png")       // crafting symbol
var todo_list = LoadImageFromPath("assets/todo_list.png") // crafting symbol

var crafting_divider = LoadImageFromPath("assets/ui/crafting_divider.png")

var hotbar_slot = LoadImageFromPath("assets/ui/hotbar_slot.png")
var hotbar_slot_unselected = LoadImageFromPath("assets/ui/hotbar_slot_unselected.png")

var item_string = LoadImageFromPath("assets/items/string.png")
var item_rod = LoadImageFromPath("assets/items/rod.png")
var item_screwdriver = LoadImageFromPath("assets/items/screwdriver.png")
var item_hammer = LoadImageFromPath("assets/items/hammer.png")
var item_wire_cutters = LoadImageFromPath("assets/items/wire_cutters.png")
var item_bundle = LoadImageFromPath("assets/items/bundle.png")
var item_box = LoadImageFromPath("assets/box.png")
var item_auth_chip = LoadImageFromPath("assets/items/auth_chip.png")
var item_template_machine = LoadImageFromPath("assets/items/template_machine.png")
var item_template = LoadImageFromPath("assets/items/circuit_board_template.png")
var item_auth_card = LoadImageFromPath("assets/items/auth_card.png")
var recipe_slot = LoadImageFromPath("assets/ui/recipe_slot.png")
var recipe_slot_active = LoadImageFromPath("assets/ui/recipe_slot_active.png")

var item_hacking_usb = LoadImageFromPath("assets/items/hacking_usb.png")
var item_hacking_chip = LoadImageFromPath("assets/conveyor_items/hacking_chip.png")
var reprogramming_chip = LoadImageFromPath("assets/conveyor_items/reprogramming_chip.png")
var broken_chip = LoadImageFromPath("assets/conveyor_items/broken_chip.png")

var vent = LoadImageFromPath("assets/vent.png")
var wire = LoadImageFromPath("assets/wire.png")
var vent_open = LoadImageFromPath("assets/vent_open.png")
var vent_shadow = LoadImageFromPath("assets/vent_shadow.png")
var vent_open_shadow = LoadImageFromPath("assets/vent_open_shadow.png")
var machine = LoadImageFromPath("assets/machine.png")
var conveyor_left = LoadImageFromPath("assets/conveyor_left.png")
var conveyor_left_flipbook = LoadImageFromPath("assets/conveyor_left_flipbook.png")
var conveyor_down_flipbook = LoadImageFromPath("assets/conveyor_down_flipbook.png")

var zap_flipbook = LoadImageFromPath("assets/zap_flipbook.png")

var circuit_board_finished = LoadImageFromPath("assets/conveyor_items/circuit_board_finished.png")
var circuit_board_uncut = LoadImageFromPath("assets/conveyor_items/circuit_board_uncut.png")
var copper_sheet = LoadImageFromPath("assets/conveyor_items/copper_sheet.png")
var resin_board = LoadImageFromPath("assets/conveyor_items/resin_board.png")
var battery = LoadImageFromPath("assets/conveyor_items/battery.png")
var led = LoadImageFromPath("assets/conveyor_items/led.png")
var chip = LoadImageFromPath("assets/conveyor_items/chip.png")
var antenna = LoadImageFromPath("assets/conveyor_items/antenna.png")
var casing = LoadImageFromPath("assets/conveyor_items/casing.png")
var final_chip = LoadImageFromPath("assets/conveyor_items/final_chip.png")

var metal_sheet = LoadImageFromPath("assets/conveyor_items/metal_sheet.png")
var metal_sheet_held = LoadImageFromPath("assets/items/metal_sheet.png")

var right_wall_conveyor_overlay = LoadImageFromPath("assets/right_wall_conveyor_overlay.png")
var left_wall_conveyor_overlay = LoadImageFromPath("assets/left_wall_conveyor_overlay.png")
var top_wall_conveyor_overlay = LoadImageFromPath("assets/top_wall_conveyor_overlay.png")

var shootSound = [][]byte{
	ReadOggBytesFromPath("assets/sounds/explosionCrunch_000.ogg"),
	ReadOggBytesFromPath("assets/sounds/explosionCrunch_001.ogg"),
}

var concreteFootstepSound = [][]byte{
	ReadOggBytesFromPath("assets/sounds/footstep_concrete_000.ogg"),
	ReadOggBytesFromPath("assets/sounds/footstep_concrete_001.ogg"),
	ReadOggBytesFromPath("assets/sounds/footstep_concrete_002.ogg"),
	ReadOggBytesFromPath("assets/sounds/footstep_concrete_003.ogg"),
	ReadOggBytesFromPath("assets/sounds/footstep_concrete_004.ogg"),
}
var footstepSounds [][]byte = [][]byte{
	ReadOggBytesFromPath("assets/sounds/footstep00.ogg"),
	ReadOggBytesFromPath("assets/sounds/footstep01.ogg"),
	ReadOggBytesFromPath("assets/sounds/footstep02.ogg"),
	ReadOggBytesFromPath("assets/sounds/footstep03.ogg"),
	ReadOggBytesFromPath("assets/sounds/footstep04.ogg"),
	ReadOggBytesFromPath("assets/sounds/footstep05.ogg"),
	ReadOggBytesFromPath("assets/sounds/footstep06.ogg"),
	ReadOggBytesFromPath("assets/sounds/footstep07.ogg"),
	ReadOggBytesFromPath("assets/sounds/footstep08.ogg"),
	ReadOggBytesFromPath("assets/sounds/footstep09.ogg"),
}
var robotDeathSound [][]byte = [][]byte{
	ReadOggBytesFromPath("assets/sounds/impactPlate_heavy_000.ogg"),
	ReadOggBytesFromPath("assets/sounds/impactPlate_heavy_001.ogg"),
	ReadOggBytesFromPath("assets/sounds/impactPlate_heavy_002.ogg"),
}

var impactPlankSound = [][]byte{
	ReadOggBytesFromPath("assets/sounds/impactPlank_medium_000.ogg"),
	ReadOggBytesFromPath("assets/sounds/impactPlank_medium_001.ogg"),
	ReadOggBytesFromPath("assets/sounds/impactPlank_medium_002.ogg"),
	ReadOggBytesFromPath("assets/sounds/impactPlank_medium_003.ogg"),
	ReadOggBytesFromPath("assets/sounds/impactPlank_medium_004.ogg"),
}

func RandomSound(sounds [][]byte) []byte {
	L := len(sounds)
	index := rand.IntN(L)

	return sounds[index]
}

func PlaySound(context *audio.Context, sound []byte, volume float64) {
	sePlayer := context.NewPlayerFromBytes(sound)
	sePlayer.SetVolume(
		volume * 0.2,
	)
	sePlayer.Play()
}
