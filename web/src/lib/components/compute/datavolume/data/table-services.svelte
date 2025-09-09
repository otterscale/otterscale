<script lang="ts" module>
	import { type Writable } from 'svelte/store';

	import type { Application } from '$lib/api/application/v1/application_pb';
	import * as Table from '$lib/components/custom/table';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let {
		application,
	}: {
		application: Writable<Application>;
	} = $props();
</script>

<Table.Root>
	<Table.Header>
		<Table.Row>
			<Table.Head>
				{m.name()}
			</Table.Head>
			<Table.Head>
				{m.type()}
			</Table.Head>
			<Table.Head>
				{m.cluster_ip()}
			</Table.Head>
			<Table.Head>
				{m.ports()}
			</Table.Head>
		</Table.Row>
	</Table.Header>
	<Table.Body>
		{#each $application.services as service}
			<Table.Row>
				<Table.Cell>{service.name}</Table.Cell>
				<Table.Cell>
					<Badge variant="outline">{service.type}</Badge>
				</Table.Cell>
				<Table.Cell>
					{service.clusterIp}
				</Table.Cell>
				<Table.Cell>
					<div class="flex items-center">
						{#each service.ports as port}
							<Badge variant="outline">{port.protocol}</Badge>
							{port.port}:{port.nodePort}
							{port.targetPort}
						{/each}
					</div>
				</Table.Cell>
			</Table.Row>
		{/each}
		{#if $application.services.length === 0}
			<Table.Row>
				<Table.Cell colspan={4}>
					<Table.Empty />
				</Table.Cell>
			</Table.Row>
		{/if}
	</Table.Body>
</Table.Root>
