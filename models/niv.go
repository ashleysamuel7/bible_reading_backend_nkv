package models

type NIV struct {
	BookID  int    `gorm:"column:book_id;primaryKey"`
	Book    string `gorm:"column:book;size:255;not null"`
	Chapter int    `gorm:"column:chapter;primaryKey"`
	Verse   int    `gorm:"column:verse;primaryKey"`
	Text    string `gorm:"column:text;size:1000;not null"`
}

// TableName overrides the default pluralized table name
func (NIV) TableName() string {
	return "niv"
}
