# sge-tech-book-vol3-part1
技術書展に出版するSGE Go Tech Book Vol.03の1章のサンプルコードです。

## 構成
- cmd: 起動処理や依存性の注入
- config: 設定値
- ddl: データベースのデータ定義
- docker: Dockerfileを管理
- pkg: ロジックの実装
  - entity
  - handler
  - ifrepository: リポジトリのインターフェース
  - infra
    - mysql
      - cachedb: データベース操作のキャッシュに関する処理
      - cachemodel: 各テーブルにおけるデータベース操作のキャッシュ
      - cacherepository: 各テーブルにおけるデータベース操作のキャッシュに関する処理
      - repository: 各テーブルに対するデータベース操作の処理
      - integrationrepositry: cacherepositoryとrepositoryを切り替え
      - template: sqlboilerのためのテンプレート
      - datamodel: sqlboilerによる自動生成されたコード
      - client.go: データベースのクライアント
  - masterdata: マスターデータ
  - service
  - usecase
- script

## 提供する機能
- ユーザ作成API
- ユーザデータ取得API
- モンスター強化API

上記では、最初にユーザ作成APIを実行する必要があります。

モンスター強化APIでは、アプローチの評価のためにデータベースアクセスを発生させる不要な処理を加えています。

## PostmanでAPIリクエストを実行できるようにする
Postmanに`Title.postman_collection.json`と`GameServerExample環境.postman_environment.json`をimportして、
CollectionとEnvironmentsに上記のデータを追加します。
そして、GameServerExample環境のEnvironmentを有効化します。

有効化の手順は以下を参考にしてください。

https://learning.postman.com/docs/sending-requests/managing-environments/#selecting-an-active-environment

## 動かし方
コンテナについて、CPU制限・メモリ制限をしていますが、
スペックは適宜調整していただいて問題ございません。

初回のコンテナ起動
```shell
make docker_up_db
# 上記でmysqlコンテナが立ち上がって
# リクエストを受け付けられるようになったら、以下を実行します
# 数秒程度待てば、リクエストを受け付けられるようになると思います
# 動作環境のスペックによっては、リクエストを受け付けられるようになるまで少し時間がかかる可能性があります
make exec_db_setting
make docker_up
```

2回目以降のコンテナ起動
```shell
make docker_up
```

## データベース操作のコンテキストキャッシュの評価方法
以下の手順で評価します。

### 1. 提案手法を無効化したサーバに対してリクエストを行う
まず初めに、docker-compose.ymlのgameサービスの以下の環境変数の値を0に更新します。
```text
ENABLE_CACHE_REPOSITORY
```

以下を実行します。
```shell
# コンテナを起動
make docker_up
```

上記でmysqlコンテナがリクエストを受け付けられるようになったら、以下を実行します。
```shell
# general_logの設定
make exec_general_log_setting
```

Postmanを利用して、以下の順番でAPIを実行します。
1. ユーザ作成API
2. ユーザデータ取得API
3. モンスター強化API

次に、以下を実行します。
```shell
# データベースのgeneral_logの保存
make save_result FILEPATH=保存先のファイル名
```

### 2. 提案手法を有効化したサーバに対して処理を行う
以下を実行します。
```shell
# 前回起動したコンテナを削除
make docker_down
```

次に、docker-compose.ymlのgameサービスの以下の環境変数の値を1に更新します。
```text
ENABLE_CACHE_REPOSITORY
```

以降は、1の手順を、コンテナを起動させる所から実施します。

### 3. 結果の比較
2と3で保存したデータベースのgeneral_logの結果を比較します。
以下は、著者の環境での実行した際の結果です。
提案手法を有効化することで、データベースで実行されたSQLの数が減っていることがわかります。

#### 提案手法を無効化した際のgeneral_log
```text
2023-05-09T07:06:44.329850Z	    8 Query	START TRANSACTION
2023-05-09T07:06:44.336075Z	    8 Prepare	INSERT INTO `user_coin` (`user_id`,`num`) VALUES (?,?)
2023-05-09T07:06:44.336236Z	    8 Execute	INSERT INTO `user_coin` (`user_id`,`num`) VALUES (2590299773,0)
2023-05-09T07:06:44.338138Z	    8 Close stmt	
2023-05-09T07:06:44.338337Z	    8 Prepare	UPDATE `user_coin` SET `num`=? WHERE `user_id`=?
2023-05-09T07:06:44.338444Z	    8 Execute	UPDATE `user_coin` SET `num`=100000 WHERE `user_id`=2590299773
2023-05-09T07:06:44.338716Z	    8 Close stmt	
2023-05-09T07:06:44.339523Z	    8 Prepare	INSERT INTO `user_monster` (`user_id`,`monster_id`,`exp`) VALUES (?,?,?)
2023-05-09T07:06:44.339620Z	    8 Execute	INSERT INTO `user_monster` (`user_id`,`monster_id`,`exp`) VALUES (2590299773,1,0)
2023-05-09T07:06:44.340785Z	    8 Close stmt	
2023-05-09T07:06:44.340906Z	    8 Prepare	INSERT INTO `user_monster` (`user_id`,`monster_id`,`exp`) VALUES (?,?,?)
2023-05-09T07:06:44.341107Z	    8 Execute	INSERT INTO `user_monster` (`user_id`,`monster_id`,`exp`) VALUES (2590299773,2,0)
2023-05-09T07:06:44.341407Z	    8 Close stmt	
2023-05-09T07:06:44.341517Z	    8 Prepare	INSERT INTO `user_monster` (`user_id`,`monster_id`,`exp`) VALUES (?,?,?)
2023-05-09T07:06:44.341613Z	    8 Execute	INSERT INTO `user_monster` (`user_id`,`monster_id`,`exp`) VALUES (2590299773,3,0)
2023-05-09T07:06:44.341840Z	    8 Close stmt	
2023-05-09T07:06:44.342512Z	    8 Prepare	INSERT INTO `user_item` (`user_id`,`item_id`,`count`) VALUES (?,?,?)
2023-05-09T07:06:44.342652Z	    8 Execute	INSERT INTO `user_item` (`user_id`,`item_id`,`count`) VALUES (2590299773,1,0)
2023-05-09T07:06:44.343754Z	    8 Close stmt	
2023-05-09T07:06:44.343860Z	    8 Prepare	UPDATE `user_item` SET `count`=? WHERE `user_id`=? AND `item_id`=?
2023-05-09T07:06:44.343974Z	    8 Execute	UPDATE `user_item` SET `count`=1000 WHERE `user_id`=2590299773 AND `item_id`=1
2023-05-09T07:06:44.344231Z	    8 Close stmt	
2023-05-09T07:06:44.344326Z	    8 Prepare	INSERT INTO `user_item` (`user_id`,`item_id`,`count`) VALUES (?,?,?)
2023-05-09T07:06:44.344423Z	    8 Execute	INSERT INTO `user_item` (`user_id`,`item_id`,`count`) VALUES (2590299773,2,0)
2023-05-09T07:06:44.344646Z	    8 Close stmt	
2023-05-09T07:06:44.344744Z	    8 Prepare	UPDATE `user_item` SET `count`=? WHERE `user_id`=? AND `item_id`=?
2023-05-09T07:06:44.344839Z	    8 Execute	UPDATE `user_item` SET `count`=1000 WHERE `user_id`=2590299773 AND `item_id`=2
2023-05-09T07:06:44.345059Z	    8 Close stmt	
2023-05-09T07:06:44.345150Z	    8 Prepare	INSERT INTO `user_item` (`user_id`,`item_id`,`count`) VALUES (?,?,?)
2023-05-09T07:06:44.345247Z	    8 Execute	INSERT INTO `user_item` (`user_id`,`item_id`,`count`) VALUES (2590299773,3,0)
2023-05-09T07:06:44.345466Z	    8 Close stmt	
2023-05-09T07:06:44.345577Z	    8 Prepare	UPDATE `user_item` SET `count`=? WHERE `user_id`=? AND `item_id`=?
2023-05-09T07:06:44.345689Z	    8 Execute	UPDATE `user_item` SET `count`=1000 WHERE `user_id`=2590299773 AND `item_id`=3
2023-05-09T07:06:44.345899Z	    8 Close stmt	
2023-05-09T07:06:44.345938Z	    8 Query	COMMIT
2023-05-09T07:06:48.710906Z	    8 Query	START TRANSACTION
2023-05-09T07:06:48.712365Z	   11 Connect	root@***.***.***.*** on game_server_example using TCP/IP
2023-05-09T07:06:48.712715Z	   11 Query	SET NAMES utf8mb4
2023-05-09T07:06:48.713425Z	   11 Prepare	SELECT `user_monster`.* FROM `user_monster` WHERE (`user_monster`.`user_id` = ?)
2023-05-09T07:06:48.713513Z	   11 Execute	SELECT `user_monster`.* FROM `user_monster` WHERE (`user_monster`.`user_id` = 2590299773)
2023-05-09T07:06:48.713971Z	   11 Close stmt	
2023-05-09T07:06:48.714154Z	   11 Prepare	SELECT `user_item`.* FROM `user_item` WHERE (`user_item`.`user_id` = ?)
2023-05-09T07:06:48.714351Z	   11 Execute	SELECT `user_item`.* FROM `user_item` WHERE (`user_item`.`user_id` = 2590299773)
2023-05-09T07:06:48.714920Z	   11 Close stmt	
2023-05-09T07:06:48.715166Z	   11 Prepare	SELECT `user_coin`.* FROM `user_coin` WHERE (`user_coin`.`user_id` = ?) LIMIT 1
2023-05-09T07:06:48.715343Z	   11 Execute	SELECT `user_coin`.* FROM `user_coin` WHERE (`user_coin`.`user_id` = 2590299773) LIMIT 1
2023-05-09T07:06:48.715651Z	   11 Close stmt	
2023-05-09T07:06:48.715709Z	    8 Query	COMMIT
2023-05-09T07:06:52.155337Z	    8 Query	START TRANSACTION
2023-05-09T07:06:52.155914Z	   11 Prepare	SELECT `user_item`.* FROM `user_item` WHERE (`user_item`.`user_id` = ?) AND (`user_item`.`item_id` IN (?,?))
2023-05-09T07:06:52.156030Z	   11 Execute	SELECT `user_item`.* FROM `user_item` WHERE (`user_item`.`user_id` = 2590299773) AND (`user_item`.`item_id` IN (1,2))
2023-05-09T07:06:52.156370Z	   11 Close stmt	
2023-05-09T07:06:52.156563Z	    8 Prepare	UPDATE `user_item` SET `count`=? WHERE `user_id`=? AND `item_id`=?
2023-05-09T07:06:52.156678Z	    8 Execute	UPDATE `user_item` SET `count`=985 WHERE `user_id`=2590299773 AND `item_id`=1
2023-05-09T07:06:52.156994Z	    8 Close stmt	
2023-05-09T07:06:52.157159Z	    8 Prepare	UPDATE `user_item` SET `count`=? WHERE `user_id`=? AND `item_id`=?
2023-05-09T07:06:52.157179Z	    8 Execute	UPDATE `user_item` SET `count`=970 WHERE `user_id`=2590299773 AND `item_id`=2
2023-05-09T07:06:52.157429Z	    8 Close stmt	
2023-05-09T07:06:52.157633Z	   11 Prepare	SELECT `user_monster`.* FROM `user_monster` WHERE (`user_monster`.`user_id` = ?) AND (`user_monster`.`monster_id` = ?) LIMIT 1
2023-05-09T07:06:52.157739Z	   11 Execute	SELECT `user_monster`.* FROM `user_monster` WHERE (`user_monster`.`user_id` = 2590299773) AND (`user_monster`.`monster_id` = 1) LIMIT 1
2023-05-09T07:06:52.158009Z	   11 Close stmt	
2023-05-09T07:06:52.158180Z	    8 Prepare	UPDATE `user_monster` SET `exp`=? WHERE `user_id`=? AND `monster_id`=?
2023-05-09T07:06:52.158271Z	    8 Execute	UPDATE `user_monster` SET `exp`=1700 WHERE `user_id`=2590299773 AND `monster_id`=1
2023-05-09T07:06:52.158537Z	    8 Close stmt	
2023-05-09T07:06:52.158681Z	   11 Prepare	SELECT `user_coin`.* FROM `user_coin` WHERE (`user_coin`.`user_id` = ?) LIMIT 1
2023-05-09T07:06:52.158778Z	   11 Execute	SELECT `user_coin`.* FROM `user_coin` WHERE (`user_coin`.`user_id` = 2590299773) LIMIT 1
2023-05-09T07:06:52.159026Z	   11 Close stmt	
2023-05-09T07:06:52.159209Z	    8 Prepare	UPDATE `user_coin` SET `num`=? WHERE `user_id`=?
2023-05-09T07:06:52.159318Z	    8 Execute	UPDATE `user_coin` SET `num`=99700 WHERE `user_id`=2590299773
2023-05-09T07:06:52.159614Z	    8 Close stmt	
2023-05-09T07:06:52.159761Z	   11 Prepare	SELECT `user_monster`.* FROM `user_monster` WHERE (`user_monster`.`user_id` = ?) AND (`user_monster`.`monster_id` = ?) LIMIT 1
2023-05-09T07:06:52.159865Z	   11 Execute	SELECT `user_monster`.* FROM `user_monster` WHERE (`user_monster`.`user_id` = 2590299773) AND (`user_monster`.`monster_id` = 1) LIMIT 1
2023-05-09T07:06:52.160125Z	   11 Close stmt	
2023-05-09T07:06:52.160274Z	    8 Prepare	UPDATE `user_monster` SET `exp`=? WHERE `user_id`=? AND `monster_id`=?
2023-05-09T07:06:52.160372Z	    8 Execute	UPDATE `user_monster` SET `exp`=1 WHERE `user_id`=2590299773 AND `monster_id`=1
2023-05-09T07:06:52.160633Z	    8 Close stmt	
2023-05-09T07:06:52.160748Z	    8 Prepare	UPDATE `user_monster` SET `exp`=? WHERE `user_id`=? AND `monster_id`=?
2023-05-09T07:06:52.160825Z	    8 Execute	UPDATE `user_monster` SET `exp`=2 WHERE `user_id`=2590299773 AND `monster_id`=1
2023-05-09T07:06:52.161052Z	    8 Close stmt	
2023-05-09T07:06:52.161164Z	    8 Prepare	UPDATE `user_monster` SET `exp`=? WHERE `user_id`=? AND `monster_id`=?
2023-05-09T07:06:52.161243Z	    8 Execute	UPDATE `user_monster` SET `exp`=3 WHERE `user_id`=2590299773 AND `monster_id`=1
2023-05-09T07:06:52.161509Z	    8 Close stmt	
2023-05-09T07:06:52.161616Z	    8 Prepare	UPDATE `user_monster` SET `exp`=? WHERE `user_id`=? AND `monster_id`=?
2023-05-09T07:06:52.161692Z	    8 Execute	UPDATE `user_monster` SET `exp`=4 WHERE `user_id`=2590299773 AND `monster_id`=1
2023-05-09T07:06:52.161918Z	    8 Close stmt	
2023-05-09T07:06:52.162024Z	    8 Prepare	UPDATE `user_monster` SET `exp`=? WHERE `user_id`=? AND `monster_id`=?
2023-05-09T07:06:52.162103Z	    8 Execute	UPDATE `user_monster` SET `exp`=5 WHERE `user_id`=2590299773 AND `monster_id`=1
2023-05-09T07:06:52.162322Z	    8 Close stmt	
2023-05-09T07:06:52.162468Z	    8 Prepare	UPDATE `user_monster` SET `exp`=? WHERE `user_id`=? AND `monster_id`=?
2023-05-09T07:06:52.162548Z	    8 Execute	UPDATE `user_monster` SET `exp`=6 WHERE `user_id`=2590299773 AND `monster_id`=1
2023-05-09T07:06:52.162758Z	    8 Close stmt	
2023-05-09T07:06:52.162884Z	    8 Prepare	UPDATE `user_monster` SET `exp`=? WHERE `user_id`=? AND `monster_id`=?
2023-05-09T07:06:52.162969Z	    8 Execute	UPDATE `user_monster` SET `exp`=7 WHERE `user_id`=2590299773 AND `monster_id`=1
2023-05-09T07:06:52.163189Z	    8 Close stmt	
2023-05-09T07:06:52.163299Z	    8 Prepare	UPDATE `user_monster` SET `exp`=? WHERE `user_id`=? AND `monster_id`=?
2023-05-09T07:06:52.163398Z	    8 Execute	UPDATE `user_monster` SET `exp`=8 WHERE `user_id`=2590299773 AND `monster_id`=1
2023-05-09T07:06:52.163597Z	    8 Close stmt	
2023-05-09T07:06:52.163693Z	    8 Prepare	UPDATE `user_monster` SET `exp`=? WHERE `user_id`=? AND `monster_id`=?
2023-05-09T07:06:52.163733Z	    8 Execute	UPDATE `user_monster` SET `exp`=9 WHERE `user_id`=2590299773 AND `monster_id`=1
2023-05-09T07:06:52.163930Z	    8 Close stmt	
2023-05-09T07:06:52.164029Z	    8 Prepare	UPDATE `user_monster` SET `exp`=? WHERE `user_id`=? AND `monster_id`=?
2023-05-09T07:06:52.164114Z	    8 Execute	UPDATE `user_monster` SET `exp`=10 WHERE `user_id`=2590299773 AND `monster_id`=1
2023-05-09T07:06:52.164317Z	    8 Close stmt	
2023-05-09T07:06:52.164484Z	   11 Prepare	SELECT `user_item`.* FROM `user_item` WHERE (`user_item`.`user_id` = ?) AND (`user_item`.`item_id` IN (?,?))
2023-05-09T07:06:52.164566Z	   11 Execute	SELECT `user_item`.* FROM `user_item` WHERE (`user_item`.`user_id` = 2590299773) AND (`user_item`.`item_id` IN (1,2))
2023-05-09T07:06:52.164897Z	   11 Close stmt	
2023-05-09T07:06:52.165030Z	    8 Prepare	UPDATE `user_item` SET `count`=? WHERE `user_id`=? AND `item_id`=?
2023-05-09T07:06:52.165109Z	    8 Execute	UPDATE `user_item` SET `count`=999 WHERE `user_id`=2590299773 AND `item_id`=1
2023-05-09T07:06:52.165326Z	    8 Close stmt	
2023-05-09T07:06:52.165485Z	    8 Prepare	UPDATE `user_item` SET `count`=? WHERE `user_id`=? AND `item_id`=?
2023-05-09T07:06:52.165502Z	    8 Execute	UPDATE `user_item` SET `count`=999 WHERE `user_id`=2590299773 AND `item_id`=2
2023-05-09T07:06:52.165696Z	    8 Close stmt	
2023-05-09T07:06:52.165830Z	    8 Query	COMMIT
```
#### 提案手法を有効化した際のアクセス内容
```text
2023-05-10T03:02:33.819480Z	   11 Query	START TRANSACTION
2023-05-10T03:02:33.827235Z	   11 Prepare	INSERT INTO `user_coin` (`user_id`,`num`) VALUES (?,?)
2023-05-10T03:02:33.827469Z	   11 Execute	INSERT INTO `user_coin` (`user_id`,`num`) VALUES (2712740013,100000)
2023-05-10T03:02:33.828249Z	   11 Close stmt	
2023-05-10T03:02:33.829501Z	   11 Prepare	INSERT INTO `user_monster` (`user_id`,`monster_id`,`exp`) VALUES (?,?,?),(?,?,?),(?,?,?)
2023-05-10T03:02:33.829707Z	   11 Execute	INSERT INTO `user_monster` (`user_id`,`monster_id`,`exp`) VALUES (2712740013,1,0),(2712740013,2,0),(2712740013,3,0)
2023-05-10T03:02:33.830336Z	   11 Close stmt	
2023-05-10T03:02:33.831296Z	   11 Prepare	INSERT INTO `user_item` (`user_id`,`item_id`,`count`) VALUES (?,?,?),(?,?,?),(?,?,?)
2023-05-10T03:02:33.831488Z	   11 Execute	INSERT INTO `user_item` (`user_id`,`item_id`,`count`) VALUES (2712740013,1,1000),(2712740013,2,1000),(2712740013,3,1000)
2023-05-10T03:02:33.832040Z	   11 Close stmt	
2023-05-10T03:02:33.832146Z	   11 Query	COMMIT
2023-05-10T03:02:36.001713Z	   11 Prepare	SELECT * FROM user_monster WHERE user_id = ?
2023-05-10T03:02:36.001969Z	   11 Execute	SELECT * FROM user_monster WHERE user_id = 2712740013
2023-05-10T03:02:36.002561Z	   11 Close stmt	
2023-05-10T03:02:36.002990Z	   11 Prepare	SELECT * FROM user_item WHERE user_id = ?
2023-05-10T03:02:36.003182Z	   11 Execute	SELECT * FROM user_item WHERE user_id = 2712740013
2023-05-10T03:02:36.003543Z	   11 Close stmt	
2023-05-10T03:02:36.003716Z	   11 Prepare	SELECT * FROM user_coin WHERE user_id=?
2023-05-10T03:02:36.003879Z	   11 Execute	SELECT * FROM user_coin WHERE user_id=2712740013
2023-05-10T03:02:36.004242Z	   11 Close stmt	
2023-05-10T03:02:36.004326Z	   11 Query	START TRANSACTION
2023-05-10T03:02:36.004621Z	   11 Query	COMMIT
2023-05-10T03:02:38.548702Z	   11 Prepare	SELECT * FROM user_item WHERE user_id = ? AND item_id IN (?,?)
2023-05-10T03:02:38.548930Z	   11 Execute	SELECT * FROM user_item WHERE user_id = 2712740013 AND item_id IN (1,2)
2023-05-10T03:02:38.549400Z	   11 Close stmt	
2023-05-10T03:02:38.549597Z	   11 Prepare	SELECT * FROM user_monster WHERE user_id=? AND monster_id=?
2023-05-10T03:02:38.549762Z	   11 Execute	SELECT * FROM user_monster WHERE user_id=2712740013 AND monster_id=1
2023-05-10T03:02:38.550147Z	   11 Close stmt	
2023-05-10T03:02:38.550337Z	   11 Prepare	SELECT * FROM user_coin WHERE user_id=?
2023-05-10T03:02:38.550518Z	   11 Execute	SELECT * FROM user_coin WHERE user_id=2712740013
2023-05-10T03:02:38.550886Z	   11 Close stmt	
2023-05-10T03:02:38.551015Z	   11 Query	START TRANSACTION
2023-05-10T03:02:38.551346Z	   11 Query	UPDATE user_item SET count = case WHEN user_id=2712740013 AND item_id=1 THEN 984 WHEN user_id=2712740013 AND item_id=2 THEN 969 Else 0 End WHERE (user_id=2712740013 AND item_id=1) OR (user_id=2712740013 AND item_id=2)
2023-05-10T03:02:38.552040Z	   11 Query	UPDATE user_monster SET exp = case WHEN user_id=2712740013 AND monster_id=1 THEN 1710 Else 0 End WHERE (user_id=2712740013 AND monster_id=1)
2023-05-10T03:02:38.552536Z	   11 Query	UPDATE user_coin SET num = case WHEN user_id=2712740013 THEN 99700 Else 0 End WHERE (user_id=2712740013)
2023-05-10T03:02:38.552939Z	   11 Query	COMMIT
```



