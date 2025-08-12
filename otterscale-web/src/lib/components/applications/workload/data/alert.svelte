<script lang="ts" module>
	import type { Application } from '$lib/api/application/v1/application_pb';
	import * as Alert from '$lib/components/ui/alert/index.js';
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

{#if $application.healthies !== $application.pods.length}
	<Alert.Root variant="destructive">
		<span class="flex items-center gap-3">
			<Icon icon="radix-icons:exclamation-triangle" class="size-5" />

			<span>
				<Alert.Title>WARNING</Alert.Title>
				<Alert.Description>
					{$application.pods.length - $application.healthies} out of {$application.pods.length} pods
					are unhealthy.
				</Alert.Description>
			</span>
		</span>
	</Alert.Root>
{/if}
