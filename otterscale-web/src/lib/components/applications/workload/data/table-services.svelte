<script lang="ts" module>
	import type { Application } from '$lib/api/application/v1/application_pb';
	import * as Table from '$lib/components/custom/table';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import { buttonVariants } from '$lib/components/ui/button';
	import Separator from '$lib/components/ui/separator/separator.svelte';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import Icon from '@iconify/svelte';
	import { type Writable } from 'svelte/store';
</script>

<script lang="ts">
	let {
		application
	}: {
		application: Writable<Application>;
	} = $props();
</script>

<Table.Root>
	<Table.Header>
		<Table.Row>
			<Table.Head>NAME</Table.Head>
			<Table.Head>TYPE</Table.Head>
			<Table.Head>CLUSTER IP</Table.Head>
			<Table.Head>PORTS</Table.Head>
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
	</Table.Body>
</Table.Root>
