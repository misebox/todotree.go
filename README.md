# todotree


## タスクランナーについて

makeの代わりにtaskコマンドを使います。 

- 定義ファイル: [Taskfile.yml](./Taskfile.yml)

- Taskコマンドの詳細: [taskfile.dev](https://taskfile.dev)

例えばコンテナのログを表示するには次のようにします。
```
$ task logs
task: [logs] docker compose logs --tail=100 -f
todotreego-api-1  | 
todotreego-api-1  |   __    _   ___  
todotreego-api-1  |  / /\  | | | |_) 
todotreego-api-1  | /_/--\ |_| |_| \_ , built with Go 
todotreego-api-1  | 
...
...
```


### 各タスクコマンドについて
```
$ task -l
task: Available tasks for this project:
* build:             Build docker image to deploy
* build-local:       Build docker image to development
* down:              down
* logs:              logs
* migrate:           Migrate local DB
* ps:                ps
* test:              Run Test
* up:                up with hot relaod
```

### デフォルトタスク

`task test` と `task up` を同時に実行し、正常に完了したら `task logs` を実行します

```
$ task
task: [test] go test -race -shuffle=on -v ./...
...
[+] Running 2/2                                    
 ⠿ Container todotreego-api-1  Started        0.6s 
 ⠿ Container todo-db           Running        0.0s
...
todotreego-api-1  | running...
todotreego-api-1  | 2022/11/27 08:24:11 start with: http://[::]:8001

```
