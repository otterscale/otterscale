<script lang="ts">
	import '@xterm/xterm/css/xterm.css';

	import { createClient, type Transport } from '@connectrpc/connect';
	import type { ITerminalInitOnlyOptions, ITerminalOptions, Terminal } from '@xterm/xterm';
	import { getContext, onMount } from 'svelte';

	import { ApplicationService } from '$lib/api/application/v1/application_pb';

	let {
		scope,
		facility,
		namespace,
		podName,
		containerName,
		command
	}: {
		scope: string;
		facility: string;
		namespace: string;
		podName: string;
		containerName: string;
		command: string[];
	} = $props();

	// Types
	interface TerminalState {
		terminal?: Terminal;
		sessionId: string;
		isConnected: boolean;
	}

	interface TerminalAddons {
		clipboard: any;
		fit: any;
		search: any;
		unicode11: any;
		webLinks: any;
		webgl: any;
	}

	// Configuration
	const CONTROL_SEQUENCES = {
		CTRL_C: '\x03',
		CTRL_D: '\x04',
		CTRL_Z: '\x1a',
		CTRL_L: '\x0c',
		CTRL_A: '\x01',
		CTRL_E: '\x05',
		CTRL_K: '\x0b',
		CTRL_U: '\x15',
		CTRL_W: '\x17',
		CTRL_R: '\x12',
		ALT_B: '\x1bb',
		ALT_F: '\x1bf',
		ALT_D: '\x1bd',
		ALT_BACKSPACE: '\x1b\x7f'
	} as const;

	const TERMINAL_OPTIONS: ITerminalOptions & ITerminalInitOnlyOptions = {
		allowProposedApi: true,
		fontFamily: 'Consolas, Monaco, "Lucida Console", monospace',
		fontSize: 14,
		cursorBlink: true,
		theme: {
			background: '#1e1e1e',
			foreground: '#d4d4d4'
		}
	};

	const KEYBOARD_MAPPINGS = {
		ctrl: {
			c: CONTROL_SEQUENCES.CTRL_C,
			d: CONTROL_SEQUENCES.CTRL_D,
			z: CONTROL_SEQUENCES.CTRL_Z,
			l: CONTROL_SEQUENCES.CTRL_L,
			a: CONTROL_SEQUENCES.CTRL_A,
			e: CONTROL_SEQUENCES.CTRL_E,
			k: CONTROL_SEQUENCES.CTRL_K,
			u: CONTROL_SEQUENCES.CTRL_U,
			w: CONTROL_SEQUENCES.CTRL_W,
			r: CONTROL_SEQUENCES.CTRL_R
		},
		alt: {
			b: CONTROL_SEQUENCES.ALT_B,
			f: CONTROL_SEQUENCES.ALT_F,
			d: CONTROL_SEQUENCES.ALT_D,
			Backspace: CONTROL_SEQUENCES.ALT_BACKSPACE
		}
	} as const;

	// State
	const transport: Transport = getContext('transport');
	const client = createClient(ApplicationService, transport);

	let container = $state<HTMLElement>();
	let handleResize = $state<() => void>();
	let terminalState = $state<TerminalState>({
		sessionId: '',
		isConnected: false
	});

	// Terminal utilities
	function writeToTerminal(message: string, isError = false): void {
		if (!terminalState.terminal) return;

		const color = isError ? '\x1b[31m' : '\x1b[33m';
		const formattedMessage = `\r\n${color}${message}\x1b[0m\r\n`;
		terminalState.terminal.write(formattedMessage);
	}

	// Terminal setup
	async function initializeTerminal(): Promise<void> {
		const { Terminal } = await import('@xterm/xterm');

		terminalState.terminal = new Terminal(TERMINAL_OPTIONS);
		terminalState.terminal.onData(handleTerminalData);
		terminalState.terminal.onKey(handleTerminalKey);

		if (container) {
			terminalState.terminal.open(container);
			await loadTerminalAddons();
		}
	}

	async function loadTerminalAddons(): Promise<void> {
		if (!terminalState.terminal) return;

		try {
			const addons = await loadAddons();
			setupAddons(addons);
		} catch (error) {
			writeToTerminal(`Failed to load terminal addons: ${error}`, true);
		}
	}

	async function loadAddons(): Promise<TerminalAddons> {
		const [clipboard, fit, search, unicode11, webLinks, webgl] = await Promise.all([
			import('@xterm/addon-clipboard').then((m) => new m.ClipboardAddon()),
			import('@xterm/addon-fit').then((m) => new m.FitAddon()),
			import('@xterm/addon-search').then((m) => new m.SearchAddon()),
			import('@xterm/addon-unicode11').then((m) => new m.Unicode11Addon()),
			import('@xterm/addon-web-links').then((m) => new m.WebLinksAddon()),
			import('@xterm/addon-webgl').then((m) => new m.WebglAddon())
		]);

		return { clipboard, fit, search, unicode11, webLinks, webgl };
	}

	function setupAddons(addons: TerminalAddons): void {
		if (!terminalState.terminal) return;

		const { terminal } = terminalState;

		// Load addons
		terminal.loadAddon(addons.clipboard);
		terminal.loadAddon(addons.fit);
		terminal.loadAddon(addons.search);
		terminal.loadAddon(addons.unicode11);
		terminal.loadAddon(addons.webLinks);

		// Configure fit addon
		handleResize = () => addons.fit.fit();
		handleResize();

		// Configure unicode
		terminal.unicode.activeVersion = '11';

		// Configure WebGL with error handling
		addons.webgl.onContextLoss(() => addons.webgl.dispose());
		terminal.loadAddon(addons.webgl);
	}

	// TTY communication
	async function startTTYSession(): Promise<void> {
		try {
			const stream = client.executeTTY({
				scope: scope,
				facility: facility,
				namespace: namespace,
				podName: podName,
				containerName: containerName,
				command: command
			});
			terminalState.isConnected = true;

			for await (const response of stream) {
				handleTTYResponse(response);
			}
		} catch (error) {
			terminalState.isConnected = false;
			writeToTerminal(`Connection failed: ${error}`, true);
		}
	}

	function handleTTYResponse(response: any): void {
		if (!terminalState.sessionId && response.sessionId) {
			terminalState.sessionId = response.sessionId;
		}

		if (response.stdout) {
			const data = new TextDecoder().decode(response.stdout);
			terminalState.terminal?.write(data);

			if (data.includes('exit') && data.includes('\r')) {
				writeToTerminal('Session terminated. Connection closed.');
				terminalState.isConnected = false;
			}
		}
	}

	function sendToTTY(data: string): void {
		if (!terminalState.sessionId || !terminalState.isConnected) return;

		try {
			client.writeTTY({
				sessionId: terminalState.sessionId,
				stdin: new TextEncoder().encode(data)
			});
		} catch (error) {
			writeToTerminal(`Failed to send data: ${error}`, true);
		}
	}

	// Event handlers
	function handleTerminalData(data: string): void {
		sendToTTY(data);
	}

	function handleTerminalKey(event: { key: string; domEvent: KeyboardEvent }): void {
		const { domEvent } = event;

		if (domEvent.ctrlKey) {
			handleModifierKey(KEYBOARD_MAPPINGS.ctrl, domEvent);
		} else if (domEvent.altKey) {
			handleModifierKey(KEYBOARD_MAPPINGS.alt, domEvent);
		}
	}

	function handleModifierKey(mapping: Record<string, string>, event: KeyboardEvent): void {
		const sequence = mapping[event.key as keyof typeof mapping];

		if (sequence) {
			sendToTTY(sequence);
			event.preventDefault();
		}
	}

	// Lifecycle
	onMount(async () => {
		await initializeTerminal();
		await startTTYSession();
	});
</script>

<svelte:window onresize={handleResize} />

<div class="h-screen w-screen" bind:this={container}></div>

<style>
	:global(.xterm) {
		height: 100% !important;
	}

	:global(.xterm-viewport) {
		overflow-y: hidden !important;
	}
</style>
