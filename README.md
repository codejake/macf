# macf

This small program reformats 48-bit OUIs, such as Ethernet addresses, in useful
ways for network engineers.

The OUI is normally represented as a set of octets in hexadecimal notation
separated by hyphens or colons, but some tools and vendors also use plain hex
or Cisco-style dotted notation (`AABBCCDDEEFF` and `aabb.ccdd.eeff`).

## Usage

```
Usage: macf [-f format] [-u] <ethernet_address>

Formats:
  cisco   a1b2.c3d4.e5f6
  colon   a1:b2:c3:d4:e5:f6
  hyphen  a1-b2-c3-d4-e5-f6
  none    a1b2c3d4e5f6
```

## Releases

GitHub Releases are built automatically with GoReleaser when a tag like
`v0.1.0` is pushed.

Each release publishes archives for:

- macOS arm64
- Linux amd64
- Linux arm64
- Windows arm64

To cut a release:

```sh
git tag v0.1.0
git push origin v0.1.0
```

For a local dry run, if you have GoReleaser installed:

```sh
goreleaser release --snapshot --clean
```
