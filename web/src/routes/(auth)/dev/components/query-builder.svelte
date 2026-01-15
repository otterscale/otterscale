<script lang="ts" module>
	export type FilterOperator = 'equals' | 'contains' | 'greater than' | 'less than';
	export type LogicOperator = 'and' | 'or';
	export type FilterRule = {
		id: string;
		operator: FilterOperator;
		value: any;
	};
	export type FilterGroup = {
		logic: LogicOperator;
		filters: (FilterRule | FilterGroup)[];
	};
	export const checkFilter = (rowValue: any, filter: FilterRule | FilterGroup): boolean => {
		// 如果是群組，遞迴處理
		if ('logic' in filter) {
			if (filter.logic === 'and') {
				return filter.filters.every((f) => checkFilter(rowValue, f));
			} else {
				// OR
				return filter.filters.some((f) => checkFilter(rowValue, f));
			}
		}
		// 如果是單一規則，執行比對
		const { id, operator, value } = filter;
		const cellValue = rowValue[id]; // 取得該 Row 對應欄位的值
		switch (operator) {
			case 'contains':
				return String(cellValue).toLowerCase().includes(String(value).toLowerCase());
			case 'equals':
				return cellValue == value;
			case 'greater than':
				return Number(cellValue) > Number(value);
			case 'less than':
				return Number(cellValue) < Number(value);
			default:
				return true;
		}
	};
</script>

<script lang="ts">
	// 這是雙向綁定的資料節點，可以是 Group 或 Rule
	export let node: FilterGroup | FilterRule;
	// 如果是根節點，不需要移除按鈕
	export let isRoot: boolean = false;
	// 用來通知父層刪除自己
	export let onDelete: () => void = () => {};
	// 判斷是否為 Group 的 Type Guard
	const isGroup = (n: any): n is FilterGroup => 'logic' in n;
	function addRule() {
		if (isGroup(node)) {
			node.filters = [
				...node.filters,
				{ id: 'Name', operator: 'equals', value: '' } // 預設空規則
			];
		}
	}
	function addGroup() {
		if (isGroup(node)) {
			node.filters = [
				...node.filters,
				{ logic: 'and', filters: [] } // 預設新群組
			];
		}
	}
	function removeChild(index: number) {
		if (isGroup(node)) {
			node.filters = node.filters.filter((_, i) => i !== index);
			node.filters = node.filters;
		}
	}
</script>

<div class="flex flex-col gap-2 rounded-lg border p-3 shadow-sm">
	{#if isGroup(node)}
		<div class="mb-2 flex items-center gap-2 rounded p-2">
			<select class="rounded border px-2 py-1 text-sm font-bold" bind:value={node.logic}>
				<option value="and">and</option>
				<option value="or">or</option>
			</select>
			<div class="ml-auto flex gap-1">
				<button
					class="rounded bg-green-600 px-2 py-1 text-xs text-white hover:bg-green-700"
					on:click={addRule}
				>
					+ Rule
				</button>
				<button
					class="rounded bg-indigo-600 px-2 py-1 text-xs text-white hover:bg-indigo-700"
					on:click={addGroup}
				>
					+ Group
				</button>
				{#if !isRoot}
					<button
						class="rounded bg-red-500 px-2 py-1 text-xs text-white hover:bg-red-600"
						on:click={onDelete}
					>
						Remove Group
					</button>
				{/if}
			</div>
		</div>
		<div class="flex flex-col gap-2 border-l-2 border-gray-300 pl-6">
			{#each node.filters as child, index (index)}
				<svelte:self bind:node={child} onDelete={() => removeChild(index)} />
			{/each}
		</div>
	{:else}
		<div class="flex items-center gap-2">
			<span class="text-sm text-gray-400">Field:</span>
			<select class="rounded border px-2 py-1 text-sm" bind:value={node.id}>
				<option value="Name">Name</option>
				<option value="Namespace">Namespace</option>
			</select>
			<select class="rounded border px-2 py-1 text-sm" bind:value={node.operator}>
				<option value="equals">equals</option>
				<option value="contains">contains</option>
				<option value="greater than">greater than</option>
				<option value="less than">less than</option>
			</select>
			<input
				class="w-32 rounded border px-2 py-1 text-sm"
				type="text"
				bind:value={node.value}
				placeholder="Value..."
			/>
			<button
				class="ml-auto px-2 font-bold text-red-500 hover:text-red-700"
				on:click={onDelete}
				aria-label="Remove rule"
			>
				✕
			</button>
		</div>
	{/if}
</div>
