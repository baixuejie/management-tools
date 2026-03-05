package services

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/baixuejie/key-management-tool/backend/internal/models"
	"gorm.io/gorm"
)

const (
	ChannelXianyu = "xianyu"
	ChannelWechat = "wechat"
)

type CostRecordItem struct {
	ID           uint      `json:"id"`
	Amount       float64   `json:"amount"`
	Note         string    `json:"note"`
	RecordedBy   uint      `json:"recorded_by"`
	RecorderName string    `json:"recorder_name"`
	CreatedAt    time.Time `json:"created_at"`
}

type TransactionItem struct {
	ID               uint      `json:"id"`
	CustomerID       uint      `json:"customer_id"`
	CustomerName     string    `json:"customer_name"`
	Amount           float64   `json:"amount"`
	Channel          string    `json:"channel"`
	CommissionRate   float64   `json:"commission_rate"`
	CommissionAmount float64   `json:"commission_amount"`
	IsNewCustomer    bool      `json:"is_new_customer"`
	RecordedBy       uint      `json:"recorded_by"`
	RecorderName     string    `json:"recorder_name"`
	CreatedAt        time.Time `json:"created_at"`
}

type Statistics struct {
	TotalCost        float64 `json:"total_cost"`
	TotalRevenue     float64 `json:"total_revenue"`
	TotalCommission  float64 `json:"total_commission"`
	NetProfit        float64 `json:"net_profit"`
	NewCustomers     int64   `json:"new_customers"`
	RenewalCustomers int64   `json:"renewal_customers"`
}

type CreateTransactionInput struct {
	CustomerID    *uint
	CustomerName  string
	Amount        float64
	Channel       string
	IsNewCustomer bool
	RecordedBy    uint
}

type LedgerService struct {
	db *gorm.DB
}

func NewLedgerService(db *gorm.DB) *LedgerService {
	return &LedgerService{db: db}
}

func (s *LedgerService) ListCosts(page, limit int) (int64, []CostRecordItem, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit

	var total int64
	if err := s.db.Model(&models.CostRecord{}).Count(&total).Error; err != nil {
		return 0, nil, err
	}

	var items []CostRecordItem
	err := s.db.Table("cost_records AS cr").
		Select("cr.id, cr.amount, cr.note, cr.recorded_by, u.display_name AS recorder_name, cr.created_at").
		Joins("JOIN users AS u ON u.id = cr.recorded_by").
		Order("cr.created_at DESC").
		Limit(limit).
		Offset(offset).
		Scan(&items).Error
	if err != nil {
		return 0, nil, err
	}

	return total, items, nil
}

func (s *LedgerService) CreateCost(amount float64, note string, recordedBy uint) (*models.CostRecord, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}
	if recordedBy == 0 {
		return nil, errors.New("recorded_by is required")
	}

	cost := &models.CostRecord{
		Amount:     amount,
		Note:       strings.TrimSpace(note),
		RecordedBy: recordedBy,
	}
	if err := s.db.Create(cost).Error; err != nil {
		return nil, err
	}
	return cost, nil
}

func (s *LedgerService) DeleteCost(id uint) error {
	result := s.db.Delete(&models.CostRecord{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("cost record not found")
	}
	return nil
}

func (s *LedgerService) ListCustomers(search string) ([]models.Customer, error) {
	query := s.db.Model(&models.Customer{})
	if trimmed := strings.TrimSpace(search); trimmed != "" {
		query = query.Where("name LIKE ?", "%"+trimmed+"%")
	}

	var customers []models.Customer
	if err := query.Order("updated_at DESC").Find(&customers).Error; err != nil {
		return nil, err
	}
	return customers, nil
}

func (s *LedgerService) CreateCustomer(name string, createdBy uint) (*models.Customer, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("name is required")
	}
	if createdBy == 0 {
		return nil, errors.New("created_by is required")
	}

	customer, err := s.findOrCreateCustomer(name, createdBy)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func (s *LedgerService) UpdateCustomer(id uint, name string) (*models.Customer, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("name is required")
	}

	var customer models.Customer
	if err := s.db.First(&customer, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("customer not found")
		}
		return nil, err
	}

	customer.Name = name
	if err := s.db.Save(&customer).Error; err != nil {
		return nil, err
	}

	return &customer, nil
}

func (s *LedgerService) ListTransactions(page, limit int, customerID *uint, isNewCustomer *bool) (int64, []TransactionItem, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit

	countQuery := s.db.Model(&models.Transaction{})
	if customerID != nil {
		countQuery = countQuery.Where("customer_id = ?", *customerID)
	}
	if isNewCustomer != nil {
		countQuery = countQuery.Where("is_new_customer = ?", *isNewCustomer)
	}

	var total int64
	if err := countQuery.Count(&total).Error; err != nil {
		return 0, nil, err
	}

	query := s.db.Table("transactions AS t").
		Select(`t.id, t.customer_id, c.name AS customer_name, t.amount, t.channel,
                t.commission_rate, t.commission_amount, t.is_new_customer,
                t.recorded_by, u.display_name AS recorder_name, t.created_at`).
		Joins("JOIN customers AS c ON c.id = t.customer_id").
		Joins("JOIN users AS u ON u.id = t.recorded_by")

	if customerID != nil {
		query = query.Where("t.customer_id = ?", *customerID)
	}
	if isNewCustomer != nil {
		query = query.Where("t.is_new_customer = ?", *isNewCustomer)
	}

	var items []TransactionItem
	if err := query.Order("t.created_at DESC").Limit(limit).Offset(offset).Scan(&items).Error; err != nil {
		return 0, nil, err
	}

	return total, items, nil
}

func (s *LedgerService) CreateTransaction(input CreateTransactionInput) (*TransactionItem, error) {
	if input.RecordedBy == 0 {
		return nil, errors.New("recorded_by is required")
	}
	if input.Amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}

	rate, err := commissionRate(strings.TrimSpace(strings.ToLower(input.Channel)))
	if err != nil {
		return nil, err
	}

	var txModel models.Transaction
	err = s.db.Transaction(func(tx *gorm.DB) error {
		var customer models.Customer
		if input.IsNewCustomer {
			customer, err = s.findOrCreateCustomerTx(tx, input.CustomerName, input.RecordedBy)
			if err != nil {
				return err
			}
		} else {
			if input.CustomerID == nil || *input.CustomerID == 0 {
				return errors.New("customer_id is required for renewal transactions")
			}
			if err := tx.First(&customer, *input.CustomerID).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return errors.New("customer not found")
				}
				return err
			}
		}

		txModel = models.Transaction{
			CustomerID:       customer.ID,
			Amount:           input.Amount,
			Channel:          strings.ToLower(strings.TrimSpace(input.Channel)),
			CommissionRate:   rate,
			CommissionAmount: roundTo2(input.Amount * rate),
			IsNewCustomer:    input.IsNewCustomer,
			RecordedBy:       input.RecordedBy,
		}
		return tx.Create(&txModel).Error
	})
	if err != nil {
		return nil, err
	}

	return s.GetTransactionByID(txModel.ID)
}

func (s *LedgerService) GetTransactionByID(id uint) (*TransactionItem, error) {
	var item TransactionItem
	err := s.db.Table("transactions AS t").
		Select(`t.id, t.customer_id, c.name AS customer_name, t.amount, t.channel,
                t.commission_rate, t.commission_amount, t.is_new_customer,
                t.recorded_by, u.display_name AS recorder_name, t.created_at`).
		Joins("JOIN customers AS c ON c.id = t.customer_id").
		Joins("JOIN users AS u ON u.id = t.recorded_by").
		Where("t.id = ?", id).
		Scan(&item).Error
	if err != nil {
		return nil, err
	}
	if item.ID == 0 {
		return nil, errors.New("transaction not found")
	}
	return &item, nil
}

func (s *LedgerService) GetStatistics() (*Statistics, error) {
	stats := &Statistics{}

	if err := s.db.Model(&models.CostRecord{}).Select("COALESCE(SUM(amount), 0)").Scan(&stats.TotalCost).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&models.Transaction{}).Select("COALESCE(SUM(amount), 0)").Scan(&stats.TotalRevenue).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&models.Transaction{}).Select("COALESCE(SUM(commission_amount), 0)").Scan(&stats.TotalCommission).Error; err != nil {
		return nil, err
	}

	var newCustomers int64
	if err := s.db.Model(&models.Transaction{}).
		Where("is_new_customer = ?", true).
		Distinct("customer_id").
		Count(&newCustomers).Error; err != nil {
		return nil, err
	}
	stats.NewCustomers = newCustomers

	var renewalCustomers int64
	if err := s.db.Model(&models.Transaction{}).
		Where("is_new_customer = ?", false).
		Distinct("customer_id").
		Count(&renewalCustomers).Error; err != nil {
		return nil, err
	}
	stats.RenewalCustomers = renewalCustomers

	stats.NetProfit = roundTo2(stats.TotalRevenue - stats.TotalCost - stats.TotalCommission)
	return stats, nil
}

func (s *LedgerService) findOrCreateCustomer(name string, createdBy uint) (*models.Customer, error) {
	customer, err := s.findOrCreateCustomerTx(s.db, name, createdBy)
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (s *LedgerService) findOrCreateCustomerTx(tx *gorm.DB, name string, createdBy uint) (models.Customer, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return models.Customer{}, errors.New("customer_name is required for new customer transactions")
	}

	var existing models.Customer
	err := tx.Where("name = ? AND created_by = ?", name, createdBy).First(&existing).Error
	if err == nil {
		return existing, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return models.Customer{}, err
	}

	customer := models.Customer{
		Name:      name,
		CreatedBy: createdBy,
	}
	if err := tx.Create(&customer).Error; err != nil {
		return models.Customer{}, err
	}
	return customer, nil
}

func commissionRate(channel string) (float64, error) {
	switch channel {
	case ChannelXianyu:
		return 0.007, nil
	case ChannelWechat:
		return 0, nil
	default:
		return 0, fmt.Errorf("unsupported channel: %s", channel)
	}
}

func roundTo2(v float64) float64 {
	return math.Round(v*100) / 100
}
