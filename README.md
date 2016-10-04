# infrakit-plugin-sakuracloud


infrakitのインスタンスプラグイン サンプル実装です。
さくらのクラウド上にインスタンスを作成します。

**サンプル実装です。最低限の動きは実装していますが、実用には向きません。**

フレーバープラグインに`vanilla`を使う場合のサンプルは以下の通りです。
事前に以下の環境変数を設定しておいてください。

  - SAKURACLOUD_ACCESS_TOKEN
  - SAKURACLOUD_ACCESS_TOKEN_SECRET

```
$ infrakit/cli group --name group watch <<EOF
{
    "ID": "cattle",
    "Properties": {
        "Instance": {
            "Plugin": "instance-sakuracloud",
            "Properties": {
                "SourceArchiveID": 112800673084 ,
                "Zone": "tk1a",
                "Password": "Put Your Password"
            }
        },
        "Flavor": {
            "Plugin": "flavor-vanilla",
            "Properties": {
                "Size": 1,
                "UserData": [
                    "sudo apt-get update -y",
                    "sudo apt-get install -y nginx",
                    "sudo service nginx start"
                ],

                "Labels": {
                    "tier": "web",
                    "project": "infrakit"
                }
            }
        }
    }
}
EOF
```
