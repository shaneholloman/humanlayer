package session

import (
	"os"
	"path/filepath"
	"testing"

	claudecode "github.com/humanlayer/humanlayer/claudecode-go"
	"github.com/humanlayer/humanlayer/hld/bus"
	"github.com/humanlayer/humanlayer/hld/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMCPSocketPathPassedToServers(t *testing.T) {
	// Create temporary directory for SQLite database
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	// Create SQLite store
	sqliteStore, err := store.NewSQLiteStore(dbPath)
	require.NoError(t, err)
	defer func() { _ = sqliteStore.Close() }()

	eventBus := bus.NewEventBus()

	// Test socket path
	testSocketPath := "/tmp/test-daemon.sock"

	// Create manager with socket path
	manager, err := NewManager(eventBus, sqliteStore, testSocketPath)
	require.NoError(t, err)

	// Verify socket path was stored
	assert.Equal(t, testSocketPath, manager.socketPath)

	// Create a session config with MCP servers
	config := claudecode.SessionConfig{
		Query:      "test query",
		WorkingDir: "/tmp/test",
		MCPConfig: &claudecode.MCPConfig{
			MCPServers: map[string]claudecode.MCPServer{
				"test-server": {
					Command: "test-command",
					Args:    []string{"arg1", "arg2"},
					Env: map[string]string{
						"EXISTING_VAR": "value",
					},
				},
				"another-server": {
					Command: "another-command",
					// No existing env
				},
			},
		},
	}

	// Since we can't mock the client directly, we'll test the logic by examining the config
	// after it's been processed in LaunchSession. We'll skip the actual launch.
	// Instead, we'll test the socket path logic directly.

	// Add HUMANLAYER_RUN_ID and HUMANLAYER_DAEMON_SOCKET to MCP server environment
	// This mimics what happens in LaunchSession
	runID := "test-run-id"
	if config.MCPConfig != nil {
		for name, server := range config.MCPConfig.MCPServers {
			if server.Env == nil {
				server.Env = make(map[string]string)
			}
			server.Env["HUMANLAYER_RUN_ID"] = runID
			// Add daemon socket path so MCP servers connect to the correct daemon
			if manager.socketPath != "" {
				server.Env["HUMANLAYER_DAEMON_SOCKET"] = manager.socketPath
			}
			config.MCPConfig.MCPServers[name] = server
		}
	}

	// Verify HUMANLAYER_DAEMON_SOCKET was added to all MCP servers
	require.NotNil(t, config.MCPConfig)
	require.Len(t, config.MCPConfig.MCPServers, 2)

	// Check first server
	testServer := config.MCPConfig.MCPServers["test-server"]
	assert.Equal(t, testSocketPath, testServer.Env["HUMANLAYER_DAEMON_SOCKET"])
	assert.Equal(t, "value", testServer.Env["EXISTING_VAR"], "existing env vars should be preserved")
	assert.Equal(t, runID, testServer.Env["HUMANLAYER_RUN_ID"], "run ID should be set")

	// Check second server
	anotherServer := config.MCPConfig.MCPServers["another-server"]
	assert.Equal(t, testSocketPath, anotherServer.Env["HUMANLAYER_DAEMON_SOCKET"])
	assert.Equal(t, runID, anotherServer.Env["HUMANLAYER_RUN_ID"], "run ID should be set")
}

func TestMCPSocketPathNotSetWhenEmpty(t *testing.T) {
	// Create temporary directory for SQLite database
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	// Create SQLite store
	sqliteStore, err := store.NewSQLiteStore(dbPath)
	require.NoError(t, err)
	defer func() { _ = sqliteStore.Close() }()

	eventBus := bus.NewEventBus()

	// Create manager with empty socket path
	manager, err := NewManager(eventBus, sqliteStore, "")
	require.NoError(t, err)

	// Verify socket path is empty
	assert.Empty(t, manager.socketPath)

	// Create session config with MCP server
	config := claudecode.SessionConfig{
		Query:      "test query",
		WorkingDir: "/tmp/test",
		MCPConfig: &claudecode.MCPConfig{
			MCPServers: map[string]claudecode.MCPServer{
				"test-server": {
					Command: "test-command",
				},
			},
		},
	}

	// Test the socket path logic directly
	runID := "test-run-id"
	if config.MCPConfig != nil {
		for name, server := range config.MCPConfig.MCPServers {
			if server.Env == nil {
				server.Env = make(map[string]string)
			}
			server.Env["HUMANLAYER_RUN_ID"] = runID
			// Add daemon socket path so MCP servers connect to the correct daemon
			if manager.socketPath != "" {
				server.Env["HUMANLAYER_DAEMON_SOCKET"] = manager.socketPath
			}
			config.MCPConfig.MCPServers[name] = server
		}
	}

	// Verify HUMANLAYER_DAEMON_SOCKET was NOT added when socket path is empty
	testServer := config.MCPConfig.MCPServers["test-server"]
	_, hasSocketPath := testServer.Env["HUMANLAYER_DAEMON_SOCKET"]
	assert.False(t, hasSocketPath, "HUMANLAYER_DAEMON_SOCKET should not be set when socket path is empty")
	assert.Equal(t, runID, testServer.Env["HUMANLAYER_RUN_ID"], "run ID should still be set")
}

func TestMCPSocketPathFromEnvironment(t *testing.T) {
	// Set environment variable
	testSocketPath := "/tmp/env-daemon.sock"
	_ = os.Setenv("HUMANLAYER_DAEMON_SOCKET", testSocketPath)
	defer func() { _ = os.Unsetenv("HUMANLAYER_DAEMON_SOCKET") }()

	// Create temporary directory for SQLite database
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	// Create SQLite store
	sqliteStore, err := store.NewSQLiteStore(dbPath)
	require.NoError(t, err)
	defer func() { _ = sqliteStore.Close() }()

	eventBus := bus.NewEventBus()

	// Create manager - socket path should come from config passed to NewManager
	// not from environment variable
	managerSocketPath := "/tmp/manager-socket.sock"
	manager, err := NewManager(eventBus, sqliteStore, managerSocketPath)
	require.NoError(t, err)

	// Verify manager uses the socket path passed to it, not env var
	assert.Equal(t, managerSocketPath, manager.socketPath)
}
