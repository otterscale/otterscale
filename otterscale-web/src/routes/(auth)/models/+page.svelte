<script lang="ts">
	import * as Alert from '$lib/components/ui/alert';
	import { siteConfig } from '$lib/config/site';
	import { m } from '$lib/paraglide/messages';
	import Icon from '@iconify/svelte';
</script>

<Alert.Root variant="destructive">
	<Icon icon="ph:hammer" class="size-16" />
	<Alert.Title>{m.under_development()}</Alert.Title>
	<Alert.Description>{m.under_development_description()}</Alert.Description>
</Alert.Root>

<div class="pointer-events-none fixed inset-0 flex flex-col items-center justify-center">
	<Icon icon="ph:barricade" class="text-9xl" />
	{m.current_version({ version: siteConfig.version })}
</div>
