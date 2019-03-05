package takeadump

import (
	"encoding/csv"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"
)

const (
	columnNPCName             = "Name[0x40] /* ie priest_of_discord00 */"
	columnNPCLevel            = "Level"
	columnNPCClass            = "GetClass()"
	columnNPCLastName         = "Lastname[0x20]"
	columnNPCBodyType         = "BodyType /* this really should be renamed to charprops or something its broken anyway*/"
	columnNPCDrakkinDetails   = "mActorClient.Details"
	columnNPCDrakkinHeritage  = "mActorClient.Heritage"
	columnNPCDrakkinTattoo    = "mActorClient.Tattoo"
	columnNPCFace             = "mActorClient.FaceStyle"
	columnNPCGender           = "mActorClient.Gender"
	columnNPCLuclinBeard      = "mActorClient.FacialHair"
	columnNPCLuclinBeardColor = "mActorClient.FacialHairColor"
	columnNPCLuclinEyeColor   = "mActorClient.EyeColor1"
	columnNPCLuclinEyeColor2  = "mActorClient.EyeColor2"
	columnNPCLuclinHairColor  = "mActorClient.HairColor"
	columnNPCLuclinHairStyle  = "mActorClient.HairStyle"
	columnNPCRace             = "mActorClient.Race"
	columnNPCSize             = "Height"
	columnNPCTexture          = "Equipment.Chest.ID //idfile on Lucy"
	columnNPCMeleeTexture1    = "Equipment.Primary.ID //idfile on Lucy"
	columnNPCMeleeTexture2    = "Equipment.Offhand.ID //idfile on Lucy"
	columnNPCRunspeed         = "RunSpeed"
	columnNPCTargetable       = "Targetable /* true if mob is targetable */"
	columnNPCSeeInvis         = "SeeInvis[1]" //0 and 2, but monocle is using 1
	columnNPCX                = "X"
	columnNPCY                = "Y"
	columnNPCZ                = "Z"
	columnNPCHeading          = "Heading"
	columnNPCType             = "Type"
	columnNPCPlayerState      = "PlayerState /* 0=Idle 1=Open 2=WeaponSheathed 4=Aggressive 8=ForcedAggressive 0x10=InstrumentEquipped 0x20=Stunned 0x40=PrimaryWeaponEquipped 0x80=SecondaryWeaponEquipped */"
	columnNPCDeity            = "Deity"
)

// parser represents a csv loader
type parser struct {
	records  [][]string
	filename string
}

// Newparser will load a CSV
func (c *Client) newParser(filename string) (p *parser, zoneShortName string, zoneID int64, zoneInstanceVersion int64, err error) {
	p = &parser{}

	f, err := os.Open(filename)
	if err != nil {
		err = errors.Wrap(err, "open")
		return
	}

	fileName := path.Base(f.Name())
	if len(fileName) < 3 {
		err = fmt.Errorf("file name must be greater than 3 characters long")
		return
	}
	if !strings.Contains(fileName, "_") {
		err = fmt.Errorf("file name should have a _ seperator")
		return
	}
	zoneShortName = fileName[0:strings.Index(fileName, "_")]
	if len(zoneShortName) < 1 {
		err = fmt.Errorf("zone name failed to parse")
		return
	}

	row := c.db.QueryRow("SELECT zoneidnumber FROM zone WHERE short_name = ?", strings.ToLower(zoneShortName))
	err = row.Scan(&zoneID)
	if err != nil {
		return
	}
	fmt.Println("zone ID:", zoneID, ", name:", zoneShortName)

	r := csv.NewReader(f)
	p.records, err = r.ReadAll()
	if err != nil {
		err = errors.Wrap(err, "failed to parse csv")
		return
	}
	return
}

// TotalLines returns the total number of lines in the CSV file loaded
func (p *parser) TotalLines() (lines int, err error) {
	if len(p.records) == 0 {
		err = fmt.Errorf("not loaded")
		return
	}
	lines = len(p.records)
	return
}

// ColumnIndex returns the index of a specified column
func (p *parser) ColumnIndex(name string) (index int, err error) {
	if len(p.records) == 0 {
		err = fmt.Errorf("not loaded")
		return
	}
	var column string
	for index, column = range p.records[0] {
		if strings.ToLower(column) == strings.ToLower(name) {
			return
		}
	}
	err = fmt.Errorf("column with name %s not found", name)
	return
}

// EntryByIndex returns a column by their integer value
func (p *parser) EntryByIndex(line int, columnIndex int) (value string, err error) {
	if len(p.records) == 0 {
		err = fmt.Errorf("not loaded")
		return
	}
	if len(p.records) < line {
		err = fmt.Errorf("line is out of range")
		return
	}
	value = p.records[line][columnIndex]
	return
}

// EntryByName returns a column by their column name
func (p *parser) EntryByName(line int, columnName string) (value string, err error) {
	if len(p.records) == 0 {
		err = fmt.Errorf("not loaded")
		return
	}
	if len(p.records) < line {
		err = fmt.Errorf("line is out of range")
		return
	}
	columnIndex, err := p.ColumnIndex(columnName)
	if err != nil {
		return
	}

	value = p.records[line][columnIndex]
	return
}
