# Terraform Provider Gravity

Terraform provider for https://gravity.beryju.io

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.18

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

### Running the tests

The acceptance tests run against a real Gravity instance. The devcontainer
brings one up alongside the workspace; to run one directly:

```shell
docker run -d --name gravity \
  -p 8008:8008 -p 8009:8009 -p 8010:8010 \
  -e LISTEN_ONLY=true -e INSTANCE_IP=0.0.0.0 \
  -e ADMIN_TOKEN=test -e ADMIN_PASSWORD=test \
  ghcr.io/beryju/gravity:latest

# Wait for it to come up
curl --retry 30 --retry-all-errors -sf http://localhost:8009/healthz/live
```

Then point the tests at it:

```shell
TF_ACC=1 GRAVITY_URL=http://localhost:8008 GRAVITY_TOKEN=test go test -v ./internal/provider/
```

*Note:* Acceptance tests create real resources, and often cost money to run.

`TestAccResourceRoleDHCP` is skipped unless `GRAVITY_TEST_DHCP_ROLE_CONFIG` is
set. Applying the DHCP role configuration restarts the DHCP role, and on
Gravity 0.33.0 and earlier that leaves every subsequent `gravity_dhcp_scope`
write failing until Gravity is restarted, so it needs an instance it can
render unusable:

```shell
GRAVITY_TEST_DHCP_ROLE_CONFIG=1 TF_ACC=1 GRAVITY_URL=http://localhost:8008 \
  GRAVITY_TOKEN=test go test -v -run TestAccResourceRoleDHCP ./internal/provider/
```
