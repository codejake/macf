# macf

This small program reformats 48-bit OUIs, such as Ethernet addresses, in useful
ways for network engineers.

The OUI is normally represented as a set of octets in hexadecimal notation 
sperated by hyphens or colons, but then you have companies, like Microsoft and
Cisco, who do things their own way (`AABBCCDDEEFF` and `aabb.ccdd.eeff` 
respectively)

## Usage

```
Usage: macf [-f/--format <format>] [-u/--uppercase] <ethernet_address>
  Options:

    --format options:
        cisco        => a1b2.c3d4.e5f6
        colon        => a1:b2:c3:d4:e5:f6
        hyphen       => a1-b2-c3-d4-e5-f6
        none         => a1b2c3d4e5f6
```