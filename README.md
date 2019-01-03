# tkradar

[![Build Status](https://travis-ci.org/jiro4989/tkradar.svg?branch=master)](https://travis-ci.org/jiro4989/tkradar)

Classes.jsonをレーダーチャート画像に変換する

## できること

ツクールMVのClasses.jsonから以下のような能力値のSVGレーダーチャートを生成する。

![能力値](./testdata/out/class001.svg)

## インストール

`go get github.com/jiro4989/tkradar`

または

GitHub Releaseからダウンロード。

## 使い方

以下のコマンドを実行するとカレントディレクトリは以下にclass%03d.svgが生成される。

`tkradar Classes.json`

