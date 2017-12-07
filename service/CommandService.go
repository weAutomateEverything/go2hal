package service

import (
	"github.com/zamedic/go2hal/database"
	"io/ioutil"
	"os/exec"
	"fmt"
	"gopkg.in/kyokomi/emoji.v1"
	"time"
	"log"
	"errors"
	"runtime/debug"
)

/*
ExecuteRemoteCommand will run the command against the supplied address
 */
func ExecuteRemoteCommand(commandName, address string) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Print(err)
			SendError(errors.New(fmt.Sprint(err)))
			SendError(errors.New(string(debug.Stack())))

		}
	}()
	command, err := database.FindCommand(commandName)
	if err != nil {
		SendError(err)
		return err
	}

	key, err := database.GetKey()
	if err != nil {
		SendError(err)
		return err
	}

	err = ioutil.WriteFile("/tmp/key", []byte(key.Key), 0600)
	if err != nil {
		SendError(err)
		return err
	}

	SendAlert(emoji.Sprintf(":ghost: Executing Remote Command %s on machine %s", command, address))
	cmd := exec.Command("sh", "-c", fmt.Sprintf("ssh -i /tmp/key -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no %s@%s \"%s\" > /tmp/ssh.log", key.Username, address, command))
	log.Println(cmd.Args)
	stdout, err := cmd.StdoutPipe()

	go func() {
		err = cmd.Run()
		if err != nil {
			log.Println(err.Error())
			SendError(err)
		}
	}()

	count := 0

	for cmd.ProcessState == nil || !cmd.ProcessState.Exited() {
		time.Sleep(time.Second * 1)
		count++
		if count%60 == 0 {
			SendAlert(fmt.Sprintf("Still waiting for command %s to complete on %s", command, address))
		}
		if count > 600 {
			SendAlert(emoji.Sprintf(":bangbang: Timed Out waiting for command %s to complete on %s", command, address))
			return fmt.Errorf("timed out waiting for %s to complete on %s", command, address)
		}
	}

	if !cmd.ProcessState.Success() {
		SendAlert(emoji.Sprintf(":bangbang:  command %s did not complete successfully on %s", command, address))
	} else {
		SendAlert(emoji.Sprintf(":white_check_mark: command %s complete successfully on %s", command, address))
	}
	log.Println("output")
	log.Println(stdout)

	return nil
}
