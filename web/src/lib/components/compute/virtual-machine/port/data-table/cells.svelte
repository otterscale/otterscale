<script lang="ts" module>
	import { timestampDate } from '@bufbuild/protobuf/wkt';
	import Icon from '@iconify/svelte';
	import type { Row } from '@tanstack/table-core';

	import type { Application_Service } from '$lib/api/application/v1/application_pb';
	import { Cells } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { Badge } from '$lib/components/ui/badge';
	import * as HoverCard from '$lib/components/ui/hover-card';
	import * as Table from '$lib/components/ui/table';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { formatTimeAgo } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';

	export const cells = {
		row_picker,
		name,
		type,
		clusterIp,
		port,
		createTime,
	};
</script>

{#snippet row_picker(row: Row<Application_Service>)}
	<Layout.Cell class="items-center">
		<Cells.RowPicker {row} />
	</Layout.Cell>
{/snippet}

{#snippet name(row: Row<Application_Service>)}
	<Layout.Cell class="items-start">
		{row.original.name}
	</Layout.Cell>
{/snippet}

{#snippet type(row: Row<Application_Service>)}
	<Layout.Cell class="items-start">
		{#if row.original.type}
			<Badge variant="outline">
				{row.original.type}
			</Badge>
		{/if}
	</Layout.Cell>
{/snippet}

{#snippet clusterIp(row: Row<Application_Service>)}
	<Layout.Cell class="items-start">
		{row.original.clusterIp}
	</Layout.Cell>
{/snippet}

{#snippet port(row: Row<Application_Service>)}
	<Layout.Cell class="items-end">
		{#each row.original.ports as port}
			{#if port.nodePort > 0}
				<div class="flex items-center gap-1">
					<Badge variant="outline">
						{port.port}:{port.nodePort} → {port.targetPort}
						{#if port.protocol}
							({port.protocol})
						{/if}
					</Badge>
					<HoverCard.Root>
						<HoverCard.Trigger>
							<Icon icon="ph:info" />
						</HoverCard.Trigger>
						<HoverCard.Content class="min-w-[300px]">
							<Table.Root>
								<Table.Body class="text-xs">
									{#if port.name}
										<Table.Row>
											<Table.Head class="text-right">{m.name()}</Table.Head>
											<Table.Cell>{port.name}</Table.Cell>
										</Table.Row>
									{/if}
									{#if port.protocol}
										<Table.Row>
											<Table.Head class="text-right">{m.protocol()}</Table.Head>
											<Table.Cell>{port.protocol}</Table.Cell>
										</Table.Row>
									{/if}
									<Table.Row>
										<Table.Head class="text-right">{m.ports()}</Table.Head>
										<Table.Cell>{port.port}</Table.Cell>
									</Table.Row>
									<Table.Row>
										<Table.Head class="text-right">{m.node_port()}</Table.Head>
										<Table.Cell>{port.nodePort}</Table.Cell>
									</Table.Row>
									<Table.Row>
										<Table.Head class="text-right">{m.target_port()}</Table.Head>
										<Table.Cell>{port.targetPort}</Table.Cell>
									</Table.Row>
								</Table.Body>
							</Table.Root>
						</HoverCard.Content>
					</HoverCard.Root>
				</div>
			{/if}
		{/each}
	</Layout.Cell>
{/snippet}

{#snippet createTime(row: Row<Application_Service>)}
	<Layout.Cell class="items-end">
		{#if row.original.createdAt}
			<Tooltip.Provider>
				<Tooltip.Root>
					<Tooltip.Trigger>
						{formatTimeAgo(timestampDate(row.original.createdAt))}
					</Tooltip.Trigger>
					<Tooltip.Content>
						{timestampDate(row.original.createdAt)}
					</Tooltip.Content>
				</Tooltip.Root>
			</Tooltip.Provider>
		{/if}
	</Layout.Cell>
{/snippet}
