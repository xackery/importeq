package takeadump

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/xackery/eqemuconfig"
	"github.com/xackery/importeq/model"

	// mysql is used for support
	_ "github.com/go-sql-driver/mysql"
)

// Client represents a takeadump parser
type Client struct {
	config              *eqemuconfig.Config
	db                  *sql.DB
	ZoneShortName       string
	ZoneID              int64
	ZoneInstanceVersion int64
	Npcs                map[string]*model.NpcType
}

// New creates a new takeadump client based on provided file
func New(dirname string, zoneShortName string) (c *Client, err error) {

	fmt.Println("----TakeADump----")
	c = &Client{
		Npcs: map[string]*model.NpcType{},
	}
	c.config, err = eqemuconfig.GetConfig()
	if err != nil {
		return
	}
	c.db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", c.config.Database.Username, c.config.Database.Password, c.config.Database.Host, c.config.Database.Port, c.config.Database.Db))
	if err != nil {
		err = errors.Wrap(err, "failed to open sql")
		return
	}

	fmt.Println("walking", dirname)

	err = filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			err = errors.Wrapf(err, "failed to access %s", path)
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".csv" {
			return nil
		}

		if !strings.Contains(strings.ToLower(path), fmt.Sprintf("%s_npc_", zoneShortName)) {
			return nil
		}

		err = c.parse(path)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return
	}

	fmt.Printf("%d total npcs extracted\n", len(c.Npcs))

	return
}

func (c *Client) parse(filename string) (err error) {
	parser, zoneShortName, zoneID, zoneInstanceVersion, err := c.newParser(filename)
	if err != nil {
		err = errors.Wrapf(err, "failed to load %s", filename)
		return
	}

	isNewNpc := false
	totalLines, err := parser.TotalLines()
	if err != nil {
		return
	}
	if totalLines < 3 {
		err = fmt.Errorf("empty file")
		return
	}
	var npc *model.NpcType
	fmt.Println("source:", filename, "has", totalLines-2, "total lines")
	for line := 0; line < totalLines; line++ {
		if line < 2 {
			continue
		}
		npc, err = newNpcType(parser, line, zoneShortName, zoneID, zoneInstanceVersion)
		if err != nil {
			err = errors.Wrapf(err, "failed to parse npc on line %d", line)
			return
		}
		if strings.Contains(npc.LastName, "'s Mercenary") {
			//fmt.Printf("%d: skipped %s (mercenary)\n", line, npc.Name)
			continue
		}
		if strings.Contains(npc.LastName, "'s Pet") {
			//fmt.Printf("%d: skipped %s (pet)\n", line, npc.Name)
			continue
		}
		if strings.Contains(npc.LastName, "`s Mount") {
			//fmt.Printf("%d: skipped %s (mount)\n", line, npc.Name)
			continue
		}
		if npc.Deity > 0 {
			fmt.Printf("%d: skipped %s (player)\n", line, npc.Name)
			continue
		}
		if npc.Type == 0 {
			//fmt.Printf("%d: skipped %s (corpse)\n", line, npc.Name)
			//continue
		}

		isNewNpc = true
		for _, oldNpc := range c.Npcs {
			if npc.Name == oldNpc.Name {
				isNewNpc = false
				if npc.Level < oldNpc.Level {
					oldNpc.Level = npc.Level
				}
				if npc.MaxLevel > oldNpc.MaxLevel {
					oldNpc.MaxLevel = npc.MaxLevel
				}
				for _, sg := range npc.Spawngroups {
					oldNpc.Spawngroups = append(oldNpc.Spawngroups, sg)
				}
				for _, source := range npc.Sources {
					oldNpc.Sources = append(oldNpc.Sources, source)
				}
				break
			}
		}
		if !isNewNpc {
			continue
		}

		c.Npcs[npc.Name] = npc
	}
	return
}

// Close releases the database connection and cleans up any lingering data
func (c *Client) Close() (err error) {
	if c.db != nil {
		err = c.db.Close()
		c.db = nil
	}

	return
}
