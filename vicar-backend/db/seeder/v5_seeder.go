package seeder

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"vicar-backend/db/entities"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type V5Seeder struct {
	AssetsRoot string
}

func NewV5Seeder(assetsRoot string) *V5Seeder {
	return &V5Seeder{AssetsRoot: assetsRoot}
}

func (s *V5Seeder) Seed(ctx context.Context, db *gorm.DB) error {
	if db == nil {
		return errors.New("db is nil")
	}

	seeded, err := s.isAlreadySeeded(ctx, db)
	if err != nil {
		return err
	}
	if seeded {
		return nil
	}

	assignments, err := s.loadAssignments()
	if err != nil {
		return err
	}

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		bookMap, err := s.seedBooks(tx)
		if err != nil {
			return err
		}

		if err := s.seedBloodPotencyTable(tx); err != nil {
			return err
		}

		disciplineMap, abilityMap, err := s.seedDisciplines(tx)
		if err != nil {
			return err
		}

		_, err = s.seedClans(tx, bookMap, assignments, disciplineMap)
		if err != nil {
			return err
		}

		if err := s.seedPredatorTypes(tx, bookMap, assignments); err != nil {
			return err
		}

		if err := s.seedTraitPacks(tx, bookMap, assignments); err != nil {
			return err
		}

		if err := s.seedBloodRituals(tx, bookMap); err != nil {
			return err
		}

		if err := s.seedOblivionCeremonies(tx, bookMap); err != nil {
			return err
		}

		if err := s.seedItems(tx); err != nil {
			return err
		}

		_ = abilityMap

		return nil
	})
}

func (s *V5Seeder) isAlreadySeeded(ctx context.Context, db *gorm.DB) (bool, error) {
	var count int64
	err := db.WithContext(ctx).
		Model(&entities.V5Book{}).
		Where("is_official = ?", true).
		Where("old_vicar_id IN ?", []int{1, 2, 3, 4, 5, 6, 7}).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count == 7, nil
}

func (s *V5Seeder) seedBooks(tx *gorm.DB) (map[int]uuid.UUID, error) {
	type bookDef struct {
		OldID int
		Name  string
	}

	books := []bookDef{
		{OldID: 1, Name: "Grundregelwerk"},
		{OldID: 2, Name: "Anarchen"},
		{OldID: 3, Name: "Chicago by Night"},
		{OldID: 4, Name: "Kompendium"},
		{OldID: 5, Name: "Cults of the Blood Gods"},
		{OldID: 6, Name: "Let the Streets Run Red"},
		{OldID: 7, Name: "Sonderregeln"},
	}

	out := make(map[int]uuid.UUID, len(books))

	for _, b := range books {
		rec := entities.V5Book{
			OldVicarID: uint(b.OldID),
			Name:       b.Name,
			IsOfficial: true,
		}
		if err := tx.Create(&rec).Error; err != nil {
			return nil, err
		}
		out[b.OldID] = rec.ID
	}

	return out, nil
}

func (s *V5Seeder) dePath(file string) string {
	return filepath.Join(s.AssetsRoot, "de-DE", file)
}

func (s *V5Seeder) metaPath() string {
	return filepath.Join(s.AssetsRoot, "Meta.json")
}

func readJSONFile[T any](path string) (T, error) {
	var zero T
	b, err := os.ReadFile(path)
	if err != nil {
		return zero, err
	}
	b = stripUTF8BOM(b)
	if err := json.Unmarshal(b, &zero); err != nil {
		return zero, err
	}
	return zero, nil
}

func stripUTF8BOM(b []byte) []byte {
	if len(b) >= 3 && b[0] == 0xEF && b[1] == 0xBB && b[2] == 0xBF {
		return b[3:]
	}
	return b
}

type vicarMeta struct {
	Languages []string        `json:"languages"`
	Books     []vicarMetaBook `json:"books"`
}

type vicarMetaBook struct {
	ID            int   `json:"id"`
	Clans         []int `json:"clans"`
	Merits        []int `json:"merits"`
	Backgrounds   []int `json:"backgrounds"`
	PredatorTypes []int `json:"predatorTypes"`
}

type v5BookAssignments struct {
	ClanToBook       map[int]int
	MeritToBook      map[int]int
	BackgroundToBook map[int]int
	PredatorToBook   map[int]int
}

func (s *V5Seeder) loadAssignments() (*v5BookAssignments, error) {
	meta, err := readJSONFile[vicarMeta](s.metaPath())
	if err != nil {
		return nil, err
	}

	asg := &v5BookAssignments{
		ClanToBook:       map[int]int{},
		MeritToBook:      map[int]int{},
		BackgroundToBook: map[int]int{},
		PredatorToBook:   map[int]int{},
	}

	putUnique := func(m map[int]int, key int, bookID int, what string) error {
		if existing, ok := m[key]; ok && existing != bookID {
			return fmt.Errorf("%s id %d appears in multiple books (%d and %d) but schema supports only one BookID", what, key, existing, bookID)
		}
		m[key] = bookID
		return nil
	}

	for _, b := range meta.Books {
		for _, id := range b.Clans {
			if err := putUnique(asg.ClanToBook, id, b.ID, "clan"); err != nil {
				return nil, err
			}
		}
		for _, id := range b.Merits {
			if err := putUnique(asg.MeritToBook, id, b.ID, "merit traitpack"); err != nil {
				return nil, err
			}
		}
		for _, id := range b.Backgrounds {
			if err := putUnique(asg.BackgroundToBook, id, b.ID, "background traitpack"); err != nil {
				return nil, err
			}
		}
		for _, id := range b.PredatorTypes {
			if err := putUnique(asg.PredatorToBook, id, b.ID, "predator type"); err != nil {
				return nil, err
			}
		}
	}

	return asg, nil
}

func (s *V5Seeder) resolveBookUUID(bookMap map[int]uuid.UUID, bookOldID *int) (uuid.UUID, error) {
	if bookOldID != nil {
		if u, ok := bookMap[*bookOldID]; ok {
			return u, nil
		}
	}
	u, ok := bookMap[1]
	if !ok {
		return uuid.Nil, fmt.Errorf("book 1 missing (Grundregelwerk)")
	}
	return u, nil
}

type oldBloodPotencyRow struct {
	Value                      int    `json:"value"`
	BleedingSpurt              int    `json:"bleedingSpurt"`
	HealedDamage               int    `json:"healedDamage"`
	DisciplineBonus            int    `json:"disciplineBonus"`
	RouseRepeatDisciplineLevel int    `json:"rouseRepeatDisciplineLevel"`
	BaneLevel                  int    `json:"baneLevel"`
	Pray                       string `json:"pray"`
}

func (s *V5Seeder) seedBloodPotencyTable(tx *gorm.DB) error {
	rows, err := readJSONFile[[]oldBloodPotencyRow](s.dePath("BloodPotencyTable.json"))
	if err != nil {
		return err
	}

	for _, r := range rows {
		rec := entities.V5BloodPotencyData{
			Value:                      r.Value,
			BleedingSpurt:              r.BleedingSpurt,
			HealedDamage:               r.HealedDamage,
			DisciplineBonus:            r.DisciplineBonus,
			RouseRepeatDisciplineLevel: r.RouseRepeatDisciplineLevel,
			BaneLevel:                  r.BaneLevel,
			Pray:                       r.Pray,
		}
		if err := tx.Create(&rec).Error; err != nil {
			return err
		}
	}

	return nil
}

type oldDisciplineAbility struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Combination *struct {
		ID    int `json:"id"`
		Level int `json:"level"`
	} `json:"combination"`
	Requirement     *int     `json:"requirement"`
	MinBloodPotency *int     `json:"minBloodPotency"`
	Summary         string   `json:"summary"`
	Costs           string   `json:"costs"`
	DiceSupplies    *string  `json:"diceSupplies"`
	System          string   `json:"system"`
	Alternatives    []string `json:"alternatives"`
	Duration        string   `json:"duration"`
}

type oldDiscipline struct {
	ID      int                               `json:"id"`
	Name    string                            `json:"name"`
	Summary *string                           `json:"summary"`
	Note    *string                           `json:"note"`
	Levels  map[string][]oldDisciplineAbility `json:"levels"`
}

type pendingCombination struct {
	AbilityOldID     int
	CombinationOldID int
	CombinationLevel int
}

func (s *V5Seeder) seedDisciplines(tx *gorm.DB) (map[int]uuid.UUID, map[int]uuid.UUID, error) {
	items, err := readJSONFile[[]oldDiscipline](s.dePath("Disciplines.json"))
	if err != nil {
		return nil, nil, err
	}

	discMap := make(map[int]uuid.UUID, len(items))
	abilityMap := make(map[int]uuid.UUID, 1024)
	pending := make([]pendingCombination, 0, 256)

	for _, d := range items {
		rec := entities.V5Discipline{
			OldVicarID: uint(d.ID),
			Name:       d.Name,
			Summary:    d.Summary,
			Note:       d.Note,
			IsHomebrew: false,
			Creator:    nil,
		}
		if err := tx.Create(&rec).Error; err != nil {
			return nil, nil, err
		}
		discMap[d.ID] = rec.ID

		levelKeys := make([]int, 0, len(d.Levels))
		for k := range d.Levels {
			n, err := strconv.Atoi(k)
			if err != nil {
				return nil, nil, fmt.Errorf("invalid discipline level key %q for discipline %d", k, d.ID)
			}
			levelKeys = append(levelKeys, n)
		}
		sort.Ints(levelKeys)

		for _, lvl := range levelKeys {
			abilities := d.Levels[strconv.Itoa(lvl)]
			for _, a := range abilities {
				var alts pq.StringArray
				if len(a.Alternatives) > 0 {
					alts = a.Alternatives
				} else {
					alts = pq.StringArray{}
				}

				ab := entities.V5DisciplineAbility{
					OldVicarID:       uint(a.ID),
					DisciplineID:     rec.ID,
					Level:            lvl,
					Name:             a.Name,
					CombinationRefID: nil,
					CombinationLevel: nil,
					Requirement:      a.Requirement,
					MinBloodPotency:  a.MinBloodPotency,
					Summary:          a.Summary,
					Costs:            a.Costs,
					DiceSupplies:     a.DiceSupplies,
					System:           a.System,
					Alternatives:     alts,
					Duration:         a.Duration,
				}

				if err := tx.Create(&ab).Error; err != nil {
					return nil, nil, err
				}

				abilityMap[a.ID] = ab.ID

				if a.Combination != nil {
					pending = append(pending, pendingCombination{
						AbilityOldID:     a.ID,
						CombinationOldID: a.Combination.ID,
						CombinationLevel: a.Combination.Level,
					})
				}
			}
		}
	}

	for _, p := range pending {
		abID, ok := abilityMap[p.AbilityOldID]
		if !ok {
			return nil, nil, fmt.Errorf("pending combination: missing ability old id %d", p.AbilityOldID)
		}
		combID, ok := abilityMap[p.CombinationOldID]
		if !ok {
			return nil, nil, fmt.Errorf("pending combination: missing combination ability old id %d", p.CombinationOldID)
		}
		level := p.CombinationLevel
		if err := tx.Model(&entities.V5DisciplineAbility{}).
			Where("id = ?", abID).
			Updates(map[string]any{
				"combination_ref_id": combID,
				"combination_level":  level,
			}).Error; err != nil {
			return nil, nil, err
		}
	}

	return discMap, abilityMap, nil
}

type oldClan struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Slogan      string  `json:"slogan"`
	Description string  `json:"description"`
	Curse       string  `json:"curse"`
	Symbol      *string `json:"symbol"`
	Disciplines []int   `json:"disciplines"`
}

func (s *V5Seeder) seedClans(tx *gorm.DB, bookMap map[int]uuid.UUID, asg *v5BookAssignments, disciplineMap map[int]uuid.UUID) (map[int]uuid.UUID, error) {
	items, err := readJSONFile[[]oldClan](s.dePath("Clans.json"))
	if err != nil {
		return nil, err
	}

	out := make(map[int]uuid.UUID, len(items))

	for _, c := range items {
		var bookOldIDPtr *int
		if v, ok := asg.ClanToBook[c.ID]; ok {
			bookOldIDPtr = &v
		}
		bookUUID, err := s.resolveBookUUID(bookMap, bookOldIDPtr)
		if err != nil {
			return nil, err
		}

		rec := entities.V5Clan{
			OldVicarID:  uint(c.ID),
			BookID:      bookUUID,
			Name:        c.Name,
			Slogan:      c.Slogan,
			Description: c.Description,
			Curse:       c.Curse,
			Symbol:      c.Symbol,
			IsHomebrew:  false,
			Creator:     nil,
		}

		if err := tx.Create(&rec).Error; err != nil {
			return nil, err
		}

		if len(c.Disciplines) > 0 {
			discUUIDs := make([]entities.V5Discipline, 0, len(c.Disciplines))
			for _, oldID := range c.Disciplines {
				u, ok := disciplineMap[oldID]
				if !ok {
					return nil, fmt.Errorf("unknown discipline old id %d for clan %d", oldID, c.ID)
				}
				discUUIDs = append(discUUIDs, entities.V5Discipline{ID: u})
			}
			if err := tx.Model(&rec).Association("Disciplines").Replace(discUUIDs); err != nil {
				return nil, err
			}
		}

		out[c.ID] = rec.ID
	}

	return out, nil
}

type oldPredatorAction struct {
	Description string          `json:"description"`
	Type        string          `json:"type"`
	Data        json.RawMessage `json:"data"`
	Restriction *struct {
		Type string          `json:"type"`
		Data json.RawMessage `json:"data"`
	} `json:"restriction"`
}

type oldPredatorType struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Restriction *struct {
		Type string          `json:"type"`
		Data json.RawMessage `json:"data"`
	} `json:"restriction"`
	Actions []oldPredatorAction `json:"actions"`
}

func (s *V5Seeder) seedPredatorTypes(tx *gorm.DB, bookMap map[int]uuid.UUID, asg *v5BookAssignments) error {
	items, err := readJSONFile[[]oldPredatorType](s.dePath("PredatorTypes.json"))
	if err != nil {
		return err
	}

	for _, p := range items {
		var bookOldIDPtr *int
		if v, ok := asg.PredatorToBook[p.ID]; ok {
			bookOldIDPtr = &v
		}
		bookUUID, err := s.resolveBookUUID(bookMap, bookOldIDPtr)
		if err != nil {
			return err
		}

		var restriction *entities.V5Restriction
		if p.Restriction != nil {
			rt := entities.RestrictionType(p.Restriction.Type)
			restriction = &entities.V5Restriction{
				Type: &rt,
				Data: datatypes.JSON(p.Restriction.Data),
			}
		}

		rec := entities.V5PredatorType{
			OldVicarID:  uint(p.ID),
			BookID:      bookUUID,
			Name:        p.Name,
			Description: p.Description,
			Restriction: restriction,
		}

		if err := tx.Create(&rec).Error; err != nil {
			return err
		}

		for _, a := range p.Actions {
			var ar *entities.V5Restriction
			if a.Restriction != nil {
				rt := entities.RestrictionType(a.Restriction.Type)
				ar = &entities.V5Restriction{
					Type: &rt,
					Data: datatypes.JSON(a.Restriction.Data),
				}
			}

			act := entities.V5PredatorAction{
				PredatorTypeID: rec.ID,
				Description:    a.Description,
				Type:           entities.PTActionType(a.Type),
				Data:           datatypes.JSON(a.Data),
				Restriction:    ar,
			}

			if err := tx.Create(&act).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

type oldTrait struct {
	ID           int    `json:"id"`
	Level        int    `json:"level"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	IsRepeatable bool   `json:"isRepeatable"`

	Actions     any `json:"actions"`
	Requirement any `json:"requirement"`

	RestrictRepeats *struct {
		Size   *int `json:"size"`
		Amount int  `json:"amount"`
	} `json:"restrictRepeats"`

	Restriction *struct {
		Type string          `json:"type"`
		Data json.RawMessage `json:"data"`
	} `json:"restriction"`
}

type oldTraitPack struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`

	Restriction *struct {
		Type string          `json:"type"`
		Data json.RawMessage `json:"data"`
	} `json:"restriction"`

	Advantages    []oldTrait `json:"advantages"`
	Disadvantages []oldTrait `json:"disadvantages"`
}

func (s *V5Seeder) seedTraitPacks(tx *gorm.DB, bookMap map[int]uuid.UUID, asg *v5BookAssignments) error {
	if err := s.seedTraitPackFile(tx, bookMap, asg, "Merits.json", entities.TraitPackTypeMerits); err != nil {
		return err
	}
	if err := s.seedTraitPackFile(tx, bookMap, asg, "Backgrounds.json", entities.TraitPackTypeBackgrounds); err != nil {
		return err
	}
	return nil
}

func (s *V5Seeder) seedTraitPackFile(tx *gorm.DB, bookMap map[int]uuid.UUID, asg *v5BookAssignments, filename string, packType entities.TraitPackType) error {
	packs, err := readJSONFile[[]oldTraitPack](s.dePath(filename))
	if err != nil {
		return err
	}

	for _, p := range packs {
		var bookOldIDPtr *int
		if packType == entities.TraitPackTypeMerits {
			if v, ok := asg.MeritToBook[p.ID]; ok {
				bookOldIDPtr = &v
			}
		} else {
			if v, ok := asg.BackgroundToBook[p.ID]; ok {
				bookOldIDPtr = &v
			}
		}

		bookUUID, err := s.resolveBookUUID(bookMap, bookOldIDPtr)
		if err != nil {
			return err
		}

		var restriction *entities.V5Restriction
		if p.Restriction != nil {
			rt := entities.RestrictionType(p.Restriction.Type)
			restriction = &entities.V5Restriction{
				Type: &rt,
				Data: datatypes.JSON(p.Restriction.Data),
			}
		}

		rec := entities.V5TraitPack{
			OldVicarID:   uint(p.ID),
			BookID:       bookUUID,
			Type:         packType,
			Name:         p.Name,
			Description:  p.Description,
			SpecialRules: entities.TraitSpecialRulesNone,
			Restriction:  restriction,
		}

		if err := tx.Create(&rec).Error; err != nil {
			return err
		}

		if err := s.seedTraitsForPack(tx, rec.ID, p.ID, p.Advantages, entities.TraitPackSideAdvantage); err != nil {
			return err
		}
		if err := s.seedTraitsForPack(tx, rec.ID, p.ID, p.Disadvantages, entities.TraitPackSideDisadvantage); err != nil {
			return err
		}
	}

	return nil
}

func (s *V5Seeder) seedTraitsForPack(tx *gorm.DB, packUUID uuid.UUID, packOldID int, traits []oldTrait, side entities.TraitPackSide) error {
	for _, t := range traits {
		oldComposite := uint(packOldID*1000 + t.ID)

		var restriction *entities.V5Restriction
		if t.Restriction != nil {
			rt := entities.RestrictionType(t.Restriction.Type)
			restriction = &entities.V5Restriction{
				Type: &rt,
				Data: datatypes.JSON(t.Restriction.Data),
			}
		}

		var repeatAmount *int
		var repeatSize *int
		if t.RestrictRepeats != nil {
			repeatAmount = &t.RestrictRepeats.Amount
			repeatSize = t.RestrictRepeats.Size
		}

		reqJSON, _ := json.Marshal(t.Requirement)
		actJSON, _ := json.Marshal(t.Actions)

		tr := entities.V5Trait{
			OldVicarID:   oldComposite,
			Level:        t.Level,
			Name:         t.Name,
			Description:  t.Description,
			IsRepeatable: t.IsRepeatable,
			RepeatAmount: repeatAmount,
			RepeatSize:   repeatSize,
			Requirement:  datatypes.JSON(reqJSON),
			Actions:      datatypes.JSON(actJSON),
			Restriction:  restriction,
		}

		if err := tx.Create(&tr).Error; err != nil {
			return err
		}

		link := entities.V5TraitPackTrait{
			TraitPackID: packUUID,
			TraitID:     tr.ID,
			Side:        side,
		}
		if err := tx.Create(&link).Error; err != nil {
			return err
		}
	}

	return nil
}

type oldBloodRitual struct {
	ID          int    `json:"id"`
	Level       int    `json:"level"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Ingredients string `json:"ingredients"`
	Execution   string `json:"execution"`
	System      string `json:"system"`
}

func (s *V5Seeder) seedBloodRituals(tx *gorm.DB, bookMap map[int]uuid.UUID) error {
	items, err := readJSONFile[[]oldBloodRitual](s.dePath("BloodRituals.json"))
	if err != nil {
		return err
	}

	bookUUID, ok := bookMap[1]
	if !ok {
		return fmt.Errorf("book 1 missing for blood rituals")
	}

	for _, r := range items {
		rec := entities.V5BloodRitual{
			OldVicarID:  uint(r.ID),
			BookID:      bookUUID,
			Level:       r.Level,
			Name:        r.Name,
			Description: r.Description,
			Ingredients: r.Ingredients,
			Execution:   r.Execution,
			System:      r.System,
		}
		if err := tx.Create(&rec).Error; err != nil {
			return err
		}
	}

	return nil
}

type oldOblivionCeremony struct {
	ID          int     `json:"id"`
	Level       int     `json:"level"`
	Name        string  `json:"name"`
	Cost        string  `json:"cost"`
	Roll        string  `json:"roll"`
	Summary     string  `json:"summary"`
	Requires    *int    `json:"requires"`
	Cult        *string `json:"cult"`
	Ingredients string  `json:"ingredients"`
	Execution   string  `json:"execution"`
	System      string  `json:"system"`
	Duration    *string `json:"duration"`
}

func (s *V5Seeder) seedOblivionCeremonies(tx *gorm.DB, bookMap map[int]uuid.UUID) error {
	items, err := readJSONFile[[]oldOblivionCeremony](s.dePath("OblivionCeremonies.json"))
	if err != nil {
		return err
	}

	bookUUID, ok := bookMap[5]
	if !ok {
		bookUUID, ok = bookMap[1]
		if !ok {
			return fmt.Errorf("book 1 missing for oblivion ceremonies")
		}
	}

	for _, r := range items {
		rec := entities.V5OblivionCeremony{
			OldVicarID:  uint(r.ID),
			BookID:      bookUUID,
			Level:       r.Level,
			Name:        r.Name,
			Cost:        r.Cost,
			Roll:        r.Roll,
			Summary:     r.Summary,
			Requires:    r.Requires,
			Cult:        r.Cult,
			Ingredients: r.Ingredients,
			Execution:   r.Execution,
			System:      r.System,
			Duration:    r.Duration,
		}
		if err := tx.Create(&rec).Error; err != nil {
			return err
		}
	}

	return nil
}

type oldItem struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

func (s *V5Seeder) seedItems(tx *gorm.DB) error {
	items, err := readJSONFile[[]oldItem](s.dePath("Items.json"))
	if err != nil {
		return err
	}

	for _, it := range items {
		rec := entities.BaseItem{
			Name:        it.Name,
			Description: it.Description,
			Category:    it.Category,
		}
		if err := tx.Create(&rec).Error; err != nil {
			return err
		}
	}

	return nil
}
