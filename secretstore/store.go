package secretstore

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/crypto/argon2"
)

const RefPrefix = "secret://"

const maxSecretSize = 64 * 1024

type encryptedPayload struct {
	Version    int    `json:"version"`
	Nonce      string `json:"nonce"`
	Ciphertext string `json:"ciphertext"`
}

type Store struct {
	path      string
	masterKey []byte
}

func ParseRef(value string) (string, bool) {
	if !strings.HasPrefix(value, RefPrefix) {
		return "", false
	}
	key := strings.TrimPrefix(value, RefPrefix)
	key = strings.TrimSpace(key)
	if key == "" {
		return "", false
	}
	return key, true
}

func New(path string, masterKey []byte) (*Store, error) {
	if strings.TrimSpace(path) == "" {
		return nil, errors.New("secret store path is required")
	}
	if len(masterKey) == 0 {
		return nil, errors.New("master key is required")
	}

	resolved := filepath.Clean(path)
	return &Store{path: resolved, masterKey: normalizeMasterKey(masterKey)}, nil
}

func (s *Store) Resolve(ref string) (string, error) {
	key, ok := ParseRef(ref)
	if !ok {
		return "", fmt.Errorf("invalid secret reference: %q", ref)
	}
	secrets, err := s.loadAll()
	if err != nil {
		return "", err
	}
	value, exists := secrets[key]
	if !exists {
		return "", fmt.Errorf("secret not found: %s", key)
	}
	return value, nil
}

func (s *Store) Set(key string, value string) error {
	if err := validateKey(key); err != nil {
		return err
	}
	if err := validateValue(value); err != nil {
		return err
	}
	secrets, err := s.loadAll()
	if err != nil {
		return err
	}
	secrets[key] = value
	return s.saveAll(secrets)
}

func (s *Store) Delete(key string) (bool, error) {
	if err := validateKey(key); err != nil {
		return false, err
	}
	secrets, err := s.loadAll()
	if err != nil {
		return false, err
	}
	if _, exists := secrets[key]; !exists {
		return false, nil
	}
	delete(secrets, key)
	if err := s.saveAll(secrets); err != nil {
		return false, err
	}
	return true, nil
}

func (s *Store) ListKeys() ([]string, error) {
	secrets, err := s.loadAll()
	if err != nil {
		return nil, err
	}
	keys := make([]string, 0, len(secrets))
	for key := range secrets {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys, nil
}

func (s *Store) loadAll() (map[string]string, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return map[string]string{}, nil
		}
		return nil, fmt.Errorf("failed to read secret store: %w", err)
	}
	if len(strings.TrimSpace(string(data))) == 0 {
		return map[string]string{}, nil
	}

	var payload encryptedPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse secret store: %w", err)
	}
	if payload.Version != 1 {
		return nil, fmt.Errorf("unsupported secret store version: %d", payload.Version)
	}

	nonce, err := base64.StdEncoding.DecodeString(payload.Nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to decode nonce: %w", err)
	}
	ciphertext, err := base64.StdEncoding.DecodeString(payload.Ciphertext)
	if err != nil {
		return nil, fmt.Errorf("failed to decode ciphertext: %w", err)
	}

	block, err := aes.NewCipher(s.masterKey)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize AEAD: %w", err)
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, errors.New("failed to decrypt secret store: invalid master key or corrupted data")
	}

	secrets := map[string]string{}
	if len(plaintext) == 0 {
		return secrets, nil
	}
	if err := json.Unmarshal(plaintext, &secrets); err != nil {
		return nil, fmt.Errorf("failed to parse decrypted secret data: %w", err)
	}
	return secrets, nil
}

func (s *Store) saveAll(secrets map[string]string) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o700); err != nil {
		return fmt.Errorf("failed to create secret store directory: %w", err)
	}

	plaintext, err := json.Marshal(secrets)
	if err != nil {
		return fmt.Errorf("failed to encode secret data: %w", err)
	}

	block, err := aes.NewCipher(s.masterKey)
	if err != nil {
		return fmt.Errorf("failed to initialize cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("failed to initialize AEAD: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return fmt.Errorf("failed to generate nonce: %w", err)
	}

	payload := encryptedPayload{
		Version:    1,
		Nonce:      base64.StdEncoding.EncodeToString(nonce),
		Ciphertext: base64.StdEncoding.EncodeToString(gcm.Seal(nil, nonce, plaintext, nil)),
	}

	encoded, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to encode encrypted payload: %w", err)
	}

	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, encoded, 0o600); err != nil {
		return fmt.Errorf("failed to write secret store temp file: %w", err)
	}
	if err := os.Rename(tmp, s.path); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("failed to atomically update secret store: %w", err)
	}
	return nil
}

func MasterKeyFromEnv() ([]byte, error) {
	raw := strings.TrimSpace(os.Getenv("FOGHORN_SECRET_MASTER_KEY"))
	if raw == "" {
		return nil, errors.New("FOGHORN_SECRET_MASTER_KEY is required")
	}

	var input []byte
	if decoded, err := base64.StdEncoding.DecodeString(raw); err == nil && len(decoded) > 0 {
		input = decoded
	} else {
		input = []byte(raw)
	}

	if len(input) < 32 {
		return nil, errors.New("master key must be at least 32 characters long")
	}

	return normalizeMasterKey(input), nil
}

func normalizeMasterKey(key []byte) []byte {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		panic(err)
	}
	return argon2.IDKey(key, salt, 1, 64*1024, 4, 32)
}

func validateKey(key string) error {
	key = strings.TrimSpace(key)
	if key == "" {
		return errors.New("secret key is required")
	}
	if strings.Contains(key, "..") {
		return errors.New("secret key cannot contain '..'")
	}
	if strings.HasPrefix(key, "/") {
		return errors.New("secret key must be relative")
	}
	return nil
}

func validateValue(value string) error {
	if len(value) > maxSecretSize {
		return fmt.Errorf("secret value exceeds maximum size of %d bytes", maxSecretSize)
	}
	return nil
}
