# DDNS

Provides DDNS services for a single domain.

## Prerequisites

Create NS records for the domain to be controlled by DDNS. For example, if the
domain is `myddns.domain.com`, and the DDNS server component is running on
`domain.com` (which resolves to `100.100.100.100`), then create the following
DNS records for `domain.com`:

| Type | Host   | Value           |
|------|--------|-----------------|
| A    | ns1    | 100.100.100.100 |
| A    | ns2    | 100.100.100.100 |
| NS   | myddns | ns1.domain.com  |
| NS   | myddns | ns2.domain.com  |

The key here is at a minimum of two NS records is *required*. To accomodate
this, two A records were created that point to the host on which the DDNS
server component is running. As a different example, let's suppose that
`myddns.com` is the DDNS domain, and the server is running at `server.com`
(same IP as before). Then, the DDNS records would need to be edited for both
domains.

For `server.com`:

| Type | Host   | Value           |
|------|--------|-----------------|
| A    | ns1    | 100.100.100.100 |
| A    | ns2    | 100.100.100.100 |

For `myddns.com`

| Type | Host   | Value           |
|------|--------|-----------------|
| NS   | myddns | ns1.server.com  |
| NS   | myddns | ns2.server.com  |

## Installation and setup

### Server setup

To start quickly:

```
docker run -e DDNS_SERVER_API_KEY=createatoken ddns server
```

See `server.env.example` for more options.

### Agent setup

To start quickly:

```
docker run -e DDNS_API_SERVER=ddns.myserver.site -e DDNS_API_KEY=createatoken ddns agent yourdomain.site 1.2.3.4
```

To update the IP address to the public IP of the box running the agent, simply
omit the IP and it will be calculated automatically.

See `agent.env.example` for more options.

