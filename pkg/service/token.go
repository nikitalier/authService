package service

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid"
	"golang.org/x/crypto/bcrypt"
)

var secretToken = []byte("secret")

func (s *Service) GenerateTokenPair(guid string) (map[string]string, error) {
	uuid := uuid.NewV4().String()

	token := jwt.New(jwt.SigningMethodHS512)
	claims := token.Claims.(jwt.MapClaims)
	claims["type"] = "access"
	claims["guid"] = guid
	claims["uuid"] = uuid
	claims["exp"] = time.Now().Add(time.Minute * 60).Unix()

	t, err := token.SignedString(secretToken)
	if err != nil {
		s.logger.Error().Msg(err.Error())
		return nil, err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS512)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["type"] = "refresh"
	rtClaims["guid"] = guid
	rtClaims["uuid"] = uuid
	rtClaims["exp"] = time.Now().Add(time.Hour * 48).Unix()

	rt, err := refreshToken.SignedString(secretToken)
	if err != nil {
		s.logger.Error().Msg(err.Error())
		return nil, err
	}

	tokens := map[string]string{
		"access_token":  t,
		"refresh_token": rt,
	}

	hashedRT, err := bcrypt.GenerateFromPassword([]byte(rt), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error().Msg(err.Error())
		return nil, err
	}

	s.repository.AddRefreshToken(hashedRT, uuid, guid)

	return tokens, nil
}

func (s *Service) RefreshToken(tokenString string) (map[string]string, error) {
	token, err := s.parseToken(tokenString)
	if err != nil {
		s.logger.Error().Msg(err.Error())
		return nil, err
	}

	claims := token.Claims.(jwt.MapClaims)
	uuid := fmt.Sprintf("%v", claims["uuid"])
	guid := fmt.Sprintf("%v", claims["guid"])
	tokenType := fmt.Sprintf("%v", claims["type"])

	if !token.Valid && tokenType != "refresh" {
		s.logger.Error().Msg("Invalid token")
		return nil, fmt.Errorf("Invalid token")
	}

	rt, err := s.repository.FindRefreshTokenByUUID(uuid)
	if err != nil {
		return nil, fmt.Errorf("Invalid token")
	}
	err = bcrypt.CompareHashAndPassword(rt.TokenHash, []byte(tokenString))
	if err != nil {
		s.logger.Error().Msg("Invalid token")
		return nil, fmt.Errorf("Invalid token")
	}

	err = s.repository.DeleteRefreshTokenByUUID(uuid)
	if err != nil {
		return nil, err
	}
	tokens, err := s.GenerateTokenPair(guid)
	if err != nil {
		s.logger.Error().Msg(err.Error())
		return nil, err
	}
	return tokens, err
}

func (s *Service) DeleteRefreshToken(tokenString string) error {
	token, err := s.parseToken(tokenString)
	if err != nil {
		s.logger.Error().Msg(err.Error())
		return err
	}

	claims := token.Claims.(jwt.MapClaims)
	if !token.Valid && fmt.Sprintf("%v", claims["type"]) != "refresh" {
		s.logger.Error().Msg("Invalid token")
		return fmt.Errorf("Invalid token")
	}

	err = s.repository.DeleteRefreshTokenByUUID(fmt.Sprintf("%v", claims["uuid"]))
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) DeleteAllRefreshTokens(guid string) error {
	err := s.repository.DeleteAllRefreshTokensByGUID(guid)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) parseToken(tokenString string) (token *jwt.Token, err error) {
	token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return secretToken, nil
	})
	return token, err
}
