package ssh

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/weAutomateEverything/go2hal/alert"
	"github.com/weAutomateEverything/go2hal/telegram"
	"golang.org/x/crypto/ssh"
	"golang.org/x/net/context"
	"gopkg.in/kyokomi/emoji.v1"
	"io"
	"log"
	"net"
	"os"
	"runtime/debug"
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
	d := make([]byte, base64.RawStdEncoding.DecodedLen(len(key.EncryptedBase64Key)))
	_, err = base64.RawStdEncoding.Decode(d, []byte(key.EncryptedBase64Key))
	if err != nil {
		err = fmt.Errorf("unable to base64 decode encrypted value. %v", err.Error())
		return
	}

	base64Key, err := decrypt(d, os.Getenv("ENCRYPTION_KEY"))
	if err != nil {
		err = fmt.Errorf("unable to decrypt data. %v", err.Error())
		return
	}

	decryptedKey := make([]byte, base64.RawStdEncoding.DecodedLen(len(base64Key)))
	_, err = base64.RawStdEncoding.Decode(decryptedKey, base64Key)

	if err != nil {
		err = fmt.Errorf("unable to base64 decode your key. %v", err.Error())
		return
	}

	s.alert.SendAlert(ctx, chatId, emoji.Sprintf(":ghost: Executing Remote Command %s on machine %s", command, address))
	resp, err := remoteRun(key.Username, address, string(decryptedKey), command)

	if err != nil {
		log.Println(err)
		s.alert.SendAlert(ctx, chatId, emoji.Sprintf(":bangbang:  command %s did not complete successfully on %s. %v. %v", command, address, telegram.Escape(resp), err))
	} else {
		s.alert.SendAlert(ctx, chatId, emoji.Sprintf(":white_check_mark: command %s complete successfully on %s. %v", command, address, telegram.Escape(resp)))
	}
	log.Println("output")
	log.Println(string(resp))

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

	v := base64.RawStdEncoding.EncodeToString(key)

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
		log.Printf("Bad Key %v", err)
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Printf("Bad GCM %v", err)
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	return plaintext, err
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func remoteRun(user string, addr string, privateKey string, cmd string) (string, error) {
	// privateKey could be read from a file, or retrieved from another storage
	// source, such as the Secret Service / GNOME Keyring
	key, err := ssh.ParsePrivateKey([]byte(privateKey))
	if err != nil {
		return "", err
	}
	// Authentication
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		//alternatively, you could use a password
		/*
			Auth: []ssh.AuthMethod{
				ssh.Password("PASSWORD"),
			},
		*/
	}
	// Connect
	client, err := ssh.Dial("tcp", addr+":22", config)
	if err != nil {
		return "", err
	}
	// Create a session. It is one session per command.
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()
	var b bytes.Buffer // import "bytes"
	var e bytes.Buffer
	session.Stdout = &b // get output
	session.Stderr = &e
	// you can also pass what gets input to the stdin, allowing you to pipe
	// content from client to server
	//      session.Stdin = bytes.NewBufferString("My input")

	// Finally, run the command
	err = session.Run(cmd)

	return e.String() + "\n" + b.String(), err
}
