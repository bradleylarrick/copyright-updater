# How to use:
#
# This example Makefile system is a template for building applications that
# use libsidekiq. The user is expected to set either PLATFORM or BUILD_CONFIG
# to match the target where the application should run. Based on PLATFORM or
# BUILD_CONFIG, this file will select a matching library (LIB_CONFIG) and a
# toolchain (CROSS_COMPILE). The selected toolchain is only intended to be a
# known working example.
#
# The libraries named after a Matchstiq product are intended to be compatible
# with the root filesystem provided for that product. The libraries named with
# x86_64, x86_32, or aarch64 are intended to run on a Linux OS with a minimum
# GLIBC/GLIBC++ runtime installed and support all Sidekiq radios that are able
# to interface with that host system (i.e. Sidekiq M.2 or NV100 over PCIe).
#
# A shared-library libsidekiq is available for aarch64, Matchstiq G20/G40, and
# Matchstiq X40. Set SO_LINK=enabled to link using libsidekiq.so.
#
# For more details about using the Sidekiq SDK, see the Sidekiq Software
# Development Manual section "Target Devices".
#
# Target Device     | BUILD_CONFIG/PLATFORM         | Library
# ==================|===============================|==========================
# x86-64            | BUILD_CONFIG=x86_64.gcc       | libsidekiq__x86_64.gcc.a
# x86-32            | BUILD_CONFIG=x86_32.gcc       | libsidekiq__x86_32.gcc.a
# aarch64           | BUILD_CONFIG=aarch64          | libsidekiq__aarch64.a
# Matchstiq S1x/S2x | PLATFORM=msiq-sx              | libsidekiq__msiq-sx.a
# Sidekiq Z2        | PLATFORM=z2-armhf             | libsidekiq__z2-armhf.a
# Matchstiq Z3u     | PLATFORM=z3u                  | libsidekiq__z3u.a
# Matchstiq G20/G40 | PLATFORM=msiq-g20g40          | libsidekiq__msiq-g20g40.a
# Matchstiq X40     | PLATFORM=msiq-x40             | libsidekiq__msiq-x40.a
#
# Note: Target device "aarch64" is only for general aarch64 devices with a
#       Sidekiq card (i.e. NVIDIA Jetson with Sidekiq NV100) and not for
#       Matchstiq Z3u/G20/G40/X40.
#
#
# The following BUILD_CONFIG/PLATFORM are deprecated as of v4.21.0:
#
#   * BUILD_CONFIG=aarch64.gcc7.3.1
#     - Use BUILD_CONFIG=aarch64 instead
#
#   * BUILD_CONFIG=aarch64.native
#     - This was an undocumented option, only intended for Z3u
#     - Use PLATFORM=z3u instead
#
# The following BUILD_CONFIG/PLATFORM are removed as of v4.20.0:
#
#   * BUILD_CONFIG=arm_cortex-a9.gcc4.8_gnueabihf_linux
#     - Use PLATFORM=z2-armhf instead (see note for PLATFORM=z2)
#
#   * BUILD_CONFIG=arm_cortex-a9.gcc4.9.2_gnueabi
#     - Use PLATFORM=z2-armhf instead (see note for PLATFORM=z2)
#
#   * PLATFORM=z2
#     - Old platform/library for Z2 BSP v2.0.0 and earlier
#     - Z2 hardware should be upgraded to BSP > v3.0.0 and use PLATFORM=z2-armhf
#
# The following BUILD_CONFIG/PLATFORM are deprecated as of v4.19.0:
#
#   * BUILD_CONFIG=arm_cortex-a9.gcc5.2_glibc_openwrt
#     - Use PLATFORM=msiq-sx instead
#
#   * BUILD_CONFIG=arm_cortex-a9.gcc7.2.1_gnueabihf
#     - Use PLATFORM=z2-armhf instead
#
#   * BUILD_CONFIG=aarch64.gcc6.3
#     - Use PLATFORM=z3u instead
#
#   * BUILD_CONFIG=x86_32.gcc
#     - No longer supported and may be removed in future release

define ALLOWED_BUILD_CONFIGS
    aarch64.gcc6.3
    aarch64
    x86_32.gcc
    x86_64.gcc
    arm_cortex-a9.gcc5.2_glibc_openwrt
    arm_cortex-a9.gcc7.2.1_gnueabihf
endef

define ALLOWED_PLATFORMS
    msiq-g20g40
    msiq-sx
    msiq-x40
    z2-armhf
    z3u
endef

BUILD_CONFIG?= aarch64
PLATFORM?= msiq-g20g40

# BUILD_MACHINE is used to set CROSS_COMPILE to a default value when building
# on X86_64 for an ARM target
BUILD_MACHINE ?= $(shell uname -m)

# removed PLATFORM
ifeq ($(PLATFORM),z2)
  $(error PLATFORM=z2 is no longer supported, upgrade Z2 hardware to BSP v3.0.0 or later and use PLATFORM=z2-armhf)
endif

# removed BUILD_CONFIG
ifeq ($(BUILD_CONFIG),arm_cortex-a9.gcc4.8_gnueabihf_linux)
  $(error BUILD_CONFIG=arm_cortex-a9.gcc4.8_gnueabihf_linux is no longer supported, upgrade Z2 hardware to BSP v3.0.0 or later and use PLATFORM=z2-armhf)
else ifeq ($(BUILD_CONFIG),arm_cortex-a9.gcc4.9.2_gnueabi)
  $(error BUILD_CONFIG=arm_cortex-a9.gcc4.9.2_gnueabi is no longer supported, upgrade Z2 hardware to BSP v3.0.0 or later and use PLATFORM=z2-armhf)
endif

# deprecated BUILD_CONFIG
ifeq ($(BUILD_CONFIG),arm_cortex-a9.gcc5.2_glibc_openwrt)
  $(warning BUILD_CONFIG=arm_cortex-a9.gcc5.2_glibc_openwrt is deprecated and may be removed, use PLATFORM=msiq-sx)
else ifeq ($(BUILD_CONFIG),arm_cortex-a9.gcc7.2.1_gnueabihf)
  $(warning BUILD_CONFIG=arm_cortex-a9.gcc7.2.1_gnueabihf is deprecated and may be removed, use PLATFORM=z2-armhf)
else ifeq ($(BUILD_CONFIG),aarch64.gcc6.3)
  $(warning BUILD_CONFIG=aarch64.gcc6.3 is deprecated and may be removed, use PLATFORM=z3u)
else ifeq ($(BUILD_CONFIG),x86_32.gcc)
  $(warning BUILD_CONFIG=x86_32.gcc is no longer supported and may be removed)
else ifeq ($(BUILD_CONFIG),aarch64.gcc7.3.1)
  $(warning BUILD_CONFIG=aarch64.gcc7.3.1 is deprecated, use PLATFORM=aarch64)
  override BUILD_CONFIG := aarch64
else ifeq ($(BUILD_CONFIG),aarch64.native)
  $(warning BUILD_CONFIG=aarch64.native is deprecated, use PLATFORM=z3u for Matchstiq Z3u or BUILD_CONFIG=aarch64 for aarch64 with Sidekiq card)
  override BUILD_CONFIG := aarch64.gcc6.3
endif

# check for legacy variables TOOL_DIR / CROSS_PREFIX
ifneq (,$(TOOL_DIR))
    $(warning TOOL_DIR is deprecated, use CROSS_COMPILE to specify a toolchain)
endif
ifneq (,$(CROSS_PREFIX))
    $(warning CROSS_PREFIX is deprecated, use CROSS_COMPILE to specify a toolchain)
endif

# CFLAGS common to all build configurations
CFLAGS+= -fstrict-aliasing -fPIC -Wall -O3

ifeq ($(PLATFORM),msiq-g20g40)
    # Matchstiq G20/G40
    #
    # Use Nvidia JetPack cross-compile Docker image, for example:
    #     nvcr.io/nvidia/jetpack-linux-aarch64-crosscompile-x86:5.1.2
    # or use the native compiler directly on the Matchstiq G20/G40 Jetson
    ifeq ($(BUILD_MACHINE),x86_64)
        TOOL_DIR=/usr/bin
        CROSS_PREFIX=aarch64-linux-gnu-
    endif
    ARCH=aarch64
    LIB_CONFIG=msiq-g20g40
endif

ifeq ($(PLATFORM),msiq-sx)
    # Matchstiq S1x/S2x at Release 9 or later (OpenWrt 16.02)
    BUILD_CONFIG=arm_cortex-a9.gcc5.2_glibc_openwrt
    LIB_CONFIG=msiq-sx
endif

ifeq ($(PLATFORM),msiq-x40)
    # Matchstiq X40
    #
    # Use Nvidia JetPack cross-compile Docker image, for example:
    #     nvcr.io/nvidia/jetpack-linux-aarch64-crosscompile-x86:5.1.2
    # or use the native compiler directly on the Matchstiq X40 Jetson
    ifeq ($(BUILD_MACHINE),x86_64)
        TOOL_DIR=/usr/bin
        CROSS_PREFIX=aarch64-linux-gnu-
    endif
    ARCH=aarch64
    LIB_CONFIG=msiq-x40
endif

ifeq ($(PLATFORM),z2-armhf)
    # Matchstiq Z2 (with hard-float extensions enabled)
    BUILD_CONFIG=arm_cortex-a9.gcc7.2.1_gnueabihf
    LIB_CONFIG=z2-armhf
endif

ifeq ($(PLATFORM),z3u)
    # Matchstiq Z3u
    BUILD_CONFIG=aarch64.gcc6.3
    LIB_CONFIG=z3u
endif

ifeq ($(BUILD_CONFIG),x86_32.gcc)
    # Host x86-32 with Sidekiq card
    TOOL_DIR=/usr/bin
    CROSS_PREFIX=
    CFLAGS+=-m32
    LDFLAGS+=-m32
    ARCH=x86_32
    LIB_CONFIG?=x86_32.gcc
endif

ifeq ($(BUILD_CONFIG),x86_64.gcc)
    # Host x86-64 with Sidekiq card
    TOOL_DIR=/usr/bin
    CROSS_PREFIX=
    CFLAGS+= -g -gdwarf-3
    ARCH=x86_64
    LIB_CONFIG?=x86_64.gcc
endif

ifeq ($(BUILD_CONFIG),aarch64.gcc6.3)
    # aarch64 cross-toolchain using GCC6.3.1, GLIBC 2.23, and GLIBC++ 6.0.22
    # For Matchstiq Z3u (not for Jetson or Matchstiq G20/G40/X40)
    #
    # Download from the following link and install to /opt/toolchain
    # https://releases.linaro.org/components/toolchain/binaries/6.3-2017.05/aarch64-linux-gnu/gcc-linaro-6.3.1-2017.05-x86_64_aarch64-linux-gnu.tar.xz
    ifeq ($(BUILD_MACHINE),x86_64)
        TOOL_DIR=/opt/toolchains/gcc-linaro-6.3.1-2017.05-x86_64_aarch64-linux-gnu/bin
        CROSS_PREFIX=aarch64-linux-gnu-
    endif
    ARCH=aarch64
    LIB_CONFIG?=aarch64.gcc6.3
endif

ifeq ($(BUILD_CONFIG),aarch64)
    # aarch64 cross-toolchain using GCC 7.3.1, GLIBC 2.24, and GLIBC++ 6.0.24
    # For aarch64 with Sidekiq card (not for Matchstiq Z3u/G20/G40/X40)
    #
    # Download from the following link and install to /opt/toolchain
    # https://releases.linaro.org/components/toolchain/binaries/7.3-2018.05/aarch64-linux-gnu/gcc-linaro-7.3.1-2018.05-x86_64_aarch64-linux-gnu.tar.xz
    ifeq ($(BUILD_MACHINE),x86_64)
        TOOL_DIR=/opt/toolchains/gcc-linaro-7.3.1-2018.05-x86_64_aarch64-linux-gnu/bin
        CROSS_PREFIX=aarch64-linux-gnu-
    endif
    ARCH=aarch64
    LIB_CONFIG?=aarch64
endif

ifeq ($(BUILD_CONFIG),arm_cortex-a9.gcc5.2_glibc_openwrt)
    ifeq ($(BUILD_MACHINE),x86_64)
        TOOL_DIR=/opt/toolchains/arm_cortex-a9.gcc5.2_glibc_openwrt/bin
        CROSS_PREFIX=arm-openwrt-linux-
        export STAGING_DIR=$(TOOL_DIR)
        CONFIG_FLAGS="--host=arm-openwrt-linux"
    endif
    CFLAGS+=-mfpu=neon -mfloat-abi=hard -mno-unaligned-access -ggdb
    ARCH=arm_cortex-a9.glibc
    LIB_CONFIG?=arm_cortex-a9.gcc5.2_glibc_openwrt
endif

ifeq ($(BUILD_CONFIG),arm_cortex-a9.gcc7.2.1_gnueabihf)
    ifeq ($(BUILD_MACHINE),x86_64)
       TOOL_VERSION=2018.2
       # Xilinx Tools installed by default to this location
       TOOL_DIR = /opt/Xilinx/SDK/$(TOOL_VERSION)/gnu/aarch32/lin/gcc-arm-linux-gnueabi/bin
       CROSS_PREFIX=arm-linux-gnueabihf-
       VIVADO_SETTINGS=/opt/Xilinx/Vivado/$(TOOL_VERSION)/settings64.sh
       CROSS_COMPILE=$(CROSS_PREFIX)
       CONFIG_FLAGS="--host=arm-linux-gnueabihf"
    endif
    LIB_CONFIG?=arm_cortex-a9.gcc7.2.1_gnueabihf
endif

# BUILD_CONFIG/PLATFORM must be in allowed list, but not configured above
ifeq (,$(LIB_CONFIG))
    $(error Bug in tools.mk! Unhandled BUILD_CONFIG=$(BUILD_CONFIG) or PLATFORM=$(PLATFORM))
endif

ifneq ($(PLATFORM),undefined)
  ifeq ($(filter $(PLATFORM),$(strip $(ALLOWED_PLATFORMS))),)
    $(info PLATFORM is not available, use one of the following:)
    $(info )
    $(info $(ALLOWED_PLATFORMS))
    $(info )
    $(error choose a supported PLATFORM (or BUILD_CONFIG))
  endif
else ifneq ($(BUILD_CONFIG),undefined)
  ifeq ($(filter $(BUILD_CONFIG),$(strip $(ALLOWED_BUILD_CONFIGS))),)
    $(info BUILD_CONFIG is not available, use one of the following:)
    $(info )
    $(info $(ALLOWED_BUILD_CONFIGS))
    $(info )
    $(info Or define PLATFORM to one of the following:)
    $(info )
    $(info $(ALLOWED_PLATFORMS))
    $(info )
    $(error choose a supported BUILD_CONFIG or PLATFORM)
  endif
else
  $(info BUILD_CONFIG is not defined, use one of the following:)
  $(info )
  $(info $(ALLOWED_BUILD_CONFIGS))
  $(info )
  $(info Or define PLATFORM to one of the following:)
  $(info )
  $(info $(ALLOWED_PLATFORMS))
  $(info )
  $(error choose a supported BUILD_CONFIG or PLATFORM)
endif

##############################################
#
# Define toolchain based on CROSS_COMPILE
#
##############################################

# support legacy TOOL_DIR / CROSS_PREFIX
# prior users of tools.mk may have set TOOL_DIR or CROSS_PREFIX but not both,
# so the configuration above still sets the default values and combines them
# into the new variable CROSS_COMPILE (but this also allows a user to override
# with 'make CROSS_COMPILE=...')
ifeq (,$(TOOL_DIR))
  CROSS_COMPILE = $(CROSS_PREFIX)
else
  CROSS_COMPILE = $(TOOL_DIR)/$(CROSS_PREFIX)
endif

CC      := $(CROSS_COMPILE)gcc
CXX     := $(CROSS_COMPILE)g++
LD      := $(CROSS_COMPILE)ld
NM      := $(CROSS_COMPILE)nm
OBJCOPY := $(CROSS_COMPILE)objcopy
AR      := $(CROSS_COMPILE)ar
RANLIB  := $(CROSS_COMPILE)ranlib
STRIP   := $(CROSS_COMPILE)strip

ifeq (,$(shell command -v $(CC) 2>/dev/null))
    $(error $(CC) not found, check that cross compiler is installed or specify alternative CROSS_COMPILE)
endif
