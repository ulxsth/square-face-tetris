# square-face-tetris
物理エンジンをつけたテトリスのミノを人の顔にしたゲーム

# build
wasm としてビルドする。

```sh
GOOS=js GOARCH=wasm go build -o dist/main.wasm app/main.go
```