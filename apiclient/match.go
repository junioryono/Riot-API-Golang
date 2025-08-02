package apiclient

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/junioryono/Riot-API-Golang/apiclient/ratelimiter"
	"github.com/junioryono/Riot-API-Golang/constants/continent"
	"github.com/junioryono/Riot-API-Golang/constants/patch"
	"github.com/junioryono/Riot-API-Golang/constants/queue"
	"github.com/junioryono/Riot-API-Golang/constants/region"
	"github.com/junioryono/Riot-API-Golang/constants/summoner_spell"
)

// Matchlist is an array of strings that represent the match IDs.
type Matchlist []string

func (m Matchlist) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

func (m *Matchlist) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

type GetMatchlistOptions struct {
	StartTime *time.Time `json:"startTime"`
	EndTime   *time.Time `json:"endTime"`
	Queue     *queue.ID  `json:"queue"`
	Type      *string    `json:"type"`
	Start     *int16     `json:"start"`
	Count     *int16     `json:"count"`
}

func (c *uniqueClient) GetMatchlist(continent continent.Continent, puuid string, opts *GetMatchlistOptions) (*Matchlist, error) {
	var params url.Values = make(map[string][]string)

	if opts != nil {
		if opts.StartTime != nil {
			params.Add("startTime", fmt.Sprintf("%d", opts.StartTime.UnixNano()/int64(time.Second)))
		}

		if opts.EndTime != nil {
			params.Add("endTime", fmt.Sprintf("%d", opts.EndTime.UnixNano()/int64(time.Second)))
		}

		if opts.Queue != nil {
			params.Add("queue", fmt.Sprintf("%d", *opts.Queue))
		}

		if opts.Type != nil {
			params.Add("type", *opts.Type)
		}

		if opts.Start != nil {
			params.Add("start", fmt.Sprintf("%d", *opts.Start))
		}

		if opts.Count != nil {
			params.Add("count", fmt.Sprintf("%d", *opts.Count))
		}
	}

	var res Matchlist
	err := c.dispatchAndUnmarshal(continent, "/lol/match/v5/matches/by-puuid", fmt.Sprintf("/%s/ids", puuid), params, ratelimiter.GetMatchlist, &res)
	return &res, err
}

type Match struct {
	Metadata MatchMetadata `json:"metadata"`
	Info     MatchInfo     `json:"info"`
}

type MatchMetadata struct {
	DataVersion  string   `json:"dataVersion"` // ex: 2
	MatchID      string   `json:"matchId"`     // ex: NA1_1234567890
	Participants []string `json:"participants"`
}

type MatchInfo struct {
	GameCreation       int64                  `json:"gameCreation"`       // ex: 1483140696030
	GameDuration       int32                  `json:"gameDuration"`       // ex: 1561
	GameEndTimestamp   int64                  `json:"gameEndTimestamp"`   // ex: 1483142251591
	GameID             int64                  `json:"gameId"`             // ex: 1234567890
	GameMode           string                 `json:"gameMode"`           // ex: CLASSIC
	GameName           string                 `json:"gameName"`           // ex: teambuilder-match-4742129795
	GameStartTimestamp int64                  `json:"gameStartTimestamp"` // ex: 1483140696030
	GameType           string                 `json:"gameType"`           // ex: MATCHED_GAME
	GameVersion        patch.Patch            `json:"gameVersion"`        // ex: 6.24.204.6436
	MapID              int16                  `json:"mapId"`              // ex: 11
	Participants       []MatchInfoParticipant `json:"participants"`
	PlatformID         region.Region          `json:"platformId"` // ex: NA1
	QueueID            queue.ID               `json:"queueId"`    // ex: 420
	Teams              []MatchInfoTeam        `json:"teams"`
	TournamentCode     string                 `json:"tournamentCode"`
}

type MatchInfoParticipant struct {
	AllInPings                     int16                           `json:"allInPings"`
	AssistMePings                  int16                           `json:"assistMePings"`
	Assists                        int16                           `json:"assists"`
	BaitPings                      int16                           `json:"baitPings"`
	BaronKills                     int16                           `json:"baronKills"`
	BasicPings                     int16                           `json:"basicPings"`
	BountyLevel                    int16                           `json:"bountyLevel"`
	Challenges                     *MatchInfoParticipantChallenges `json:"challenges"`
	ChampExperience                int32                           `json:"champExperience"`
	ChampLevel                     int16                           `json:"champLevel"`
	ChampionID                     int32                           `json:"championId"`
	ChampionName                   string                          `json:"championName"`
	ChampionTransform              int16                           `json:"championTransform"`
	CommandPings                   int16                           `json:"commandPings"`
	ConsumablesPurchased           int16                           `json:"consumablesPurchased"`
	DamageDealtToBuildings         int32                           `json:"damageDealtToBuildings"`
	DamageDealtToObjectives        int32                           `json:"damageDealtToObjectives"`
	DamageDealtToTurrets           int32                           `json:"damageDealtToTurrets"`
	DamageSelfMitigated            int32                           `json:"damageSelfMitigated"`
	DangerPings                    int16                           `json:"dangerPings"`
	Deaths                         int16                           `json:"deaths"`
	DetectorWardsPlaced            int16                           `json:"detectorWardsPlaced"`
	DoubleKills                    int16                           `json:"doubleKills"`
	DragonKills                    int16                           `json:"dragonKills"`
	EligibleForProgression         bool                            `json:"eligibleForProgression"`
	EnemyMissingPings              int16                           `json:"enemyMissingPings"`
	EnemyVisionPings               int16                           `json:"enemyVisionPings"`
	FirstBloodAssist               bool                            `json:"firstBloodAssist"`
	FirstBloodKill                 bool                            `json:"firstBloodKill"`
	FirstTowerAssist               bool                            `json:"firstTowerAssist"`
	FirstTowerKill                 bool                            `json:"firstTowerKill"`
	GameEndedInEarlySurrender      bool                            `json:"gameEndedInEarlySurrender"`
	GameEndedInSurrender           bool                            `json:"gameEndedInSurrender"`
	GetBackPings                   int16                           `json:"getBackPings"`
	GoldEarned                     int32                           `json:"goldEarned"`
	GoldSpent                      int32                           `json:"goldSpent"`
	HoldPings                      int16                           `json:"holdPings"`
	IndividualPosition             string                          `json:"individualPosition"`
	InhibitorKills                 int16                           `json:"inhibitorKills"`
	InhibitorTakedowns             int16                           `json:"inhibitorTakedowns"`
	InhibitorsLost                 int16                           `json:"inhibitorsLost"`
	Item0                          int32                           `json:"item0"`
	Item1                          int32                           `json:"item1"`
	Item2                          int32                           `json:"item2"`
	Item3                          int32                           `json:"item3"`
	Item4                          int32                           `json:"item4"`
	Item5                          int32                           `json:"item5"`
	Item6                          int32                           `json:"item6"`
	ItemsPurchased                 int16                           `json:"itemsPurchased"`
	KillingSprees                  int16                           `json:"killingSprees"`
	Kills                          int16                           `json:"kills"`
	Lane                           string                          `json:"lane"`
	LargestCriticalStrike          int32                           `json:"largestCriticalStrike"`
	LargestKillingSpree            int16                           `json:"largestKillingSpree"`
	LargestMultiKill               int16                           `json:"largestMultiKill"`
	LongestTimeSpentLiving         int32                           `json:"longestTimeSpentLiving"`
	MagicDamageDealt               int32                           `json:"magicDamageDealt"`
	MagicDamageDealtToChampions    int32                           `json:"magicDamageDealtToChampions"`
	MagicDamageTaken               int32                           `json:"magicDamageTaken"`
	NeedVisionPings                int16                           `json:"needVisionPings"`
	NeutralMinionsKilled           int16                           `json:"neutralMinionsKilled"`
	NexusKills                     int16                           `json:"nexusKills"`
	NexusLost                      int16                           `json:"nexusLost"`
	NexusTakedowns                 int16                           `json:"nexusTakedowns"`
	ObjectivesStolen               int16                           `json:"objectivesStolen"`
	ObjectivesStolenAssists        int16                           `json:"objectivesStolenAssists"`
	OnMyWayPings                   int16                           `json:"onMyWayPings"`
	ParticipantID                  int16                           `json:"participantId"`
	PentaKills                     int16                           `json:"pentaKills"`
	Perks                          MatchInfoParticipantPerks       `json:"perks"`
	PhysicalDamageDealt            int32                           `json:"physicalDamageDealt"`
	PhysicalDamageDealtToChampions int32                           `json:"physicalDamageDealtToChampions"`
	PhysicalDamageTaken            int32                           `json:"physicalDamageTaken"`
	Placement                      int16                           `json:"placement"`
	PlayerAugment1                 int32                           `json:"playerAugment1"`
	PlayerAugment2                 int32                           `json:"playerAugment2"`
	PlayerAugment3                 int32                           `json:"playerAugment3"`
	PlayerAugment4                 int32                           `json:"playerAugment4"`
	PlayerSubteamID                int16                           `json:"playerSubteamId"`
	PlayedChampSelectPosition      int16                           `json:"playedChampSelectPosition"`
	ProfileIcon                    int32                           `json:"profileIcon"`
	PushPings                      int16                           `json:"pushPings"`
	QuadraKills                    int16                           `json:"quadraKills"`
	RiotIdGameName                 string                          `json:"riotIdGameName"`
	RiotIdTagline                  string                          `json:"riotIdTagline"`
	Role                           string                          `json:"role"`
	SightWardsBoughtInGame         int16                           `json:"sightWardsBoughtInGame"`
	Spell1Casts                    int32                           `json:"spell1Casts"`
	Spell2Casts                    int32                           `json:"spell2Casts"`
	Spell3Casts                    int32                           `json:"spell3Casts"`
	Spell4Casts                    int32                           `json:"spell4Casts"`
	SubteamPlacement               int16                           `json:"subteamPlacement"`
	Summoner1Casts                 int16                           `json:"summoner1Casts"`
	Summoner1ID                    summoner_spell.ID               `json:"summoner1Id"`
	Summoner2Casts                 int16                           `json:"summoner2Casts"`
	Summoner2ID                    summoner_spell.ID               `json:"summoner2Id"`
	SummonerID                     string                          `json:"summonerId"`
	SummonerLevel                  int32                           `json:"summonerLevel"`
	SummonerName                   string                          `json:"summonerName"`
	SummonerPuuid                  string                          `json:"puuid"`
	TeamEarlySurrendered           bool                            `json:"teamEarlySurrendered"`
	TeamID                         int16                           `json:"teamId"`
	TeamPosition                   string                          `json:"teamPosition"`
	TimeCCingOthers                int16                           `json:"timeCCingOthers"`
	TimePlayed                     int16                           `json:"timePlayed"`
	TotalAllyJungleMinionsKilled   int16                           `json:"totalAllyJungleMinionsKilled"`
	TotalDamageDealt               int32                           `json:"totalDamageDealt"`
	TotalDamageDealtToChampions    int32                           `json:"totalDamageDealtToChampions"`
	TotalDamageShieldedOnTeammates int32                           `json:"totalDamageShieldedOnTeammates"`
	TotalDamageTaken               int32                           `json:"totalDamageTaken"`
	TotalEnemyJungleMinionsKilled  int16                           `json:"totalEnemyJungleMinionsKilled"`
	TotalHeal                      int32                           `json:"totalHeal"`
	TotalHealsOnTeammates          int32                           `json:"totalHealsOnTeammates"`
	TotalMinionsKilled             int16                           `json:"totalMinionsKilled"`
	TotalTimeCCDealt               int16                           `json:"totalTimeCCDealt"`
	TotalTimeSpentDead             int16                           `json:"totalTimeSpentDead"`
	TotalUnitsHealed               int16                           `json:"totalUnitsHealed"`
	TripleKills                    int16                           `json:"tripleKills"`
	TrueDamageDealt                int32                           `json:"trueDamageDealt"`
	TrueDamageDealtToChampions     int32                           `json:"trueDamageDealtToChampions"`
	TrueDamageTaken                int32                           `json:"trueDamageTaken"`
	TurretKills                    int16                           `json:"turretKills"`
	TurretTakedowns                int16                           `json:"turretTakedowns"`
	TurretsLost                    int16                           `json:"turretsLost"`
	UnrealKills                    int16                           `json:"unrealKills"`
	VisionClearedPings             int16                           `json:"visionClearedPings"`
	VisionScore                    int16                           `json:"visionScore"`
	VisionWardsBoughtInGame        int16                           `json:"visionWardsBoughtInGame"`
	WardsKilled                    int16                           `json:"wardsKilled"`
	WardsPlaced                    int16                           `json:"wardsPlaced"`
	Win                            bool                            `json:"win"`
}

type MatchInfoParticipantChallenges struct {
	Assist12StreakCount                      int16   `json:"12AssistStreakCount"`
	AbilityUses                              int32   `json:"abilityUses"`
	AcesBefore15Minutes                      int16   `json:"acesBefore15Minutes"`
	AlliedJungleMonsterKills                 int16   `json:"alliedJungleMonsterKills"`
	BaronTakedowns                           int16   `json:"baronTakedowns"`
	BlastConeOppositeOpponentCount           int16   `json:"blastConeOppositeOpponentCount"`
	BountyGold                               float64 `json:"bountyGold"`
	BuffsStolen                              int16   `json:"buffsStolen"`
	CompleteSupportQuestInTime               int16   `json:"completeSupportQuestInTime"`
	ControlWardsPlaced                       int16   `json:"controlWardsPlaced"`
	DamagePerMinute                          float64 `json:"damagePerMinute"`
	DamageTakenOnTeamPercentage              float64 `json:"damageTakenOnTeamPercentage"`
	DancedWithRiftHerald                     int16   `json:"dancedWithRiftHerald"`
	DeathsByEnemyChamps                      int16   `json:"deathsByEnemyChamps"`
	DodgeSkillShotsSmallWindow               int16   `json:"dodgeSkillShotsSmallWindow"`
	DoubleAces                               int16   `json:"doubleAces"`
	DragonTakedowns                          int16   `json:"dragonTakedowns"`
	EffectiveHealAndShielding                float64 `json:"effectiveHealAndShielding"`
	ElderDragonKillsWithOpposingSoul         int16   `json:"elderDragonKillsWithOpposingSoul"`
	ElderDragonMultikills                    int16   `json:"elderDragonMultikills"`
	EnemyChampionImmobilizations             int16   `json:"enemyChampionImmobilizations"`
	EnemyJungleMonsterKills                  int16   `json:"enemyJungleMonsterKills"`
	EpicMonsterKillsNearEnemyJungler         int16   `json:"epicMonsterKillsNearEnemyJungler"`
	EpicMonsterKillsWithin30SecondsOfSpawn   int16   `json:"epicMonsterKillsWithin30SecondsOfSpawn"`
	EpicMonsterSteals                        int16   `json:"epicMonsterSteals"`
	FastestLegendary                         float64 `json:"fastestLegendary"`
	EpicMonsterStolenWithoutSmite            int16   `json:"epicMonsterStolenWithoutSmite"`
	FirstTurretKilled                        float64 `json:"firstTurretKilled"`
	FirstTurretKilledTime                    float64 `json:"firstTurretKilledTime"`
	FistBumpParticipation                    int16   `json:"fistBumpParticipation"`
	FlawlessAces                             int16   `json:"flawlessAces"`
	FullTeamTakedown                         int16   `json:"fullTeamTakedown"`
	GameLength                               float64 `json:"gameLength"`
	GetTakedownsInAllLanesEarlyJungleAsLaner int16   `json:"getTakedownsInAllLanesEarlyJungleAsLaner"`
	GoldPerMinute                            float64 `json:"goldPerMinute"`
	HadOpenNexus                             int16   `json:"hadOpenNexus"`
	HealFromMapSources                       float64 `json:"HealFromMapSources"`
	ImmobilizeAndKillWithAlly                int16   `json:"immobilizeAndKillWithAlly"`
	InfernalScalePickup                      int16   `json:"InfernalScalePickup"`
	InitialBuffCount                         int16   `json:"initialBuffCount"`
	InitialCrabCount                         int16   `json:"initialCrabCount"`
	JungleCsBefore10Minutes                  float64 `json:"jungleCsBefore10Minutes"`
	JunglerTakedownsNearDamagedEpicMonster   int16   `json:"junglerTakedownsNearDamagedEpicMonster"`
	KTurretsDestroyedBeforePlatesFall        int16   `json:"kTurretsDestroyedBeforePlatesFall"`
	Kda                                      float64 `json:"kda"`
	KillAfterHiddenWithAlly                  int16   `json:"killAfterHiddenWithAlly"`
	KillParticipation                        float64 `json:"killParticipation"`
	KilledChampTookFullTeamDamageSurvived    int16   `json:"killedChampTookFullTeamDamageSurvived"`
	KillingSprees                            int16   `json:"killingSprees"`
	KillsNearEnemyTurret                     int16   `json:"killsNearEnemyTurret"`
	KillsOnOtherLanesEarlyJungleAsLaner      int16   `json:"killsOnOtherLanesEarlyJungleAsLaner"`
	KillsOnRecentlyHealedByAramPack          int16   `json:"killsOnRecentlyHealedByAramPack"`
	KillsUnderOwnTurret                      int16   `json:"killsUnderOwnTurret"`
	KillsWithHelpFromEpicMonster             int16   `json:"killsWithHelpFromEpicMonster"`
	KnockEnemyIntoTeamAndKill                int16   `json:"knockEnemyIntoTeamAndKill"`
	LandSkillShotsEarlyGame                  int32   `json:"landSkillShotsEarlyGame"`
	LaneMinionsFirst10Minutes                int16   `json:"laneMinionsFirst10Minutes"`
	LegendaryCount                           int16   `json:"legendaryCount"`
	LegendaryItemUsed                        []int32 `json:"legendaryItemUsed"`
	LostAnInhibitor                          int16   `json:"lostAnInhibitor"`
	MaxKillDeficit                           int16   `json:"maxKillDeficit"`
	MejaisFullStackInTime                    int16   `json:"mejaisFullStackInTime"`
	MoreEnemyJungleThanOpponent              float64 `json:"moreEnemyJungleThanOpponent"`
	MultiKillOneSpell                        int16   `json:"multiKillOneSpell"`
	MultiTurretRiftHeraldCount               int16   `json:"multiTurretRiftHeraldCount"`
	Multikills                               int16   `json:"multikills"`
	MultikillsAfterAggressiveFlash           int16   `json:"multikillsAfterAggressiveFlash"`
	OuterTurretExecutesBefore10Minutes       int16   `json:"outerTurretExecutesBefore10Minutes"`
	OutnumberedKills                         int16   `json:"outnumberedKills"`
	OutnumberedNexusKill                     int16   `json:"outnumberedNexusKill"`
	PerfectDragonSoulsTaken                  int16   `json:"perfectDragonSoulsTaken"`
	PerfectGame                              int16   `json:"perfectGame"`
	PickKillWithAlly                         int16   `json:"pickKillWithAlly"`
	PoroExplosions                           int16   `json:"poroExplosions"`
	QuickCleanse                             int16   `json:"quickCleanse"`
	QuickFirstTurret                         int16   `json:"quickFirstTurret"`
	QuickSoloKills                           int16   `json:"quickSoloKills"`
	RiftHeraldTakedowns                      int16   `json:"riftHeraldTakedowns"`
	SaveAllyFromDeath                        int16   `json:"saveAllyFromDeath"`
	ScuttleCrabKills                         int16   `json:"scuttleCrabKills"`
	ShortestTimeToAceFromFirstTakedown       float64 `json:"shortestTimeToAceFromFirstTakedown"`
	SkillshotsDodged                         int32   `json:"skillshotsDodged"`
	SkillshotsHit                            int32   `json:"skillshotsHit"`
	SnowballsHit                             int16   `json:"snowballsHit"`
	SoloBaronKills                           int16   `json:"soloBaronKills"`
	SoloKills                                int16   `json:"soloKills"`
	StealthWardsPlaced                       int16   `json:"stealthWardsPlaced"`
	SurvivedSingleDigitHpCount               int16   `json:"survivedSingleDigitHpCount"`
	SurvivedThreeImmobilizesInFight          int16   `json:"survivedThreeImmobilizesInFight"`
	SWARM_DefeatAatrox                       int16   `json:"SWARM_DefeatAatrox"`
	SWARM_DefeatBriar                        int16   `json:"SWARM_DefeatBriar"`
	SWARM_DefeatMiniBosses                   int16   `json:"SWARM_DefeatMiniBosses"`
	SWARM_EvolveWeapon                       int16   `json:"SWARM_EvolveWeapon"`
	SWARM_Have3Passives                      int16   `json:"SWARM_Have3Passives"`
	SWARM_KillEnemy                          int16   `json:"SWARM_KillEnemy"`
	SWARM_PickupGold                         int16   `json:"SWARM_PickupGold"`
	SWARM_ReachLevel50                       int16   `json:"SWARM_ReachLevel50"`
	SWARM_Survive15Min                       int16   `json:"SWARM_Survive15Min"`
	SWARM_WinWith5EvolvedWeapons             int16   `json:"SWARM_WinWith5EvolvedWeapons"`
	TakedownOnFirstTurret                    int16   `json:"takedownOnFirstTurret"`
	Takedowns                                int16   `json:"takedowns"`
	TakedownsAfterGainingLevelAdvantage      int16   `json:"takedownsAfterGainingLevelAdvantage"`
	TakedownsBeforeJungleMinionSpawn         int16   `json:"takedownsBeforeJungleMinionSpawn"`
	TakedownsFirstXMinutes                   int16   `json:"takedownsFirstXMinutes"`
	TakedownsInAlcove                        int16   `json:"takedownsInAlcove"`
	TakedownsInEnemyFountain                 int16   `json:"takedownsInEnemyFountain"`
	TeamBaronKills                           int16   `json:"teamBaronKills"`
	TeamDamagePercentage                     float64 `json:"teamDamagePercentage"`
	TeamElderDragonKills                     int16   `json:"teamElderDragonKills"`
	TeamRiftHeraldKills                      int16   `json:"teamRiftHeraldKills"`
	TookLargeDamageSurvived                  int16   `json:"tookLargeDamageSurvived"`
	TurretPlatesTaken                        int16   `json:"turretPlatesTaken"`
	TurretTakedowns                          int16   `json:"turretTakedowns"`
	TurretsTakenWithRiftHerald               int16   `json:"turretsTakenWithRiftHerald"`
	TwentyMinionsIn3SecondsCount             int16   `json:"twentyMinionsIn3SecondsCount"`
	TwoWardsOneSweeperCount                  int16   `json:"twoWardsOneSweeperCount"`
	UnseenRecalls                            int16   `json:"unseenRecalls"`
	VisionScorePerMinute                     float32 `json:"visionScorePerMinute"`
	VoidMonsterKill                          int16   `json:"voidMonsterKill"`
	WardTakedowns                            int16   `json:"wardTakedowns"`
	WardTakedownsBefore20M                   int16   `json:"wardTakedownsBefore20M"`
	WardsGuarded                             int16   `json:"wardsGuarded"`
}

type MatchInfoParticipantPerks struct {
	StatPerks MatchInfoParticipantPerksStatPerks `json:"statPerks"`
	Styles    []MatchInfoParticipantPerksStyles  `json:"styles"`
}

type MatchInfoParticipantPerksStatPerks struct {
	Defense int16 `json:"defense"`
	Flex    int16 `json:"flex"`
	Offense int16 `json:"offense"`
}

type MatchInfoParticipantPerksStyles struct {
	Description string                                     `json:"description"`
	Selections  []MatchInfoParticipantPerksStylesSelection `json:"selections"`
	Style       int16                                      `json:"style"`
}

type MatchInfoParticipantPerksStylesSelection struct {
	Perk int16 `json:"perk"`
	Var1 int32 `json:"var1"`
	Var2 int32 `json:"var2"`
	Var3 int32 `json:"var3"`
}

type MatchInfoTeam struct {
	Bans       []MatchInfoTeamBan      `json:"bans"`
	Objectives MatchInfoTeamObjectives `json:"objectives"`
	TeamID     int16                   `json:"teamId"`
	Win        bool                    `json:"win"`
}

type MatchInfoTeamBan struct {
	PickTurn   int16 `json:"pickTurn"`
	ChampionID int16 `json:"championId"`
}

type MatchInfoTeamObjectives struct {
	Baron      MatchInfoTeamObjectiveType `json:"baron"`
	Champion   MatchInfoTeamObjectiveType `json:"champion"`
	Dragon     MatchInfoTeamObjectiveType `json:"dragon"`
	Horde      MatchInfoTeamObjectiveType `json:"horde"`
	Inhibitor  MatchInfoTeamObjectiveType `json:"inhibitor"`
	RiftHerald MatchInfoTeamObjectiveType `json:"riftHerald"`
	Tower      MatchInfoTeamObjectiveType `json:"tower"`
}

type MatchInfoTeamObjectiveType struct {
	First bool  `json:"first"`
	Kills int16 `json:"kills"`
}

func (m MatchInfo) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

func (m *MatchInfo) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

func (c *uniqueClient) GetMatch(continent continent.Continent, matchID string) (*Match, error) {
	var res Match
	err := c.dispatchAndUnmarshal(continent, "/lol/match/v5/matches", fmt.Sprintf("/%s", matchID), nil, ratelimiter.GetMatch, &res)
	return &res, err
}

type MatchTimeline struct {
	Metadata MatchTimelineMetadata `json:"metadata"`
	Info     MatchTimelineInfo     `json:"info"`
}

type MatchTimelineMetadata struct {
	DataVersion  string   `json:"dataVersion"`
	MatchID      string   `json:"matchId"`
	Participants []string `json:"participants"`
}

type MatchTimelineInfo struct {
	FrameInterval int32                      `json:"frameInterval"`
	Frames        []MatchTimelineFrame       `json:"frames"`
	GameID        int64                      `json:"gameId"`
	Participants  []MatchTimelineParticipant `json:"participants"`
}

type MatchTimelineParticipant struct {
	ParticipantID int16  `json:"participantId"`
	Puuid         string `json:"puuid"`
}

type MatchTimelineFrame struct {
	Timestamp         int32                           `json:"timestamp"`
	ParticipantFrames []MatchTimelineParticipantFrame `json:"participantFrames"`
	Events            []interface{}                   `json:"events"`
}

type MatchTimelineFrameEventType string

const (
	AscendedEvent           MatchTimelineFrameEventType = "ASCENDED_EVENT"
	BuildingKill            MatchTimelineFrameEventType = "BUILDING_KILL"
	CapturePoint            MatchTimelineFrameEventType = "CAPTURE_POINT"
	ChampionKill            MatchTimelineFrameEventType = "CHAMPION_KILL"
	ChampionSpecialKill     MatchTimelineFrameEventType = "CHAMPION_SPECIAL_KILL"
	ChampionTransform       MatchTimelineFrameEventType = "CHAMPION_TRANSFORM"
	DragonSoulGiven         MatchTimelineFrameEventType = "DRAGON_SOUL_GIVEN"
	EliteMonsterKill        MatchTimelineFrameEventType = "ELITE_MONSTER_KILL"
	GameEnd                 MatchTimelineFrameEventType = "GAME_END"
	ItemDestroyed           MatchTimelineFrameEventType = "ITEM_DESTROYED"
	ItemPurchased           MatchTimelineFrameEventType = "ITEM_PURCHASED"
	ItemSold                MatchTimelineFrameEventType = "ITEM_SOLD"
	ItemUndo                MatchTimelineFrameEventType = "ITEM_UNDO"
	LevelUp                 MatchTimelineFrameEventType = "LEVEL_UP"
	ObjectiveBountyFinish   MatchTimelineFrameEventType = "OBJECTIVE_BOUNTY_FINISH"
	ObjectiveBountyPreStart MatchTimelineFrameEventType = "OBJECTIVE_BOUNTY_PRESTART"
	PauseEnd                MatchTimelineFrameEventType = "PAUSE_END"
	PoroKingSummon          MatchTimelineFrameEventType = "PORO_KING_SUMMON"
	SkillLevelUp            MatchTimelineFrameEventType = "SKILL_LEVEL_UP"
	TurretPlateDestroyed    MatchTimelineFrameEventType = "TURRET_PLATE_DESTROYED"
	WardKill                MatchTimelineFrameEventType = "WARD_KILL"
	WardPlaced              MatchTimelineFrameEventType = "WARD_PLACED"
)

type MatchTimelineEvent_AscendedEvent struct {
	Timestamp int32                       `json:"timestamp"`
	Type      MatchTimelineFrameEventType `json:"type"`
	// TODO
}

type MatchTimelineEvent_BuildingKill struct {
	AssistingParticipantIDs []int16                         `json:"assistingParticipantIds"`
	Bounty                  int16                           `json:"bounty"`
	BuildingType            MatchTimelineEvent_BuildingType `json:"buildingType"`
	KillerID                int16                           `json:"killerId"`
	LaneType                MatchTimelineEvent_LaneType     `json:"laneType"`
	Position                MatchTimelinePosition           `json:"position"`
	TeamID                  int16                           `json:"teamId"`
	Timestamp               int32                           `json:"timestamp"`
	TowerType               *MatchTimelineEvent_TowerType   `json:"towerType"`
	Type                    MatchTimelineFrameEventType     `json:"type"`
}

type MatchTimelineEvent_CapturePoint struct {
	Timestamp int32                       `json:"timestamp"`
	Type      MatchTimelineFrameEventType `json:"type"`
	// TODO
}

type MatchTimelineEvent_ChampionKill struct {
	AssistingParticipantIDs []int16                     `json:"assistingParticipantIds"`
	Bounty                  int16                       `json:"bounty"`
	KillStreakLength        int16                       `json:"killStreakLength"`
	KillerID                int16                       `json:"killerId"`
	Position                MatchTimelinePosition       `json:"position"`
	ShutdownBounty          int16                       `json:"shutdownBounty"`
	Timestamp               int32                       `json:"timestamp"`
	Type                    MatchTimelineFrameEventType `json:"type"`
	VictimDamageDealt       []MatchTimelineDamage       `json:"victimDamageDealt"`
	VictimDamageReceived    []MatchTimelineDamage       `json:"victimDamageReceived"`
	VictimID                int16                       `json:"victimId"`
}

type MatchTimelineEvent_ChampionSpecialKill struct {
	KillType        MatchTimelineEvent_KillType `json:"killType"`
	KillerID        int16                       `json:"killerId"`
	MultiKillLength *int16                      `json:"multiKillLength"`
	Position        MatchTimelinePosition       `json:"position"`
	Timestamp       int32                       `json:"timestamp"`
	Type            MatchTimelineFrameEventType `json:"type"`
}

type MatchTimelineEvent_ChampionTransform struct {
	ParticipantID int16                                    `json:"participantId"`
	Timestamp     int32                                    `json:"timestamp"`
	TransformType MatchTimelineEvent_ChampionTransformType `json:"transformType"`
	Type          MatchTimelineFrameEventType              `json:"type"`
}

type MatchTimelineEvent_ChampionTransformType string

const (
	MatchTimelineEvent_ChampionTransformType_Assassin MatchTimelineEvent_ChampionTransformType = "ASSASSIN"
	MatchTimelineEvent_ChampionTransformType_Slayer   MatchTimelineEvent_ChampionTransformType = "SLAYER"
)

type MatchTimelineEvent_DragonSoulGiven struct {
	Name      MatchTimelineEvent_DragonSoul `json:"name"`
	TeamID    int16                         `json:"teamId"`
	Timestamp int32                         `json:"timestamp"`
	Type      MatchTimelineFrameEventType   `json:"type"`
}

type MatchTimelineEvent_EliteMonsterKill struct {
	AssistingParticipantIDs []int16                           `json:"assistingParticipantIds"`
	Bounty                  int16                             `json:"bounty"`
	KillerID                int16                             `json:"killerId"`
	KillerTeamID            int16                             `json:"killerTeamId"`
	MonsterSubType          MatchTimelineEvent_MonsterSubType `json:"monsterSubType"`
	MonsterType             MatchTimelineEvent_MonsterType    `json:"monsterType"`
	Position                MatchTimelinePosition             `json:"position"`
	Timestamp               int32                             `json:"timestamp"`
	Type                    MatchTimelineFrameEventType       `json:"type"`
}

type MatchTimelineEvent_GameEnd struct {
	GameID        int64                       `json:"gameId"`
	RealTimestamp int64                       `json:"realTimestamp"`
	Timestamp     int32                       `json:"timestamp"`
	Type          MatchTimelineFrameEventType `json:"type"`
	WinningTeam   int16                       `json:"winningTeam"`
}

type MatchTimelineEvent_ItemDestroyed struct {
	ItemID        int32                       `json:"itemId"`
	ParticipantID int16                       `json:"participantId"`
	Timestamp     int32                       `json:"timestamp"`
	Type          MatchTimelineFrameEventType `json:"type"`
}

type MatchTimelineEvent_ItemPurchased struct {
	ItemID        int32                       `json:"itemId"`
	ParticipantID int16                       `json:"participantId"`
	Timestamp     int32                       `json:"timestamp"`
	Type          MatchTimelineFrameEventType `json:"type"`
}

type MatchTimelineEvent_ItemSold struct {
	ItemID        int32                       `json:"itemId"`
	ParticipantID int16                       `json:"participantId"`
	Timestamp     int32                       `json:"timestamp"`
	Type          MatchTimelineFrameEventType `json:"type"`
}

type MatchTimelineEvent_ItemUndo struct {
	AfterID       int32                       `json:"afterId"`
	BeforeID      int32                       `json:"beforeId"`
	GoldGain      int16                       `json:"goldGain"`
	ParticipantID int16                       `json:"participantId"`
	Timestamp     int32                       `json:"timestamp"`
	Type          MatchTimelineFrameEventType `json:"type"`
}

type MatchTimelineEvent_LevelUp struct {
	Level         int16                       `json:"level"`
	ParticipantID int16                       `json:"participantId"`
	Timestamp     int32                       `json:"timestamp"`
	Type          MatchTimelineFrameEventType `json:"type"`
}

type MatchTimelineEvent_ObjectiveBountyFinish struct {
	TeamID    int16                       `json:"teamId"`
	Timestamp int32                       `json:"timestamp"`
	Type      MatchTimelineFrameEventType `json:"type"`
}

type MatchTimelineEvent_ObjectiveBountyPreStart struct {
	ActualStartTime int32                       `json:"actualStartTime"`
	TeamID          int16                       `json:"teamId"`
	Timestamp       int32                       `json:"timestamp"`
	Type            MatchTimelineFrameEventType `json:"type"`
}

type MatchTimelineEvent_PauseEnd struct {
	RealTimestamp int64                       `json:"realTimestamp"`
	Timestamp     int32                       `json:"timestamp"`
	Type          MatchTimelineFrameEventType `json:"type"`
}

type MatchTimelineEvent_PoroKingSummon struct {
	Timestamp int32 `json:"timestamp"`
	Type      MatchTimelineFrameEventType
	// TODO
}

type MatchTimelineEvent_SkillLevelUp struct {
	LevelUpType   MatchTimelineEvent_LevelUpType `json:"levelUpType"`
	ParticipantID int16                          `json:"participantId"`
	SkillSlot     int16                          `json:"skillSlot"`
	Timestamp     int32                          `json:"timestamp"`
	Type          MatchTimelineFrameEventType    `json:"type"`
}

type MatchTimelineEvent_TurretPlateDestroyed struct {
	KillerID  int16                       `json:"killerId"`
	LaneType  MatchTimelineEvent_LaneType `json:"laneType"`
	Position  MatchTimelinePosition       `json:"position"`
	TeamID    int16                       `json:"teamId"`
	Timestamp int32                       `json:"timestamp"`
	Type      MatchTimelineFrameEventType `json:"type"`
}

type MatchTimelineEvent_WardKill struct {
	KillerID  int16                       `json:"killerId"`
	Timestamp int32                       `json:"timestamp"`
	Type      MatchTimelineFrameEventType `json:"type"`
	WardType  MatchTimelineEvent_WardType `json:"wardType"`
}

type MatchTimelineEvent_WardPlaced struct {
	CreatorID int16                       `json:"creatorId"`
	Timestamp int32                       `json:"timestamp"`
	Type      MatchTimelineFrameEventType `json:"type"`
	WardType  MatchTimelineEvent_WardType `json:"wardType"`
}

type MatchTimelineEvent_BuildingType string

const (
	MatchTimelineEvent_BuildingType_InhibitorBuilding MatchTimelineEvent_BuildingType = "INHIBITOR_BUILDING"
	MatchTimelineEvent_BuildingType_TowerBuilding     MatchTimelineEvent_BuildingType = "TOWER_BUILDING"
)

type MatchTimelineEvent_DragonSoul string

const (
	MatchTimelineEvent_DragonSoul_Chemtech MatchTimelineEvent_DragonSoul = "Chemtech"
	MatchTimelineEvent_DragonSoul_Cloud    MatchTimelineEvent_DragonSoul = "Cloud"
	MatchTimelineEvent_DragonSoul_Hextech  MatchTimelineEvent_DragonSoul = "Hextech"
	MatchTimelineEvent_DragonSoul_Infernal MatchTimelineEvent_DragonSoul = "Infernal"
	MatchTimelineEvent_DragonSoul_Mountain MatchTimelineEvent_DragonSoul = "Mountain"
	MatchTimelineEvent_DragonSoul_Ocean    MatchTimelineEvent_DragonSoul = "Ocean"
)

type MatchTimelineEvent_KillType string

const (
	MatchTimelineEvent_KillType_Ace             MatchTimelineEvent_KillType = "ACE"
	MatchTimelineEvent_KillType_KillMulti       MatchTimelineEvent_KillType = "KILL_MULTI"
	MatchTimelineEvent_KillType_KillFirst_Blood MatchTimelineEvent_KillType = "KILL_FIRST_BLOOD"
)

type MatchTimelineEvent_LaneType string

const (
	MatchTimelineEvent_LaneType_BotLane MatchTimelineEvent_LaneType = "BOT_LANE"
	MatchTimelineEvent_LaneType_MidLane MatchTimelineEvent_LaneType = "MID_LANE"
	MatchTimelineEvent_LaneType_TopLane MatchTimelineEvent_LaneType = "TOP_LANE"
)

type MatchTimelineEvent_LevelUpType string

const (
	MatchTimelineEvent_LevelUpType_Evolve MatchTimelineEvent_LevelUpType = "EVOLVE"
	MatchTimelineEvent_LevelUpType_Normal MatchTimelineEvent_LevelUpType = "NORMAL"
)

type MatchTimelineEvent_MonsterSubType string

const (
	MatchTimelineEvent_MonsterSubType_AirDragon      MatchTimelineEvent_MonsterSubType = "AIR_DRAGON"
	MatchTimelineEvent_MonsterSubType_ChemtechDragon MatchTimelineEvent_MonsterSubType = "CHEMTECH_DRAGON"
	MatchTimelineEvent_MonsterSubType_EarthDragon    MatchTimelineEvent_MonsterSubType = "EARTH_DRAGON"
	MatchTimelineEvent_MonsterSubType_ElderDragon    MatchTimelineEvent_MonsterSubType = "ELDER_DRAGON"
	MatchTimelineEvent_MonsterSubType_FireDragon     MatchTimelineEvent_MonsterSubType = "FIRE_DRAGON"
	MatchTimelineEvent_MonsterSubType_WaterDragon    MatchTimelineEvent_MonsterSubType = "WATER_DRAGON"
	MatchTimelineEvent_MonsterSubType_HextechDragon  MatchTimelineEvent_MonsterSubType = "HEXTECH_DRAGON"
)

type MatchTimelineEvent_MonsterType string

const (
	MatchTimelineEvent_MonsterType_Baron      MatchTimelineEvent_MonsterType = "BARON_NASHOR"
	MatchTimelineEvent_MonsterType_Dragon     MatchTimelineEvent_MonsterType = "DRAGON"
	MatchTimelineEvent_MonsterType_RiftHerald MatchTimelineEvent_MonsterType = "RIFTHERALD"
)

type MatchTimelineEvent_TowerType string

const (
	MatchTimelineEvent_TowerType_BaseTurret  MatchTimelineEvent_TowerType = "BASE_TURRET"
	MatchTimelineEvent_TowerType_InnerTurret MatchTimelineEvent_TowerType = "INNER_TURRET"
	MatchTimelineEvent_TowerType_NexusTurret MatchTimelineEvent_TowerType = "NEXUS_TURRET"
	MatchTimelineEvent_TowerType_OuterTurret MatchTimelineEvent_TowerType = "OUTER_TURRET"
)

type MatchTimelineEvent_WardType string

const (
	MatchTimelineEvent_WardType_ControlWard MatchTimelineEvent_WardType = "CONTROL_WARD"
	MatchTimelineEvent_WardType_SightWard   MatchTimelineEvent_WardType = "SIGHT_WARD"
	MatchTimelineEvent_WardType_TeemoMush   MatchTimelineEvent_WardType = "TEEMO_MUSHROOM"
	MatchTimelineEvent_WardType_YellowTrink MatchTimelineEvent_WardType = "YELLOW_TRINKET"
)

// Need to unmarshal MatchTimelineFrame because ParticipantFrames comes as an object with keys as participant IDs.
func (m *MatchTimelineFrame) UnmarshalJSON(data []byte) error {
	// Define a temporary struct with the same fields as MatchTimelineFrame
	temp := &struct {
		Timestamp         int32                                    `json:"timestamp"`
		ParticipantFrames map[string]MatchTimelineParticipantFrame `json:"participantFrames"`
		Events            []json.RawMessage                        `json:"events"`
	}{}

	// Unmarshal data into the temporary struct
	if err := json.Unmarshal(data, temp); err != nil {
		return err
	}

	// Manually assign the values to the MatchTimelineFrame
	m.Timestamp = temp.Timestamp

	for _, rawMsg := range temp.Events {
		// Unmarshal the type field to determine the event type
		var typeHolder struct {
			Type MatchTimelineFrameEventType `json:"type"`
		}
		if err := json.Unmarshal(rawMsg, &typeHolder); err != nil {
			return err
		}

		var event interface{}
		switch typeHolder.Type {
		case AscendedEvent:
			event = new(MatchTimelineEvent_AscendedEvent)
		case BuildingKill:
			event = new(MatchTimelineEvent_BuildingKill)
		case CapturePoint:
			event = new(MatchTimelineEvent_CapturePoint)
		case ChampionKill:
			event = new(MatchTimelineEvent_ChampionKill)
		case ChampionSpecialKill:
			event = new(MatchTimelineEvent_ChampionSpecialKill)
		case ChampionTransform:
			event = new(MatchTimelineEvent_ChampionTransform)
		case DragonSoulGiven:
			event = new(MatchTimelineEvent_DragonSoulGiven)
		case EliteMonsterKill:
			event = new(MatchTimelineEvent_EliteMonsterKill)
		case GameEnd:
			event = new(MatchTimelineEvent_GameEnd)
		case ItemDestroyed:
			event = new(MatchTimelineEvent_ItemDestroyed)
		case ItemPurchased:
			event = new(MatchTimelineEvent_ItemPurchased)
		case ItemSold:
			event = new(MatchTimelineEvent_ItemSold)
		case ItemUndo:
			event = new(MatchTimelineEvent_ItemUndo)
		case LevelUp:
			event = new(MatchTimelineEvent_LevelUp)
		case ObjectiveBountyFinish:
			event = new(MatchTimelineEvent_ObjectiveBountyFinish)
		case ObjectiveBountyPreStart:
			event = new(MatchTimelineEvent_ObjectiveBountyPreStart)
		case PauseEnd:
			event = new(MatchTimelineEvent_PauseEnd)
		case PoroKingSummon:
			event = new(MatchTimelineEvent_PoroKingSummon)
		case SkillLevelUp:
			event = new(MatchTimelineEvent_SkillLevelUp)
		case TurretPlateDestroyed:
			event = new(MatchTimelineEvent_TurretPlateDestroyed)
		case WardKill:
			event = new(MatchTimelineEvent_WardKill)
		case WardPlaced:
			event = new(MatchTimelineEvent_WardPlaced)
		default:
			continue // Skip unknown event types
		}

		// Unmarshal the event
		if err := json.Unmarshal(rawMsg, event); err != nil {
			return err
		}

		m.Events = append(m.Events, event)
	}

	// Convert map to slice
	m.ParticipantFrames = make([]MatchTimelineParticipantFrame, 0, len(temp.ParticipantFrames))
	for _, pf := range temp.ParticipantFrames {
		m.ParticipantFrames = append(m.ParticipantFrames, pf)
	}

	return nil
}

type MatchTimelineParticipantFrame struct {
	ChampionStats            MatchTimelineChampionStats `json:"championStats"`
	CurrentGold              int32                      `json:"currentGold"`
	DamageStats              MatchTimelineDamageStats   `json:"damageStats"`
	JungleMinionsKilled      int16                      `json:"jungleMinionsKilled"`
	Level                    int16                      `json:"level"`
	MinionsKilled            int16                      `json:"minionsKilled"`
	ParticipantID            int16                      `json:"participantId"`
	Position                 MatchTimelinePosition      `json:"position"`
	TimeEnemySpentControlled int32                      `json:"timeEnemySpentControlled"`
	TotalGold                int32                      `json:"totalGold"`
	XP                       int32                      `json:"xp"`
}

type MatchTimelineChampionStats struct {
	AbilityHaste         int16 `json:"abilityHaste"`
	AbilityPower         int16 `json:"abilityPower"`
	Armor                int16 `json:"armor"`
	ArmorPen             int16 `json:"armorPen"`
	ArmorPenPercent      int16 `json:"armorPenPercent"`
	AttackDamage         int16 `json:"attackDamage"`
	AttackSpeed          int16 `json:"attackSpeed"`
	BonusArmorPenPercent int16 `json:"bonusArmorPenPercent"`
	BonusMagicPenPercent int16 `json:"bonusMagicPenPercent"`
	CCReduction          int16 `json:"ccReduction"`
	CooldownReduction    int16 `json:"cooldownReduction"`
	Health               int32 `json:"health"`
	HealthMax            int32 `json:"healthMax"`
	HealthRegen          int32 `json:"healthRegen"`
	Lifesteal            int16 `json:"lifesteal"`
	MagicPen             int16 `json:"magicPen"`
	MagicPenPercent      int16 `json:"magicPenPercent"`
	MagicResist          int16 `json:"magicResist"`
	MovementSpeed        int16 `json:"movementSpeed"`
	Omnivamp             int16 `json:"omnivamp"`
	PhysicalVamp         int16 `json:"physicalVamp"`
	Power                int16 `json:"power"`
	PowerMax             int16 `json:"powerMax"`
	PowerRegen           int16 `json:"powerRegen"`
	SpellVamp            int16 `json:"spellVamp"`
}

type MatchTimelineDamageStats struct {
	MagicDamageDone               int32 `json:"magicDamageDone"`
	MagicDamageDoneToChampions    int32 `json:"magicDamageDoneToChampions"`
	MagicDamageTaken              int32 `json:"magicDamageTaken"`
	PhysicalDamageDone            int32 `json:"physicalDamageDone"`
	PhysicalDamageDoneToChampions int32 `json:"physicalDamageDoneToChampions"`
	PhysicalDamageTaken           int32 `json:"physicalDamageTaken"`
	TotalDamageDone               int32 `json:"totalDamageDone"`
	TotalDamageDoneToChampions    int32 `json:"totalDamageDoneToChampions"`
	TotalDamageTaken              int32 `json:"totalDamageTaken"`
	TrueDamageDone                int32 `json:"trueDamageDone"`
	TrueDamageDoneToChampions     int32 `json:"trueDamageDoneToChampions"`
	TrueDamageTaken               int32 `json:"trueDamageTaken"`
}

type MatchTimelinePosition struct {
	X int16 `json:"x"`
	Y int16 `json:"y"`
}

type MatchTimelineDamage struct {
	Basic          bool                    `json:"basic"`
	MagicDamage    int32                   `json:"magicDamage"`
	Name           string                  `json:"name"`
	ParticipantID  int16                   `json:"participantId"`
	PhysicalDamage int32                   `json:"physicalDamage"`
	SpellName      string                  `json:"spellName"`
	SpellSlot      int32                   `json:"spellSlot"`
	TrueDamage     int32                   `json:"trueDamage"`
	Type           MatchTimelineDamageType `json:"type"`
}

type MatchTimelineDamageType string

const (
	MatchTimelineDamageType_Minion  MatchTimelineDamageType = "MINION"
	MatchTimelineDamageType_Monster MatchTimelineDamageType = "MONSTER"
	MatchTimelineDamageType_Tower   MatchTimelineDamageType = "TOWER"
	MatchTimelineDamageType_Other   MatchTimelineDamageType = "OTHER"
)

func (c *uniqueClient) GetMatchTimeline(continent continent.Continent, matchID string) (*MatchTimeline, error) {
	var res MatchTimeline
	err := c.dispatchAndUnmarshal(continent, "/lol/match/v5/matches", fmt.Sprintf("/%s/timeline", matchID), nil, ratelimiter.GetMatchTimeline, &res)
	return &res, err
}
