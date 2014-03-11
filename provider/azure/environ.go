// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package azure

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"launchpad.net/gwacl"

	"launchpad.net/juju-core/constraints"
	"launchpad.net/juju-core/environs"
	"launchpad.net/juju-core/environs/cloudinit"
	"launchpad.net/juju-core/environs/config"
	"launchpad.net/juju-core/environs/imagemetadata"
	"launchpad.net/juju-core/environs/instances"
	"launchpad.net/juju-core/environs/simplestreams"
	"launchpad.net/juju-core/environs/storage"
	envtools "launchpad.net/juju-core/environs/tools"
	"launchpad.net/juju-core/instance"
	"launchpad.net/juju-core/provider/common"
	"launchpad.net/juju-core/state"
	"launchpad.net/juju-core/state/api"
	"launchpad.net/juju-core/tools"
	"launchpad.net/juju-core/utils/parallel"
	"launchpad.net/juju-core/utils/set"
)

const (
	// deploymentSlot says in which slot to deploy instances.  Azure
	// supports 'Production' or 'Staging'.
	// This provider always deploys to Production.  Think twice about
	// changing that: DNS names in the staging slot work differently from
	// those in the production slot.  In Staging, Azure assigns an
	// arbitrary hostname that we can then extract from the deployment's
	// URL.  In Production, the hostname in the deployment URL does not
	// actually seem to resolve; instead, the service name is used as the
	// DNS name, with ".cloudapp.net" appended.
	deploymentSlot = "Production"

	// Address space of the virtual network used by the nodes in this
	// environement, in CIDR notation. This is the network used for
	// machine-to-machine communication.
	networkDefinition = "10.0.0.0/8"
)

type azureEnviron struct {
	// Except where indicated otherwise, all fields in this object should
	// only be accessed using a lock or a snapshot.
	sync.Mutex

	// name is immutable; it does not need locking.
	name string

	// ecfg is the environment's Azure-specific configuration.
	ecfg *azureEnvironConfig

	// storage is this environ's own private storage.
	storage storage.Storage

	// storageAccountKey holds an access key to this environment's
	// private storage.  This is automatically queried from Azure on
	// startup.
	storageAccountKey string
}

// azureEnviron implements Environ and HasRegion.
var _ environs.Environ = (*azureEnviron)(nil)
var _ simplestreams.HasRegion = (*azureEnviron)(nil)
var _ imagemetadata.SupportsCustomSources = (*azureEnviron)(nil)
var _ envtools.SupportsCustomSources = (*azureEnviron)(nil)

// NewEnviron creates a new azureEnviron.
func NewEnviron(cfg *config.Config) (*azureEnviron, error) {
	env := azureEnviron{name: cfg.Name()}
	err := env.SetConfig(cfg)
	if err != nil {
		return nil, err
	}

	// Set up storage.
	env.storage = &azureStorage{
		storageContext: &environStorageContext{environ: &env},
	}
	return &env, nil
}

// extractStorageKey returns the primary account key from a gwacl
// StorageAccountKeys struct, or if there is none, the secondary one.
func extractStorageKey(keys *gwacl.StorageAccountKeys) string {
	if keys.Primary != "" {
		return keys.Primary
	}
	return keys.Secondary
}

// queryStorageAccountKey retrieves the storage account's key from Azure.
func (env *azureEnviron) queryStorageAccountKey() (string, error) {
	azure, err := env.getManagementAPI()
	if err != nil {
		return "", err
	}
	defer env.releaseManagementAPI(azure)

	accountName := env.getSnapshot().ecfg.storageAccountName()
	keys, err := azure.GetStorageAccountKeys(accountName)
	if err != nil {
		return "", fmt.Errorf("cannot obtain storage account keys: %v", err)
	}

	key := extractStorageKey(keys)
	if key == "" {
		return "", fmt.Errorf("no keys available for storage account")
	}

	return key, nil
}

// Name is specified in the Environ interface.
func (env *azureEnviron) Name() string {
	return env.name
}

// getSnapshot produces an atomic shallow copy of the environment object.
// Whenever you need to access the environment object's fields without
// modifying them, get a snapshot and read its fields instead.  You will
// get a consistent view of the fields without any further locking.
// If you do need to modify the environment's fields, do not get a snapshot
// but lock the object throughout the critical section.
func (env *azureEnviron) getSnapshot() *azureEnviron {
	env.Lock()
	defer env.Unlock()

	// Copy the environment.  (Not the pointer, the environment itself.)
	// This is a shallow copy.
	snap := *env
	// Reset the snapshot's mutex, because we just copied it while we
	// were holding it.  The snapshot will have a "clean," unlocked mutex.
	snap.Mutex = sync.Mutex{}
	return &snap
}

// getAffinityGroupName returns the name of the affinity group used by all
// the Services in this environment.
func (env *azureEnviron) getAffinityGroupName() string {
	return env.getEnvPrefix() + "ag"
}

func (env *azureEnviron) createAffinityGroup() error {
	affinityGroupName := env.getAffinityGroupName()
	azure, err := env.getManagementAPI()
	if err != nil {
		return err
	}
	defer env.releaseManagementAPI(azure)
	snap := env.getSnapshot()
	location := snap.ecfg.location()
	cag := gwacl.NewCreateAffinityGroup(affinityGroupName, affinityGroupName, affinityGroupName, location)
	return azure.CreateAffinityGroup(&gwacl.CreateAffinityGroupRequest{
		CreateAffinityGroup: cag})
}

func (env *azureEnviron) deleteAffinityGroup() error {
	affinityGroupName := env.getAffinityGroupName()
	azure, err := env.getManagementAPI()
	if err != nil {
		return err
	}
	defer env.releaseManagementAPI(azure)
	return azure.DeleteAffinityGroup(&gwacl.DeleteAffinityGroupRequest{
		Name: affinityGroupName})
}

// getVirtualNetworkName returns the name of the virtual network used by all
// the VMs in this environment.
func (env *azureEnviron) getVirtualNetworkName() string {
	return env.getEnvPrefix() + "vnet"
}

func (env *azureEnviron) createVirtualNetwork() error {
	vnetName := env.getVirtualNetworkName()
	affinityGroupName := env.getAffinityGroupName()
	azure, err := env.getManagementAPI()
	if err != nil {
		return err
	}
	defer env.releaseManagementAPI(azure)
	virtualNetwork := gwacl.VirtualNetworkSite{
		Name:          vnetName,
		AffinityGroup: affinityGroupName,
		AddressSpacePrefixes: []string{
			networkDefinition,
		},
	}
	return azure.AddVirtualNetworkSite(&virtualNetwork)
}

func (env *azureEnviron) deleteVirtualNetwork() error {
	azure, err := env.getManagementAPI()
	if err != nil {
		return err
	}
	defer env.releaseManagementAPI(azure)
	vnetName := env.getVirtualNetworkName()
	return azure.RemoveVirtualNetworkSite(vnetName)
}

// getContainerName returns the name of the private storage account container
// that this environment is using.
func (env *azureEnviron) getContainerName() string {
	return env.getEnvPrefix() + "private"
}

// Bootstrap is specified in the Environ interface.
func (env *azureEnviron) Bootstrap(ctx environs.BootstrapContext, cons constraints.Value) (err error) {
	// The creation of the affinity group and the virtual network is specific to the Azure provider.
	err = env.createAffinityGroup()
	if err != nil {
		return err
	}
	// If we fail after this point, clean up the affinity group.
	defer func() {
		if err != nil {
			env.deleteAffinityGroup()
		}
	}()
	err = env.createVirtualNetwork()
	if err != nil {
		return err
	}
	// If we fail after this point, clean up the virtual network.
	defer func() {
		if err != nil {
			env.deleteVirtualNetwork()
		}
	}()
	err = common.Bootstrap(ctx, env, cons)
	return err
}

// StateInfo is specified in the Environ interface.
func (env *azureEnviron) StateInfo() (*state.Info, *api.Info, error) {
	return common.StateInfo(env)
}

// Config is specified in the Environ interface.
func (env *azureEnviron) Config() *config.Config {
	snap := env.getSnapshot()
	return snap.ecfg.Config
}

// SetConfig is specified in the Environ interface.
func (env *azureEnviron) SetConfig(cfg *config.Config) error {
	ecfg, err := azureEnvironProvider{}.newConfig(cfg)
	if err != nil {
		return err
	}

	env.Lock()
	defer env.Unlock()

	if env.ecfg != nil {
		_, err = azureEnvironProvider{}.Validate(cfg, env.ecfg.Config)
		if err != nil {
			return err
		}
	}

	env.ecfg = ecfg

	// Reset storage account key.  Even if we had one before, it may not
	// be appropriate for the new config.
	env.storageAccountKey = ""

	return nil
}

// attemptCreateService tries to create a new hosted service on Azure, with a
// name it chooses (based on the given prefix), but recognizes that the name
// may not be available.  If the name is not available, it does not treat that
// as an error but just returns nil.
//
// If label is non-empty, it will be used for the cloud service's label;
// otherwise, the randomly generated name will be used.
func attemptCreateService(azure *gwacl.ManagementAPI, prefix, affinityGroupName, label string) (*gwacl.CreateHostedService, error) {
	var err error
	name := gwacl.MakeRandomHostedServiceName(prefix)
	err = azure.CheckHostedServiceNameAvailability(name)
	if err != nil {
		// The calling function should retry.
		return nil, nil
	}
	if label == "" {
		label = name
	}
	req := gwacl.NewCreateHostedServiceWithLocation(name, label, "")
	req.AffinityGroup = affinityGroupName
	err = azure.AddHostedService(req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

// architectures lists the CPU architectures supported by Azure.
var architectures = []string{"amd64", "i386"}

// newHostedService creates a hosted service.  It will make up a unique name,
// starting with the given prefix.
func newHostedService(azure *gwacl.ManagementAPI, prefix, affinityGroupName, label string) (*gwacl.CreateHostedService, error) {
	var err error
	var svc *gwacl.CreateHostedService
	for tries := 10; tries > 0 && err == nil && svc == nil; tries-- {
		svc, err = attemptCreateService(azure, prefix, affinityGroupName, label)
	}
	if err != nil {
		return nil, fmt.Errorf("could not create hosted service: %v", err)
	}
	if svc == nil {
		return nil, fmt.Errorf("could not come up with a unique hosted service name - is your randomizer initialized?")
	}
	return svc, nil
}

// selectInstanceTypeAndImage returns the appropriate instance-type name and
// the OS image name for launching a virtual machine with the given parameters.
func (env *azureEnviron) selectInstanceTypeAndImage(cons constraints.Value, series, location string) (string, string, error) {
	ecfg := env.getSnapshot().ecfg
	sourceImageName := ecfg.forceImageName()
	if sourceImageName != "" {
		// Configuration forces us to use a specific image.  There may
		// not be a suitable image in the simplestreams database.
		// This means we can't use Juju's normal selection mechanism,
		// because it combines instance-type and image selection: if
		// there are no images we can use, it won't offer us an
		// instance type either.
		//
		// Select the instance type using simple, Azure-specific code.
		machineType, err := selectMachineType(gwacl.RoleSizes, defaultToBaselineSpec(cons))
		if err != nil {
			return "", "", err
		}
		return machineType.Name, sourceImageName, nil
	}

	// Choose the most suitable instance type and OS image, based on
	// simplestreams information.
	//
	// This should be the normal execution path.  The user is not expected
	// to configure a source image name in normal use.
	constraint := instances.InstanceConstraint{
		Region:      location,
		Series:      series,
		Arches:      architectures,
		Constraints: cons,
	}
	spec, err := findInstanceSpec(env, constraint)
	if err != nil {
		return "", "", err
	}
	return spec.InstanceType.Id, spec.Image.Id, nil
}

// ensureCloudService returns the cloud service with
// the specified label, or creates one if there is none
// or no label is specified.
func (env *azureEnviron) ensureCloudService(azure *gwacl.ManagementAPI, label string) (service *gwacl.HostedService, serviceLabel string, err error) {
	affinityGroup := env.getAffinityGroupName()
	var serviceName string
	if label != "" {
		labelBase64 := base64.StdEncoding.EncodeToString([]byte(label))
		services, err := azure.ListHostedServices()
		if err != nil {
			return nil, "", err
		}
		for _, service := range services {
			if service.AffinityGroup != affinityGroup {
				continue
			}
			if service.Label == labelBase64 {
				serviceName = service.ServiceName
				break
			}
		}
	}
	if serviceName == "" {
		createdService, err := newHostedService(azure, env.getEnvPrefix(), affinityGroup, label)
		if err != nil {
			return nil, "", err
		}
		serviceName = createdService.ServiceName
	}
	service, err = azure.GetHostedServiceProperties(serviceName, true)
	if err != nil {
		return nil, "", err
	}
	if label == "" {
		labelBytes, err := base64.StdEncoding.DecodeString(service.Label)
		if err != nil {
			return nil, "", err
		}
		label = string(labelBytes)
	}
	return service, label, nil
}

func (env *azureEnviron) createRole(azure *gwacl.ManagementAPI, role *gwacl.Role, label string) (resultInst instance.Instance, resultErr error) {
	var inst instance.Instance
	defer func() {
		if inst != nil && resultErr != nil {
			if err := env.StopInstances([]instance.Instance{inst}); err != nil {
				// Failure upon failure. Log it, but return the original error.
				logger.Errorf("error releasing failed instance: %v", err)
			}
		}
	}()
	service, label, err := env.ensureCloudService(azure, label)
	if err != nil {
		return nil, err
	}
	if len(service.Deployments) == 0 {
		// This is a newly created cloud service, so we
		// should destroy it if anything below fails.
		defer func() {
			if resultErr != nil {
				azure.DestroyHostedService(&gwacl.DestroyHostedServiceRequest{
					ServiceName: service.ServiceName,
				})
				// Destroying the hosted service destroys the instance,
				// so ensure StopInstances isn't called.
				inst = nil
			}
		}()
		// Create an initial deployment.
		deployment := gwacl.NewDeploymentForCreateVMDeployment(
			deploymentNameV2(service.ServiceName),
			deploymentSlot,
			label,
			[]gwacl.Role{*role},
			env.getVirtualNetworkName(),
		)
		if err := azure.AddDeployment(deployment, service.ServiceName); err != nil {
			return nil, err
		}
		service.Deployments = append(service.Deployments, *deployment)
	} else {
		// Update the deployment.
		deployment := &service.Deployments[0]
		if err := azure.AddRole(&gwacl.AddRoleRequest{
			ServiceName:      service.ServiceName,
			DeploymentName:   deployment.Name,
			PersistentVMRole: (*gwacl.PersistentVMRole)(role),
		}); err != nil {
			return nil, err
		}
		deployment.RoleList = append(deployment.RoleList, *role)
	}
	return env.getInstance(service, role.RoleName)
}

// deploymentNameV1 returns the deployment name used
// in the original implementation of the Azure provider.
func deploymentNameV1(serviceName string) string {
	return serviceName
}

// deploymentNameV2 returns the deployment name used
// in the current implementation of the Azure provider.
func deploymentNameV2(serviceName string) string {
	return serviceName + "-v2"
}

// StartInstance is specified in the InstanceBroker interface.
func (env *azureEnviron) StartInstance(cons constraints.Value, possibleTools tools.List,
	machineConfig *cloudinit.MachineConfig) (_ instance.Instance, _ *instance.HardwareCharacteristics, err error) {

	// Declaring "err" in the function signature so that we can "defer"
	// any cleanup that needs to run during error returns.

	err = environs.FinishMachineConfig(machineConfig, env.Config(), cons)
	if err != nil {
		return nil, nil, err
	}

	// Pick envtools.  Needed for the custom data (which is what we normally
	// call userdata).
	machineConfig.Tools = possibleTools[0]
	logger.Infof("picked tools %q", machineConfig.Tools)

	// Compose userdata.
	userData, err := makeCustomData(machineConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("custom data: %v", err)
	}

	azure, err := env.getManagementAPI()
	if err != nil {
		return nil, nil, err
	}
	defer env.releaseManagementAPI(azure)

	location := env.getSnapshot().ecfg.location()
	series := possibleTools.OneSeries()
	instanceType, sourceImageName, err := env.selectInstanceTypeAndImage(cons, series, location)
	if err != nil {
		return nil, nil, err
	}

	// We use the cloud service label as a way to group instances with
	// the same affinity, so that machines can be be allocated to the
	// same availability set.
	var label string
	if machineConfig.StateServer {
		// TODO(axw) label should be blank by default, have a special value
		// if provisioning a state server, else whatever's in the Affinity.
		label = "juju-state-server"
	} else {
		// TODO(axw) 2014-03-10 #1229411
		// Choose a label based on the service name of the unit
		// that is deployed to the machine.
	}

	vhd := env.newOSDisk(sourceImageName)
	// If we're creating machine-0, we'll want to expose port 22.
	// All other machines get an auto-generated public port for SSH.
	role := env.newRole(instanceType, vhd, userData, machineConfig.StateServer)
	inst, err := env.createRole(azure.ManagementAPI, role, label)
	if err != nil {
		return nil, nil, err
	}
	// TODO(bug 1193998) - return instance hardware characteristics as well
	return inst, &instance.HardwareCharacteristics{}, nil
}

// getInstance returns an up-to-date version of the instance with the given
// name.
func (env *azureEnviron) getInstance(hostedService *gwacl.HostedService, roleName string) (instance.Instance, error) {
	if n := len(hostedService.Deployments); n != 1 {
		return nil, fmt.Errorf("expected one deployment for %q, got %d", hostedService.ServiceName, n)
	}
	deployment := &hostedService.Deployments[0]

	var instanceId instance.Id
	switch deployment.Name {
	case deploymentNameV1(hostedService.ServiceName):
		// Old style instance.
		instanceId = instance.Id(hostedService.ServiceName)
		if n := len(deployment.RoleList); n != 1 {
			return nil, fmt.Errorf("expected one role for %q, got %d", deployment.Name, n)
		}
		roleName = deployment.RoleList[0].RoleName
	case deploymentNameV2(hostedService.ServiceName):
		instanceId = instance.Id(fmt.Sprintf("%s-%s", hostedService.ServiceName, roleName))
	}

	var roleInstance *gwacl.RoleInstance
	for _, role := range deployment.RoleInstanceList {
		if role.RoleName == roleName {
			roleInstance = &role
			break
		}
	}

	instance := &azureInstance{
		environ:        env,
		hostedService:  &hostedService.HostedServiceDescriptor,
		instanceId:     instanceId,
		deploymentName: deployment.Name,
		roleName:       roleName,
		roleInstance:   roleInstance,
	}
	return instance, nil
}

// newOSDisk creates a gwacl.OSVirtualHardDisk object suitable for an
// Azure Virtual Machine.
func (env *azureEnviron) newOSDisk(sourceImageName string) *gwacl.OSVirtualHardDisk {
	vhdName := gwacl.MakeRandomDiskName("juju")
	vhdPath := fmt.Sprintf("vhds/%s", vhdName)
	snap := env.getSnapshot()
	storageAccount := snap.ecfg.storageAccountName()
	mediaLink := gwacl.CreateVirtualHardDiskMediaLink(storageAccount, vhdPath)
	// The disk label is optional and the disk name can be omitted if
	// mediaLink is provided.
	return gwacl.NewOSVirtualHardDisk("", "", "", mediaLink, sourceImageName, "Linux")
}

// getInitialEndpoints returns a slice of the endpoints every instance should have open
// (ssh port, etc).
func (env *azureEnviron) getInitialEndpoints(stateServer bool) []gwacl.InputEndpoint {
	// TODO(axw) either proxy ssh traffic through one of the
	// randomly chosen VMs to the internal address, or otherwise
	// don't load balance SSH and provie a way of getting the
	// local port.
	cfg := env.Config()
	endpoints := []gwacl.InputEndpoint{{
		LocalPort: 22,
		Name:      "sshport",
		Port:      22,
		Protocol:  "tcp",
	}}
	if stateServer {
		endpoints = append(endpoints, []gwacl.InputEndpoint{{
			LocalPort: cfg.StatePort(),
			Port:      cfg.StatePort(),
			Protocol:  "tcp",
			Name:      "stateport",
		}, {
			LocalPort: cfg.APIPort(),
			Port:      cfg.APIPort(),
			Protocol:  "tcp",
			Name:      "apiport",
		}}...)
	}
	for i, endpoint := range endpoints {
		endpoint.LoadBalancedEndpointSetName = endpoint.Name
		endpoint.LoadBalancerProbe = &gwacl.LoadBalancerProbe{
			Port:     endpoint.Port,
			Protocol: "TCP",
		}
		endpoints[i] = endpoint
	}
	return endpoints
}

// newRole creates a gwacl.Role object (an Azure Virtual Machine) which uses
// the given Virtual Hard Drive.
//
// The VM will have:
// - an 'ubuntu' user defined with an unguessable (randomly generated) password
// - its ssh port (TCP 22) open
// (if a state server)
// - its state port (TCP mongoDB) port open
// - its API port (TCP) open
//
// roleSize is the name of one of Azure's machine types, e.g. ExtraSmall,
// Large, A6 etc.
func (env *azureEnviron) newRole(roleSize string, vhd *gwacl.OSVirtualHardDisk, userData string, stateServer bool) *gwacl.Role {
	roleName := gwacl.MakeRandomRoleName("juju")
	// Create a Linux Configuration with the username and the password
	// empty and disable SSH with password authentication.
	hostname := roleName
	username := "ubuntu"
	password := gwacl.MakeRandomPassword()
	linuxConfigurationSet := gwacl.NewLinuxProvisioningConfigurationSet(hostname, username, password, userData, "true")
	// Generate a Network Configuration with the initially required ports open.
	networkConfigurationSet := gwacl.NewNetworkConfigurationSet(env.getInitialEndpoints(stateServer), nil)
	role := gwacl.NewRole(
		roleSize, roleName, vhd,
		[]gwacl.ConfigurationSet{*linuxConfigurationSet, *networkConfigurationSet},
	)
	role.AvailabilitySetName = "juju"
	return role
}

// Spawn this many goroutines to issue requests for destroying services.
// TODO: this is currently set to 1 because of a problem in Azure:
// removing Services in the same affinity group concurrently causes a conflict.
// This conflict is wrongly reported by Azure as a BadRequest (400).
// This has been reported to Windows Azure.
const maxConcurrentDeletes = 1

// StartInstance is specified in the InstanceBroker interface.
func (env *azureEnviron) StopInstances(instances []instance.Instance) error {
	context, err := env.getManagementAPI()
	if err != nil {
		return err
	}
	defer env.releaseManagementAPI(context)

	// Destroy all the roles in parallel. Record services for which
	// roles are destroyed, so we can garbage collect later.
	services := make(map[string]bool)
	run := parallel.NewRun(maxConcurrentDeletes)
	for _, instance := range instances {
		instance, ok := instance.(*azureInstance)
		if !ok {
			continue
		}
		serviceName := instance.hostedService.ServiceName
		deploymentName := instance.deploymentName
		roleName := instance.roleName
		services[serviceName] = true
		run.Do(func() error {
			return context.DeleteRole(&gwacl.DeleteRoleRequest{
				ServiceName:    serviceName,
				DeploymentName: deploymentName,
				RoleName:       roleName,
				DeleteMedia:    true,
			})
		})
	}
	if err := run.Wait(); err != nil {
		return fmt.Errorf("failed to delete roles", err)
	}

	// Destroy services now bereft of roles.
	run = parallel.NewRun(maxConcurrentDeletes)
	for serviceName := range services {
		serviceName := serviceName // copy for closure
		run.Do(func() error {
			service, err := context.GetHostedServiceProperties(serviceName, true)
			if err != nil {
				return err
			} else if len(service.Deployments) != 1 {
				return nil
			} else if len(service.Deployments[0].RoleList) != 0 {
				return nil
			}
			return context.DeleteHostedService(serviceName)
		})
	}
	return run.Wait()
}

// destroyAllServices destroys all Cloud Services and deployments contained.
// This is needed to clean up broken environments, in which there are cloud
// services with no deployments.
func (env *azureEnviron) destroyAllServices() error {
	context, err := env.getManagementAPI()
	if err != nil {
		return err
	}
	defer env.releaseManagementAPI(context)

	request := &gwacl.ListPrefixedHostedServicesRequest{ServiceNamePrefix: env.getEnvPrefix()}
	services, err := context.ListPrefixedHostedServices(request)
	if err != nil {
		return err
	}

	run := parallel.NewRun(maxConcurrentDeletes)
	for _, service := range services {
		run.Do(func() error {
			return context.DestroyHostedService(&gwacl.DestroyHostedServiceRequest{
				ServiceName: service.ServiceName,
			})
		})
	}
	return run.Wait()
}

// Instances is specified in the Environ interface.
func (env *azureEnviron) Instances(ids []instance.Id) ([]instance.Instance, error) {
	context, err := env.getManagementAPI()
	if err != nil {
		return nil, err
	}
	defer env.releaseManagementAPI(context)

	type instanceId struct {
		serviceName, roleName string
	}

	prefix := env.getEnvPrefix()
	instancesIds := make([]instanceId, len(ids))
	var serviceNames set.Strings
	for i, id := range ids {
		if !strings.HasPrefix(string(id), prefix) {
			continue
		}
		fields := strings.Split(string(id)[len(prefix):], "-")
		serviceName := prefix + fields[0]
		var roleName string
		if len(fields) > 1 {
			roleName = fields[1]
		}
		instancesIds[i] = instanceId{
			serviceName: serviceName,
			roleName:    roleName,
		}
		serviceNames.Add(serviceName)
	}

	// Map service names to gwacl.HostedServices.
	services, err := context.ListSpecificHostedServices(&gwacl.ListSpecificHostedServicesRequest{
		ServiceNames: serviceNames.Values(),
	})
	if err != nil {
		return nil, err
	}
	if len(services) == 0 {
		return nil, environs.ErrNoInstances
	}
	hostedServices := make(map[string]*gwacl.HostedService)
	for _, s := range services {
		hostedService, err := context.GetHostedServiceProperties(s.ServiceName, true)
		if err != nil {
			return nil, err
		}
		hostedServices[s.ServiceName] = hostedService
	}

	err = nil
	instances := make([]instance.Instance, len(ids))
	for i, id := range instancesIds {
		if id.serviceName == "" {
			// Previously determined to be an invalid instance ID.
			continue
		}
		hostedService := hostedServices[id.serviceName]
		instance, err := env.getInstance(hostedService, id.roleName)
		if err == nil {
			instances[i] = instance
		} else {
			logger.Debugf("failed to get instance for role %q in service %q: %v", id.roleName, hostedService.ServiceName, err)
		}
	}
	for _, instance := range instances {
		if instance == nil {
			err = environs.ErrPartialInstances
		}
	}
	return instances, err
}

// AllInstances is specified in the InstanceBroker interface.
func (env *azureEnviron) AllInstances() ([]instance.Instance, error) {
	// The instance list is built using the list of all the Azure
	// Services (instance==service).
	// Acquire management API object.
	context, err := env.getManagementAPI()
	if err != nil {
		return nil, err
	}
	defer env.releaseManagementAPI(context)

	request := &gwacl.ListPrefixedHostedServicesRequest{ServiceNamePrefix: env.getEnvPrefix()}
	serviceDescriptors, err := context.ListPrefixedHostedServices(request)
	if err != nil {
		return nil, err
	}

	var instances []instance.Instance
	for _, sd := range serviceDescriptors {
		hostedService, err := context.GetHostedServiceProperties(sd.ServiceName, true)
		if err != nil {
			return nil, err
		} else if len(hostedService.Deployments) != 1 {
			continue
		}
		deployment := &hostedService.Deployments[0]
		for _, role := range deployment.RoleList {
			instance, err := env.getInstance(hostedService, role.RoleName)
			if err != nil {
				return nil, err
			}
			instances = append(instances, instance)
		}
	}
	return instances, nil
}

// getEnvPrefix returns the prefix used to name the objects specific to this
// environment.
func (env *azureEnviron) getEnvPrefix() string {
	return fmt.Sprintf("juju-%s-", env.Name())
}

// Storage is specified in the Environ interface.
func (env *azureEnviron) Storage() storage.Storage {
	return env.getSnapshot().storage
}

// Destroy is specified in the Environ interface.
func (env *azureEnviron) Destroy() error {
	logger.Debugf("destroying environment %q", env.name)

	// Stop all instances.
	if err := env.destroyAllServices(); err != nil {
		return fmt.Errorf("cannot destroy instances: %v", err)
	}

	// Delete vnet and affinity group.
	if err := env.deleteVirtualNetwork(); err != nil {
		return fmt.Errorf("cannot delete the environment's virtual network: %v", err)
	}
	if err := env.deleteAffinityGroup(); err != nil {
		return fmt.Errorf("cannot delete the environment's affinity group: %v", err)
	}

	// Delete storage.
	// Deleting the storage is done last so that if something fails
	// half way through the Destroy() method, the storage won't be cleaned
	// up and thus an attempt to re-boostrap the environment will lead to
	// a "error: environment is already bootstrapped" error.
	if err := env.Storage().RemoveAll(); err != nil {
		return fmt.Errorf("cannot clean up storage: %v", err)
	}
	return nil
}

// OpenPorts is specified in the Environ interface. However, Azure does not
// support the global firewall mode.
func (env *azureEnviron) OpenPorts(ports []instance.Port) error {
	return nil
}

// ClosePorts is specified in the Environ interface. However, Azure does not
// support the global firewall mode.
func (env *azureEnviron) ClosePorts(ports []instance.Port) error {
	return nil
}

// Ports is specified in the Environ interface.
func (env *azureEnviron) Ports() ([]instance.Port, error) {
	// TODO: implement this.
	return []instance.Port{}, nil
}

// Provider is specified in the Environ interface.
func (env *azureEnviron) Provider() environs.EnvironProvider {
	return azureEnvironProvider{}
}

// azureManagementContext wraps two things: a gwacl.ManagementAPI (effectively
// a session on the Azure management API) and a tempCertFile, which keeps track
// of the temporary certificate file that needs to be deleted once we're done
// with this particular session.
// Since it embeds *gwacl.ManagementAPI, you can use it much as if it were a
// pointer to a ManagementAPI object.  Just don't forget to release it after
// use.
type azureManagementContext struct {
	*gwacl.ManagementAPI
	certFile *tempCertFile
}

var (
	retryPolicy = gwacl.RetryPolicy{
		NbRetries: 6,
		HttpStatusCodes: []int{
			http.StatusConflict,
			http.StatusRequestTimeout,
			http.StatusInternalServerError,
			http.StatusServiceUnavailable,
		},
		Delay: 10 * time.Second}
)

// getManagementAPI obtains a context object for interfacing with Azure's
// management API.
// For now, each invocation just returns a separate object.  This is probably
// wasteful (each context gets its own SSL connection) and may need optimizing
// later.
func (env *azureEnviron) getManagementAPI() (*azureManagementContext, error) {
	snap := env.getSnapshot()
	subscription := snap.ecfg.managementSubscriptionId()
	certData := snap.ecfg.managementCertificate()
	certFile, err := newTempCertFile([]byte(certData))
	if err != nil {
		return nil, err
	}
	// After this point, if we need to leave prematurely, we should clean
	// up that certificate file.
	location := snap.ecfg.location()
	mgtAPI, err := gwacl.NewManagementAPIWithRetryPolicy(subscription, certFile.Path(), location, retryPolicy)
	if err != nil {
		certFile.Delete()
		return nil, err
	}
	context := azureManagementContext{
		ManagementAPI: mgtAPI,
		certFile:      certFile,
	}
	return &context, nil
}

// releaseManagementAPI frees up a context object obtained through
// getManagementAPI.
func (env *azureEnviron) releaseManagementAPI(context *azureManagementContext) {
	// Be tolerant to incomplete context objects, in case we ever get
	// called during cleanup of a failed attempt to create one.
	if context == nil || context.certFile == nil {
		return
	}
	// For now, all that needs doing is to delete the temporary certificate
	// file.  We may do cleverer things later, such as connection pooling
	// where this method returns a context to the pool.
	context.certFile.Delete()
}

// updateStorageAccountKey queries the storage account key, and updates the
// version cached in env.storageAccountKey.
//
// It takes a snapshot in order to preserve transactional integrity relative
// to the snapshot's starting state, without having to lock the environment
// for the duration.  If there is a conflicting change to env relative to the
// state recorded in the snapshot, this function will fail.
func (env *azureEnviron) updateStorageAccountKey(snapshot *azureEnviron) (string, error) {
	// This method follows an RCU pattern, an optimistic technique to
	// implement atomic read-update transactions: get a consistent snapshot
	// of state; process data; enter critical section; check for conflicts;
	// write back changes.  The advantage is that there are no long-held
	// locks, in particular while waiting for the request to Azure to
	// complete.
	// "Get a consistent snapshot of state" is the caller's responsibility.
	// The caller can use env.getSnapshot().

	// Process data: get a current account key from Azure.
	key, err := env.queryStorageAccountKey()
	if err != nil {
		return "", err
	}

	// Enter critical section.
	env.Lock()
	defer env.Unlock()

	// Check for conflicts: is the config still what it was?
	if env.ecfg != snapshot.ecfg {
		// The environment has been reconfigured while we were
		// working on this, so the key we just get may not be
		// appropriate any longer.  So fail.
		// Whatever we were doing isn't likely to be right any more
		// anyway.  Otherwise, it might be worth returning the key
		// just in case it still works, and proceed without updating
		// env.storageAccountKey.
		return "", fmt.Errorf("environment was reconfigured")
	}

	// Write back changes.
	env.storageAccountKey = key
	return key, nil
}

// getStorageContext obtains a context object for interfacing with Azure's
// storage API.
// For now, each invocation just returns a separate object.  This is probably
// wasteful (each context gets its own SSL connection) and may need optimizing
// later.
func (env *azureEnviron) getStorageContext() (*gwacl.StorageContext, error) {
	snap := env.getSnapshot()
	key := snap.storageAccountKey
	if key == "" {
		// We don't know the storage-account key yet.  Request it.
		var err error
		key, err = env.updateStorageAccountKey(snap)
		if err != nil {
			return nil, err
		}
	}
	context := gwacl.StorageContext{
		Account:       snap.ecfg.storageAccountName(),
		Key:           key,
		AzureEndpoint: gwacl.GetEndpoint(snap.ecfg.location()),
		RetryPolicy:   retryPolicy,
	}
	return &context, nil
}

// baseURLs specifies an Azure specific location where we look for simplestreams information.
// It contains the central databases for the released and daily streams, but this may
// become more configurable.  This variable is here as a placeholder, but also
// as an injection point for tests.
var baseURLs = []string{}

// GetImageSources returns a list of sources which are used to search for simplestreams image metadata.
func (env *azureEnviron) GetImageSources() ([]simplestreams.DataSource, error) {
	sources := make([]simplestreams.DataSource, 1+len(baseURLs))
	sources[0] = storage.NewStorageSimpleStreamsDataSource("cloud storage", env.Storage(), storage.BaseImagesPath)
	for i, url := range baseURLs {
		sources[i+1] = simplestreams.NewURLDataSource("Azure base URL", url, simplestreams.VerifySSLHostnames)
	}
	return sources, nil
}

// GetToolsSources returns a list of sources which are used to search for simplestreams tools metadata.
func (env *azureEnviron) GetToolsSources() ([]simplestreams.DataSource, error) {
	// Add the simplestreams source off the control bucket.
	sources := []simplestreams.DataSource{
		storage.NewStorageSimpleStreamsDataSource("cloud storage", env.Storage(), storage.BaseToolsPath)}
	return sources, nil
}

// getImageMetadataSigningRequired returns whether this environment requires
// image metadata from Simplestreams to be signed.
func (env *azureEnviron) getImageMetadataSigningRequired() bool {
	// Hard-coded to true for now.  Once we support custom base URLs,
	// this may have to change.
	return true
}

// Region is specified in the HasRegion interface.
func (env *azureEnviron) Region() (simplestreams.CloudSpec, error) {
	ecfg := env.getSnapshot().ecfg
	return simplestreams.CloudSpec{
		Region:   ecfg.location(),
		Endpoint: string(gwacl.GetEndpoint(ecfg.location())),
	}, nil
}
