package services

import (
	"errors"
	"fmt"

	"github.com/Davethompson01/rialo_hub_backend/config"
	auth "github.com/Davethompson01/rialo_hub_backend/internal/Auth"
	repository "github.com/Davethompson01/rialo_hub_backend/internal/Repository"
	"github.com/Davethompson01/rialo_hub_backend/internal/models"
	"github.com/Davethompson01/rialo_hub_backend/internal/validation"
)

func Register(api *config.ApiConfig, register models.Register) (string, error) {

	if err := validation.ValidateRegister(register); err != nil {
		return "", err
	}

	if repository.CheckDiscordExist(api, register.DiscordUserName) {
		return "", errors.New("Discord Username already exists")
	}

	hashedPassword, err := auth.HashPassword(register.Password)
	if err != nil {
		return err.Error(), err
	}
	register.Password = hashedPassword
	if err := repository.CreateUser(api, register); err != nil {
		return err.Error(), err
	}

	return "Rialo Account successfully created", nil
}

func LoginInto_AsStudent(apicfg *config.ApiConfig, login models.Login) (models.LoginTokens, error) {

	err := validation.ValidateLogin(login)
	if err != nil {
		return models.LoginTokens{}, fmt.Errorf("Invalid Credentials")
	}
	checkMailExist, err := repository.GetUserByUsername(apicfg, login.Username)
	if err != nil {
		return models.LoginTokens{}, err
	}
	// fmt.Println("%v", checkMailExist)

	comparePassword := auth.ComparePassword(checkMailExist.Password, login.Password)
	if comparePassword != nil {

		return models.LoginTokens{}, fmt.Errorf("Incorrect password")
	}

	generateToken, err := auth.GenerateToken(checkMailExist.User_id, checkMailExist.Role)
	if err != nil {
		return models.LoginTokens{}, err
	}

	refreshToken, err := auth.RefreshToken(checkMailExist.User_id, checkMailExist.Role)
	if err != nil {
		return models.LoginTokens{}, err
	}

	fmt.Println("%w", refreshToken)

	return models.LoginTokens{
		AccessToken:  generateToken,
		RefreshToken: refreshToken,
	}, nil

}
