# ⏳ 時間管理 Slack アプリ time-management-slack

![Screen Shot 2021-12-25 at 0 19 37](https://user-images.githubusercontent.com/66525257/147361800-31611653-fde8-47c5-afb1-1ccdbaf9eecf.png)
![Screen Shot 2021-12-25 at 0 19 52](https://user-images.githubusercontent.com/66525257/147361804-2eabf20d-40d7-4d2c-876f-7681e71cb651.png)
![Screen Shot 2021-12-25 at 0 20 12](https://user-images.githubusercontent.com/66525257/147361819-72f3954c-081b-4d9f-8ecb-56bd4eae0650.png)



## 背景・課題
私の勤務するインターン先では、自身の勤務時間をある勤怠管理サービスに打刻し、そのデータに応じて給与が振り込まれる。普通に考えれば業務の開始・終了時にその勤怠管理サービスに打刻を行えば良いのだが、学生である私は1日に 3, 4 回離脱と復帰を繰り返している。そのような狂った勤怠を直接報告するわけにはいかないので1日ごとに働いた時間を集計し、合計時間を勤怠管理サービスに打刻するようにしていた。（勤務の開始・終了時に slack で発言している。）

しかし毎日の勤務時間を計算するのがめんどくさく、ひどい時だと勤務時間の報告を1週間くらい溜めた時もあった。

そこで自身の狂った勤怠を管理すべく、slack の勤務開始、終了の宣言をから 1 日あたりの勤務時間を計算する slack app を作ることにした。


## デザイン
このアプリの機能は大きく分けて2つある。

- 打刻
- 時間取得

打刻は業務の開始、終了時に行い、時間取得は日毎の合計時間を取得する時に用いる。

また、slack チャンネル内からこのアプリを利用するインターフェイスとしては slack slush command を選んだ。
サポートしているコマンドは以下の通りである。

- `/start` - 開始
- `/end` - 終了
- `/times` - 時間取得

![Screen Shot 2021-12-25 at 0 31 39](https://user-images.githubusercontent.com/66525257/147362282-908b9741-d9f9-4f13-af2d-22a0f46e6bbe.png)

### `/start`, `/end` の仕様
業務の開始、終了時に用いる。`/start はじめ` と入力すると、アプリ内で打刻が保存され、slack チャンネルには引数にとった文字列を表示する。この文字列を自由に指定することでその打刻がその日初めてなのか、復帰してきたのかわかるようにしてある。`/end` についても同様である。

`/start はじめ` と入力した例

![Screen Shot 2021-12-25 at 0 19 37](https://user-images.githubusercontent.com/66525257/147361800-31611653-fde8-47c5-afb1-1ccdbaf9eecf.png)

また、無効な入力についてはエラーメッセージを返すようにしてある。例えば「`/start` を2回連続で入力する」などである。

### `/times` の仕様
勤務時間の合計を計算して表示する。1日に何度 `/start` と `/end` を繰り返しても適切な合計時間を返すようにしている。  
また `/times 3` のように引数に数字を取ることで n 日前 ~ 今日までの勤怠時間を表示することもできる。(0 ≦ n ≦ 6）

`/times 3` と入力した例

![Screen Shot 2021-12-25 at 0 45 38](https://user-images.githubusercontent.com/66525257/147362898-58afe733-0c41-497b-9eaa-42a986699d7c.png)

また、その日の勤務がまだ終わっていない時には、現在時刻までの時間で計算するようにしている。

## アーキテクチャ
このアプリはDB に打刻を保存したり、勤務時間を計算するためのサーバ（このリポジトリ）と slack slush command を受け付けるためのサーバ（[igsr5/time-management-bolt](https://github.com/igsr5/time-management-bolt)）の2つから構成される。
それぞれのサーバは AWS の ECS 上にデプロイされている。

![勤怠slackapp (1)](https://user-images.githubusercontent.com/66525257/147362062-970befb5-9e13-4a05-b457-7e5118925220.jpg)
https://miro.com/app/board/o9J_lk2M0YU=/


## Development
### Setup
- you need to install docker, docker-compose, golang.
- you need to setting these path.
```sh
$ git clone (URL)
$ make setup
```

### Run
```sh
$ docker compose start && go run ./cmd/server.go
```

### Build
```sh
# output bin/api-server
$ make build
```

## Generate
```sh
$ make gen
```

### Migrate
```sh
# when not permission error
$ docker compose start && docker compose exec app chmod +x ./script/migrate-XX

#create
$ docker compose start && docker compose exec app ./script/migrate-create XXXXXX
# up
$ docker compose start && docker compose exec app ./script/migrate-up
# down
$ docker compose start && docker compose exec app ./script/migrate-down
```
