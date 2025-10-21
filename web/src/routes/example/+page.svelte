<script lang="ts">
	import Icon from '@iconify/svelte';

	import * as Accordion from '$lib/components/ui/accordion/index.js';
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Card from '$lib/components/ui/card';
	import { Progress } from '$lib/components/ui/progress/index.js';
	import { formatProgressColor } from '$lib/formatter';
	import { cn } from '$lib/utils';

	function getNodeBackgroundByStatus(status: 'success' | 'warning' | 'fail') {
		if (status == 'success') {
			return 'bg-green-50';
		} else if (status == 'warning') {
			return 'bg-yellow-50';
		} else if (status == 'fail') {
			return 'bg-red-50';
		}
	}
	function getStateIconByStatus(status: 'success' | 'warning' | 'fail') {
		if (status == 'success') {
			return 'ph:check-circle';
		} else if (status == 'warning') {
			return 'ph:minus-circle';
		} else if (status == 'fail') {
			return 'ph:x-circle';
		} else {
			return 'ph:circle';
		}
	}
	function getNodeTextByStatus(status: 'success' | 'warning' | 'fail') {
		if (status == 'success') {
			return 'text-green-500';
		} else if (status == 'warning') {
			return 'text-yellow-500';
		} else if (status == 'fail') {
			return 'text-red-500';
		}
	}
</script>

<section class="container mx-auto space-y-12 px-4 py-24 md:px-6 2xl:max-w-[1400px]">
	<div class="mx-auto max-w-3xl space-y-4 text-center">
		<h2 class="text-3xl font-bold tracking-tight">Plugins</h2>
		<p class="text-muted-foreground">
			Get started with the installation, install dependencies, and run the dev server. See the full installation
			guide in the docs for production setup, environment variables, and troubleshooting.
		</p>
	</div>
	<Accordion.Root
		type="multiple"
		class="group bg-card text-card-foreground w-full overflow-hidden rounded-lg border transition-all duration-300"
	>
		<Accordion.Item value="item-1" class="p-6">
			<Accordion.Trigger>
				{@render Thumbnial(
					'ph:robot',
					'Models',
					'Enable vLLM plugins (tokenizers, backends, batching, logging).',
					1,
					3,
				)}
			</Accordion.Trigger>
			<Accordion.Content class="flex flex-col gap-4 text-balance">
				<div class="relative mx-auto hidden w-full space-y-12 pt-10 md:block">
					{@render Node(
						'right',
						'ph:gauge',
						'Prometheus',
						'Prometheus is an open-source monitoring system that scrapes metrics, stores time-series data, and supports queries and alerts.',
						'success',
					)}
				</div>
				<div class="relative mx-auto hidden w-full space-y-12 pt-10 md:block">
					<div class="bg-border absolute left-1/2 h-full w-0.5 -translate-x-1/2 transform"></div>
					{@render Node(
						'left',
						'ph:graphics-card',
						'HAMi',
						'Released our flagship product that changed the industry — simplifying workflows and improving efficiency.',
						'fail',
					)}
					{@render Node(
						'right',
						'ph:robot',
						'llm-d',
						'Opened London and Tokyo offices to expand our global presence and improve regional support.',
						'warning',
					)}
				</div>
			</Accordion.Content>
		</Accordion.Item>

		<Accordion.Item value="item-2" class="p-6">
			<Accordion.Trigger>
				{@render Thumbnial(
					'ph:desktop-tower',
					'Virtual Machines',
					'Provision and manage virtual machines with scalable resource allocation, snapshots, and secure networking.',
					1,
					2,
				)}
			</Accordion.Trigger>
			<Accordion.Content class="flex flex-col gap-4 text-balance">
				<div class="relative mx-auto hidden w-full space-y-12 pt-10 md:block">
					{@render Node(
						'right',
						'ph:gauge',
						'Prometheus',
						'Prometheus is an open-source monitoring system that scrapes metrics, stores time-series data, and supports queries and alerts.',
						'success',
					)}
					{@render Node(
						'left',
						'ph:graphics-card',
						'HAMi',
						'Released our flagship product that changed the industry — simplifying workflows and improving efficiency.',
						'fail',
					)}
				</div>
			</Accordion.Content>
		</Accordion.Item>

		<Accordion.Item value="item-3" class="p-6 py-2">
			<Accordion.Trigger>
				{@render Thumbnial(
					'ph:gauge',
					'Dashboard',
					'Create interactive dashboards with customizable widgets, real-time charts, and drill-down insights from metrics and logs, plus role-based access controls for secure collaboration.',
					1,
					1,
				)}
			</Accordion.Trigger>
			<Accordion.Content class="flex flex-col gap-4 text-balance">
				<div class="relative mx-auto hidden w-full space-y-12 pt-10 md:block">
					{@render Node(
						'right',
						'ph:gauge',
						'Prometheus',
						'Prometheus is an open-source monitoring system that scrapes metrics, stores time-series data, and supports queries and alerts.',
						'success',
					)}
				</div>
			</Accordion.Content>
		</Accordion.Item>
	</Accordion.Root>
</section>

{#snippet Thumbnial(icon: string, title: string, description: string, installed: number, necessaries: number)}
	{@const percentage = (installed * 100) / necessaries}
	<div class="flex w-full flex-col gap-4">
		<Progress value={percentage} class={formatProgressColor(percentage)} />
		<div class="flex gap-2">
			<div class="flex items-start gap-4">
				<div class="bg-primary/10 flex-shrink-0 rounded-md p-3">
					<Icon {icon} class="size-6" />
				</div>
				<div>
					<h3 class="text-lg font-bold">{title}</h3>
					<p class="text-muted-foreground mt-1 text-sm">
						{description}
					</p>
				</div>
			</div>
			<div class="ml-auto">
				<p class="text-muted-foreground whitespace-nowrap">{installed} over {necessaries}</p>
			</div>
		</div>
	</div>
{/snippet}

{#snippet Node(
	alignment: 'left' | 'right' = 'right',
	icon: string,
	name: string,
	description: string,
	status: 'success' | 'warning' | 'fail',
	action?: string,
)}
	<div
		class={alignment == 'right'
			? 'relative flex flex-row-reverse items-center gap-8'
			: 'relative flex items-center gap-8'}
	>
		<div
			class="bg-primary text-primary-foreground absolute left-1/2 z-10 flex h-12 w-12 -translate-x-1/2 transform items-center justify-center rounded-full font-bold"
		>
			<Icon {icon} class="size-6" />
		</div>
		<div class={alignment == 'right' ? 'w-1/2 pr-16' : 'w-1/2 pl-16'}>
			<Card.Root class={cn(getNodeBackgroundByStatus(status), 'p-0')}>
				<Card.Content class="space-y-2 p-5">
					<div class="flex flex-row-reverse items-center justify-end gap-2">
						<Icon icon={getStateIconByStatus(status)} class={cn(getNodeTextByStatus(status), 'size-6')} />
						<h3 class="font-bold">{name}</h3>
					</div>
					<p class="text-muted-foreground text-sm">
						{description}
					</p>
					{#if status == 'fail'}
						<div class="ml-auto">
							<Button class="w-full" size="sm" href={action}>Install</Button>
						</div>
					{:else if status == 'warning'}
						<div class="ml-auto">
							<Button class="w-full" size="sm" disabled href={action}>Install</Button>
						</div>
					{/if}
				</Card.Content>
			</Card.Root>
		</div>
		<div class="w-1/2"></div>
	</div>
{/snippet}
