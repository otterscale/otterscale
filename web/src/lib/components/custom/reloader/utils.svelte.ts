class ReloadManager {
    private reloadFunction: () => void;
    private identifier: number | NodeJS.Timeout | undefined;

    state: boolean = $state(true);
    interval: number | undefined = $state(5);

    constructor(reloadFunction: () => void) {
        this.reloadFunction = reloadFunction;
    }

    get isReloading() {
        return this.state && this.interval && this.interval > 0;
    }

    force() {
        this.reloadFunction()
    }

    start() {
        if (this.state && this.interval && this.interval > 0) {
            this.identifier = setInterval(() => this.reloadFunction(), this.interval * 1000);
        }
    }

    stop() {
        if (this.identifier) {
            clearInterval(this.identifier);
            this.identifier = undefined;
        }
    }

    restart() {
        this.stop();
        this.start();
    }
}

export { ReloadManager }