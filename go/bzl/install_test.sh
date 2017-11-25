set -euox pipefail

if  ! ./go/bzl/bzl install | grep '0.7.0'; then
    exit 1
else
    echo "PASS"
fi

if  ! ./go/bzl/bzl install 0.7.0 | grep 'Installed'; then
    exit 1
else
    echo "PASS"
fi

