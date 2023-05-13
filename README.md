<!-- <p align="center">
日本語 | <a href="./README_en.md">English</a>
</p> -->

# ALINCO DJ-X100 CommandLine Tool

- 非公式のアルインコ [DJ-X100](https://www.alinco.co.jp/product/electron/detail/djx100.html) 用のコマンドラインツールです。

## :beginner: 使い方

- [Homebrew](https://brew.sh/index_ja)がインストールされていれば`brew install bellx2/tap/x100cmd`でインストールできます。
- ビルド済みバイナリーを[Releases](https://github.com/bellx2/x100cmd/releases/)よりダウンロードし任意の場所に置くか、ソースからのビルドも可能です。
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
- [`read`] - チャンネルデータ読み込み
- [`write`] - チャンネルデータ書き込み
- [`clear`] - チャンネルデータクリア
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

### `x100cmd read <channel>`

チャンネルデータを読み込みます。

```sh
x100cmd read 10

{"freq":433.000000, "mode":"FM", "step":"20k", "name":"430メイン", "empty": false}
```

### `x100cmd write <channel>`

チャンネルデータを書き込みます。指定したチャンネルのデータは再起動するまで反映されません。指定したデータ以外（コード等）は指定チャンネルの情報が保持されます。`-r`オプションを付けると書き込み後に再起動を行います。

```sh
x100cmd write 10 -f 433.00 -m FM -n "430メイン" -s "20k" -r
```

| フラグ            | 初期値 | 説明                                                                                                                                                                                   |
| ----------------- | ------ | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `-f`, `--freq`    | NULL   | 周波数 (e.g. 433.0)                                                                                                                                                                    |
| `-m`, `--mode`    | FM     | モード: FM, NFM, AM, NAM, T98 , T102_B54, DMR, T61_typ1, T61_typ2, T61_typ3, T61_typ4, dPMR,DSTAR, C4FM, AIS, ACERS, POCSAG, 12KIF_W, 12KIF_N <br />※対応してないモードは表示されない? |
| `-s`, `--step`    | 20k    | 周波数ステップ: 1k, 5k, 6k25, 8k33, 10k, 12k5,15k, 20k, 25k, 30k, 50k, 100k, 125k,200k                                                                                                 |
| `-n`, `--name`    | NULL   | 名称 (UTF-8) 最大: 30byte, NONE で空白                                                                                                                                                 |
| `-y`, `--yes`     | false  | 書き込み確認をしない                                                                                                                                                                   |
| `-r`, `--restart` | false  | 実行後再起動                                                                                                                                                                           |

### `x100cmd clear <channel>`

チャンネルデータを消去します。

```sh
x100cmd clear 10
OK
```

| フラグ            | 初期値 | 説明             |
| ----------------- | ------ | ---------------- |
| `-y`, `--yes`     | false  | 消去確認をしない |
| `-r`, `--restart` | false  | 実行後再起動     |

### `x100cmd exec <command>`

コントロールコマンドを送信します。

```sh
x100cmd exec sql 0  # SQL OFF
x100cmd exec restart # 再起動
```

| コマンド     | 説明                  |
| ------------ | --------------------- |
| freq \<freq> | 現在の周波数取得/設定 |
| gps          | 本体の GPS 情報の取得 |
| restart      | 再起動                |
| sql \<level> | SQL 設定              |
| vol \<level> | ボリューム設定        |

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
