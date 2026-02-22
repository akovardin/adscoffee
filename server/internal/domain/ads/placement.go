package ads

import "gorm.io/gorm"

// Входнаяточка для запроса. На плейсменте может
// быть настроена медиаця или он может использоваться
// как ендпоинт для RTB
type Placement struct {
	gorm.Model

	Units []Unit `gorm:"many2many:placement_units;"`
}
