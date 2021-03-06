package nvml

// See https://docs.nvidia.com/deploy/nvml-api/group__nvmlDeviceQueries.html

/*
#cgo CPPFLAGS: -I/usr/include/nvidia-367/ -I/usr/include/nvidia-375/ -I/usr/include/nvidia-378/ -I/usr/include/nvidia-381/ -I/usr/include/nvidia-384/
#cgo LDFLAGS: -l nvidia-ml -L/usr/lib/nvidia-367/ -L/usr/lib/nvidia-375/ -L/usr/lib/nvidia-378/ -L/usr/lib/nvidia-381/ -L/usr/lib/nvidia-384/

#include "nvmlbridge.h"
*/
import "C"

import (
	"errors"
	"log"
	"unsafe"
)

type Device struct {
	nvmldevice C.nvmlDevice_t
	index      uint
	pcibus     string
	name       string
	uuid       string
}

// NewDevice is a contstructor function for Device structs. Given an nvmlDevice_t
// object as input, it populates some static property fields and returns a Device
func NewDevice(cdevice C.nvmlDevice_t) (*Device, error) {
	device := Device{
		nvmldevice: cdevice,
	}

	uuid, err := device.UUID()
	if err != nil {
		return nil, errors.New("Cannot retrieve UUID property")
	}
	device.uuid = uuid

	name, err := device.Name()
	if err != nil {
		return nil, errors.New("Cannot retrieve Name property")
	}
	device.name = name

	index, err := device.Index()
	if err != nil {
		return nil, errors.New("Cannot retrieve Index property")
	}
	device.index = index

	return &device, nil
}

func (gpu *Device) PowerState() (int, error) {
	var pstate C.nvmlPstates_t
	var result C.nvmlReturn_t

	result = C.nvmlDeviceGetPowerState(gpu.nvmldevice, &pstate)
	if result != C.NVML_SUCCESS {
		return -1, errors.New("GetPowerState returned error")
	}

	return int(pstate), nil
}

// Temp returns the current temperature of the card in degrees Celsius
func (gpu *Device) Temp() (uint, error) {
	var result C.nvmlReturn_t
	var ctemp C.uint

	result = C.nvmlDeviceGetTemperature(gpu.nvmldevice, C.NVML_TEMPERATURE_GPU, &ctemp)
	if result != C.NVML_SUCCESS {
		return 0, errors.New("GetPowerState returned error")
	}

	return uint(ctemp), nil
}

type cIntPropFunc struct {
	f C.getintProperty
}

var intpropfunctions = map[string]*cIntPropFunc{
	"Index":                        {C.getintProperty(C.nvmlDeviceGetIndex)},
	"MinorNumber":                  {C.getintProperty(C.nvmlDeviceGetMinorNumber)},
	"InforomConfigurationChecksum": {C.getintProperty(C.nvmlDeviceGetInforomConfigurationChecksum)},
	"MaxPCIeLinkGeneration":        {C.getintProperty(C.nvmlDeviceGetMaxPcieLinkGeneration)},
	"MaxPCIeLinkWidth":             {C.getintProperty(C.nvmlDeviceGetMaxPcieLinkWidth)},
	"CurrPCIeLinkGeneration":       {C.getintProperty(C.nvmlDeviceGetCurrPcieLinkGeneration)},
	"CurrPCIeLinkWidth":            {C.getintProperty(C.nvmlDeviceGetCurrPcieLinkWidth)},
	"PCIeReplayCounter":            {C.getintProperty(C.nvmlDeviceGetPcieReplayCounter)},
	"FanSpeed":                     {C.getintProperty(C.nvmlDeviceGetFanSpeed)},
	"PowerManagementLimit":         {C.getintProperty(C.nvmlDeviceGetPowerManagementLimit)},
	"PowerManagementDefaultLimit":  {C.getintProperty(C.nvmlDeviceGetPowerManagementDefaultLimit)},
	"PowerUsage":                   {C.getintProperty(C.nvmlDeviceGetPowerUsage)},
	"EnforcedPowerLimit":           {C.getintProperty(C.nvmlDeviceGetEnforcedPowerLimit)},
	"BoardId":                      {C.getintProperty(C.nvmlDeviceGetBoardId)},
	"MultiGpuBoard":                {C.getintProperty(C.nvmlDeviceGetMultiGpuBoard)},
}

func (gpu *Device) intProperty(property string) (uint, error) {
	var cuintproperty C.uint

	ipf, ok := intpropfunctions[property]
	if ok == false {
		return 0, errors.New("property not found")
	}

	result := C.bridge_get_int_property(ipf.f, gpu.nvmldevice, &cuintproperty)
	if result != C.EXIT_SUCCESS {
		return 0, errors.New("getintProperty bridge returned error")
	}

	return uint(cuintproperty), nil
}

// Index returns the NVML index of the device.
func (gpu *Device) Index() (uint, error) {
	return gpu.intProperty("Index")
}

// MinorNumber returns the minor number of the device. The minor number
// is the integer such that the device node file for the GPU will be
// /dev/nvidia[Device.MinorNumber]
func (gpu *Device) MinorNumber() (uint, error) {
	return gpu.intProperty("MinorNumber")
}

// InforomConfigurationChecksum returns the checksum of the configuration
// stored in the device's inforom. (Can be used to verify identical configuration
// between devices.)
func (gpu *Device) InforomConfigurationChecksum() (uint, error) {
	return gpu.intProperty("InforomConfigurationChecksum")
}

// MaxPCIeLinkGeneration returns the maximum PCIe link generation possible with this
// device and system.
func (gpu *Device) MaxPCIeLinkGeneration() (uint, error) {
	return gpu.intProperty("MaxPCIeLinkGeneration")
}

// MaxPCIeLinkWidth returns the maximum PCIe link width possible with this device
// and system
func (gpu *Device) MaxPCIeLinkWidth() (uint, error) {
	return gpu.intProperty("MaxPCIeLinkWidth")
}

// CurrPCIeLinkGeneration returns the current PCIe link generation number
func (gpu *Device) CurrPCIeLinkGeneration() (uint, error) {
	return gpu.intProperty("CurrPCIeLinkGeneration")
}

// CurrPCIeLinkWidth returns the current PCIe link width
func (gpu *Device) CurrPCIeLinkWidth() (uint, error) {
	return gpu.intProperty("CurrPCIeLinkWidth")
}

// PCIeReplayCounter returns the replay counter and rollover info.
func (gpu *Device) PCIeReplayCounter() (uint, error) {
	return gpu.intProperty("PCIeReplayCounter")
}

// FanSpeed returns the current fan speed of the device, on devices that
// have fans.
func (gpu *Device) FanSpeed() (uint, error) {
	return gpu.intProperty("FanSpeed")
}

// PowerManagementLimit returns the power management limit for the device, in mW
func (gpu *Device) PowerManagementLimit() (uint, error) {
	return gpu.intProperty("PowerManagementLimit")
}

// PowerManagementDefaultLimit returns the upper limit for the amount of power
// the card is allowed to draw, in mW.
func (gpu *Device) PowerManagementDefaultLimit() (uint, error) {
	return gpu.intProperty("PowerManagementDefaultLimit")
}

// PowerUsage returns the current power usage of the device, in mW.
func (gpu *Device) PowerUsage() (uint, error) {
	return gpu.intProperty("PowerUsage")
}

// EnforcedPowerLimit returns the effective power limit that the driver enforces after
// taking into account all limiters.
func (gpu *Device) EnforcedPowerLimit() (uint, error) {
	return gpu.intProperty("EnforcedPowerLimit")
}

// BoardID returns the device boardId, which will be identical for GPUs connected to
// the same PLX
func (gpu *Device) BoardId() (uint, error) {
	return gpu.intProperty("BoardId")
}

// GetDecoderUtilization retrieves the current utilization and sampling size in
// microseconds for the Decoder
func (gpu *Device) GetDecoderUtilization() (utilization uint, samplingPeriosUs uint, err error) {
	var result C.nvmlReturn_t
	var ctemp C.uint
	var ctemp2 C.uint

	result = C.nvmlDeviceGetDecoderUtilization(gpu.nvmldevice, &ctemp, &ctemp2)
	if result != C.NVML_SUCCESS {
		return 0, 0, errors.New("GetDecoderUtilization returned error")
	}

	return uint(ctemp), uint(ctemp2), nil
}

// GetEncoderUtilization retrieves the current utilization and sampling size in microseconds
// for the Encoder
func (gpu *Device) GetEncoderUtilization() (utilization uint, samplingPeriosUs uint, err error) {
	var result C.nvmlReturn_t
	var ctemp C.uint
	var ctemp2 C.uint

	result = C.nvmlDeviceGetEncoderUtilization(gpu.nvmldevice, &ctemp, &ctemp2)
	if result != C.NVML_SUCCESS {
		return 0, 0, errors.New("GetEncoderUtilization returned error")
	}

	return uint(ctemp), uint(ctemp2), nil
}

// GetUtilizationRates retrieves the current utilization rates for the device's major subsystems.
func (gpu *Device) GetUtilizationRates() (gpuUtilization uint, memoryUtilization uint, err error) {
	var result C.nvmlReturn_t
	var ctemp C.nvmlUtilization_t

	result = C.nvmlDeviceGetUtilizationRates(gpu.nvmldevice, &ctemp)
	if result != C.NVML_SUCCESS {
		return 0, 0, errors.New("GetUtilizationRates returned error")
	}

	return uint(ctemp.gpu), uint(ctemp.memory), nil
}

// MultiGpuBoard
func (gpu *Device) MultiGpuBoard() (bool, error) {
	p, err := gpu.intProperty("MultiGpuBoard")
	if int(p) == 0 {
		return true, err
	} else {
		return false, err
	}
}

type cTextPropFunc struct {
	f      C.gettextProperty
	length C.uint
}

var textpropfunctions = map[string]*cTextPropFunc{
	"Name":                {C.gettextProperty(C.nvmlDeviceGetName), C.NVML_DEVICE_NAME_BUFFER_SIZE},
	"Serial":              {C.gettextProperty(C.nvmlDeviceGetSerial), C.NVML_DEVICE_SERIAL_BUFFER_SIZE},
	"UUID":                {C.gettextProperty(C.nvmlDeviceGetUUID), C.NVML_DEVICE_UUID_BUFFER_SIZE},
	"InforomImageVersion": {C.gettextProperty(C.nvmlDeviceGetInforomImageVersion), C.NVML_DEVICE_INFOROM_VERSION_BUFFER_SIZE},
	"VbiosVersion":        {C.gettextProperty(C.nvmlDeviceGetVbiosVersion), C.NVML_DEVICE_VBIOS_VERSION_BUFFER_SIZE},
}

// textProperty takes a propertyname as input and then runs the corresponding
// function in the textpropfunctions map, returning the result as a Go string.
//
// textProperty takes care of allocating (and freeing) the text buffers of
// proper size.
func (gpu *Device) textProperty(property string) (string, error) {
	var propvalue string

	// If there isn't a valid entry for this property in the map, there's no reason
	// to process any further
	tpf, ok := textpropfunctions[property]
	if ok == false {
		return "", errors.New("property not found")
	}

	var buf *C.char = genCStringBuffer(uint(tpf.length))
	defer C.free(unsafe.Pointer(buf))

	result := C.bridge_get_text_property(tpf.f, gpu.nvmldevice, buf, tpf.length)
	if result != C.EXIT_SUCCESS {
		return propvalue, errors.New("gettextProperty bridge returned error")
	}

	propvalue = strndup(buf, uint(tpf.length))
	if len(propvalue) > 0 {
		return propvalue, nil
	} else {
		return "", errors.New("textProperty returned empty string")
	}
}

// InforomImageVersion returns the global inforom image version
func (gpu *Device) InforomImageVersion() (string, error) {
	return gpu.textProperty("InforomImageVersion")
}

// VbiosVersion returns the VBIOS version of the device
func (gpu *Device) VbiosVersion() (string, error) {
	return gpu.textProperty("VbiosVersion")
}

// Return the product name of the device, e.g. "Tesla K40m"
func (gpu *Device) Name() (string, error) {
	return gpu.textProperty("Name")
}

// Return the UUID of the device
func (gpu *Device) UUID() (string, error) {
	return gpu.textProperty("UUID")
}

// Return the serial number of the device
func (gpu *Device) Serial() (string, error) {
	return gpu.textProperty("Serial")
}

// Go correspondent of the C.nvmlMemory_t struct. Memory in bytes
type NVMLMemory struct {
	Free  uint64
	Total uint64
	Used  uint64
}

// MemoryInfo returns a NVMLMemory struct populated with the amount of memory used,
// free, and in total on the device, in bytes.
func (gpu *Device) MemoryInfo() (NVMLMemory, error) {
	var result C.nvmlReturn_t
	var cmeminfo C.nvmlMemory_t
	var meminfo NVMLMemory

	result = C.nvmlDeviceGetMemoryInfo(gpu.nvmldevice, &cmeminfo)
	if result != C.NVML_SUCCESS {
		return meminfo, errors.New("GetPowerState returned error")
	}

	meminfo.Free = uint64(cmeminfo.free)
	meminfo.Total = uint64(cmeminfo.total)
	meminfo.Used = uint64(cmeminfo.used)

	return meminfo, nil
}

// Return a proper golang error of representation of the nvmlReturn_t error
func (gpu *Device) Error(cerror C.nvmlReturn_t) error {
	var cerrorstring *C.char

	// No need to process anything further if the nvml call succeeded
	if cerror == C.NVML_SUCCESS {
		return nil
	}

	cerrorstring = C.nvmlErrorString(cerror)
	if cerrorstring == nil {
		// I'm not sure how this could happen, but it's easy to check for
		return errors.New("Error not found in nvml.h")
	}

	return errors.New(C.GoString(cerrorstring))
}

func nvmlDeviceGetCount() (int, error) {
	var count C.uint

	result := C.nvmlDeviceGetCount(&count)
	if result != C.NVML_SUCCESS {
		return -1, errors.New("nvmlDeviceGetCount failed")
	}

	return int(count), nil
}

// GetAllGPUs will return a slice of type Device for all NVML devices present on
// the host system
func GetAllGPUs() ([]Device, error) {
	var devices []Device
	cdevices, err := getAllDevices()
	if err != nil {
		return devices, err
	}

	for _, cdevice := range cdevices {
		device, err := NewDevice(cdevice)
		if err != nil {
			break
		}

		devices = append(devices, *device)
	}

	return devices, nil
}

// getAllDevices returns an array of nvmlDevice_t structs representing all GPU
// devices in the system.
func getAllDevices() ([]C.nvmlDevice_t, error) {
	var devices []C.nvmlDevice_t

	device_count, err := nvmlDeviceGetCount()
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < device_count; i++ {
		var device C.nvmlDevice_t
		result := C.nvmlDeviceGetHandleByIndex(C.uint(i), &device)
		if result != C.NVML_SUCCESS {
			return devices, errors.New("nvmlDeviceGetHandleByIndex returns error")
		}

		devices = append(devices, device)
	}

	if len(devices) == 0 {
		return devices, errors.New("No devices found")
	}

	return devices, nil
}
