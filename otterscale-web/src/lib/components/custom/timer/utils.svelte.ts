enum Speed {
    FAST = 1,
    NORMAL = 2,
    SLOW = 3,
}

class TimerManager {
    interval: number | undefined = $state(undefined);

    private job: () => void;
    private identifier: number | NodeJS.Timeout | undefined;

    constructor(job: () => void, interval?: number | undefined) {
        this.interval = interval
        this.job = job;
    }

    get isProcessing() {
        return this.interval && this.interval > 0;
    }

    start() {
        if (this.interval && this.interval > 0) {
            this.identifier = setInterval(() => this.job(), this.interval * 1000);
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

export { TimerManager, Speed }