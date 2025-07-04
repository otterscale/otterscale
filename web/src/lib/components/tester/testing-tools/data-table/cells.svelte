<script lang="ts" module>
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import type { Row } from '@tanstack/table-core';
	import { formatTimeAgo } from '$lib/formatter';
	import { TestResult_Type, type TestResult } from '$gen/api/bist/v1/bist_pb'

	export const cells = {
		_row_picker: _row_picker,
		type: type,
		name: name,
		input: input,
	};
</script>

{#snippet _row_picker(row: Row<TestResult>)}
	<Checkbox
		checked={row.getIsSelected()}
		onCheckedChange={(value) => row.toggleSelected(!!value)}
		class="border-secondary-950"
		aria-label="Select row"
	/>
{/snippet}

{#snippet type(row: Row<TestResult>)}
	<p>
		{TestResult_Type[row.original.type]}
	</p>
{/snippet}

{#snippet name(row: Row<TestResult>)}
	<p>
		{row.original.name}
	</p>
{/snippet}

{#snippet input(row: Row<TestResult>)}
	<p>
		{row.original.input.case}
	</p>
{/snippet}

<!-- {#snippet startTime(row: Row<TestResult>)}
	<p>
		{formatTimeAgo(row.original.startTime)}
	</p>
{/snippet} -->
