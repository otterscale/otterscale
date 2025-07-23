<script lang="ts">
	import { activeScope, scopeLoading } from '$lib/stores';
</script>

applications dashboard

<div>
	{#if $scopeLoading}
		loading scope
	{:else}
		current scope: {$activeScope.uuid}
	{/if}
</div>
