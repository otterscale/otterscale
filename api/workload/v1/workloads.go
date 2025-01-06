package workload

type Workloads struct {
	srcs []*Workload
	dsts []*Workload
	trfs []*Workload
}

func (w *Workloads) Sources() []*Workload {
	return w.srcs
}

func (w *Workloads) Destinations() []*Workload {
	return w.dsts
}

func (w *Workloads) Transformers() []*Workload {
	return w.trfs
}

func (w *Workloads) AppendSources(wls ...*Workload) {
	w.srcs = append(w.srcs, wls...)
}

func (w *Workloads) AppendDestinations(wls ...*Workload) {
	w.dsts = append(w.dsts, wls...)
}

func (w *Workloads) AppendTransformers(wls ...*Workload) {
	w.trfs = append(w.trfs, wls...)
}
