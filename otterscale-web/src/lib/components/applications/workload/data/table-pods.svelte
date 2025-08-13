<script lang="ts" module>
	import type { Application } from '$lib/api/application/v1/application_pb';
	import * as Table from '$lib/components/custom/table';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import { buttonVariants } from '$lib/components/ui/button';
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { cn } from '$lib/utils';
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
						{#if pod.lastCondition.reason || pod.lastCondition.message}
							<div class="text-destructive flex items-center gap-2">
								<Badge
									variant="destructive"
									class={pod.lastCondition.reason ? 'visible' : 'hidden'}
								>
									{pod.lastCondition.reason}
								</Badge>
								<Tooltip.Provider>
									<Tooltip.Root>
										<Tooltip.Trigger>
											<p
												class={cn(pod.lastCondition.message ? 'max-w-[1000px] truncate' : 'hidden')}
											>
												{pod.lastCondition.message}
											</p>
										</Tooltip.Trigger>
										<Tooltip.Content class="max-w-[77vw] overflow-auto">
											{pod.lastCondition.message}
										</Tooltip.Content>
									</Tooltip.Root>
								</Tooltip.Provider>
							</div>
						{:else}
							<Badge variant="outline">
								{pod.lastCondition.type}
							</Badge>
						{/if}
					{/if}
				</Table.Cell>
			</Table.Row>
		{/each}
	</Table.Body>
</Table.Root>
