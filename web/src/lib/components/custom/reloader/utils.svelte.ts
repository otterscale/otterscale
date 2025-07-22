class ReloadManager {
    private reloadFn: () => void;
    private identifier: number | NodeJS.Timeout | undefined;

    state: boolean = $state(false);
    interval: number | undefined = $state(15);

    constructor(reloadFn: () => void) {
        this.reloadFn = reloadFn;
    }

    get isReloading() {
        return this.state && this.interval && this.interval > 0;
    }

    start() {
        if (this.state && this.interval && this.interval > 0) {
            this.identifier = setInterval(() => this.reloadFn(), this.interval * 1000);
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