<script lang="ts">
	import Icon from '@iconify/svelte';

	interface Props {
		icon: string;
		color?: string;
	}

	let { icon, color }: Props = $props();
</script>

<div class="bg-background relative flex size-16 rounded-xl dark:bg-transparent">
	<div
		role="presentation"
		class="absolute inset-0 rounded-xl border border-black/20 dark:border-white/25"
	></div>
	<div class="relative z-20 m-auto size-fit">
		<Icon {icon} class="size-8 {color ? `text-[${color}]` : ''}" />
	</div>
</div>
