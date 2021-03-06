package repository

import (
	"crypto/md5"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/globalsign/mgo/bson"

	radiobot "github.com/Oleg-MBO/Radio-en-Esperanto"
	"github.com/globalsign/mgo"
)

type mongoPodcastRepository struct {
	Collection *mgo.Collection
}

// NewMongoPodcastRepository mongo repository for podcasts
func NewMongoPodcastRepository(collection *mgo.Collection) radiobot.PodcastRepository {
	return &mongoPodcastRepository{Collection: collection}
}

// AddPocast add new podcast to DB
func (rpod *mongoPodcastRepository) AddPocast(p radiobot.Podcast) error {
	p.CalcID()
	if p.ChannelID == (uuid.UUID{}) {
		return fmt.Errorf("channel id can`t be empty")
	}
	if p.ParsedOn == (time.Time{}) || p.CreatedOn == (time.Time{}) {
		return fmt.Errorf("time values must be not empty")
	}
	return rpod.Collection.Insert(p)
}

// IsNewPodcast check if podcast exist
func (rpod *mongoPodcastRepository) IsNewPodcast(p radiobot.Podcast) (bool, error) {
	p.CalcID()
	n, err := rpod.Collection.Find(bson.M{"_id": p.ID}).Limit(1).Count()
	if err != nil {
		return false, err
	}
	if n > 0 {
		return false, nil
	}
	return true, nil
}

// GetUnsendedPodcasts check if podcast exist
func (rpod *mongoPodcastRepository) FindPodcastByID(p *radiobot.Podcast) error {
	if p.ID == ([md5.Size]byte{}) {
		p.CalcID()
	}
	return rpod.Collection.FindId(p.ID).One(p)
}

// GetUnsendedPodcasts check if podcast exist
func (rpod *mongoPodcastRepository) FindUnsendedPodcasts(count, offset int) ([]radiobot.Podcast, error) {
	podcastList := make([]radiobot.Podcast, 0, count)
	err := rpod.Collection.Find(bson.M{"recipient": ""}).Sort("created_on").Skip(offset).Limit(count).All(&podcastList)
	return podcastList, err
}

// UpdatePodcast update podcast in db
func (rpod *mongoPodcastRepository) UpdatePodcast(p radiobot.Podcast) error {
	p.CalcID()
	if p.ChannelID == (uuid.UUID{}) {
		return fmt.Errorf("channel id can`t be empty")
	}
	if p.ParsedOn == (time.Time{}) || p.CreatedOn == (time.Time{}) {
		return fmt.Errorf("time values must be not empty")
	}
	if !p.IsSended() {
		return fmt.Errorf("podcast must be sended")
	}

	return rpod.Collection.UpdateId(p.ID, p)
}
