SLO monitoring
---

## SQL モニタリング チュートリアル

Cloud Shell で以下のコマンドを実行してください。

```sh
git clone https://github.com/pottava/slo-monitoring.git ~/slo
teachme ~/slo/tutorial.md
```

## サンプル アプリケーションのローカル起動

```sh
go run src/main.go
```

### Websocket アプリなどを扱う場合

Google のプロキシーが Websocket 未対応のため  
ローカルマシンから TCP トンネルを確立することで動作を確認できます。

Cloud Shell の場合

```sh
app_port=8080
gcloud alpha cloud-shell ssh -- -nNT -L "${app_port}:localhost:${app_port}"
```

Cloud Workstations の場合

```sh
cluster_name=
config_name=
workstation_name=
app_port=8080

gcloud workstations start-tcp-tunnel --cluster "${cluster_name}" \
    --config "${config_name}" --region asia-northeast1 \
    "${workstation_name}" "${app_port}" \
    --local-host-port ":${app_port}"
```
