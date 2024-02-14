/*
Copyright Brendan Thompson

Licensed under the PolyForm Internal Use License, Version 1.0.0 (the "License");
you may not use this file except in compliance with the License.
A copy of the License may be obtained at

https://polyformproject.org/licenses/internal-use/1.0.0/
*/

package azure

import (
	"context"
	"instancescheduler/internal/patchwindow"
	"instancescheduler/internal/schedule"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	compute "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"
	"github.com/rs/zerolog/log"
)

func NewComputeClient(subscriptionID, tagsConfigPath string) (*ComputeClient, error) {
	var computeClient ComputeClient

	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}

	client, err := compute.NewVirtualMachinesClient(subscriptionID, credential, nil)
	if err != nil {
		return nil, err
	}

	tags, err := NewTagsFromConfig(tagsConfigPath)
	if err != nil {
		return nil, err
	}

	computeClient.client = client
	computeClient.ctx = context.Background()
	computeClient.Tags = tags

	return &computeClient, nil
}

type ComputeClient struct {
	Tags *Tags

	client *compute.VirtualMachinesClient
	ctx    context.Context
}

// ListInstances returns a list of all instances within an Azure subscription
func (c *ComputeClient) ListInstances() ([]*compute.VirtualMachine, error) {
	var instances []*compute.VirtualMachine
	pager := c.client.NewListAllPager(nil)

	for pager.More() {
		page, err := pager.NextPage(c.ctx)
		if err != nil {
			return nil, err
		}

		for _, instance := range page.Value {
			instances = append(instances, instance)
		}
	}

	return instances, nil
}

// AssessInstancesAndAction iterates through the `instances` passed into the method to ascertain
// if the instance should be; powered-off, powered-on, or no action
func (c *ComputeClient) AssessInstancesAndAction(instances []*compute.VirtualMachine) {
	for _, instance := range instances {
		var isWithinPatchWindow bool
		var isCurrentTimeWithinPatchWindow bool

		resourceID, err := arm.ParseResourceID(*instance.ID)
		if err != nil {
			log.Error().Stack().Err(err).Str("instance", *instance.Name).Msg("Unable to parse resource ID")
			continue
		}

		enabled, stringSchedule, stringPatchWindow := c.Tags.LoadValues(instance.Tags)

		log.Debug().Msgf("String patch window: %s", stringPatchWindow)

		if enabled {
			schedule, err := schedule.NewSchedule([]byte(stringSchedule))
			if err != nil {
				log.Error().Stack().Err(err).Msg("Unable to create a schedule based on input")
			}

			patchWindow, err := patchwindow.New([]byte(stringPatchWindow))
			if err != nil {
				log.Error().Stack().Err(err).Msg("Failed to get patch window")
			}

			nextPatchWindowStart, err := patchWindow.NextWindowStart()
			if err != nil {
				log.Error().Stack().Err(err).Msg("Failed to get the next patch window start date")
			}

			log.Debug().Msgf("Next patch window start: %s", nextPatchWindowStart.String())

			if !schedule.Validate() || !schedule.ValidateOverrides() {
				continue
			}

			isInstanceRunning := c.IsInstanceRunning(resourceID.ResourceGroupName, resourceID.Name)
			shouldShutdown := schedule.ShouldShutdown()

			if patchWindow != nil {
				isWithinPatchWindow = schedule.IsWithinPatchWindow(
					patchWindow.Timeslice.Start, patchWindow.Timeslice.End, patchWindow.IsToday(),
				)
				isCurrentTimeWithinPatchWindow = patchWindow.CurrentTimeWithinRange()
			} else {
				isWithinPatchWindow = false
			}

			c.ManagePowerState(isInstanceRunning, isWithinPatchWindow, isCurrentTimeWithinPatchWindow,
				shouldShutdown, resourceID.ResourceGroupName, resourceID.Name)
		}
	}
}

// ManagePowerState will power-off, power-on or leave an instance alone
func (c *ComputeClient) ManagePowerState(isInstanceRunning, isWithinPatchWindow,
	isCurrentTimeWithinPatchWindow, shouldShutdown bool, resourceGroupName, instanceName string) {
	if shouldShutdown && isInstanceRunning && !isWithinPatchWindow {
		c.ShutdownInstance(resourceGroupName, instanceName)
	} else if !shouldShutdown && !isInstanceRunning {
		c.StartInstance(resourceGroupName, instanceName)
	} else if isCurrentTimeWithinPatchWindow && !isInstanceRunning {
		log.Info().Str("instance", instanceName).Msg("Instance is starting as its within 1 hour of the patch window.")
		c.StartInstance(resourceGroupName, instanceName)
	} else {
		log.Info().Str("instance", instanceName).Msg("No action required")
	}
}

// ShutdownInstance will shutdown a given instance
func (c *ComputeClient) ShutdownInstance(resourceGroupName string, instanceName string) {
	opts := &compute.VirtualMachinesClientBeginPowerOffOptions{
		SkipShutdown: to.Ptr(false),
	}

	log.Info().Str("instance", instanceName).Msg("Shutting down instance")

	poller, err := c.client.BeginPowerOff(c.ctx, resourceGroupName, instanceName, opts)
	if err != nil {
		log.Error().Stack().Err(err).Str("instance", instanceName).Msg("Failed to execute power off")
		return
	}

	_, err = poller.PollUntilDone(c.ctx, nil)
	if err != nil {
		log.Error().Stack().Err(err).Str("instance", instanceName).Msg("Failed to complete power off")
		return
	}

	log.Info().Str("instance", instanceName).Msg("Shutting down successful")
}

// StartInstance will power-on a given instance
func (c *ComputeClient) StartInstance(resourceGroupName string, instanceName string) {
	opts := &compute.VirtualMachinesClientBeginStartOptions{}

	log.Info().Str("instance", instanceName).Msg("Starting up instance")

	poller, err := c.client.BeginStart(c.ctx, resourceGroupName, instanceName, opts)
	if err != nil {
		log.Error().Stack().Err(err).Str("instance", instanceName).Msg("Failed to execute startup")
		return
	}

	_, err = poller.PollUntilDone(c.ctx, nil)
	if err != nil {
		log.Error().Stack().Err(err).Str("instance", instanceName).Msg("Failed to complete startup")
		return
	}

	log.Info().Str("instance", instanceName).Msg("Startup successful")
}

// IsInstanceRunning determines if a given instance is in a running state
//
// parameters:
//   - `resourceGroupName` – the name of the resource group
//   - `instanceName` - name of the instance in Azure
//
// returns:
//   - `bool`
func (c *ComputeClient) IsInstanceRunning(resourceGroupName string, instanceName string) bool {
	instance, err := c.client.InstanceView(c.ctx, resourceGroupName, instanceName, nil)
	if err != nil {

	}

	for _, status := range instance.Statuses {
		if strings.Contains(*status.Code, "PowerState/") {
			powerState := ParsePowerState(*status.Code)

			log.Debug().Str("parsed", powerState.String()).Str("raw", *status.Code).Msg("Instance Power State")

			if powerState == PowerStateRunning || powerState == PowerStateStarting {
				return true
			}
		}
	}

	return false
}
