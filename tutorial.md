# **Cloud Monitoring による SLO のモニタリング**

## **概要**

本ハンズオンでは、Google Cloud で SLO モニタリングを体験できます。

- コンテナのビルド・デプロイ
- 不具合バージョンのリリース
- ロールバック
- カナリアリリース

## Google Cloud プロジェクトの設定、確認

### **プロジェクトの課金が有効化されていることを確認する**

```bash
gcloud beta billing projects describe ${GOOGLE_CLOUD_PROJECT} | grep billingEnabled
```

**Cloud Shell の承認** という確認メッセージが出た場合は **承認** をクリックします。

出力結果の `billingEnabled` が **true** になっていることを確認してください。**false** の場合は、こちらのプロジェクトではハンズオンが進められません。別途、課金を有効化したプロジェクトを用意してからやり直してください。

## **環境準備**

<walkthrough-tutorial-duration duration=10></walkthrough-tutorial-duration>

最初に、ハンズオンを進めるための環境準備を行います。

下記の設定を進めていきます。

- gcloud コマンドラインツール設定
- Google Cloud 機能（API）有効化設定

## **gcloud コマンドラインツール**

Google Cloud は、コマンドライン（CLI）、GUI から操作が可能です。

### **1. gcloud コマンドラインツールとは？**

このツールを使用するとコマンドラインから、またはスクリプトや他の自動化により、多くの一般的な運用タスクを実行できます。

**ヒント**: gcloud についての詳細は[こちら](https://cloud.google.com/sdk/gcloud?hl=ja)をご参照ください

### **2. gcloud からの Cloud Run のデフォルト設定**

Cloud Run の利用するリージョン、プラットフォームのデフォルト値を設定します。

```bash
gcloud config set run/region asia-northeast1
gcloud config set run/platform managed
```

ここではリージョンを東京、プラットフォームをフルマネージドに設定しました。この設定を行うことで、gcloud コマンドから Cloud Run を操作するときに毎回指定する必要がなくなります。

### **3. ユーザ ID の設定**

1 つのプロジェクトでも複数人でリソースが競合しないよう、個人の ID を設定しておきます。

```bash
export user_id="$( git config user.email  | awk '{ split($0, a, "@"); print a[1] }' )"
```

## **参考: Cloud Shell の接続が途切れてしまったときは?**

一定時間非アクティブ状態になる、またはブラウザが固まってしまったなどで `Cloud Shell` の接続が切れてしまう場合があります。

その場合は `再接続` をクリックした後、以下の対応を行い、チュートリアルを再開してください。

![再接続画面](https://raw.githubusercontent.com/GoogleCloudPlatform/gcp-getting-started-lab-jp/master/workstations_with_generative_ai/images/reconnect_cloudshell.png)

### **1. チュートリアルを開く**

```bash
teachme ~/slo-monitoring/tutorial.md
```

### **2. ユーザ ID をセット**

```bash
export user_id="$( git config user.email  | awk '{ split($0, a, "@"); print a[1] }' )"
```

途中まで進めていたチュートリアルのページまで `次へ` ボタンを押し、進めてください。

## **Google Cloud 環境設定**

Google Cloud では利用したい機能（API）ごとに、有効化を行う必要があります。  
ここでは、以降のハンズオンで利用する機能を事前に有効化しておきます。

```bash
gcloud services enable compute.googleapis.com run.googleapis.com cloudbuild.googleapis.com artifactregistry.googleapis.com
```

**GUI**: [API ライブラリ](https://console.cloud.google.com/apis/library)

<walkthrough-footnote>必要な機能が使えるようになりました。次に実際に Cloud Run にアプリケーションをデプロイする方法を学びます。</walkthrough-footnote>

## **Cloud Run へのデプロイ**

<walkthrough-tutorial-duration duration=15></walkthrough-tutorial-duration>

### **準備**

下記のように GUI を操作し Cloud Run の管理画面を開いておきましょう。

<walkthrough-spotlight-pointer spotlightId="console-nav-menu">ナビゲーションメニュー</walkthrough-spotlight-pointer> -> サーバーレス -> Cloud Run

また以降の手順で Cloud Run の管理画面は何度も開くことになるため、ピン留め (Cloud Run メニューにマウスオーバーし、ピンのアイコンをクリック) しておくと便利です。

### **1. リポジトリを作成**

```bash
gcloud artifacts repositories create "apps-${user_id}" --repository-format "docker" --location "asia-northeast1" --description "Docker repository for ${user_id}"
```

### **2. アプリのビルド & プッシュ**

Cloud Build を使い、一連の操作を一気に行います。 `--pack` オプションを指定することで [Buildpacks](https://github.com/GoogleCloudPlatform/buildpacks) が内部的に利用され、Dockerfile なしにコンテナをビルドできます。

```bash
gcloud builds submit --pack "builder=gcr.io/buildpacks/builder,image=asia-northeast1-docker.pkg.dev/${GOOGLE_CLOUD_PROJECT}/apps-${user_id}/sample:v0.1" src
```

### **3. Cloud Run にデプロイ**

```bash
gcloud run deploy svc-${user_id} --image "asia-northeast1-docker.pkg.dev/${GOOGLE_CLOUD_PROJECT}/apps-${user_id}/sample:v0.1" --allow-unauthenticated
```

## **不具合バージョンのリリース**

### **1. 新リビジョンのデプロイ**

100% の確率でレスポンスコード 403 を返す環境変数を設定し、新しいリビジョンをデプロイします。

```bash
gcloud run deploy svc-${user_id} --image "asia-northeast1-docker.pkg.dev/${GOOGLE_CLOUD_PROJECT}/apps-${user_id}/sample:v0.1" --set-env-vars "ERROR_RATE=1.0"
```

### **2. アプリケーションにアクセス**

デプロイ後、どのようなレスポンスが返ってくるかを確認します。 `Forbidden` と返ってくるはずです。

```bash
curl -iXGET $(gcloud run services describe svc-${user_id} --format json | jq -r '.status.address.url')
```

## **ロールバック**

### **1. 旧リビジョンのロールバック**

不具合のなかった前のリビジョンにもどします。

```bash
OLD_REV=$(gcloud run revisions list --format json | jq -r '.[].metadata.name' | grep "svc-${user_id}-" | sort -r | sed -n 2p)
gcloud run services update-traffic svc-${user_id} --to-revisions=${OLD_REV}=100
```

### **2. アプリケーションにアクセス**

デプロイ後、どのようなレスポンスが返ってくるかを確認します。正しい応答が返ってくるはずです。

```bash
curl -iXGET $(gcloud run services describe svc-${user_id} --format json | jq -r '.status.address.url')
```

## **カナリアリリース**

カナリアリリースは新リビジョンをトラフィックを流さない状態でデプロイし、徐々にトラフィックを流すように設定することで実現します。

### **1. 新リビジョンのデプロイ**

半分の確率でレスポンスコード 403 を返す環境変数を添えて、新しいリビジョンをデプロイします。

```bash
gcloud run deploy svc-${user_id} --image "asia-northeast1-docker.pkg.dev/${GOOGLE_CLOUD_PROJECT}/apps-${user_id}/sample:v0.1" --set-env-vars "ERROR_RATE=0.5" --no-traffic
```

**ヒント**: 新リビジョンにトラフィックを流さないよう、`--no-traffic` のオプションをつけています。これがない場合、デプロイされた瞬間にすべてのトラフィックが新リビジョンに流れます。

### **2. カナリアリリース**

以下のコマンドで新リビジョンに 20 %, 現行リビジョンに 80 % のアクセスを割り振ります。

```bash
gcloud run services update-traffic svc-${user_id} --to-revisions "LATEST=20"
```

ターミナルに出力された URL をクリックするとブラウザが開きます。そこでリロードを繰り返してみます。10 回に 1 回 `Forbidden` と表示されます。

### **3.すべてのアクセスを新リビジョンに**

状況的に正しい処置ではなく、エラーレートが悪化することになりますが、すべて最新のリビジョンに流してみます。

```bash
gcloud run services update-traffic svc-${user_id} --to-latest
```

2 回に 1 回 `Forbidden` と表示されます。

## **Congraturations!**

<walkthrough-conclusion-trophy></walkthrough-conclusion-trophy>

以上で SLO モニタリングの学習は完了です。
