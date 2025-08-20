<script lang="ts">
	import type { Application } from '$lib/api/application/v1/application_pb';
	import * as Alert from '$lib/components/ui/alert/index.js';
	import Icon from '@iconify/svelte';
	import { type Writable } from 'svelte/store';

	let {
		application
	}: {
		application: Writable<Application>;
	} = $props();
</script>

{#if $application.healthies !== $application.pods.length}
	<Alert.Root variant="destructive" class="border-destructive">
		<Icon icon="ph:warning" class="size-5" />
		<Alert.Title>WARNING</Alert.Title>
		<Alert.Description>
			{$application.pods.length - $application.healthies} out of {$application.pods.length} pods are
			unhealthy.
		</Alert.Description>
	</Alert.Root>
{/if}
