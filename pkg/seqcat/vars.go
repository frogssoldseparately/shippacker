package seqcat

// Not currently used
// const (
// 	GENERAL_SFX                 string = "100"
// 	AMBIENCE                    string = "101"
// 	TERMINA_FIELD               string = "102"
// 	CHASE                       string = "103"
// 	MAJORAS_THEME               string = "104"
// 	CLOCK_TOWER                 string = "105"
// 	STONE_TOWER_TEMPLE          string = "106"
// 	INV_STONE_TOWER_TEMPLE      string = "107"
// 	FAILURE_0                   string = "108"
// 	FAILURE_1                   string = "109"
// 	HAPPY_MASK_SALESMAN         string = "10A"
// 	SONG_OF_HEALING             string = "10B"
// 	SWAMP_REGION                string = "10C"
// 	ALIEN_INVASION              string = "10D"
// 	SWAMP_CRUISE                string = "10E"
// 	SHARPS_CURSE                string = "10F"
// 	GREAT_BAY_REGION            string = "110"
// 	IKANA_REGION                string = "111"
// 	DEKU_PALACE                 string = "112"
// 	MOUNTAIN_REGION             string = "113"
// 	PIRATES_FORTRESS            string = "114"
// 	CLOCK_TOWN_DAY_1            string = "115"
// 	CLOCK_TOWN_DAY_2            string = "116"
// 	CLOCK_TOWN_DAY_3            string = "117"
// 	FILE_SELECT                 string = "118"
// 	CLEAR_EVENT                 string = "119"
// 	ENEMY                       string = "11A"
// 	BOSS                        string = "11B"
// 	WOODFALL_TEMPLE             string = "11C"
// 	CLOCK_TOWN_MAIN_SEQUENCE    string = "11D"
// 	OPENING                     string = "11E"
// 	INSIDE_A_HOUSE              string = "11F"
// 	GAME_OVER                   string = "120"
// 	CLEAR_BOSS                  string = "121"
// 	GET_ITEM                    string = "122"
// 	CLOCK_TOWN_DAY_2_PTR        string = "123"
// 	GET_HEART                   string = "124"
// 	TIMED_MINI_GAME             string = "125"
// 	GORON_RACE                  string = "126"
// 	MUSIC_BOX_HOUSE             string = "127"
// 	FAIRY_FOUNTAIN              string = "128"
// 	ZELDAS_LULLABY              string = "129"
// 	ROSA_SISTERS                string = "12A"
// 	OPEN_CHEST                  string = "12B"
// 	MARINE_RESEARCH_LAB         string = "12C"
// 	GIANTS_THEME                string = "12D"
// 	SONG_OF_STORMS              string = "12E"
// 	ROMANI_RANCH                string = "12F"
// 	GORON_VILLAGE               string = "130"
// 	MAYORS_OFFICE               string = "131"
// 	OCARINA_EPONA               string = "132"
// 	OCARINA_SUNS                string = "133"
// 	OCARINA_TIME                string = "134"
// 	OCARINA_STORM               string = "135"
// 	ZORA_HALL                   string = "136"
// 	GET_NEW_MASK                string = "137"
// 	MINI_BOSS                   string = "138"
// 	GET_SMALL_ITEM              string = "139"
// 	ASTRAL_OBSERVATORY          string = "13A"
// 	CAVERN                      string = "13B"
// 	MILK_BAR                    string = "13C"
// 	ZELDA_APPEAR                string = "13D"
// 	SARIAS_SONG                 string = "13E"
// 	GORON_GOAL                  string = "13F"
// 	HORSE                       string = "140"
// 	HORSE_GOAL                  string = "141"
// 	INGO                        string = "142"
// 	KOTAKE_POTION_SHOP          string = "143"
// 	SHOP                        string = "144"
// 	OWL                         string = "145"
// 	SHOOTING_GALLERY            string = "146"
// 	OCARINA_SOARING             string = "147"
// 	OCARINA_HEALING             string = "148"
// 	INVERTED_SONG_OF_TIME       string = "149"
// 	SONG_OF_DOUBLE_TIME         string = "14A"
// 	SONATA_OF_AWAKENING         string = "14B"
// 	GORON_LULLABY               string = "14C"
// 	NEW_WAVE_BOSSA_NOVA         string = "14D"
// 	ELEGY_OF_EMPTINESS          string = "14E"
// 	OATH_TO_ORDER               string = "14F"
// 	SWORD_TRAINING_HALL         string = "150"
// 	OCARINA_LULLABY_INTRO       string = "151"
// 	LEARNED_NEW_SONG            string = "152"
// 	BREMEN_MARCH                string = "153"
// 	BALLAD_OF_THE_WIND_FISH     string = "154"
// 	SONG_OF_SOARING             string = "155"
// 	MILK_BAR_DUPLICATE          string = "156"
// 	FINAL_HOURS                 string = "157"
// 	MIKAU_RIFF                  string = "158"
// 	MIKAU_FINALE                string = "159"
// 	FROG_SONG                   string = "15A"
// 	OCARINA_SONATA              string = "15B"
// 	OCARINA_LULLABY             string = "15C"
// 	OCARINA_NEW_WAVE            string = "15D"
// 	OCARINA_ELEGY               string = "15E"
// 	OCARINA_OATH                string = "15F"
// 	MAJORAS_LAIR                string = "160"
// 	OCARINA_LULLABY_INTRO_PTR   string = "161"
// 	OCARINA_GUITAR_BASS_SESSION string = "162"
// 	PIANO_SESSION               string = "163"
// 	INDIGO_GO_SESSION           string = "164"
// 	SNOWHEAD_TEMPLE             string = "165"
// 	GREAT_BAY_TEMPLE            string = "166"
// 	NEW_WAVE_SAXOPHONE          string = "167"
// 	NEW_WAVE_VOCAL              string = "168"
// 	MAJORAS_WRATH               string = "169"
// 	MAJORAS_INCARNATION         string = "16A"
// 	MAJORAS_MASK                string = "16B"
// 	BASS_PLAY                   string = "16C"
// 	DRUMS_PLAY                  string = "16D"
// 	PIANO_PLAY                  string = "16E"
// 	IKANA_CASTLE                string = "16F"
// 	GATHERING_GIANTS            string = "170"
// 	KAMARO_DANCE                string = "171"
// 	CREMIA_CARRIAGE             string = "172"
// 	KEATON_QUIZ                 string = "173"
// 	END_CREDITS                 string = "174"
// 	OPENING_LOOP                string = "175"
// 	TITLE_THEME                 string = "176"
// 	DUNGEON_APPEAR              string = "177"
// 	WOODFALL_CLEAR              string = "178"
// 	SNOWHEAD_CLEAR              string = "179"
// 	MUSIC_BOX_HOUSE_INTERIOR    string = "17A"
// 	INTO_THE_MOON               string = "17B"
// 	GOODBYE_GIANT               string = "17C"
// 	TATL_AND_TAEL               string = "17D"
// 	MOONS_DESTRUCTION           string = "17E"
// 	END_CREDITS_SECOND_HALF     string = "17F"
// )

var fanfareCategories = []string{"8", "9", "10"}

// Not currently used
// var categoryReplacements = map[string]string{
// 	TERMINA_FIELD:          "0",
// 	STONE_TOWER_TEMPLE:     "2",
// 	INV_STONE_TOWER_TEMPLE: "2",
// 	SWAMP_REGION:           "0-1",
// 	SWAMP_CRUISE:           "4-3-6",
// 	IKANA_REGION:           "0-1",
// 	DEKU_PALACE:            "1-0",
// 	CLOCK_TOWN_DAY_1:       "1",
// 	FILE_SELECT:            "6-3",
// 	WOODFALL_TEMPLE:        "2",
// 	INSIDE_A_HOUSE:         "3",
// 	BOSS:                   "7",
// 	GAME_OVER:              "9",
// 	GORON_RACE:             "4",
// 	MUSIC_BOX_HOUSE:        "3-4-6",
// 	ROMANI_RANCH:           "0-1",
// 	GORON_VILLAGE:          "1-3",
// 	MINI_BOSS:              "7",
// 	CAVERN:                 "2",
// 	SARIAS_SONG:            "6-1-3-4-0",
// 	SHOP:                   "3",
// 	SHOOTING_GALLERY:       "3-4",
// 	FINAL_HOURS:            "5",
// 	PIANO_PLAY:             "8",
// 	IKANA_CASTLE:           "2",
// }
