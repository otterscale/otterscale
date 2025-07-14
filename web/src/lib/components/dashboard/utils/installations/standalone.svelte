<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import Icon from '@iconify/svelte';
	import Template from './template.svelte';
</script>

<Template>
	{#snippet icon()}
		<span class="relative flex w-full items-end">
			<Icon icon="simple-icons:ceph" class="size-24" style="color: #f0424d" />
			<Icon
				icon="simple-icons:kubernetes"
				class="absolute right-0 size-24"
				style="color: #326de6"
			/>
		</span>
	{/snippet}
	{#snippet title()}
		Standalone
	{/snippet}
	{#snippet description()}
		<p>
			Deploy a standalone Ceph storage cluster with Kubernetes. This option provides a complete
			storage solution for testing and small-scale deployments.
		</p>
	{/snippet}
	{#snippet action()}
		<Button variant="outline" class="w-full">Deploy</Button>
	{/snippet}
</Template>
