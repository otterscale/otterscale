<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import { type Writable } from 'svelte/store';

	import type { Application } from '$lib/api/application/v1/application_pb';
	import { Layout } from '$lib/components/custom/instance';
	import * as Table from '$lib/components/custom/table';
	import { Badge } from '$lib/components/ui/badge';
	import Button from '$lib/components/ui/button/button.svelte';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let {
		application
	}: {
		application: Writable<Application>;
	} = $props();

	let isExpand = $state(false);
</script>

<Layout.Statistic.Root class={isExpand ? 'col-span-3' : 'col-span-1'}>
	<Layout.Statistic.Header>
		<Layout.Statistic.Title>{m.persistent_volume_claims()}</Layout.Statistic.Title>
		<Layout.Statistic.Action>
			<Button
				disabled={$application.persistentVolumeClaims.length === 0}
				variant="ghost"
				onclick={() => {
					isExpand = !isExpand;
				}}
			>
				<Icon icon="ph:resize" />
			</Button>
		</Layout.Statistic.Action>
	</Layout.Statistic.Header>
	<Layout.Statistic.Content class={isExpand ? 'flex h-full flex-col justify-evenly gap-8' : ''}>
		{#if !isExpand}
			{$application.persistentVolumeClaims.length}
		{:else}
			<div class="max-h-30 w-full overflow-y-auto">
				<Table.Root>
					<Table.Header>
						<Table.Row>
							<Table.Head>
								{m.name()}
							</Table.Head>
							<Table.Head>
								{m.status()}
							</Table.Head>
							<Table.Head>
								{m.capacity()}
							</Table.Head>
							<Table.Head>
								{m.access_modes()}
							</Table.Head>
						</Table.Row>
					</Table.Header>
					<Table.Body>
						{#each $application.persistentVolumeClaims as persistentVolumeClaim}
							<Table.Row>
								<Table.Cell>{persistentVolumeClaim.name}</Table.Cell>
								<Table.Cell>
									<Badge variant="outline">
										{persistentVolumeClaim.status}
									</Badge>
								</Table.Cell>
								<Table.Cell>{persistentVolumeClaim.capacity}</Table.Cell>
								<Table.Cell>
									<span class="flex items-center gap-1">
										{#each persistentVolumeClaim.accessModes as mode}
											<Badge variant="outline">{mode}</Badge>
										{/each}
									</span>
								</Table.Cell>
							</Table.Row>
						{/each}
					</Table.Body>
				</Table.Root>
			</div>
		{/if}
	</Layout.Statistic.Content>
</Layout.Statistic.Root>
