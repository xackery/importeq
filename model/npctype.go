package model

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/pkg/errors"
)

// NpcType repesents npc entries
type NpcType struct {
	Name                string
	Level               int64
	MaxLevel            int64
	Class               int64
	LastName            string
	BodyType            int64
	DrakkinDetails      int64
	DrakkinHeritage     int64
	DrakkinTattoo       int64
	Face                int64
	Gender              int64
	LuclinBeard         int64
	LuclinBeardColor    int64
	LuclinEyeColor      int64
	LuclinEyeColor2     int64
	LuclinHairColor     int64
	LuclinHairStyle     int64
	Race                int64
	Size                float64
	Texture             int64
	MeleeTexture1       int64
	MeleeTexture2       int64
	Runspeed            float64
	Targetable          int64
	SeeInvis            int64
	Trackable           int64
	Findable            int64
	SpecialAbilities    string
	Spawngroups         []*Spawngroup
	ZoneShortName       string
	ZoneID              int64
	ZoneInstanceVersion int64
	Type                int64
	PlayerState         int64
	Deity               int64
	Sources             []string
}

// GenerateInsert will create an INSERT statement
func (npc *NpcType) GenerateInsert() (out string, err error) {
	insertString := `#{{.Name}} - {{range .Sources}}{{.}} {{end}}
SELECT @npc_type_id := MAX(id) + 1 FROM npc_types WHERE id > {{.ZoneID}}*1000-1 AND id < ({{.ZoneID}}+1)*1000;
INSERT INTO npc_types SET id=@npc_type_id, name='{{.Name}}', level='{{.Level}}', 
maxlevel='{{.MaxLevel}}', class='{{.Class}}', lastname='{{.LastName}}', 
bodytype='{{.BodyType}}', drakkin_details='{{.DrakkinDetails}}', drakkin_heritage='{{.DrakkinHeritage}}', 
drakkin_tattoo='{{.DrakkinTattoo}}', face='{{.Face}}', gender='{{.Gender}}', luclin_beard='{{.LuclinBeard}}', 
luclin_beardcolor='{{.LuclinBeardColor}}', luclin_eyecolor='{{.LuclinEyeColor}}', luclin_eyecolor2='{{.LuclinEyeColor2}}', 
luclin_haircolor='{{.LuclinHairColor}}', luclin_hairstyle='{{.LuclinHairStyle}}', race='{{.Race}}', size='{{.Size}}', 
texture='{{.Texture}}', d_melee_texture1='{{.MeleeTexture1}}', d_melee_texture2='{{.MeleeTexture2}}', 
runspeed='{{.Runspeed}}', targetable='{{.Targetable}}', see_invis='{{.SeeInvis}}', 
trackable='{{.Trackable}}', findable='{{.Findable}}', special_abilities='{{.SpecialAbilities}}';
`

	t, err := template.New("insert").Parse(insertString)
	if err != nil {
		err = errors.Wrap(err, "failed to prepare npc template")
		return
	}
	buf := bytes.NewBufferString(out)

	err = t.Execute(buf, npc)
	if err != nil {
		err = errors.Wrap(err, "failed to generate template for npc")
		return
	}

	counter := 1
	for _, sg := range npc.Spawngroups {
		sg.Name = fmt.Sprintf("%s_%s_shin_%d", npc.ZoneShortName, npc.Name, counter)
		sg.ZoneShortName = npc.ZoneShortName
		sg.ZoneInstanceVersion = npc.ZoneInstanceVersion
		insertString = "SELECT @spawngroup_id := MAX(id) + 1 FROM spawngroup;\n"
		insertString += "INSERT INTO spawngroup SET id=@spawngroup_id, name='{{.Name}}', delay='45000', mindelay='15000', despawn_timer='100';\n"
		insertString += "SELECT @spawn2_id := MAX(id) +1 FROM spawn2;\n"
		insertString += "INSERT INTO spawn2 SET id=@spawn2_id, spawngroupID=@spawngroup_id, zone='{{.ZoneShortName}}', version='{{.ZoneInstanceVersion}}', x='{{.X}}', y='{{.Y}}', z='{{.Z}}', heading='{{.Heading}}';\n"
		insertString += "INSERT INTO spawnentry SET id=@spawngroup_id, npcid=@npc_type_id, chance='100';\n"
		counter++
		t, err = template.New("insert").Parse(insertString)
		if err != nil {
			err = errors.Wrap(err, "failed to prepare npc template")
			return
		}

		err = t.Execute(buf, sg)
		if err != nil {
			err = errors.Wrap(err, "failed to generate template for npc")
			return
		}
		//out = buf.String()
	}

	out = buf.String()
	out += "\n"

	return
}
