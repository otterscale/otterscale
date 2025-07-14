<script lang="ts">
	import Enterprise from '$lib/components/custom/enterprise/enterprise.svelte';
	import { Button } from '$lib/components/ui/button';
	import Icon from '@iconify/svelte';
	import Template from './template.svelte';
</script>

<Template>
	{#snippet icon()}
		<Icon icon="logos:openai-icon" class="size-24" style="color: #326de6" />
		<span class="relative flex w-full items-end justify-center">
			<Icon
				icon="simple-icons:kubernetes"
				class="absolute right-8 size-16"
				style="color: #326de6"
			/>
			<Icon icon="simple-icons:ceph" class="absolute right-0 size-16" style="color: #f0424d" />
		</span>
	{/snippet}
	{#snippet title()}
		Model
	{/snippet}
	{#snippet description()}
		<p>
			Deploy and manage your ML models with integrated GPU orchestration and scalable storage
			solutions powered by Kubernetes and Ceph
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
