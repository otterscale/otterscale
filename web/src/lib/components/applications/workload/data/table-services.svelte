<script lang="ts" module>
	import { type Writable } from 'svelte/store';

	import type { Application } from '$lib/api/application/v1/application_pb';
	import { Empty } from '$lib/components/custom/table';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import * as Table from '$lib/components/ui/table';
	import { m } from '$lib/paraglide/messages';
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
		{#each $application.services as service (service.name)}
			<Table.Row>
				<Table.Cell>{service.name}</Table.Cell>
				<Table.Cell>
					<Badge variant="outline">{service.type}</Badge>
				</Table.Cell>
				<Table.Cell>
					{service.clusterIp}
				</Table.Cell>
				<!-- <Table.Cell class="space-y-1">
					{#each service.ports as port}
						<div class="flex items-center gap-1">
							<Badge variant="outline">{port.protocol}</Badge>
							{port.port}:{port.nodePort}
						</div>
					{/each} -->
				<Table.Cell class="space-y-1">
					<div class="flex flex-col gap-1.5">
						{#each service.ports as port (port.port)}
							<div class="flex items-center gap-2">
								<Badge variant="outline" class="uppercase">
									{port.protocol}
								</Badge>
								<div class="flex items-center">
									<span>{port.port}</span>
									{#if port.nodePort > 0}
										<span class="text-muted-foreground">:{port.nodePort}</span>
									{/if}
									<span class="mx-1.5 text-muted-foreground/60">→</span>
									<span class="text-blue-500/80">{port.targetPort}</span>
								</div>
								{#if port.name}
									<span class="rounded bg-muted px-1 text-xs text-muted-foreground uppercase">
										{port.name}
									</span>
								{/if}
							</div>
						{/each}
					</div>
				</Table.Cell>
			</Table.Row>
		{/each}
		{#if $application.services.length === 0}
			<Table.Row>
				<Table.Cell colspan={4}>
					<Empty />
				</Table.Cell>
			</Table.Row>
		{/if}
	</Table.Body>
</Table.Root>
