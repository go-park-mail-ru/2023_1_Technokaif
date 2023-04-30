#!/bin/bash

apt update
apt install myspell-ru
cd /usr/share/postgresql/15/tsearch_data
DICT=/usr/share/hunspell/ru_RU

iconv -f koi8-r -t utf-8 -o russian.affix $DICT.aff
iconv -f koi8-r -t utf-8 -o russian.dict $DICT.dic
