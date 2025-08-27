<script lang="ts">
	import { useSidebar } from '$lib/components/ui/sidebar';
	import { m } from '$lib/paraglide/messages';
	import { cn } from '$lib/utils';
	import Icon from '@iconify/svelte';

	const sidebar = useSidebar();
</script>

<div
	class={cn(
		'pointer-events-none fixed inset-0 flex flex-col items-center justify-center gap-2 transition-all duration-200 ease-linear',
		sidebar.open ? 'ml-(--sidebar-width)' : ''
	)}
>
	<div class="border-border flex size-24 items-center justify-center rounded-full border-2">
		<Icon icon="ph:magic-wand" class="size-14" />
	</div>

	<p class="pt-2 text-3xl">{m.under_development()}</p>
	<p class="text-muted-foreground text-base">
		{m.under_development_description()}
	</p>
</div>
