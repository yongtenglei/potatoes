package dao

import "log"

func LoadPotatoes() ([]Potato, error) {
	var potatoes []Potato

	if err := db.Order("type desc").Order("id").Find(&potatoes).Error; err != nil {
		return nil, err
	}

	return potatoes, nil
}

func AddEntry(entry string, potatoType PotatoType) error {
	potato := &Potato{Entry: entry, Type: potatoType}

	if err := db.Create(potato).Error; err != nil {
		return err
	}

	return nil
}

func ToggleCheck(id uint) {
	var potato Potato

	if err := db.First(&potato, id).Error; err != nil {
		log.Println(err)
		return
	}

	if err := db.Model(&potato).Update("checked", !potato.Checked).Error; err != nil {
		log.Println(err)
	}
}

func DeleteEntry(id uint) {
	if err := db.Delete(&Potato{}, id).Error; err != nil {
		log.Println(err)
	}
}
