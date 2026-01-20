<script lang="ts">
	// 需安裝 jexl、jsep，並在 main.js 引入主題 css
	import Jexl from 'jexl';
	import jsep from 'jsep';
	import { tick } from 'svelte';
	import Monaco from 'svelte-monaco';

	let filter = '';
	let ghost = '';
	let isValid = true;

	// 假設可選的自動補全詞
	const completions = ['Name', 'Namespace', 'Labels', 'Annotations', 'Configurations'];

	// ghost autocomplete 實現（根據目前輸入給出建議）
	function computeGhost(value: string) {
		const input = document.activeElement as HTMLInputElement;
		if (!input) return '';
		const cursor = input.selectionStart ?? value.length;

		// 取得游標前的內容
		const before = value.slice(0, cursor);
		const after = value.slice(cursor);

		// 找出游標前的單字起點
		const startMatch = before.match(/(\w+)$/);
		const wordStart = startMatch ? cursor - startMatch[1].length : cursor;

		// 找出游標後的單字終點
		const endMatch = after.match(/^(\w+)/);
		const wordEnd = endMatch ? cursor + endMatch[1].length : cursor;

		// 取得游標所在單字
		const word = value.slice(wordStart, wordEnd);
		if (!word) return '';

		const candidate = completions.find(
			(c) => c.toLowerCase().startsWith(word.toLowerCase()) && c !== word
		);
		if (candidate) {
			// 直接顯示完整補全詞（用於取代整個單字）
			return candidate.slice(0);
		}
		return '';
	}

	// 修改 acceptGhost，直接取代游標所在單字
	function acceptGhost(e: KeyboardEvent) {
		if (e.key === 'Tab' && ghost) {
			e.preventDefault();
			const input = e.target as HTMLInputElement;
			const cursor = input.selectionStart ?? filter.length;

			const before = filter.slice(0, cursor);
			const after = filter.slice(cursor);

			const startMatch = before.match(/(\w+)$/);
			const wordStart = startMatch ? cursor - startMatch[1].length : cursor;

			const endMatch = after.match(/^(\w+)/);
			const wordEnd = endMatch ? cursor + endMatch[1].length : cursor;

			// 取代游標所在單字
			filter = filter.slice(0, wordStart) + ghost + filter.slice(wordEnd);

			// 移動游標到補全詞尾
			tick().then(() => {
				input.setSelectionRange(wordStart + ghost.length, wordStart + ghost.length);
			});

			ghost = '';
			checkFilter(filter);
		}
	}

	// 檢查 filter 合法性
	function checkFilter(expr: string) {
		let valid = true;
		try {
			Jexl.compile(expr);
		} catch {
			valid = false;
		}
		isValid = valid;
	}

	// 觸發 ghost 補全 & 語法檢查
	function onInput(e: Event) {
		filter = (e.target as HTMLInputElement).value;
		ghost = computeGhost(filter);
		checkFilter(filter);
	}

	// // 按 Tab 自動接受 ghost 補全
	// function acceptGhost(e: KeyboardEvent) {
	// 	if (e.key === 'Tab' && ghost) {
	// 		e.preventDefault();
	// 		filter += ghost;
	// 		ghost = '';
	// 		checkFilter(filter);
	// 	}
	// }

	// 高亮展示 (用 jsep parse AST 為成功則致上色，demo 簡化處理)
	function highlight(expr: string): string {
		try {
			const ast = jsep(expr);
			// 簡單: 把字串做如下替換
			let html = expr.replace(
				/(==|!=|>|<|>=|<=|and|or)/g,
				'<span class="text-blue-500 font-bold">$1</span>'
			);
			html = html.replace(/(["][^"]*["]|'[^']*')/g, '<span class="text-green-600">$1</span>');
			return html;
		} catch {
			return expr;
		}
	}
</script>

<!-- 樣式仿 shadcn (tailwind or 自訂 css) -->
<div class="mx-auto mt-10 w-full max-w-lg">
	<label class="mb-2 block text-sm font-medium text-gray-700">Filter Expression</label>
	<div class="group relative">
		<!-- Ghost autocomplete input -->
		<input
			type="text"
			class="w-full rounded-xl border border-gray-300 bg-white px-4 py-2 font-mono text-lg
             shadow-sm transition-colors focus:border-blue-500 focus:ring-2 focus:outline-none"
			bind:value={filter}
			placeholder="輸入 filter（例如 age > 18）"
			autocomplete="off"
			spellcheck="false"
			on:input={onInput}
			on:keydown={acceptGhost}
		/>
		<!-- ghost: 以灰色疊在原始文字之後作虛影 -->
		{#if ghost}
			<span
				class="pointer-events-none absolute top-1/2 left-4 font-mono text-lg text-gray-400 opacity-60 select-none"
				style="transform: translateY(-50%);"
				aria-hidden="true"
			>
				<span style="visibility:hidden">{filter}</span>{ghost}
			</span>
		{/if}
	</div>
	<!-- filter合法性提示 -->
	{#if filter !== ''}
		<p class="mt-1 text-sm {isValid ? 'text-green-600' : 'text-red-500'}">
			{isValid ? 'filter 語法正確，可以套用搜尋' : 'filter 語法錯誤，請檢查表達式'}
		</p>
	{/if}
	<!-- 高亮展示 -->

	<div
		class="mt-4 h-full rounded border bg-gray-50 px-4 py-2 font-mono text-base"
		style="min-height:2.5em;"
	>
		<!-- {@html highlight(filter)} -->
	</div>
</div>
<div class="min-h-48 border">
	<Monaco
		options={{
			language: 'javascript',
			automaticLayout: true,
			padding: { top: 8, bottom: 8 },
			overviewRulerBorder: false,
			hideCursorInOverviewRuler: true
		}}
		value={filter}
		theme="vs-dark"
	/>
</div>
