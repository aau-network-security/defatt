package store

import (
	"errors"
	"regexp"
	"time"
)

var (
	tagRawRegexp        = `^[a-z0-9][a-z0-9-]*[a-z0-9]$`
	tagRegex            = regexp.MustCompile(tagRawRegexp)
	TagEmptyErr         = errors.New("Tag cannot be empty")
	TeamExistsErr       = errors.New("Team already exists")
	UnknownTeamErr      = errors.New("Unknown team")
	UnknownTokenErr     = errors.New("Unknown token")
	InvalidFlagValueErr = errors.New("Incorrect value for flag")
	UnknownChallengeErr = errors.New("Unknown challenge")
)

const (
	ID_KEY       = "I"
	TEAMNAME_KEY = "TN"
	token_key    = "testing"
)

type GameConfig struct {
	Name       string     `yaml:"name"`
	Tag        string     `yaml:"tag"`
	ScenarioID int        `yaml:"scenario_id"`
	StartedAt  *time.Time `yaml:"started-at, omitempty"`
	FinishedAt *time.Time `yaml:"finished-at, omitempty"`
}

//
// type RawGameFile struct {
// 	GameConfig `yaml:",inline"`
// 	Teams      []Team `yaml:"teams,omitempty"`
// }
//
// type EmptyVarErr struct {
// 	Var  string
// 	Type string
// }
//
// type InvalidTagSyntaxErr struct {
// 	tag string
// }
//
// func (ite *InvalidTagSyntaxErr) Error() string {
// 	return fmt.Sprintf("Invalid syntax for tag \"%s\", allowed syntax: %s", ite.tag, tagRawRegexp)
// }
//
// func (eve *EmptyVarErr) Error() string {
// 	if eve.Type == "" {
// 		return fmt.Sprintf("%s cannot be empty", eve.Var)
// 	}
//
// 	return fmt.Sprintf("%s cannot be empty for %s", eve.Var, eve.Type)
// }
//
// func (g GameConfig) Validate() error {
// 	if g.Name == "" {
// 		return &EmptyVarErr{Var: "Name", Type: "Game"}
// 	}
//
// 	if g.Tag == "" {
// 		return &EmptyVarErr{Var: "Tag", Type: "Game"}
// 	}
// 	return nil
// }
//
// type Tag string
//
// func (t Tag) Validate() error {
// 	s := string(t)
// 	if s == "" {
// 		return TagEmptyErr
// 	}
//
// 	if !tagRegex.MatchString(s) {
// 		return &InvalidTagSyntaxErr{s}
// 	}
//
// 	return nil
// }
//
// type Challenge struct {
// 	OwnerID     string     `yaml:"-"`
// 	FlagTag     Tag        `yaml:"tag"`
// 	FlagValue   string     `yaml:"-"`
// 	CompletedAt *time.Time `yaml:"completed-at,omitempty"`
// }
//
// func NewTag(s string) (Tag, error) {
// 	t := Tag(s)
// 	if err := t.Validate(); err != nil {
// 		return "", err
// 	}
//
// 	return t, nil
// }
//
// type Team struct {
// 	ID               string            `yaml:"id"`
// 	Email            string            `yaml:"email"`
// 	Name             string            `yaml:"name"`
// 	HashedPassword   string            `yaml:"hashed-password"`
// 	SolvedChallenges []Challenge       `yaml:"solved-challenges,omitempty"`
// 	Metadata         map[string]string `yaml:"metadata,omitempty"`
// 	CreatedAt        *time.Time        `yaml:"created-at,omitempty"`
// 	ChalMap          map[Tag]Challenge `yaml:"-"`
// 	RedTeam          bool              `yaml:"is-red-teamm"`
// 	VPNConfig        string            // todo: this is just a temp
// }
//
// func NewTeam(email, name, password string, chals ...Challenge) Team {
// 	now := time.Now()
//
// 	hashedPassword := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
// 	email = strings.ToLower(email)
//
// 	t := Team{
// 		ID:             uuid.New().String()[0:8],
// 		Email:          email,
// 		Name:           name,
// 		HashedPassword: hashedPassword,
// 		CreatedAt:      &now,
// 	}
// 	for _, chal := range chals {
// 		t.AddChallenge(chal)
// 		log.Info().Str("chal-tag", string(chal.FlagTag)).
// 			Str("chal-val", chal.FlagValue).
// 			Msgf("Flag is created for team %s ", t.Name)
// 	}
// 	return t
// }
//
// func (t *Team) IsCorrectFlag(tag Tag, v string) error {
// 	c, ok := t.ChalMap[tag]
// 	if !ok {
// 		return UnknownChallengeErr
// 	}
//
// 	if c.FlagValue != v {
// 		return InvalidFlagValueErr
// 	}
//
// 	return nil
// }
//
// func (t *Team) SolveChallenge(tag Tag, v string) error {
// 	now := time.Now()
//
// 	if err := t.IsCorrectFlag(tag, v); err != nil {
// 		return err
// 	}
//
// 	c := t.ChalMap[tag]
// 	c.CompletedAt = &now
//
// 	t.SolvedChallenges = append(t.SolvedChallenges, c)
// 	t.AddChallenge(c)
//
// 	return nil
// }
//
// func (t *Team) AddMetadata(key, value string) {
// 	if t.Metadata == nil {
// 		t.Metadata = map[string]string{}
// 	}
// 	t.Metadata[key] = value
// }
//
// func (t *Team) DataCollection() bool {
// 	if t.Metadata == nil {
// 		return false
// 	}
//
// 	v, ok := t.Metadata["consent"]
// 	if !ok {
// 		return false
// 	}
//
// 	return v == "ok"
// }
//
// func (t *Team) AddChallenge(c Challenge) {
// 	if t.ChalMap == nil {
// 		t.ChalMap = map[Tag]Challenge{}
// 	}
// 	t.ChalMap[c.FlagTag] = c
// }
//
// func (t *Team) DataConsent() bool {
// 	if t.Metadata == nil {
// 		return false
// 	}
// 	v, ok := t.Metadata["consent"]
// 	if !ok {
// 		return false
// 	}
// 	return v == "ok"
// }
//
// type TeamStore interface {
// 	CreateTeam(Team) error
// 	GetTeamByToken(string) (Team, error)
// 	GetTeamByEmail(string) (Team, error)
// 	GetTeamByName(string) (Team, error)
// 	GetTeams() []Team
// 	SaveTeam(Team) error
// 	CreateTokenForTeam(string, Team) error
// 	SaveTokenForTeam(token string, in *Team) error
// 	DeleteToken(string) error
// }
//
// func GetTokenForTeam(key []byte, t *Team) (string, error) {
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		ID_KEY:       t.ID,
// 		TEAMNAME_KEY: t.Name,
// 	})
// 	tokenStr, err := token.SignedString(key)
// 	if err != nil {
// 		return "", err
// 	}
// 	return tokenStr, nil
// }
//
// type teamstore struct {
// 	m sync.RWMutex
//
// 	hooks  []func([]Team) error
// 	teams  map[string]Team
// 	tokens map[string]string
// 	emails map[string]string
// 	names  map[string]string
// }
//
// type TeamStoreOpt func(ts *teamstore)
//
// func WithTeams(teams []Team) func(ts *teamstore) {
// 	return func(ts *teamstore) {
// 		for _, t := range teams {
// 			ts.CreateTeam(t)
// 		}
// 	}
// }
//
// func WithPostTeamHook(hook func(teams []Team) error) func(ts *teamstore) {
// 	return func(ts *teamstore) {
// 		ts.hooks = append(ts.hooks, hook)
// 	}
// }
// func (es *teamstore) SaveTokenForTeam(token string, in *Team) error {
// 	es.m.Lock()
// 	defer es.m.Unlock()
// 	if token == "" {
// 		return &EmptyVarErr{Var: "Token"}
// 	}
// 	if in.ID == "" {
// 		return fmt.Errorf("SaveTokenForTeam function error %v", UnknownTeamErr)
// 	}
// 	es.tokens[token] = in.ID
// 	return nil
// }
//
// func NewTeamStore(opts ...TeamStoreOpt) *teamstore {
// 	ts := &teamstore{
// 		hooks:  []func(teams []Team) error{},
// 		teams:  map[string]Team{},
// 		tokens: map[string]string{},
// 		names:  map[string]string{},
// 		emails: map[string]string{},
// 	}
//
// 	for _, opt := range opts {
// 		opt(ts)
// 	}
//
// 	return ts
// }
//
// func (es *teamstore) CreateTeam(t Team) error {
// 	es.m.Lock()
// 	defer es.m.Unlock()
//
// 	if _, ok := es.teams[t.ID]; ok {
// 		return TeamExistsErr
// 	}
//
// 	es.teams[t.ID] = t
// 	es.emails[t.Email] = t.ID
// 	es.names[t.Name] = t.ID
//
// 	return es.RunHooks()
// }
//
// func (es *teamstore) SaveTeam(t Team) error {
// 	es.m.Lock()
// 	defer es.m.Unlock()
//
// 	if _, ok := es.teams[t.ID]; !ok {
// 		return UnknownTeamErr
// 	}
//
// 	es.teams[t.ID] = t
//
// 	return es.RunHooks()
// }
//
// func (es *teamstore) CreateTokenForTeam(token string, in Team) error {
// 	es.m.Lock()
// 	defer es.m.Unlock()
//
// 	if token == "" {
// 		return &EmptyVarErr{Var: "Token"}
// 	}
//
// 	t, ok := es.teams[in.ID]
// 	if !ok {
// 		return UnknownTeamErr
// 	}
//
// 	es.tokens[token] = t.ID
//
// 	return nil
// }
//
// func (es *teamstore) DeleteToken(token string) error {
// 	es.m.Lock()
// 	defer es.m.Unlock()
//
// 	delete(es.tokens, token)
//
// 	return nil
// }
//
// func (es *teamstore) GetTeams() []Team {
// 	var teams []Team
// 	for _, t := range es.teams {
// 		teams = append(teams, t)
// 	}
//
// 	return teams
// }
//
// func (es *teamstore) GetTeamByEmail(email string) (Team, error) {
// 	es.m.RLock()
// 	defer es.m.RUnlock()
//
// 	id, ok := es.emails[email]
// 	if !ok {
// 		return Team{}, UnknownTokenErr
// 	}
//
// 	t, ok := es.teams[id]
// 	if !ok {
// 		return Team{}, UnknownTeamErr
// 	}
//
// 	return t, nil
// }
//
// func (es *teamstore) GetTeamByName(name string) (Team, error) {
// 	es.m.RLock()
// 	defer es.m.RUnlock()
//
// 	id, ok := es.names[name]
// 	if !ok {
// 		return Team{}, UnknownTokenErr
// 	}
//
// 	t, ok := es.teams[id]
// 	if !ok {
// 		return Team{}, UnknownTeamErr
// 	}
//
// 	return t, nil
// }
//
// func (es *teamstore) GetTeamByToken(token string) (Team, error) {
// 	es.m.RLock()
// 	defer es.m.RUnlock()
//
// 	id, ok := es.tokens[token]
// 	if !ok {
// 		return Team{}, UnknownTokenErr
// 	}
//
// 	t, ok := es.teams[id]
// 	if !ok {
// 		return Team{}, UnknownTeamErr
// 	}
//
// 	return t, nil
// }
//
// func (es *teamstore) RunHooks() error {
// 	teams := es.GetTeams()
// 	for _, h := range es.hooks {
// 		if err := h(teams); err != nil {
// 			return err
// 		}
// 	}
//
// 	return nil
// }
//
// type GameConfigStore interface {
// 	Read() GameConfig
// 	Finish(time.Time) error
// }
//
// type gameConfigStore struct {
// 	m     sync.Mutex
// 	conf  GameConfig
// 	hooks []func(GameConfig) error
// }
//
// func NewGameConfigStore(conf GameConfig, hooks ...func(GameConfig) error) *gameConfigStore {
// 	return &gameConfigStore{
// 		conf:  conf,
// 		hooks: hooks,
// 	}
// }
//
// func (es *gameConfigStore) Read() GameConfig {
// 	es.m.Lock()
// 	defer es.m.Unlock()
//
// 	return es.conf
// }
//
// func (es *gameConfigStore) Finish(t time.Time) error {
// 	es.m.Lock()
// 	defer es.m.Unlock()
//
// 	es.conf.FinishedAt = &t
//
// 	return es.runHooks()
// }
//
// func (es *gameConfigStore) runHooks() error {
// 	for _, h := range es.hooks {
// 		if err := h(es.conf); err != nil {
// 			return err
// 		}
// 	}
//
// 	return nil
// }
//
// type GameFileHub interface {
// 	CreateGameFile(GameConfig) (GameFile, error)
// 	GetUnfinishedGames() ([]GameFile, error)
// }
//
// type Gamefilehub struct {
// 	m    sync.Mutex
// 	path string
// }
//
// func NewGameFileHub(path string) (GameFileHub, error) {
// 	if _, err := os.Stat(path); os.IsNotExist(err) {
// 		if err := os.MkdirAll(path, os.ModePerm); err != nil {
// 			return nil, err
// 		}
// 	}
//
// 	return &Gamefilehub{
// 		path: path,
// 	}, nil
// }
//
// type Archiver interface {
// 	ArchiveDir() string
// 	Archive() error
// }
//
// type GameFile interface {
// 	TeamStore
// 	GameConfigStore
// 	Archiver
// }
//
// type Gamefile struct {
// 	m        sync.Mutex
// 	file     RawGameFile
// 	dir      string
// 	filename string
//
// 	TeamStore
// 	GameConfigStore
// }
//
// func NewGameFile(dir string, filename string, file RawGameFile) *Gamefile {
// 	ef := &Gamefile{
// 		dir:      dir,
// 		filename: filename,
// 		file:     file,
// 	}
//
// 	ef.TeamStore = NewTeamStore(WithTeams(file.Teams), WithPostTeamHook(ef.saveTeams))
// 	ef.GameConfigStore = NewGameConfigStore(file.GameConfig, ef.saveGameConfig)
//
// 	return ef
// }
//
// func (ef *Gamefile) save() error {
// 	bytes, err := yaml.Marshal(ef.file)
// 	if err != nil {
// 		return err
// 	}
//
// 	return ioutil.WriteFile(ef.path(), bytes, 0644)
// }
//
// func (ef *Gamefile) delete() error {
// 	return os.Remove(ef.path())
// }
//
// func (ef *Gamefile) saveTeams(teams []Team) error {
// 	ef.m.Lock()
// 	defer ef.m.Unlock()
//
// 	ef.file.Teams = teams
//
// 	return ef.save()
// }
//
// func (ef *Gamefile) saveGameConfig(conf GameConfig) error {
// 	ef.m.Lock()
// 	defer ef.m.Unlock()
//
// 	ef.file.GameConfig = conf
//
// 	return ef.save()
// }
//
// func (ef *Gamefile) path() string {
// 	return filepath.Join(ef.dir, ef.filename)
// }
//
// func (ef *Gamefile) ArchiveDir() string {
// 	parts := strings.Split(ef.filename, ".")
// 	relativeDir := strings.Join(parts[:len(parts)-1], ".")
// 	return filepath.Join(ef.dir, relativeDir)
// }
//
// func (ef *Gamefile) Archive() error {
// 	ef.m.Lock()
// 	defer ef.m.Unlock()
//
// 	if _, err := os.Stat(ef.ArchiveDir()); os.IsNotExist(err) {
// 		if err := os.MkdirAll(ef.ArchiveDir(), os.ModePerm); err != nil {
// 			return err
// 		}
// 	}
//
// 	cpy := Gamefile{
// 		file:     ef.file,
// 		dir:      ef.ArchiveDir(),
// 		filename: "config.yml",
// 	}
//
// 	cpy.file.Teams = []Team{}
// 	for _, t := range ef.GetTeams() {
// 		t.Name = ""
// 		t.Email = ""
// 		t.HashedPassword = ""
// 		cpy.file.Teams = append(cpy.file.Teams, t)
// 	}
// 	cpy.save()
//
// 	if err := ef.delete(); err != nil {
// 		log.Warn().Msgf("Failed to delete old Game file: %s", err)
// 	}
//
// 	return nil
// }
//
// func getFileNameForGame(path string, tag Tag) (string, error) {
// 	now := time.Now().Format("02-01-06")
// 	dirname := fmt.Sprintf("%s-%s", tag, now)
// 	filename := fmt.Sprintf("%s.yml", dirname)
//
// 	_, dirErr := os.Stat(filepath.Join(path, dirname))
// 	_, fileErr := os.Stat(filepath.Join(path, filename))
//
// 	if os.IsNotExist(fileErr) && os.IsNotExist(dirErr) {
// 		return filename, nil
// 	}
//
// 	for i := 1; i < 999; i++ {
// 		dirname := fmt.Sprintf("%s-%s-%d", tag, now, i)
// 		filename := fmt.Sprintf("%s.yml", dirname)
//
// 		_, dirErr := os.Stat(filepath.Join(path, dirname))
// 		_, fileErr := os.Stat(filepath.Join(path, filename))
//
// 		if os.IsNotExist(fileErr) && os.IsNotExist(dirErr) {
// 			return filename, nil
// 		}
// 	}
//
// 	return "", fmt.Errorf("unable to get filename for Game")
// }
//
// func (esh *Gamefilehub) CreateGameFile(conf GameConfig) (GameFile, error) {
// 	filename, err := getFileNameForGame(esh.path, Tag(conf.Tag))
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	ef := NewGameFile(esh.path, filename, RawGameFile{GameConfig: conf})
// 	if err := ef.save(); err != nil {
// 		return nil, err
// 	}
//
// 	return ef, nil
// }
//
// func (esh *Gamefilehub) GetUnfinishedGames() ([]GameFile, error) {
// 	var games []GameFile
// 	err := filepath.Walk(esh.path, func(path string, info os.FileInfo, err error) error {
// 		if filepath.Ext(path) == ".yml" {
// 			f, err := ioutil.ReadFile(path)
// 			if err != nil {
// 				return err
// 			}
//
// 			var ef RawGameFile
// 			err = yaml.Unmarshal(f, &ef)
// 			if err != nil {
// 				return err
// 			}
//
// 			if ef.FinishedAt == nil {
// 				dir, filename := filepath.Split(path)
//
// 				log.Debug().Str("name", ef.Name).Msg("Found unfinished Game")
// 				games = append(games, NewGameFile(dir, filename, ef))
// 			}
// 		}
//
// 		return nil
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return games, nil
// }
