<script lang="ts">
	import LayersIcon from '@lucide/svelte/icons/layers';
	import { onMount } from 'svelte';

	import { version } from '$app/environment';
	import { startTour } from '$lib/components/layout';
	import { useSidebar } from '$lib/components/ui/sidebar';
	import { m } from '$lib/paraglide/messages';
	import { cn } from '$lib/utils';

	const STORAGE_KEY = 'tutorial_version';
	const sidebar = useSidebar();

	onMount(() => {
		const seenVersion = localStorage.getItem(STORAGE_KEY);

		if (seenVersion !== version) {
			startTour();
			localStorage.setItem(STORAGE_KEY, version);
		}
	});
</script>

<div
	class={cn(
		'pointer-events-none fixed inset-0 flex flex-col items-center justify-center gap-2 transition-all duration-200 ease-linear',
		sidebar.open ? 'ml-(--sidebar-width)' : ''
	)}
>
	<div class="absolute top-12 right-18 flex flex-col items-end opacity-90 md:right-20">
		<svg
			width="58"
			height="43"
			viewBox="0 0 58 43"
			fill="none"
			xmlns="http://www.w3.org/2000/svg"
			class="size-14 rotate-10 transform text-muted-foreground"
		>
			<path
				d="M53.7685 0.146447C53.5733 -0.0488155 53.2567 -0.0488155 53.0614 0.146447L49.8794 3.32843C49.6842 3.52369 49.6842 3.84027 49.8794 4.03553C50.0747 4.2308 50.3913 4.2308 50.5866 4.03553L53.415 1.20711L56.2434 4.03553C56.4387 4.2308 56.7552 4.2308 56.9505 4.03553C57.1458 3.84027 57.1458 3.52369 56.9505 3.32843L53.7685 0.146447ZM0.414978 42.5L0.829954 42.7789C10.9459 27.7284 23.972 26.623 34.6945 24.6165C40.0335 23.6173 44.88 22.3807 48.3672 18.9502C51.8715 15.5028 53.915 9.93657 53.915 0.5H53.415H52.915C52.915 9.81343 50.896 15.0597 47.6659 18.2373C44.4187 21.4318 39.8589 22.6327 34.5105 23.6335C23.8579 25.627 10.3841 26.7716 2.44379e-06 42.2211L0.414978 42.5Z"
				stroke="currentColor"
				fill="black"
				class="animate-pulse"
			/>
		</svg>
		<p class="font-mono text-xs text-muted-foreground/80">
			{m.console_guide()}
		</p>
	</div>

	<LayersIcon class="size-10" />

	<div class="space-y-2 text-center">
		<p class="text-lg font-semibold tracking-tight">{m.console_title()}</p>
		<p class="max-w-md text-sm text-muted-foreground">
			{m.console_description()}
		</p>
	</div>
</div>
