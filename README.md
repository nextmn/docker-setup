# Docker-setup
Docker-setup is a program that allow configure a container (for networking) via environment variables.

## Usage
This program only use environment variables for its configuration:
- `ONESHOT`: when this environment variable is equal to `true`, the program will *not* sleep until a signal (SIGINT or SIGTERM) is received, and not perform cleaning scripts on exit
- `NAT4_IFACES` is a list of interfaces where MASQUERADE will be enabled
- `ROUTES_INIT` is a list of routes modifications that will be performed on init
- `ROUTES_EXIT` is a list of routes modifications that will be performed on exit
- `PRE_INIT_HOOK` is a command that is run before the init (nat & routes init). The command can takes some arguments from `PRE_INIT_HOOK_0`, `PRE_INIT_HOOK_1`, and so on.
- `POST_INIT_HOOK` is a command that is run after the init (nat & routes init). The command can takes some arguments from `POST_INIT_HOOK_0`, `POST_INIT_HOOK_1`, and so on.
- `PRE_EXIT_HOOK` is a command that is run before the exit (nat & routes cleaning). The command can takes some arguments from `PRE_EXIT_HOOK_0`, `PRE_EXIT_HOOK_1`, and so on.
- `POST_EXIT_HOOK` is a command that is run after the exit (nat & routes cleaning). The command can takes some arguments from `POST_EXIT_HOOK_0`, `POST_EXIT_HOOK_1`, and so on.

### Example
In Docker Compose:
```yaml
volumes:
    - "./config_init.sh:/usr/local/bin/config_init.sh:ro"
    - "./config_exit.sh:/usr/local/bin/config_exit.sh:ro"
environment:
    ONESHOT: "false"
    ROUTES_INIT: |-
        - add 10.0.1.0/24 via 10.0.0.2
        - add 10.0.2.0/24 via 10.0.0.3
    ROUTES_EXIT: |-
        - del 10.0.1.0/24
        - del 10.0.2.0/24
    NAT4_IFACES: |-
        - eth0
        - eth1
    PRE_INIT_HOOK: config_init.sh
    PRE_INIT_HOOK_0: "Hello"
    PRE_INIT_HOOK_1: "World"
    PRE_INIT_HOOK_2: "!"
    PRE_EXIT_HOOK: config_exit.sh
    PRE_EXIT_HOOK_0: "Goodbye"
```

## Getting started
### Runtime dependencies
- iproute2
- iptables

### Build
Run `go build`

### Docker
- The container requires the `NET_ADMIN` capability;

This can be done in `docker-compose.yaml` by defining the following for the service:

```yaml
cap_add:
    - NET_ADMIN
```

## Author
Louis Royer

## License
MIT
