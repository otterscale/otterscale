<script lang="ts">
	import { SiteFooter, SiteHeader } from '$lib/components';

	let { children } = $props();
</script>

<div class="relative flex min-h-screen flex-col bg-background">
	<SiteHeader />
	<div class="flex-1">
		{@render children()}
	</div>
	<SiteFooter />
</div>
