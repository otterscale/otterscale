<script lang="ts" module>
	import type { Application } from '$lib/api/application/v1/application_pb';
	import { Layout } from '$lib/components/custom/instance';
	import { Badge } from '$lib/components/ui/badge';
	import Icon from '@iconify/svelte';
	import { type Writable } from 'svelte/store';
	import Alert from './alert.svelte';
	import StatisticContainers from './statistic-containers.svelte';
	import StatisticPersistVolumeClaims from './statistic-persist-volume-claims.svelte';
	import StatisticStorageClasses from './statistic-storage-classes.svelte';
	import TablePods from './table-pods.svelte';
	import TableServices from './table-services.svelte';
</script>

<script lang="ts">
	let {
		application
	}: {
		application: Writable<Application>;
	} = $props();
</script>

<main class="space-y-4 py-4">
	<Alert {application} />

	<div class="space-y-4 py-4">
		<div class="flex items-end gap-2 text-5xl">
			<p class="text-muted-foreground">{$application.namespace}</p>
			{$application.name}
		</div>
		<Badge variant="outline">
			{$application.type}
		</Badge>
		<div class="flex flex-wrap gap-1 overflow-visible">
			{#each Object.entries($application.labels) as [key, value]}
				<Badge variant="secondary">
					{key}: {value}
				</Badge>
			{/each}
		</div>
	</div>

	<Layout.Statistics>
		<StatisticContainers {application} />
		<StatisticPersistVolumeClaims {application} />
		<StatisticStorageClasses {application} />
	</Layout.Statistics>

	<Layout.Tables>
		<Layout.Table.Root open={true}>
			<Layout.Table.Trigger>
				<Icon icon="ph:cube" />
				Pods
			</Layout.Table.Trigger>
			<Layout.Table.Content>
				<TablePods {application} />
			</Layout.Table.Content>
		</Layout.Table.Root>

		<Layout.Table.Root open={true}>
			<Layout.Table.Trigger>
				<Icon icon="ph:cube" />
				Services
			</Layout.Table.Trigger>
			<Layout.Table.Content>
				<TableServices {application} />
			</Layout.Table.Content>
		</Layout.Table.Root>
	</Layout.Tables>
</main>
