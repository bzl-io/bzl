package plugin

import (
	"os"
	"path"
	"log"
	"plugin"
	"google.golang.org/grpc"
	"github.com/bzl-io/bzl/api"
	"github.com/bzl-io/bzl/command"
	"github.com/bzl-io/bzl/config"
)

type Manager struct {
	conn *grpc.ClientConn
	Client *api.PluginApiClient
}

// Given the name of a plugin, return true if the plugin
// already exists in the cache dir.
func (p *Manager) HasPlugin(name string) bool {
	log.Println("Plugin not exists", name)
	return false
}

// Given the name of a plugin, return it as a command
func (m *Manager) GetCommandPlugin(name string) (command.Command, error) {
	home, err := config.GetHome()
	if err != nil {
		return nil, err
	}
	
	filename := path.Join(home, "plugin", name + ".so")
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Println("Plugin not found", name)
		return nil, nil // plugin not exists
	}

	p, err := plugin.Open(filename)
	if err != nil {
		return nil, err
	}

	cmd, err := p.Lookup("Execute")
	if err != nil {
		return nil, err
	}
	
	return cmd.(command.Command), nil
}

func (m *Manager) Dispose() {
	if m.conn != nil {
		m.conn.Close()
		m.conn = nil
	}
}

func NewManager(address string) (*Manager, error) {
	if address == "" {
		address = "localhost:6060"
	}
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	client := api.NewPluginApiClient(conn)
	return &Manager{
		conn: conn,
		Client: &client,
	}, nil
}
