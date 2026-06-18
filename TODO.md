# TODO.md

go-netconf プロジェクトの改善タスク管理ファイル。

---

## Phase 1: デリミタ欠損問題の再現・対策・エビデンス取得

**統合ブランチ**: `phase/1-delimiter-split-test-env`
**目的**: TCPセグメント分割によるデリミタ欠損問題の再現環境整備・対策実装・エビデンス取得

### 進捗サマリー

| Issue | タイトル | ステータス | PR |
|-------|---------|-----------|-----|
| #5 | REQ-4: デリミタ欠損の再現環境整備（ループ可視化含む） | 進行中 | - |
| #6 | REQ-5: デリミタ欠損への対策実装 | 未着手（着手時に設計方式を相談） | - |
| #7 | REQ-6: デリミタ欠損対策のエビデンス取得 | 未着手 | - |
| #8 | REQ-8: TestWaitForBytesEmpty 既存テスト失敗修正 | 未着手（REQ-4実装中に発覚、別Issue管理） | - |

### 背景・課題

`netconf/transport.go` の `WaitForFunc` において、TCPセグメント分割によって
デリミタ `]]>]]>` が受信欠損した場合、forループが終了しない課題が発覚。
現状は再現環境が存在しないため、再現条件の整備から着手する。

### REQ-4: デリミタ欠損の再現環境整備 (Issue #5)

- [ ] 対象ファイル特定・実装方針確定（方針A: Goユニットテスト / 方針B: netconf-mock API追加）
- [ ] `WaitForFunc` にイテレーションごとのデバッグログを追加（ループ空回りの可視化）
  - iteration カウンタ追加、デリミタ未発見時に `iteration=N read=M bytes total=K delimiter not found` を出力
- [ ] 再現環境実装（デリミタなしでデータを送信するシナリオ）
- [ ] 再現時のログ確認・手順ドキュメント化
- [ ] feature/req-5-delimiter-split-repro ブランチで PR 作成 → phase/1 にマージ

### REQ-6: デリミタ欠損への対策実装 (Issue #6)

- [ ] **着手時に設計方式を相談する**（方針A/B/C から選定）
- [ ] 設計方式を Issue コメントに記載
- [ ] 対策実装
- [ ] feature/req-6-delimiter-split-fix ブランチで PR 作成 → phase/1 にマージ

### REQ-8: TestWaitForBytesEmpty 既存テスト失敗修正 (Issue #8)

**発覚経緯**: REQ-4（Issue #5）実装中に発見。PR #118（e572641）が `WaitForFunc` の EOF 処理を変更した際に壊れ、CI未実行のため長期間気づかれなかった。

- [ ] `WaitForFunc` の `n==0` 判定に `out.Len()==0 && pos==0` ガードを追加（B案）
- [ ] `TestWaitForBytesEmpty` がパスすることを確認
- [ ] feature/req-8-fix-wait-for-bytes-empty ブランチで PR 作成 → phase/1 にマージ

### REQ-7: デリミタ欠損対策のエビデンス取得 (Issue #7)

- [ ] 対策前後の動作比較ログ取得
- [ ] 正常系テスト全件パスのエビデンス取得
- [ ] docs/ 配下に検証レポートまとめ
- [ ] feature/req-7-delimiter-split-evidence ブランチで PR 作成 → phase/1 にマージ

---

## 完了済みタスク

| Issue | タイトル | PR | 完了日 |
|-------|---------|-----|--------|
| #1 | Add debug logging and netconf-mock timeout reproduction environment | - | 2026-06-15 |
| #3 | REQ-3: netconf-client のタイムアウト機構を time.After() ゴルーチンパターンに置き換える | - | 2026-06-16 |
