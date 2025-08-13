<script lang="ts" module>
	import type { Application } from '$lib/api/application/v1/application_pb';
	import * as Table from '$lib/components/custom/table';
	import { Badge } from '$lib/components/ui/badge';
	import Button from '$lib/components/ui/button/button.svelte';
	import Icon from '@iconify/svelte';
	import { type Writable } from 'svelte/store';
	import { Statistic } from '../layout';
</script>

<script lang="ts">
	let {
		application
	}: {
		application: Writable<Application>;
	} = $props();

	let isExpand = $state(false);
</script>

<Statistic.Root class={isExpand ? 'col-span-3' : 'col-span-1'}>
	<Statistic.Header>
		<Statistic.Title>Storage Classes</Statistic.Title>
		<Statistic.Action>
			<Button
				disabled={$application.persistentVolumeClaims.length === 0}
				variant="ghost"
				onclick={() => {
					isExpand = !isExpand;
				}}
			>
				<Icon icon="ph:resize" />
			</Button>
		</Statistic.Action>
	</Statistic.Header>
	<Statistic.Content class={isExpand ? 'flex h-full flex-col justify-evenly gap-8' : ''}>
		{#if !isExpand}
			{$application.persistentVolumeClaims.length}
		{:else}
			<div class="max-h-30 w-full overflow-y-auto">
				<Table.Root>
					<Table.Header>
						<Table.Row>
							<Table.Head>Name</Table.Head>
							<Table.Head>Provisioner</Table.Head>
							<Table.Head>Reclaim Policy</Table.Head>
							<Table.Head>Volume Binding Mode</Table.Head>
						</Table.Row>
					</Table.Header>
					<Table.Body>
						{#each $application.persistentVolumeClaims as persistentVolumeClaim}
							{#if persistentVolumeClaim.storageClass}
								<Table.Row>
									<Table.Cell>{persistentVolumeClaim.storageClass.name}</Table.Cell>
									<Table.Cell>{persistentVolumeClaim.storageClass.provisioner}</Table.Cell>
									<Table.Cell>
										{persistentVolumeClaim.storageClass.reclaimPolicy}
									</Table.Cell>
									<Table.Cell>
										<Badge variant="outline">
											{persistentVolumeClaim.storageClass.volumeBindingMode}
										</Badge>
									</Table.Cell>
								</Table.Row>

								<Table.Row>
									<Table.Cell>{persistentVolumeClaim.storageClass.name}</Table.Cell>
									<Table.Cell>{persistentVolumeClaim.storageClass.provisioner}</Table.Cell>
									<Table.Cell>
										{persistentVolumeClaim.storageClass.reclaimPolicy}
									</Table.Cell>
									<Table.Cell>
										<Badge variant="outline">
											{persistentVolumeClaim.storageClass.volumeBindingMode}
										</Badge>
									</Table.Cell>
								</Table.Row>
								<Table.Row>
									<Table.Cell>{persistentVolumeClaim.storageClass.name}</Table.Cell>
									<Table.Cell>{persistentVolumeClaim.storageClass.provisioner}</Table.Cell>
									<Table.Cell>
										{persistentVolumeClaim.storageClass.reclaimPolicy}
									</Table.Cell>
									<Table.Cell>
										<Badge variant="outline">
											{persistentVolumeClaim.storageClass.volumeBindingMode}
										</Badge>
									</Table.Cell>
								</Table.Row>
								<Table.Row>
									<Table.Cell>{persistentVolumeClaim.storageClass.name}</Table.Cell>
									<Table.Cell>{persistentVolumeClaim.storageClass.provisioner}</Table.Cell>
									<Table.Cell>
										{persistentVolumeClaim.storageClass.reclaimPolicy}
									</Table.Cell>
									<Table.Cell>
										<Badge variant="outline">
											{persistentVolumeClaim.storageClass.volumeBindingMode}
										</Badge>
									</Table.Cell>
								</Table.Row>
							{/if}
						{/each}
					</Table.Body>
				</Table.Root>
			</div>
		{/if}
	</Statistic.Content>
</Statistic.Root>
