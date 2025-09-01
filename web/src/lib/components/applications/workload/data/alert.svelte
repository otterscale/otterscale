<script lang="ts" module>
	import type { Application } from '$lib/api/application/v1/application_pb';
	import * as Alert from '$lib/components/ui/alert/index.js';
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
</script>

{#if $application.healthies !== $application.pods.length}
	<Alert.Root variant="destructive" class="border-destructive">
		<Icon icon="ph:warning" class="size-5" />
		<Alert.Title>
			{m.workload_health()}
		</Alert.Title>
		<Alert.Description>
			{m.workload_health_description({
				unhealthies: $application.pods.length - $application.healthies,
				total: $application.pods.length,
			})}
		</Alert.Description>
	</Alert.Root>
{/if}
