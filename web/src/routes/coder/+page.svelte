<script lang="ts">
	import { onMount } from 'svelte';

	let editorContainer: HTMLDivElement = $state({} as HTMLDivElement);

	let value = `\
# ---------------------------------------------------------
# 測試用的伺服器配置 YAML
# ---------------------------------------------------------

server:
  host: "127.0.0.1"
  port: 8080
  timeout: 30.5
  enable_logging: true
  environment: production

# 測試陣列與嵌套物件
database:
  driver: postgres
  endpoints:
    - "db-primary.internal"
    - "db-replica.internal"
  credentials:
    username: admin
    # 這裡可以測試你的 Marker，例如故意輸入錯誤格式
    password: p@ssw0rd

# 測試多行字串 (Literal Block)
description: |
  這是一段多行文字測試。
  Monaco 應該要能正確處理縮進。
  YAML 著色器通常會給予不同的顏色。

# 測試複雜列表
services:
  - id: 1
    name: "Auth Service"
    tags: [web, security, oauth2]
  - id: 2
    name: "Payment Gateway"
    tags:
      - financial
      - stripe
      - api

# 故意留一個縮進錯誤的範例（你可以取消註解來測試錯誤提示）
#  error_test:
#  bad_indent: true
    `;

	onMount(async () => {
		const monaco = await import('monaco-editor');

		// 1. 初始化編輯器
		const editor = monaco.editor.create(editorContainer, {
			value: value,
			language: 'yaml',
			theme: 'vs-dark'
		});

		// 2. 監聽內容變更 (雙向綁定)
		editor.onDidChangeModelContent(() => {
			value = editor.getValue();
		});

		// 3. 測試：手動設置 Marker (語法提示)
		const model = editor.getModel();
		if (model) {
			monaco.editor.setModelMarkers(model, 'owner', [
				{
					startLineNumber: 1,
					startColumn: 1,
					endLineNumber: 10,
					endColumn: 10,
					message: '自定義語法警告',
					severity: monaco.MarkerSeverity.Warning
				}
			]);
		}
	});
</script>

<div class="editor-wrapper">
	<div bind:this={editorContainer} class="h-screen w-screen" />
</div>
