<script lang="ts" module>
	import { page } from '$app/state';
	import type { Application } from '$lib/api/application/v1/application_pb';
	import { RowPickers } from '$lib/components/custom/data-table';
	import * as Progress from '$lib/components/custom/progress/index.js';
	import { Badge } from '$lib/components/ui/badge';
	import Button from '$lib/components/ui/button/button.svelte';
	import Icon from '@iconify/svelte';
	import type { Row } from '@tanstack/table-core';

	export const cells = {
		_row_picker,
		name,
		type,
		namespace,
		health,
		service,
		pod,
		replica,
		container,
		volume,
		nodeport
	};
</script>

{#snippet _row_picker(row: Row<Application>)}
	<RowPickers.Cell {row} />
{/snippet}

{#snippet name(row: Row<Application>)}
	<span class="flex items-center">
		{row.original.name}
		<Button variant="ghost" href={`${page.url}/${row.original.namespace}/${row.original.name}`}>
			<Icon icon="ph:arrow-square-out" />
		</Button>
	</span>
{/snippet}

{#snippet type(row: Row<Application>)}
	<Badge variant="outline">
		{row.original.type}
	</Badge>
{/snippet}

{#snippet namespace(row: Row<Application>)}
	{row.original.namespace}
{/snippet}

{#snippet health(row: Row<Application>)}
	<Progress.Root
		numerator={Number(row.original.healthies)}
		denominator={Number(row.original.pods.length)}
	>
		{#snippet ratio({ numerator, denominator })}
			{Progress.formatRatio(numerator, denominator)}
		{/snippet}
		{#snippet detail({ numerator, denominator })}
			{numerator}/{denominator}
		{/snippet}
	</Progress.Root>
{/snippet}

{#snippet service(row: Row<Application>)}
	<span class="flex justify-end">
		{row.original.services.length}
	</span>
{/snippet}

{#snippet pod(row: Row<Application>)}
	<span class="flex justify-end">
		{row.original.pods.length}
	</span>
{/snippet}

{#snippet replica(row: Row<Application>)}
	<span class="flex justify-end">
		{row.original.replicas}
	</span>
{/snippet}

{#snippet container(row: Row<Application>)}
	<span class="flex justify-end">
		{row.original.containers.length}
	</span>
{/snippet}

{#snippet volume(row: Row<Application>)}
	<span class="flex justify-end">
		{row.original.persistentVolumeClaims.length}
	</span>
{/snippet}

{#snippet nodeport(row: Row<Application>)}
	{#each row.original.services as service}
		{#each service.ports as port}
			{#if port.nodePort > 0}
				<span class="flex items-center">
					<Badge variant="outline">{port.targetPort}</Badge>
					<Button
						variant="ghost"
						target="_blank"
						href={`http://${row.original.publicAddress}:${port.nodePort}`}
					>
						<Icon icon="ph:arrow-square-out" />
					</Button>
				</span>
			{/if}
		{/each}
	{/each}
{/snippet}
