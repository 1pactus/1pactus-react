#!/bin/sh
set -e

PACTUS_HOME=/root/pactus

if [ ! -f "$PACTUS_HOME/config.toml" ]; then
    echo "init pactus node ..."
    
    WALLET_PASSWORD=${PACTUS_WALLET_PASSWORD?}
    
    yes | pactus-daemon init --password "$WALLET_PASSWORD" --val-num 1
    
    echo "init completed"
fi

echo "start node ..."
exec pactus-daemon start --password "${PACTUS_WALLET_PASSWORD}"
