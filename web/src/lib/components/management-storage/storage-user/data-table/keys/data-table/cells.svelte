<script lang="ts" module>
	import type { User_Key } from '$gen/api/storage/v1/storage_pb';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import type { Row } from '@tanstack/table-core';

	export const cells = {
		_row_picker: _row_picker,
		access: access
	};
</script>

{#snippet _row_picker(row: Row<User_Key>)}
	<Checkbox
		checked={row.getIsSelected()}
		onCheckedChange={(value) => row.toggleSelected(!!value)}
		class="border-secondary-950"
		aria-label="Select row"
	/>
{/snippet}

{#snippet access(row: Row<User_Key>)}
	{row.original.accessKey}
{/snippet}
