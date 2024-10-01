# macf

This small program reformats 48-bit OUIs, such as Ethernet addresses, in useful
ways for network engineers.

The OUI is normally represented as a set of octets in hexadecimal notation 
sperated by hyphens or colons, but then you have companies, like Microsoft and
Cisco, who do things their own way (`AABBCCDDEEFF` and `aabb.ccdd.eeff` 
respectively)