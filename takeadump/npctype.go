package takeadump

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/xackery/importeq/model"
)

// newNpcType returns a new npc type based on provided data
func newNpcType(parser *parser, line int, zoneShortName string, zoneID int64, zoneInstanceVersion int64) (npc *model.NpcType, err error) {
	npc = &model.NpcType{
		ZoneShortName:       zoneShortName,
		ZoneID:              zoneID,
		ZoneInstanceVersion: zoneInstanceVersion,
		Sources:             []string{parser.filename},
	}
	npc.Name, err = parser.EntryByName(line, columnNPCName)
	if err != nil {
		err = errors.Wrap(err, "name")
		return
	}
	npc.Name = model.StringCleanName(npc.Name)

	data, err := parser.EntryByName(line, columnNPCLevel)
	if err != nil {
		err = errors.Wrap(err, "level")
		return
	}
	npc.Level, err = strconv.ParseInt(data, 10, 64)
	if err != nil {
		err = errors.Wrap(err, "level conversion")
		return
	}
	npc.MaxLevel = npc.Level

	data, err = parser.EntryByName(line, columnNPCClass)
	if err != nil {
		err = errors.Wrap(err, "class")
		return
	}
	npc.Class, err = strconv.ParseInt(data, 10, 64)

	npc.LastName, err = parser.EntryByName(line, columnNPCLastName)
	if err != nil {
		err = errors.Wrap(err, "lastName")
		return
	}

	data, err = parser.EntryByName(line, columnNPCBodyType)
	if err != nil {
		err = errors.Wrap(err, "bodytype")
		return
	}
	npc.BodyType, err = strconv.ParseInt(data, 10, 64)
	if err != nil {
		err = errors.Wrap(err, "bodytype conversion")
		return
	}

	data, err = parser.EntryByName(line, columnNPCDrakkinDetails)
	if err != nil {
		err = errors.Wrap(err, "DrakkinDetails")
		return
	}
	npc.DrakkinDetails, err = strconv.ParseInt(data, 10, 64)
	if err != nil {
		err = errors.Wrap(err, "DrakkinDetails conversion")
		return
	}
	data, err = parser.EntryByName(line, columnNPCDrakkinHeritage)
	if err != nil {
		err = errors.Wrap(err, "DrakkinHeritage")
		return
	}
	npc.DrakkinHeritage, err = strconv.ParseInt(data, 10, 64)
	if err != nil {
		err = errors.Wrap(err, "DrakkinHeritage conversion")
		return
	}
	data, err = parser.EntryByName(line, columnNPCDrakkinTattoo)
	if err != nil {
		err = errors.Wrap(err, "DrakkinTattoo")
		return
	}
	npc.DrakkinTattoo, err = strconv.ParseInt(data, 10, 64)
	if err != nil {
		err = errors.Wrap(err, "DrakkinTattoo conversion")
		return
	}
	data, err = parser.EntryByName(line, columnNPCFace)
	if err != nil {
		err = errors.Wrap(err, "Face")
		return
	}
	npc.Face, err = strconv.ParseInt(data, 10, 64)
	if err != nil {
		err = errors.Wrap(err, "Face conversion")
		return
	}
	data, err = parser.EntryByName(line, columnNPCGender)
	if err != nil {
		err = errors.Wrap(err, "Gender")
		return
	}
	npc.Gender, err = strconv.ParseInt(data, 10, 64)
	if err != nil {
		err = errors.Wrap(err, "Gender conversion")
		return
	}
	data, err = parser.EntryByName(line, columnNPCLuclinBeard)
	if err != nil {
		err = errors.Wrap(err, "LuclinBeard")
		return
	}
	npc.LuclinBeard, err = strconv.ParseInt(data, 10, 64)
	if err != nil {
		err = errors.Wrap(err, "LuclinBeard conversion")
		return
	}
	data, err = parser.EntryByName(line, columnNPCLuclinBeardColor)
	if err != nil {
		err = errors.Wrap(err, "LuclinBeardColor")
		return
	}
	npc.LuclinBeardColor, err = strconv.ParseInt(data, 10, 64)
	if err != nil {
		err = errors.Wrap(err, "LuclinBeardColor conversion")
		return
	}
	data, err = parser.EntryByName(line, columnNPCLuclinEyeColor)
	if err != nil {
		err = errors.Wrap(err, "LuclinEyeColor")
		return
	}
	npc.LuclinEyeColor, err = strconv.ParseInt(data, 10, 64)
	if err != nil {
		err = errors.Wrap(err, "LuclinEyeColor conversion")
		return
	}
	data, err = parser.EntryByName(line, columnNPCLuclinEyeColor2)
	if err != nil {
		err = errors.Wrap(err, "LuclinEyeColor2")
		return
	}
	npc.LuclinEyeColor2, err = strconv.ParseInt(data, 10, 64)
	if err != nil {
		err = errors.Wrap(err, "LuclinEyeColor2 conversion")
		return
	}
	data, err = parser.EntryByName(line, columnNPCLuclinHairColor)
	if err != nil {
		err = errors.Wrap(err, "LuclinHairColor")
		return
	}
	npc.LuclinHairColor, err = strconv.ParseInt(data, 10, 64)
	if err != nil {
		err = errors.Wrap(err, "LuclinHairColor conversion")
		return
	}
	data, err = parser.EntryByName(line, columnNPCLuclinHairStyle)
	if err != nil {
		err = errors.Wrap(err, "LuclinHairStyle")
		return
	}
	npc.LuclinHairStyle, err = strconv.ParseInt(data, 10, 64)
	if err != nil {
		err = errors.Wrap(err, "LuclinHairStyle conversion")
		return
	}
	data, err = parser.EntryByName(line, columnNPCRace)
	if err != nil {
		err = errors.Wrap(err, "Race")
		return
	}
	npc.Race, err = strconv.ParseInt(data, 10, 64)
	if err != nil {
		err = errors.Wrap(err, "Race conversion")
		return
	}
	if npc.Race == 127 || npc.Race == 240 {
		npc.BodyType = 0
		npc.Trackable = 0
		npc.Findable = 0
		npc.SpecialAbilities = "24,35"
	}
	data, err = parser.EntryByName(line, columnNPCSize)
	if err != nil {
		err = errors.Wrap(err, "Size")
		return
	}
	npc.Size, err = strconv.ParseFloat(data, 64)
	if err != nil {
		err = errors.Wrap(err, "Size conversion")
		return
	}
	data, err = parser.EntryByName(line, columnNPCTexture)
	if err != nil {
		err = errors.Wrap(err, "Texture")
		return
	}
	npc.Texture, err = strconv.ParseInt(data, 10, 64)
	if err != nil {
		err = errors.Wrap(err, "Texture conversion")
		return
	}
	data, err = parser.EntryByName(line, columnNPCMeleeTexture1)
	if err != nil {
		err = errors.Wrap(err, "MeleeTexture1")
		return
	}
	npc.MeleeTexture1, err = strconv.ParseInt(data, 10, 64)
	if err != nil {
		err = errors.Wrap(err, "MeleeTexture1 conversion")
		return
	}
	data, err = parser.EntryByName(line, columnNPCMeleeTexture2)
	if err != nil {
		err = errors.Wrap(err, "MeleeTexture2")
		return
	}
	npc.MeleeTexture2, err = strconv.ParseInt(data, 10, 64)
	if err != nil {
		err = errors.Wrap(err, "MeleeTexture2 conversion")
		return
	}
	data, err = parser.EntryByName(line, columnNPCRunspeed)
	if err != nil {
		err = errors.Wrap(err, "Runspeed")
		return
	}
	npc.Runspeed, err = strconv.ParseFloat(data, 64)
	if err != nil {
		err = errors.Wrap(err, "Runspeed conversion")
		return
	}
	data, err = parser.EntryByName(line, columnNPCTargetable)
	if err != nil {
		err = errors.Wrap(err, "Targetable")
		return
	}
	if strings.ToLower(data) == "true" {
		npc.Targetable = 1
	}

	data, err = parser.EntryByName(line, columnNPCSeeInvis)
	if err != nil {
		err = errors.Wrap(err, "SeeInvis")
		return
	}
	npc.SeeInvis, err = strconv.ParseInt(data, 10, 64)
	if err != nil {
		err = errors.Wrap(err, "SeeInvis conversion")
		return
	}
	if npc.SeeInvis > 1 {
		npc.SeeInvis = 1
		//fmt.Println("warn: line", line, "see_invis is", npc.SeeInvis, "seems out of range, setting to 0")

		return
	}

	data, err = parser.EntryByName(line, columnNPCType)
	if err != nil {
		err = errors.Wrap(err, "type")
		return
	}
	npc.Type, err = strconv.ParseInt(data, 10, 64)
	if err != nil {
		err = errors.Wrap(err, "type conversion")
		return
	}

	data, err = parser.EntryByName(line, columnNPCPlayerState)
	if err != nil {
		err = errors.Wrap(err, "playerState")
		return
	}
	npc.PlayerState, err = strconv.ParseInt(data, 10, 64)
	if err != nil {
		err = errors.Wrap(err, "playerState conversion")
		return
	}

	data, err = parser.EntryByName(line, columnNPCDeity)
	if err != nil {
		err = errors.Wrap(err, "deity")
		return
	}
	npc.Deity, err = strconv.ParseInt(data, 10, 64)
	if err != nil {
		err = errors.Wrap(err, "deity conversion")
		return
	}

	sg := &model.Spawngroup{}
	data, err = parser.EntryByName(line, columnNPCX)
	if err != nil {
		err = errors.Wrap(err, "X position")
		return
	}
	sg.X, err = strconv.ParseFloat(data, 64)
	if err != nil {
		err = errors.Wrap(err, "X position conversion")
		return
	}
	data, err = parser.EntryByName(line, columnNPCY)
	if err != nil {
		err = errors.Wrap(err, "Y position")
		return
	}
	sg.Y, err = strconv.ParseFloat(data, 64)
	if err != nil {
		err = errors.Wrap(err, "Y position conversion")
		return
	}
	data, err = parser.EntryByName(line, columnNPCZ)
	if err != nil {
		err = errors.Wrap(err, "Z position")
		return
	}
	sg.Z, err = strconv.ParseFloat(data, 64)
	if err != nil {
		err = errors.Wrap(err, "Z position conversion")
		return
	}
	data, err = parser.EntryByName(line, columnNPCHeading)
	if err != nil {
		err = errors.Wrap(err, "heading")
		return
	}
	sg.Heading, err = strconv.ParseFloat(data, 64)
	if err != nil {
		err = errors.Wrap(err, "heading conversion")
		return
	}

	npc.Spawngroups = append(npc.Spawngroups, sg)
	return
}
