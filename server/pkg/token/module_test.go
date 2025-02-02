package token

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestSetupLogger(t *testing.T) {
	scope := "jwt"
	m := Module{
		scope: "jwt",
	}
	p := Params{
		Logger: zap.NewExample(),
	}
	m.logger = m.setupLogger(scope, p)

	assert.NotNil(t, m.logger)
	assert.IsType(t, &zap.Logger{}, m.logger)
}

func TestSetupConfig(t *testing.T) {
	m := Module{
		logger: zap.NewExample(),
	}

	t.Run("TestNoTokenScopesPassed", func(t *testing.T) {
		tokenScopes := []string{}
		configs := m.setupConfig(tokenScopes...)

		assert.Equal(t, 1, len(configs))
		// Check if the configurations are correctly set for the default scope
		config, exists := configs[defaultTokenScope]
		assert.True(t, exists)
		assert.Equal(t, defaultTokenLookup, config.TokenLookup)
		assert.Equal(t, defaultSigningKey, config.SigningKey)
		assert.Equal(t, defaultSigningMethod, config.SigningMethod)
		assert.Equal(t, defaultExpInHours, config.ExpInHours)
	})

	t.Run("TestTokenScopesWithNoConfig", func(t *testing.T) {
		tokenScopes := []string{"scope1", "scope2"}
		configs := m.setupConfig(tokenScopes...)

		assert.Equal(t, 2, len(configs))
		// Check if default configurations are set for each scope
		for _, scope := range tokenScopes {
			config, exists := configs[scope]
			assert.True(t, exists)
			assert.Equal(t, defaultTokenLookup, config.TokenLookup)
			assert.Equal(t, defaultSigningKey, config.SigningKey)
			assert.Equal(t, defaultSigningMethod, config.SigningMethod)
			assert.Equal(t, defaultExpInHours, config.ExpInHours)
		}
	})

	t.Run("TestTokenScopesWithConfig", func(t *testing.T) {
		tokenScopes := []string{"scope1", "scope2"}

		viper.Set("scope1.token_lookup", "header:Authorization")
		viper.Set("scope1.signing_key", "my_secret")
		viper.Set("scope1.signing_method", "HS512")
		viper.Set("scope1.exp_in_hours", 24)

		viper.Set("scope2.token_lookup", "query:token")
		viper.Set("scope2.signing_key", "my_secret2")
		viper.Set("scope2.signing_method", "HS256")
		viper.Set("scope2.exp_in_hours", 48)

		configs := m.setupConfig(tokenScopes...)

		assert.Equal(t, 2, len(configs))

		for _, scope := range tokenScopes {
			config, exists := configs[scope]
			assert.True(t, exists)
			assert.Equal(t, viper.GetString(scope+".token_lookup"), config.TokenLookup)
			assert.Equal(t, viper.GetString(scope+".signing_key"), config.SigningKey)
			assert.Equal(t, viper.GetString(scope+".signing_method"), config.SigningMethod)
			assert.Equal(t, viper.GetInt(scope+".exp_in_hours"), config.ExpInHours)
		}
	})

	t.Run("TestSetupWithPartialConfig", func(t *testing.T) {
		viper.Reset()
		defer viper.Reset()

		viper.Set("scope1.token_lookup", "header:Authorization")
		viper.Set("scope1.signing_key", "test_secret")

		configs := m.setupConfig("scope1", "scope2")

		assert.Equal(t, 2, len(configs))

		config, exists := configs["scope1"]
		assert.True(t, exists)
		assert.Equal(t, viper.GetString("scope1.token_lookup"), config.TokenLookup)
		assert.Equal(t, viper.GetString("scope1.signing_key"), config.SigningKey)
		assert.Equal(t, defaultSigningMethod, config.SigningMethod)
		assert.Equal(t, defaultExpInHours, config.ExpInHours)

		config, exists = configs["scope2"]
		assert.True(t, exists)
		assert.Equal(t, defaultTokenLookup, config.TokenLookup)
		assert.Equal(t, defaultSigningKey, config.SigningKey)
		assert.Equal(t, defaultSigningMethod, config.SigningMethod)
		assert.Equal(t, defaultExpInHours, config.ExpInHours)

	})
}

func TestGetConfigHelper(t *testing.T) {
	m := Module{
		configs: map[string]*Config{
			"scope1": {
				TokenLookup:   "header:Authorization",
				SigningKey:    "my_secret",
				SigningMethod: "HS512",
				ExpInHours:    24,
			},
			"scope2": {
				TokenLookup:   "query:token",
				SigningKey:    "my_secret2",
				SigningMethod: "HS256",
				ExpInHours:    48,
			},
		},
	}

	t.Run("ExistingScope", func(t *testing.T) {
		scope := "scope1"
		config, err := m.getConfigHelper(scope)
		assert.NoError(t, err)
		assert.Equal(t, "header:Authorization", config.TokenLookup)
		assert.Equal(t, "my_secret", config.SigningKey)
		assert.Equal(t, "HS512", config.SigningMethod)
		assert.Equal(t, 24, config.ExpInHours)
	})

	t.Run("NonExistingScope", func(t *testing.T) {
		scope := "non_existing_scope"
		config, err := m.getConfigHelper(scope)
		assert.Error(t, err)
		assert.Nil(t, config)
	})
}

func TestGenerateToken(t *testing.T) {
	m := Module{
		configs: map[string]*Config{
			"scope1": {
				TokenLookup:   "header:Authorization",
				SigningKey:    "my_secret",
				SigningMethod: "HS512",
				ExpInHours:    24,
			},
		},
		logger: zap.NewExample(),
	}
	scope := "scope1"
	additionalClaims := jwt.MapClaims{
		"sub":  "user123",
		"role": "admin",
	}
	token, err := m.GenerateToken(scope, additionalClaims)
	assert.NoError(t, err)
	assert.NotNil(t, token)
	//todo: currently, only asserting that token exists
	//todo: need to check if token is correct?
}

func TestGetJWTMiddleware(t *testing.T) {
	m := Module{
		configs: map[string]*Config{
			"scope1": {
				TokenLookup:   "header:Authorization",
				SigningKey:    "my_secret",
				SigningMethod: "HS512",
				ExpInHours:    24,
			},
		},
		logger: zap.NewExample(),
	}
	scope := "scope1"
	middleware := m.GetJWTMiddleware(scope)
	assert.NotNil(t, middleware)
	//todo: currently, only asserting that middleware exists
	//todo: need to check if middleware is correct?
}
