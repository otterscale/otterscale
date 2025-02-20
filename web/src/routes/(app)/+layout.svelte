<script lang="ts">
	import { SiteFooter, SiteHeader } from '$lib/components';
</script>

<SiteHeader />
<div class="flex-1">
	<slot />
</div>
<SiteFooter />
