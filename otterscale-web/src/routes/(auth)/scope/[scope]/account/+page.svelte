<script lang="ts">
	import Icon from '@iconify/svelte';
	import { page } from '$app/state';
	import * as Alert from '$lib/components/ui/alert';
	import { m } from '$lib/paraglide/messages';
	import { dynamicPaths } from '$lib/path';
	import { breadcrumb } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({ parents: [], current: dynamicPaths.account(page.params.scope) });
</script>

<Alert.Root variant="default">
	<Icon icon="ph:airplane-takeoff" />
	<Alert.Title>{m.migrating()}</Alert.Title>
	<Alert.Description>{m.migrating_description()}</Alert.Description>
</Alert.Root>

<div class="pointer-events-none fixed inset-0 flex flex-col items-center justify-center">
	<Icon icon="ph:barricade" class="text-9xl" />
	{m.current_version({ version: import.meta.env.PACKAGE_VERSION })}
</div>
