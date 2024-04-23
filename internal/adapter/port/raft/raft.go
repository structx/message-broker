package raft

import (
	"fmt"
	"os"
	"path/filepath"

	transport "github.com/Jille/raft-grpc-transport"
	"github.com/hashicorp/raft"
	boltdb "github.com/hashicorp/raft-boltdb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// New
func New(fsm raft.FSM) (*raft.Raft, *transport.Manager, error) {

	c := raft.DefaultConfig()
	c.LocalID = raft.ServerID(config.Raft.LocalID)

	baseDir := filepath.Join(config.Raft.BaseDir, config.Raft.LocalID)

	logStore, err := boltdb.NewBoltStore(filepath.Join(baseDir, "logs.dat"))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create bolt store 1 %v", err)
	}

	stableStore, err := boltdb.NewBoltStore(filepath.Join(baseDir, "stable.dat"))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create bolt store 2 %v", err)
	}

	snapStore, err := raft.NewFileSnapshotStore(baseDir, 3, os.Stderr)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create snapstor store %v", err)
	}

	tm := transport.New(raft.ServerAddress(config.Server.ListenAddr),
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})

	r, err := raft.NewRaft(
		c,
		fsm,
		logStore,
		stableStore,
		snapStore,
		tm.Transport(),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create new raft %v", err)
	}

	if config.Raft.Bootstrap {
		cfg := raft.Configuration{
			Servers: []raft.Server{
				{
					Suffrage: raft.Voter,
					ID:       raft.ServerID(config.Raft.LocalID),
					Address:  raft.ServerAddress(config.Server.ListenAddr),
				},
			},
		}

		fut := r.BootstrapCluster(cfg)
		err := fut.Error()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to bootstrap raft %v", err)
		}
	}

	return r, tm, nil
}
