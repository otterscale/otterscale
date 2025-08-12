<script lang="ts" module>
	import type { Application } from '$lib/api/application/v1/application_pb';
	import * as Table from '$lib/components/custom/table';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import { buttonVariants } from '$lib/components/ui/button';
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
			<Table.Head>PHASE</Table.Head>
			<Table.Head>READY</Table.Head>
			<Table.Head>RESTARTS</Table.Head>
			<Table.Head>LAST CONDITION</Table.Head>
		</Table.Row>
	</Table.Header>
	<Table.Body>
		{#each $application.pods as pod}
			<Table.Row>
				<Table.Cell>{pod.name}</Table.Cell>
				<Table.Cell>
					<Badge variant="outline">{pod.phase}</Badge>
				</Table.Cell>
				<Table.Cell>
					{pod.ready}
				</Table.Cell>
				<Table.Cell>{pod.restarts}</Table.Cell>
				<Table.Cell>
					{#if pod.lastCondition}
						<Badge variant="outline">
							{pod.lastCondition.type}: {pod.lastCondition.status}
							<p class={pod.lastCondition.reason ? 'text-muted-foreground italic' : 'hidden'}>
								({pod.lastCondition.reason})
							</p>
						</Badge>
						<span class={pod.lastCondition.message ? 'flex items-center' : 'hidden'}>
							<Tooltip.Provider>
								<Tooltip.Root>
									<Tooltip.Trigger class={buttonVariants({ variant: 'ghost' })}>
										<Icon icon="ph:info" />
									</Tooltip.Trigger>
									<Tooltip.Content class="max-w-[77vw] overflow-auto">
										{pod.lastCondition.message}
									</Tooltip.Content>
								</Tooltip.Root>
							</Tooltip.Provider>
							<p class="text-muted-foreground max-w-[700px] truncate">
								{pod.lastCondition.message}
							</p>
						</span>
					{/if}
				</Table.Cell>
			</Table.Row>
		{/each}
	</Table.Body>
</Table.Root>
