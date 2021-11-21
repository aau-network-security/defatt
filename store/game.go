package store

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

var (
	tagRawRegexp        = `^[a-z0-9][a-z0-9-]*[a-z0-9]$`
	tagRegex            = regexp.MustCompile(tagRawRegexp)
	ErrTagEmpty         = errors.New("tag cannot be empty")
	ErrUnknownToken     = errors.New("unknown token")
	ErrInvalidFlagValue = errors.New("incorrect value for flag")
	ErrUnknownChallenge = errors.New("unknown challenge")
)

type GameConfig struct {
	Name       string     `yaml:"name"`
	Tag        string     `yaml:"tag"`
	ScenarioID int        `yaml:"scenario_id"`
	StartedAt  *time.Time `yaml:"started-at, omitempty"`
	FinishedAt *time.Time `yaml:"finished-at, omitempty"`
}

type RawGameFile struct {
	GameConfig `yaml:",inline"`
}

type EmptyVarErr struct {
	Var  string
	Type string
}

type InvalidTagSyntaxErr struct {
	tag string
}

func (ite *InvalidTagSyntaxErr) Error() string {
	return fmt.Sprintf("Invalid syntax for tag \"%s\", allowed syntax: %s", ite.tag, tagRawRegexp)
}

func (eve *EmptyVarErr) Error() string {
	if eve.Type == "" {
		return fmt.Sprintf("%s cannot be empty", eve.Var)
	}

	return fmt.Sprintf("%s cannot be empty for %s", eve.Var, eve.Type)
}

func (g GameConfig) Validate() error {
	if g.Name == "" {
		return &EmptyVarErr{Var: "Name", Type: "Game"}
	}

	if g.Tag == "" {
		return &EmptyVarErr{Var: "Tag", Type: "Game"}
	}
	return nil
}

type Tag string

func (t Tag) Validate() error {
	s := string(t)
	if s == "" {
		return ErrTagEmpty
	}

	if !tagRegex.MatchString(s) {
		return &InvalidTagSyntaxErr{s}
	}

	return nil
}

type Challenge struct {
	OwnerID     string     `yaml:"-"`
	FlagTag     Tag        `yaml:"tag"`
	FlagValue   string     `yaml:"-"`
	CompletedAt *time.Time `yaml:"completed-at,omitempty"`
}

func NewTag(s string) (Tag, error) {
	t := Tag(s)
	if err := t.Validate(); err != nil {
		return "", err
	}

	return t, nil
}

type GameConfigStore interface {
	Read() GameConfig
	Finish(time.Time) error
}

type gameConfigStore struct {
	m     sync.Mutex
	conf  GameConfig
	hooks []func(GameConfig) error
}

func NewGameConfigStore(conf GameConfig, hooks ...func(GameConfig) error) *gameConfigStore {
	return &gameConfigStore{
		conf:  conf,
		hooks: hooks,
	}
}

func (es *gameConfigStore) Read() GameConfig {
	es.m.Lock()
	defer es.m.Unlock()

	return es.conf
}

func (es *gameConfigStore) Finish(t time.Time) error {
	es.m.Lock()
	defer es.m.Unlock()

	es.conf.FinishedAt = &t

	return es.runHooks()
}

func (es *gameConfigStore) runHooks() error {
	for _, h := range es.hooks {
		if err := h(es.conf); err != nil {
			return err
		}
	}

	return nil
}

type GameFileHub interface {
	CreateGameFile(GameConfig) (GameFile, error)
	GetUnfinishedGames() ([]GameFile, error)
}

type Gamefilehub struct {
	m    sync.Mutex
	path string
}

func NewGameFileHub(path string) (GameFileHub, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return nil, err
		}
	}

	return &Gamefilehub{
		path: path,
	}, nil
}

type Archiver interface {
	ArchiveDir() string
	Archive() error
}

type GameFile interface {
	GameConfigStore
	Archiver
}

type Gamefile struct {
	m        sync.Mutex
	file     RawGameFile
	dir      string
	filename string

	GameConfigStore
}

func NewGameFile(dir string, filename string, file RawGameFile) *Gamefile {
	ef := &Gamefile{
		dir:      dir,
		filename: filename,
		file:     file,
	}
	ef.GameConfigStore = NewGameConfigStore(file.GameConfig, ef.saveGameConfig)

	return ef
}

func (ef *Gamefile) save() error {
	bytes, err := yaml.Marshal(ef.file)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(ef.path(), bytes, 0644)
}

func (ef *Gamefile) delete() error {
	return os.Remove(ef.path())
}

func (ef *Gamefile) saveGameConfig(conf GameConfig) error {
	ef.m.Lock()
	defer ef.m.Unlock()

	ef.file.GameConfig = conf

	return ef.save()
}

func (ef *Gamefile) path() string {
	return filepath.Join(ef.dir, ef.filename)
}

func (ef *Gamefile) ArchiveDir() string {
	parts := strings.Split(ef.filename, ".")
	relativeDir := strings.Join(parts[:len(parts)-1], ".")
	return filepath.Join(ef.dir, relativeDir)
}

func (ef *Gamefile) Archive() error {
	ef.m.Lock()
	defer ef.m.Unlock()

	if _, err := os.Stat(ef.ArchiveDir()); os.IsNotExist(err) {
		if err := os.MkdirAll(ef.ArchiveDir(), os.ModePerm); err != nil {
			return err
		}
	}

	cpy := Gamefile{
		file:     ef.file,
		dir:      ef.ArchiveDir(),
		filename: "config.yml",
	}

	cpy.save()

	if err := ef.delete(); err != nil {
		log.Warn().Msgf("Failed to delete old Game file: %s", err)
	}

	return nil
}

func getFileNameForGame(path string, tag Tag) (string, error) {
	now := time.Now().Format("02-01-06")
	dirname := fmt.Sprintf("%s-%s", tag, now)
	filename := fmt.Sprintf("%s.yml", dirname)

	_, dirErr := os.Stat(filepath.Join(path, dirname))
	_, fileErr := os.Stat(filepath.Join(path, filename))

	if os.IsNotExist(fileErr) && os.IsNotExist(dirErr) {
		return filename, nil
	}

	for i := 1; i < 999; i++ {
		dirname := fmt.Sprintf("%s-%s-%d", tag, now, i)
		filename := fmt.Sprintf("%s.yml", dirname)

		_, dirErr := os.Stat(filepath.Join(path, dirname))
		_, fileErr := os.Stat(filepath.Join(path, filename))

		if os.IsNotExist(fileErr) && os.IsNotExist(dirErr) {
			return filename, nil
		}
	}

	return "", fmt.Errorf("unable to get filename for Game")
}

func (esh *Gamefilehub) CreateGameFile(conf GameConfig) (GameFile, error) {
	filename, err := getFileNameForGame(esh.path, Tag(conf.Tag))
	if err != nil {
		return nil, err
	}

	ef := NewGameFile(esh.path, filename, RawGameFile{GameConfig: conf})
	if err := ef.save(); err != nil {
		return nil, err
	}

	return ef, nil
}

func (esh *Gamefilehub) GetUnfinishedGames() ([]GameFile, error) {
	var games []GameFile
	err := filepath.Walk(esh.path, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".yml" {
			f, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			var ef RawGameFile
			err = yaml.Unmarshal(f, &ef)
			if err != nil {
				return err
			}

			if ef.FinishedAt == nil {
				dir, filename := filepath.Split(path)

				log.Debug().Str("name", ef.Name).Msg("Found unfinished Game")
				games = append(games, NewGameFile(dir, filename, ef))
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return games, nil
}
