package profiler

//go:generate mockgen -source vm_statistics_printer_test.go -destination vm_statistics_printer_mocks_test.go -package profiler

import (
	"testing"

	"github.com/Fantom-foundation/Aida/executor"
	"github.com/Fantom-foundation/Aida/utils"
	"github.com/Fantom-foundation/Tosca/go/vm/registry"
	"go.uber.org/mock/gomock"
)

func TestVirtualMachineStatisticsPrinter_WorksWithDefaultSetup(t *testing.T) {
	cfg := utils.Config{}
	ext := MakeVirtualMachineStatisticsPrinter[any](&cfg)
	ext.PostRun(executor.State[any]{}, nil, nil)
}

func TestVirtualMachineStatisticsPrinter_TriggersStatPrintingAtEndOfRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	vm := NewMockProfilingVm(ctrl)
	registry.RegisterVirtualMachine("test-vm", vm)

	vm.EXPECT().DumpProfile()

	cfg := utils.Config{}
	cfg.VmImpl = "test-vm"
	ext := MakeVirtualMachineStatisticsPrinter[any](&cfg)

	ext.PostRun(executor.State[any]{}, nil, nil)
}

type ProfilingVm interface {
	registry.VirtualMachine
	registry.ProfilingVM
}