package command

import (
	"testing"

	"github.com/hashicorp/nomad/ci"
	"github.com/hashicorp/nomad/command/agent"
	"github.com/hashicorp/nomad/nomad/structs"
	"github.com/hashicorp/nomad/testutil"
	"github.com/mitchellh/cli"
	"github.com/shoenig/test/must"
)

func TestLoginCommand_Run(t *testing.T) {
	ci.Parallel(t)

	// Build a test server with ACLs enabled.
	srv, _, agentURL := testServer(t, false, func(c *agent.Config) {
		c.ACL.Enabled = true
	})
	defer srv.Shutdown()

	// Wait for the server to start fully.
	testutil.WaitForLeader(t, srv.Agent.RPC)

	ui := cli.NewMockUi()
	cmd := &LoginCommand{
		Meta: Meta{
			Ui:          ui,
			flagAddress: agentURL,
		},
	}

	// Store a default auth method
	state := srv.Agent.Server().State()
	method := &structs.ACLAuthMethod{
		Name:    "test-auth-method",
		Default: true,
		Type:    "JWT",
		Config: &structs.ACLAuthMethodConfig{
			OIDCDiscoveryURL: "http://example.com",
		},
	}
	method.SetHash()
	must.NoError(t, state.UpsertACLAuthMethods(1000, []*structs.ACLAuthMethod{method}))

	// Test the basic validation on the command.
	must.Eq(t, 1, cmd.Run([]string{"-address=" + agentURL, "this-command-does-not-take-args"}))
	must.StrContains(t, ui.ErrorWriter.String(), "This command takes no arguments")

	ui.OutputWriter.Reset()
	ui.ErrorWriter.Reset()

	// Specify an incorrect type of default method
	must.Eq(t, 1, cmd.Run([]string{"-address=" + agentURL, "-type=OIDC"}))
	must.StrContains(t, ui.ErrorWriter.String(), "Specified type: OIDC does not match the type of the default method: JWT")

	ui.OutputWriter.Reset()
	ui.ErrorWriter.Reset()

	// Try logging in with non-OIDC method and no token (expected error)
	must.Eq(t, 1, cmd.Run([]string{"-address=" + agentURL, "-type=JWT"}))
	must.StrContains(t, ui.ErrorWriter.String(), "You need to provide a bearer token.")

	ui.OutputWriter.Reset()
	ui.ErrorWriter.Reset()

	// TODO(jrasell) find a way to test the full login flow from the CLI
	//  perspective.
}
