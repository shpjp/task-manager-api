package config

import (
	"log"

	"task-manager-api/internal/models"

	"gorm.io/gorm"
)

// MigrateAttachments upgrades the attachments table from local-disk storage to
// Cloudinary. Existing rows without a URL are removed because those files lived
// on ephemeral Render disk and cannot be recovered.
func MigrateAttachments(db *gorm.DB) error {
	m := db.Migrator()

	if !m.HasColumn(&models.Attachment{}, "url") {
		if err := db.Exec(`ALTER TABLE attachments ADD COLUMN url varchar(512)`).Error; err != nil {
			return err
		}
	}
	if !m.HasColumn(&models.Attachment{}, "cloudinary_public_id") {
		if err := db.Exec(`ALTER TABLE attachments ADD COLUMN cloudinary_public_id varchar(255)`).Error; err != nil {
			return err
		}
	}
	if !m.HasColumn(&models.Attachment{}, "resource_type") {
		if err := db.Exec(`ALTER TABLE attachments ADD COLUMN resource_type varchar(32) DEFAULT 'image'`).Error; err != nil {
			return err
		}
	}

	result := db.Exec(`DELETE FROM attachments WHERE url IS NULL OR url = ''`)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected > 0 {
		log.Printf("Removed %d legacy attachment row(s) without Cloudinary URLs", result.RowsAffected)
	}

	if err := db.Exec(`UPDATE attachments SET resource_type = 'image' WHERE resource_type IS NULL OR resource_type = ''`).Error; err != nil {
		return err
	}

	if m.HasColumn(&models.Attachment{}, "stored_name") {
		if err := m.DropColumn(&models.Attachment{}, "stored_name"); err != nil {
			return err
		}
	}

	return db.AutoMigrate(&models.Attachment{})
}
