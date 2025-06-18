<script lang="ts" module>
	import { Badge } from '$lib/components/ui/badge';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import Icon from '@iconify/svelte';
	import type { Row } from '@tanstack/table-core';
	import type { User } from './types';

	export const cells = {
		_row_picker: _row_picker,
		username: username,
		tenant: tenant,
		fullName: fullName,
		emailAddress: emailAddress,
		suspended: suspended,
		maximumBuckets: maximumBuckets,
		capacityLimit: capacityLimit,
		objectLimit: objectLimit
	};
</script>

{#snippet _row_picker(row: Row<User>)}
	<Checkbox
		checked={row.getIsSelected()}
		onCheckedChange={(value) => row.toggleSelected(!!value)}
		class="border-secondary-950"
		aria-label="Select row"
	/>
{/snippet}

{#snippet username(row: Row<User>)}
	<p>{row.original.userId}</p>
{/snippet}

{#snippet tenant(row: Row<User>)}
	<Badge variant="outline">{row.original.tenant}</Badge>
{/snippet}

{#snippet fullName(row: Row<User>)}
	<p>{row.original.fullName}</p>
{/snippet}

{#snippet emailAddress(row: Row<User>)}
	<p>{row.original.emailAddress}</p>
{/snippet}

{#snippet suspended(row: Row<User>)}
	<p class="flex justify-end">
		{#if row.original.suspended}
			<Icon icon="ph:x" class="text-destructive" />
		{:else}
			<Icon icon="ph:check" class="text-success" />
		{/if}
	</p>
{/snippet}

{#snippet maximumBuckets(row: Row<User>)}
	<p class="text-right">{row.original.maximumBuckets}</p>
{/snippet}

{#snippet capacityLimit(row: Row<User>)}
	<p class="text-right">{row.original.capacityLimit}</p>
{/snippet}

{#snippet objectLimit(row: Row<User>)}
	<p class="text-right">{row.original.objectLimit}</p>
{/snippet}
