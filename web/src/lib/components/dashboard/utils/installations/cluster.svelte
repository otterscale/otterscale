<script lang="ts">
	import Enterprise from '$lib/components/custom/enterprise/enterprise.svelte';
	import { Button } from '$lib/components/ui/button';
	import Icon from '@iconify/svelte';
	import Template from './template.svelte';
</script>

<Template>
	{#snippet icon()}
		<span class="relative flex w-full items-end">
			<Icon icon="simple-icons:ceph" class="absolute right-0  size-24" style="color: #f0424d" />
			<Icon icon="simple-icons:kubernetes" class="size-24" style="color: #326de6" />
		</span>
	{/snippet}
	{#snippet title()}
		Cluster
	{/snippet}
	{#snippet description()}
		<p>
			Set up and manage a robust storage infrastructure with integrated Ceph and Kubernetes. Perfect
			for enterprise-grade distributed storage solutions with high availability.
		</p>
	{/snippet}
	{#snippet action()}
		<Enterprise class="w-full">
			{#snippet children({ disable })}
				<Button disabled={disable} class="w-full">Deploy</Button>
			{/snippet}
		</Enterprise>
	{/snippet}
</Template>
