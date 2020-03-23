package grpc

import (
	"asd/common/api"
	"asd/common/helpers"
	"asd/common/zfs"
	"context"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"time"
)

func (s *Server) SolutionList(ctx context.Context, in *api.Void) (*api.Solutions, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: SolutionList")

	var apiSolutionsVal api.Solutions
	apiSolutions := &apiSolutionsVal
	var apiOutcome api.Outcome
	apiSolutions.Outcome = &apiOutcome

	// get pool name from config
	pool := viper.GetString("pool")
	if pool == "" {
		message := "Init config value is empty for master pool"
		_log.Error(message)
		err := errors.New(message)
		return apiSolutions, err
	} else {
		_log.Info("got pool from config file:" + pool)
	}

	// create default master dataset name and get it via zfs wrap
	dataset, err := zfs.GetDataset(pool + "/asd")
	if err != nil {
		message := "Unable to locate the master dataset: did you run 'asd master init'?"
		_log.Error(message)
		_log.Error(err)
		return apiSolutions, err
	} else {
		_log.Info("Master dataset found:" + dataset.Name)
	}

	// get the actual mountpoint
	mountpoint, err := dataset.GetProperty("mountpoint")
	if err != nil {
		message := "Unable to locate the mountpoint of the master dataset"
		_log.Error(message)
		_log.Error(err)
		return apiSolutions, err
	} else {
		_log.Info("Got mountpoint: " + mountpoint)
	}

	// open or create the k/v db
	db, err := bolt.Open(mountpoint+"/asd.db", 0600, &bolt.Options{Timeout: 3 * time.Second})
	if err != nil {
		message := "Unable to open the main db for persisting node info"
		_log.Error(message)
		_log.Error(err)
		return apiSolutions, err
	} else {
		_log.Info("succesfully opened main db")
	}
	defer db.Close()

	// add node info to the k/v db
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("solutions"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var apiSolutionVal api.Solution
			apiSolution := &apiSolutionVal
			err = proto.Unmarshal(v, apiSolution)
			if err != nil {
				return fmt.Errorf("unmarshaling solution proto: %s", err)
			}
			apiSolutions.Solutions = append(apiSolutions.Solutions, apiSolution)
		}
		return nil
	})
	if err != nil {
		message := "Unable to list solutions from k/v db"
		_log.Error(message)
		_log.Error(err)
		apiSolutions.Outcome.Message = message
		return apiSolutions, err
	}

	message := "Succesfully obtained solution list"
	_log.Info(message)
	apiSolutions.Outcome.Message = message
	return apiSolutions, nil
}

func (s *Server) SolutionCopy(ctx context.Context, in *api.CopyArgs) (*api.Outcome, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Copy")

	copyArgs := in
	var apiOutcomeVal api.Outcome
	apiOutcome := &apiOutcomeVal

	var pool string
	pool = viper.GetString("pool")
	if pool == "" {
		message := "The master pool is unconfigured"
		_log.Error(message)
		apiOutcome.Error = true
		apiOutcome.Message = message
		err := errors.New(message)
		return apiOutcome, err
	}

	var dbPath string
	dbPath = viper.GetString("mountpoint") + "/asd.db"
	if dbPath == "" {
		message := "The master pool is unconfigured (dbpath)"
		_log.Error(message)
		apiOutcome.Error = true
		apiOutcome.Message = message
		err := errors.New(message)
		return apiOutcome, err
	}

	sourceName := pool + "/asd/" + copyArgs.Source
	destName := pool + "/asd/" + copyArgs.Destination

	//- [ ]  check if A exists and A.ORIG and B names are available
	sourceDataset, err := zfs.GetDataset(sourceName)
	if err != nil {
		message := "Unable to locate the source dataset"
		_log.Error(message)
		_log.Error(err)
		apiOutcome.Message = message
		apiOutcome.Error = true
		return apiOutcome, err
	} else {
		_log.Info("Source dataset found:" + sourceName)
	}
	_, err = zfs.GetDataset(destName)
	if err == nil {
		message := "The destination dataset exists: " + destName
		_log.Error(message)
		_log.Error(err)
		apiOutcome.Message = message
		apiOutcome.Error = true
		return apiOutcome, err
	} else {
		_log.Info("Destination dataset not already present:" + destName)
	}
	_, err = zfs.GetDataset(sourceName + ".orig")
	if err == nil {
		message := "The destination dataset exists: " + sourceName + ".orig"
		_log.Error(message)
		_log.Error(err)
		apiOutcome.Message = message
		apiOutcome.Error = true
		return apiOutcome, err
	} else {
		_log.Info("Source .origin temporary dataset name not already present:" + sourceName)
	}

	//- [ ]  snap A@passage
	sourceDataset, err = sourceDataset.Snapshot("passage", false)
	if err != nil {
		message := "Unable to snap the source dataset"
		_log.Error(message)
		_log.Error(err)
		apiOutcome.Message = message
		apiOutcome.Error = true
		return apiOutcome, err
	} else {
		_log.Info("Source dataset snapshot made:" + copyArgs.Source + "@clone")
	}

	//- [ ]  rename A → A.ORIG
	_, err = zfs.Command("rename", sourceName, sourceName+".orig")
	if err != nil {
		message := "Unable to rename the source dataset"
		_log.Error(message)
		_log.Error(err)
		apiOutcome.Message = message
		apiOutcome.Error = true
		return apiOutcome, err
	} else {
		_log.Info("Source dataset renamed")
	}

	//- [ ]  clone A.ORIG@clone → A (a get all snapshots)
	_, err = zfs.Command("clone", sourceName+".orig@passage", sourceName)
	if err != nil {
		message := "Unable to clone A.ORIG@clone → A"
		_log.Error(message)
		_log.Error(err)
		apiOutcome.Message = message
		apiOutcome.Error = true
		return apiOutcome, err
	} else {
		_log.Info("Source dataset cloned")
	}

	//- [ ]  promote A (a retain all snapshots)
	_, err = zfs.Command("promote", sourceName)
	if err != nil {
		message := "Unable to clone A.ORIG@clone → A"
		_log.Error(message)
		_log.Error(err)
		apiOutcome.Message = message
		apiOutcome.Error = true
		return apiOutcome, err
	} else {
		_log.Info("Source dataset cloned")
	}

	//- [ ]  rename A.ORIG → B
	_, err = zfs.Command("rename", sourceName+".orig", destName)
	if err != nil {
		message := "Unable to rename A.ORIG → B"
		_log.Error(message)
		_log.Error(err)
		apiOutcome.Message = message
		apiOutcome.Error = true
		return apiOutcome, err
	} else {
		_log.Info("Source dataset renamed")
	}

	//- [ ]  promote B
	_, err = zfs.Command("promote", destName)
	if err != nil {
		message := "Unable to rename A.ORIG → B"
		_log.Error(message)
		_log.Error(err)
		apiOutcome.Message = message
		apiOutcome.Error = true
		return apiOutcome, err
	} else {
		_log.Info("Destination dataset promoted")
	}

	// open or create the k/v db
	db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 3 * time.Second})
	if err != nil {
		message := "Unable to open the main db for persisting master info"
		_log.Error(message)
		_log.Error(err)
		apiOutcome.Message = message
		return apiOutcome, err
	} else {
		_log.Info("Main db opened succesfully")
	}
	defer db.Close()

	// add solution info to the k/v db
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("solutions"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		var apiSolutionVal api.Solution
		apiSolution := &apiSolutionVal
		apiSolution.Name = copyArgs.Destination
		apiSolution.Hostname = "master"
		apiSolution.Status = "available"

		var encoded []byte
		encoded, err = proto.Marshal(apiSolution)
		if err != nil {
			return err
		}
		_log.Debugf("encoded lenght before put: %v\n", len(encoded))

		err = b.Put([]byte(apiSolution.Name), encoded)
		if err != nil {
			return fmt.Errorf("put: %s", err)
		}

		encoded2 := b.Get([]byte(apiSolution.Name))
		_log.Debugf("encoded lenght after get: %v\n", len(encoded2))

		return nil
	})
	if err != nil {
		message := "Unable to update db to persist solution info"
		_log.Error(message)
		_log.Error(err)
		apiOutcome.Message = message
		return apiOutcome, err
	} else {
		_log.Info("Main db updated with master info")
	}

	message := "Succesfully copied " + copyArgs.Source + " into " + copyArgs.Destination
	_log.Info(message)
	apiOutcome.Error = false
	apiOutcome.Message = message
	return apiOutcome, nil
}

func (s *Server) SolutionCreate(ctx context.Context, in *api.Solution) (*api.Solution, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Create")

	apiSolution := in
	var apiOutcome api.Outcome
	apiSolution.Outcome = &apiOutcome

	var pool string
	pool = viper.GetString("pool")
	if pool == "" {
		message := "The master pool is unconfigured"
		_log.Error(message)
		apiSolution.Outcome.Error = true
		apiSolution.Outcome.Message = message
		err := errors.New(message)
		return apiSolution, err
	}

	var dbPath string
	dbPath = viper.GetString("mountpoint") + "/asd.db"
	if dbPath == "" {
		message := "The master pool is unconfigured (dbpath)"
		_log.Error(message)
		apiSolution.Outcome.Error = true
		apiSolution.Outcome.Message = message
		err := errors.New(message)
		return apiSolution, err
	}

	datasetName := pool + "/asd/" + apiSolution.Name

	dataset, err := zfs.CreateFilesystem(datasetName, nil)
	if err != nil {
		message := "Error creating dataset for new solutions" + datasetName
		_log.Error(message)
		apiSolution.Outcome.Error = true
		apiSolution.Outcome.Message = message
		return apiSolution, err
	} else {
		_log.Info("Succesfully created dataset:" + datasetName)
	}

	// get the actual mountpoint
	mountpoint, err := dataset.GetProperty("mountpoint")
	if err != nil {
		message := "Unable to locate the mountpoint of the master dataset"
		_log.Error(message)
		_log.Error(err)
		apiSolution.Outcome.Message = message
		return apiSolution, err
	} else {
		_log.Info("Got mountpoint:" + mountpoint)
	}

	// hostname
	hostname, err := os.Hostname()
	if err != nil {
		message := "Unable to get master hostname"
		_log.Error(message)
		_log.Error(err)
		apiSolution.Outcome.Message = message
		return apiSolution, err
	} else {
		_log.Info("Got master hostname:" + hostname)
		apiSolution.Hostname = hostname
	}

	// open or create the k/v db
	db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 3 * time.Second})
	if err != nil {
		message := "Unable to open the main db for persisting master info"
		_log.Error(message)
		_log.Error(err)
		apiSolution.Outcome.Message = message
		return apiSolution, err
	} else {
		_log.Info("Main db opened succesfully")
	}
	defer db.Close()

	// add solution info to the k/v db
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("solutions"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		apiSolution.Hostname = hostname
		apiSolution.Status = "available"

		var encoded []byte
		encoded, err = proto.Marshal(apiSolution)
		if err != nil {
			return err
		}
		_log.Debugf("encoded lenght before put: %v\n", len(encoded))

		err = b.Put([]byte(apiSolution.Name), encoded)
		if err != nil {
			return fmt.Errorf("put: %s", err)
		}

		encoded2 := b.Get([]byte(apiSolution.Name))
		_log.Debugf("encoded lenght after get: %v\n", len(encoded2))

		return nil
	})
	if err != nil {
		message := "Unable to update db to persist solution info"
		_log.Error(message)
		_log.Error(err)
		apiSolution.Outcome.Message = message
		return apiSolution, err
	} else {
		_log.Info("Main db updated with master info")
	}

	message := "Succesfully created new solution" + apiSolution.Name
	_log.Info(message)
	apiSolution.Outcome.Error = false
	apiSolution.Outcome.Message = message
	return apiSolution, nil
}

func (s *Server) SolutionDestroy(ctx context.Context, in *api.Solution) (*api.Solution, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Destroy")

	apiSolution := in
	var apiOutcome api.Outcome
	apiSolution.Outcome = &apiOutcome

	var pool string
	pool = viper.GetString("pool")
	if pool == "" {
		message := "The master pool is unconfigured"
		_log.Error(message)
		apiSolution.Outcome.Error = true
		apiSolution.Outcome.Message = message
		err := errors.New(message)
		return apiSolution, err
	}

	var dbPath string
	dbPath = viper.GetString("mountpoint") + "/asd.db"
	if dbPath == "" {
		message := "The master pool is unconfigured (dbpath)"
		_log.Error(message)
		apiSolution.Outcome.Error = true
		apiSolution.Outcome.Message = message
		err := errors.New(message)
		return apiSolution, err
	}

	datasetName := pool + "/asd/" + apiSolution.Name

	dataset, err := zfs.GetDataset(datasetName)
	if err != nil {
		message := "Error getting dataset with the given name:" + datasetName
		_log.Error(message)
		apiSolution.Outcome.Error = true
		apiSolution.Outcome.Message = message
		return apiSolution, err
	} else {
		_log.Info("Succesfully got dataset object for dataset:" + datasetName)
	}

	// destroy recursively solution
	err = dataset.Destroy(zfs.DestroyRecursive)
	if err != nil {
		message := "Error recursively destroying dataset named:" + datasetName
		_log.Error(message)
		apiSolution.Outcome.Error = true
		apiSolution.Outcome.Message = message
		return apiSolution, err
	} else {
		_log.Info("Succesfully destroyed dataset named:" + datasetName)
	}

	// open or create the k/v db
	db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 3 * time.Second})
	if err != nil {
		message := "Unable to open the main db for persisting master info"
		_log.Error(message)
		_log.Error(err)
		apiSolution.Outcome.Message = message
		return apiSolution, err
	} else {
		_log.Info("Main db opened succesfully")
	}
	defer db.Close()

	// add solution info to the k/v db
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("solutions"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		c := b.Cursor()

		//_log.Debug("Entering delete loop...")
		//_log.Debugf("datasetname=%s\n", datasetName)

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			//_log.Debugf("key=%s\n", k)
			if string(k) == apiSolution.Name {
				//_log.Debug("Found")
				err = c.Delete()
				if err != nil {
					return fmt.Errorf("deleting key %s: error: %s", string(k), err)
				}
			}
		}
		return nil
	})
	if err != nil {
		message := "Unable to update db to delete solution info"
		_log.Error(message)
		_log.Error(err)
		apiSolution.Outcome.Message = message
		return apiSolution, err
	} else {
		_log.Info("Main db updated with master info")
	}

	message := "Succesfully destroyed solution:" + apiSolution.Name
	_log.Info(message)
	apiSolution.Outcome.Error = false
	apiSolution.Outcome.Message = message
	return apiSolution, nil
}

func (s *Server) SolutionDeploy(ctx context.Context, in *api.Solution) (*api.Solution, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Deploy")

	apiSolution := in
	var apiOutcomeVal api.Outcome
	apiOutcome := &apiOutcomeVal
	apiSolution.Outcome = apiOutcome

	var pool string
	pool = viper.GetString("pool")
	if pool == "" {
		message := "The master pool is unconfigured"
		_log.Error(message)
		apiOutcome.Error = true
		apiOutcome.Message = message
		err := errors.New(message)
		return apiSolution, err
	}

	var dbPath string
	dbPath = viper.GetString("mountpoint") + "/asd.db"
	if dbPath == "" {
		message := "The master pool is unconfigured (dbpath)"
		_log.Error(message)
		apiOutcome.Error = true
		apiOutcome.Message = message
		err := errors.New(message)
		return apiSolution, err
	}

	//- [ ]  check if the dataset exists

	sourceName := pool + "/asd/" + apiSolution.Name

	sourceDataset, err := zfs.GetDataset(sourceName)
	if err != nil {
		message := "Unable to locate the source dataset"
		_log.Error(message)
		_log.Error(err)
		apiOutcome.Message = message
		apiOutcome.Error = true
		return apiSolution, err
	} else {
		_log.Info("Source dataset found:" + apiSolution.Name)
	}

	//- [ ]  check if the dataset in in the "available" state

	// open or create the k/v db
	db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 3 * time.Second})
	if err != nil {
		message := "Unable to open the main db for persisting master info"
		_log.Error(message)
		_log.Error(err)
		apiSolution.Outcome.Message = message
		return apiSolution, err
	} else {
		_log.Info("Main db opened succesfully")
	}
	defer db.Close()

	var apiSolutionVal api.Solution
	apiSolutionTmp := &apiSolutionVal

	var apiNodeVal api.Node
	apiNode := &apiNodeVal

	// get node info to the k/v db
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("nodes"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		c := b.Cursor()
		var found bool

		for k, v := c.First(); k != nil; k, v = c.Next() {
			if string(k) == apiSolution.Hostname {
				found = true
				err = proto.Unmarshal(v, apiNode)
				if err != nil {
					return fmt.Errorf("unmarshaling node proto: %s", err)
				}
				apiSolution.Ip = apiNode.Ip
				_log.Debugf("node ip=%s\n", apiSolution.Ip)
				apiSolution.Poolname = apiNode.Poolname
				_log.Debugf("node poolname=%s\n", apiSolution.Poolname)
				break
			}
		}

		if !found {
			return fmt.Errorf("unmarshaling node proto: %s", err)
		}
		return nil
	})
	if err != nil {
		message := "Unable to open the main db to get node info"
		_log.Error(message)
		_log.Error(err)
		apiSolution.Outcome.Message = message
		return apiSolution, err
	} else {
		_log.Info("node info obtained")
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("solutions"))

		encodedSolution := b.Get([]byte(apiSolution.Name))

		// solution exists?
		if encodedSolution == nil {
			// no solution found in k/v
			err := errors.New("solution not found")
			return err
		}

		err = proto.Unmarshal(encodedSolution, apiSolutionTmp)
		if err != nil {
			// no solution found in k/v
			err := errors.New("error unmarshaling solution from k/v")
			return err
		}

		if apiSolutionTmp.Status != "available" {
			// impossible to proceed, wrong status
			err := errors.New("solution must be 'available' to be deployed")
			return err
		}

		return nil
	})
	if err != nil {
		message := "unable to get solution info or solution not in 'available' state"
		_log.Error(message)
		_log.Error(err)
		apiSolution.Outcome.Message = message
		return apiSolution, err
	} else {
		_log.Info("got solution info")
	}

	//- [ ]  makes a @deploy snapshot
	sourceDataset, err = sourceDataset.Snapshot("deploy", false)
	if err != nil {
		message := "Unable to snap the source dataset"
		_log.Error(message)
		_log.Error(err)
		apiSolution.Outcome.Message = message
		apiSolution.Outcome.Error = true
		return apiSolution, err
	} else {
		_log.Info("Source dataset snapshot made prior to deploy")
	}

	//- [ ]  sends the deploy snapshot to the worker node via ssh command, effectively creating a new dataset on the worker node

	//func Command(arg ...string) ([][]string, error) {

	cmd := "zfs send " + sourceName + "@deploy | ssh root@" + "asdo.avero.it" + " 'zfs recv -F " + apiSolution.Poolname + "/asd/" + apiSolution.Name + "'"
	//cmd := "zfs send " + sourceName + "@deploy | zfs recv " + sourceName + "@deployed"
	//cmd := "zfs snapshot " + sourceName + "@deployed"
	//cmd := "date"
	_log.Infof("cmd=%s\n", cmd)
	output, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	_log.Notice(string(output))
	if err != nil {
		message := "error executing zfs send to deploy solution"
		_log.Error(message)
		_log.Error(err)
		apiSolution.Outcome.Message = message
		apiSolution.Outcome.Error = true
		return apiSolution, err
	} else {
		_log.Info("zfs send command executed succesfully")
	}

	//lines, err := zfs.ShCommand("/sbin/bash -c", "'zfs send", sourceName+"@deploy", "| ssh root@"+apiSolution.Ip+" 'zfs recv "+apiSolution.Poolname+"/asd/"+apiSolution.Name + "'")
	//lines, err := zfs.ShCommand("/sbin/bash -c", "echo hello")
	//
	//if err != nil {
	//	message := "error executing zfs send to deploy solution"
	//	_log.Error(message)
	//	_log.Error(err)
	//	apiSolution.Outcome.Message = message
	//	apiSolution.Outcome.Error = true
	//	return apiSolution, err
	//} else {
	//	for _, line := range lines {
	//		fmt.Println(line)
	//	}
	//	_log.Info("zfs send command executed succesfully")
	//}

	//- [ ]  update k/v, sets solution as "deployed"
	// add solution info to the k/v db
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("solutions"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		encodedSolution := b.Get([]byte(apiSolution.Name))

		// solution exists?
		if encodedSolution == nil {
			// no solution found in k/v
			err := errors.New("solution not found")
			return err
		}

		err = proto.Unmarshal(encodedSolution, apiSolutionTmp)
		if err != nil {
			// no solution found in k/v
			err := errors.New("error unmarshaling solution from k/v")
			return err
		}

		if apiSolutionTmp.Status != "available" {
			// impossible to proceed, wrong status
			err := errors.New("solution must be 'available' to be deployed")
			return err
		}

		return nil
	})

	return apiSolution, nil
}

func (s *Server) SolutionRetire(ctx context.Context, in *api.Solution) (*api.Solution, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Retreat")

	apiSolution := in

	err := errors.New("not implemented")
	// destroy a solution datased, with ALL backups, snapshots

	err = errors.Wrap(err, "really not")
	return apiSolution, err
}

func (s *Server) SolutionStart(ctx context.Context, in *api.Solution) (*api.Solution, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Start")

	apiSolution := in

	err := errors.New("not implemented")
	// destroy a solution datased, with ALL backups, snapshots

	err = errors.Wrap(err, "really not")
	return apiSolution, err
}

func (s *Server) SolutionStop(ctx context.Context, in *api.Solution) (*api.Solution, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Stop")

	apiSolution := in

	err := errors.New("not implemented")
	// destroy a solution datased, with ALL backups, snapshots

	err = errors.Wrap(err, "really not")
	return apiSolution, err
}

func (s *Server) SolutionSnapshot(ctx context.Context, in *api.Solution) (*api.Solution, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Snapshot")

	apiSolution := in

	err := errors.New("not implemented")
	// destroy a solution datased, with ALL backups, snapshots

	err = errors.Wrap(err, "really not")
	return apiSolution, err
}

func (s *Server) SolutionRollback(ctx context.Context, in *api.Solution) (*api.Solution, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Rollback")

	apiSolution := in

	err := errors.New("not implemented")
	// destroy a solution datased, with ALL backups, snapshots

	err = errors.Wrap(err, "really not")
	return apiSolution, err
}

func (s *Server) SolutionBackup(ctx context.Context, in *api.Solution) (*api.Solution, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Backup")

	apiSolution := in

	err := errors.New("not implemented")
	// destroy a solution datased, with ALL backups, snapshots

	err = errors.Wrap(err, "really not")
	return apiSolution, err
}

func (s *Server) SolutionRestore(ctx context.Context, in *api.Solution) (*api.Solution, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Restore")

	apiSolution := in

	err := errors.New("not implemented")
	// destroy a solution datased, with ALL backups, snapshots

	err = errors.Wrap(err, "really not")
	return apiSolution, err
}

func (s *Server) SolutionExpose(ctx context.Context, in *api.Solution) (*api.Solution, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Expose")

	apiSolution := in

	err := errors.New("not implemented")
	// destroy a solution datased, with ALL backups, snapshots

	err = errors.Wrap(err, "really not")
	return apiSolution, err
}

func (s *Server) SolutionContain(ctx context.Context, in *api.Solution) (*api.Solution, error) {
	_log := helpers.InitLogs(true)
	_log.Debug("gRPC call: Contain")

	apiSolution := in

	err := errors.New("not implemented")
	// destroy a solution datased, with ALL backups, snapshots

	err = errors.Wrap(err, "really not")
	return apiSolution, err
}
