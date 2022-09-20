# DDNS

Provides DDNS services for a single domain.

## Installation and setup

### Agent setup

Generate an SSH keypair for communication between the agent and server:
```bash
mkdir agent/conf
ssh-keygen -t ed25519 -C ddns -f agent/conf/ssh.key -N ''
```

> NOTE: Stop here and make sure the server component is running before
continuing!

Set the `SSH_HOST` and `SSH_PORT` environment variables in `agent/ddns.env`.
Use `agent/ddns.env.example` as a template/guide.

Build and run the agent to make sure it is working:
```bash
cd agent
docker-compose build && docker-compose up
```

Create a cron entry to run the agent regularly with `/etc/cron.d/ddns`:
```cron
# Run the DDNS agent container every 4 minutes
*/4 * * * * root cd /app/ddns/agent && /usr/bin/docker-compose up
```

### Server setup

Create `server/conf` and copy `agent/conf/ssh.key.pub` from the previous step
into it.

Set the `DDNS_DOMAIN` environment variable in `server/ddns.env`. Use
`server/ddns.env.example` as a template/guide.

Build and start the container:
```bash
cd server
docker-compose build && docker-compose up -d
```

## How it works

For this example, **lanhost** is a server or computer running on someone's home
network (or another environment that requires DDNS services) and **dnshost** is
the DDNS server which faces the internet.

The **dnshost** runs the **server** component, which consists of an SSH server
(sshd) and a DNS server (dnsmasq). The **lanhost** runs the **agent** component
with a cronjob. Each time the agent is run, it makes a request to icanhazip.com
to obtain the public IP for **lanhost**. For a home server, this would be the
WAN address, usually assigned by the ISPs DHCP service. It checks the newly
received address against the last one it received, and if it has changed (or
`$FORCE` is set), the agent makes an SSH request to the server component which
updates the DNS entry.

The SSH server is locked down to ignore/overwrite the command provided by the
agent (`$0`) other than its first parameter (the new IP).

The DNS server is locked down to *only* respond to the DDNS domain. It does not
recurse or refer to any upstream nameserver for a result. This prevents it from
being used as a "normal" DNS server (as it is only useful if trying to resolve
the DDNS domain).

## Future features

- [ ] Multiple DDNS domains
- [ ] Different services to get public IP
- [ ] Proper host key verification

