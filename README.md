# wold
simply Wake On Lan Daemon

# API

|  |URL|BODY|note|
|--|---|----|----|
|POST|/wake|{"mac": "12:34:56:78:90:af"}|Wakeup the machine having the MAC address.|

# Install

# Run
./wold --listen-address=127.0.0.1 --broad-cast-address=10.0.0.255 --iface eth0