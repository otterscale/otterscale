<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import { ApplicationService } from '$lib/api/application/v1/application_pb';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import { Label } from '$lib/components/ui/label';
	import { formatCapacity } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';
	import { cn } from '$lib/utils';

	import type { UploadedFile } from './types';
</script>

<script lang="ts">
	const handleUpload = (selectedFile: File) => {
		const maxFileSize = 20 * 1024 * 1024;

		if (!selectedFile) return;

		if (selectedFile.size > maxFileSize) {
			const { value: maxFileSizeValue, unit: maxFileSizeUnit } = formatCapacity(maxFileSize);
			toast.error(
				`${selectedFile.name} is too large. Maximum size is ${maxFileSizeValue} ${maxFileSizeUnit}`
			);
			return;
		}

		if (!selectedFile.name.endsWith('.tgz') && !selectedFile.name.endsWith('.tar.gz')) {
			toast.error(
				`${selectedFile.name} is not a valid Helm chart file. Only .tgz and .tar.gz files are supported`
			);
			return;
		}

		uploadedFile = {
			name: selectedFile.name,
			size: selectedFile.size,
			type: selectedFile.type,
			lastModifiedAt: selectedFile.lastModified,
			url: Promise.resolve(URL.createObjectURL(selectedFile)),
			uploadedAt: Date.now()
		};
	};

	async function getChartContent(uploadedFile: UploadedFile) {
		const url = await uploadedFile.url;
		const response = await fetch(url);
		const arrayBuffer = await response.arrayBuffer();
		const chartContent = new Uint8Array(arrayBuffer);
		return chartContent;
	}

	const upload = async () => {
		try {
			if (!uploadedFile) {
				toast.error('Please select a chart file to upload');
				return;
			}

			if (!uploadedFile.name.endsWith('.tgz') && !uploadedFile.name.endsWith('.tar.gz')) {
				toast.error('Please select a valid Helm chart file (.tgz or .tar.gz)');
				return;
			}

			const chartContent = await getChartContent(uploadedFile);
			await client.uploadChart({
				chartContent: chartContent
			});

			toast.success(`Chart uploaded successfully!`, {
				description: `Saved to local charts directory`
			});
		} catch (error) {
			if (error instanceof ConnectError) {
				toast.error(`Upload failed: ${error.message}`);
			} else {
				toast.error('Unknown error occurred during upload');
			}
			console.error('Error uploading charts:', error);
		}
	};

	const handleDragOver = (e: DragEvent) => {
		e.preventDefault();
		e.stopPropagation();
		isDragging = true;
	};

	const handleDragEnter = (e: DragEvent) => {
		e.preventDefault();
		e.stopPropagation();
		isDragging = true;
	};

	const handleDragLeave = (e: DragEvent) => {
		e.preventDefault();
		e.stopPropagation();

		const currentTarget = e.currentTarget as HTMLElement;
		const relatedTarget = e.relatedTarget as HTMLElement;
		if (!currentTarget?.contains(relatedTarget)) {
			isDragging = false;
		}
	};

	const handleDrop = (e: DragEvent) => {
		e.preventDefault();
		e.stopPropagation();
		isDragging = false;

		const droppedFiles = e.dataTransfer?.files;
		if (droppedFiles) {
			handleUpload(droppedFiles[0]);
		}
	};

	const handleKeydown = (e: KeyboardEvent) => {
		if (e.key === 'Enter' || e.key === ' ') {
			e.preventDefault();
			document.getElementById('file-upload')?.click();
		}
	};

	let { class: className }: { class?: string } = $props();

	const transport: Transport = getContext('transport');
	const client = createClient(ApplicationService, transport);

	let open = $state(false);
	function close() {
		open = false;
	}

	let uploadedFile = $state<UploadedFile | undefined>(undefined);
	function reset() {
		uploadedFile = undefined;
	}

	let isDragging = $state(false);
</script>

<Modal.Root bind:open>
	<Modal.Trigger variant="primary" class={cn(className)}>
		<Icon icon="ph:upload" />
		{m.upload()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>
			{m.upload()}
			<div class="text-sm font-normal text-muted-foreground">
				{m.applications_store_chart_upload_description()}
			</div>
		</Modal.Header>
		<div
			class={cn(
				'flex w-full flex-col gap-2 rounded-lg border-2 border-dashed p-6 text-center transition-colors hover:bg-muted',
				isDragging ? 'border-blue-400 bg-blue-50' : 'border-gray-300 hover:border-gray-400'
			)}
			role="button"
			tabindex="0"
			ondragover={handleDragOver}
			ondragenter={handleDragEnter}
			ondragleave={handleDragLeave}
			ondrop={handleDrop}
			onkeydown={handleKeydown}
			aria-label="Upload Helm chart file"
		>
			<input
				type="file"
				accept=".tgz,.tar.gz,application/gzip,application/x-gzip"
				class="hidden"
				id="file-upload"
				onchange={(e) => {
					const target = e.target as HTMLInputElement;
					if (target.files) {
						handleUpload(target.files[0]);
					}
				}}
			/>
			<Label
				for="file-upload"
				class="flex min-h-36 cursor-pointer flex-col items-center justify-center space-y-2"
			>
				{#if !uploadedFile}
					<Icon icon="ph:upload" class="size-8 text-gray-400" />
					<div>
						<p class="text-sm text-gray-600">
							{isDragging
								? m.applications_store_chart_upload_is_dragging()
								: m.applications_store_chart_upload_is_not_dragging()}
						</p>
						<p class="text-xs text-gray-400">{m.applications_store_chart_upload_constraint()}</p>
					</div>
				{:else}
					{@const { value: fileSizeValue, unit: fileSizeUnit } = formatCapacity(uploadedFile.size)}
					<Icon icon="ph:file-archive" class="size-8 text-gray-600" />
					<div class="space-y-1">
						<p class="text-sm text-gray-600">
							{uploadedFile.name}
						</p>
						<p class="text-xs text-gray-400">
							{new Date(uploadedFile.lastModifiedAt).toLocaleString()}．{fileSizeValue}
							{fileSizeUnit}
						</p>
					</div>
				{/if}
			</Label>
		</div>

		<Modal.Footer>
			<Modal.Cancel
				onclick={() => {
					reset();
				}}>{m.cancel()}</Modal.Cancel
			>
			<Modal.Action
				disabled={!uploadedFile}
				onclick={async () => {
					await upload();
					reset();
					close();
				}}
			>
				{m.confirm()}
			</Modal.Action>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
