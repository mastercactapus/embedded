# embedded

Collection of tools, drivers, etc... for embedded development.

## OpenOCD

## Quickstart

To debug on a RP2040, you will need openocd compiled from the rpi-openocd repository.

```sh
# deps for debian systems
sudo apt install automake autoconf build-essential texinfo libtool libftdi-dev libusb-1.0-0-dev

# special repo, special branch (or you will get "Error: The specified debug interface was not found (picoprobe)")
git clone https://github.com/raspberrypi/openocd.git --branch rp2040 --depth=1 --no-single-branch

cd openocd
./bootstrap

# add --enable-sysfsgpio --enable-bcm2835gpio for raspberry pi GPIO/bitbang support
./configure --enable-picoprobe
make -j4
sudo make install
```

Serial: `picocom -b 115200 /dev/ttyACM0`
Flash: `tinygo flash -target=pico -programmer=ocd ./cmd/bustool-pico/`