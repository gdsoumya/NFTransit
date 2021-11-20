package database

import (
	"time"

	"gorm.io/gorm"
)

type IDatabase interface {
	Insert(tx *gorm.DB, payload *MintRequest) error
	Find(tx *gorm.DB, payload *MintRequest) ([]MintRequest, error)
	Update(tx *gorm.DB, payload *MintRequest, all bool) error
	Delete(tx *gorm.DB, payload *MintRequest) error
	Transaction(tx func(tx *gorm.DB) error) error
	DropTable() error
}

type DBConfig struct {
	User     string
	Password string
	DBName   string
	Host     string
	Port     uint64
}

type MintStatus string

const (
	Pending   MintStatus = "pending"
	Completed MintStatus = "completed"
	Failed    MintStatus = "failed"
)

type MintRequest struct {
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	Contract    string    `gorm:"primaryKey; not null; default: null"`
	Counter     uint64    `gorm:"primaryKey; not null"`
	URIs        string    `gorm:"not null;default: null"`
	ToAddrs     string    `gorm:"not null;default: null"`
	Signature   string    `gorm:"not null;default: null"`
	Sender      string    `gorm:"not null;default: null"`
	TxHash      *string
	Status      MintStatus `gorm:"not null;default: null"`
	Error       *string
	BlockNumber uint64
	Cost        uint64
}
