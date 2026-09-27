package service

import (
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"manadgasht/internal/models"
	"manadgasht/internal/repository"
)

var ErrInvalidCredentials = errors.New("ایمیل یا رمز عبور اشتباه است")
var ErrReferralCodeNotFound = errors.New("کد معرف واردشده معتبر نیست")

const sessionDuration = 30 * 24 * time.Hour // ۳۰ روز

type RegisterUserInput struct {
	FullName     string
	Email        string
	Phone        string
	Password     string
	ReferralCode string
}

type RegisterAgentInput struct {
	FullName   string
	NationalID string
	Email      string
	Phone      string
	CardNumber string
	Password   string
}

type AuthService struct {
	users    repository.UserRepository
	agents   repository.AgentRepository
	admins   repository.AdminRepository
	sessions repository.SessionRepository
}

func NewAuthService(users repository.UserRepository, agents repository.AgentRepository, admins repository.AdminRepository, sessions repository.SessionRepository) *AuthService {
	return &AuthService{users: users, agents: agents, admins: admins, sessions: sessions}
}

// ---------- ثبت‌نام ----------

func (s *AuthService) RegisterUser(in RegisterUserInput) (*models.User, error) {
	verr := newValidationError()
	in.FullName = strings.TrimSpace(in.FullName)
	in.Email = strings.TrimSpace(strings.ToLower(in.Email))
	in.Phone = strings.TrimSpace(in.Phone)
	in.ReferralCode = strings.TrimSpace(strings.ToUpper(in.ReferralCode))

	if in.FullName == "" {
		verr.add("FullName", "نام و نام‌خانوادگی را وارد کنید")
	}
	if !strings.Contains(in.Email, "@") {
		verr.add("Email", "ایمیل معتبر نیست")
	}
	if in.Phone == "" {
		verr.add("Phone", "شماره تماس را وارد کنید")
	}
	if len(in.Password) < 6 {
		verr.add("Password", "رمز عبور باید حداقل ۶ کاراکتر باشد")
	}
	if verr.hasErrors() {
		return nil, verr
	}

	var referredBy *uint
	if in.ReferralCode != "" {
		agent, err := s.agents.FindByReferralCode(in.ReferralCode)
		if err != nil {
			verr.add("ReferralCode", "کد معرف واردشده معتبر نیست")
			return nil, verr
		}
		referredBy = &agent.ID
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		FullName:          in.FullName,
		Email:             in.Email,
		Phone:             in.Phone,
		PasswordHash:      string(hash),
		ReferredByAgentID: referredBy,
	}

	if err := s.users.Create(user); err != nil {
		if isDuplicateKeyError(err) {
			verr.add("Email", "این ایمیل یا شماره تماس قبلاً ثبت شده است")
			return nil, verr
		}
		return nil, err
	}

	return user, nil
}

func (s *AuthService) RegisterAgent(in RegisterAgentInput) (*models.Agent, error) {
	verr := newValidationError()
	in.FullName = strings.TrimSpace(in.FullName)
	in.NationalID = strings.TrimSpace(in.NationalID)
	in.Email = strings.TrimSpace(strings.ToLower(in.Email))
	in.Phone = strings.TrimSpace(in.Phone)
	in.CardNumber = strings.TrimSpace(in.CardNumber)

	if in.FullName == "" {
		verr.add("FullName", "نام و نام‌خانوادگی را وارد کنید")
	}
	if len(in.NationalID) < 8 {
		verr.add("NationalID", "کد ملی معتبر نیست")
	}
	if !strings.Contains(in.Email, "@") {
		verr.add("Email", "ایمیل معتبر نیست")
	}
	if in.Phone == "" {
		verr.add("Phone", "شماره تماس را وارد کنید")
	}
	if in.CardNumber == "" {
		verr.add("CardNumber", "شماره کارت را وارد کنید")
	}
	if len(in.Password) < 6 {
		verr.add("Password", "رمز عبور باید حداقل ۶ کاراکتر باشد")
	}
	if verr.hasErrors() {
		return nil, verr
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	var code string
	for i := 0; i < 5; i++ {
		candidate, err := generateReferralCode()
		if err != nil {
			return nil, err
		}
		exists, err := s.agents.ExistsByReferralCode(candidate)
		if err != nil {
			return nil, err
		}
		if !exists {
			code = candidate
			break
		}
	}
	if code == "" {
		return nil, errors.New("ساخت کد معرف یکتا ناموفق بود، دوباره تلاش کنید")
	}

	agent := &models.Agent{
		FullName:     in.FullName,
		NationalID:   in.NationalID,
		Email:        in.Email,
		Phone:        in.Phone,
		CardNumber:   in.CardNumber,
		PasswordHash: string(hash),
		ReferralCode: code,
	}

	if err := s.agents.Create(agent); err != nil {
		if isDuplicateKeyError(err) {
			verr.add("Email", "این ایمیل، کد ملی یا شماره تماس قبلاً ثبت شده است")
			return nil, verr
		}
		return nil, err
	}

	return agent, nil
}

// ---------- ورود / خروج ----------

func (s *AuthService) LoginUser(email, password string) (token string, expiresAt time.Time, err error) {
	user, err := s.users.FindByEmail(strings.TrimSpace(strings.ToLower(email)))
	if err != nil {
		return "", time.Time{}, ErrInvalidCredentials
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return "", time.Time{}, ErrInvalidCredentials
	}
	return s.createSession(models.SubjectUser, user.ID)
}

func (s *AuthService) LoginAgent(email, password string) (token string, expiresAt time.Time, err error) {
	agent, err := s.agents.FindByEmail(strings.TrimSpace(strings.ToLower(email)))
	if err != nil {
		return "", time.Time{}, ErrInvalidCredentials
	}
	if bcrypt.CompareHashAndPassword([]byte(agent.PasswordHash), []byte(password)) != nil {
		return "", time.Time{}, ErrInvalidCredentials
	}
	return s.createSession(models.SubjectAgent, agent.ID)
}

// LoginAdmin برخلاف کاربر/همکار، با نام‌کاربری (نه ایمیل) وارد می‌شود.
func (s *AuthService) LoginAdmin(username, password string) (token string, expiresAt time.Time, err error) {
	admin, err := s.admins.FindByUsername(strings.TrimSpace(username))
	if err != nil {
		return "", time.Time{}, ErrInvalidCredentials
	}
	if bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password)) != nil {
		return "", time.Time{}, ErrInvalidCredentials
	}
	return s.createSession(models.SubjectAdmin, admin.ID)
}

func (s *AuthService) createSession(subjectType models.SubjectType, subjectID uint) (string, time.Time, error) {
	token, err := generateToken(32)
	if err != nil {
		return "", time.Time{}, err
	}
	expiresAt := time.Now().Add(sessionDuration)
	session := &models.Session{
		Token:       token,
		SubjectType: subjectType,
		SubjectID:   subjectID,
		ExpiresAt:   expiresAt,
	}
	if err := s.sessions.Create(session); err != nil {
		return "", time.Time{}, err
	}
	return token, expiresAt, nil
}

func (s *AuthService) Logout(token string) error {
	if token == "" {
		return nil
	}
	return s.sessions.Delete(token)
}

// ---------- گرفتن کاربر/همکار/ادمین فعلی از روی توکن نشست ----------

func (s *AuthService) CurrentUser(token string) (*models.User, error) {
	if token == "" {
		return nil, ErrInvalidCredentials
	}
	session, err := s.sessions.FindValidByToken(token)
	if err != nil || session.SubjectType != models.SubjectUser {
		return nil, ErrInvalidCredentials
	}
	return s.users.FindByID(session.SubjectID)
}

func (s *AuthService) CurrentAgent(token string) (*models.Agent, error) {
	if token == "" {
		return nil, ErrInvalidCredentials
	}
	session, err := s.sessions.FindValidByToken(token)
	if err != nil || session.SubjectType != models.SubjectAgent {
		return nil, ErrInvalidCredentials
	}
	return s.agents.FindByID(session.SubjectID)
}

func (s *AuthService) CurrentAdmin(token string) (*models.Admin, error) {
	if token == "" {
		return nil, ErrInvalidCredentials
	}
	session, err := s.sessions.FindValidByToken(token)
	if err != nil || session.SubjectType != models.SubjectAdmin {
		return nil, ErrInvalidCredentials
	}
	return s.admins.FindByID(session.SubjectID)
}

func isDuplicateKeyError(err error) bool {
	return err != nil && strings.Contains(strings.ToLower(err.Error()), "duplicate key")
}