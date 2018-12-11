package ssh

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/weAutomateEverything/go2hal/alert"
	"golang.org/x/net/context"
	"gopkg.in/kyokomi/emoji.v1"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime/debug"
	"time"
)

type Service interface {
	ExecuteRemoteCommand(ctx context.Context, chatId uint32, commandName, address string) error
	addCommand(chatId uint32, name, command string) error
	addKey(chatId uint32, userName, base64Key string) error
	addServer(chatId uint32, address, description string) error
}

func NewService(alert alert.Service, store Store) Service {
	return &service{alert, store}
}

type service struct {
	alert alert.Service
	store Store
}

func (s *service) addServer(chatId uint32, address, description string) error {
	return s.store.addServer(chatId, address, description)
}

/*
ExecuteRemoteCommand will run the command against the supplied address
*/
func (s *service) ExecuteRemoteCommand(ctx context.Context, chatId uint32, commandName, address string) (err error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Print(err)
			s.alert.SendError(ctx, errors.New(fmt.Sprint(err)))
			s.alert.SendError(ctx, errors.New(string(debug.Stack())))

		}
	}()
	command, err := s.store.findCommand(chatId, commandName)
	if err != nil {
		s.alert.SendError(ctx, err)
		return err
	}

	key, err := s.store.getKey(chatId)
	if err != nil {
		s.alert.SendError(ctx, err)
		return err
	}
	d := make([]byte, base64.StdEncoding.DecodedLen(len(key.EncryptedBase64Key)))
	base64.StdEncoding.Decode(d, []byte(key.EncryptedBase64Key))

	base64Key, err := decrypt(d, os.Getenv("ENCRYPTION_KEY"))
	if err != nil {
		err = fmt.Errorf("unable to decrypt data. %v", err.Error())
		return
	}

	decryptedKey := make([]byte, len(base64Key))
	_, err = base64.StdEncoding.Decode(decryptedKey, base64Key)

	if err != nil {
		err = fmt.Errorf("unable to base64 decode your key. %v", err.Error())
		return
	}

	err = ioutil.WriteFile("/tmp/key", []byte(decryptedKey), 0600)
	if err != nil {
		s.alert.SendError(ctx, err)
		return err
	}

	s.alert.SendAlert(ctx, chatId, emoji.Sprintf(":ghost: Executing Remote Command %s on machine %s", command, address))
	cmd := exec.Command("sh", "-c", fmt.Sprintf("ssh -i /tmp/key -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no %s@%s \"%s\" > /tmp/ssh.log", key.Username, address, command))
	log.Println(cmd.Args)
	stdout, err := cmd.StdoutPipe()

	go func() {
		err = cmd.Run()
		if err != nil {
			log.Println(err.Error())
			s.alert.SendError(ctx, err)
		}
	}()

	count := 0

	for cmd.ProcessState == nil || !cmd.ProcessState.Exited() {
		time.Sleep(time.Second * 1)
		count++
		if count%60 == 0 {
			s.alert.SendAlert(ctx, chatId, fmt.Sprintf("Still waiting for command %s to complete on %s", command, address))
		}
		if count > 600 {
			s.alert.SendAlert(ctx, chatId, emoji.Sprintf(":bangbang: Timed Out waiting for command %s to complete on %s", command, address))
			return fmt.Errorf("timed out waiting for %s to complete on %s", command, address)
		}
	}

	if !cmd.ProcessState.Success() {
		s.alert.SendAlert(ctx, chatId, emoji.Sprintf(":bangbang:  command %s did not complete successfully on %s", command, address))
	} else {
		s.alert.SendAlert(ctx, chatId, emoji.Sprintf(":white_check_mark: command %s complete successfully on %s", command, address))
	}
	log.Println("output")
	b, err := ioutil.ReadAll(stdout)
	log.Println(string(b))

	return nil
}

func (s *service) addCommand(chatId uint32, name, command string) error {
	return s.store.addCommand(chatId, name, command)
}

func (s *service) addKey(chatId uint32, userName, base64Key string) error {
	key, err := encrypt([]byte(base64Key), os.Getenv("ENCRYPTION_KEY"))
	if err != nil {
		return err
	}

	v := base64.StdEncoding.EncodeToString(key)

	return s.store.addKey(chatId, userName, v)
}

func encrypt(data []byte, passphrase string) ([]byte, error) {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func decrypt(data []byte, passphrase string) ([]byte, error) {
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext, err
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}
