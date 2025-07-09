import type { StepManagerState } from './types';

class StepManager {
    step = $state(1);
    state: { isUpdating: boolean };
    steps: number

    constructor(steps: number, state: StepManagerState) {
        this.steps = steps;
        this.state = state;
    }

    get isFirstStep() {
        return this.step === 1
    }

    get isLastStep() {
        return this.step === this.steps
    }

    areStepsActive(value: number) {
        return this.step >= value
    }

    isStepActive(value: number) {
        return this.step == value
    }

    async update(value: number) {
        if (this.state.isUpdating) return;

        if (this.step === value) return;

        const direction = this.step < value ? 1 : -1;
        this.state.isUpdating = true;
        while (this.step !== value) {
            this.step = this.step + direction;
            await new Promise((response) => setTimeout(response, 100));
        }
        this.state.isUpdating = false;
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

export { IndexManager, StepManager };
