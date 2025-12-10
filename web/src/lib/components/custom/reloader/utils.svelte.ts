class ReloadManager {
	private reloadFunction: () => void;
	private identifier: number | NodeJS.Timeout | undefined;

	state: boolean = $state(false);
	interval: number = $state(10);

	constructor(reloadFunction: () => void, initialState: boolean = false) {
		this.reloadFunction = reloadFunction;
		this.state = initialState;
	}

	get isReloading() {
		return this.state && this.interval && this.interval > 0;
	}

	force() {
		this.reloadFunction();
	}

	start() {
		if (this.identifier) return;
		this.state = true;
		this.identifier = setInterval(() => this.reloadFunction(), this.interval * 1000);
	}

	stop() {
		if (this.identifier) {
			clearInterval(this.identifier);
			this.identifier = undefined;
		}
		this.state = false;
	}

	restart() {
		this.stop();
		this.start();
	}
}

export { ReloadManager };
