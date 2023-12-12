SLO monitoring
---

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
