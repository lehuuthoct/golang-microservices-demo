package db

import (
	"lehuuthoct/lht-microservices/leader/model"
	"github.com/boltdb/bolt"
	"github.com/sirupsen/logrus"
	"encoding/json"
	"fmt"
	"time"
	"strconv"
)

const LEADER_BUCKET = "LEADER_BUCKET"
type IBoltClient interface {
	OpenBoltDB()
	FindLeaderById(leaderID string) (model.Leader, error)
	MockLeaderData()
	Check() bool
}

type BoltClient struct {
	boltDB *bolt.DB
}

// init leaders.db data file in current directory
func (bc *BoltClient) OpenBoltDB()  {
	var err error
	options := &bolt.Options{Timeout: 10 * time.Second}
	bc.boltDB, err = bolt.Open("leaders.db", 0600, options)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Infof("Init BoltDB Successfully [%#v] \n", bc.boltDB.Path())
}

// mock 10 leaders
func (bc *BoltClient) MockLeaderData() {
	logrus.Info("MockLeaderData")

	bc.initLeaderBucket()
	bc.initMockData()

	logrus.Printf("Mocked %d leaders successfully \n", 10)
}

func (bc *BoltClient) initMockData() {
	for i := 0; i < 10; i++  {
		key := strconv.Itoa(1 + i)
		leader := model.Leader{
			Id: key,
			Name: "Leader_" + key,
		}
		//	serialize leader data to []bytes
		jsonLeaderBytes, err := json.Marshal(leader)
		if err != nil {
			logrus.Info("error marshalling leader info ", err)
		}

		bc.boltDB.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(LEADER_BUCKET))

			err := b.Put([]byte(key), jsonLeaderBytes)

			return err
		})
	}
}

func (bc *BoltClient) initLeaderBucket() {
	bc.boltDB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(LEADER_BUCKET))
		if err != nil {
			return fmt.Errorf("Error creating %s bucket [%s]", LEADER_BUCKET, err)
		}
		return nil
	})
}

func (bc *BoltClient) FindLeaderById(leaderID string) (model.Leader, error)  {
	//	init leader
	leader := model.Leader{}

	// read leader from Bolt Bucket
	err := bc.boltDB.View(func(tx *bolt.Tx) error {
		// init bucket
		b := tx.Bucket([]byte(LEADER_BUCKET))

		//	get leader by id
		leaderBytes := b.Get([]byte(leaderID))
		if leaderBytes == nil {
			return fmt.Errorf("No Leader is found for %s", leaderID)
		}

		//	convert leader bytes to leader struct
		json.Unmarshal(leaderBytes, &leader)

		//	return no error
		return nil
	})

	//	case 1 return error
	if err != nil {
		return model.Leader{}, err
	}

	//	case 2 return leader & no error
	return leader, nil
}

func (bc *BoltClient) Check() bool  {
	return bc.boltDB != nil
}
