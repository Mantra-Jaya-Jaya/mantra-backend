package models

func (FraudStatus) TableName() string {
	return "fraud_status"
}

type FraudStatus struct {
	IdFraudStatus uint   `gorm:"primaryKey;column:id" json:"id_fraud_status"`
	NamaStatus    string `gorm:"column:nama_status;unique" json:"nama_status"`
}
