# cat-audio

Concat Audio Files

## Usage

```prompt
$ cat-audio -out result.wav before.wav after.wav ...
```

## Install

Pre-built binaryies are available on: https://github.com/karupanerura/cat-audio/releases/tag/v0.0.1

```prompt
$ VERSION=0.0.1
$ curl -sfLO https://github.com/karupanerura/cat-audio/releases/download/v${VERSION}/cat-audio_${VERSION}_$(go env GOOS)_$(go env GOARCH).tar.gz
$ tar zxf cat-audio_${VERSION}_$(go env GOOS)_$(go env GOARCH).tar.gz
$ install -m 0755 cat-audio $PREFIX
$ rm cat-audio cat-audio_${VERSION}_$(go env GOOS)_$(go env GOARCH).tar.gz
```

## Restriction

Supported 44100Hz/2ch/16bit WAVE and MP3 files ONLY.

Patches welcome :)
