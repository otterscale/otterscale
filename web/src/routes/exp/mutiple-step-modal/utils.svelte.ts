class StepManager {
    step = $state(1);
    steps: number

    constructor(steps: number) {
        this.steps = steps;
    }

    get isFirstStep() {
        return this.step === 1
    }

    get isLastStep() {
        return this.step === this.steps
    }

    isStepActive(value: number) {
        return this.step >= value
    }

    update(value: number) {
        this.step = value
    }

    next() {
        this.step = Math.min(this.steps, this.step + 1)
    }

    back() {
        this.step = Math.max(1, this.step - 1)
    }

    reset() {
        this.step = 1
    }
}

class IndexManager {
    index = 0;

    get() {
        this.index = this.index + 1
        return this.index
    }

    reset() {
        this.index = 0
    }
}

export { StepManager, IndexManager }