package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

type ArgonSettings struct {
	hashLength uint32
	memory     uint32
	time       uint32
	threads    uint8
	saltLength uint32
}

func NewAuthSettings() ArgonSettings {
	hashLength, _ := strconv.ParseUint(os.Getenv("ARGON_HASH_LENGTH"), 10, 32)
	memory, _ := strconv.ParseUint(os.Getenv("ARGON_MEMORY"), 10, 32)
	time, _ := strconv.ParseUint(os.Getenv("ARGON_TIME"), 10, 32)
	threads, _ := strconv.ParseUint(os.Getenv("ARGON_THREADS"), 10, 32)
	saltLength, _ := strconv.ParseUint(os.Getenv("ARGON_SALT_LENGTH"), 10, 32)
	a := ArgonSettings{
		hashLength: uint32(hashLength),
		memory:     uint32(memory),
		time:       uint32(time),
		threads:    uint8(threads),
		saltLength: uint32(saltLength),
	}
	slog.Info("Argon2ID settings for the API",
		slog.Uint64("hashLength", uint64(a.hashLength)),
		slog.Uint64("memoryInBytes", uint64(a.memory * 1024)),
		slog.Uint64("times", uint64(a.time)),
		slog.Uint64("threads", uint64(a.threads)),
		slog.Uint64("saltLength", uint64(a.saltLength)),
	)
	return a
}

func (a ArgonSettings) generateSalt() []byte {
	saltBytes := make([]byte, a.saltLength)
	_, err := rand.Read(saltBytes)
	if err != nil {
		slog.Error("Unable to generate salt.", slog.String("error", err.Error()))
		os.Exit(1)
	}
	return saltBytes
}

func (a ArgonSettings) CreateNewHash(password string) string {
	saltBytes := a.generateSalt()
	hashBytes := argon2.IDKey([]byte(password), saltBytes, a.time, a.memory*1024, a.threads, a.hashLength)

	hash := base64.RawStdEncoding.EncodeToString(hashBytes)
	salt := base64.RawStdEncoding.EncodeToString(saltBytes)

	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, a.memory*1024, a.time, a.threads, salt, hash)
}

func (a ArgonSettings) CompareHash(password, hash string) bool {
	splitHash := strings.Split(hash, "$")

	hashVersion, err := strconv.Atoi(strings.Split(splitHash[2], "=")[1])
	if err != nil {
		slog.Error("Strconv error with parsing hash version", slog.String("error", err.Error()))
		os.Exit(1)
	}

	if hashVersion != argon2.Version {
		slog.Error("Mis-matching Argon2Id Version", slog.Int("currentVersion", argon2.Version), slog.Int("hashsVersion", hashVersion))
		os.Exit(1)
	}

	hashParams := strings.Split(splitHash[3], ",")

	hashMemory, err := strconv.ParseUint(strings.Split(hashParams[0], "=")[1], 10, 32)
	if err != nil {
		slog.Error("Strconv error with parsing hash memory", slog.String("error", err.Error()))
		os.Exit(1)
	}
	hashTime, err := strconv.ParseUint(strings.Split(hashParams[1], "=")[1], 10, 32)
	if err != nil {
		slog.Error("Strconv error with parsing hash time", slog.String("error", err.Error()))
		os.Exit(1)
	}

	hashThreads, err := strconv.ParseUint(strings.Split(hashParams[2], "=")[1], 10, 32)
	if err != nil {
		slog.Error("Strconv error with parsing hash threads", slog.String("error", err.Error()))
		os.Exit(1)
	}

	parsedHashSetting := ArgonSettings {
		hashLength: a.hashLength,
		memory: uint32(hashMemory),
		time: uint32(hashTime),
		threads: uint8(hashThreads),
		saltLength: a.saltLength,
	}

	slog.Debug("Argon2ID settings from parsed hash",
		slog.Uint64("hashLength", uint64(parsedHashSetting.hashLength)),
		slog.Uint64("memoryInBytes", uint64(parsedHashSetting.memory)),
		slog.Uint64("times", uint64(parsedHashSetting.time)),
		slog.Uint64("threads", uint64(parsedHashSetting.threads)),
		slog.Uint64("saltLength", uint64(parsedHashSetting.saltLength)),
	)

	saltBytes, err := base64.RawStdEncoding.DecodeString(splitHash[len(splitHash)-2])
	if err != nil {
		slog.Error("Unable to grab salt from argon2 hash.", slog.String("error", err.Error()))
		os.Exit(1)
	}

	hashBytes := argon2.IDKey([]byte(password), saltBytes, parsedHashSetting.time, parsedHashSetting.memory, parsedHashSetting.threads, parsedHashSetting.hashLength)

	passwordHash := base64.RawStdEncoding.EncodeToString(hashBytes)
	salt := base64.RawStdEncoding.EncodeToString(saltBytes)

	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, parsedHashSetting.memory, parsedHashSetting.time, parsedHashSetting.threads, salt, passwordHash) == hash
}
