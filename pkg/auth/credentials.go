package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/scrypt"
)

type Credentials struct {
	ID       int    `json:"id"`
	Server   string `json:"server"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type CredentialStore struct {
	db     *sql.DB
	cipher cipher.AEAD
}

const (
	createTableSQL = `
	CREATE TABLE IF NOT EXISTS credentials (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		server TEXT NOT NULL,
		username TEXT NOT NULL,
		password TEXT NOT NULL,
		role TEXT,
		UNIQUE(server, username)
	);`
)

func NewCredentialStore(masterKey string) (*CredentialStore, error) {
	// Create the .infra directory in user's home if it doesn't exist
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	dbDir := filepath.Join(home, ".infra")
	if err := os.MkdirAll(dbDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// Open SQLite database
	dbPath := filepath.Join(dbDir, "credentials.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Create table if it doesn't exist
	if _, err := db.Exec(createTableSQL); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	// Initialize encryption
	key, err := deriveKey(masterKey)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to derive key: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	return &CredentialStore{
		db:     db,
		cipher: gcm,
	}, nil
}

func (cs *CredentialStore) Close() error {
	return cs.db.Close()
}

func (cs *CredentialStore) SaveCredentials(creds Credentials) error {
	// Encrypt password
	encryptedPass, err := cs.encrypt([]byte(creds.Password))
	if err != nil {
		return fmt.Errorf("failed to encrypt password: %w", err)
	}

	// Save to database
	_, err = cs.db.Exec(
		"INSERT OR REPLACE INTO credentials (server, username, password, role) VALUES (?, ?, ?, ?)",
		creds.Server,
		creds.Username,
		encryptedPass,
		creds.Role,
	)
	if err != nil {
		return fmt.Errorf("failed to save credentials: %w", err)
	}

	return nil
}

func (cs *CredentialStore) GetCredentials(server, username string) (*Credentials, error) {
	var creds Credentials
	var encryptedPass string

	err := cs.db.QueryRow(
		"SELECT id, server, username, password, role FROM credentials WHERE server = ? AND username = ?",
		server,
		username,
	).Scan(&creds.ID, &creds.Server, &creds.Username, &encryptedPass, &creds.Role)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get credentials: %w", err)
	}

	// Decrypt password
	decryptedPass, err := cs.decrypt(encryptedPass)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt password: %w", err)
	}

	creds.Password = string(decryptedPass)
	return &creds, nil
}

func (cs *CredentialStore) ListCredentials() ([]Credentials, error) {
	rows, err := cs.db.Query("SELECT id, server, username, password, role FROM credentials")
	if err != nil {
		return nil, fmt.Errorf("failed to list credentials: %w", err)
	}
	defer rows.Close()

	var credentials []Credentials
	for rows.Next() {
		var cred Credentials
		var encryptedPass string
		if err := rows.Scan(&cred.ID, &cred.Server, &cred.Username, &encryptedPass, &cred.Role); err != nil {
			return nil, fmt.Errorf("failed to scan credential: %w", err)
		}

		// Decrypt password
		decryptedPass, err := cs.decrypt(encryptedPass)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt password: %w", err)
		}

		cred.Password = string(decryptedPass)
		credentials = append(credentials, cred)
	}

	return credentials, nil
}

func (cs *CredentialStore) DeleteCredentials(server, username string) error {
	_, err := cs.db.Exec("DELETE FROM credentials WHERE server = ? AND username = ?", server, username)
	if err != nil {
		return fmt.Errorf("failed to delete credentials: %w", err)
	}
	return nil
}

func (cs *CredentialStore) encrypt(data []byte) (string, error) {
	nonce := make([]byte, cs.cipher.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := cs.cipher.Seal(nonce, nonce, data, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (cs *CredentialStore) decrypt(encodedData string) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < cs.cipher.NonceSize() {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce := ciphertext[:cs.cipher.NonceSize()]
	ciphertext = ciphertext[cs.cipher.NonceSize():]

	plaintext, err := cs.cipher.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func deriveKey(masterKey string) ([]byte, error) {
	salt := []byte("infra-cli-salt") // In production, use a random salt and store it
	return scrypt.Key([]byte(masterKey), salt, 32768, 8, 1, 32)
}
