# macf: MAC Address Formatter

This small program reformats 48-bit OUIs, such as Ethernet addresses, in ways that are 
useful for network engineers. Tools, vendors, and individuals have different preferred 
formats.

## Usage

```
Usage: macf [-f format] [-u] [ethernet_address]

If no ethernet_address is provided, macf reads one from stdin.

Formats:
  cisco   a1b2.c3d4.e5f6
  colon   a1:b2:c3:d4:e5:f6
  hyphen  a1-b2-c3-d4-e5-f6
  none    a1b2c3d4e5f6
```

Examples:

```sh
# Show all common formats for a MAC address passed as an argument
macf c6:89:f2:d2:dc:3e

# Produce a single output format
macf -f none c6:89:f2:d2:dc:3e

# Read the MAC address from stdin
echo c6:89:f2:d2:dc:3e | macf

# Combine stdin input with formatting flags
echo c6:89:f2:d2:dc:3e | macf -f hyphen -u
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
