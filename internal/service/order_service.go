package service

import (
	"errors"

	"gorm.io/gorm"

	"manadgasht/internal/models"
	"manadgasht/internal/repository"
)

var ErrPackageNotFound = errors.New("تور/هتل پیدا نشد")

type OrderService struct {
	db          *gorm.DB
	packages    repository.PackageRepository
	orders      repository.OrderRepository
	agents      repository.AgentRepository
	users       repository.UserRepository
	commissions repository.CommissionRepository
}

func NewOrderService(
	db *gorm.DB,
	packages repository.PackageRepository,
	orders repository.OrderRepository,
	agents repository.AgentRepository,
	users repository.UserRepository,
	commissions repository.CommissionRepository,
) *OrderService {
	return &OrderService{
		db:          db,
		packages:    packages,
		orders:      orders,
		agents:      agents,
		users:       users,
		commissions: commissions,
	}
}

// PlaceOrderForUser خرید یک مشتری عادی را ثبت می‌کند. اگر آن کاربر
// موقع ثبت‌نام یک کد معرف معتبر وارد کرده باشد، کمیسیون بر اساس نوع
// پکیج محاسبه و خودکار به موجودی همان همکار اضافه می‌شود.
//
// نکته: پرداخت واقعی (درگاه بانکی) هنوز وصل نیست؛ این متد سفارش را
// مستقیم «تأییدشده» ثبت می‌کند تا جریان کمیسیون از همین حالا قابل
// تست باشد. وصل‌کردن درگاه پرداخت واقعی، گام بعدی است.
func (s *OrderService) PlaceOrderForUser(userID uint, packageSlug string) (*models.Order, error) {
	pkg, err := s.packages.FindBySlug(packageSlug)
	if err != nil {
		return nil, ErrPackageNotFound
	}

	user, err := s.users.FindByID(userID)
	if err != nil {
		return nil, err
	}

	uid := userID
	order := &models.Order{
		UserID:     &uid,
		PackageID:  pkg.ID,
		TotalPrice: pkg.Price,
		Status:     models.OrderStatusConfirmed,
	}

	if user.ReferredByAgentID != nil {
		agent, err := s.agents.FindByID(*user.ReferredByAgentID)
		if err == nil {
			percentage, err := s.commissions.PercentageFor(pkg.Type)
			if err == nil && percentage > 0 {
				order.ReferringAgentID = user.ReferredByAgentID
				order.ReferralCodeUsed = agent.ReferralCode
				order.CommissionPercentage = percentage
				order.CommissionAmount = pkg.Price * percentage / 100
			}
		}
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.orders.Create(tx, order); err != nil {
			return err
		}
		if order.CommissionAmount > 0 && order.ReferringAgentID != nil {
			if err := s.agents.IncreaseBalance(tx, *order.ReferringAgentID, order.CommissionAmount); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return order, nil
}

// PlaceOrderForAgentWallet خرید خودِ همکار از موجودی کیف‌پولش را ثبت
// می‌کند (بدون کمیسیون، چون خریدار خودِ همکاره).
func (s *OrderService) PlaceOrderForAgentWallet(agentID uint, packageSlug string) (*models.Order, error) {
	pkg, err := s.packages.FindBySlug(packageSlug)
	if err != nil {
		return nil, ErrPackageNotFound
	}

	aid := agentID
	order := &models.Order{
		PurchasedByAgentID: &aid,
		PackageID:          pkg.ID,
		TotalPrice:         pkg.Price,
		Status:             models.OrderStatusConfirmed,
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.agents.DecreaseBalance(tx, agentID, pkg.Price); err != nil {
			return err
		}
		return s.orders.Create(tx, order)
	})
	if err != nil {
		if errors.Is(err, repository.ErrInsufficientBalance) {
			return nil, repository.ErrInsufficientBalance
		}
		return nil, err
	}

	return order, nil
}

func (s *OrderService) ListByUser(userID uint) ([]models.Order, error) {
	return s.orders.ListByUser(userID)
}

func (s *OrderService) ListByReferringAgent(agentID uint) ([]models.Order, error) {
	return s.orders.ListByReferringAgent(agentID)
}
