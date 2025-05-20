package service

import (
	"testing"

	mocks "github.com/openhdc/otterscale/internal/domain/service/mocks"
	"go.uber.org/mock/gomock"
)

func TestNexusService_FieldAssignment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	server := mocks.NewMockMAASServer(ctrl)
	pkgRepo := mocks.NewMockMAASPackageRepository(ctrl)
	bootResource := mocks.NewMockMAASBootResource(ctrl)
	bootSource := mocks.NewMockMAASBootSource(ctrl)
	bootSourceSelection := mocks.NewMockMAASBootSourceSelection(ctrl)
	fabric := mocks.NewMockMAASFabric(ctrl)
	vlan := mocks.NewMockMAASVLAN(ctrl)
	subnet := mocks.NewMockMAASSubnet(ctrl)
	ipRange := mocks.NewMockMAASIPRange(ctrl)
	machine := mocks.NewMockMAASMachine(ctrl)
	tag := mocks.NewMockMAASTag(ctrl)
	sshKey := mocks.NewMockMAASSSHKey(ctrl)
	keyManager := mocks.NewMockJujuKey(ctrl)
	client := mocks.NewMockJujuClient(ctrl)
	machineManager := mocks.NewMockJujuMachine(ctrl)
	scope := mocks.NewMockJujuModel(ctrl)
	scopeConfig := mocks.NewMockJujuModelConfig(ctrl)
	facility := mocks.NewMockJujuApplication(ctrl)
	action := mocks.NewMockJujuAction(ctrl)
	charmhub := mocks.NewMockJujuCharmHub(ctrl)
	kubernetes := mocks.NewMockKubeClient(ctrl)
	apps := mocks.NewMockKubeApps(ctrl)
	batch := mocks.NewMockKubeBatch(ctrl)
	core := mocks.NewMockKubeCore(ctrl)
	storage := mocks.NewMockKubeStorage(ctrl)
	helm := mocks.NewMockKubeHelm(ctrl)

	ns := &NexusService{
		server:              server,
		packageRepository:   pkgRepo,
		bootResource:        bootResource,
		bootSource:          bootSource,
		bootSourceSelection: bootSourceSelection,
		fabric:              fabric,
		vlan:                vlan,
		subnet:              subnet,
		ipRange:             ipRange,
		machine:             machine,
		tag:                 tag,
		sshKey:              sshKey,
		keyManager:          keyManager,
		client:              client,
		machineManager:      machineManager,
		scope:               scope,
		scopeConfig:         scopeConfig,
		facility:            facility,
		action:              action,
		charmhub:            charmhub,
		kubernetes:          kubernetes,
		apps:                apps,
		batch:               batch,
		core:                core,
		storage:             storage,
		helm:                helm,
	}

	if ns.server != server {
		t.Error("server field not assigned correctly")
	}
	if ns.packageRepository != pkgRepo {
		t.Error("packageRepository field not assigned correctly")
	}
	if ns.bootResource != bootResource {
		t.Error("bootResource field not assigned correctly")
	}
	if ns.bootSource != bootSource {
		t.Error("bootSource field not assigned correctly")
	}
	if ns.bootSourceSelection != bootSourceSelection {
		t.Error("bootSourceSelection field not assigned correctly")
	}
	if ns.fabric != fabric {
		t.Error("fabric field not assigned correctly")
	}
	if ns.vlan != vlan {
		t.Error("vlan field not assigned correctly")
	}
	if ns.subnet != subnet {
		t.Error("subnet field not assigned correctly")
	}
	if ns.ipRange != ipRange {
		t.Error("ipRange field not assigned correctly")
	}
	if ns.machine != machine {
		t.Error("machine field not assigned correctly")
	}
	if ns.tag != tag {
		t.Error("tag field not assigned correctly")
	}
	if ns.sshKey != sshKey {
		t.Error("sshKey field not assigned correctly")
	}
	if ns.keyManager != keyManager {
		t.Error("keyManager field not assigned correctly")
	}
	if ns.client != client {
		t.Error("client field not assigned correctly")
	}
	if ns.machineManager != machineManager {
		t.Error("machineManager field not assigned correctly")
	}
	if ns.scope != scope {
		t.Error("scope field not assigned correctly")
	}
	if ns.scopeConfig != scopeConfig {
		t.Error("scopeConfig field not assigned correctly")
	}
	if ns.facility != facility {
		t.Error("facility field not assigned correctly")
	}
	if ns.action != action {
		t.Error("action field not assigned correctly")
	}
	if ns.charmhub != charmhub {
		t.Error("charmhub field not assigned correctly")
	}
	if ns.kubernetes != kubernetes {
		t.Error("kubernetes field not assigned correctly")
	}
	if ns.apps != apps {
		t.Error("apps field not assigned correctly")
	}
	if ns.batch != batch {
		t.Error("batch field not assigned correctly")
	}
	if ns.core != core {
		t.Error("core field not assigned correctly")
	}
	if ns.storage != storage {
		t.Error("storage field not assigned correctly")
	}
	if ns.helm != helm {
		t.Error("helm field not assigned correctly")
	}
}

func TestNexusService_InterfaceSatisfies(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	_ = mocks.NewMockMAASServer(ctrl)
	_ = mocks.NewMockMAASPackageRepository(ctrl)
	_ = mocks.NewMockMAASBootResource(ctrl)
	_ = mocks.NewMockMAASBootSource(ctrl)
	_ = mocks.NewMockMAASBootSourceSelection(ctrl)
	_ = mocks.NewMockMAASFabric(ctrl)
	_ = mocks.NewMockMAASVLAN(ctrl)
	_ = mocks.NewMockMAASSubnet(ctrl)
	_ = mocks.NewMockMAASIPRange(ctrl)
	_ = mocks.NewMockMAASMachine(ctrl)
	_ = mocks.NewMockMAASTag(ctrl)
	_ = mocks.NewMockMAASSSHKey(ctrl)
	_ = mocks.NewMockJujuKey(ctrl)
	_ = mocks.NewMockJujuMachine(ctrl)
	_ = mocks.NewMockJujuClient(ctrl)
	_ = mocks.NewMockJujuModel(ctrl)
	_ = mocks.NewMockJujuModelConfig(ctrl)
	_ = mocks.NewMockJujuApplication(ctrl)
	_ = mocks.NewMockJujuAction(ctrl)
	_ = mocks.NewMockJujuCharmHub(ctrl)
	_ = mocks.NewMockKubeClient(ctrl)
	_ = mocks.NewMockKubeApps(ctrl)
	_ = mocks.NewMockKubeBatch(ctrl)
	_ = mocks.NewMockKubeCore(ctrl)
	_ = mocks.NewMockKubeStorage(ctrl)
	_ = mocks.NewMockKubeHelm(ctrl)
}
