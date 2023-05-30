<!-- <p align="center">
日本語 | <a href="./README_en.md">English</a>
</p> -->

# ALINCO DJ-X100 CommandLine Tool

- 非公式のアルインコ [DJ-X100](https://www.alinco.co.jp/product/electron/detail/djx100.html) 用のコマンドラインツールです。

## :beginner: 使い方

- macos の場合[Homebrew](https://brew.sh/index_ja)がインストール済みであれば`brew install bellx2/tap/x100cmd`でインストール可能です。アップデートは`brew upgrade x100cmd`で行えます。
- Windows/macos 用のビルド済みバイナリーは[Releases](https://github.com/bellx2/x100cmd/releases/)よりダウンロードできます。任意の場所に置いて実行してください。
- [DJ-X100](https://www.alinco.co.jp/product/electron/detail/djx100.html) を USB ケーブルで接続します。

`read` コマンドを使って指定チャンネルのデータを読み込み表示します。ポートは自動認識します。

```sh
x100cmd read 10
```

`write`コマンドを使って指定チャンネルへデータを書き込みます。
データは再起動するまで反映されません。`-r` オプションを付けると書き込み後に再起動を行います。

```sh
x100cmd write 10 -f 433.00 -m FM -n "430メイン" -s "20k" -r
```

その他、コマンド制御などが可能です。

```sh
x100cmd exec restart
```

## :rocket: コマンド

コマンド一覧:

- [`check`] - シリアルポートと接続確認
- [`ch`] - チャンネルコマンド（省略可）

  - [`read`] - チャンネルデータ読み込み
  - [`write`] - チャンネルデータ書き込み
  - [`clear`] - チャンネルデータクリア
  - [`export`] - チャンネルデータのファイル出力
  - [`import`] - チャンネルデータのファイル読み込み

- [`bank`] - バンクコマンド(初期ファームのみ)

  - [`read`] - バンク名読み込み
  - [`write`] - バンク名書き込み

- [`exec`] - 制御コマンド実行

| グローバルフラグ | デフォルト | 説明                                               |
| ---------------- | ---------- | -------------------------------------------------- |
| `-p`, `--port`   | `auto`     | シリアルポート名 <br/>`auto`の場合は自動検索します |
| `--debug`        | false      | デバッグ表示                                       |

### `x100cmd check`

接続状態を確認します。

```sh
x100cmd check

** scan ports **
/dev/cu.wlan-debug
/dev/cu.Bluetooth-Incoming-Port
/dev/cu.usbmodem00000000000001 [3614:D001] DJ-X100!

** check connection **
PortName: auto
DJ-X100 PortName: /dev/cu.usbmodem00000000000001

** send command **
OK

** current freq **
433.000000
```

### `x100cmd read <channel_no>`<br/>`x100cmd ch read <channel_no>`

チャンネルデータを読み込みます。

```sh
x100cmd read 10

{"freq":433.000000, "mode":"FM", "step":"1k", "name":"430メイン", "offset":"OFF", "shift_freq":"0.000000", "att":"OFF", "sq":"OFF", "tone":"670", "dcs":"017", "bank":"ABCTYZ", "empty": false}
```

### `x100cmd write <channel>`<br/>`x100cmd ch write <channel_no>`

チャンネルデータを書き込みます。指定したチャンネルのデータは再起動するまで反映されません。指定したデータ以外（コード等）は指定チャンネルの情報が保持されます。`-r`オプションを付けると書き込み後に再起動を行います。

```sh
x100cmd write 10 -f 433.00 -m FM -n "430メイン" -s "20k" -r
```

※ 指定しない場合は現状のチャンネル設定を保持します
| フラグ | 初期値 | 説明 |
| ----------------- | ------ | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `-f`, `--freq` | | 周波数 (e.g. 433.0) |
| `-m`, `--mode` | | モード: FM, NFM, AM, NAM, T98 , T102_B54, DMR, T61_typ1, T61_typ2, T61_typ3, T61_typ4,T61_typx,ICUD,dPMR,DSTAR, C4FM, AIS, ACARS, POCSAG, 12KIF_W, 12KIF_N <br />※対応してないモードは表示されない? |
| `-s`, `--step` | | 周波数ステップ: 1k, 5k, 6k25, 8k33, 10k, 12k5,15k, 20k, 25k, 30k, 50k, 100k, 125k,200k |
| `-n`, `--name` | | 名称 (UTF-8) 最大: 30byte, NONE で空白 |
| `--offset` | | オフセット : ON, OFF |
| `--shift_freq` | | シフト周波数 |
| `--att` | | アッテネータ: OFF, 10db, 20db |
| `--sq` | | スケルチ: OFF,CTCSS,DCS,R_CTCSS,R_DCS,JR,MSK |
| `--tone` | | CTSS トーン: 670,693...2503,2541 |
| `--dcs` | | DCS コード: 017-754 |
| `--bank` | | バンク: A-Z ex. `ABCZ` のように複数指定可。`NONE`で消去|
| `--skip` | | メモリースキップ: ON, OFF|
| `--lat` | | 緯度 ex.35.681382 ※緯度経度=0 で消去|
| `--lon` | | 経度 ex.139.766084 ※緯度経度=0 で消去|
| `--ext` | | 拡張情報(0x50-0x7F) 半角 96 文字 |
| `-y`, `--yes` | false | 書き込み確認をしない |
| `-r`, `--restart` | false | 実行後再起動 |

### `x100cmd clear <channel_no>`<br/>`x100cmd ch clear <channel_no>`

チャンネルデータを初期データで消去します。

```sh
x100cmd clear 10
OK
```

| フラグ            | 初期値 | 説明             |
| ----------------- | ------ | ---------------- |
| `-y`, `--yes`     | false  | 消去確認をしない |
| `-r`, `--restart` | false  | 実行後再起動     |

### `x100cmd export <csv_filename>`<br/>`x100cmd ch export <csv_filename>`

チャンネルデータ(1-999)をエクスポートします。
フォーマットはカンマ区切りの CSV (UTF-8 BOM 付き)です。

```sh
x100cmd export channels.csv
```

| フラグ        | 初期値 | 説明                                       |
| ------------- | ------ | ------------------------------------------ |
| `-y`, `--yes` | false  | 上書き警告を表示しない                     |
| `-a`, `--all` | false  | 空のチャンネルデータ(周波数=0)も出力する   |
| `--ext`       | false  | 拡張情報(0x50-0x7F) 半角 96 文字を出力する |

#### ファイル形式

```:csv
Channel,Freq,Mode,Step,Name,offset,shift_freq,att,sq,tone,dcs,bank,lat,lon,skip
001,433.000000,FM,10k,430メイン,OFF,0.000000,OFF,OFF,670,017,A,0.000000,0.000000,OFF
002,145.000000,FM,10k,144メイン,OFF,0.000000,OFF,OFF,670,017,Z,0.000000,0.000000,OFF
....
```

- `--ext`オプション付きの場合

```:csv
Channel,Freq,Mode,Step,Name,offset,shift_freq,att,sq,tone,dcs,bank,lat,lon,skip,ext
001,433.000000,FM,10k,430メイン,OFF,0.000000,OFF,OFF,670,017,A,0.000000,0.000000,OFF,0000e4000000e400000000000000000000000180018001800180010000800100008001000080000080008000807b1700
002,145.000000,FM,10k,144メイン,OFF,0.000000,OFF,OFF,670,017,Z,0.000000,0.000000,OFF,0000e4000000e400000000000000000000000180018001800180010000800100008001000080000080008000807b1700
....
```

### `x100cmd import <csv_filename>`<br/>`x100cmd ch import <csv_filename>`

チャンネルデータをインポートします。フォーマットはカンマ区切りの CSV (UTF-8 BOM 付き)です。
Export したデータを基準にしてください。

- 指定したもの以外のデータは保持されます。
- 周波数が０の場合はチャンネルデータを消去します。
- 位置情報は Lat,Lon がいずれも 0.0 の場合データを消去します。
- ヘッダー行が異なるとインポートできないばあいがあります
- ファイルの文字コードは UTF-8 である必要があります
- ext 項目についてはそのまま書き込まれます

```sh
x100cmd import channels.csv
```

#### ファイル形式

- シンプル

```:csv
Channel,Freq,Mode,Step,Name
001,433.000000,FM,10k,430メイン
002,145.000000,FM,10k,144メイン
....
```

- 12 項目

```:csv
Channel,Freq,Mode,Step,Name,offset,shift_freq,att,sq,tone,dcs,bank
001,433.000000,FM,10k,430メイン,OFF,0.000000,OFF,OFF,670,017,A
002,145.000000,FM,10k,144メイン,OFF,0.000000,OFF,OFF,670,017,Z
```

- 位置情報付き（現行バージョンの export データ）

```:csv
Channel,Freq,Mode,Step,Name,offset,shift_freq,att,sq,tone,dcs,bank,lat,lon,skip
001,433.000000,FM,10k,430メイン,OFF,0.000000,OFF,OFF,670,017,A,0.000000,0.000000,OFF
002,145.000000,FM,10k,144メイン,OFF,0.000000,OFF,OFF,670,017,Z,0.000000,0.000000,OFF
....
```

- `--ext`オプション付きデータ

```:csv
Channel,Freq,Mode,Step,Name,offset,shift_freq,att,sq,tone,dcs,bank,lat,lon,skip,ext
001,433.000000,FM,10k,430メイン,OFF,0.000000,OFF,OFF,670,017,A,0.000000,0.000000,OFF,0000e4000000e400000000000000000000000180018001800180010000800100008001000080000080008000807b1700
002,145.000000,FM,10k,144メイン,OFF,0.000000,OFF,OFF,670,017,Z,0.000000,0.000000,OFF,0000e4000000e400000000000000000000000180018001800180010000800100008001000080000080008000807b1700
....
```

| フラグ            | 初期値 | 説明                             |
| ----------------- | ------ | -------------------------------- |
| `-v`, `--verbose` | false  | 書き込み中データの詳細を表示する |
| `-r`, `--restart` | false  | 書き込み後に再起動               |

### `x100cmd bank read <A-Z>`

**(初期ファームウェア v1.00 のみ)**

バンク名を読み込みます。バンクは A-Z で指定します。省略した場合はすべてのバンクを出力します

```sh
x100cmd bank read A
"A","羽田空港"

x100cmd bank read ABC
"A","羽田空港"
"B","成田空港"
"C","横田基地"
```

### `x100cmd bank write <A-Z> <bank_name>`

**(ファームウェア v1.00 のみ)**

バンク名を書き込みます。バンクは A-Z で指定します。再起動するまで反映されません。`-r`オプションを付けると書き込み後に再起動を行います。名称に`NONE`を指定すると名称を消去します（表示は`バンク-A~Z`となります）。

```sh
x100cmd bank write A "羽田空港" -r
```

| フラグ            | 初期値 | 説明               |
| ----------------- | ------ | ------------------ |
| `-y`, `--yes`     | false  | 上書き確認をしない |
| `-r`, `--restart` | false  | 実行後再起動       |

### `x100cmd exec <command>`

コントロールコマンドを送信します。

```sh
x100cmd exec restart # 再起動
```

| コマンド                | 説明                     |
| ----------------------- | ------------------------ |
| version                 | バージョン情報の取得     |
| restart                 | 再起動                   |
| read \<address>         | メモリー読み込み         |
| write \<address> <data> | メモリー書き込み 265Byte |

**(ファームウェア v1.00 のみ)**

| コマンド     | 説明                   |
| ------------ | ---------------------- |
| freq \<freq> | 現在の周波数取得/設定  |
| gps          | 本体の GPS 情報の取得  |
| sql \<level> | SQL 設定 (00-32)       |
| vol \<level> | ボリューム設定 (00-32) |

## 制限事項など

- 非公式ツールであり動作保証はありません。
- コマンド引数やレスポンスは予告なく変更する場合があります。
- 本体のデータを書き換えるため、データが破損する可能性があります。自己責任でご利用ください。
- 開発は MacOS で行っています。それ以外のプラットフォームの積極的な動作確認は行なっていません。
- コマンドライン文字列は UTF-8 です。
- [Windows] DJ-X100 を接続しても認識しない場合があります。その場合は`x100cmd check`コマンドでシリアルポート状況を確認するか、[メーカ提供の DJ-X100 ソフトウエア
  ](https://www.alinco.co.jp/product/electron/soft/softdl02/index.html)でまず接続を確認してください。認識できない場合は動作しません。
- 対応していない周波数やモードを書き込んだ場合の動作は不明です

## :memo: ライセンス

[MIT License](./LICENSE)

## 謝辞

<https://github.com/musen23872/djx100-commandline-tools>

- メモリーデータ構造の一部は`djx100-unofficial-memory-data.hexpat`を参考にさせていただきました。
