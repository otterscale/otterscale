<script lang="ts">
	import Icon from '@iconify/svelte';

	import { siteConfig } from '$lib/config/site.js';

	import Features from './features.svelte';
</script>

<nav
	class="hidden flex-col gap-6 whitespace-nowrap text-lg font-medium md:flex md:flex-row md:items-center md:gap-5 md:text-sm lg:gap-6"
>
	<a href="/" class="flex items-center gap-2 text-lg font-semibold md:text-base [&_svg]:size-8">
		<Icon icon="ph:polygon" />
		<span class="sr-only">{siteConfig.name}</span>
	</a>
	<Features />
</nav>
