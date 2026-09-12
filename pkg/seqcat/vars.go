package seqcat

const (
	MMCatField       string = "0"
	MMCatTown        string = "1"
	MMCatDungeon     string = "2"
	MMCatIndoor      string = "3"
	MMCatMinigame    string = "4"
	MMCatAction      string = "5"
	MMCatCalm        string = "6"
	MMCatBoss        string = "7"
	MMCatFanGetItem  string = "8"
	MMCatFanGameOver string = "9"
	MMCatFanClear    string = "10"
	MMCatTitle       string = "16"

	MMSongTerminaField             string = "102"
	MMSongChase                    string = "103"
	MMSongMajorasTheme             string = "104"
	MMSongClockTower               string = "105"
	MMSongStoneTowerTemple         string = "106"
	MMSongInvertedStoneTowerTemple string = "107"
	MMSongFailure0                 string = "108"
	MMSongFailure1                 string = "109"
	MMSongHappyMaskSalesman        string = "10A"
	MMSongSongOfHealing            string = "10B"
	MMSongSwampRegion              string = "10C"
	MMSongAlienInvasion            string = "10D"
	MMSongSwampCruise              string = "10E"
	MMSongSharpsCurse              string = "10F"
	MMSongGreatBayRegion           string = "110"
	MMSongIkanaRegion              string = "111"
	MMSongDekuPalace               string = "112"
	MMSongMountainRegion           string = "113"
	MMSongPiratesFortress          string = "114"
	MMSongClockTownDay1            string = "115"
	MMSongClockTownDay2            string = "116"
	MMSongClockTownDay3            string = "117"
	MMSongFileSelect               string = "118"
	MMSongClearEvent               string = "119"
	MMSongEnemy                    string = "11A"
	MMSongBoss                     string = "11B"
	MMSongWoodfallTemple           string = "11C"
	MMSongClockTownMainSequence    string = "11D"
	MMSongOpening                  string = "11E"
	MMSongInsideAHouse             string = "11F"
	MMSongGameOver                 string = "120"
	MMSongClearBoss                string = "121"
	MMSongGetItem                  string = "122"
	MMSongClockTownDay2Ptr         string = "123"
	MMSongGetHeart                 string = "124"
	MMSongTimedMinigame            string = "125"
	MMSongGoronRace                string = "126"
	MMSongMusicBoxHouse            string = "127"
	MMSongFairyFountain            string = "128"
	MMSongZeldasLullaby            string = "129"
	MMSongRosaSisters              string = "12A"
	MMSongOpenChest                string = "12B"
	MMSongMarineResearchLab        string = "12C"
	MMSongGiantsTheme              string = "12D"
	MMSongSongOfStorms             string = "12E"
	MMSongRomaniRanch              string = "12F"
	MMSongGoronVillage             string = "130"
	MMSongMayorsOffice             string = "131"
	MMSongOcarinaEpona             string = "132"
	MMSongOcarinaSuns              string = "133"
	MMSongOcarinaTime              string = "134"
	MMSongOcarinaStorm             string = "135"
	MMSongZoraHall                 string = "136"
	MMSongGetNewMask               string = "137"
	MMSongMiniboss                 string = "138"
	MMSongGetSmallItem             string = "139"
	MMSongAstralObservatory        string = "13A"
	MMSongCavern                   string = "13B"
	MMSongMilkBar                  string = "13C"
	MMSongZeldaAppear              string = "13D"
	MMSongSariasSong               string = "13E"
	MMSongGoronGoal                string = "13F"
	MMSongHorse                    string = "140"
	MMSongHorseGoal                string = "141"
	MMSongIngo                     string = "142"
	MMSongKotakePotionShop         string = "143"
	MMSongShop                     string = "144"
	MMSongOwl                      string = "145"
	MMSongShootingGallery          string = "146"
	MMSongOcarinaSoaring           string = "147"
	MMSongOcarinaHealing           string = "148"
	MMSongInvertedSongOfTime       string = "149"
	MMSongSongOfDoubleTime         string = "14A"
	MMSongSonataOfAwakening        string = "14B"
	MMSongGoronLullaby             string = "14C"
	MMSongNewWaveBossaNova         string = "14D"
	MMSongElegyOfEmptiness         string = "14E"
	MMSongOathToOrder              string = "14F"
	MMSongSwordTrainingHall        string = "150"
	MMSongOcarinaLullabyIntro      string = "151"
	MMSongLearnedNewSong           string = "152"
	MMSongBremenMarch              string = "153"
	MMSongBalladOfTheWindFish      string = "154"
	MMSongSongOfSoaring            string = "155"
	MMSongMilkBarDuplicate         string = "156"
	MMSongFinalHours               string = "157"
	MMSongMikauRiff                string = "158"
	MMSongMikauFinale              string = "159"
	MMSongFrogSong                 string = "15A"
	MMSongOcarinaSonata            string = "15B"
	MMSongOcarinaLullaby           string = "15C"
	MMSongOcarinaNewWave           string = "15D"
	MMSongOcarinaElegy             string = "15E"
	MMSongOcarinaOath              string = "15F"
	MMSongMajorasLair              string = "160"
	MMSongOcarinaLullabyIntroPtr   string = "161"
	MMSongOcarinaGuitarBassSession string = "162"
	MMSongPianoSession             string = "163"
	MMSongIndigoGoSession          string = "164"
	MMSongSnowheadTemple           string = "165"
	MMSongGreatBayTemple           string = "166"
	MMSongNewWaveSaxophone         string = "167"
	MMSongNewWaveVocal             string = "168"
	MMSongMajorasWrath             string = "169"
	MMSongMajorasIncarnation       string = "16A"
	MMSongMajorasMask              string = "16B"
	MMSongBassPlay                 string = "16C"
	MMSongDrumsPlay                string = "16D"
	MMSongPianoPlay                string = "16E"
	MMSongIkanaCastle              string = "16F"
	MMSongGatheringGiants          string = "170"
	MMSongKamaroDance              string = "171"
	MMSongCremiaCarriage           string = "172"
	MMSongKeatonQuiz               string = "173"
	MMSongEndCredits               string = "174"
	MMSongOpeningLoop              string = "175"
	MMSongTitleTheme               string = "176"
	MMSongDungeonAppear            string = "177"
	MMSongWoodfallClear            string = "178"
	MMSongSnowheadClear            string = "179"
	MMSongMusicBoxHouseInterior    string = "17A"
	MMSongIntoTheMoon              string = "17B"
	MMSongGoodbyeGiant             string = "17C"
	MMSongTatlAndTael              string = "17D"
	MMSongMoonsDestruction         string = "17E"
	MMSongEndCreditsSecondHalf     string = "17F"
)

var mmCatFanfare = []string{
	MMCatFanGetItem,
	MMCatFanGameOver,
	MMCatFanClear,
}

var ootCatConversion = map[string][]string{
	"Overworld":      {"Indoors", "Outdoors"},
	"Dungeon":        {"AdultDungeon", "ChildDungeon"},
	"Fight":          {"BigFight", "SmallFight"},
	"CharacterTheme": {"HeroTheme", "VillainTheme"},
	"EventFanfare":   {"ItemFanfare", "SuccessFanfare", "BigFanfare", "GameOver"},
	"SongFanfare":    {"WarpSong", "UtilitySong"},

	"Outdoors":       {"Fields", "Town", "Fun"},
	"Indoors":        {"Fun", "MagicalPlace", "House", "SalesArea", "WindmillHut"},
	"ChildDungeon":   {MMCatDungeon},
	"AdultDungeon":   {"AncientDungeon", "MysticalDungeon", "SpookyDungeon"},
	"SmallFight":     {MMCatAction},
	"BigFight":       {"BossFight", "FinalFight"},
	"HeroTheme":      {MMCatCalm},
	"VillainTheme":   {MMCatAction},
	"ItemFanfare":    {MMCatFanGetItem},
	"SuccessFanfare": {MMCatFanClear},
	"BigFanfare":     {MMCatFanClear, MMCatFanGetItem},
	"GameOver":       {MMCatFanGameOver},
	"WarpSong":       mmCatFanfare,
	"UtilitySong":    mmCatFanfare,

	"Fields":          {MMCatField},
	"Town":            {MMCatTown},
	"Fun":             {MMCatMinigame},
	"MagicalPlace":    {MMCatCalm},
	"House":           {MMCatIndoor},
	"SalesArea":       {MMCatIndoor},
	"AncientDungeon":  {MMCatDungeon},
	"MysticalDungeon": {MMCatDungeon},
	"SpookyDungeon":   {MMCatDungeon},
	"BossFight":       {MMCatBoss},
	"FinalFight":      {MMCatBoss},

	"HyruleField":       {MMCatField},
	"LostWoods":         {MMSongSariasSong},
	"GerudoValley":      {MMSongIkanaRegion},
	"Market":            {MMCatTown},
	"KakarikoChild":     {MMCatTown},
	"KakarikoAdult":     {MMCatTown},
	"LonLonRanch":       {MMSongRomaniRanch},
	"KokiriForest":      {MMSongSwampRegion},
	"GoronCity":         {MMSongGoronVillage},
	"ZorasDomain":       {MMSongZoraHall},
	"CastleCourtyard":   {MMCatTown},
	"HorseRace":         {MMSongHorse},
	"Mini-game":         {MMCatMinigame},
	"ShootingGallery":   {MMSongShootingGallery},
	"FairyFountain":     {MMSongFairyFountain},
	"TempleOfTime":      {MMCatCalm},
	"ChamberOfSages":    {MMCatCalm},
	"Shop":              {MMSongShop},
	"PotionShop":        {MMSongKotakePotionShop},
	"WindmillHut":       {MMSongMusicBoxHouse},
	"InsideDekuTree":    {MMSongWoodfallTemple},
	"DodongosCavern":    {MMCatDungeon},
	"JabuJabu":          {MMSongGreatBayTemple},
	"ForestTemple":      {MMSongWoodfallTemple},
	"FireTemple":        {MMCatDungeon},
	"IceCavern":         {MMCatDungeon},
	"WaterTemple":       {MMCatDungeon},
	"SpiritTemple":      {MMSongStoneTowerTemple},
	"ShadowTemple":      {MMSongInvertedStoneTowerTemple},
	"CastleUnderground": {MMSongCavern},
	"CastleEscape":      {MMCatAction},
	"Battle":            {MMCatAction},
	"MinibossBattle":    {MMCatBoss},
	"BossBattle":        {MMCatBoss},
	"FireBoss":          {MMCatBoss},
	"GanondorfBattle":   {MMCatBoss},
	"GanonBattle":       {MMCatBoss},
	"TitleTheme":        {MMCatTitle},
	"ZeldaTheme":        {MMSongZeldasLullaby},
	"SheikTheme":        {MMCatCalm},
	"DekuTree":          {MMCatDungeon},
	"KaeporaGaebora":    {MMSongOwl},
	"FairyFlying":       {MMCatCalm},
	"GanondorfTheme":    {MMCatAction},
	"KotakeAndKoume":    {MMCatBoss},
	"IngoTheme":         {MMSongIngo},

	"ItemGet":           {MMSongGetItem},
	"HeartContainerGet": {MMSongGetHeart},
	"SpiritStoneGet":    {MMSongGetNewMask},
	"HeartPieceGet":     {MMSongGetSmallItem},
	"MedallionGet":      {MMSongGetNewMask},
	"LearnSong":         {MMSongLearnedNewSong},
	"BossDefeated":      {MMSongClearBoss},
	"EponaRaceGoal":     {MMSongHorseGoal},
	"EscapeFromRanch":   {MMCatFanClear},
	"ZeldaTurnsAround":  {MMSongZeldaAppear},
	"TreasureChest":     {MMSongOpenChest},
	"MasterSword":       {MMCatFanGetItem},
	"DoorOfTime":        mmCatFanfare,
	"GanondorfAppears":  mmCatFanfare,

	"PreludeOfLight":   {"WarpSong"},
	"BoleroOfFire":     {"WarpSong"},
	"MinuetOfForest":   {"WarpSong"},
	"SerenadeOfWater":  {"WarpSong"},
	"RequiemOfSpirt":   {"WarpSong"},
	"NocturneOfShadow": {"WarpSong"},
	"SariasSong":       {"UtilitySong"},
	"EponasSong":       {"UtilitySong"},
	"ZeldasLullaby":    {"UtilitySong"},
	"SunsSong":         {"UtilitySong"},
	"SongOfTime":       {"UtilitySong"},
	"SongOfStorms":     {"UtilitySong"},
}
