package service

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid"
)

var secretToken = []byte("secret")

func (s *Service) GenerateTokenPair(GUID string) (map[string]string, error) {
	uuid := uuid.NewV4().String()

	token := jwt.New(jwt.SigningMethodHS512)
	claims := token.Claims.(jwt.MapClaims)
	claims["type"] = "access"
	claims["GUID"] = GUID
	claims["uuid"] = uuid
	claims["exp"] = time.Now().Add(time.Minute * 60).Unix()

	t, err := token.SignedString(secretToken)
	if err != nil {
		s.logger.Error().Msg(err.Error())
		return nil, err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS512)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["exp"] = time.Now().Add(time.Hour * 48).Unix()
	rtClaims["type"] = "refresh"
	rtClaims["uuid"] = uuid
	rtClaims["GUID"] = GUID

	rt, err := refreshToken.SignedString(secretToken)
	if err != nil {
		s.logger.Error().Msg(err.Error())
		return nil, err
	}
	tokens := map[string]string{
		"access_token":  t,
		"refresh_token": rt,
	}

	s.repository.AddTokenPairs(tokens, GUID, uuid)

	return tokens, nil
}

func (s *Service) RefreshToken(tokenString string) (map[string]string, error) {
	token, err := s.parseToken(tokenString)
	if err != nil {
		s.logger.Error().Msg(err.Error())
		return nil, err
	}

	claims := token.Claims.(jwt.MapClaims)
	if token.Valid && fmt.Sprintf("%v", claims["type"]) == "refresh" {
		err = s.repository.DeleteRefreshToken(fmt.Sprintf("%v", claims["uuid"]))
		if err != nil {
			return nil, err
		}
		err = s.repository.DeleteAccessToken(fmt.Sprintf("%v", claims["uuid"]))
		if err != nil {
			return nil, err
		}
		tokens, err := s.GenerateTokenPair(fmt.Sprintf("%v", claims["uuid"]))
		if err != nil {
			s.logger.Error().Msg(err.Error())
			return nil, err
		}
		return tokens, err
	}
	s.logger.Error().Msg("Invalid token")
	return nil, fmt.Errorf("Invalid token")
}

func (s *Service) DeleteToken(tokenString string) error {
	token, err := s.parseToken(tokenString)
	if err != nil {
		s.logger.Error().Msg(err.Error())
		return err
	}

	claims := token.Claims.(jwt.MapClaims)
	if token.Valid && fmt.Sprintf("%v", claims["type"]) == "refresh" {
		err = s.repository.DeleteRefreshToken(fmt.Sprintf("%v", claims["uuid"]))
		if err != nil {
			return err
		}
		return nil
	}
	s.logger.Error().Msg("Invalid token")
	return fmt.Errorf("Invalid token")
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
