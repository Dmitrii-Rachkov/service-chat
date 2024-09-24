package encryption

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_mySecret(t *testing.T) {
	exSecret := "abc&1*~#^2^#s0^=)^^7%b34"
	acSecret, err := mySecret()

	// Сравниваем ожидаемый и актуальный результат
	assert.NoError(t, err)
	assert.Equal(t, exSecret, acSecret)
}

func Test_encode(t *testing.T) {
	testTable := []struct {
		name     string
		input    []byte
		expected string
	}{
		{
			name:     "OK",
			input:    []byte("hello world"),
			expected: "aGVsbG8gd29ybGQ=",
		},
		{
			name:     "Empty input",
			input:    []byte(""),
			expected: "",
		},
		{
			name:     "Nil input",
			expected: "",
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, encode(tt.input))
		})
	}
}

func Test_decode(t *testing.T) {
	testTable := []struct {
		name     string
		input    string
		expected []byte
	}{
		{
			name:     "OK",
			input:    "aGVsbG8gd29ybGQ=",
			expected: []byte("hello world"),
		},
		{
			name:     "Empty input",
			input:    "",
			expected: []byte(""),
		},
		{
			name:     "Nil input",
			expected: []byte(""),
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := decode(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestEncrypt(t *testing.T) {
	testTable := []struct {
		name     string
		input    string
		expected string
		setEnv   bool
		secret   string
	}{
		{
			name:     "OK",
			input:    "hello world",
			expected: "AyVL7wcMVPcTbss=",
		},
		{
			name:     "OK empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "OK nil string",
			expected: "",
		},
		{
			name:     "Random secret",
			input:    "hello world",
			expected: "encryption.Encrypt: crypto/aes: invalid key size 6",
			setEnv:   true,
			secret:   "random",
		},
		{
			name:     "Empty secret",
			input:    "hello world",
			expected: "encryption.Encrypt: crypto/aes: invalid key size 0",
			setEnv:   true,
			secret:   "",
		},
		{
			name:     "Nil secret",
			input:    "hello world",
			expected: "encryption.Encrypt: crypto/aes: invalid key size 0",
			setEnv:   true,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				t.Setenv("MY_SECRET", tt.secret)
				actual, err := Encrypt(tt.input)
				assert.Empty(t, actual)
				assert.Equal(t, err.Error(), tt.expected)
			} else {
				actual, err := Encrypt(tt.input)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, actual)
			}
		})
	}
}

func TestDecrypt(t *testing.T) {
	testTable := []struct {
		name     string
		input    string
		expected string
		setEnv   bool
		secret   string
	}{
		{
			name:     "OK",
			input:    "AyVL7wcMVPcTbss=",
			expected: "hello world",
		},
		{
			name:     "OK empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "OK nil string",
			expected: "",
		},
		{
			name:     "Random secret",
			input:    "AyVL7wcMVPcTbss=",
			expected: "encryption.Decrypt: crypto/aes: invalid key size 6",
			setEnv:   true,
			secret:   "random",
		},
		{
			name:     "Empty secret",
			input:    "hello world",
			expected: "encryption.Decrypt: crypto/aes: invalid key size 0",
			setEnv:   true,
			secret:   "",
		},
		{
			name:     "Nil secret",
			input:    "AyVL7wcMVPcTbss=",
			expected: "encryption.Decrypt: crypto/aes: invalid key size 0",
			setEnv:   true,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				t.Setenv("MY_SECRET", tt.secret)
				actual, err := Decrypt(tt.input)
				assert.Empty(t, actual)
				assert.Equal(t, err.Error(), tt.expected)
			} else {
				actual, err := Decrypt(tt.input)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, actual)
			}
		})
	}
}
