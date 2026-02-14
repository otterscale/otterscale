<script lang="ts">
	import Icon from '@iconify/svelte';

	import { getPathIcon, type Path } from '$lib/path';

	interface Props {
		background: string;
		path: Path;
		description: string;
		disabled?: boolean;
	}

	let { background, path, description, disabled = false }: Props = $props();
</script>

<!-- eslint-disable svelte/no-navigation-without-resolve -->
<a
	href={path.url}
	class="group flex aspect-square flex-col overflow-hidden rounded-lg border bg-card shadow-sm transition-all hover:shadow-md
	{disabled ? 'pointer-events-none cursor-not-allowed opacity-50' : ''}"
>
	<header class="relative flex aspect-video">
		<div
			class="absolute inset-0 flex h-full items-center justify-center saturate-60 transition-opacity group-hover:opacity-80 dark:saturate-100 {background}"
		>
			<div
				class="mx-auto flex size-18 overflow-hidden rounded-full bg-primary/30 shadow-md group-hover:scale-105"
			>
				<Icon icon="{getPathIcon(path.url)}-fill" class="m-auto size-8 text-muted" />
			</div>
		</div>
	</header>

	<div class="flex flex-col items-center justify-center space-y-1 p-4 text-center">
		<h3 class="text-xl font-semibold">
			{path.title}
		</h3>
		<p class="text-sm text-muted-foreground">
			{description}
		</p>
	</div>
</a>
<!-- eslint-enable svelte/no-navigation-without-resolve -->
