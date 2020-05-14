// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/environs (interfaces: BootstrapEnviron,NetworkingEnviron)

// Package space is a generated GoMock package.
package space

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	constraints "github.com/juju/juju/core/constraints"
	instance "github.com/juju/juju/core/instance"
	network "github.com/juju/juju/core/network"
	environs "github.com/juju/juju/environs"
	config "github.com/juju/juju/environs/config"
	context "github.com/juju/juju/environs/context"
	instances "github.com/juju/juju/environs/instances"
	network0 "github.com/juju/juju/network"
	storage "github.com/juju/juju/storage"
	names_v3 "github.com/juju/names/v4"
	version "github.com/juju/version"
)

// MockBootstrapEnviron is a mock of BootstrapEnviron interface
type MockBootstrapEnviron struct {
	ctrl     *gomock.Controller
	recorder *MockBootstrapEnvironMockRecorder
}

// MockBootstrapEnvironMockRecorder is the mock recorder for MockBootstrapEnviron
type MockBootstrapEnvironMockRecorder struct {
	mock *MockBootstrapEnviron
}

// NewMockBootstrapEnviron creates a new mock instance
func NewMockBootstrapEnviron(ctrl *gomock.Controller) *MockBootstrapEnviron {
	mock := &MockBootstrapEnviron{ctrl: ctrl}
	mock.recorder = &MockBootstrapEnvironMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBootstrapEnviron) EXPECT() *MockBootstrapEnvironMockRecorder {
	return m.recorder
}

// Bootstrap mocks base method
func (m *MockBootstrapEnviron) Bootstrap(arg0 environs.BootstrapContext, arg1 context.ProviderCallContext, arg2 environs.BootstrapParams) (*environs.BootstrapResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Bootstrap", arg0, arg1, arg2)
	ret0, _ := ret[0].(*environs.BootstrapResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Bootstrap indicates an expected call of Bootstrap
func (mr *MockBootstrapEnvironMockRecorder) Bootstrap(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Bootstrap", reflect.TypeOf((*MockBootstrapEnviron)(nil).Bootstrap), arg0, arg1, arg2)
}

// Config mocks base method
func (m *MockBootstrapEnviron) Config() *config.Config {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Config")
	ret0, _ := ret[0].(*config.Config)
	return ret0
}

// Config indicates an expected call of Config
func (mr *MockBootstrapEnvironMockRecorder) Config() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Config", reflect.TypeOf((*MockBootstrapEnviron)(nil).Config))
}

// ConstraintsValidator mocks base method
func (m *MockBootstrapEnviron) ConstraintsValidator(arg0 context.ProviderCallContext) (constraints.Validator, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConstraintsValidator", arg0)
	ret0, _ := ret[0].(constraints.Validator)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ConstraintsValidator indicates an expected call of ConstraintsValidator
func (mr *MockBootstrapEnvironMockRecorder) ConstraintsValidator(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConstraintsValidator", reflect.TypeOf((*MockBootstrapEnviron)(nil).ConstraintsValidator), arg0)
}

// Create mocks base method
func (m *MockBootstrapEnviron) Create(arg0 context.ProviderCallContext, arg1 environs.CreateParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create
func (mr *MockBootstrapEnvironMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockBootstrapEnviron)(nil).Create), arg0, arg1)
}

// Destroy mocks base method
func (m *MockBootstrapEnviron) Destroy(arg0 context.ProviderCallContext) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Destroy", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Destroy indicates an expected call of Destroy
func (mr *MockBootstrapEnvironMockRecorder) Destroy(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Destroy", reflect.TypeOf((*MockBootstrapEnviron)(nil).Destroy), arg0)
}

// DestroyController mocks base method
func (m *MockBootstrapEnviron) DestroyController(arg0 context.ProviderCallContext, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DestroyController", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DestroyController indicates an expected call of DestroyController
func (mr *MockBootstrapEnvironMockRecorder) DestroyController(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DestroyController", reflect.TypeOf((*MockBootstrapEnviron)(nil).DestroyController), arg0, arg1)
}

// PrepareForBootstrap mocks base method
func (m *MockBootstrapEnviron) PrepareForBootstrap(arg0 environs.BootstrapContext, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PrepareForBootstrap", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// PrepareForBootstrap indicates an expected call of PrepareForBootstrap
func (mr *MockBootstrapEnvironMockRecorder) PrepareForBootstrap(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PrepareForBootstrap", reflect.TypeOf((*MockBootstrapEnviron)(nil).PrepareForBootstrap), arg0, arg1)
}

// SetConfig mocks base method
func (m *MockBootstrapEnviron) SetConfig(arg0 *config.Config) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetConfig", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetConfig indicates an expected call of SetConfig
func (mr *MockBootstrapEnvironMockRecorder) SetConfig(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetConfig", reflect.TypeOf((*MockBootstrapEnviron)(nil).SetConfig), arg0)
}

// StorageProvider mocks base method
func (m *MockBootstrapEnviron) StorageProvider(arg0 storage.ProviderType) (storage.Provider, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StorageProvider", arg0)
	ret0, _ := ret[0].(storage.Provider)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StorageProvider indicates an expected call of StorageProvider
func (mr *MockBootstrapEnvironMockRecorder) StorageProvider(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StorageProvider", reflect.TypeOf((*MockBootstrapEnviron)(nil).StorageProvider), arg0)
}

// StorageProviderTypes mocks base method
func (m *MockBootstrapEnviron) StorageProviderTypes() ([]storage.ProviderType, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StorageProviderTypes")
	ret0, _ := ret[0].([]storage.ProviderType)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StorageProviderTypes indicates an expected call of StorageProviderTypes
func (mr *MockBootstrapEnvironMockRecorder) StorageProviderTypes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StorageProviderTypes", reflect.TypeOf((*MockBootstrapEnviron)(nil).StorageProviderTypes))
}

// MockNetworkingEnviron is a mock of NetworkingEnviron interface
type MockNetworkingEnviron struct {
	ctrl     *gomock.Controller
	recorder *MockNetworkingEnvironMockRecorder
}

// MockNetworkingEnvironMockRecorder is the mock recorder for MockNetworkingEnviron
type MockNetworkingEnvironMockRecorder struct {
	mock *MockNetworkingEnviron
}

// NewMockNetworkingEnviron creates a new mock instance
func NewMockNetworkingEnviron(ctrl *gomock.Controller) *MockNetworkingEnviron {
	mock := &MockNetworkingEnviron{ctrl: ctrl}
	mock.recorder = &MockNetworkingEnvironMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockNetworkingEnviron) EXPECT() *MockNetworkingEnvironMockRecorder {
	return m.recorder
}

// AdoptResources mocks base method
func (m *MockNetworkingEnviron) AdoptResources(arg0 context.ProviderCallContext, arg1 string, arg2 version.Number) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AdoptResources", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AdoptResources indicates an expected call of AdoptResources
func (mr *MockNetworkingEnvironMockRecorder) AdoptResources(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AdoptResources", reflect.TypeOf((*MockNetworkingEnviron)(nil).AdoptResources), arg0, arg1, arg2)
}

// AllInstances mocks base method
func (m *MockNetworkingEnviron) AllInstances(arg0 context.ProviderCallContext) ([]instances.Instance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllInstances", arg0)
	ret0, _ := ret[0].([]instances.Instance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllInstances indicates an expected call of AllInstances
func (mr *MockNetworkingEnvironMockRecorder) AllInstances(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllInstances", reflect.TypeOf((*MockNetworkingEnviron)(nil).AllInstances), arg0)
}

// AllRunningInstances mocks base method
func (m *MockNetworkingEnviron) AllRunningInstances(arg0 context.ProviderCallContext) ([]instances.Instance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllRunningInstances", arg0)
	ret0, _ := ret[0].([]instances.Instance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllRunningInstances indicates an expected call of AllRunningInstances
func (mr *MockNetworkingEnvironMockRecorder) AllRunningInstances(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllRunningInstances", reflect.TypeOf((*MockNetworkingEnviron)(nil).AllRunningInstances), arg0)
}

// AllocateContainerAddresses mocks base method
func (m *MockNetworkingEnviron) AllocateContainerAddresses(arg0 context.ProviderCallContext, arg1 instance.Id, arg2 names_v3.MachineTag, arg3 network.InterfaceInfos) (network.InterfaceInfos, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllocateContainerAddresses", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(network.InterfaceInfos)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllocateContainerAddresses indicates an expected call of AllocateContainerAddresses
func (mr *MockNetworkingEnvironMockRecorder) AllocateContainerAddresses(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllocateContainerAddresses", reflect.TypeOf((*MockNetworkingEnviron)(nil).AllocateContainerAddresses), arg0, arg1, arg2, arg3)
}

// AreSpacesRoutable mocks base method
func (m *MockNetworkingEnviron) AreSpacesRoutable(arg0 context.ProviderCallContext, arg1, arg2 *environs.ProviderSpaceInfo) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AreSpacesRoutable", arg0, arg1, arg2)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AreSpacesRoutable indicates an expected call of AreSpacesRoutable
func (mr *MockNetworkingEnvironMockRecorder) AreSpacesRoutable(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AreSpacesRoutable", reflect.TypeOf((*MockNetworkingEnviron)(nil).AreSpacesRoutable), arg0, arg1, arg2)
}

// Bootstrap mocks base method
func (m *MockNetworkingEnviron) Bootstrap(arg0 environs.BootstrapContext, arg1 context.ProviderCallContext, arg2 environs.BootstrapParams) (*environs.BootstrapResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Bootstrap", arg0, arg1, arg2)
	ret0, _ := ret[0].(*environs.BootstrapResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Bootstrap indicates an expected call of Bootstrap
func (mr *MockNetworkingEnvironMockRecorder) Bootstrap(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Bootstrap", reflect.TypeOf((*MockNetworkingEnviron)(nil).Bootstrap), arg0, arg1, arg2)
}

// Config mocks base method
func (m *MockNetworkingEnviron) Config() *config.Config {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Config")
	ret0, _ := ret[0].(*config.Config)
	return ret0
}

// Config indicates an expected call of Config
func (mr *MockNetworkingEnvironMockRecorder) Config() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Config", reflect.TypeOf((*MockNetworkingEnviron)(nil).Config))
}

// ConstraintsValidator mocks base method
func (m *MockNetworkingEnviron) ConstraintsValidator(arg0 context.ProviderCallContext) (constraints.Validator, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConstraintsValidator", arg0)
	ret0, _ := ret[0].(constraints.Validator)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ConstraintsValidator indicates an expected call of ConstraintsValidator
func (mr *MockNetworkingEnvironMockRecorder) ConstraintsValidator(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConstraintsValidator", reflect.TypeOf((*MockNetworkingEnviron)(nil).ConstraintsValidator), arg0)
}

// ControllerInstances mocks base method
func (m *MockNetworkingEnviron) ControllerInstances(arg0 context.ProviderCallContext, arg1 string) ([]instance.Id, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ControllerInstances", arg0, arg1)
	ret0, _ := ret[0].([]instance.Id)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ControllerInstances indicates an expected call of ControllerInstances
func (mr *MockNetworkingEnvironMockRecorder) ControllerInstances(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ControllerInstances", reflect.TypeOf((*MockNetworkingEnviron)(nil).ControllerInstances), arg0, arg1)
}

// Create mocks base method
func (m *MockNetworkingEnviron) Create(arg0 context.ProviderCallContext, arg1 environs.CreateParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create
func (mr *MockNetworkingEnvironMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockNetworkingEnviron)(nil).Create), arg0, arg1)
}

// Destroy mocks base method
func (m *MockNetworkingEnviron) Destroy(arg0 context.ProviderCallContext) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Destroy", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Destroy indicates an expected call of Destroy
func (mr *MockNetworkingEnvironMockRecorder) Destroy(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Destroy", reflect.TypeOf((*MockNetworkingEnviron)(nil).Destroy), arg0)
}

// DestroyController mocks base method
func (m *MockNetworkingEnviron) DestroyController(arg0 context.ProviderCallContext, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DestroyController", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DestroyController indicates an expected call of DestroyController
func (mr *MockNetworkingEnvironMockRecorder) DestroyController(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DestroyController", reflect.TypeOf((*MockNetworkingEnviron)(nil).DestroyController), arg0, arg1)
}

// InstanceTypes mocks base method
func (m *MockNetworkingEnviron) InstanceTypes(arg0 context.ProviderCallContext, arg1 constraints.Value) (instances.InstanceTypesWithCostMetadata, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InstanceTypes", arg0, arg1)
	ret0, _ := ret[0].(instances.InstanceTypesWithCostMetadata)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InstanceTypes indicates an expected call of InstanceTypes
func (mr *MockNetworkingEnvironMockRecorder) InstanceTypes(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InstanceTypes", reflect.TypeOf((*MockNetworkingEnviron)(nil).InstanceTypes), arg0, arg1)
}

// Instances mocks base method
func (m *MockNetworkingEnviron) Instances(arg0 context.ProviderCallContext, arg1 []instance.Id) ([]instances.Instance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Instances", arg0, arg1)
	ret0, _ := ret[0].([]instances.Instance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Instances indicates an expected call of Instances
func (mr *MockNetworkingEnvironMockRecorder) Instances(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Instances", reflect.TypeOf((*MockNetworkingEnviron)(nil).Instances), arg0, arg1)
}

// MaintainInstance mocks base method
func (m *MockNetworkingEnviron) MaintainInstance(arg0 context.ProviderCallContext, arg1 environs.StartInstanceParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MaintainInstance", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// MaintainInstance indicates an expected call of MaintainInstance
func (mr *MockNetworkingEnvironMockRecorder) MaintainInstance(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MaintainInstance", reflect.TypeOf((*MockNetworkingEnviron)(nil).MaintainInstance), arg0, arg1)
}

// NetworkInterfaces mocks base method
func (m *MockNetworkingEnviron) NetworkInterfaces(arg0 context.ProviderCallContext, arg1 []instance.Id) ([]network.InterfaceInfos, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NetworkInterfaces", arg0, arg1)
	ret0, _ := ret[0].([]network.InterfaceInfos)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NetworkInterfaces indicates an expected call of NetworkInterfaces
func (mr *MockNetworkingEnvironMockRecorder) NetworkInterfaces(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NetworkInterfaces", reflect.TypeOf((*MockNetworkingEnviron)(nil).NetworkInterfaces), arg0, arg1)
}

// PrecheckInstance mocks base method
func (m *MockNetworkingEnviron) PrecheckInstance(arg0 context.ProviderCallContext, arg1 environs.PrecheckInstanceParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PrecheckInstance", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// PrecheckInstance indicates an expected call of PrecheckInstance
func (mr *MockNetworkingEnvironMockRecorder) PrecheckInstance(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PrecheckInstance", reflect.TypeOf((*MockNetworkingEnviron)(nil).PrecheckInstance), arg0, arg1)
}

// PrepareForBootstrap mocks base method
func (m *MockNetworkingEnviron) PrepareForBootstrap(arg0 environs.BootstrapContext, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PrepareForBootstrap", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// PrepareForBootstrap indicates an expected call of PrepareForBootstrap
func (mr *MockNetworkingEnvironMockRecorder) PrepareForBootstrap(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PrepareForBootstrap", reflect.TypeOf((*MockNetworkingEnviron)(nil).PrepareForBootstrap), arg0, arg1)
}

// Provider mocks base method
func (m *MockNetworkingEnviron) Provider() environs.EnvironProvider {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Provider")
	ret0, _ := ret[0].(environs.EnvironProvider)
	return ret0
}

// Provider indicates an expected call of Provider
func (mr *MockNetworkingEnvironMockRecorder) Provider() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Provider", reflect.TypeOf((*MockNetworkingEnviron)(nil).Provider))
}

// ProviderSpaceInfo mocks base method
func (m *MockNetworkingEnviron) ProviderSpaceInfo(arg0 context.ProviderCallContext, arg1 *network.SpaceInfo) (*environs.ProviderSpaceInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProviderSpaceInfo", arg0, arg1)
	ret0, _ := ret[0].(*environs.ProviderSpaceInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ProviderSpaceInfo indicates an expected call of ProviderSpaceInfo
func (mr *MockNetworkingEnvironMockRecorder) ProviderSpaceInfo(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProviderSpaceInfo", reflect.TypeOf((*MockNetworkingEnviron)(nil).ProviderSpaceInfo), arg0, arg1)
}

// ReleaseContainerAddresses mocks base method
func (m *MockNetworkingEnviron) ReleaseContainerAddresses(arg0 context.ProviderCallContext, arg1 []network0.ProviderInterfaceInfo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReleaseContainerAddresses", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReleaseContainerAddresses indicates an expected call of ReleaseContainerAddresses
func (mr *MockNetworkingEnvironMockRecorder) ReleaseContainerAddresses(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReleaseContainerAddresses", reflect.TypeOf((*MockNetworkingEnviron)(nil).ReleaseContainerAddresses), arg0, arg1)
}

// SSHAddresses mocks base method
func (m *MockNetworkingEnviron) SSHAddresses(arg0 context.ProviderCallContext, arg1 network.SpaceAddresses) (network.SpaceAddresses, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SSHAddresses", arg0, arg1)
	ret0, _ := ret[0].(network.SpaceAddresses)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SSHAddresses indicates an expected call of SSHAddresses
func (mr *MockNetworkingEnvironMockRecorder) SSHAddresses(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SSHAddresses", reflect.TypeOf((*MockNetworkingEnviron)(nil).SSHAddresses), arg0, arg1)
}

// SetConfig mocks base method
func (m *MockNetworkingEnviron) SetConfig(arg0 *config.Config) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetConfig", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetConfig indicates an expected call of SetConfig
func (mr *MockNetworkingEnvironMockRecorder) SetConfig(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetConfig", reflect.TypeOf((*MockNetworkingEnviron)(nil).SetConfig), arg0)
}

// Spaces mocks base method
func (m *MockNetworkingEnviron) Spaces(arg0 context.ProviderCallContext) ([]network.SpaceInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Spaces", arg0)
	ret0, _ := ret[0].([]network.SpaceInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Spaces indicates an expected call of Spaces
func (mr *MockNetworkingEnvironMockRecorder) Spaces(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Spaces", reflect.TypeOf((*MockNetworkingEnviron)(nil).Spaces), arg0)
}

// StartInstance mocks base method
func (m *MockNetworkingEnviron) StartInstance(arg0 context.ProviderCallContext, arg1 environs.StartInstanceParams) (*environs.StartInstanceResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StartInstance", arg0, arg1)
	ret0, _ := ret[0].(*environs.StartInstanceResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StartInstance indicates an expected call of StartInstance
func (mr *MockNetworkingEnvironMockRecorder) StartInstance(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartInstance", reflect.TypeOf((*MockNetworkingEnviron)(nil).StartInstance), arg0, arg1)
}

// StopInstances mocks base method
func (m *MockNetworkingEnviron) StopInstances(arg0 context.ProviderCallContext, arg1 ...instance.Id) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "StopInstances", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// StopInstances indicates an expected call of StopInstances
func (mr *MockNetworkingEnvironMockRecorder) StopInstances(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StopInstances", reflect.TypeOf((*MockNetworkingEnviron)(nil).StopInstances), varargs...)
}

// StorageProvider mocks base method
func (m *MockNetworkingEnviron) StorageProvider(arg0 storage.ProviderType) (storage.Provider, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StorageProvider", arg0)
	ret0, _ := ret[0].(storage.Provider)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StorageProvider indicates an expected call of StorageProvider
func (mr *MockNetworkingEnvironMockRecorder) StorageProvider(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StorageProvider", reflect.TypeOf((*MockNetworkingEnviron)(nil).StorageProvider), arg0)
}

// StorageProviderTypes mocks base method
func (m *MockNetworkingEnviron) StorageProviderTypes() ([]storage.ProviderType, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StorageProviderTypes")
	ret0, _ := ret[0].([]storage.ProviderType)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StorageProviderTypes indicates an expected call of StorageProviderTypes
func (mr *MockNetworkingEnvironMockRecorder) StorageProviderTypes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StorageProviderTypes", reflect.TypeOf((*MockNetworkingEnviron)(nil).StorageProviderTypes))
}

// Subnets mocks base method
func (m *MockNetworkingEnviron) Subnets(arg0 context.ProviderCallContext, arg1 instance.Id, arg2 []network.Id) ([]network.SubnetInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Subnets", arg0, arg1, arg2)
	ret0, _ := ret[0].([]network.SubnetInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Subnets indicates an expected call of Subnets
func (mr *MockNetworkingEnvironMockRecorder) Subnets(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subnets", reflect.TypeOf((*MockNetworkingEnviron)(nil).Subnets), arg0, arg1, arg2)
}

// SuperSubnets mocks base method
func (m *MockNetworkingEnviron) SuperSubnets(arg0 context.ProviderCallContext) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SuperSubnets", arg0)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SuperSubnets indicates an expected call of SuperSubnets
func (mr *MockNetworkingEnvironMockRecorder) SuperSubnets(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SuperSubnets", reflect.TypeOf((*MockNetworkingEnviron)(nil).SuperSubnets), arg0)
}

// SupportsContainerAddresses mocks base method
func (m *MockNetworkingEnviron) SupportsContainerAddresses(arg0 context.ProviderCallContext) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SupportsContainerAddresses", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SupportsContainerAddresses indicates an expected call of SupportsContainerAddresses
func (mr *MockNetworkingEnvironMockRecorder) SupportsContainerAddresses(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SupportsContainerAddresses", reflect.TypeOf((*MockNetworkingEnviron)(nil).SupportsContainerAddresses), arg0)
}

// SupportsSpaceDiscovery mocks base method
func (m *MockNetworkingEnviron) SupportsSpaceDiscovery(arg0 context.ProviderCallContext) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SupportsSpaceDiscovery", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SupportsSpaceDiscovery indicates an expected call of SupportsSpaceDiscovery
func (mr *MockNetworkingEnvironMockRecorder) SupportsSpaceDiscovery(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SupportsSpaceDiscovery", reflect.TypeOf((*MockNetworkingEnviron)(nil).SupportsSpaceDiscovery), arg0)
}

// SupportsSpaces mocks base method
func (m *MockNetworkingEnviron) SupportsSpaces(arg0 context.ProviderCallContext) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SupportsSpaces", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SupportsSpaces indicates an expected call of SupportsSpaces
func (mr *MockNetworkingEnvironMockRecorder) SupportsSpaces(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SupportsSpaces", reflect.TypeOf((*MockNetworkingEnviron)(nil).SupportsSpaces), arg0)
}
