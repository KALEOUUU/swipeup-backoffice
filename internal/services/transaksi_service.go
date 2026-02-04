package services

import (
	"swipeup-be/internal/models"
	"time"

	"gorm.io/gorm"
)

type TransaksiService struct {
	*BaseService[models.Transaksi]
}

func NewTransaksiService(db *gorm.DB) *TransaksiService {
	return &TransaksiService{
		BaseService: NewBaseService[models.Transaksi](db),
	}
}

func (s *TransaksiService) CreateWithDetails(transaksi *models.Transaksi, details []models.DetailTransaksi) error {
	return s.GetDB().Transaction(func(tx *gorm.DB) error {
		transaksi.Tanggal = time.Now()
		if err := tx.Create(transaksi).Error; err != nil {
			return err
		}

		for i := range details {
			details[i].IDTransaksi = transaksi.ID
			if err := tx.Create(&details[i]).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *TransaksiService) GetBySiswaID(siswaID uint) ([]models.Transaksi, error) {
	return s.FindWithCondition(map[string]interface{}{"id_siswa": siswaID}, "Stan", "Siswa", "DetailTransaksi", "DetailTransaksi.Menu")
}

func (s *TransaksiService) GetByStanID(stanID uint) ([]models.Transaksi, error) {
	return s.FindWithCondition(map[string]interface{}{"id_stan": stanID}, "Stan", "Siswa", "DetailTransaksi", "DetailTransaksi.Menu")
}

func (s *TransaksiService) GetByStatus(status models.StatusTransaksi) ([]models.Transaksi, error) {
	return s.FindWithCondition(map[string]interface{}{"status": status}, "Stan", "Siswa", "DetailTransaksi")
}

func (s *TransaksiService) UpdateStatus(id uint, status models.StatusTransaksi) error {
	return s.GetDB().Model(&models.Transaksi{}).Where("id = ?", id).Update("status", status).Error
}

func (s *TransaksiService) GetWithFullDetails(id uint) (*models.Transaksi, error) {
	return s.FindByID(id, "Stan", "Siswa", "DetailTransaksi", "DetailTransaksi.Menu")
}

func (s *TransaksiService) GetByDateRange(startDate, endDate time.Time) ([]models.Transaksi, error) {
	var transaksi []models.Transaksi
	err := s.GetDB().Preload("Stan").Preload("Siswa").Preload("DetailTransaksi").
		Where("tanggal BETWEEN ? AND ?", startDate, endDate).
		Find(&transaksi).Error
	return transaksi, err
}

// DeleteTransaksi soft deletes a transaction and its details (admin only)
func (s *TransaksiService) DeleteTransaksi(id uint) error {
	return s.GetDB().Transaction(func(tx *gorm.DB) error {
		// Delete detail transaksi first
		if err := tx.Where("id_transaksi = ?", id).Delete(&models.DetailTransaksi{}).Error; err != nil {
			return err
		}
		// Delete transaksi
		return tx.Delete(&models.Transaksi{}, id).Error
	})
}

// UpdateTransaksi updates transaction fields (admin only)
func (s *TransaksiService) UpdateTransaksi(id uint, updates map[string]interface{}) error {
	return s.GetDB().Model(&models.Transaksi{}).Where("id = ?", id).Updates(updates).Error
}
