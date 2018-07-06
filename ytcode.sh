#!/bin/bash

INDIR=$HOME/Music
OUTDIR=$HOME/Plex/Music

TRGTS=(`find $INDIR -type f`)

for IN in "${TRGTS[@]}" ; do
	TMP=${IN//$INDIR}
	OUT=$OUTDIR${TMP//.opus}.m4a
	DIR=$(dirname $OUT)
	mkdir -p $(dirname $OUT)
	ffmpeg -i $IN -c:a libfdk_aac $OUT
done
