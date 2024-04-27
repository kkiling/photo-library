package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
	"github.com/kkiling/photo-library/backend/api/pkg/common/log"
	"github.com/kkiling/photo-library/backend/api/pkg/common/utils"
	"time"

	"github.com/google/uuid"

	"github.com/kkiling/photo-library/backend/api/internal/service"
	"github.com/kkiling/photo-library/backend/api/internal/service/auth/codes"
)

type Storage interface {
	service.Transactor
	GetPeopleCount(ctx context.Context) (int64, error)
	CreatePerson(ctx context.Context, person Person) error
	AddPersonAuth(ctx context.Context, auth Auth) error
	GetPersonFull(ctx context.Context, personID uuid.UUID) (PersonFull, error)
	GetPerson(ctx context.Context, personID uuid.UUID) (Person, error)
	UpdatePersonAuth(ctx context.Context, personID uuid.UUID, updateAuth UpdateAuth) error
	GetAuthByEmail(ctx context.Context, email string) (Auth, error)
}

// ConfirmCodeService сервис кодов подтверждения
type ConfirmCodeService interface {
	GetActiveConfirmCode(ctx context.Context, code string) (codes.ConfirmCode, error)
	SendConfirmCode(ctx context.Context, personID uuid.UUID, confirmType codes.ConfirmCodeType) error
	DeactivateCode(ctx context.Context, personID uuid.UUID, confirmType codes.ConfirmCodeType) error
}

type SessionManagerService interface {
	CreateTokenBySession(ctx context.Context, session Session) (string, time.Time, error)
}

// PasswordService сервис для работы с паролями
type PasswordService interface {
	HashPassword(password string) ([]byte, error)
	CompareHashAndPassword(hash []byte, password string) error
}

// Service сервис пользователей
type Service struct {
	logger             log.Logger
	storage            Storage
	validate           *validator.Validate
	confirmCodeService ConfirmCodeService
	passwordService    PasswordService
	sessionService     SessionManagerService
}

// CheckAdminExist зарегистрированы ли люди в системе
func (s *Service) CheckAdminExist(ctx context.Context) (adminExist bool, err error) {
	count, err := s.storage.GetPeopleCount(ctx)
	if err != nil {
		return false, serviceerr.MakeErr(err, "s.storage.GetPeopleCount")
	}
	return count > 0, nil
}

func transformForm(form *InitAdminForm) {
	form.Name = utils.TransformToName(form.Name)
	form.Surname = utils.TransformToName(form.Surname)
	form.Patronymic = utils.TransformToNamePtr(form.Patronymic)
}

// InitFirstAdmin инициализация администратора в системе (если еще не зарегистрированы люди)
func (s *Service) InitFirstAdmin(ctx context.Context, form InitAdminForm) (PersonFull, error) {
	if exist, err := s.CheckAdminExist(ctx); err != nil {
		return PersonFull{}, serviceerr.MakeErr(err, "s.CheckAdminExist")
	} else if exist {
		return PersonFull{}, serviceerr.Conflictf("Admin has already been created")
	}

	transformForm(&form)
	if err := s.validate.Struct(form); err != nil {
		return PersonFull{}, serviceerr.InvalidInputErr(err, "Error in creating administrator account")
	}

	newPerson := Person{
		ID:         uuid.New(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Name:       form.Name,
		Surname:    form.Surname,
		Patronymic: form.Patronymic,
	}

	newAuth := Auth{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		PersonID:  newPerson.ID,
		Email:     form.Email,
		Status:    AuthStatusNotActivated,
	}

	err := s.storage.RunTransaction(ctx, func(ctxTx context.Context) error {
		if saveErr := s.storage.CreatePerson(ctxTx, newPerson); saveErr != nil {
			return fmt.Errorf("createPerson: %w", saveErr)
		}
		if saveErr := s.storage.AddPersonAuth(ctxTx, newAuth); saveErr != nil {
			return fmt.Errorf("AddPersonAuth: %w", saveErr)
		}
		return nil
	})

	if err != nil {
		return PersonFull{}, serviceerr.MakeErr(err, " s.storage.RunTransaction")
	}

	err = s.confirmCodeService.SendConfirmCode(ctx, newPerson.ID, codes.ConfirmCodeTypeActivateAuth)
	if err != nil {
		return PersonFull{}, serviceerr.MakeErr(err, "s.confirmCodeService.SendConfirmCode")
	}

	fullPerson, err := s.storage.GetPersonFull(ctx, newPerson.ID)
	if err != nil {
		return PersonFull{}, serviceerr.MakeErr(err, "s.storage.GetPersonFull")
	}

	return fullPerson, nil
}

// ActivateAuth активация авторизации
func (s *Service) ActivateAuth(ctx context.Context, form ActivateAuthForm) error {
	if err := s.validate.Struct(form); err != nil {
		return serviceerr.InvalidInputErr(err, "Error in activation account")
	}

	// Поиск кода подтверждения в базе
	code, err := s.confirmCodeService.GetActiveConfirmCode(ctx, form.CodeConfirm)
	switch {
	case err == nil:
	case errors.Is(err, serviceerr.ErrNotFound):
		return serviceerr.NotFoundf("Confirm code not found")
	default:
		return serviceerr.MakeErr(err, "s.confirmCodeService.GetActiveConfirmCode")
	}

	fullPerson, err := s.storage.GetPersonFull(ctx, code.PersonID)
	switch {
	case err == nil:
	case errors.Is(err, serviceerr.ErrNotFound):
		return serviceerr.NotFoundf("Person code not found")
	default:
		return serviceerr.MakeErr(err, "s.confirmCodeService.GetPersonFull")
	}

	if fullPerson.Auth == nil {
		return serviceerr.Conflictf("Person has no auth role assigned")
	}

	if fullPerson.Auth.Status == AuthStatusActivated {
		return serviceerr.Conflictf("Person already activated")
	}

	if fullPerson.Auth.Status == AuthStatusBlocked {
		return serviceerr.PermissionDeniedErr("Person blocked")
	}

	// Генерация соли
	hash, err := s.passwordService.HashPassword(form.Password)
	if err != nil {
		return serviceerr.MakeErr(err, "s.passwordService.HashPassword")
	}

	updateAuth := UpdateAuth{
		UpdatedAt:    time.Now(),
		PasswordHash: hash,
		Status:       AuthStatusActivated,
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

func (s *Service) Login(ctx context.Context, form LoginForm) (AuthData, error) {
	if err := s.validate.Struct(form); err != nil {
		return AuthData{}, serviceerr.InvalidInputErr(err, "Error in login data")
	}

	personAuth, err := s.storage.GetAuthByEmail(ctx, form.Email)

	switch {
	case err == nil:
	case errors.Is(err, serviceerr.ErrNotFound):
		return AuthData{}, serviceerr.PermissionDeniedErr("Incorrect username or password")
	default:
		return AuthData{}, serviceerr.MakeErr(err, "s.storage.GetPersonAuthByEmail")
	}

	if compareErr := s.passwordService.CompareHashAndPassword(personAuth.PasswordHash, form.Password); compareErr != nil {
		return AuthData{}, serviceerr.PermissionDeniedErr("Incorrect username or password")
	}

	if personAuth.Status == AuthStatusBlocked {
		return AuthData{}, serviceerr.PermissionDeniedErr("Person is blocked")
	}
	if personAuth.Status == AuthStatusNotActivated {
		return AuthData{}, serviceerr.PermissionDeniedErr("Not activated person account")
	}

	// Пока излишне запрашивать person, но возможно потом будет роль (Admin, User ит)
	person, err := s.storage.GetPerson(ctx, personAuth.PersonID)
	if err != nil {
		return AuthData{}, serviceerr.MakeErr(err, "s.storage.GetPersonFull")
	}

	session := Session{
		PersonID: person.ID,
	}

	accessToken, expiresAt, err := s.sessionService.CreateTokenBySession(ctx, session)
	if err != nil {
		return AuthData{}, serviceerr.MakeErr(err, "s.sessionService.CreateTokenBySession")
	}

	return AuthData{
		Email:                  form.Email,
		AccessToken:            accessToken,
		AccessTokenExpiration:  expiresAt,
		RefreshToken:           "", // TODO: добавить поддержкуs
		RefreshTokenExpiration: time.Time{},
	}, nil
}

func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (AuthData, error) {
	panic("implement me")
}
