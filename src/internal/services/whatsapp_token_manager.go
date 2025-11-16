package services

import (
	"context"
	"easy-queue-go/src/internal/log"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap"
)

// WhatsAppTokenManager manages access token refresh and validation
type WhatsAppTokenManager struct {
	currentToken  string
	tokenExpiry   time.Time
	appID         string
	appSecret     string
	httpClient    *http.Client
	mu            sync.RWMutex
	refreshTicker *time.Ticker
	stopChan      chan struct{}
}

// TokenDebugResponse represents the response from token debug endpoint
type TokenDebugResponse struct {
	Data struct {
		AppID       string `json:"app_id"`
		Type        string `json:"type"`
		Application string `json:"application"`
		ExpiresAt   int64  `json:"expires_at"`
		IsValid     bool   `json:"is_valid"`
		Scopes      []string `json:"scopes"`
	} `json:"data"`
}

// NewWhatsAppTokenManager creates a new token manager
func NewWhatsAppTokenManager(initialToken, appID, appSecret string) *WhatsAppTokenManager {
	return &WhatsAppTokenManager{
		currentToken: initialToken,
		appID:        appID,
		appSecret:    appSecret,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		stopChan: make(chan struct{}),
	}
}

// Start begins the automatic token refresh process
func (tm *WhatsAppTokenManager) Start(ctx context.Context) error {
	log.Info(ctx, "Starting WhatsApp token manager")

	// Validate current token and get expiry
	if err := tm.validateAndUpdateExpiry(ctx); err != nil {
		log.Warn(ctx, "Failed to validate initial token", zap.Error(err))
	}

	// Start refresh ticker (check every 6 hours)
	tm.refreshTicker = time.NewTicker(6 * time.Hour)

	go tm.refreshLoop(ctx)

	return nil
}

// Stop stops the token manager
func (tm *WhatsAppTokenManager) Stop() {
	if tm.refreshTicker != nil {
		tm.refreshTicker.Stop()
	}
	close(tm.stopChan)
}

// GetToken returns the current valid token
func (tm *WhatsAppTokenManager) GetToken() string {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.currentToken
}

// refreshLoop continuously checks and refreshes the token
func (tm *WhatsAppTokenManager) refreshLoop(ctx context.Context) {
	for {
		select {
		case <-tm.refreshTicker.C:
			if err := tm.checkAndRefresh(ctx); err != nil {
				log.Error(ctx, "Failed to refresh token", zap.Error(err))
			}
		case <-tm.stopChan:
			log.Info(ctx, "Stopping token manager")
			return
		}
	}
}

// checkAndRefresh checks token validity and refreshes if needed
func (tm *WhatsAppTokenManager) checkAndRefresh(ctx context.Context) error {
	tm.mu.RLock()
	expiry := tm.tokenExpiry
	tm.mu.RUnlock()

	// If token expires in less than 7 days, try to extend it
	if time.Until(expiry) < 7*24*time.Hour {
		log.Info(ctx, "Token expiring soon, attempting to extend",
			zap.Time("current_expiry", expiry),
		)

		if err := tm.extendToken(ctx); err != nil {
			log.Error(ctx, "Failed to extend token", zap.Error(err))
			return err
		}
	}

	return nil
}

// validateAndUpdateExpiry validates the current token and updates expiry time
func (tm *WhatsAppTokenManager) validateAndUpdateExpiry(ctx context.Context) error {
	tm.mu.RLock()
	token := tm.currentToken
	tm.mu.RUnlock()

	url := fmt.Sprintf("https://graph.facebook.com/debug_token?input_token=%s&access_token=%s|%s",
		token, tm.appID, tm.appSecret)

	resp, err := tm.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("failed to validate token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var debugResp TokenDebugResponse
	if err := json.Unmarshal(body, &debugResp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if !debugResp.Data.IsValid {
		return fmt.Errorf("token is invalid")
	}

	tm.mu.Lock()
	if debugResp.Data.ExpiresAt > 0 {
		tm.tokenExpiry = time.Unix(debugResp.Data.ExpiresAt, 0)
		log.Info(ctx, "Token validated",
			zap.Time("expires_at", tm.tokenExpiry),
			zap.Duration("time_until_expiry", time.Until(tm.tokenExpiry)),
		)
	} else {
		// Token never expires (system user token)
		tm.tokenExpiry = time.Now().Add(100 * 365 * 24 * time.Hour)
		log.Info(ctx, "Token is permanent (never expires)")
	}
	tm.mu.Unlock()

	return nil
}

// extendToken extends the current token's lifetime
func (tm *WhatsAppTokenManager) extendToken(ctx context.Context) error {
	tm.mu.RLock()
	token := tm.currentToken
	tm.mu.RUnlock()

	// Exchange short-lived token for long-lived token
	url := fmt.Sprintf("https://graph.facebook.com/oauth/access_token?grant_type=fb_exchange_token&client_id=%s&client_secret=%s&fb_exchange_token=%s",
		tm.appID, tm.appSecret, token)

	resp, err := tm.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("failed to extend token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var result struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int64  `json:"expires_in"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if result.AccessToken == "" {
		return fmt.Errorf("no access token in response: %s", string(body))
	}

	tm.mu.Lock()
	tm.currentToken = result.AccessToken
	if result.ExpiresIn > 0 {
		tm.tokenExpiry = time.Now().Add(time.Duration(result.ExpiresIn) * time.Second)
	}
	tm.mu.Unlock()

	log.Info(ctx, "Token extended successfully",
		zap.Time("new_expiry", tm.tokenExpiry),
	)

	return nil
}

// GetTokenInfo returns information about the current token
func (tm *WhatsAppTokenManager) GetTokenInfo() map[string]interface{} {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	return map[string]interface{}{
		"expires_at":        tm.tokenExpiry,
		"time_until_expiry": time.Until(tm.tokenExpiry).String(),
		"is_valid":          time.Now().Before(tm.tokenExpiry),
	}
}
