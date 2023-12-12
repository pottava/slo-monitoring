SLO monitoring
---

## ローカル開発

### アプリケーションの依存解決と起動

```sh
cd src
pip install poetry
export PATH="$HOME/.local/bin:$PATH"
poetry install --no-root
poetry run streamlit run app.py
```

Cloud Shell で試している場合は、Google のプロキシーが Websocket 未対応のため  
ローカルマシンから TCP トンネルを確立することで動作を確認できます。

```sh
gcloud alpha cloud-shell ssh -- -nNT -L 8501:localhost:8501
```

Cloud Workstations の場合も同様に、Google のプロキシーが Websocket 未対応のため  
ローカルマシンから TCP トンネルを確立することで動作を確認できます。

```sh
cluster_name=
config_name=
workstation_name=
streamlit_port=8501

gcloud workstations start-tcp-tunnel --cluster "${cluster_name}" \
    --config "${config_name}" --region asia-northeast1 \
    "${workstation_name}" "${streamlit_port}" \
    --local-host-port ":${streamlit_port}"
```
