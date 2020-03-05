package commands

import (
	"asd/common/helpers"
	"errors"
)

func Init(pool string) error {

	// init logs ====================================
	_log := helpers.InitLogs(true)

	client, err := dialGrpc()
	if err != nil {

	}

	var err = errors.New("not implemented")
	return err
}
