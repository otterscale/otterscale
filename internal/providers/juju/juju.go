package juju

import (
	"context"
	"fmt"
	"sync"
	"time"

	"connectrpc.com/connect"
	"github.com/juju/juju/api"
	"github.com/juju/juju/api/client/action"
	"github.com/juju/juju/api/client/application"
	"github.com/juju/juju/api/client/modelmanager"
	"github.com/juju/juju/api/connector"
	"github.com/juju/names/v5"

	"github.com/otterscale/otterscale/internal/config"
)

type Juju struct {
	conf *config.Config

	connections sync.Map
}

func New(conf *config.Config) *Juju {
	return &Juju{
		conf: conf,
	}
}

func (m *Juju) newConnection(uuid string) (api.Connection, error) {
	juju := m.conf.Juju
	opts := connector.SimpleConfig{
		ModelUUID:           uuid,
		ControllerAddresses: juju.ControllerAddresses,
		Username:            juju.Username,
		Password:            juju.Password,
		CACert:              juju.CACert,
	}

	sc, err := connector.NewSimple(opts)
	if err != nil {
		return nil, err
	}
	return sc.Connect()
}

func (m *Juju) getUUID(scope string) (string, error) {
	if scope == "controller" {
		return "", nil
	}

	client, err := m.newConnection("")
	if err != nil {
		return "", err
	}
	defer client.Close()

	models, err := modelmanager.NewClient(client).ListModels(m.username())
	if err != nil {
		return "", err
	}

	for _, model := range models {
		if model.Name == scope {
			return model.UUID, nil
		}
	}

	return "", connect.NewError(connect.CodeNotFound, fmt.Errorf("scope %q not found", scope))
}

func (m *Juju) connection(scope string) (api.Connection, error) {
	if v, ok := m.connections.Load(scope); ok {
		conn := v.(api.Connection)
		if !conn.IsBroken() {
			return conn, nil
		}
		conn.Close()
	}

	uuid, err := m.getUUID(scope)
	if err != nil {
		return nil, err
	}

	conn, err := m.newConnection(uuid)
	if err != nil {
		return nil, err
	}

	m.connections.Store(scope, conn)

	return conn, nil
}

func (m *Juju) waitForCompleted(ctx context.Context, scope, id string, tickInterval, timeoutDuration time.Duration) (map[string]any, error) {
	ticker := time.NewTicker(tickInterval)
	defer ticker.Stop()

	timeout := time.After(timeoutDuration)

	for {
		select {
		case <-ticker.C:
			conn, err := m.connection(scope)
			if err != nil {
				return nil, err
			}

			results, err := action.NewClient(conn).Actions([]string{id})
			if err != nil {
				return nil, err
			}

			if len(results) == 0 || results[0].Action == nil {
				return nil, fmt.Errorf("failed to get action result %q", id)
			}

			if results[0].Status == "completed" {
				return results[0].Output, nil
			}

			continue

		case <-timeout:
			return nil, fmt.Errorf("timeout waiting for action %s to become completed", id)

		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func (m *Juju) Run(ctx context.Context, scope, appName, actionName string, params map[string]any) (map[string]any, error) {
	conn, err := m.connection(scope)
	if err != nil {
		return nil, err
	}

	leader, err := application.NewClient(conn).Leader(appName)
	if err != nil {
		return nil, err
	}

	actions := []action.Action{
		{
			Receiver:   names.NewUnitTag(leader).String(),
			Name:       actionName,
			Parameters: params,
		},
	}

	enqueued, err := action.NewClient(conn).EnqueueOperation(actions)
	if err != nil {
		return nil, err
	}

	if len(enqueued.Actions) == 0 || enqueued.Actions[0].Action == nil {
		return nil, fmt.Errorf("failed to run action %q on %s", actionName, leader)
	}

	id := enqueued.Actions[0].Action.ID

	return m.waitForCompleted(ctx, scope, id, time.Second, time.Minute)
}

func (m *Juju) Execute(ctx context.Context, scope, appName, command string) (map[string]any, error) {
	conn, err := m.connection(scope)
	if err != nil {
		return nil, err
	}

	leader, err := application.NewClient(conn).Leader(appName)
	if err != nil {
		return nil, err
	}

	parallel := true
	args := action.RunParams{
		Commands: command,
		Units:    []string{leader},
		Parallel: &parallel,
	}

	enqueued, err := action.NewClient(conn).Run(args)
	if err != nil {
		return nil, err
	}

	if len(enqueued.Actions) == 0 || enqueued.Actions[0].Action == nil {
		return nil, fmt.Errorf("failed to run command %q on %s", command, leader)
	}

	id := enqueued.Actions[0].Action.ID

	return m.waitForCompleted(ctx, scope, id, time.Second, time.Minute)
}

func (m *Juju) GetEndpoint(_ context.Context, scope, appName string) (string, error) {
	conn, err := m.connection(scope)
	if err != nil {
		return "", err
	}

	leader, err := application.NewClient(conn).Leader(appName)
	if err != nil {
		return "", err
	}

	tag := names.NewUnitTag(leader)
	unitTags := []names.UnitTag{
		tag,
	}

	units, err := application.NewClient(conn).UnitsInfo(unitTags)
	if err != nil {
		return "", err
	}

	for i := range units {
		if units[i].Tag == tag.String() {
			return units[i].PublicAddress, nil
		}
	}

	return "", connect.NewError(connect.CodeNotFound, fmt.Errorf("endpoint %q not found", appName))
}

func (m *Juju) controller() string {
	return m.conf.Juju.Controller
}

func (m *Juju) username() string {
	return m.conf.Juju.Username
}

func (m *Juju) cloudName() string {
	return m.conf.Juju.CloudName
}

func (m *Juju) cloudRegion() string {
	return m.conf.Juju.CloudRegion
}

func (m *Juju) charmhubAPIURL() string {
	return m.conf.Juju.CharmhubAPIURL
}
