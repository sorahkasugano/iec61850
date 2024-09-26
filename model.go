package iec61850

// #include <iec61850_server.h>
import "C"
import (
	"os"
	"unsafe"
)

type IedModel struct {
	model *C.IedModel
}

func NewIedModel(name string) *IedModel {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return &IedModel{
		model: C.IedModel_create(cname),
	}
}

func (m *IedModel) Destroy() {
	C.IedModel_destroy(m.model)
}

func CreateModelFromConfigFileEx(filepath string) (*IedModel, error) {
	if _, err := os.Stat(filepath); err != nil {
		if os.IsNotExist(err) {
			return nil, err
		}
	}
	cFilepath := C.CString(filepath)
	// 释放内存
	defer C.free(unsafe.Pointer(cFilepath))
	model := &IedModel{
		model: C.ConfigFileParser_createModelFromConfigFileEx(cFilepath),
	}
	return model, nil
}

type LogicalDevice struct {
	device *C.LogicalDevice
}

func (m *IedModel) CreateLogicalDevice(name string) *LogicalDevice {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return &LogicalDevice{
		device: C.LogicalDevice_create(cname, m.model),
	}
}

// GetModelNodeByObjectReference
func (m *IedModel) GetModelNodeByObjectReference(objectReference string) *DataAttribute {
	return  &DataAttribute{
		attribute: (*C.DataAttribute)(unsafe.Pointer(C.IedModel_getModelNodeByObjectReference(m.model, C.CString(objectReference)))),
	}
}

// GetModelNodeByShortAddress
func (m *IedModel) GetModelNodeByShortAddress(addr uint32) *DataAttribute {
	return  &DataAttribute{
		attribute: (*C.DataAttribute)(unsafe.Pointer(C.IedModel_getModelNodeByShortAddress(m.model, C.uint32_t(addr)))),
	}
}

type LogicalNode struct {
	node *C.LogicalNode
}

func (d *LogicalDevice) CreateLogicalNode(name string) *LogicalNode {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return &LogicalNode{
		node: C.LogicalNode_create(cname, d.device),
	}
}

type DataObject struct {
	object *C.DataObject
}

// ENS: EnumerationString
// VSS: Visible String Setting
// SAV: Sampled Value
// APC: Analogue Process Control

func (n *LogicalNode) CreateDataObjectCDC_ENS(name string) *DataObject {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return &DataObject{
		object: C.CDC_ENS_create(cname, (*C.ModelNode)(n.node), 0),
	}
}

func (n *LogicalNode) CreateDataObjectCDC_VSS(name string) *DataObject {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return &DataObject{
		object: C.CDC_VSS_create(cname, (*C.ModelNode)(n.node), 0),
	}
}

func (n *LogicalNode) CreateDataObjectCDC_SAV(name string, isInteger bool) *DataObject {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return &DataObject{
		object: C.CDC_SAV_create(cname, (*C.ModelNode)(n.node), 0, C.bool(isInteger)),
	}
}

func (n *LogicalNode) CreateDataObjectCDC_APC(name string, ctlModel int) *DataObject {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return &DataObject{
		object: C.CDC_APC_create(cname, (*C.ModelNode)(n.node), 0, C.uint(ctlModel), C.bool(false)),
	}
}

type DataAttribute struct {
	attribute *C.DataAttribute
}

func (do *DataObject) GetChild(name string) *DataAttribute {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return &DataAttribute{
		attribute: (*C.DataAttribute)(unsafe.Pointer(C.ModelNode_getChild((*C.ModelNode)(unsafe.Pointer(do.object)), cname))),
	}
}

type DataSet struct {
	dataSet *C.DataSet
}

// CreateDataSet creates a new DataSet under this LogicalNode.
func (ln *LogicalNode) CreateDataSet(name string) *DataSet {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	cDataSet := C.DataSet_create(cName, ln.node)
	return &DataSet{dataSet: cDataSet}
}

// AddDataSetEntry adds a new DataSetEntry to this DataSet.
func (ds *DataSet) AddDataSetEntry(ref string) {
	cRef := C.CString(ref)
	defer C.free(unsafe.Pointer(cRef))

	C.DataSetEntry_create(ds.dataSet, cRef, -1, nil)
}
