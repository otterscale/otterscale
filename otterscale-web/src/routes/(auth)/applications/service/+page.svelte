<script lang="ts">
	import Icon from '@iconify/svelte';
	import * as Alert from '$lib/components/ui/alert';
	import { siteConfig } from '$lib/config/site';
	import { m } from '$lib/paraglide/messages';
	import { applicationsPath, applicationsServicePath } from '$lib/path';
	import { activeScope, breadcrumb } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({ parents: [applicationsPath], current: applicationsServicePath });
</script>

{#if $activeScope}
	current scope: {$activeScope.uuid}
{/if}

<Alert.Root variant="default">
	<Icon icon="ph:airplane-takeoff" />
	<Alert.Title>{m.migrating()}</Alert.Title>
	<Alert.Description>{m.migrating_description()}</Alert.Description>
</Alert.Root>

<div class="pointer-events-none fixed inset-0 flex flex-col items-center justify-center">
	<Icon icon="ph:barricade" class="text-9xl" />
	{m.current_version({ version: siteConfig.version })}
</div>
