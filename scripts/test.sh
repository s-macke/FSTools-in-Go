#!/bin/bash
set -e

RED='\033[0;31m'
NC='\033[0m'

function exec {
    cmd="$1"
    printf "${RED}Execute ${cmd}${NC}\n"
    proot -r testenv -w / ${cmd}
    echo ============================
}

rm -rf testenv

mkdir -p testenv
cp bin/*go testenv/

cd testenv
mkdir -p {a,b}/{e,f,g}/{h,i,j}
cd ..

cmd="./findgo -maxdepth 1"
exec "${cmd}"

cmd="./dugo"
exec "${cmd}"

cmd="./rmgo -r -v a b"
exec "${cmd}"

cmd="./dugo -s"
exec "${cmd}"


