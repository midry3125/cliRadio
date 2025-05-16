# cliRadio
CLI、つまりコンソール上でラジオを聴くことが出来ます

# 使う前に
本ソフトはffplayに依存しています  
したがって、本ソフトを使用する前にffmpegをインストールする必要があります  
(ffplayはffmpegに付属しています)

また、ラジコのAPIを使用しているため、前提として日本国内にいる必要があります

# ダウンロード
こちらの[リリースページ](https://github.com/midry3125/cliRadio/releases/latest)より、cliRadio.exeをダウンロードしてください

# 使い方
1.  ダウンロードしたcliRadio.exeをダブルクリック等で起動
2.  十字キーで聴きたい放送局を選び、エンターキーを押す
3.  ラジオが聴ける! 
4.  終了したい場合は、Ctrl+Cキーを押してください

# ビルド方法
※go言語のコンパイラが必要です
```bash
$ git clone https://github.com/midry3125/cliRadio && cd cliRadio && go build
```

# ライセンス
[MITライセンス](./LICENSE)です
