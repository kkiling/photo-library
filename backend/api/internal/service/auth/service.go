package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/kkiling/photo-library/backend/api/internal/ctxutils"
	"github.com/kkiling/photo-library/backend/api/internal/service/form"
	"github.com/kkiling/photo-library/backend/api/internal/service/utils"
	"sync/atomic"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/kkiling/photo-library/backend/api/internal/service"
	"github.com/kkiling/photo-library/backend/api/internal/service/auth/session_manager"
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
)

type Config struct {
	AllowRegistration bool `yaml:"allow_registration"`
}

type Storage interface {
	service.Transactor
	GetPeopleCount(ctx context.Context) (int64, error)
	SavePerson(ctx context.Context, person model.Person) error
	SavePersonAuth(ctx context.Context, auth model.Auth) error
	EmailExists(ctx context.Context, email string) (bool, error)
	GetAuth(ctx context.Context, personID uuid.UUID) (model.Auth, error)
	GetPerson(ctx context.Context, personID uuid.UUID) (model.Person, error)
	UpdatePerson(ctx context.Context, personID uuid.UUID, updatePerson model.UpdatePerson) error
	UpdatePersonAuth(ctx context.Context, personID uuid.UUID, updateAuth model.UpdateAuth) error
	GetAuthByEmail(ctx context.Context, email string) (model.Auth, error)
	GetLastActiveRefreshToken(ctx context.Context, refreshTokenID uuid.UUID) (model.RefreshToken, error)
	SaveRefreshToken(ctx context.Context, refreshToken model.RefreshToken) error
	UpdateRefreshTokenStatus(ctx context.Context, refreshTokenID uuid.UUID, status model.RefreshTokenStatus) error
}

// ConfirmCodeService сервис кодов подтверждения
type ConfirmCodeService interface {
	GetActiveConfirmCode(ctx context.Context, code string, confirmType model.ConfirmCodeType) (model.ConfirmCode, error)
	SendConfirmCode(ctx context.Context, personID uuid.UUID, confirmType model.ConfirmCodeType) error
	DeactivateCode(ctx context.Context, personID uuid.UUID, confirmType model.ConfirmCodeType) error
}

type SessionManagerService interface {
	CreateTokenBySession(session model.Session) (session_manager.Token, error)
	CreateTokenByRefresh(refresh model.RefreshSession) (session_manager.Token, error)
	GetRefreshSessionByToken(token string) (model.RefreshSession, error)
}

// PasswordService сервис для работы с паролями
type PasswordService interface {
	HashPassword(password string) ([]byte, error)
	CompareHashAndPassword(hash []byte, password string) error
}

// Service сервис пользователей
type Service struct {
	PersonsExists      atomic.Bool
	cfg                Config
	logger             log.Logger
	storage            Storage
	validate           *validator.Validate
	confirmCodeService ConfirmCodeService
	passwordService    PasswordService
	sessionService     SessionManagerService
}

func NewService(logger log.Logger,
	storage Storage,
	cfg Config,
	confirmCodeService ConfirmCodeService,
	passwordService PasswordService,
	sessionService SessionManagerService,
) *Service {
	return &Service{
		cfg:                cfg,
		logger:             logger,
		storage:            storage,
		validate:           utils.NewValidator(),
		confirmCodeService: confirmCodeService,
		passwordService:    passwordService,
		sessionService:     sessionService,
	}
}

// CheckPersonsExists есть ли хоть один пользователь
func (s *Service) CheckPersonsExists(ctx context.Context) (bool, error) {
	if s.PersonsExists.Load() {
		return true, nil
	}
	count, err := s.storage.GetPeopleCount(ctx)
	if err != nil {
		return false, serviceerr.MakeErr(err, "s.storage.GetPeopleCount")
	}
	s.PersonsExists.Store(count > 0)
	return count > 0, nil
}

func (s *Service) AdminInitInvite(ctx context.Context, form form.AdminInitInviteForm) error {
	if err := s.validate.Struct(form); err != nil {
		return serviceerr.InvalidInputErr(err, "error in creating administrator account")
	}

	// Его может отправить только администратор, либо если нет ни одного пользователя
	if exist, err := s.CheckPersonsExists(ctx); err != nil {
		return serviceerr.MakeErr(err, "s.CheckAdminExist")
	} else if exist {
		return serviceerr.Conflictf("people already exist")
	}

	return s.sendInvite(ctx, form.Email, model.AuthRoleAdmin)
}

func (s *Service) SendInvite(ctx context.Context, form form.SendInviteForm) error {
	if _, err := utils.GetSession(ctx, model.AuthRoleAdmin); err != nil {
		return err
	}
	if err := s.validate.Struct(form); err != nil {
		return serviceerr.InvalidInputErr(err, "error in creating administrator account")
	}
	// TODO: Нужно указать кто отправил invite (session.PersonID)
	return s.sendInvite(ctx, form.Email, form.Role)
}

// ActivateInvite активация инвайта
func (s *Service) ActivateInvite(ctx context.Context, form form.ActivateInviteForm) error {
	_, err := ctxutils.Get[model.Session](ctx, ctxutils.Session)
	switch {
	case err == nil:
		return serviceerr.Conflictf("person is already authorized")
	case errors.Is(err, ctxutils.ErrNotFound):
	default:
		return serviceerr.MakeErr(err, "ctxutils.Get")
	}

	if err = s.validate.Struct(form); err != nil {
		return serviceerr.InvalidInputErr(err, "error in activation account")
	}

	// Поиск кода подтверждения в базе
	code, err := s.confirmCodeService.GetActiveConfirmCode(ctx, form.CodeConfirm, model.ConfirmCodeTypeActivateInvite)
	switch {
	case err == nil:
	case errors.Is(err, serviceerr.ErrNotFound):
		return serviceerr.NotFoundf("confirm code not found")
	default:
		return serviceerr.MakeErr(err, "s.confirmCodeService.GetActiveConfirmCode")
	}

	auth, err := s.storage.GetAuth(ctx, code.PersonID)
	switch {
	case err == nil:
	case errors.Is(err, serviceerr.ErrNotFound):
		return serviceerr.NotFoundf("Person code not found")
	default:
		return serviceerr.MakeErr(err, "s.confirmCodeService.GetPersonFull")
	}

	if auth.Status == model.AuthStatusActivated {
		return serviceerr.Conflictf("person already activated")
	}

	if auth.Status == model.AuthStatusBlocked {
		return serviceerr.PermissionDeniedf("person blocked")
	}

	// Генерация соли
	hash, err := s.passwordService.HashPassword(form.Password)
	if err != nil {
		return serviceerr.MakeErr(err, "s.passwordService.HashPassword")
	}

	form.FirstName = utils.TransformToName(form.FirstName)
	form.Surname = utils.TransformToName(form.Surname)
	form.Patronymic = utils.TransformToNamePtr(form.Patronymic)

	updatePerson := model.UpdatePerson{
		BaseUpdate: model.NewBaseUpdate(),
		FirstName:  model.NewUpdateField(form.FirstName),
		Surname:    model.NewUpdateField(form.Surname),
		Patronymic: model.NewUpdateField(form.Patronymic),
	}

	updateAuth := model.UpdateAuth{
		BaseUpdate:   model.NewBaseUpdate(),
		PasswordHash: model.NewUpdateField(hash),
		Status:       model.NewUpdateField(model.AuthStatusActivated),
	}

	err = s.storage.RunTransaction(ctx, func(ctxTx context.Context) error {
		if err = s.storage.UpdatePerson(ctxTx, code.PersonID, updatePerson); err != nil {
			return serviceerr.MakeErr(err, "s.storage.UpdatePerson")
		}

		if err = s.storage.UpdatePersonAuth(ctxTx, code.PersonID, updateAuth); err != nil {
			return serviceerr.MakeErr(err, "s.storage.UpdatePersonAuth")
		}

		if err = s.confirmCodeService.DeactivateCode(ctxTx, code.PersonID, code.Type); err != nil {
			return serviceerr.MakeErr(err, "s.confirmCodeService.DeactivateCode")
		}
		return nil
	})

	if err != nil {
		return serviceerr.MakeErr(err, "s.storage.UpdatePersonAuth")
	}

	return nil
}

func (s *Service) Login(ctx context.Context, form form.LoginForm) (model.AuthDataDTO, error) {
	if err := s.validate.Struct(form); err != nil {
		return model.AuthDataDTO{}, serviceerr.InvalidInputErr(err, "Error in login data")
	}

	personAuth, err := s.storage.GetAuthByEmail(ctx, form.Email)
	switch {
	case err == nil:
	case errors.Is(err, serviceerr.ErrNotFound):
		return model.AuthDataDTO{}, serviceerr.PermissionDeniedf("Incorrect username or password")
	default:
		return model.AuthDataDTO{}, serviceerr.MakeErr(err, "s.storage.GetPersonAuthByEmail")
	}

	if compareErr := s.passwordService.CompareHashAndPassword(personAuth.PasswordHash, form.Password); compareErr != nil {
		return model.AuthDataDTO{}, serviceerr.PermissionDeniedf("Incorrect username or password")
	}

	if personAuth.Status == model.AuthStatusBlocked {
		return model.AuthDataDTO{}, serviceerr.PermissionDeniedf("Person is blocked")
	}
	if personAuth.Status == model.AuthStatusNotActivated {
		return model.AuthDataDTO{}, serviceerr.PermissionDeniedf("Not activated person account")
	}

	return s.createAuthData(ctx, personAuth)
}

func (s *Service) Registration(ctx context.Context, form form.RegisterForm) error {
	if !s.cfg.AllowRegistration {
		return serviceerr.FailPreconditionf("registration is not available")
	}

	// Его может отправить только администратор, либо если нет ни одного пользователя
	if exist, err := s.CheckPersonsExists(ctx); err != nil {
		return serviceerr.MakeErr(err, "s.CheckAdminExist")
	} else if !exist {
		return serviceerr.FailPreconditionf("first create an administrator")
	}

	if err := s.validate.Struct(form); err != nil {
		return serviceerr.InvalidInputErr(err, "error in creating administrator account")
	}

	if emailExists, err := s.storage.EmailExists(ctx, form.Email); err != nil {
		return serviceerr.MakeErr(err, "s.storage.EmailExists")
	} else if emailExists {
		return serviceerr.Conflictf("email already exists")
	}

	// Генерация соли
	hash, err := s.passwordService.HashPassword(form.Password)
	if err != nil {
		return serviceerr.MakeErr(err, "s.passwordService.HashPassword")
	}

	newPerson := model.Person{
		Base: model.NewBase(),
		ID:   uuid.New(),
	}

	newAuth := model.Auth{
		Base:         model.NewBase(),
		PersonID:     newPerson.ID,
		Email:        form.Email,
		PasswordHash: hash,
		Status:       model.AuthStatusNotActivated,
		Role:         model.AuthRoleUser,
	}

	err = s.storage.RunTransaction(ctx, func(ctxTx context.Context) error {
		if saveErr := s.storage.SavePerson(ctxTx, newPerson); saveErr != nil {
			return fmt.Errorf("s.storage.CreatePerson: %w", saveErr)
		}
		if saveErr := s.storage.SavePersonAuth(ctxTx, newAuth); saveErr != nil {
			return fmt.Errorf("s.storage.AddPersonAuth: %w", saveErr)
		}
		return nil
	})

	if err != nil {
		return serviceerr.MakeErr(err, " s.storage.RunTransaction")
	}

	err = s.confirmCodeService.SendConfirmCode(ctx, newPerson.ID, model.ConfirmCodeTypeActivateRegistration)
	if err != nil {
		return serviceerr.MakeErr(err, "s.confirmCodeService.SendConfirmCode")
	}

	return nil
}

// ActivateRegistration активация инвайта
func (s *Service) ActivateRegistration(ctx context.Context, form form.ActivateRegisterForm) error {
	// TODO: пользовтель должне быть не авторизован
	// TODO: проверка авторизации из контекста

	if err := s.validate.Struct(form); err != nil {
		return serviceerr.InvalidInputErr(err, "error in activation account")
	}

	// Поиск кода подтверждения в базе
	code, err := s.confirmCodeService.GetActiveConfirmCode(ctx, form.CodeConfirm, model.ConfirmCodeTypeActivateRegistration)
	switch {
	case err == nil:
	case errors.Is(err, serviceerr.ErrNotFound):
		return serviceerr.NotFoundf("Confirm code not found")
	default:
		return serviceerr.MakeErr(err, "s.confirmCodeService.GetActiveConfirmCode")
	}

	auth, err := s.storage.GetAuth(ctx, code.PersonID)
	switch {
	case err == nil:
	case errors.Is(err, serviceerr.ErrNotFound):
		return serviceerr.NotFoundf("Person code not found")
	default:
		return serviceerr.MakeErr(err, "s.confirmCodeService.GetPersonFull")
	}

	if auth.Status == model.AuthStatusActivated {
		return serviceerr.Conflictf("person already activated")
	}

	if auth.Status == model.AuthStatusBlocked {
		return serviceerr.PermissionDeniedf("person blocked")
	}

	updateAuth := model.UpdateAuth{
		BaseUpdate: model.NewBaseUpdate(),
		Status:     model.NewUpdateField(model.AuthStatusActivated),
	}

	err = s.storage.RunTransaction(ctx, func(ctxTx context.Context) error {
		if err = s.storage.UpdatePersonAuth(ctxTx, code.PersonID, updateAuth); err != nil {
			return serviceerr.MakeErr(err, "s.storage.UpdatePersonAuth")
		}

		if err = s.confirmCodeService.DeactivateCode(ctxTx, code.PersonID, code.Type); err != nil {
			return serviceerr.MakeErr(err, "s.confirmCodeService.DeactivateCode")
		}
		return nil
	})

	if err != nil {
		return serviceerr.MakeErr(err, "s.storage.UpdatePersonAuth")
	}

	return nil
}

func (s *Service) EmailAvailable(ctx context.Context, form form.EmailAvailableForm) (bool, error) {
	if err := s.validate.Struct(form); err != nil {
		return false, serviceerr.InvalidInputErr(err, "error in creating administrator account")
	}

	if emailExists, err := s.storage.EmailExists(ctx, form.Email); err != nil {
		return false, serviceerr.MakeErr(err, "s.storage.EmailExists")
	} else {
		return emailExists, nil
	}
}

func (s *Service) Logout(ctx context.Context, token string) error {
	refreshSession, err := s.sessionService.GetRefreshSessionByToken(token)
	if err != nil {
		return serviceerr.PermissionDeniedErr(fmt.Errorf("invalid token"))
	}
	err = s.storage.UpdateRefreshTokenStatus(ctx, refreshSession.RefreshTokenID, model.RefreshTokenStatusLogout)
	switch {
	case err == nil:
		return nil
	case errors.Is(err, serviceerr.ErrNotFound):
		return serviceerr.PermissionDeniedErr(fmt.Errorf("invalid token"))
	default:
		return serviceerr.MakeErr(err, "s.storage.UpdateRefreshTokenStatus")
	}
}

func (s *Service) RefreshToken(ctx context.Context, token string) (model.AuthDataDTO, error) {
	refreshSession, err := s.sessionService.GetRefreshSessionByToken(token)
	if err != nil {
		return model.AuthDataDTO{}, serviceerr.PermissionDeniedErr(fmt.Errorf("invalid token"))
	}
	refreshToken, err := s.storage.GetLastActiveRefreshToken(ctx, refreshSession.RefreshTokenID)
	switch {
	case err == nil:
	case errors.Is(err, serviceerr.ErrNotFound):
		return model.AuthDataDTO{}, serviceerr.PermissionDeniedErr(fmt.Errorf("invalid token"))
	default:
		return model.AuthDataDTO{}, serviceerr.MakeErr(err, "s.storage.GetLastActiveRefreshToken")
	}

	personAuth, err := s.storage.GetAuth(ctx, refreshToken.PersonID)
	switch {
	case err == nil:
	case errors.Is(err, serviceerr.ErrNotFound):
		return model.AuthDataDTO{}, serviceerr.PermissionDeniedf("Incorrect username or password")
	default:
		return model.AuthDataDTO{}, serviceerr.MakeErr(err, "s.storage.GetPersonAuthByEmail")
	}

	if personAuth.Status == model.AuthStatusBlocked {
		return model.AuthDataDTO{}, serviceerr.PermissionDeniedf("Person is blocked")
	}

	if personAuth.Status == model.AuthStatusNotActivated {
		return model.AuthDataDTO{}, serviceerr.PermissionDeniedf("Not activated person account")
	}

	err = s.storage.UpdateRefreshTokenStatus(ctx, refreshSession.RefreshTokenID, model.RefreshTokenStatusRevoked)
	switch {
	case err == nil:
	case errors.Is(err, serviceerr.ErrNotFound):
		return model.AuthDataDTO{}, serviceerr.PermissionDeniedErr(fmt.Errorf("invalid token"))
	default:
		return model.AuthDataDTO{}, serviceerr.MakeErr(err, "s.storage.UpdateRefreshTokenStatus")
	}

	return s.createAuthData(ctx, personAuth)
}
