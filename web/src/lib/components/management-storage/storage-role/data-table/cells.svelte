<script lang="ts" module>
	import { Badge } from '$lib/components/ui/badge';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import type { Row } from '@tanstack/table-core';
	import type { Role } from './types';

	export const cells = {
		_row_picker: _row_picker,
		roleName: roleName,
		path: path,
		arn: arn,
		createDate: createDate,
		maximumSessionDuration: maximumSessionDuration
	};
</script>

{#snippet _row_picker(row: Row<Role>)}
	<Checkbox
		checked={row.getIsSelected()}
		onCheckedChange={(value) => row.toggleSelected(!!value)}
		class="border-secondary-950"
		aria-label="Select row"
	/>
{/snippet}

{#snippet roleName(row: Row<Role>)}
	{row.original.roleName}
{/snippet}

{#snippet path(row: Row<Role>)}
	<Badge variant="outline">{row.original.path}</Badge>
{/snippet}

{#snippet arn(row: Row<Role>)}
	<p class="font-mono text-sm">{row.original.arn}</p>
{/snippet}

{#snippet createDate(row: Row<Role>)}
	{row.original.createDate.toLocaleDateString()}
{/snippet}

{#snippet maximumSessionDuration(row: Row<Role>)}
	<p class="text-right">{row.original.maximumSessionDuration}</p>
{/snippet}
