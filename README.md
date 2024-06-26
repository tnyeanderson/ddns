# ddns

Go package and CLI which provides server and agent components for a simple DDNS
solution. A single executable with no dependencies. Can also be used with
docker or kubernetes.

## How it works

The `ddns server` command starts an HTTP API server and a DNS server which will
serve as authoritative for the domains that are [configured to use
it](#prerequisites). Then another machine somewhere else (usually the backend
server itself) can use the `ddns update` command to set the IP address for
a given domain. Both commands are configured using environment variables.

API keys are used for authentication, and API keys can be restricted to only
update certain domains based on a regex matcher.

## Prerequisites

Feel free to test the program all you want locally, but if you want other
resolvers to start actually handing out your dynamic IP address, you need to
tell them that this DDNS server you are running should be used to resolve your
DDNS domain. This is done by creating NS records.

For example, you own `domain.com` and have already set up A records for it
which point to your public server at `100.100.100.100`, which is running the
DDNS server component. You want to have a subdomain that resolves using DDNS to
lead to your homelab, like `myddns.domain.com`. To handle this, create the
following DNS records for `domain.com`.

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

For `myddns.com`:

| Type | Host   | Value           |
|------|--------|-----------------|
| NS   | @      | ns1.server.com  |
| NS   | @      | ns2.server.com  |

## Installation and configuration

Install the program by downloading the release binary directly, or by running:

```
go install github.com/tnyeanderson/ddns@latest
```

View supported environment variables and their descriptions with:

```
ddns help
```

### Server setup

Using the binary:

```
DDNS_SERVER_API_KEY=createatoken ddns server
```

Using docker:

```
docker run -e DDNS_SERVER_API_KEY=createatoken ghcr.io/tnyeanderson/ddns server
```

### Agent setup

Using the binary:

```
DDNS_API_SERVER=ddns.myserver.site DDNS_API_KEY=createatoken ddns update yourdomain.site 1.2.3.4
```

Using docker:

```
docker run -e DDNS_API_SERVER=ddns.myserver.site -e DDNS_API_KEY=createatoken ghcr.io/tnyeanderson/ddns update yourdomain.site 1.2.3.4
```

> NOTE: To update the IP address to the public IP of the box making the
request, simply omit the IP argument and it will be calculated automatically by
the API server.

Updating an IP can also be done directly with `curl`:

```
curl -X POST -H "Authorization: Bearer $DDNS_API_KEY" 'yourserver.com/api/v1/update?domain=yourdomain.site&ip=1.2.3.4'
```

> NOTE: To update the IP address to the public IP of the box making the
request, set the IP parameter to "auto" and it will be calculated automatically
by the API server.

