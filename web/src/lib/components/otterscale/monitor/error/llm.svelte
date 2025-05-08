<script lang="ts">
	import Icon from '@iconify/svelte';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
</script>

<Card.Root>
	<Card.Content class="flex aspect-square items-center justify-center">
		<div class="flex-col items-center justify-center space-y-2">
			<div class="flex items-center justify-center space-x-2">
				<div class="flex items-end">
					<Icon icon="logos:openai-icon" class="size-24 p-2" style="color: #326de6" />
				</div>
				<div class="flex-col space-y-1">
					<p class="text-xl font-semibold leading-none tracking-tight">[WIP] LLM Platform</p>
					<p class="text-sm font-medium leading-none tracking-tight text-muted-foreground">
						with Kubernetes & Ceph Storage
					</p>
					<div class="flex space-x-2">
						<Icon icon="simple-icons:ceph" class="size-6" style="color: #f0424d" />
						<Icon icon="simple-icons:kubernetes" class="size-6" style="color: #326de6" />
					</div>
				</div>
			</div>
			<div class="flex-col text-center">
				<p class="text-sm text-muted-foreground">Deploy Your LLM and Container Platform</p>
				<p class="text-sm text-muted-foreground">with Storage</p>
			</div>
			<div class="flex items-center justify-center pt-6">
				<Button disabled={true}>Deploy</Button>
			</div>
		</div>
	</Card.Content>
</Card.Root>
