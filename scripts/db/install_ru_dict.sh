#!/bin/bash

POSTGRES_VER=14

apt update -y && apt install myspell-ru -y
cd /usr/share/postgresql/$POSTGRES_VER/tsearch_data
DICT=/usr/share/hunspell/ru_RU

iconv -f koi8-r -t utf-8 -o russian.affix $DICT.aff
iconv -f koi8-r -t utf-8 -o russian.dict $DICT.dic
