<script lang="ts">
	import Icon from '@iconify/svelte';
	import { toggleMode } from 'mode-watcher';

	import type { LightSwitchProps } from './types';

	import { Button } from '$lib/components/ui/button';

	let { variant = 'ghost' }: LightSwitchProps = $props();
</script>

<Button onclick={toggleMode} {variant} size="icon">
	<Icon icon="ph:sun-bold" class="size-4.5 scale-100 rotate-0 !transition-all dark:scale-0 dark:-rotate-90" />
	<Icon icon="ph:moon" class="absolute size-4.5 scale-0 rotate-90 !transition-all dark:scale-100 dark:rotate-0" />
	<span class="sr-only">Toggle theme</span>
</Button>
