package database

import (
	"fmt"

	"gorm.io/gorm"
)

func repairLegacyForeignKeyViolations(db *gorm.DB) error {
	hasUsers := db.Migrator().HasTable(&User{})
	hasPools := db.Migrator().HasTable(&Pool{})
	hasUploads := db.Migrator().HasTable(&Upload{})
	hasSessions := db.Migrator().HasTable(&Session{})
	hasVerifications := db.Migrator().HasTable(&EmailVerification{})

	return db.Transaction(func(tx *gorm.DB) error {
		if hasUploads && hasUsers && hasPools {
			var orphanUploads int64
			if err := tx.Raw(`
				SELECT COUNT(*)
				FROM uploads AS upload
				LEFT JOIN users AS user ON user.id = upload.user_id
				LEFT JOIN pools AS pool ON pool.id = upload.pool_id
				WHERE user.id IS NULL OR pool.id IS NULL
			`).Scan(&orphanUploads).Error; err != nil {
				return err
			}
			if orphanUploads > 0 {
				return fmt.Errorf("foreign-key repair blocked: found %d orphaned uploads requiring manual storage reconciliation", orphanUploads)
			}
		}

		if hasPools && hasUsers {
			if hasUploads {
				var orphanPoolsWithUploads int64
				if err := tx.Raw(`
					SELECT COUNT(*)
					FROM pools AS pool
					LEFT JOIN users AS user ON user.id = pool.user_id
					WHERE user.id IS NULL
					  AND EXISTS (SELECT 1 FROM uploads AS upload WHERE upload.pool_id = pool.id)
				`).Scan(&orphanPoolsWithUploads).Error; err != nil {
					return err
				}
				if orphanPoolsWithUploads > 0 {
					return fmt.Errorf("foreign-key repair blocked: found %d orphaned pools containing uploads", orphanPoolsWithUploads)
				}
			}

			if err := tx.Exec(`
				DELETE FROM pools
				WHERE NOT EXISTS (SELECT 1 FROM users WHERE users.id = pools.user_id)
			`).Error; err != nil {
				return err
			}
		}

		if hasSessions && hasUsers {
			if err := tx.Exec(`
				DELETE FROM sessions
				WHERE NOT EXISTS (SELECT 1 FROM users WHERE users.id = sessions.user_id)
			`).Error; err != nil {
				return err
			}
		}

		if hasVerifications && hasUsers {
			if err := tx.Exec(`
				UPDATE email_verifications
				SET user_id = NULL
				WHERE user_id IS NOT NULL
				  AND NOT EXISTS (SELECT 1 FROM users WHERE users.id = email_verifications.user_id)
			`).Error; err != nil {
				return err
			}
		}

		if hasUsers && hasPools {
			if err := tx.Exec(`
				UPDATE users
				SET default_pool_id = NULL
				WHERE default_pool_id IS NOT NULL
				  AND NOT EXISTS (SELECT 1 FROM pools WHERE pools.id = users.default_pool_id)
			`).Error; err != nil {
				return err
			}
		}

		var remainingViolations int64
		if err := tx.Raw("SELECT COUNT(*) FROM pragma_foreign_key_check").Scan(&remainingViolations).Error; err != nil {
			return err
		}
		if remainingViolations > 0 {
			return fmt.Errorf("foreign-key repair incomplete: %d violations remain", remainingViolations)
		}

		return nil
	})
}

func validateUniqueUploadIdentifiers(db *gorm.DB) error {
	if !db.Migrator().HasTable(&Upload{}) {
		return nil
	}

	var duplicateGroups int64
	if err := db.Raw(`
		SELECT COUNT(*)
		FROM (
			SELECT identifier
			FROM uploads
			GROUP BY identifier
			HAVING COUNT(*) > 1
		)
	`).Scan(&duplicateGroups).Error; err != nil {
		return err
	}
	if duplicateGroups > 0 {
		return fmt.Errorf("cannot create unique upload identifier index: found %d duplicate identifier groups", duplicateGroups)
	}

	return nil
}

func dropObsoleteIndexes(db *gorm.DB) error {
	indexes := []struct {
		model any
		name  string
	}{
		{model: &Upload{}, name: "idx_uploads_identifier"},
		{model: &InvitationKey{}, name: "idx_invitation_keys_key"},
		{model: &EmailVerification{}, name: "idx_email_verifications_key"},
	}

	for _, index := range indexes {
		if db.Migrator().HasIndex(index.model, index.name) {
			if err := db.Migrator().DropIndex(index.model, index.name); err != nil {
				return fmt.Errorf("drop obsolete index %s: %w", index.name, err)
			}
		}
	}

	return nil
}
