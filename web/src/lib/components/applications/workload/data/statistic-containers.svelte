<script lang="ts" module>
	import type { Application } from '$lib/api/application/v1/application_pb';
	import { Layout } from '$lib/components/custom/instance';
	import * as Table from '$lib/components/custom/table';
	import { Badge } from '$lib/components/ui/badge';
	import Button from '$lib/components/ui/button/button.svelte';
	import { m } from '$lib/paraglide/messages';
	import Icon from '@iconify/svelte';
	import { type Writable } from 'svelte/store';
</script>

<script lang="ts">
	let {
		application,
	}: {
		application: Writable<Application>;
	} = $props();

	let isExpand = $state(false);
</script>

<Layout.Statistic.Root class={isExpand ? 'col-span-2' : 'col-span-1'}>
	<Layout.Statistic.Header>
		<Layout.Statistic.Title>
			{m.containers()}
		</Layout.Statistic.Title>
		<Layout.Statistic.Action>
			<Button
				disabled={$application.containers.length === 0}
				variant="ghost"
				onclick={() => {
					isExpand = !isExpand;
				}}
			>
				<Icon icon="ph:resize" />
			</Button>
		</Layout.Statistic.Action>
	</Layout.Statistic.Header>
	<Layout.Statistic.Content>
		{#if !isExpand}
			{$application.containers.length}
		{:else}
			<div class="max-h-30 w-full overflow-y-auto">
				<Table.Root>
					<Table.Header>
						<Table.Row>
							<Table.Head>{m.image()}</Table.Head>
							<Table.Head>{m.pull_policy()}</Table.Head>
						</Table.Row>
					</Table.Header>
					<Table.Body>
						{#each $application.containers as container}
							<Table.Row>
								<Table.Cell>{container.imageName}</Table.Cell>
								<Table.Cell>
									<Badge variant="outline">{container.imagePullPolicy}</Badge>
								</Table.Cell>
							</Table.Row>
						{/each}
					</Table.Body>
				</Table.Root>
			</div>
		{/if}
	</Layout.Statistic.Content>
</Layout.Statistic.Root>
