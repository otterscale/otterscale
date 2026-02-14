<script lang="ts">
	import WandSparklesIcon from '@lucide/svelte/icons/wand-sparkles';

	import { useSidebar } from '$lib/components/ui/sidebar';
	import { m } from '$lib/paraglide/messages';
	import { cn } from '$lib/utils';

	const sidebar = useSidebar();
</script>

<div
	class={cn(
		'pointer-events-none fixed inset-0 flex flex-col items-center justify-center gap-2 transition-all duration-200 ease-linear',
		sidebar.open ? 'ml-(--sidebar-width)' : ''
	)}
>
	<WandSparklesIcon class="size-10" />

	<div class="space-y-2 text-center">
		<p class="text-lg font-semibold tracking-tight">{m.under_development()}</p>
		<p class="max-w-md text-sm text-muted-foreground">
			{m.under_development_description()}
		</p>
	</div>
</div>
