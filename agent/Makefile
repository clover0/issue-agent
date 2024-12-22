# Go関連のターゲット
.PHONY: test test-coverage clean

# テストの実行
test:
	go test ./... -v

# テストカバレッジの計測と結果の出力
test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

# 生成されたファイルの削除
clean:
	rm -f coverage.out coverage.html
