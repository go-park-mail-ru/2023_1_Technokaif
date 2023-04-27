#!/bin/bash

sudo apt install myspell-ru
cd /usr/share/postgresql/14/tsearch_data
DICT=/usr/share/hunspell/ru_RU

sudo iconv -f koi8-r -t utf-8 -o russian.affix $DICT.aff
sudo iconv -f koi8-r -t utf-8 -o russian.dict $DICT.dic
