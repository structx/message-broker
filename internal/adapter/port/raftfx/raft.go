// Package raftfx hashicorp raft provider
package raftfx

import (
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"

	transport "github.com/Jille/raft-grpc-transport"
	"github.com/hashicorp/raft"
	boltdb "github.com/hashicorp/raft-boltdb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/trevatk/mora/internal/adapter/setup"
)

// New constructor
func New(config *setup.Config, fsm raft.FSM) (*raft.Raft, *transport.Manager, error) {

	c := raft.DefaultConfig()
	c.LocalID = raft.ServerID(config.Raft.LocalID)

	baseDir, err := mkdirs(config.Raft.BaseDir, config.Raft.LocalID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create directory and files %v", err)
	}

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

	addr := net.JoinHostPort(
		config.Server.BindAddr,
		fmt.Sprintf("%d", config.Server.Ports.GRPC),
	)

	tm := transport.New(raft.ServerAddress(addr),
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
					Address:  raft.ServerAddress(addr),
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

func mkdirs(baseDir, localID string) (string, error) {

	nd := filepath.Join(baseDir, localID)
	err := os.Mkdir(filepath.Clean(nd), os.ModePerm)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return "", fmt.Errorf("failed to create directory %v", err)
		}
	}

	fp1 := filepath.Join(nd, "logs.dat")
	f1, err := os.Create(filepath.Clean(fp1))
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return "", fmt.Errorf("failed to create file %v", err)
		}
	}
	defer func() { _ = f1.Close() }()

	fp2 := filepath.Join(nd, "stable.dat")
	f2, err := os.Create(filepath.Clean(fp2))
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return "", fmt.Errorf("failed to create file %v", err)
		}
	}
	defer func() { _ = f2.Close() }()

	return nd, nil
}
