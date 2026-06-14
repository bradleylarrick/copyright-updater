include tools.mk
ROOTDIR = $(CURDIR)

# Python Module Build Settings
PY_MOD_NAME = rxcapture
PY_SRC = $(SRC_DIR)rx_capture_module.c
PY_SO = script/lib/lib$(PY_MOD_NAME).so

PY_CFLAGS = -fPIC -std=gnu11 -Wall
PY_CFLAGS += -I$(ROOTDIR)/include -I$(ROOTDIR)/sidekiq_core/inc -I$(ROOTDIR)/arg_parser/inc

PY_LDFLAGS = -shared
PY_LIBS = $(STATIC_LIBS)
PY_LDLIBS = -lusb-1.0 -lglib-2.0 -ltirpc -lpthread -lrt -lm -lstdc++

CFLAGS += -std=gnu11
CFLAGS += -I$(ROOTDIR)/include -I$(ROOTDIR)/sidekiq_core/inc -I$(ROOTDIR)/arg_parser/inc
#CFLAGS += -I$(ROOTDIR)/tflite_aarch64/include/tensorflow/lite -I$(ROOTDIR)/include/tensorflow/lite/c -I$(ROOTDIR)/include/tensorflow/lite/core -ltensorflow-lite

OUT_DIR = ./bin/
SRC_DIR = ./src/

# LIB_CONFIG holds the name of the libsidekiq library that matches the selected
# PLATFORM or BUILD_CONFIG. See tools.mk for more details.

# define any libraries to link into executable:
STATIC_LIBS = $(ROOTDIR)/lib/libsidekiq__$(LIB_CONFIG).a $(ROOTDIR)/arg_parser/lib/arg_parser__$(LIB_CONFIG).a
SUPPORT = $(ROOTDIR)/lib/support/$(LIB_CONFIG)
LIBS+= $(STATIC_LIBS)

# glib2 / libusb / libiio locations / flags
LDFLAGS+= -L$(SUPPORT)/usr/lib/epiq 

#-L$(ROOTDIR)/tflite_aarch64/lib 
#-ltensorflowlite_c

LIBS+= -lusb-1.0 -lglib-2.0
ifeq ($(BUILD_CONFIG),arm_cortex-a9.gcc7.2.1_gnueabihf)
LIBS+= -liio
else ifeq ($(BUILD_CONFIG),aarch64.gcc6.3)
LIBS+= -liio
endif

ifeq ($(PLATFORM),msiq-x40)
LIBS+= $(shell PKG_CONFIG_PATH=$(SUPPORT)/usr/lib/epiq/pkgconfig pkg-config --libs-only-l grpc++ protobuf) -lgpiod
endif

# define any libraries to link into executable:
LIBS+= -ltirpc -lpthread -lrt -lm -lstdc++

#
# build executables with RUNPATH configured to look for shared libraries in the:
#   1. same directory as the executable
#   2. current SDK's lib/ directory, relative to the executable
#   3. system path /usr/lib/epiq (installed by Sidekiq runtime installer)
#
# The goal of RUNPATH being relative to $ORIGIN is to help with testing and to
# ensure the example executables load the shared libraries packaged in this
# version of the SDK.
#
# For custom applications, the following flags are recommended for informing
# GCC where to find the shared libraries used by libsidekiq when linking:
#     -L$(SUPPORT)/usr/lib/epiq -Wl,-rpath-link=$(SUPPORT)/usr/lib/epiq
# along with the following flags to ensure the libraries are found at runtime:
#     -Wl,--enable-new-dtags -Wl,-rpath=/usr/lib/epiq
#
ifeq ($(BUILD_CONFIG),arm_cortex-a9.gcc4.8_uclibc_openwrt)
LDFLAGS+= -Wl,--enable-new-dtags -Wl,-rpath,/usr/lib/epiq
else
LDFLAGS+= -Wl,--enable-new-dtags -Wl,-rpath='$$ORIGIN:$$ORIGIN/../../../lib/support/$(LIB_CONFIG)/usr/lib/epiq:/usr/lib/epiq' -Wl,-z,origin
endif

# the below apps are released to customers as part of the Sidekiq SDK
#TESTCSRCS=  src/inference.c
TESTCSRCS+= src/rx_capture.c
# TESTCSRCS+= src/tx_from_file.c
# TESTCSRCS+= src/tx_tone.c

TESTAPPS= $(patsubst src/%.c,bin/%,$(filter src/%.c,$(TESTCSRCS)))
TESTAPPS+= $(patsubst %.c,%,$(filter %.c,$(filter-out src/%.c,$(TESTCSRCS))))

all: $(TESTAPPS)

clean:
	rm -f src/*.o
	rm -f bin/*
	rm -f script/lib/lib$(PY_MOD_NAME).so

ifeq ($(SO_LINK),enabled)
ifeq ($(wildcard $(SUPPORT)/usr/lib/epiq/libsidekiq.so),)
$(error SO_LINK=enabled not supported for PLATFORM/BUILD_CONFIG at this time)
endif

# search path for inter-library dependencies
LDFLAGS+= -Wl,-rpath-link=$(SUPPORT)/usr/lib/epiq

# reset LIBS for shared object libsidekiq, only arg_parser is required
LIBS= $(ROOTDIR)/arg_parser/lib/arg_parser__$(LIB_CONFIG).a

# arg_parser dependency
LDLIBS+= -lglib-2.0
# libsidekiq dynamic library
LDLIBS+= -lsidekiq
# general example app libraries
LDLIBS+= -lpthread -lrt -lm
endif

# build the test executable in bin/ from src/
bin/%: src/%.o
	@mkdir -p $(dir $@)
	$(CC) $(LDFLAGS) -o $@ -Wl,--start-group $^ $(LIBS) $(LDLIBS) -Wl,--end-group
ifeq ($(STRIP_ENABLED),enabled)
	$(STRIP) $@
endif

python: $(PY_SO)

$(PY_SO): $(PY_SRC)
	@echo "Building Python shared object: $@"
	@mkdir -p $(dir $@)
	$(CC) $(PY_CFLAGS) -shared -o $@ $^ $(LDFLAGS) $(LIBS)
