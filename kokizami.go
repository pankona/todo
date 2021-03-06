package kokizami

import (
	"fmt"
	"time"

	// go-sqlite3 is imported only here
	_ "github.com/mattn/go-sqlite3"
)

// Kokizami represents a instance of kokizami
// Kokizami provides most APIs of kokizami library
type Kokizami struct {
	now func() time.Time

	KizamiRepo  KizamiRepository
	TagRepo     TagRepository
	SummaryRepo SummaryRepository
}

// initialTime is used to insert a time value that indicates initial value of time.
func initialTime() time.Time {
	t, err := time.Parse("2006-01-02 15:04:05", "1970-01-01 00:00:00")
	if err != nil {
		panic(fmt.Sprintf("failed to parse time for initial value for time: %v", err))
	}
	return t.UTC()
}

// Start starts a new kizami with specified desc
func (k *Kokizami) Start(desc string) (*Kizami, error) {
	if len(desc) == 0 {
		return nil, fmt.Errorf("desc must not be empty")
	}

	return k.KizamiRepo.Insert(desc)
}

// Get returns a Kizami by specified ID
func (k *Kokizami) Get(id int) (*Kizami, error) {
	return k.KizamiRepo.FindByID(id)
}

// Edit edits a specified kizami and update its model
func (k *Kokizami) Edit(ki *Kizami) (*Kizami, error) {
	m, err := k.KizamiRepo.FindByID(ki.ID)
	if err != nil {
		return nil, err
	}

	m.Desc = ki.Desc
	m.StartedAt = ki.StartedAt.UTC()
	m.StoppedAt = ki.StoppedAt.UTC()

	err = k.KizamiRepo.Update(m)
	if err != nil {
		return nil, err
	}

	return k.KizamiRepo.FindByID(ki.ID)
}

// Stop stops a on-going kizami by specified ID
func (k *Kokizami) Stop(id int) error {
	ki, err := k.KizamiRepo.FindByID(id)
	if err != nil {
		return err
	}
	ki.StoppedAt = k.now().UTC()
	return k.KizamiRepo.Update(ki)
}

// StopAll stops all on-going kizamis
func (k *Kokizami) StopAll() error {
	ks, err := k.KizamiRepo.FindByStoppedAt(initialTime())
	if err != nil {
		return err
	}
	now := k.now().UTC()
	for i := range ks {
		ks[i].StoppedAt = now
		if err := k.KizamiRepo.Update(ks[i]); err != nil {
			return err
		}
	}
	return nil
}

// Delete deletes a kizami by specified ID
func (k *Kokizami) Delete(id int) error {
	ki, err := k.KizamiRepo.FindByID(id)
	if err != nil {
		return err
	}
	return k.KizamiRepo.Delete(ki)
}

// List returns all Kizamis
func (k *Kokizami) List() ([]*Kizami, error) {
	return k.KizamiRepo.FindAll()
}

// SummaryByTag returns total elapsed time of Kizamis in specified month grouped by tag
func (k *Kokizami) SummaryByTag(yyyymm string) ([]*Elapsed, error) {
	// validate input
	_, err := time.Parse("2006-01", yyyymm)
	if err != nil {
		return nil, fmt.Errorf("invalid argument format. should be yyyy-mm: %v", err)
	}

	return k.SummaryRepo.ElapsedOfMonthByTag(yyyymm)
}

// SummaryByDesc returns total elapsed time of Kizamis in specified month grouped by desc
func (k *Kokizami) SummaryByDesc(yyyymm string) ([]*Elapsed, error) {
	// validate input
	_, err := time.Parse("2006-01", yyyymm)
	if err != nil {
		return nil, fmt.Errorf("invalid argument format. should be yyyy-mm: %v", err)
	}

	return k.SummaryRepo.ElapsedOfMonthByDesc(yyyymm)
}

// AddTags adds a new tags
func (k *Kokizami) AddTags(labels []string) error {
	return k.TagRepo.Insert(labels)
}

// DeleteTag deletes a specified tag
func (k *Kokizami) DeleteTag(id int) error {
	return k.TagRepo.Delete(id)
}

// Tags returns list of tags
func (k *Kokizami) Tags() ([]*Tag, error) {
	ms, err := k.TagRepo.FindAll()
	if err != nil {
		return nil, err
	}

	ts := make([]Tag, len(ms))
	for i := range ms {
		ts[i].ID = ms[i].ID
		ts[i].Label = ms[i].Label
	}

	ret := make([]*Tag, len(ts))
	for i := range ts {
		ret[i] = &ts[i]
	}

	return ret, nil
}

// Tagging makes relation between specified kizami and tags
func (k *Kokizami) Tagging(kizamiID int, tagIDs []int) error {
	return k.KizamiRepo.Tagging(kizamiID, tagIDs)
}

// Untagging removes all tags from specified kizami
func (k *Kokizami) Untagging(kizamiID int) error {
	return k.KizamiRepo.Untagging(kizamiID)
}

// TagsByKizamiID returns tags of specified kizami
func (k *Kokizami) TagsByKizamiID(kizamiID int) ([]*Tag, error) {
	ms, err := k.TagRepo.FindByKizamiID(kizamiID)
	if err != nil {
		return nil, err
	}

	ts := make([]Tag, len(ms))
	for i := range ms {
		ts[i].ID = ms[i].ID
		ts[i].Label = ms[i].Label
	}

	ret := make([]*Tag, len(ts))
	for i := range ts {
		ret[i] = &ts[i]
	}

	return ret, nil
}

// TagsByLabels returns tags by specified tags
func (k *Kokizami) TagsByLabels(labels []string) ([]*Tag, error) {
	return k.TagRepo.FindByLabels(labels)
}
