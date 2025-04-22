#!/bin/bash

set -e

docker run -i --rm \
      --cap-add=SYS_ADMIN \
      --security-opt seccomp=unconfined \
      -v "$(realpath "$1")":/sandbox/script.py \
      judge-python \
      nsjail \
          --mode=o \
          --quiet \
          --chroot / \
          --proc_rw \
          --disable_clone_newcgroup \
          -- \
          /usr/bin/python3 /sandbox/script.py