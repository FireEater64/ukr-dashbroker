package dashbroker

import (
	"time"

	log "github.com/cihub/seelog"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

type Housemate struct {
	ID          uint `gorm:"primary_key"`
	FirstName   string
	LastName    string
	PhoneNumber string
	Active      bool
}

type Button struct {
	ID         uint
	MacAddress string
	Name       string
}

type ButtonLogEntry struct {
	ID        uint   `gorm:"primary_key"`
	ButtonID  string `gorm:"ForeignKey:Button"`
	PressedAt time.Time
	Reason    string
}

func (ButtonLogEntry) TableName() string {
	return "buttonlog"
}

func LogButtonPress(macAddress string, reason string) {
	checkDBSession()

	toLog := ButtonLogEntry{ButtonID: macAddress, PressedAt: time.Now(), Reason: reason}

	db.Create(&toLog)
}

func GetAllActiveHousemates() []Housemate {
	var housemates []Housemate

	checkDBSession()
	db.Where("active = ?", true).Find(&housemates)

	return housemates
}

func GetAllButtons() []Button {
	var buttons []Button

	checkDBSession()
	db.Find(&buttons)

	return buttons
}

func newDatabaseSession() *gorm.DB {
	connection, err := gorm.Open(Configuration.DatabaseType, Configuration.DatabaseConnectionString)

	if err != nil {
		log.Criticalf("Error connecting to DB: %s", err)
		panic(err)
	}

	return connection
}

func checkDBSession() {
	if db == nil {
		db = newDatabaseSession()
	}
}
