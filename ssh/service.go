package ssh

import (
	"errors"
	"fmt"
	"github.com/zamedic/go2hal/alert"
	"gopkg.in/kyokomi/emoji.v1"
	"io/ioutil"
	"log"
	"os/exec"
	"runtime/debug"
	"time"
)

type Service interface {
	ExecuteRemoteCommand(commandName, address string) error
}

type service struct {
	alert alert.Service
	store Store
}

func NewService(alert alert.Service, store Store) Service {
	return &service{alert, store}
}

/*
ExecuteRemoteCommand will run the command against the supplied address
*/
func (s *service) ExecuteRemoteCommand(commandName, address string) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Print(err)
			s.alert.SendError(errors.New(fmt.Sprint(err)))
			s.alert.SendError(errors.New(string(debug.Stack())))

		}
	}()
	command, err := s.store.findCommand(commandName)
	if err != nil {
		s.alert.SendError(err)
		return err
	}

	key, err := s.store.getKey()
	if err != nil {
		s.alert.SendError(err)
		return err
	}

	err = ioutil.WriteFile("/tmp/key", []byte(key.Key), 0600)
	if err != nil {
		s.alert.SendError(err)
		return err
	}

	s.alert.SendAlert(emoji.Sprintf(":ghost: Executing Remote Command %s on machine %s", command, address))
	cmd := exec.Command("sh", "-c", fmt.Sprintf("ssh -i /tmp/key -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no %s@%s \"%s\" > /tmp/ssh.log", key.Username, address, command))
	log.Println(cmd.Args)
	stdout, err := cmd.StdoutPipe()

	go func() {
		err = cmd.Run()
		if err != nil {
			log.Println(err.Error())
			s.alert.SendError(err)
		}
	}()

	count := 0

	for cmd.ProcessState == nil || !cmd.ProcessState.Exited() {
		time.Sleep(time.Second * 1)
		count++
		if count%60 == 0 {
			s.alert.SendAlert(fmt.Sprintf("Still waiting for command %s to complete on %s", command, address))
		}
		if count > 600 {
			s.alert.SendAlert(emoji.Sprintf(":bangbang: Timed Out waiting for command %s to complete on %s", command, address))
			return fmt.Errorf("timed out waiting for %s to complete on %s", command, address)
		}
	}

	if !cmd.ProcessState.Success() {
		s.alert.SendAlert(emoji.Sprintf(":bangbang:  command %s did not complete successfully on %s", command, address))
	} else {
		s.alert.SendAlert(emoji.Sprintf(":white_check_mark: command %s complete successfully on %s", command, address))
	}
	log.Println("output")
	log.Println(stdout)

	return nil
}
