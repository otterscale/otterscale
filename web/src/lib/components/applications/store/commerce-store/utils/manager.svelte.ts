import type { Application_Chart } from "$lib/api/application/v1/application_pb";

class FilterManager {
    charts: Application_Chart[] = $state([])

    searchedName: string = $state('')
    selectedKeywords: string[] = $state([])
    selectedMaintainerNames: string[] = $state([])
    selectedLicences: string[] = $state([])
    selectedDeprecation: boolean | null = $state(false)

    isFiltered = $derived(this.searchedName !== '' || this.selectedKeywords.length > 0 || this.selectedMaintainerNames.length > 0 || this.selectedLicences.length > 0 || this.selectedDeprecation !== false)
    filteredCharts = $derived(
        this.charts
            .filter((chart) => (this.searchedName ? chart.name.includes(this.searchedName) : true))
            .filter((chart) => (this.selectedKeywords.length ? chart.keywords.some((keyword) => (this.selectedKeywords.includes(keyword))) : true))
            .filter((chart) => (this.selectedMaintainerNames.length ? chart.maintainers.some((maintainer) => (this.selectedMaintainerNames.includes(maintainer.name))) : true))
            .filter((chart) => (this.selectedLicences.length ? this.selectedLicences.includes(chart.license) : true))
            .filter((chart) => {
                if (this.selectedDeprecation === null) { return true }
                else if (this.selectedDeprecation === true) { return chart.deprecated }
                else if (this.selectedDeprecation === false) { return !chart.deprecated }
            })
    )

    constructor(charts: Application_Chart[]) {
        this.charts = charts
    }

    isKeywordSelected(keyword: string) {
        return this.selectedKeywords.includes(keyword)
    }
    isMaintainerSelected(maintainerName: string) {
        return this.selectedMaintainerNames.includes(maintainerName)
    }
    isLicenceSelected(licence: string) {
        return this.selectedLicences.includes(licence)
    }
    isDeprecationSelected(deprecations: boolean | null) {
        return this.selectedDeprecation === deprecations
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
    toggleLicence(licence: string) {
        if (this.isLicenceSelected(licence)) {
            this.selectedLicences =
                this.selectedLicences.filter(
                    (selectedLicence) => selectedLicence !== licence
                );
        } else {
            this.selectedLicences.push(licence);
        }
    }
    toggleDeprecation(deprecations: boolean | null) {
        this.selectedDeprecation = deprecations
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
    resetLicence() {
        this.selectedLicences = []
    }
    resetDeprecatedr() {
        this.selectedDeprecation = false
    }
    reset() {
        this.resetName()
        this.resetKeyword()
        this.resetMaintainer()
        this.resetLicence()
        this.resetDeprecatedr()
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
