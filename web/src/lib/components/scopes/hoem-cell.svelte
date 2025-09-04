<script lang="ts">
	import Icon from '@iconify/svelte';

	import LogoImage from '$lib/assets/logo.png';

	interface Props {
		icon?: string;
		color?: string;
		us?: boolean;
	}

	let { icon, color, us }: Props = $props();
</script>

<div class="bg-background relative flex size-16 rounded-xl dark:bg-transparent">
	<div role="presentation" class="absolute inset-0 rounded-xl border border-black/20 dark:border-white/25"></div>
	<div class="relative z-20 m-auto size-fit">
		{#if icon}
			<Icon {icon} class="size-8 {color ? `text-[${color}]` : ''}" />
		{/if}
		{#if us}
			<img src={LogoImage} alt="logo" class="size-12" />
		{/if}
	</div>
</div>
