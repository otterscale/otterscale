<script lang="ts" module>
	import * as Progress from '$lib/components/custom/progress';
	import { Badge } from '$lib/components/ui/badge';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import type { Row } from '@tanstack/table-core';
	import type { Pool } from './types';

	export const cells = {
		_row_picker: _row_picker,
		name: name,
		dataProtection: dataProtection,
		applications: applications,
		PGStatus: PGStatus,
		usage: usage
	};
</script>

{#snippet _row_picker(row: Row<Pool>)}
	<Checkbox
		checked={row.getIsSelected()}
		onCheckedChange={(value) => row.toggleSelected(!!value)}
		class="border-secondary-950"
		aria-label="Select row"
	/>
{/snippet}

{#snippet name(row: Row<Pool>)}
	{row.original.name}
{/snippet}

{#snippet dataProtection(row: Row<Pool>)}
	<Badge variant="outline">
		{row.original.dataProtection}
	</Badge>
{/snippet}

{#snippet applications(row: Row<Pool>)}
	<Badge variant="outline">
		{row.original.applications}
	</Badge>
{/snippet}

{#snippet PGStatus(row: Row<Pool>)}
	{row.original.PGStatus}
{/snippet}

{#snippet usage(row: Row<Pool>)}
	<Progress.Root numerator={row.original.usage} denominator={100}>
		{#snippet ratio({ numerator, denominator })}
			{(numerator * 100) / denominator}%
		{/snippet}
	</Progress.Root>
{/snippet}
