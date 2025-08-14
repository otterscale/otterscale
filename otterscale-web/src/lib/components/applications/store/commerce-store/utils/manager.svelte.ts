import type { Application_Chart } from "$lib/api/application/v1/application_pb";

class FilterManager {
    charts: Application_Chart[]

    searchedName: string = $state('')
    selectedKeywords: string[] = $state([])
    selectedMaintainerNames: string[] = $state([])

    constructor(charts: Application_Chart[]) {
        this.charts = charts
    }

    get isFiltered() {
        if (this.searchedName !== '' || this.selectedKeywords.length > 0 || this.selectedMaintainerNames.length > 0) {
            return true
        }
    }

    get filteredCharts() {
        return this.charts
            .filter((chart) => (this.searchedName ? chart.name.includes(this.searchedName) : true))
            .filter((chart) => (this.selectedKeywords.length ? chart.keywords.some((keyword) => (this.selectedKeywords.includes(keyword))) : true))
            .filter((chart) => (this.selectedMaintainerNames.length ? chart.maintainers.some((maintainer) => (this.selectedMaintainerNames.includes(maintainer.name))) : true))
    }

    isKeywordSelected(keyword: string) {
        return this.selectedKeywords.includes(keyword)
    }
    isMaintainerSelected(maintainerName: string) {
        return this.selectedMaintainerNames.includes(maintainerName)
    }

    toggleKeyword(keyword: string) {
        if (this.isKeywordSelected(keyword)) {
            this.selectedKeywords = this.selectedKeywords.filter(
                (selectedKeyword) => selectedKeyword !== keyword
            );
        } else {
            this.selectedKeywords.push(keyword);
        }
    }
    toggleMaintainer(maintainerName: string) {
        if (this.isMaintainerSelected(maintainerName)) {
            this.selectedMaintainerNames =
                this.selectedMaintainerNames.filter(
                    (selectedMaintainerName) => selectedMaintainerName !== maintainerName
                );
        } else {
            this.selectedMaintainerNames.push(maintainerName);
        }
    }

    resetName() {
        this.searchedName = ''
    }
    resetKeyword() {
        this.selectedKeywords = []
    }
    resetMaintainer() {
        this.selectedMaintainerNames = []
    }
    reset() {
        this.resetName()
        this.resetKeyword()
        this.resetMaintainer()
    }
}

class PaginationManager {
    activePage: number = $state(0)
    count: number = $state(0)
    perPage: number = 6
    siblingCount: number = 1

    constructor(charts: Application_Chart[]) {
        this.count = charts.length
    }
}

export { FilterManager, PaginationManager };
